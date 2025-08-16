package src

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

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
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			fmt.Println(err)
		}
	}()

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
				fmt.Println(err)
			}
		}(link, collection)
	}
	wg.Wait()

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
	file, err := utils.GetRobots(startURL)
	if err != nil {
		return err
	}

	normURL, err := utils.Normalize(startURL)
	if err != nil {
		return err
	}

	rules, err := utils.ParseRobots(normURL, file)
	if err != nil {
		return err
	}

	dom, err := url.Parse(startURL)
	if err != nil {
		return err
	}

	visited := map[string]struct{}{}
	queue := &utils.Queue{}
	content := []any{}
	options := options.InsertMany().SetOrdered(false)

	queue.Enqueue(startURL)

	re, err := regexp.Compile(`[^a-zA-Z0-9 ]+`)
	if err != nil {
		return err
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
			continue
		}
		if !ok {
			continue
		}

		currURL, err := utils.Normalize(popped)
		if err != nil {
			continue
		}

		ok = utils.CheckAbility(visited, rules, currURL)
		if !ok {
			continue
		}

		page, err := utils.GetHTML(popped)
		if err != nil {
			continue
		}

		res, err := utils.ParseHTML(dom, page)
		if err != nil {
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
	}

	result, err := collection.InsertMany(context.TODO(), content, options)
	if err != nil {
		return err
	}
	if len(result.InsertedIDs) != len(content) {
		return errors.New("couldn't insert some content")
	}

	return nil
}
