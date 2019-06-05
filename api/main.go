package main

import (
	"log"

	"github.com/valyala/fasthttp"
	"github.com/vincentLiuxiang/lu"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	// API STUFF
	api := lu.New()

	api.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte("hello world\n"))
	})

	api.Listen(":1337")

	// MESSAGE QUEUE
	/* rmqHost := os.Getenv("RMQ_HOST")
	conn, err := amqp.Dial("amqp://guest:guest@" + rmqHost)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message") */
}
