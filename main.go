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

var mode string

func init() {
	flag.StringVar(&mode, "mode", defaultMode, "run the app in api or consumer mode")
}

func main() {

	flag.Parse()

	switch mode {
	case apiMode:
		//run api's server
		api()
	case consumerMode:
		//runs consumer
		consume()
	default:
		panic(fmt.Sprintf("no known mode for %s", mode))
	}
}
