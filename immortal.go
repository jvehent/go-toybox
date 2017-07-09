package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	log.Println("ohai, my pid is", os.Getpid())
	go func() {
		for tries := 0; tries < 2; tries++ {
			log.Println("goroutine says hello")
			time.Sleep(3 * time.Second)
		}
		log.Println("goroutine starts new process")
		cmd := exec.Command(os.Args[0])
		_ = cmd.Start()
		os.Exit(0)
	}()
	time.Sleep(10 * time.Second)
	log.Println("main process says bye")
}
