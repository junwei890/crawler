package src

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/junwei890/se-cli/data_structures"
	"github.com/junwei890/se-cli/parsers"
	"github.com/junwei890/se-cli/utils"
)

func Init() error {
	file, err := os.ReadFile("links.txt")
	if err != nil {
		return fmt.Errorf("init: error reading links.txt file, %v", err)
	}
	links := strings.Fields(string(file))

	wg := &sync.WaitGroup{}
	channel := make(chan struct{}, 1000)

	for _, link := range links {
		wg.Add(1)
		channel <- struct{}{}

		go func(link string) {
			defer func() {
				<-channel
				wg.Done()
			}()

			if err := crawler(link); err != nil {
				log.Println(err)
				return
			}
		}(link)
	}
	wg.Wait()

	return nil
}

func crawler(startURL string) error {
	file, err := parsers.GetRobots(startURL)
	if err != nil {
		return err
	}

	normURL, err := parsers.Normalize(startURL)
	if err != nil {
		return err
	}

	rules, err := parsers.ParseRobots(normURL, file)
	if err != nil {
		return err
	}

	dom, err := url.Parse(startURL)
	if err != nil {
		return err
	}

	visited := map[string]struct{}{}
	queue := &data_structures.Queue{}
	queue.Enqueue(startURL)

	re, err := regexp.Compile(`[^A-Za-z ]+`)
	if err != nil {
		return fmt.Errorf("crawler: error compiling regex pattern, %v", err)
	}

	for {
		if comp := queue.CheckEmpty(); comp {
			break
		}

		popped, err := queue.Dequeue()
		if err != nil {
			return err
		}

		ok, err := utils.CheckDomain(dom, popped)
		if err != nil {
			log.Println(err)
			continue
		}
		if !ok {
			continue
		}

		currURL, err := parsers.Normalize(popped)
		if err != nil {
			log.Println(err)
			continue
		}

		ok = utils.CheckAbility(visited, rules, currURL)
		if !ok {
			continue
		}

		log.Printf("crawling %s", popped)

		page, err := parsers.GetHTML(popped)
		if err != nil {
			log.Println(err)
			continue
		}

		res, err := parsers.ParseHTML(dom, page)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, link := range res.Links {
			queue.Enqueue(link)
		}

		slice := []string{}
		for _, content := range res.Content {
			slice = append(slice, re.ReplaceAllString(content, ""))
		}

		cleaned := strings.Join(slice, " ")
		if len(cleaned) < 500 {
			continue
		}

		log.Println(cleaned)
	}

	return nil
}
