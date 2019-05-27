package twerkov

import (
	"errors"
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

// Init initialises the application and connects to the database and Twitter API
func (app *App) Init(config Config) (err error) {
	app.Config = config

	app.Database, err = CreateDatabaseConnection(config)
	if err != nil {
		return
	}

	if err = app.Database.Handle.Ping(); err != nil {
		return
	}

	app.API = anaconda.NewTwitterApiWithCredentials(
		config.TwitterAccessToken,
		config.TwitterAccessTokenSecret,
		config.TwitterConsumerKey,
		config.TwitterConsumerKeySecret,
	)

	return
}

// InitialiseDatabase creates the database structure needed for the application to function
func (app *App) InitialiseDatabase() error {
	stmt, err := app.Database.Handle.Query(`CREATE TABLE tweets (
		id BIGINT(20) UNSIGNED NOT NULL,
		user_id BIGINT(20) UNSIGNED NOT NULL,
		text TEXT NOT NULL,
		PRIMARY KEY (id)
	)`)

	if err != nil {
		return err
	}

	stmt.Close()
	return nil
}

// CacheUserTweets caches the latest tweets by the user with the specified username
func (app *App) CacheUserTweets(username string) (count int, err error) {
	users, err := app.API.GetUsersLookup(username, nil)
	if err != nil {
		return 0, err
	}

	if len(users) != 1 {
		return 0, errors.New("Expected one user, found " + string(len(users)) + "!")
	}

	userID := users[0].Id

	userData := url.Values{}
	userData.Set("screen_name", username)
	userData.Set("include_rts", "false")
	userData.Set("count", "200")

	tweets, err := app.API.GetUserTimeline(userData)
	if err != nil {
		return 0, err
	}

	for _, tweet := range tweets {
		stmt, err := app.Database.Handle.Prepare("INSERT IGNORE INTO `tweets` (`id`, `user_id`, `text`) VALUES (?, ?, ?)")
		if err != nil {
			return count, err
		}

		defer stmt.Close()

		if _, err = stmt.Exec(tweet.Id, userID, tweet.FullText); err != nil {
			return count, err
		}

		count++
	}

	return count, nil
}

// CreateTweet creates a tweet using Markov chains
func (app *App) CreateTweet() (string, error) {
	chain := gomarkov.NewChain(1)

	stmt, err := app.Database.Handle.Query("SELECT `text` FROM `tweets`")
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	for stmt.Next() {
		var text string

		if err := stmt.Scan(&text); err != nil {
			return "", err
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
			return "", err
		}

		tokens = append(tokens, next)
	}

	return strings.Join(tokens[1:len(tokens)-1], " "), nil
}

// PostTweet posts a tweet on Twitter
func (app *App) PostTweet(tweet string) (err error) {
	_, err = app.API.PostTweet(tweet, url.Values{})
	return
}
