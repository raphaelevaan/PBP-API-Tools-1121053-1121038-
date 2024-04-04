package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tugas-explorasi3/models"

	"github.com/go-redis/redis/v8"
)

const latestBlogPostTitleKey = "latest_blog_post_title"

func checkForNewBlogPosts() {
	ctx := context.Background()
	latestTitle, err := client.Get(ctx, latestBlogPostTitleKey).Result()

	// If there is an error fetching from Redis (other than Redis 'nil', which means the key doesn't exist), log and return.
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching latest blog post title from Redis: %v\n", err)
		return
	}

	// Fetch the latest blog post from the database.
	latestBlogPost, err := fetchLatestBlogPostFromDB(ctx)
	if err != nil {
		log.Printf("Error fetching latest blog post from database: %v\n", err)
		return
	}

	// Check if there is a new post.
	if latestBlogPost.Title != latestTitle {
		log.Printf("New blog post found: %v\n", latestBlogPost.Title)

		// Update the title in Redis.
		err = client.Set(ctx, latestBlogPostTitleKey, latestBlogPost.Title, 0).Err()
		if err != nil {
			log.Printf("Error setting latest blog post title in Redis: %v\n", err)
			return
		}

		// Send an email notification about the new post.
		go sendEmailNotification(latestBlogPost) // using a goroutine for non-blocking
	} else {
		log.Println("No new blog post found.")
	}
}

func fetchLatestBlogPostFromDB(ctx context.Context) (*models.Blog, error) {
	var latestBlogPost models.Blog

	db := SetupDatabase()

	// Assuming you have a table 'blog_posts' with fields 'id', 'title', 'content', etc.
	query := "SELECT id, title, content FROM blog_posts ORDER BY id DESC LIMIT 1"
	row := db.QueryRowContext(ctx, query)
	err := row.Scan(&latestBlogPost.ID, &latestBlogPost.Title, &latestBlogPost.Content)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no blog posts found")
	} else if err != nil {
		return nil, err
	}

	return &latestBlogPost, nil
}

func sendEmailNotification(post *models.Blog) {

}
