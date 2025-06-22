package commands

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestVersion(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Version()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expected := "cli-tool v1.0.0\n"
	if output != expected {
		t.Errorf("Version() = %q, want %q", output, expected)
	}
}