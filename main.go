package main

import (
	"flag"
	"fmt"

	_ "github.com/joho/godotenv/autoload" //will load .env vars automatically
)

const (
	apiMode      = "api"
	consumerMode = "consumer"
	defaultMode  = apiMode
)

var (
	mode  string
	topic string
)

func init() {
	flag.StringVar(&mode, "mode", defaultMode, "run the app in api or consumer mode")
	flag.StringVar(&topic, "topic", "", "run the consumer on a specific topic")
}

func main() {

	flag.Parse()

	switch mode {
	case apiMode:
		//run api's server
		api()
	case consumerMode:
		//runs consumer
		consume(topic)
	default:
		panic(fmt.Sprintf("no known mode for %s", mode))
	}
}
