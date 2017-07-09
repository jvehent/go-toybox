package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/codingsince1985/geo-golang/google"
	strava "github.com/strava/go.strava"
)

var access_token string = "***************"
var club_id int64 = 2298
var google_api string = "*******************"

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
		aLocation := ""
		geocoder := google.Geocoder(google_api)
		address, _ := geocoder.ReverseGeocode(activity.StartLocation[0], activity.StartLocation[1])
		if address != "" {
			addressComp := strings.Split(address, ",")
			if len(addressComp) > 3 {
				aLocation = " around" + strings.Join(addressComp[len(addressComp)-3:], ",")
			} else {
				aLocation = " around " + address
			}
		}
		fmt.Printf("%s %s %s %0.1fkm in %s%s: %s\n",
			activity.Athlete.FirstName, activity.Athlete.LastName,
			aType,
			aDistance,
			aDuration,
			aLocation,
			activity.Name)
	}
}
