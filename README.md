# crawler
This is a web crawler I wrote for fun, it's fast, polite and stack safe but a little inconvenient to set up. If you find a web crawler useful, feel free to follow the installation instructions below.

## Requirements
- [Go](https://go.dev/doc/install) installed
- A [MongoDB](https://www.mongodb.com/docs/atlas/getting-started/) cluster

## Setup
Clone the repository into your working directory using:
```
git clone https://github.com/junwei890/crawler.git
```

Head into the project directory and create a .env file and crawler.txt
```
cd crawler
touch .env
touch crawler.txt
```

In the .env file, enter your MongoDB cluster connection string. There is no need to create a database or collection before running the crawler, the crawler takes care of that for you. The crawler will also create an Atlas search index after it's done crawling. In the event that you run the crawler multiple times, logic has been written to build on top of the database and index that was created when you first ran the crawler.

In the crawler.txt file, input any websites you would like to crawl, making sure each website is on a newline.

Once ready, run either of the following:
```
go build && ./crawler

go run .
```

**NOTE** ~ Some websites enforce long crawl delays and this crawler abides by them.

## Planned extensions
- [ ] UI
