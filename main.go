package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/junwei890/crawler/src"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	mongoURI := os.Getenv("MONGO_URI")

	linksInBytes, err := os.ReadFile("crawler.txt")
	if err != nil {
		log.Fatal(err)
	}
	links := strings.Fields(string(linksInBytes))

	start := time.Now()
	if err := src.StartCrawl(mongoURI, links); err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Since(start))
}
