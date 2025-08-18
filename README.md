# crawler
This is a web crawler I wrote for fun, it's fast, polite and stack safe. If you find a web crawler useful, feel free to follow the installation instructions below.

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

![image1](./images/success.png)

A nil error return will show a success message.

![image2](./images/failure.png)

A non-nil error return will show a error message

## Notes
Some sites enforce allowed and disallowed routes as well as crawl delays, this crawler abides by all of them.

Feel free to fork the repo to make changes that fit your personal needs.

Contributions are welcomed.
