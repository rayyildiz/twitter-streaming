package twitter

import "github.com/rayyildiz/twitter-streaming/util"

type StatusService struct {
	client *util.HttpClient
}

func newStatusService(httpClient *util.HttpClient) *StatusService {
	return &StatusService{
		client: httpClient,
	}
}
