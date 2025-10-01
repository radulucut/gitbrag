package internal

import (
	"image/color"
	"path/filepath"
	"strings"
)

// languageColors maps programming languages to their representative colors
// Colors are based on GitHub's language colors and popular conventions
var languageColors = map[string]color.RGBA{
	"Go":               {0, 173, 216, 255},   // #00ADD8
	"JavaScript":       {241, 224, 90, 255},  // #F1E05A
	"TypeScript":       {43, 116, 137, 255},  // #2B7489
	"Python":           {53, 114, 165, 255},  // #3572A5
	"Java":             {176, 114, 25, 255},  // #B07219
	"C":                {85, 85, 85, 255},    // #555555
	"C++":              {243, 75, 125, 255},  // #F34B7D
	"C/C++":            {243, 75, 125, 255},  // #F34B7D
	"C#":               {23, 134, 0, 255},    // #178600
	"Ruby":             {112, 21, 22, 255},   // #701516
	"PHP":              {79, 93, 149, 255},   // #4F5D95
	"Swift":            {240, 81, 56, 255},   // #F05138
	"Kotlin":           {161, 103, 224, 255}, // #A167E0
	"Rust":             {222, 165, 132, 255}, // #DEA584
	"Scala":            {194, 45, 64, 255},   // #C22D40
	"Shell":            {137, 224, 81, 255},  // #89E051
	"HTML":             {227, 76, 38, 255},   // #E34C26
	"CSS":              {86, 61, 124, 255},   // #563D7C
	"SCSS":             {198, 83, 140, 255},  // #C6538C
	"Sass":             {191, 64, 191, 255},  // #BF40BF
	"Less":             {29, 54, 93, 255},    // #1D365D
	"Vue":              {65, 184, 131, 255},  // #41B883
	"Dart":             {0, 180, 171, 255},   // #00B4AB
	"Lua":              {0, 0, 128, 255},     // #000080
	"Perl":             {2, 152, 195, 255},   // #0298C3
	"Elixir":           {110, 74, 126, 255},  // #6E4A7E
	"Clojure":          {219, 88, 85, 255},   // #DB5855
	"Elm":              {96, 181, 204, 255},  // #60B5CC
	"Erlang":           {184, 57, 152, 255},  // #B83998
	"Haskell":          {94, 80, 134, 255},   // #5E5086
	"Julia":            {162, 112, 186, 255}, // #A270BA
	"Nim":              {255, 202, 56, 255},  // #FFCA38
	"R":                {25, 140, 231, 255},  // #198CE7
	"Objective-C":      {67, 142, 255, 255},  // #438EFF
	"Protocol Buffers": {66, 66, 66, 255},    // #424242
	"GraphQL":          {225, 0, 152, 255},   // #E10098
	"Terraform":        {92, 73, 149, 255},   // #5C4D95
	"Dockerfile":       {56, 77, 84, 255},    // #384D54
	"Makefile":         {66, 120, 25, 255},   // #427819
	"Rakefile":         {112, 21, 22, 255},   // #701516
	"Gemfile":          {112, 21, 22, 255},   // #701516
	"JSON":             {41, 41, 41, 255},    // #292929
	"XML":              {0, 96, 176, 255},    // #0060B0
	"YAML":             {203, 56, 55, 255},   // #CB3837
	"Markdown":         {83, 89, 101, 255},   // #535965
	"SQL":              {224, 147, 0, 255},   // #E09300
}

// getLanguageColor returns the color for a given language.
// If the language is not found, it returns a default gray color.
func getLanguageColor(language string) color.RGBA {
	if col, ok := languageColors[language]; ok {
		return col
	}
	// Default gray for unknown languages
	return color.RGBA{150, 150, 150, 255}
}

// detectLanguage returns the language name based on file extension
func detectLanguage(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		// Check for common files without extensions
		base := filepath.Base(filename)
		switch base {
		case "Dockerfile", "Makefile", "Rakefile", "Gemfile":
			return base
		}
		return ""
	}

	ext = strings.ToLower(ext)
	langMap := map[string]string{
		".go":         "Go",
		".js":         "JavaScript",
		".ts":         "TypeScript",
		".jsx":        "JavaScript",
		".tsx":        "TypeScript",
		".py":         "Python",
		".java":       "Java",
		".c":          "C",
		".cpp":        "C++",
		".cc":         "C++",
		".cxx":        "C++",
		".h":          "C/C++",
		".hpp":        "C++",
		".cs":         "C#",
		".rb":         "Ruby",
		".php":        "PHP",
		".swift":      "Swift",
		".kt":         "Kotlin",
		".rs":         "Rust",
		".scala":      "Scala",
		".sh":         "Shell",
		".bash":       "Shell",
		".zsh":        "Shell",
		".html":       "HTML",
		".css":        "CSS",
		".scss":       "SCSS",
		".sass":       "Sass",
		".less":       "Less",
		".json":       "JSON",
		".xml":        "XML",
		".yaml":       "YAML",
		".yml":        "YAML",
		".md":         "Markdown",
		".sql":        "SQL",
		".r":          "R",
		".m":          "Objective-C",
		".vue":        "Vue",
		".dart":       "Dart",
		".lua":        "Lua",
		".pl":         "Perl",
		".ex":         "Elixir",
		".exs":        "Elixir",
		".clj":        "Clojure",
		".elm":        "Elm",
		".erl":        "Erlang",
		".hs":         "Haskell",
		".jl":         "Julia",
		".nim":        "Nim",
		".proto":      "Protocol Buffers",
		".graphql":    "GraphQL",
		".tf":         "Terraform",
		".dockerfile": "Dockerfile",
	}

	if lang, ok := langMap[ext]; ok {
		return lang
	}

	return ""
}
