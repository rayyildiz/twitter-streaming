package twitter

import (
	"github.com/rayyildiz/twitter-streaming/conf"
	"testing"
)

func TestNewClient(t *testing.T) {
	twitterConf := conf.TwitterKey{
		ConsumerSecret:    "z9xRHwVjbiDyvH2LftdyZGnkYwISl7Cfj7L8iDQAdhfCcEmllx",
		ConsumerKey:       "B2fDAgliWGT1DjuIsiWce6wtC",
		AccessTokenKey:    "13364212-TUxcRFxZACFbVDgNv6yqyEPec4pMzxdMBivImgi3B",
		AccessTokenSecret: "wA4R1yKMuvGbMmEQQj2uHOmi0RUszu40dNfik854K71zD",
	}
	client, err := NewClient(&twitterConf)

	if client == nil && err != nil {
		t.Error("Must be return error nil")
	}
}
