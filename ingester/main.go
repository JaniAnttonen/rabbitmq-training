package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
	"github.com/valyala/fasthttp"
	"github.com/vincentLiuxiang/lu"
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
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	return &q
}

func publish(channel *amqp.Channel, queue *amqp.Queue, job string) {
	err := channel.Publish(
		"",
		queue.Name,
		false,
		false,
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
	apiPort := os.Getenv("API_PORT")

	// Open the connection to RabbitMQ
	ch, conn := openRmqChannel(rmqHost)

	// Declare the queue to be pushed to
	q := declareQueue(ch, jobQueueName)

	// Defer the connection's and channel's closing
	defer conn.Close()
	defer ch.Close()

	// Initialize the API
	api := lu.New()

	// Define routes
	api.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetStatusCode(200)

		// Publish the message to queue
		publish(ch, q, "Hello world!")
	})

	// Listen to connections
	api.Listen(apiPort)
}
