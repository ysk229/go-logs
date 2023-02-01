package log

import (
	"io"
	"os"
	"sync"

	"github.com/mattn/go-colorable"
	json "github.com/neilotoole/jsoncolor"
)

// Writer provides colorable Writer to the console
type Writer struct {
	out   io.Writer
	mutex sync.Mutex
	enc   *json.Encoder
}

// NewJSONColorable returns new instance of Writer which handles escape sequence from File.
func NewJSONColorable() io.Writer {
	out := colorable.NewColorable(os.Stdout) // needed for Windows
	f := &Writer{out: out}
	f.enc = json.NewEncoder(out)
	clrs := json.DefaultColors()
	f.enc.SetEscapeHTML(false)
	f.enc.SetIndent("", "  ")
	f.enc.SetColors(clrs)
	f.enc.SetTrustRawMessage(true)
	return f
}

// Write writes data on console
func (w *Writer) Write(data []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	var v interface{}
	_ = json.Unmarshal(data, &v)
	err = w.enc.Encode(v)
	return len(data), err
}
