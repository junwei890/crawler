package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/junwei890/crawler/src"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	dbURI := os.Getenv("DB_URI")

	linksInBytes, err := os.ReadFile("crawler.txt")
	if err != nil {
		log.Fatal(err)
	}

	if err := src.StartCrawl(dbURI, strings.Fields(string(linksInBytes))); err != nil {
		log.Fatal(err)
	}
}
