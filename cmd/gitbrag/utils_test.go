package gitbrag

import (
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

var (
	defaultCurrentTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

func createGitRepo(t *testing.T) string {
	testDir := "test_gitbrag_" + t.Name()
	t.Cleanup(func() {
		os.RemoveAll(testDir)
	})

	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("git", "init", "-b", "main")
	cmd.Dir = testDir
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		t.Log(string(out))
		t.Fatal(err)
	}
	if err := os.WriteFile(path.Join(testDir, "main.go"), []byte(`
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path.Join(testDir, "main.ts"), []byte(`
console.log("Hello, World!");
console.log("Hello, World!");
`), 0644); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = testDir
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		t.Log(string(out))
		t.Fatal(err)
	}

	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = testDir
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		t.Log(string(out))
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(testDir, "main.ts"), []byte(`
console.log("Hello, World!");
`), 0644); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "checkout", "-b", "feature")
	cmd.Dir = testDir
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		t.Log(string(out))
		t.Fatal(err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = testDir
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		t.Log(string(out))
		t.Fatal(err)
	}

	cmd = exec.Command("git", "commit", "-m", "second commit", "--author", "John Doe <john.doe@example.com>")
	cmd.Dir = testDir
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		t.Log(string(out))
		t.Fatal(err)
	}

	cmd = exec.Command("git", "checkout", "main")
	cmd.Dir = testDir
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		t.Log(string(out))
		t.Fatal(err)
	}
	return testDir
}
