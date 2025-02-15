package main

import (
	"context"
	"database/sql"
	"github.com/anqiwoo/rssagg/internal/database"
	"github.com/google/uuid"
	"log"
	"strings"
	"sync"
	"time"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Starting scraping with %d workers every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error getting feeds to fetch:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %s", err)
		return
	}

	log.Printf("Scraping feed %s", feed.Url)
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error scraping feed %s: %s", feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing published date %v: %s", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: publishedAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			// Skip expected errors for parsed posts
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Error creating post: %s", err)
		}
	}
	log.Printf("Finished scraping feed %s, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
