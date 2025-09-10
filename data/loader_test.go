package data

import (
	"io"
	"strings"
	"testing"
)

func TestLoaderTokenization(t *testing.T) {
	text := "hello world\nnext line"
	l := NewLoader(strings.NewReader(text), nil)

	tests := []struct {
		value string
		path  []string
	}{
		{"hello", []string{"1", "0", "1"}},
		{"world", []string{"2", "6", "1"}},
		{"next", []string{"3", "12", "2"}},
		{"line", []string{"4", "17", "2"}},
	}

	for i, tt := range tests {
		entry, err := l.Next()
		if err != nil {
			t.Fatalf("unexpected error on token %d: %v", i, err)
		}
		if entry.Value != tt.value {
			t.Fatalf("token %d value = %q, want %q", i, entry.Value, tt.value)
		}
		for j, p := range tt.path {
			if entry.Path[j] != p {
				t.Fatalf("token %d path[%d] = %q, want %q", i, j, entry.Path[j], p)
			}
		}
	}

	if _, err := l.Next(); err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}
