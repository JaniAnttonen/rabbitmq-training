package main

import (
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
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

func main() {
	// Get the configuration from env variables
	rmqHost := os.Getenv("RMQ_HOST")
	jobQueueName := os.Getenv("RMQ_JOB_QUEUE")
	redisHost := os.Getenv("REDIS_HOST")

	// Open the connection to RabbitMQ
	ch, conn := openRmqChannel(rmqHost)

	// Declare the queue to be read from
	jobQueue := declareQueue(ch, jobQueueName)

	// Defer the connection's and channel's closing
	defer conn.Close()
	defer ch.Close()

	// Open Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	err := ch.Qos(
		1,
		0,
		false,
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		jobQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			t := float64(time.Now().UTC().UnixNano())
			response := &redis.Z{Score: t, Member: "This is a response"}

			redisClient.ZAdd("responses", response)

			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
