package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to rabbitmq
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		log.Println("Error connecting to rabbitmq:", err)
		os.Exit(1)
	}
	defer conn.Close()

	log.Println("Listening for and consuming RabbitMQ messages...")

	// create channel
	channel, err := conn.Channel()
	if err != nil {
		log.Println("Error opening channel:", err)
		os.Exit(1)
	}
	defer channel.Close()

	// create exchange
	err = channel.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // kind
		true,         // durable
		false,        // autoDelete
		false,        // internal
		false,        // noWait
		nil,          // args
	)
	if err != nil {
		log.Println("Error declaring exchange:", err)
		os.Exit(1)
	}

	// declare queue
	queue, err := channel.QueueDeclare(
		"",    // name - empty for auto-generated name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Println("Error declaring queue:", err)
		os.Exit(1)
	}

	// bind queue to exchange
	err = channel.QueueBind(
		queue.Name,       // queue name
		"#",          // routing key - # means all messages
		"logs_topic", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Println("Error binding queue:", err)
		os.Exit(1)
	}

	// create consumer
	msgs, err := channel.Consume(
		queue.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println("Error consuming messages:", err)
		os.Exit(1)
	}

	// consume messages
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Println("Listener started")
	<-forever
	os.Exit(0)
}
