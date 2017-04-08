package twitter

import (
	"github.com/rayyildiz/twitter-streaming/util"
	"testing"
)

func TestNewStreamsService(t *testing.T) {
	httpClient := util.New()

	service := newStreamService(httpClient.New())

	if service == nil {
		t.Error("Sstream Services is not initialized")
	}
}
