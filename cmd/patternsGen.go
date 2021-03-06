package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	patternNames map[string]bool
	patterns     []*pattern
	rootDir      string

	grokPatternNameRegex = `[a-zA-Z0-9_]+`
	grokPatternName      = regexp.MustCompile(grokPatternNameRegex)
)

type pattern struct {
	name, regex string
}

func main() {
	rootDir = os.Args[1]

	stat, err := os.Stat(rootDir)
	if err != nil {
		panic(err)
	}

	if !stat.IsDir() {
		fmt.Println("Path must be a directory.")
	}

	patterns = make([]*pattern, 0, 50)
	patternNames = make(map[string]bool, 50)
	// The regex for a grok variable name is defined above, the grok
	// filter module uses it to build a variable pattern.
	patterns = append(patterns, &pattern{
		name:  "GROK_VARIABLE",
		regex: grokPatternNameRegex,
	})
	patternNames["GROK_VARIABLE"] = true

	if err := filepath.Walk(rootDir, walk); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	outFile, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0655)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := writeOutput(outFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() || err != nil {
		return err
	}

	if !strings.HasSuffix(path, ".p") {
		return nil
	}

	return processfile(path)
}

func processfile(path string) error {
	path, _ = filepath.Abs(path)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	lineNum := 1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			lineNum++
			continue
		}

		parts := strings.SplitAfterN(line, " ", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Error in file %s on line %d", path, lineNum)
		}

		patternName := strings.TrimSpace(parts[0])
		if !grokPatternName.MatchString(patternName) {
			return fmt.Errorf("Invalid grok pattern name \"%s\"", patternName)
		}

		if patternNames[patternName] {
			fmt.Printf("Pattern %s already defined. File %s, line %d\n", patternName, path, lineNum)
			lineNum++
			continue
		}

		patternRegex := strings.TrimSpace(parts[1])

		patterns = append(patterns, &pattern{
			name:  patternName,
			regex: patternRegex,
		})
		patternNames[patternName] = true
		lineNum++
	}

	return nil
}

func writeOutput(outFile *os.File) error {
	fmt.Fprint(outFile, `// This file was generated by cmd/patternsGen.go.
// To change a pattern, edit the correct file under patterns and run
// "make generate" from the project root

package grok

var grokPatterns = map[string]string{
`)

	for _, pattern := range patterns {
		fmt.Fprintf(outFile, "	\"%s\": %s,\n", pattern.name, strconv.Quote(pattern.regex))
	}

	fmt.Fprint(outFile, "}\n")
	return nil
}
