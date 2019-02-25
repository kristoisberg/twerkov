package twerkov

import (
	"log"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/mb-14/gomarkov"
)

// App is the main type of the application
type App struct {
	Config   Config
	API      *anaconda.TwitterApi
	Database *Database
}

// Init initializes the application and connects to the database and Twitter API
func (app *App) Init(config Config) (err error) {
	app.Config = config
	app.Database, err = CreateDatabaseConnection(config)

	if err != nil {
		return
	}

	err = app.Database.Handle.Ping()

	if err != nil {
		return
	}

	app.API = anaconda.NewTwitterApiWithCredentials(config.TwitterAccessToken, config.TwitterAccessTokenSecret, config.TwitterConsumerKey, config.TwitterConsumerKeySecret)

	return
}

// InitializeDatabase creates the database structure needed for the application to function
func (app *App) InitializeDatabase() {
	stmt, err := app.Database.Handle.Query(`CREATE TABLE tweets (
		id BIGINT(20) UNSIGNED NOT NULL,
		user_id BIGINT(20) UNSIGNED NOT NULL,
		text TEXT NOT NULL,
		PRIMARY KEY (id)
	)
	COLLATE='utf8mb4_general_ci'
	ENGINE=InnoDB
	;
	`)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer stmt.Close()

	log.Println("Database structure successfully created!")
}

// CacheUserTweets caches the latest tweets by the user with the specified username
func (app *App) CacheUserTweets(username string) {
	users, err := app.API.GetUsersLookup(username, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	if len(users) != 1 {
		log.Fatal("Expected one user, found " + string(len(users)) + "!")
	}

	userID := users[0].Id

	userData := url.Values{}
	userData.Set("screen_name", username)
	userData.Set("include_rts", "false")
	userData.Set("count", "200")

	tweets, err := app.API.GetUserTimeline(userData)

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, tweet := range tweets {
		log.Println(tweet.FullText)

		stmt, err := app.Database.Handle.Prepare("INSERT IGNORE INTO `tweets` (`id`, `user_id`, `text`) VALUES (?, ?, ?)")

		if err != nil {
			log.Println(err.Error())
			continue
		}

		defer stmt.Close()

		_, err = stmt.Exec(tweet.Id, userID, tweet.FullText)

		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
}

// GenerateTweet generates a tweet using Markov chains and tweets it
func (app *App) GenerateTweet() {
	chain := gomarkov.NewChain(1)

	stmt, err := app.Database.Handle.Query("SELECT `text` FROM `tweets`")

	if err != nil {
		log.Fatal(err.Error())
	}

	defer stmt.Close()

	for stmt.Next() {
		var text string

		stmt.Scan(&text)

		if err != nil {
			log.Fatal(err.Error())
		}

		words := strings.Split(text, " ")

		if len(words) < 3 {
			continue
		}

		chain.Add(words)
	}

	tokens := []string{gomarkov.StartToken}

	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, err := chain.Generate(tokens[(len(tokens) - 1):])

		if err != nil {
			log.Fatal(err.Error())
		}

		tokens = append(tokens, next)
	}

	tweet := strings.Join(tokens[1:len(tokens)-1], " ")

	app.API.PostTweet(tweet, url.Values{})

	log.Println("New tweet: ", tweet)
}