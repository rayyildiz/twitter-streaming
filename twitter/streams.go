package twitter

import (
	"github.com/cenkalti/backoff"
	"github.com/rayyildiz/twitter-streaming/util"

	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	twitterPublicStreamAPI = "https://stream.twitter.com/1.1/"
)

type StreamService struct {
	client *util.HttpClient
	public *util.HttpClient
}

type Stream struct {
	client    *util.HttpClient
	Messages  chan interface{}
	ok        chan struct{}
	syncGroup *sync.WaitGroup
	body      io.Closer
}

func newStreamService(httpClient *util.HttpClient) *StreamService {
	streams := StreamService{
		client: httpClient,
		public: httpClient.New().BasePath(twitterPublicStreamAPI).Path("statuses/"),
	}

	return &streams
}

func newStream(httpClient *util.HttpClient, req *http.Request) *Stream {
	stream := &Stream{
		client:    httpClient.BasePath(twitterPublicStreamAPI).Post("statuses/filter.json"),
		ok:        make(chan struct{}),
		Messages:  make(chan interface{}),
		syncGroup: &sync.WaitGroup{},
	}

	stream.syncGroup.Add(1)
	go stream.retry(req, newExponentialBackOff(), newAggressiveExponentialBackOff())
	return stream
}

func (this *StreamService) Sample() (*Stream, error) {
	req, err := this.public.New().Get("statuses/sample.json").Request()
	if err != nil {
		return nil, err
	}
	return newStream(this.client, req), nil
}

// Stream filter . I just need language right now
type StreamFilterParams struct {
	Localtion []string `url:"locations,omitempty,comma"`
	Track     []string `url:"track,omitempty,comma"`
}

func (this *StreamService) Filter(p *StreamFilterParams) (*Stream, error) {
	req, err := this.public.New().Post("statuses/filter.json").Query(p).Request()
	if err != nil {
		return nil, err
	}
	return newStream(this.client, req), nil
}

func (this *Stream) Stop() {
	close(this.ok)

	if this.body != nil {
		this.body.Close()
	}

	this.syncGroup.Wait()
}

func (this *Stream) retry(req *http.Request, expBackOff backoff.BackOff, aggExpBackOff backoff.BackOff) {
	defer close(this.Messages)
	defer this.syncGroup.Done()

	var wait time.Duration
	for !stopped(this.ok) {
		resp, err := this.client.Doer.Do(req)
		if err != nil {
			this.Messages <- err
			return
		}

		defer resp.Body.Close()
		this.body = resp.Body
		switch resp.StatusCode {
		case 200:
			// receive stream response Body, handles closing
			this.receive(resp.Body)
			expBackOff.Reset()
			aggExpBackOff.Reset()
		case 503:
			// exponential backoff
			wait = expBackOff.NextBackOff()
		case 420, 429:
			// aggressive exponential backoff
			wait = aggExpBackOff.NextBackOff()
		default:
			// stop retrying for other response codes
			resp.Body.Close()
			return
		}
		// close response before each retry
		resp.Body.Close()
		if wait == backoff.Stop {
			return
		}
		sleepOrDone(wait, this.ok)
	}
}

func (this *Stream) receive(body io.ReadCloser) {
	defer body.Close()
	scanner := bufio.NewScanner(body)
	scanner.Split(scanLines)
	for !stopped(this.ok) && scanner.Scan() {
		token := scanner.Bytes()
		if len(token) == 0 {
			continue
		}
		select {
		case this.Messages <- getMessage(token):
			continue
		case <-this.ok:
			return
		}
	}
}

func getMessage(token []byte) interface{} {
	var data map[string]interface{}
	// unmarshal JSON encoded token into a map for
	err := json.Unmarshal(token, &data)
	if err != nil {
		return err
	}
	return decodeMessage(token, data)
}

func decodeMessage(token []byte, data map[string]interface{}) interface{} {

	if hasPath(data, "retweet_count") {
		tweet := new(Tweet)
		json.Unmarshal(token, tweet)
		// fmt.Println("tweet:", tweet.Text)
		return tweet
	} else if hasPath(data, "direct_message") {
		notice := new(directMessageNotice)
		json.Unmarshal(token, notice)
		return notice.DirectMessage
	} else if hasPath(data, "delete") {
		notice := new(statusDeletionNotice)
		json.Unmarshal(token, notice)
		return notice.Delete.StatusDeletion
	} else if hasPath(data, "scrub_geo") {
		notice := new(locationDeletionNotice)
		json.Unmarshal(token, notice)
		return notice.ScrubGeo
	} else if hasPath(data, "limit") {
		notice := new(streamLimitNotice)
		json.Unmarshal(token, notice)
		return notice.Limit
	} else if hasPath(data, "status_withheld") {
		notice := new(statusWithheldNotice)
		json.Unmarshal(token, notice)
		return notice.StatusWithheld
	} else if hasPath(data, "user_withheld") {
		notice := new(userWithheldNotice)
		json.Unmarshal(token, notice)
		return notice.UserWithheld
	} else if hasPath(data, "disconnect") {
		notice := new(streamDisconnectNotice)
		json.Unmarshal(token, notice)
		return notice.StreamDisconnect
	} else if hasPath(data, "warning") {
		notice := new(stallWarningNotice)
		json.Unmarshal(token, notice)
		return notice.StallWarning
	} else if hasPath(data, "friends") {
		friendsList := new(FriendsList)
		json.Unmarshal(token, friendsList)
		return friendsList
	} else if hasPath(data, "event") {
		event := new(Event)
		json.Unmarshal(token, event)
		return event
	}
	return data
}

func hasPath(data map[string]interface{}, key string) bool {
	_, ok := data[key]
	return ok
}
