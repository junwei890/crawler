package parsers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Normalize(rawURL string) (string, error) {
	structure, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("normalize: error parsing %s, %v", rawURL, err)
	}

	return structure.Host + strings.TrimRight(structure.Path, "/"), nil
}

func GetHTML(rawURL string) ([]byte, error) {
	client := &http.Client{}

	res, err := client.Get(rawURL)
	if err != nil {
		return []byte{}, fmt.Errorf("gethtml: error from get request to %s, %v", rawURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return []byte{}, fmt.Errorf("gethtml: %d status code returned from %s", res.StatusCode, rawURL)
	}

	mediaType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return []byte{}, fmt.Errorf("gethtml: error parsing content type for %s, %v", rawURL, err)
	}
	if mediaType != "text/html" {
		return []byte{}, fmt.Errorf("gethtml: content type not text/html for %s", rawURL)
	}

	page, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("gethtml: error reading response body, %v", err)
	}

	return page, nil
}

type Response struct {
	Title   string
	Content []string
	Links   []string
}

func ParseHTML(domain *url.URL, page []byte) (Response, error) {
	response := Response{}
	skip := true
	title := false

	tokens := html.NewTokenizer(bytes.NewReader(page))
	for {
		tn := tokens.Next()

		if tn == html.ErrorToken {
			if tokens.Err() == io.EOF {
				break
			}

			return response, fmt.Errorf("parsehtml: error tokenising html, %v", tokens.Err())
		}

		if tn == html.TextToken {
			t := tokens.Token()

			if title {
				response.Title = strings.Join(strings.Fields(t.Data), " ")
				continue
			}

			if skip {
				continue
			}

			clean := strings.ToLower(strings.Join(strings.Fields(t.Data), " "))
			if clean != "" {
				response.Content = append(response.Content, clean)
			}
			continue
		}

		if tn == html.StartTagToken {
			t := tokens.Token()

			if t.Data == "p" && t.DataAtom == atom.P {
				skip = false
				continue
			}

			if t.Data == "title" && t.DataAtom == atom.Title {
				title = true
				continue
			}

			if t.Data == "a" && t.DataAtom == atom.A {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						structure, err := url.Parse(attr.Val)
						if err != nil {
							log.Println(fmt.Errorf("parsehtml: error parsing %s, %v", attr.Val, err))
							continue
						}

						if structure.Hostname() == "" {
							fullURL := domain.ResolveReference(structure).String()
							if comp := slices.Contains(response.Links, fullURL); !comp {
								response.Links = append(response.Links, fullURL)
							}
						} else {
							if comp := slices.Contains(response.Links, attr.Val); !comp {
								response.Links = append(response.Links, attr.Val)
							}
						}
					}
				}
			}
		}

		if tn == html.EndTagToken {
			t := tokens.Token()

			if t.Data == "p" && t.DataAtom == atom.P {
				skip = true
				continue
			}

			if t.Data == "title" && t.DataAtom == atom.Title {
				title = false
				continue
			}
		}
	}

	return response, nil
}

func GetRobots(rawURL string) ([]byte, error) {
	client := &http.Client{}
	res, err := client.Get(fmt.Sprintf("%srobots.txt", rawURL))
	if err != nil {
		return []byte{}, fmt.Errorf("getrobots: error from get request to %s, %v", rawURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode == 403 {
		return []byte{}, fmt.Errorf("getrobots: %d status code returned from %s", res.StatusCode, rawURL)
	}
	if res.StatusCode == 404 {
		return []byte{}, nil
	}

	mediaType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return []byte{}, fmt.Errorf("getrobots: error parsing content type for %s, %v", rawURL, err)
	}
	if mediaType != "text/plain" {
		return []byte{}, fmt.Errorf("getrobots: robots.txt content type not text/plain for %s", rawURL)
	}

	textFile, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("getrobots: error reading response body, %v", err)
	}

	return textFile, nil
}

type Rules struct {
	Allowed    []string
	Disallowed []string
	Delay      int
}

func ParseRobots(normURL string, textFile []byte) (Rules, error) {
	rules := Rules{}
	applicable := false

	scanner := bufio.NewScanner(bytes.NewReader(textFile))
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "" || strings.HasPrefix(strings.TrimSpace(scanner.Text()), "#") {
			continue
		}

		line := strings.Split(scanner.Text(), ":")
		key := strings.TrimSpace(line[0])
		value := strings.TrimSpace(line[1])

		if key == "User-agent" {
			if value == "*" {
				applicable = true
			} else {
				applicable = false
			}
		}

		if applicable {
			switch key {
			case "Allow":
				if strings.HasPrefix(value, "/") {
					rules.Allowed = append(rules.Allowed, fmt.Sprintf("%s%s", normURL, value))
				}
			case "Disallow":
				if strings.HasPrefix(value, "/") {
					rules.Disallowed = append(rules.Disallowed, fmt.Sprintf("%s%s", normURL, value))
				}
			case "Crawl-delay":
				delay, err := strconv.Atoi(value)
				if err != nil {
					log.Println(fmt.Errorf("parserobots: error parsing crawl delay %s, %v", value, err))
					rules.Delay = 0
				} else {
					rules.Delay = delay
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(fmt.Errorf("parserobots: error reading line, %v", err))
	}

	return rules, nil
}
