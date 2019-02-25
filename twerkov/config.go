package twerkov

// Config contains the configuration data for the application
type Config struct {
	MySQLHostname string `envconfig:"mysql_hostname"`
	MySQLDatabase string `envconfig:"mysql_database"`
	MySQLUsername string `envconfig:"mysql_username"`
	MySQLPassword string `envconfig:"mysql_password"`

	TwitterAccessToken       string `envconfig:"twitter_access_token"`
	TwitterAccessTokenSecret string `envconfig:"twitter_access_token_secret"`
	TwitterConsumerKey       string `envconfig:"twitter_consumer_key"`
	TwitterConsumerKeySecret string `envconfig:"twitter_consumer_key_secret"`
}
