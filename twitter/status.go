package twitter

import "github.com/rayyildiz/twitter-streaming/util"

type StatusServie struct {
	client *util.HttpClient
}

func newStatusService(httpClient *util.HttpClient) *StatusServie {
	return &StatusServie{
		client: httpClient,
	}
}
