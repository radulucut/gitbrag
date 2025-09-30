package main

import "github.com/radulucut/gitbrag/cmd/gitbrag"

var (
	version = "dev"
)

func main() {
	gitbrag.Version = version
	gitbrag.Execute()
}
