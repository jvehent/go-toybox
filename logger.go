package main

import (
	"log"
	"log/syslog"
)

func main() {
	log.Println(syslog.LOG_EMERG)
}
