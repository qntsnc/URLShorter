package storage

import (
	"testing"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name    string
		stype   string
		wantErr bool
	}{
		{
			name:  "memory storage",
			stype: "memory",
		},
		{
			name:    "unknown type",
			stype:   "gdfmg4t",
			wantErr: true,
		},
		{
			name:    "empty type",
			stype:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStorage(tt.stype)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewStorage(%q) error = %v, wantErr %v", tt.stype, err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NewStorage(%q) returned nil", tt.stype)
			}
		})
	}
}
