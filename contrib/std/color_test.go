package std

import (
	"testing"
)

func colorFunc(color, s string) string {
	result := ""
	switch color {
	case "green":
		result = SGreen(s)
	case "white":
		result = SWhite(s)
	case "yellow":
		result = SYellow(s)
	case "red":
		result = SRed(s)
	case "blue":
		result = SBlue(s)
	case "magenta":
		result = SMagenta(s)
	case "cyan":
		result = SCyan(s)
	}
	return result
}

func TestColor(t *testing.T) {
	tests := []struct {
		s        string
		color    string
		expected string
	}{
		{"g", "green", "\x1b[32mg\x1b[0m"},
		{"w", "white", "\x1b[37mw\x1b[0m"},
		{"y", "yellow", "\x1b[33my\x1b[0m"},
		{"r", "red", "\x1b[31mr\x1b[0m"},
		{"b", "blue", "\x1b[34mb\x1b[0m"},
		{"m", "magenta", "\x1b[35mm\x1b[0m"},
		{"c", "cyan", "\x1b[36mc\x1b[0m"},
	}
	for _, row := range tests {
		actual := colorFunc(row.color, row.s)
		if actual != row.expected {
			t.Errorf("color: %s, actual: %s, expected: %s", row.color, actual, row.expected)
		}
	}
}
