package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	goquery "github.com/google/go-querystring/query"
)

type Caller interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClient struct {
	Doer     Caller
	baseUrl  string
	rawUrl   string
	method   string
	header   http.Header
	bodyJSON interface{}
	query    []interface{}
}

// New HttpClient  with default httpclient
func New() *HttpClient {
	return &HttpClient{
		Doer:   http.DefaultClient,
		method: "GET",
		header: make(http.Header),
		query:  make([]interface{}, 0),
	}
}

// Clone and return new httpclient
func (this *HttpClient) New() *HttpClient {
	copiedHeader := make(http.Header)
	for k, v := range this.header {
		copiedHeader[k] = v
	}

	return &HttpClient{
		Doer:     this.Doer,
		method:   this.method,
		baseUrl:  this.baseUrl,
		rawUrl:   this.baseUrl,
		header:   copiedHeader,
		bodyJSON: this.bodyJSON,
		query:    append([]interface{}{}, this.query...),
	}
}

// Set httpi client for tranportation OAuth impelemntation
func (this *HttpClient) Client(hc *http.Client) *HttpClient {
	if hc != nil {
		return this.Caller(hc)
	}

	return this
}

// Set basepath
func (this *HttpClient) BasePath(url string) *HttpClient {
	this.baseUrl = url
	return this
}

// Set the caller
func (this *HttpClient) Caller(caller *http.Client) *HttpClient {
	if caller != nil {
		this.Doer = caller
	} else {
		this.Doer = http.DefaultClient
	}

	return this
}

// Add header
func (this *HttpClient) AddHeader(key, value string) *HttpClient {
	this.header.Add(key, value)
	return this
}

// Set Header
func (this *HttpClient) SetHeader(key, value string) *HttpClient {
	this.header.Set(key, value)
	return this
}

// Path
func (this *HttpClient) Path(path string) *HttpClient {

	baseUrl, baseErr := url.Parse(this.baseUrl)
	httpUrl, httpErr := url.Parse(path)
	if baseErr == nil && httpErr == nil {
		this.rawUrl = baseUrl.ResolveReference(httpUrl).String()
		fmt.Println("RawUrl:", this.rawUrl, " Method:", this.method)
		return this
	}

	return this
}

// GET
func (this *HttpClient) Get(path string) *HttpClient {
	this.method = "GET"
	return this.Path(path)
}

// POST
func (this *HttpClient) Post(path string) *HttpClient {
	this.method = "POST"
	return this.Path(path)
}

// HEAD
func (this *HttpClient) Head(path string) *HttpClient {
	this.method = "HEAD"
	return this.Path(path)
}

// Set body json
func (this *HttpClient) SetBody(bodyJson interface{}) *HttpClient {
	if bodyJson != nil {
		this.bodyJSON = bodyJson
		this.SetHeader("Content-Type", "application/json")
	}

	return this
}

func (this *HttpClient) Query(params interface{}) *HttpClient {
	if params != nil {
		this.query = append(this.query, params)
	}

	return this
}

// private
func addHeaders(req *http.Request, header http.Header) {
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

// private
func encodeBodyJSON(bodyJSON interface{}) (io.Reader, error) {
	var buf = new(bytes.Buffer)
	if bodyJSON != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(bodyJSON)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

// private
func decodeResponseBodyJSON(response *http.Response, f interface{}) error {
	return json.NewDecoder(response.Body).Decode(f)
}

// private
func decodeResponseJSON(response *http.Response, successF, failureF interface{}) error {
	// fmt.Println("status:", response.StatusCode)

	if code := response.StatusCode; 200 <= code && code <= 299 {
		if successF != nil {
			return decodeResponseBodyJSON(response, successF)
		}
	} else {
		if failureF != nil {
			return decodeResponseBodyJSON(response, failureF)
		}
	}
	return nil
}

// private
func (this *HttpClient) getRequestBody() (body io.Reader, err error) {
	if this.bodyJSON != nil && this.header.Get("Content-Type") == "application/json" {
		body, err = encodeBodyJSON(this.bodyJSON)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

// Return a http.Request
func (this *HttpClient) Request() (*http.Request, error) {
	reqURL, err := url.Parse(this.rawUrl)

	if err != nil {
		return nil, err
	}

	err = addQueryStructs(reqURL, this.query)
	if err != nil {
		return nil, err
	}

	fmt.Println("RegUrl:", reqURL)

	body, err := this.getRequestBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(this.method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	addHeaders(req, this.header)
	return req, err
}

// Return Receive
func (this *HttpClient) Receive(successF, failureF interface{}) (*http.Response, error) {
	req, err := this.Request()
	if err != nil {
		return nil, err
	}

	return this.Call(req, successF, failureF)
}

func (this *HttpClient) Call(req *http.Request, successF, failureF interface{}) (*http.Response, error) {

	response, err := this.Doer.Do(req)

	if err != nil {
		return response, err
	}

	defer response.Body.Close()

	if strings.Contains(response.Header.Get("Content-Type"), "application/json") {
		err = decodeResponseJSON(response, successF, failureF)
	}
	return response, err
}

func addQueryStructs(reqURL *url.URL, queries []interface{}) error {
	urlValues, err := url.ParseQuery(reqURL.RawQuery)
	if err != nil {
		return err
	}
	for _, query := range queries {
		queryValues, err := goquery.Values(query)
		if err != nil {
			return err
		}
		for key, values := range queryValues {
			for _, value := range values {
				urlValues.Add(key, value)
			}
		}
	}
	reqURL.RawQuery = urlValues.Encode()
	return nil
}
