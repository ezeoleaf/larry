package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func twitterClient() *twitter.Client {
	config := oauth1.NewConfig("rvBUw3Jn2KDEA3YTV7AgqmL9f", "eVcymJznO9zv2C5qifJOSTcbqsuhX5At4qCX4LkvQ6CVE3kgkj")
	token := oauth1.NewToken("1328795920666931203-1btfihoHYnBdkyfvRV33Ahfr1ywUpB", "vs4Y6HcCxkOVASvzEwGOKgCXsy02rEEjWrMQNEy4wYse6")

	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	return twitter.NewClient(httpClient)
}
