package conf

import (
	"encoding/json"
	"os"
)

type TwitterKey struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessTokenKey    string
	AccessTokenSecret string
}

type Configuration struct {
	Port    int
	Twitter TwitterKey
}

func Load(filePath string) (*Configuration, error) {
	_config := Configuration{
		Port: 3000,
		Twitter: TwitterKey{
			ConsumerSecret:    "",
			ConsumerKey:       "",
			AccessTokenKey:    "",
			AccessTokenSecret: "",
		},
	}

	file, err := os.Open(filePath)
	if err != nil {
		return &_config, err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&_config)

	return &_config, err
}
