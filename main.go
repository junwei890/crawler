package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/junwei890/crawler/src"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fileToRead := "crawler.txt"
	path := path.Join(homeDir, fileToRead)

	// #nosec G304
	fullFile, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Fields(string(fullFile))

	mongoURI := ""
	links := []string{}
	for i, line := range lines {
		if i == 0 {
			mongoURI = line
		} else if i > 0 {
			links = append(links, strings.TrimSpace(line))
		}
	}

	if err := src.StartCrawl(mongoURI, links); err != nil {
		log.Fatal(err)
	}
}
