// watch a set of directories for file system events
package main

import (
	"log"

	"github.com/howeyc/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// goroutine receives events and prints them
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	// list of directories to watch events on
	err = watcher.Watch("/var/cache/mig/command/ready/")
	err = watcher.Watch("/var/cache/mig/command/inflight/")
	err = watcher.Watch("/var/cache/mig/command/returned/")
	err = watcher.Watch("/var/cache/mig/command/done/")
	if err != nil {
		log.Fatal(err)
	}

	// block while watching for events
	<-done

	watcher.Close()
}
