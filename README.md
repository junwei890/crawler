# Web crawler
This is a web crawler I wrote to gather links to articles, research papers and books regarding topics I find interesting. Though I wrote it for personal use, it slowly evolved into a piece of software that was kind of portable. If you find a web crawler useful, you can read the installation instructions below.

## Requirements
- [Go](https://go.dev/doc/install) installed.
- A [MongoDB Atlas](https://www.mongodb.com/docs/atlas/getting-started/) cluster.

## Installation
Run `go install github.com/junwei890/crawler@latest` in your terminal.

## Usage
You need to have a crawler.txt file in your home directory, it should contain your MongoDB URI and links you would like to scrape. See [crawler.txt.example](https://github.com/junwei890/crawler/blob/main/crawler.txt.example) for how it should be formatted.

Run `crawler` in your terminal to start scraping.

### Tips
- You don't have to create a MongoDB database and collection before scraping websites, the crawler takes care of that for you.
- You don't have to index the collection for [Atlas Search](https://www.mongodb.com/docs/atlas/atlas-search/) after scraping is done, the crawler takes care of that for you too.
- Logs will print in your terminal when you run the program.
