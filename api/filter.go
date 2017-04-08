package api

import (
	"fmt"
	"github.com/rayyildiz/twitter-streaming/conf"
	"github.com/rayyildiz/twitter-streaming/twitter"
)

type FilterController struct {
	conf   *conf.TwitterKey
	client *twitter.TwitterClient
}

func NewFilterController(c *conf.TwitterKey) (*FilterController, error) {
	twitterClient, err := twitter.NewClient(c)
	if err != nil {
		fmt.Println("Error creating new client")
		return nil, err
	}

	controller := FilterController{
		conf:   c,
		client: twitterClient,
	}

	return &controller, nil
}

func (f *FilterController) FilterByLocation(locations []string) {
	params := &twitter.StreamFilterParams{
		Localtion: []string{"-74,40,-73,41", "-122.75,36.8", "-121.75,37.8"}, // NYC, SF
	}

	f.client.Streams.Filter(params)
}
