package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sort"
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

func (r *PNGRenderer) RenderToFile(stats *GitStats, opts *RunOptions) error {
	if r.fontFace == nil {
		return fmt.Errorf("font not loaded")
	}

	if opts.Lang && len(stats.Languages) > 0 {
		r.height = 950 // Add extra space for language bar and labels
	}

	img := image.NewRGBA(image.Rect(0, 0, r.width, r.height))

	// Fill background more efficiently
	draw.Draw(img, img.Bounds(), &image.Uniform{r.bg}, image.Point{}, draw.Src)

	// add start padding to align numbers
	filesStr := fmt.Sprint(stats.FilesChanged)
	insertionsStr := fmt.Sprint(stats.Insertions)
	deletionsStr := fmt.Sprint(stats.Deletions)

	maxLen := max(len(filesStr), len(insertionsStr), len(deletionsStr))

	filesStr = fmt.Sprintf("%*s files changed", maxLen, filesStr)
	insertionsStr = fmt.Sprintf("%*s insertions(+)", maxLen, insertionsStr)
	deletionsStr = fmt.Sprintf("%*s deletions(-)", maxLen, deletionsStr)

	// add end padding to center text
	maxLen = max(maxLen, len(filesStr), len(insertionsStr), len(deletionsStr))
	filesStr = fmt.Sprintf("%-*s", maxLen, filesStr)
	insertionsStr = fmt.Sprintf("%-*s", maxLen, insertionsStr)
	deletionsStr = fmt.Sprintf("%-*s", maxLen, deletionsStr)

	greenColor := color.RGBA{26, 127, 55, 255} // Green for insertions
	redColor := color.RGBA{209, 36, 47, 255}   // Red for deletions

	// Draw date range if available
	yOffset := 280
	if opts.DateRange != "" {
		// Calculate text width to center it
		textWidth := font.MeasureString(r.fontFace, opts.DateRange).Ceil()
		dateRangeX := (r.width - textWidth) / 2
		r.drawTextAntialiased(img, opts.DateRange, dateRangeX, yOffset, r.fg)
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

	// Draw language breakdown if requested
	if opts.Lang && len(stats.Languages) > 0 {
		r.drawLanguageBar(img, stats, yOffset+280)
	}

	// Save to file
	f, err := os.Create(opts.Output)
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

// LanguageInfo holds information about a language's usage
type LanguageInfo struct {
	Name       string
	Lines      int
	Percentage float64
	Color      color.RGBA
}

// drawLanguageBar draws a horizontal bar chart showing language breakdown
func (r *PNGRenderer) drawLanguageBar(img *image.RGBA, stats *GitStats, yOffset int) {
	if len(stats.Languages) == 0 {
		return
	}

	// Calculate total lines
	totalLines := 0
	for _, lines := range stats.Languages {
		totalLines += lines
	}

	if totalLines == 0 {
		return
	}

	// Sort languages by lines (descending)
	var languages []LanguageInfo
	for lang, lines := range stats.Languages {
		percentage := float64(lines) / float64(totalLines) * 100
		languages = append(languages, LanguageInfo{
			Name:       lang,
			Lines:      lines,
			Percentage: percentage,
			Color:      getLanguageColor(lang),
		})
	}

	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Lines > languages[j].Lines
	})

	// Group into top 3 and others
	var displayLangs []LanguageInfo
	othersLines := 0
	othersPercentage := 0.0

	for i, lang := range languages {
		if i < 3 {
			displayLangs = append(displayLangs, lang)
		} else {
			othersLines += lang.Lines
			othersPercentage += lang.Percentage
		}
	}

	// Add "Others" category if there are more than 3 languages
	if len(languages) > 3 {
		displayLangs = append(displayLangs, LanguageInfo{
			Name:       "Other",
			Lines:      othersLines,
			Percentage: othersPercentage,
			Color:      color.RGBA{150, 150, 150, 255}, // Gray for other
		})
	}

	// Draw the bar
	barWidth := 600
	barHeight := 40
	barX := (r.width - barWidth) / 2
	barY := yOffset

	// Draw each language segment
	currentX := barX
	for _, lang := range displayLangs {
		segmentWidth := int(float64(barWidth) * lang.Percentage / 100)
		if segmentWidth > 0 {
			// Draw the colored segment
			segmentRect := image.Rect(currentX, barY, currentX+segmentWidth, barY+barHeight)
			draw.Draw(img, segmentRect, &image.Uniform{lang.Color}, image.Point{}, draw.Src)
			currentX += segmentWidth
		}
	}

	// Draw labels on the same line below the bar with colored circles
	labelY := barY + barHeight + 40
	circleRadius := 8
	circleSpacing := 10
	labelPadding := 30 // Space between different labels

	currentX = barX

	for _, lang := range displayLangs {
		// Format the label: just the language name
		label := lang.Name

		// Draw colored circle
		circleX := currentX + circleRadius
		circleY := labelY - 6 // Adjust to align with text baseline
		r.drawFilledCircle(img, circleX, circleY, circleRadius, lang.Color)

		// Draw text after the circle
		textX := currentX + (circleRadius * 2) + circleSpacing
		r.drawTextAntialiased(img, label, textX, labelY, r.fg)

		// Move to next label position
		labelWidth := font.MeasureString(r.fontFace, label).Ceil()
		currentX = textX + labelWidth + labelPadding
	}
}

// drawFilledCircle draws a filled circle at the given position
func (r *PNGRenderer) drawFilledCircle(img *image.RGBA, centerX, centerY, radius int, col color.RGBA) {
	// Use midpoint circle algorithm to draw a filled circle
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				px := centerX + x
				py := centerY + y
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, col)
				}
			}
		}
	}
}
