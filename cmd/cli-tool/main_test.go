package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
		verbose  bool
	}{
		{
			name:     "simple sequence diagram",
			input:    "sequenceDiagram\n    A->B: Hello",
			expected: "<mxGraphModel",
			wantErr:  false,
			verbose:  false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: "<mxGraphModel",
			wantErr:  false,
			verbose:  false,
		},
		{
			name:     "with participants",
			input:    "sequenceDiagram\n    participant A as Alice\n    A->B: Test",
			expected: "Alice",
			wantErr:  false,
			verbose:  false,
		},
		{
			name:     "ER diagram",
			input:    "erDiagram\n    USER {\n        int id PK\n    }",
			expected: "USER",
			wantErr:  false,
			verbose:  false,
		},
		{
			name:     "ER diagram with relationship",
			input:    "erDiagram\n    USER {}\n    ORDER {}\n    USER ||--o{ ORDER : places",
			expected: "places",
			wantErr:  false,
			verbose:  false,
		},
		{
			name:     "verbose mode",
			input:    "sequenceDiagram\n    A->B: Hello",
			expected: "<mxGraphModel",
			wantErr:  false,
			verbose:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdin
			oldStdin := os.Stdin
			r, w, _ := os.Pipe()
			os.Stdin = r
			
			// Redirect stdout
			oldStdout := os.Stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut
			
			// Write input
			go func() {
				defer w.Close()
				w.Write([]byte(tt.input))
			}()
			
			// Run the function
			err := run(tt.verbose)
			
			// Close stdout writer and read output
			wOut.Close()
			os.Stdout = oldStdout
			os.Stdin = oldStdin
			
			var buf bytes.Buffer
			buf.ReadFrom(rOut)
			output := buf.String()
			
			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, output)
			}
		})
	}
}

