# crawler
This is a web crawler I wrote that allows for large scale concurrency while being stack safe and polite.

## Requirements
- [Go](https://go.dev/doc/install) installed
- A [MongoDB](https://www.mongodb.com/docs/atlas/getting-started/) cluster

## Installation
Run the following in your terminal:
```
go install github.com/junwei890/crawler@latest
```

Now run:
```
crawler
```

This gives you access to the UI in any working directory.

## Usage
![image](./images/crawler_ui.png)

This is the main UI, in the first input field, enter your MongoDB URI. In the second input field, enter any sites you would like to scrape.

Make sure each site is on a newline and has its protocol.

## Inner workings
### Program entry
Once the MongoDB URI and sites have been entered, a database connection is established and a Goroutine is created to crawl each site, up to a **thousand**.

A channel is also shared across all Goroutines to share crawl counts and error logs in realtime to the user.

### Robots.txt
For each site, a GET request is made for its `robots.txt` file, this file outlines which routes a crawler **can and cannot access as well as the crawl delay** it should abide by.

Based on the response, one of several things could happen:
- 403: The site doesn't want us crawling so we don't.
- 404: There's no robots.txt file so we will be crawling the site.
- Malformed or no Content-Type headers: The site won't be crawled.

If all these pass, the file is passed through a parser where rules are extracted.

### Breadth First Traversal
A **breadth first traversal** was chosen over a recursive depth first one. This is because Go isn't tail call optimized, it allocates a new stack on each recursive call instead of reusing the previous one, thus using a depth first traversal could **potentially** crash our program if sites are massive.

This decision impaired performance since we couldn't crawl each route in a separate Goroutine, however, it gave us much better **stack safety** since the stack grows only with queue size.

### Early returns
The crawling takes place in an infinite for loop until the queue is empty. Before getting and parsing HTML, several checks are done:
- Checks if queue is empty.
- Checks if we are still within the same hostname.
- Checks if we have already visited this route or if are even allowed to visit this route.

If any of the above is satisfied, we head to the next for loop iteration.

### HTML
Once a route makes it through early returns, a GET request is made for the route's HTML, if the route responds with a **400 to 499 status code** or if the Content-Type in the response header is not **text/html**, we skip over to the next for loop iteration.

The retrieved HTML is then passed through a parser that extracts the title, content and outgoing links. The title and content are unmarshalled into a struct and temporarily stored in a slice while the links are enqueued.

### Post-crawling
Once each site exits the for loop, titles and content we extracted are **bulk inserted** into MongoDB, with the database and collection creation **automated**.

Once all sites have been crawled, the collection is then **automatically indexed** for Atlas Search.

The crawler builds on top of the database, collection and index that was created on the first successful run on subsequent program executions. All of this is handled by the crawler.

## Extensions
- [x] Viewport to see error logs.
- [ ] Site map crawling
