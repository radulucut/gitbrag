package internal

import (
	"fmt"
	"os"
	"path/filepath"
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
	Dirs  []string
	Since time.Time
}

func (c *Core) Run(opts *RunOptions) error {
	if len(opts.Dirs) == 0 {
		return utils.NewInternalError("no directories specified")
	}

	totalStats := GitStats{}
	repoCount := 0

	// Format the since parameter for git
	var sinceStr string
	if !opts.Since.IsZero() {
		sinceStr = opts.Since.Format(time.RFC3339)
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
			stats, err := getGitStats(absDir, sinceStr)
			if err != nil {
				c.printer.ErrPrintf("Warning: could not get git stats for '%s': %v\n", dir, err)
				continue
			}
			totalStats.Add(stats)
			repoCount++
		} else {
			// Try to find git repos in subdirectories
			err := c.processSubdirectories(absDir, &totalStats, &repoCount, sinceStr)
			if err != nil {
				c.printer.ErrPrintf("Warning: error processing subdirectories in '%s': %v\n", dir, err)
			}
		}
	}

	// Output results
	if repoCount == 0 {
		c.printer.Println("No git repositories found in the specified directories.")
		return nil
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

func (c *Core) processSubdirectories(dir string, totalStats *GitStats, repoCount *int, sinceStr string) error {
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
			stats, err := getGitStats(subDir, sinceStr)
			if err != nil {
				c.printer.ErrPrintf("Warning: could not get git stats for '%s': %v\n", subDir, err)
				continue
			}
			totalStats.Add(stats)
			*repoCount++
		}
	}

	return nil
}
