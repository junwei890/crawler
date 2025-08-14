package src

import (
	"context"
	"errors"
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init(collection *mongo.Collection) error {
	file, err := os.ReadFile("links.txt")
	if err != nil {
		return errors.New("coudn't read links.txt file")
	}
	links := strings.Fields(string(file))

	wg := &sync.WaitGroup{}
	channel := make(chan struct{}, 1000)

	for _, link := range links {
		wg.Add(1)
		channel <- struct{}{}

		go func(link string, db *mongo.Collection) {
			defer func() {
				<-channel
				wg.Done()
			}()

			if err := crawler(link, collection); err != nil {
				log.Println(err)
				return
			}
		}(link, collection)
	}
	wg.Wait()

	return nil
}

type Content struct {
	URL     string `bson:"_id"`
	Title   string `bson:"title"`
	Content string `bson:"content"`
}

func crawler(startURL string, collection *mongo.Collection) error {
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
	content := []any{}
	options := options.InsertMany().SetOrdered(false)

	queue.Enqueue(startURL)

	pattern := `[^a-zA-Z0-9 ]+`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("can't compile %s", pattern)
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

		log.Printf("crawled %s", popped)

		slice := []string{}
		for _, content := range res.Content {
			slice = append(slice, re.ReplaceAllString(content, ""))
		}

		cleaned := strings.Join(slice, " ")
		if len(cleaned) < 500 {
			continue
		}

		content = append(content, Content{
			URL:     popped,
			Title:   res.Title,
			Content: cleaned,
		})
	}

	result, _ := collection.InsertMany(context.TODO(), content, options)
	if len(result.InsertedIDs) != len(content) {
		return errors.New("couldn't insert some content")
	}

	return nil
}
