package tegami

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/hkurokawa/go-connpass"
)

// Notifier is an interface to notify event
type Notifier interface {
	Notify(event *connpass.Event) (anaconda.Tweet, error)
}

// TwitterClient manage credeintial of Twitter.
type TwitterClient struct {
	api *anaconda.TwitterApi
}

// NewTwitterClient instanciates a Twitter client.
func NewTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *TwitterClient {
	client := &TwitterClient{}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	client.api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	return client
}

// Notify posts a tweet and returns a Tweet and an error.
func (client *TwitterClient) Notify(event *connpass.Event) (anaconda.Tweet, error) {
	return client.api.PostTweet(event.Title+" "+event.Url, nil)
}
