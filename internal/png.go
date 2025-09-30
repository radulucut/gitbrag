package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/radulucut/gitbrag/internal/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type PNGRenderer struct {
	width    int
	height   int
	bg       color.Color
	fg       color.Color
	fontFace font.Face
}

func NewPNGRenderer() *PNGRenderer {
	// Load the embedded Space Mono font with size 24 for better readability
	fontFace, err := utils.LoadFont(24)
	if err != nil {
		// Fallback to basicfont if custom font fails
		fontFace = basicfont.Face7x13
	}

	return &PNGRenderer{
		width:    800,                      // Increased resolution for better quality
		height:   800,                      // Increased resolution for better quality
		bg:       color.RGBA{0, 0, 0, 0},   // Transparent by default
		fg:       color.RGBA{0, 0, 0, 255}, // Black text by default
		fontFace: fontFace,
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

func (r *PNGRenderer) RenderToFile(stats *GitStats, filepath string, dateRange string) error {
	if r.fontFace == nil {
		return fmt.Errorf("font not loaded")
	}

	img := image.NewRGBA(image.Rect(0, 0, r.width, r.height))

	// Fill background more efficiently
	draw.Draw(img, img.Bounds(), &image.Uniform{r.bg}, image.Point{}, draw.Src)

	filesStr := fmt.Sprint(stats.FilesChanged)
	insertionsStr := fmt.Sprint(stats.Insertions)
	deletionsStr := fmt.Sprint(stats.Deletions)

	maxLen := max(len(filesStr), len(insertionsStr), len(deletionsStr))

	filesStr = fmt.Sprintf("%*s files changed", maxLen, filesStr)
	insertionsStr = fmt.Sprintf("%*s insertions(+)", maxLen, insertionsStr)
	deletionsStr = fmt.Sprintf("%*s deletions(-)", maxLen, deletionsStr)

	greenColor := color.RGBA{26, 127, 55, 255} // Green for insertions
	redColor := color.RGBA{209, 36, 47, 255}   // Red for deletions

	// Draw date range if available
	yOffset := 280
	if dateRange != "" {
		// Calculate text width to center it
		textWidth := font.MeasureString(r.fontFace, dateRange).Ceil()
		dateRangeX := (r.width - textWidth) / 2
		r.drawTextAntialiased(img, dateRange, dateRangeX, yOffset, r.fg)
	}

	// Center each stat line
	filesWidth := font.MeasureString(r.fontFace, filesStr).Ceil()
	filesX := (r.width - filesWidth) / 2
	r.drawTextAntialiased(img, filesStr, filesX, yOffset+100, r.fg)

	insertionsWidth := font.MeasureString(r.fontFace, insertionsStr).Ceil()
	insertionsX := (r.width - insertionsWidth) / 2
	r.drawTextAntialiased(img, insertionsStr, insertionsX, yOffset+150, greenColor)

	deletionsWidth := font.MeasureString(r.fontFace, deletionsStr).Ceil()
	deletionsX := (r.width - deletionsWidth) / 2
	r.drawTextAntialiased(img, deletionsStr, deletionsX, yOffset+200, redColor)

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

func (r *PNGRenderer) drawTextAntialiased(img *image.RGBA, text string, x, y int, col color.Color) {
	// Use proper fixed-point positioning for better text rendering
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: r.fontFace,
		Dot:  point,
	}

	// Draw the text once with proper anti-aliasing
	d.DrawString(text)
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
