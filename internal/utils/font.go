package utils

import (
	_ "embed"
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed SpaceMonoRegular.ttf
var spaceMonoRegular []byte

// LoadFont loads the embedded Space Mono Regular font with the specified size
func LoadFont(size float64) (font.Face, error) {
	f, err := opentype.Parse(spaceMonoRegular)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %w", err)
	}

	return face, nil
}
