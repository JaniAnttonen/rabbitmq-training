package main

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func openRmqChannel(rmqHost string) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial("amqp://guest:guest@" + rmqHost)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch, conn
}

func declareQueue(channel *amqp.Channel, name string) *amqp.Queue {
	q, err := channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return &q
}

func publish(channel *amqp.Channel, queue *amqp.Queue, job string) {
	err := channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(job),
		})
	log.Printf(" [x] Sent %s", job)
	failOnError(err, "Failed to publish a message")
}

func main() {
	// Get the configuration from env variables
	rmqHost := os.Getenv("RMQ_HOST")
	jobQueueName := os.Getenv("RMQ_JOB_QUEUE")
	responseQueueName := os.Getenv("RMQ_RES_QUEUE")

	// Open the connection to RabbitMQ
	ch, conn := openRmqChannel(rmqHost)

	// Declare the queue to be read from
	jobQueue := declareQueue(ch, jobQueueName)

	// Declare the queue to be pushed to
	responseQueue := declareQueue(ch, responseQueueName)

	// Defer the connection's and channel's closing
	defer conn.Close()
	defer ch.Close()

	err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		jobQueue.Name, // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)

			publish(ch, responseQueue, "This has been donedded")

			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
