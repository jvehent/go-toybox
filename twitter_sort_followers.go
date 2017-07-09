package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("TWITTER_ACCESS_SECRET"))
	pages := api.GetFollowersListAll(nil)
	for page := range pages {
		if page.Error != nil {
			log.Fatal(page.Error)
		}
		for _, user := range page.Followers {
			fmt.Println(user.FollowersCount, user.Name)
		}
	}
}
