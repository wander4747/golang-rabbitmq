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
		failOnError(errors.New("args"), "insufficient arguments")
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

	err = channel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	for i := 0; i < 100000; i++ {
		body := fmt.Sprintf("message n° %v", i)
		err = channel.Publish(
			"logs", // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})

		failOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent %s", body)
	}
}

func receive() {
	connection, channel, err := connect()
	defer connection.Close()

	failOnError(err, "Failed to connect to RabbitMQ")

	defer channel.Close()

	err = channel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	queue, err := queueDeclare("logs", err, channel)
	failOnError(err, "Failed to declare a queue")

	err = channel.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		"logs",     // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

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
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
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
