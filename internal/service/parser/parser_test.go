package parser

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expect    string
		expectErr bool
	}{
		{
			name:   "valid http url",
			input:  "http://example.com",
			expect: "http://example.com",
		},
		{
			name:   "valid https url",
			input:  "https://example.com",
			expect: "https://example.com",
		},
		{
			name:   "uppercase",
			input:  "HTTP://EXAMPLE.COM",
			expect: "http://example.com",
		},
		{
			name:   "mixedcase",
			input:  "https://Example.Com/path",
			expect: "https://example.com/path",
		},
		{
			name:   "trailing slash removed",
			input:  "https://example.com/path/",
			expect: "https://example.com/path",
		},
		{
			name:   "params",
			input:  "https://example.com/path?key=value",
			expect: "https://example.com/path?key=value",
		},
		{
			name:   "fragment",
			input:  "https://example.com/path#section",
			expect: "https://example.com/path#section",
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:      "no scheme",
			input:     "example.com",
			expectErr: true,
		},
		{
			name:      "invalid url",
			input:     "://missing-scheme",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURL(tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("ParseURL(%q) error = %v, wantErr %v", tt.input, err, tt.expectErr)
			}
			if got != tt.expect {
				t.Errorf("ParseURL(%q) = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}
