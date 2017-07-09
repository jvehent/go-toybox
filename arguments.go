package main

import (
	"flag"
	"log"

	"bitbucket.org/kardianos/osext"
)

func main() {
	var debug = flag.Bool("d", false, "debug mode: do not fork and print all logs to stdout")
	flag.Parse()
	log.Println(osext.Executable())
	if *debug {
		log.Println("debug is set")
	}
}
