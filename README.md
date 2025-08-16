# Web crawler
This is a web crawler I wrote for personal use, though it slowly evolved into a piece of software that was kind of portable. If you find a web crawler useful, installation instructions are below.

## Requirements
- [Go](https://go.dev/doc/install) installed.
- A [MongoDB](https://www.mongodb.com/docs/atlas/getting-started/) cluster.

## Installation
Run `go install github.com/junwei890/crawler@latest` in your terminal.

## Usage
You need to have a `crawler.txt` file in your home directory, it should contain your MongoDB URI and links you would like to crawl. See [crawler.txt.example](https://github.com/junwei890/crawler/blob/main/crawler.txt.example) for how it should be formatted.

Run `crawler` in your terminal to start.

### Tips
- You don't have to create a MongoDB database and collection before crawling websites, the crawler takes care of that for you.
- You don't have to index the collection for [Atlas Search](https://www.mongodb.com/docs/atlas/atlas-search/) after crawling is done, the crawler takes care of that for you too.
- Logs will print in your terminal when you run the program.
