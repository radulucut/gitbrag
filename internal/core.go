package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/radulucut/gitbrag/internal/utils"
)

type Core struct {
	time    utils.Time
	printer *Printer
}

func NewCore(time utils.Time, printer *Printer) *Core {
	return &Core{
		time:    time,
		printer: printer,
	}
}

type RunOptions struct {
	Dirs         []string
	Since        time.Time
	Until        time.Time
	DateRange    string
	Author       string
	Output       string
	Background   string
	Color        string
	Lang         bool
	ExcludeFiles *regexp.Regexp
}

func (c *Core) Run(opts *RunOptions) error {
	if len(opts.Dirs) == 0 {
		return utils.NewInternalError("no directories specified")
	}

	totalStats := &GitStats{}

	gitOpts := &GitStatsOptions{
		Author:       opts.Author,
		ExcludeFiles: opts.ExcludeFiles,
	}
	if !opts.Since.IsZero() {
		gitOpts.Since = opts.Since.Format(time.RFC3339)
	}
	if !opts.Until.IsZero() {
		gitOpts.Until = opts.Until.Format(time.RFC3339)
	}

	// Process each directory
	for _, dir := range opts.Dirs {
		// Convert to absolute path
		absDir, err := filepath.Abs(dir)
		if err != nil {
			c.printer.ErrPrintf("Warning: could not resolve path '%s': %v\n", dir, err)
			continue
		}

		// Check if directory exists
		info, err := os.Stat(absDir)
		if err != nil {
			c.printer.ErrPrintf("Warning: could not access '%s': %v\n", dir, err)
			continue
		}

		if !info.IsDir() {
			c.printer.ErrPrintf("Warning: '%s' is not a directory\n", dir)
			continue
		}

		// Check if it's a git repository
		if isGitRepo(absDir) {
			stats, err := getGitStats(absDir, gitOpts)
			if err != nil {
				c.printer.ErrPrintf("Warning: could not get git stats for '%s': %v\n", dir, err)
				continue
			}
			totalStats.Add(stats)
			totalStats.Repositories++
		} else {
			// Try to find git repos in subdirectories
			err := c.processSubdirectories(absDir, gitOpts, totalStats)
			if err != nil {
				c.printer.ErrPrintf("Warning: error processing subdirectories in '%s': %v\n", dir, err)
			}
		}
	}

	// Output results
	if totalStats.Repositories == 0 {
		c.printer.Println("No git repositories found in the specified directories.")
		return nil
	}

	if !opts.Since.IsZero() && !opts.Until.IsZero() {
		opts.DateRange = fmt.Sprintf("%s - %s", formatOutputDate(opts.Since), formatOutputDate(opts.Until))
	} else if !opts.Since.IsZero() {
		opts.DateRange = fmt.Sprintf("Since %s", formatOutputDate(opts.Since))
	} else if !opts.Until.IsZero() {
		opts.DateRange = fmt.Sprintf("Until %s", formatOutputDate(opts.Until))
	}

	// Check if PNG output is requested
	if opts.Output != "" {
		pngRenderer := NewPNGRenderer()
		if opts.Background != "" {
			if err := pngRenderer.SetBackgroundFromHex(opts.Background); err != nil {
				return fmt.Errorf("invalid background color: %w", err)
			}
		}
		if opts.Color != "" {
			if err := pngRenderer.SetForegroundFromHex(opts.Color); err != nil {
				return fmt.Errorf("invalid text color: %w", err)
			}
		}
		if err := pngRenderer.RenderToFile(totalStats, opts); err != nil {
			return fmt.Errorf("failed to export PNG: %w", err)
		}
		c.printer.Printf("Statistics exported to %s\n", opts.Output)
		return nil
	}

	// Print date range if available
	if opts.DateRange != "" {
		c.printer.Printf("%s\n\n", opts.DateRange)
	}

	filesStr := fmt.Sprint(totalStats.FilesChanged)
	insertionsStr := fmt.Sprint(totalStats.Insertions)
	deletionsStr := fmt.Sprint(totalStats.Deletions)

	maxLen := max(len(filesStr), len(insertionsStr), len(deletionsStr))

	filesStr = fmt.Sprintf("%*s", maxLen, filesStr)
	insertionsStr = fmt.Sprintf("%*s", maxLen, insertionsStr)
	deletionsStr = fmt.Sprintf("%*s", maxLen, deletionsStr)
	c.printer.Printf(`%s files changed
%s insertions(+)
%s deletions(-)
`, filesStr, insertionsStr, deletionsStr)

	return nil
}

func (c *Core) processSubdirectories(dir string, opts *GitStatsOptions, totalStats *GitStats) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("could not read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories except .git
		if entry.Name()[0] == '.' && entry.Name() != ".git" {
			continue
		}

		subDir := filepath.Join(dir, entry.Name())

		if isGitRepo(subDir) {
			stats, err := getGitStats(subDir, opts)
			if err != nil {
				c.printer.ErrPrintf("Warning: could not get git stats for '%s': %v\n", subDir, err)
				continue
			}
			totalStats.Add(stats)
			totalStats.Repositories++
		}
	}

	return nil
}

func formatOutputDate(t time.Time) string {
	// Check if time component is zero (midnight)
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return t.Format("Jan 2, 2006")
	}
	return t.Format("Jan 2, 2006 15:04:05")
}
