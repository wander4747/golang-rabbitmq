package main

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {

	args := os.Args

	if len(args) == 1 {
		failOnError(errors.New("args"), "Failed to open a channel")
		return
	}

	if args[1] == "sender" {
		sender()
	}

	if args[1] == "receive" {
		receive()
	}
}

func sender() {
	connection, channel, err := connect()
	defer connection.Close()

	failOnError(err, "Failed to connect to RabbitMQ")

	defer channel.Close()

	queue, err := queueDeclare("hello", err, channel)
	failOnError(err, "Failed to declare a queue")

	for i := 0; i < 100000; i++ {
		body := fmt.Sprintf("message n° %v", i)

		err = channel.Publish(
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		failOnError(err, "Failed to publish a message")

		log.Printf("Send a message: %s", body)
	}
}

func receive() {
	connection, channel, err := connect()
	defer connection.Close()

	failOnError(err, "Failed to connect to RabbitMQ")

	defer channel.Close()
	queue, err := queueDeclare("hello", err, channel)
	failOnError(err, "Failed to declare a queue")

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func connect() (*amqp.Connection, *amqp.Channel, error) {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")

	return connection, channel, err
}

func queueDeclare(name string, err error, ch *amqp.Channel) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return queue, err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
