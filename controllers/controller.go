package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/smtp"
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
		//go sendEmailNotification(latestBlogPost) // pake goroutine untuk non-blocking
		go sendEmail()
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

func sendEmail() {
	// Sender's email address and password.
	from := "raphaelevaan1@gmail.com"
	password := "sdhe qiez hlzf kgex"

	// SMTP server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	ctx := context.Background()
	latestTitle, err := client.Get(ctx, latestBlogPostTitleKey).Result()
	// cek redis error atau tidak
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching latest blog post title from Redis: %v\n", err)
		return
	}

	subject := "Jovi Baru Saja Mempublish Blog Baru"
	body := latestTitle

	to := "j0njovi0jjh2710@gmail.com"

	// Message to send.
	message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	fmt.Println("Email Sudah Terkirim!")
}
