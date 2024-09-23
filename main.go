package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var link1 string = "Project"
var link2 string = "Coding"
var link3 string = "go"

// 1. Build index of links in each note.
// 2. Search for notes containing searched links
// 3. Return list of notes that contain the links. (2+ links have to match)

type Notes struct {
	Filename string
	Links    map[string]bool
}

func buildIndex(directory string) ([]Notes, error) {
	var index []Notes
	linkRegex := regexp.MustCompile(`\[\[(.*?)\]\]`)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			links := make(map[string]bool)
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				matches := linkRegex.FindAllStringSubmatch(scanner.Text(), -1)
				for _, match := range matches {
					links[match[1]] = true
				}
			}

			index = append(index, Notes{
				Filename: path,
				Links:    links,
			})
		}
		return nil
	})

	return index, err

}

func search(index []Notes, searchTerms []string) map[string][]string {
	results := make(map[string][]string)
	for _, file := range index {
		matchedTerms := []string{}
		for _, term := range searchTerms {
			if file.Links[term] {
				matchedTerms = append(matchedTerms, term)
			}
		}
		if len(matchedTerms) > 0 {
			key := fmt.Sprintf("Matches '%s':", strings.Join(matchedTerms, "', '"))
			results[key] = append(results[key], file.Filename)
		}
	}
	return results
}

func main() {
	index, err := buildIndex("path/to/your/markdown/directory")
	if err != nil {
		fmt.Println("Error building index:", err)
		return
	}

	searchTerms := []string{"go", "coding", "project"}
	results := search(index, searchTerms)

	// Sort results by number of matches (descending)
	for _, terms := range [][]string{{"project", "go", "coding"}, {"project", "go"}, {"project", "coding"}, {"go", "coding"}, {"project"}, {"go"}, {"coding"}} {
		key := fmt.Sprintf("Matches '%s':", strings.Join(terms, "', '"))
		if files, ok := results[key]; ok {
			fmt.Println(key)
			for _, file := range files {
				fmt.Println(file)
			}
		}
	}
}
