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

Head into the project directory and create a `.env` file and a `crawler.txt` file
```
cd crawler
touch .env
touch crawler.txt
```

In the `.env` file, enter your MongoDB cluster connection string. There is **no need** to create a database or collection before running the crawler, the crawler takes care of that for you. The crawler will also create an Atlas search index after it's done crawling. In the event that you run the crawler multiple times, logic has been written to build on top of the database and index that was created when you first ran the crawler.

In the `crawler.txt` file, input any websites you would like to crawl, making sure each website is on a newline.

Once ready, run either of the following:
```
go build && ./crawler
```

## Inner workings
The following description is quite simplified, check out the source code to see how I implemented certain parts.

When you run the crawler, it reads both the MongoDB URI from the `.env` file and the links from `crawler.txt`. It then starts crawling each website in a separate Goroutine, up to **1000** at any point of time.

For each website, a request is made for its `robots.txt` file. After several checks are done, the `robots.txt` file is then parsed and crawling rules such as **allowed** routes, **disallowed** routes and **crawl delay** are extracted.

The website is then put in a queue and instantly popped, several checks are done to ensure that this is not a site we've visited before or this is a site that is on a disallowed route. Once validated, a request is made for the HTML of the current site and a parser extracts the **title**, **content** and **outgoing links**.

The links are put in a queue and the entire cycle is repeated till the queue is empty. The extracted title and content are stored in a struct and bulk inserted into the MongoDB database once the entire site is crawled.

Once the crawler is done with all the websites in the `crawler.txt` file, it creates/updates the Atlas Search Index for the entire corpus.

## Extensions
I don't have much more planned for this project except for maybe a UI to ease installation. If you think any part of the crawler could use some improvements, feel free to create a pull request.
