
package main

import (
	"os/exec"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"encoding/json"
	"time"
)

type Action struct{
	Name, Target, Check, Command string
}

type FileCheckerResult struct {
	TestedFiles, ResultCount int
	Files	[]string
}

type Alert struct {
	IOC, Item string
}

func getActions(messages <- chan amqp.Delivery, actions chan []byte,
	     terminate chan bool) error {
	// range waits on the channel and returns all incoming messages
	// range will exit when the channel closes
	for m := range messages {
		log.Printf("getActions: received '%s'", m.Body)
		// Ack this message only
		err := m.Ack(true)
		if err != nil { panic(err) }
		actions <- m.Body
		log.Printf("getActions: queued in pos. %d", len(actions))
	}
	terminate <- true
	return nil
}

func parseActions(actions <- chan []byte, fCommand chan string,
		  terminate chan bool) error {
	var action Action
	for a := range actions {
		err := json.Unmarshal(a, &action)
		if err != nil { panic(err) }
		log.Printf("ParseAction: Name '%s' Target '%s' Check '%s' Command '%s'",
			   action.Name, action.Target, action.Check, action.Command)
		switch action.Check{
		case "filechecker":
			fCommand <- action.Command
			log.Printf("parseActions: queued into filechecker in pos. %d", len(fCommand))
		}
	}
	terminate <- true
	return nil
}

func runFilechecker(fCommand <- chan string, alert chan Alert,
		    terminate chan bool) error {
	for c := range fCommand {
		log.Printf("RunFilechecker: command '%s' is being executed", c)
		cmd := exec.Command("./filechecker", c)
		cmdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		st := time.Now()
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		results := make(map[string]FileCheckerResult)
		err = json.NewDecoder(cmdout).Decode(&results)
		if err != nil {
			log.Fatal(err)
		}
		cmdDone := make(chan error)
		go func() {
			cmdDone <-cmd.Wait()
		}()
		select {
		// kill the process when timeout expires
		case <-time.After(30 * time.Second):
			if err := cmd.Process.Kill(); err != nil {
				log.Fatal("failed to kill:", err)
			}
			log.Fatal("runFileChecker: command '%s' timed out", c)
		// exit normally
		case err := <-cmdDone:
			if err != nil {
				log.Fatal(err)
			}
		}
		for _, r := range results {
			log.Println("runFileChecker: command", c,"tested",
				    r.TestedFiles, "files in", time.Now().Sub(st))
			if r.ResultCount > 0 {
				for _, f := range r.Files {
					alert <- Alert{
						IOC: c,
						Item: f,
					}
				}
			}
		}
	}
	terminate <- true
	return nil
}

func raiseAlerts(alert chan Alert, terminate chan bool) error {
	for a := range alert {
		log.Printf("raiseAlerts: IOC '%s' positive match on '%s'",
			   a.IOC, a.Item)
	}
	return nil
}

func main() {
	termChan	:= make(chan bool)
	actionsChan	:= make(chan []byte, 10)
	fCommand	:= make(chan string, 10)
	alert		:= make(chan Alert, 10)
	// Connects opens an AMQP connection from the credentials in the URL.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil { panic(err) }
	defer conn.Close()

	c, err := conn.Channel()
	if err != nil { panic(err) }

	// declare a queue
	q, err := c.QueueDeclare("mig.action",	// queue name
				true,		// is durable
				false,		// is autoDelete
				false,		// is exclusive
				false,		// is noWait
				nil)		// AMQP args
	fmt.Println(q)

	// bind a queue to an exchange via the key
	err = c.QueueBind("mig.action",		// queue name
			"mig.action.create",	// exchange key
			"migexchange",		// exchange name
			false,			// is noWait
			nil)			// AMQP args
	if err != nil { panic(err) }

	// Limit the number of message the channel will receive
	err = c.Qos(1,		// prefetch count (in # of msg)
		    0,		// prefetch size (in bytes)
		    false)	// is global
	if err != nil { panic(err) }

	// Initialize a consumer than pulls messages into a channel
	tag := fmt.Sprintf("%s", time.Now())
	msgChan, err := c.Consume("mig.action",	// queue name
			tag,			// exchange key
			false,			// is autoAck
			false,			// is exclusive
			false,			// is noLocal
			false,			// is noWait
			nil)			// AMQP args
	if err != nil { panic(err) }

	// This goroutine will continously pull messages from the consumer
	// channel, print them to stdout and acknowledge them
	go getActions(msgChan, actionsChan, termChan)
	go parseActions(actionsChan, fCommand, termChan)
	go runFilechecker(fCommand, alert, termChan)
	go raiseAlerts(alert, termChan)
	// block until terminate chan is called
	<- termChan
}
