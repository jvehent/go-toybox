package main

import (
	log "github.com/Sirupsen/logrus"
	"go.mozilla.org/mozlogrus"
)

func init() {
	mozlogrus.Enable("ApplicationName")
}

func main() {
	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
}
