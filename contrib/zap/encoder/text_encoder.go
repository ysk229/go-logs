package encoder

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/mgutz/ansi"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var (
	defaultColorScheme = &ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "black+h",
	}
	noColorsColorScheme = &compiledColorScheme{
		InfoLevelColor:  ansi.ColorFunc(""),
		WarnLevelColor:  ansi.ColorFunc(""),
		ErrorLevelColor: ansi.ColorFunc(""),
		FatalLevelColor: ansi.ColorFunc(""),
		PanicLevelColor: ansi.ColorFunc(""),
		DebugLevelColor: ansi.ColorFunc(""),
		PrefixColor:     ansi.ColorFunc(""),
		TimestampColor:  ansi.ColorFunc(""),
	}
	defaultCompiledColorScheme = compileColorScheme(defaultColorScheme)
)

type ColorScheme struct {
	InfoLevelStyle  string
	WarnLevelStyle  string
	ErrorLevelStyle string
	FatalLevelStyle string
	PanicLevelStyle string
	DebugLevelStyle string
	PrefixStyle     string
	TimestampStyle  string
}

type compiledColorScheme struct {
	InfoLevelColor  func(string) string
	WarnLevelColor  func(string) string
	ErrorLevelColor func(string) string
	FatalLevelColor func(string) string
	PanicLevelColor func(string) string
	DebugLevelColor func(string) string
	PrefixColor     func(string) string
	TimestampColor  func(string) string
}

func getCompiledColor(main string, fallback string) func(string) string {
	var style string
	if main != "" {
		style = main
	} else {
		style = fallback
	}
	return ansi.ColorFunc(style)
}

func compileColorScheme(s *ColorScheme) *compiledColorScheme {
	return &compiledColorScheme{
		InfoLevelColor:  getCompiledColor(s.InfoLevelStyle, defaultColorScheme.InfoLevelStyle),
		WarnLevelColor:  getCompiledColor(s.WarnLevelStyle, defaultColorScheme.WarnLevelStyle),
		ErrorLevelColor: getCompiledColor(s.ErrorLevelStyle, defaultColorScheme.ErrorLevelStyle),
		FatalLevelColor: getCompiledColor(s.FatalLevelStyle, defaultColorScheme.FatalLevelStyle),
		PanicLevelColor: getCompiledColor(s.PanicLevelStyle, defaultColorScheme.PanicLevelStyle),
		DebugLevelColor: getCompiledColor(s.DebugLevelStyle, defaultColorScheme.DebugLevelStyle),
		PrefixColor:     getCompiledColor(s.PrefixStyle, defaultColorScheme.PrefixStyle),
		TimestampColor:  getCompiledColor(s.TimestampStyle, defaultColorScheme.TimestampStyle),
	}
}

// For Text-escaping; see textEncoder.safeAddString below.
const _hex = "0123456789abcdef"

var _textPool = sync.Pool{New: func() interface{} {
	return &textEncoder{}
}}

func getTEXTEncoder() *textEncoder {
	return _textPool.Get().(*textEncoder)
}

func putTEXTEncoder(enc *textEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.spaced = false
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	_textPool.Put(enc)
}

type textEncoder struct {
	*zapcore.EncoderConfig
	buf            *buffer.Buffer
	spaced         bool // include spaces after colons and commas
	openNamespaces int

	// for encoding generic values by reflection
	reflectBuf *buffer.Buffer
	reflectEnc *json.Encoder
	// Color scheme to use.
	colorScheme *compiledColorScheme
	isTerminal  bool
}

func NewTextNoColorEncoder() zapcore.Encoder {
	return NewTextEncoder(false)
}

func NewTextColorEncoder() zapcore.Encoder {
	return NewTextEncoder(true)
}

// NewTextEncoder creates a fast, low-allocation TEXT encoder. The encoder
// appropriately escapes all field keys and values.
//
// Note that the encoder doesn't deduplicate keys, so it's possible to produce
// a message like
// {"foo":"bar","foo":"baz"}
// This is permitted by the TEXT specification, but not encouraged. Many
// libraries will ignore duplicate key-value pairs (typically keeping the last
// pair) when unmarshaling, but users should attempt to avoid adding duplicate
// keys.
func NewTextEncoder(isTerminal bool) zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stack",
	}
	return newTextEncoder(cfg, isTerminal)
}

func newTextEncoder(cfg zapcore.EncoderConfig, isTerminal bool) *textEncoder {
	t := &textEncoder{
		EncoderConfig: &cfg,
		buf:           buffer.NewPool().Get(),
		spaced:        false,
	}
	t.SetColorScheme(&ColorScheme{
		WarnLevelStyle:  "yellow",
		InfoLevelStyle:  "38",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "41",
		PanicLevelStyle: "41",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "37",
	})
	t.isTerminal = isTerminal
	return t
}

func (enc *textEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *textEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *textEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *textEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *textEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *textEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *textEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *textEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *textEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *textEncoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = buffer.NewPool().Get()
		enc.reflectEnc = json.NewEncoder(enc.reflectBuf)

		// For consistency with our custom TEXT encoder.
		enc.reflectEnc.SetEscapeHTML(false)
	} else {
		enc.reflectBuf.Reset()
	}
}

var nullLiteralBytes = []byte("null")

// Only invoke the standard TEXT encoder if there is actually something to
// encode; otherwise write TEXT null literal directly.
func (enc *textEncoder) encodeReflected(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nullLiteralBytes, nil
	}
	enc.resetReflectBuf()
	if err := enc.reflectEnc.Encode(obj); err != nil {
		return nil, err
	}
	enc.reflectBuf.TrimNewline()
	return enc.reflectBuf.Bytes(), nil
}

func (enc *textEncoder) AddReflected(key string, obj interface{}) error {
	valueBytes, err := enc.encodeReflected(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *textEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

func (enc *textEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *textEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *textEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *textEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *textEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *textEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *textEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

func (enc *textEncoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val)) //nolint:unconvert
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *textEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	if e := enc.EncodeDuration; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// TEXT valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *textEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *textEncoder) AppendReflected(val interface{}) error {
	valueBytes, err := enc.encodeReflected(val)
	if err != nil {
		return err
	}
	enc.addElementSeparator()
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *textEncoder) AppendString(val string) {
	enc.safeAddString(val)
}

func (enc *textEncoder) AppendTimeLayout(time time.Time, layout string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.buf.AppendTime(time, layout)
	enc.buf.AppendByte('"')
}

func (enc *textEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	if e := enc.EncodeTime; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output TEXT valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *textEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *textEncoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *textEncoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *textEncoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *textEncoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *textEncoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *textEncoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *textEncoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *textEncoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *textEncoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *textEncoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *textEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	_, _ = clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *textEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()
	colorScheme := noColorsColorScheme
	if enc.isTerminal {
		colorScheme = defaultCompiledColorScheme
		if enc.colorScheme != nil {
			colorScheme = enc.colorScheme
		}
	}
	var levelColor func(string) string
	var levelText string
	switch ent.Level {
	case zapcore.InfoLevel:
		levelColor = colorScheme.InfoLevelColor
	case zapcore.WarnLevel:
		levelColor = colorScheme.WarnLevelColor
	case zapcore.ErrorLevel:
		levelColor = colorScheme.ErrorLevelColor
	case zapcore.FatalLevel:
		levelColor = colorScheme.FatalLevelColor
	case zapcore.PanicLevel:
		levelColor = colorScheme.PanicLevelColor
	default:
		levelColor = colorScheme.DebugLevelColor
	}
	levelText = "warn"
	if ent.Level != zapcore.WarnLevel {
		levelText = ent.Level.String()
	}

	levelText = strings.ToUpper(levelText)
	level := levelColor(fmt.Sprintf("%6s%3s", levelText, ""))

	if final.TimeKey != "" {
		timestamp := fmt.Sprintf("[%s]", ent.Time.Format("2006-01-02.15:04:05.000000"))
		final.buf.AppendString(colorScheme.TimestampColor(timestamp))
	}
	if final.LevelKey != "" {
		final.buf.AppendString(level)
	}
	prefix := ""
	if len(ent.LoggerName) > 0 {
		prefix = colorScheme.PrefixColor(" " + ent.LoggerName + ":")
		final.buf.AppendString(prefix)
	} else {
		prefixValue, trimmedMsg := extractPrefix(ent.Message)
		if len(prefixValue) > 0 {
			prefix = colorScheme.PrefixColor(" " + prefixValue + ":")
			final.buf.AppendString(fmt.Sprintf("%s%s", prefix, trimmedMsg))
		}
	}

	if final.MessageKey != "" {
		final.buf.AppendString(ent.Message)
	}
	for _, k := range fields {
		final.AddString(levelColor(k.Key), k.String)
	}
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.addKey(levelColor(final.StacktraceKey))
		final.buf.AppendString(ent.Stack)
	}
	if final.LineEnding != "" {
		final.buf.AppendString(final.LineEnding)
	} else {
		final.buf.AppendString(zapcore.DefaultLineEnding)
	}
	ret := final.buf
	putTEXTEncoder(final)
	return ret, nil
}

func (enc *textEncoder) clone() *textEncoder {
	clone := getTEXTEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = buffer.NewPool().Get()
	return clone
}

func (enc *textEncoder) Truncate() {
	enc.buf.Reset()
}

func (enc *textEncoder) addKey(key string) {
	enc.addElementSeparator()
	enc.buf.AppendString(key)
	enc.buf.AppendByte('=')
}

func (enc *textEncoder) addElementSeparator() {
	enc.buf.AppendByte(' ')
}

func (enc *textEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString TEXT-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or TEXTP-related problems.
func (enc *textEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *textEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		_, _ = enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *textEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *textEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func (enc *textEncoder) SetColorScheme(colorScheme *ColorScheme) {
	enc.colorScheme = compileColorScheme(colorScheme)
}

func extractPrefix(msg string) (string, string) {
	prefix := ""
	regex := regexp.MustCompile(`^\\[(.*?)\\]`)
	if regex.MatchString(msg) {
		match := regex.FindString(msg)
		prefix, msg = match[1:len(match)-1], strings.TrimSpace(msg[len(match):])
	}
	return prefix, msg
}
