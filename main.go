package main

import (
	"log"
	"tugas-explorasi3/controllers"

	"github.com/jasonlvhit/gocron"
)

func main() {

	controllers.SetupRedis()

	// Schedule the job with GoCRON
	gocron.Every(2).Minute().Do(func() {
		err := controllers.checkForNewBlogPosts()
		if err != nil {
			log.Println("Error checking for new blog posts:", err)
			return
		}
	})

	// Start the cron scheduler
	<-gocron.Start()
}
