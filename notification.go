package main

import (
	"github.com/0xAX/notificator"
)

var notify *notificator.Notificator

func main() {

	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "My test App",
	})

	notify.Push("testnotification", "hello world", "/home/ulfr/DSC_7871.jpg", notificator.UR_CRITICAL)
}
