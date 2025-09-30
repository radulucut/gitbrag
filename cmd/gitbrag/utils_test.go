package gitbrag

import (
	"os"
	"os/exec"
	"path"
	"time"
)

var (
	defaultCurrentTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

func createGitRepo(dir string) error {
	// init a repo, create a main.go and a main.ts file
	// commit files
	if err := os.Chdir(dir); err != nil {
		return err
	}
	if err := exec.Command("git", "init").Run(); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(dir, "main.go"), []byte(`
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(dir, "main.ts"), []byte(`
console.log("Hello, World!");
console.log("Hello, World!");
`), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "commit", "-m", "initial commit").Run(); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(dir, "main.ts"), []byte(`
console.log("Hello, World!");
`), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "commit", "-m", "second commit").Run(); err != nil {
		return err
	}
	return nil
}
