package internal

import (
	"image/color"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    color.RGBA
		wantErr bool
	}{
		{
			name:  "6 character hex with hash",
			input: "#282a36",
			want:  color.RGBA{40, 42, 54, 255},
		},
		{
			name:  "6 character hex without hash",
			input: "282a36",
			want:  color.RGBA{40, 42, 54, 255},
		},
		{
			name:  "8 character hex with alpha",
			input: "#282a36ff",
			want:  color.RGBA{40, 42, 54, 255},
		},
		{
			name:  "8 character hex with transparency",
			input: "282a3680",
			want:  color.RGBA{40, 42, 54, 128},
		},
		{
			name:  "3 character hex shorthand",
			input: "#fff",
			want:  color.RGBA{255, 255, 255, 255},
		},
		{
			name:  "3 character hex shorthand without hash",
			input: "000",
			want:  color.RGBA{0, 0, 0, 255},
		},
		{
			name:  "3 character hex shorthand colors",
			input: "f0a",
			want:  color.RGBA{255, 0, 170, 255},
		},
		{
			name:  "uppercase hex",
			input: "#FF5555",
			want:  color.RGBA{255, 85, 85, 255},
		},
		{
			name:  "mixed case hex",
			input: "8Be9Fd",
			want:  color.RGBA{139, 233, 253, 255},
		},
		{
			name:    "invalid length",
			input:   "12345",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			input:   "gggggg",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHexColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHexColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseHexColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetBackgroundFromHex(t *testing.T) {
	r := NewPNGRenderer()

	// Default should be transparent
	if r.bg != (color.RGBA{0, 0, 0, 0}) {
		t.Errorf("Default background should be transparent, got %v", r.bg)
	}

	// Set a valid color
	err := r.SetBackgroundFromHex("#282a36")
	if err != nil {
		t.Errorf("SetBackgroundFromHex() unexpected error: %v", err)
	}

	expected := color.RGBA{40, 42, 54, 255}
	if r.bg != expected {
		t.Errorf("SetBackgroundFromHex() bg = %v, want %v", r.bg, expected)
	}

	// Try invalid color
	err = r.SetBackgroundFromHex("invalid")
	if err == nil {
		t.Error("SetBackgroundFromHex() expected error for invalid color")
	}
}
