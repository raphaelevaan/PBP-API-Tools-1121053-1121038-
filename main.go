package main

import (
	"tugas-explorasi3/controllers"

	"github.com/jasonlvhit/gocron"
)

func main() {

	controllers.SetupRedis()

	gocron.Every(5).Second().Do(func() {
		controllers.CheckForNewBlogPosts()
	})

	<-gocron.Start()
}
