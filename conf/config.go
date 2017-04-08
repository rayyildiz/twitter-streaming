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
	config := Configuration{
		Port: 3000,
		Twitter: TwitterKey{
			ConsumerSecret:    "z9xRHwVjbiDyvH2LftdyZGnkYwISl7Cfj7L8iDQAdhfCcEmllx",
			ConsumerKey:       "B2fDAgliWGT1DjuIsiWce6wtC",
			AccessTokenKey:    "13364212-TUxcRFxZACFbVDgNv6yqyEPec4pMzxdMBivImgi3B",
			AccessTokenSecret: "wA4R1yKMuvGbMmEQQj2uHOmi0RUszu40dNfik854K71zD",
		},
	}

	file, err := os.Open(filePath)
	if err != nil {
		return &config, err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)

	return &config, err
}
