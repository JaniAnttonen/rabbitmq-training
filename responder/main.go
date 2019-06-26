package main

import (
	"os"

	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
	"github.com/vincentLiuxiang/lu"
)

//ctx.QueryArgs().Peek("haha")

func main() {
	// Get the configuration from env variables
	apiPort := os.Getenv("API_PORT")
	redisHost := os.Getenv("REDIS_HOST")

	// Open Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	// Initialize the API
	api := lu.New()

	// Define routes
	api.Use("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetStatusCode(200)

		get := redisClient.Get("responses")

		ctx.SetBody([]byte(get.String()))
	})

	// Listen to connections
	api.Listen(apiPort)
}
