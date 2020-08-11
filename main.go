package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

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
	db    *sql.DB
)

func init() {
	flag.StringVar(&mode, "mode", defaultMode, "run the app in api or consumer mode")
	flag.StringVar(&topic, "topic", "", "run the consumer on a specific topic")

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", user, password, dbname))
	if err != nil {
		panic(err.Error())
	}
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
