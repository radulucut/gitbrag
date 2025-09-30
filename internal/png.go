package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type PNGRenderer struct {
	width  int
	height int
	bg     color.Color
	fg     color.Color
}

func NewPNGRenderer() *PNGRenderer {
	return &PNGRenderer{
		width:  600,
		height: 600,
		bg:     color.RGBA{0, 0, 0, 0},   // Transparent by default
		fg:     color.RGBA{0, 0, 0, 255}, // Black text by default
	}
}

func (r *PNGRenderer) SetBackgroundFromHex(hexColor string) error {
	col, err := parseHexColor(hexColor)
	if err != nil {
		return err
	}
	r.bg = col
	return nil
}

func (r *PNGRenderer) SetForegroundFromHex(hexColor string) error {
	col, err := parseHexColor(hexColor)
	if err != nil {
		return err
	}
	r.fg = col
	return nil
}

func (r *PNGRenderer) RenderToFile(stats *GitStats, filepath string) error {
	img := image.NewRGBA(image.Rect(0, 0, r.width, r.height))

	// Fill background
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			img.Set(x, y, r.bg)
		}
	}

	filesStr := fmt.Sprint(stats.FilesChanged)
	insertionsStr := fmt.Sprint(stats.Insertions)
	deletionsStr := fmt.Sprint(stats.Deletions)

	maxLen := max(len(filesStr), len(insertionsStr), len(deletionsStr))

	filesStr = fmt.Sprintf("%*s files changed", maxLen, filesStr)
	insertionsStr = fmt.Sprintf("%*s insertions(+)", maxLen, insertionsStr)
	deletionsStr = fmt.Sprintf("%*s deletions(-)", maxLen, deletionsStr)

	greenColor := color.RGBA{26, 127, 55, 255} // Green for insertions
	redColor := color.RGBA{209, 36, 47, 255}   // Red for deletions

	r.drawText(img, filesStr, 230, 260, r.fg, 1)
	r.drawText(img, insertionsStr, 230, 290, greenColor, 1)
	r.drawText(img, deletionsStr, 230, 320, redColor, 1)

	// Save to file
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

func (r *PNGRenderer) drawText(img *image.RGBA, text string, x, y int, col color.Color, scale int) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	// Draw text multiple times for scaling effect (simple bold/larger text)
	for i := range scale {
		for j := range scale {
			d.Dot = fixed.Point26_6{
				X: fixed.Int26_6((x + i) * 64),
				Y: fixed.Int26_6((y + j) * 64),
			}
			d.DrawString(text)
		}
	}
}

func parseHexColor(s string) (color.RGBA, error) {
	s = strings.TrimPrefix(s, "#")

	var r, g, b, a uint8
	a = 255 // Default to fully opaque

	switch len(s) {
	case 6:
		// RGB format
		rgb, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("invalid hex color format")
		}
		r = uint8(rgb >> 16)
		g = uint8(rgb >> 8)
		b = uint8(rgb)
	case 8:
		// RGBA format
		rgba, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("invalid hex color format")
		}
		r = uint8(rgba >> 24)
		g = uint8(rgba >> 16)
		b = uint8(rgba >> 8)
		a = uint8(rgba)
	case 3:
		// Short RGB format (e.g., "FFF" -> "FFFFFF")
		rgb, err := strconv.ParseUint(s, 16, 16)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("invalid hex color format")
		}
		r = uint8((rgb >> 8) & 0xF * 17)
		g = uint8((rgb >> 4) & 0xF * 17)
		b = uint8((rgb & 0xF) * 17)
	default:
		return color.RGBA{}, fmt.Errorf("hex color must be 3, 6, or 8 characters")
	}

	return color.RGBA{r, g, b, a}, nil
}
