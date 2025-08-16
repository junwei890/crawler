# Web crawler
This is a web crawler I wrote to gather links to articles, research papers and books regarding topics I find interesting. Though I wrote it for personal use, it slowly evolved into a piece of software that was kind of portable. If you find a web crawler useful, you can read the installation instructions below.

## Requirements
- [Go](https://go.dev/doc/install) installed.
- A [MongoDB Atlas](https://www.mongodb.com/docs/atlas/getting-started/) cluster.

## Installation
Run `go install github.com/junwei890/crawler@latest` in your terminal.

Now you can run the command `crawler` in your terminal to bring up the UI.
![image](images/crawler_ui.png)

## Tips
- You don't have to create a MongoDB database and collection before scraping websites, the crawler takes care of it for you.
- You don't have to index the collection for [Atlas Search](https://www.mongodb.com/docs/atlas/atlas-search/) after scraping is done, the crawler takes care of it for you too.
- Make sure links you input are separated by a newline and they have their protocol (ideally https://).
- Running the crawler on an empty cluster will always be faster than running it when there is already a collection and index because creating an index is faster than reindexing.
- Some websites enforce long crawl delays, be patient.
- I've set a hard limit of 1000 concurrent crawlers, just because you can scrape 1000 websites at a time doesn't mean you should, consult your hardware first.

## Possible extensions
- [ ] Panel to see logs while scraping.
