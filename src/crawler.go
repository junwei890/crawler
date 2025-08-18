package src

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/junwei890/crawler/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartCrawl(dbURI string, links []string) error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		return err
	}

	defer client.Disconnect(context.TODO())

	// doesn't actually get created till something is inserted
	db := client.Database("crawler")
	collection := db.Collection("content")

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
			}
		}(link, collection)
	}
	wg.Wait()

	// check if anything was inserted into the collection before indexing
	names, err := client.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	if ok := slices.Contains(names, "crawler"); !ok {
		return errors.New("no sites were crawled")
	}

	// create index if it doesn't exist, update it if it does
	indexName := "search_index"
	opts := options.SearchIndexes().SetName(indexName).SetType("search")

	cursor, err := collection.SearchIndexes().List(context.TODO(), opts)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	exists := false
	for cursor.Next(context.TODO()) {
		var indexMap bson.M
		if err := cursor.Decode(&indexMap); err != nil {
			return err
		}

		if indexMap["name"] == indexName {
			exists = true
		}
	}

	searchIndexModel := mongo.SearchIndexModel{
		Definition: bson.D{
			{Key: "mappings", Value: bson.D{
				{Key: "dynamic", Value: false},
				{Key: "fields", Value: bson.D{
					{Key: "content", Value: bson.D{
						{Key: "type", Value: "string"},
					}},
				}},
			}},
		},
		Options: opts,
	}

	if exists {
		if err := collection.SearchIndexes().UpdateOne(context.TODO(), indexName, searchIndexModel.Definition); err != nil {
			return err
		}
	} else {
		if _, err := collection.SearchIndexes().CreateOne(context.TODO(), searchIndexModel); err != nil {
			return err
		}
	}

	return nil
}

type Content struct {
	URL     string `bson:"_id"`
	Title   string `bson:"title"`
	Content string `bson:"content"`
}

func crawler(startURL string, collection *mongo.Collection) error {
	// get and parse robots.txt file first
	file, err := utils.GetRobots(startURL)
	if err != nil {
		return fmt.Errorf("didn't crawl %s: %v", startURL, err)
	}

	normURL, err := utils.Normalize(startURL)
	if err != nil {
		return fmt.Errorf("didn't crawl %s: %v", startURL, err)
	}

	rules, err := utils.ParseRobots(normURL, file)
	if err != nil {
		return fmt.Errorf("didn't crawl %s: %v", startURL, err)
	}

	dom, err := url.Parse(startURL)
	if err != nil {
		return fmt.Errorf("didn't crawl %s: %v", startURL, err)
	}

	visited := map[string]struct{}{}
	queue := &utils.Queue{}
	content := []any{}

	// if id already exists in the collection don't error and continue inserting
	options := options.InsertMany().SetOrdered(false)

	queue.Enqueue(startURL)

	re, err := regexp.Compile(`[^a-zA-Z0-9 ]+`)
	if err != nil {
		return fmt.Errorf("didn't crawl %s: %v", startURL, err)
	}

	for {
		// early returns
		if comp := queue.CheckEmpty(); comp {
			break
		}

		popped, err := queue.Dequeue()
		if err != nil {
			return err
		}

		ok, err := utils.CheckDomain(dom, popped)
		if err != nil {
			log.Println(fmt.Errorf("didn't crawl %s: %v", popped, err).Error())
			continue
		}
		if !ok {
			continue
		}

		currURL, err := utils.Normalize(popped)
		if err != nil {
			log.Println(fmt.Errorf("didn't crawl %s: %v", popped, err).Error())
			continue
		}

		ok = utils.CheckAbility(visited, rules, currURL)
		if !ok {
			continue
		}

		// sleeping for crawl delay right before get request
		subWg := &sync.WaitGroup{}
		subWg.Add(1)
		go func() {
			defer subWg.Done()

			time.Sleep(time.Duration(rules.Delay) * time.Second)
		}()

		page, err := utils.GetHTML(popped)
		if err != nil {
			log.Println(fmt.Errorf("didn't crawl %s: %v", popped, err).Error())
			continue
		}

		res, err := utils.ParseHTML(dom, page)
		if err != nil {
			log.Println(fmt.Errorf("didn't crawl %s: %v", popped, err).Error())
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

		content = append(content, Content{
			URL:     popped,
			Title:   res.Title,
			Content: cleaned,
		})

		// wait for sleep to finish before proceeding
		subWg.Wait()
	}

	if _, err := collection.InsertMany(context.TODO(), content, options); err != nil {
		return err
	}

	return nil
}
