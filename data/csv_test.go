package data

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestCSVLoader(t *testing.T) {
	text := "a;b\n1;'x;y'\n"
	l := NewCSVLoader(strings.NewReader(text), ';', '\'')

	tests := []struct {
		value string
		path  []string
	}{
		{"a", []string{"1", "1"}},
		{"b", []string{"1", "2"}},
		{"1", []string{"2", "1"}},
		{"x;y", []string{"2", "2"}},
	}

	for i, tt := range tests {
		entry, err := l.Next()
		if err != nil {
			t.Fatalf("unexpected error on entry %d: %v", i, err)
		}
		if entry.Value != tt.value {
			t.Fatalf("entry %d value = %q, want %q", i, entry.Value, tt.value)
		}
		if !reflect.DeepEqual(entry.Path, tt.path) {
			t.Fatalf("entry %d path = %v, want %v", i, entry.Path, tt.path)
		}
	}

	if _, err := l.Next(); err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}
