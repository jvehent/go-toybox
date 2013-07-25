
package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	// Connects opens an AMQP connection from the credentials in the URL.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil { panic(err) }

	// This waits for a server acknowledgment which means the sockets will have
	// flushed all outbound publishings prior to returning.  It's important to
	// block on Close to not lose any publishings.
	defer conn.Close()

	c, err := conn.Channel()
	if err != nil { panic(err) }

	q, err := c.QueueDeclare("mig.scheduler",
				true,
				false,
				false,
				false,
				nil)
	fmt.Println(q)
	err = c.QueueBind("mig.scheduler",
			"mig.scheduler.action.run",
			"migexchange",
			false,
			nil)
	if err != nil { panic(err) }

	msg, err := c.Consume("mig.scheduler",
			"mig.scheduler.action.run",
			false,
			false,
			false,
			false,
			nil)
	if err != nil { panic(err) }
	for m := range msg {
		fmt.Printf("%s\n", m.Body)
	if err != nil { panic(err) }
	}
}
