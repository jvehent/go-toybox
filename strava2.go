package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	strava "github.com/strava/go.strava"
)

var access_token string = "*************************"
var club_id int64 = 2298

func main() {
	client := strava.NewClient(access_token)
	if client == nil {
		log.Fatal("Failed to create strava client")
	}
	activities, err := strava.NewClubsService(client).
		ListActivities(club_id).
		PerPage(200).
		Do()
	if err != nil {
		log.Fatal(err)
	}
	for _, activity := range activities {
		data, _ := json.Marshal(activity)
		fmt.Printf("%s\n", data)
		var aType string
		switch activity.Type.String() {
		case "Run":
			aType = "ran"
		case "Ride":
			aType = "biked"
		case "Hike":
			aType = "hiked"
		case "Kayaking":
			aType = "kayaked"
		default:
			aType = activity.Type.String()
		}
		aDistance := activity.Distance / 1000
		aDuration, err := time.ParseDuration(fmt.Sprintf("%ds", activity.ElapsedTime))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s %s %0.1fkm in %s: %s\n",
			activity.Athlete.FirstName, activity.Athlete.LastName,
			aType,
			aDistance,
			aDuration,
			activity.Name)
	}
}
