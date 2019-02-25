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
		log.Fatal("twerkov <crawl | generate>")
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
	case "crawl":
		if len(os.Args) < 3 {
			log.Fatal("twerkov crawl <username>")
		}

		app.CrawlUserTweets(os.Args[2])

	case "generate":
		app.GenerateTweet()

	default:
		log.Fatal("twerkov <crawl | generate>")
	}
}
