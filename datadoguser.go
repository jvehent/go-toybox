package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	client := datadog.NewClient(os.Getenv("DATADOG_API_KEY"), os.Getenv("DATADOG_APP_KEY"))

	user, err := client.GetUser("vintagecarigou@mozilla.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", user)

	user.Name = "Bob Kelso"
	user.Disabled = false
	user.IsAdmin = true
	user.Role = "ldapmanaged"
	err = client.UpdateUser(user)
	if err != nil {
		log.Fatal(err)
	}

	user, err = client.GetUser(user.Handle)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", user)

	user.Disabled = true
	user.IsAdmin = false
	err = client.UpdateUser(user)
	if err != nil {
		log.Fatal(err)
	}
	user, err = client.GetUser(user.Handle)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", user)

}
