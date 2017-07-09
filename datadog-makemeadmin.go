package main

import (
	"fmt"

	"github.com/zorkian/go-datadog-api"
)

func main() {
	cli := datadog.NewClient("SOMEAPIKEY", "SOMEAPPKEY")
	user, err := cli.GetUser("jvehent@mozilla.com")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", user)
	user.IsAdmin = true
	err = cli.UpdateUser(user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user.Handle, "is now datadog admin")
}
