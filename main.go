package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/kristoisberg/twerkov/twerkov"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("twerkov <init | cache | test | tweet>")
	}

	config := twerkov.Config{}
	err := envconfig.Process("TWERKOV", &config)

	if err != nil {
		log.Fatal(err.Error())
	}

	app := twerkov.App{}

	err = app.Init(config)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer app.Database.Handle.Close()

	switch os.Args[1] {
	case "init":
		err := app.InitializeDatabase()

		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Database structure successfully created!")

	case "cache":
		if len(os.Args) < 3 {
			log.Fatal("twerkov cache <Twitter username>")
		}

		count, err := app.CacheUserTweets(os.Args[2])

		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Successfully cached", count, "tweets!")

	case "test":
		tweet, err := app.CreateTweet()

		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println(tweet)

	case "tweet":
		tweet, err := app.CreateTweet()

		if err != nil {
			log.Fatal(err.Error())
		}

		app.PostTweet(tweet)
		log.Println("New tweet:", tweet)

	default:
		log.Fatal("twerkov <init | cache | test | tweet>")
	}
}
