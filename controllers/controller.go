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

func CheckForNewBlogPosts() {
	ctx := context.Background()
	latestTitle, err := client.Get(ctx, latestBlogPostTitleKey).Result()

	// cek redis error atau tidak
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching latest blog post title from Redis: %v\n", err)
		return
	} else {
		fmt.Println("aman1")
	}

	// ambil blog terakhir di database
	latestBlogPost, err := FetchLatestBlogPostFromDB(ctx)
	if err != nil {
		log.Printf("Error fetching latest blog post from database: %v\n", err)
		return
	} else {
		fmt.Println("aman2")
	}

	// cek kalo blog trakhir tidak sama brti ada blog baru
	if latestBlogPost.Title != latestTitle {
		log.Printf("New blog post found: %v\n", latestBlogPost.Title)

		// update title di redis
		err = client.Set(ctx, latestBlogPostTitleKey, latestBlogPost.Title, 0).Err()
		if err != nil {
			log.Printf("Error setting latest blog post title in Redis: %v\n", err)
			return
		}

		// send email ke subscriber
		go sendEmailNotification(latestBlogPost) // pake goroutine untuk non-blocking
	} else {
		log.Println("No new blog post found.")
	}
}

func FetchLatestBlogPostFromDB(ctx context.Context) (*models.Blog, error) {
	var latestBlogPost models.Blog

	db := connect()
	fmt.Println("udah siap setup db")
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
	fmt.Println("send email")
}
