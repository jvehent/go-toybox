
package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
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

	// We declare our topology on both the publisher and consumer to ensure they
	// are the same.  This is part of AMQP being a programmable messaging model.
	//
	// See the Channel.Consume example for the complimentary declare.
	err = c.ExchangeDeclare("migexchange",	// exchange name
				"topic",	// exchange type
				true,		// is durable
				false,		// is autodelete
				false,		// is internal
				false,		// is noWait
				nil)		// optional arguments
	if err != nil { panic(err) }

	// Prepare this message to be persistent.  Your publishing requirements may
	// be different.
	msg := amqp.Publishing{
	    DeliveryMode: amqp.Persistent,
	    Timestamp:    time.Now(),
	    ContentType:  "text/plain",
	    Body:         []byte("Go Go AMQP!"),
	}
	fmt.Println(msg)
	// This is not a mandatory delivery, so it will be dropped if there are no
	// queues bound to the logs exchange.
	err = c.Publish("migexchange",			// exchange name
			"mig.scheduler.action.run",	// exchange key
			false,				// is mandatory
			false,				// is immediate
			msg)				// AMQP message
	if err != nil { panic(err) }
}
