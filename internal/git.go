package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/radulucut/gitbrag/internal/utils"
)

type GitStats struct {
	Repositories int
	FilesChanged int
	Insertions   int
	Deletions    int
}

func (g *GitStats) Add(other GitStats) {
	g.FilesChanged += other.FilesChanged
	g.Insertions += other.Insertions
	g.Deletions += other.Deletions
}

// isGitRepo checks if a directory is a git repository
func isGitRepo(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

type GitStatsOptions struct {
	Since  string
	Until  string
	Author string
}

func getGitStats(dir string, opts *GitStatsOptions) (GitStats, error) {
	stats := GitStats{}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return stats, utils.NewInternalError("directory does not exist: " + dir)
	}

	// Check if it's a git repo
	if !isGitRepo(dir) {
		return stats, utils.NewInternalError("not a git repository: " + dir)
	}

	// Build git log command with shortstat
	args := []string{"log", "--pretty=", "--numstat", "--branches"}
	if opts.Since != "" {
		args = append(args, "--since="+opts.Since)
	}
	if opts.Until != "" {
		args = append(args, "--until="+opts.Until)
	}
	if opts.Author != "" {
		args = append(args, "--author="+opts.Author)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		// Check if it's an empty repository
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "does not have any commits yet") ||
				strings.Contains(stderr, "bad default revision") {
				// Empty repository - return zero stats
				return stats, nil
			}
		}
		return stats, utils.NewInternalError("failed to execute git command: " + err.Error())
	}

	// Parse the output
	lines := strings.Split(string(output), "\n")
	filesMap := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		insertions := parts[0]
		deletions := parts[1]
		filename := parts[2]

		// Track unique files
		filesMap[filename] = true

		// Parse insertions
		if insertions != "-" {
			if n, err := strconv.Atoi(insertions); err == nil {
				stats.Insertions += n
			}
		}

		// Parse deletions
		if deletions != "-" {
			if n, err := strconv.Atoi(deletions); err == nil {
				stats.Deletions += n
			}
		}
	}

	stats.FilesChanged = len(filesMap)

	return stats, nil
}
