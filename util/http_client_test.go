package util

import (
	"testing"
)

func TestHttpClientNew(t *testing.T) {
	client := New()

	if client == nil {
		t.Error("Client return not null")
	}
}

func TestBasePath(t *testing.T) {
	cases := []string{"http://rayyildiz.com/", "http://rayyildiz.com", "/path", "path"}

	for _, base := range cases {
		client := New().BasePath(base)
		if client.baseUrl != base {
			t.Errorf("expected %s, got %s", base, client.baseUrl)
		}

	}
}

func TestGet(t *testing.T) {
	cases := []string{"/api", "/api/hello", "/ping", "/pong"}

	for _, path := range cases {
		client := New().BasePath("http://rayyildiz.com").Get(path)

		if client.method != "GET" {
			t.Errorf("expected GET method, got %s", client.method)
		}
		newRawUrl := "http://rayyildiz.com" + path

		if client.rawUrl != newRawUrl {
			t.Errorf("expected %s, got %s", newRawUrl, client.rawUrl)
		}
	}
}

func TestPost(t *testing.T) {
	cases := []string{"/api", "/api/hello", "/ping", "/pong"}

	for _, path := range cases {
		client := New().BasePath("http://rayyildiz.com").Post(path)

		if client.method != "POST" {
			t.Errorf("expected POST method, got %s", client.method)
		}
		newRawUrl := "http://rayyildiz.com" + path

		if client.rawUrl != newRawUrl {
			t.Errorf("expected %s, got %s", newRawUrl, client.rawUrl)
		}
	}
}

func TestBodyJSON(t *testing.T) {

	type Params struct {
		Fullname string `url:"full_name"`
		Follower int    `url:"follower"`
	}

	var mockParam = Params{Fullname: "Ramazan", Follower: 100}

	cases := []struct {
		initial  interface{}
		bodyJSON interface{}
		expected interface{}
	}{
		{nil, nil, nil},
		{nil, mockParam, mockParam},
		{mockParam, nil, mockParam},
	}

	for _, body := range cases {
		client := New()

		client.bodyJSON = body.initial
		client.SetBody(body.bodyJSON)

		if client.bodyJSON != body.expected {
			t.Errorf("expected %v, got %v", body.expected, client.bodyJSON)
		}
	}
}
