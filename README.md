# Twerkov - Generate tweets using Markov chains!

Twerkov is a simple command line application that generates and tweets tweets using Markov chains. In order to work, Twerkov needs access to a MySQL database and the Twitter API. The tweets are generated from tweets by actual Twitter accounts, tweets need to be cached manually.


## Configuration

Twerkov can be configured using environment variables or an `.env` file, a sample of which can be found in the root directory of the project. All of the following values must be specified:

* `TWERKOV_MYSQL_HOSTNAME`
* `TWERKOV_MYSQL_DATABASE`
* `TWERKOV_MYSQL_USERNAME`
* `TWERKOV_MYSQL_PASSWORD`
* `TWERKOV_TWITTER_ACCESS_TOKEN`
* `TWERKOV_TWITTER_ACCESS_TOKEN_SECRET`
* `TWERKOV_TWITTER_CONSUMER_KEY`
* `TWERKOV_TWITTER_CONSUMER_KEY_SECRET`

The values of all of the configuration variables are expected to be strings. Twitter API keys can be generated [here](https://developer.twitter.com/en/apps).


## Commands

`twerkov init` - Creates the database structure needed for the application to function.

`twerkov cache <Twitter username>` - Caches the latest tweets from the user with the specified username. It is able to find the latest 200 tweets on the timeline of the user, but it excludes retweets, so the actual number of cached tweets is usually lower.

`twerkov generate` - Generates a tweet using a Markov chain and tweets it.