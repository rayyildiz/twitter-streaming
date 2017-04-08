package twitter

import (
	"errors"
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/rayyildiz/twitter-streaming/conf"
	"github.com/rayyildiz/twitter-streaming/util"
	"io/ioutil"
	"net/http"
)

const (
	twitterAPI = "https://api.twitter.com/1.1/"
	userAgent  = "rayyildiz-streaming-app"
)

var debug = false

type TwitterClient struct {
	httpClient *util.HttpClient
	Streams    *StreamService
}

func NewClient(conf *conf.TwitterKey) (*TwitterClient, error) {
	// OAuth http.Client will automatically authorize Requests
	myHttpClient, err := newHttpCLient(conf.ConsumerKey, conf.ConsumerSecret, conf.AccessTokenKey, conf.AccessTokenSecret)
	if err != nil {
		return nil, err
	}

	httpClient := util.New().Client(myHttpClient).BasePath(twitterAPI).SetHeader("User-Agent", userAgent)

	client := TwitterClient{
		httpClient: httpClient,
		//		Status:     newStatusService(httpClient.New()),
		Streams: newStreamService(httpClient.New()),
	}

	return &client, nil
}

func newHttpCLient(consumerKey, consumerSecret, accessToken, accessSecret string) (*http.Client, error) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	if httpClient != nil {
		if debug == true {
			path := "https://api.twitter.com/1.1/statuses/home_timeline.json?count=1"
			resp, _ := httpClient.Get(path)
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("Raw Response Body:\n%v\n", string(body))
		}

		return httpClient, nil
	}

	return nil, errors.New("Error creating httpClient")
}
