package storage

import (
	"testing"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name      string
		stype     string
		expectErr bool
	}{
		{
			name:  "memory storage",
			stype: "memory",
		},
		{
			name:      "unknown type",
			stype:     "gdfmg4t",
			expectErr: true,
		},
		{
			name:      "empty type",
			stype:     "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStorage(tt.stype)
			if (err != nil) != tt.expectErr {
				t.Fatalf("NewStorage(%q) error = %v, wantErr %v", tt.stype, err, tt.expectErr)
			}
			if !tt.expectErr && got == nil {
				t.Errorf("NewStorage(%q) returned nil", tt.stype)
			}
		})
	}
}
