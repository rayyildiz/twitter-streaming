package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/rayyildiz/twitter-streaming/api"
	"github.com/rayyildiz/twitter-streaming/conf"
	"github.com/rayyildiz/twitter-streaming/twitter"
	"golang.org/x/net/websocket"
	"os"
)

// VERSION number of program
const VERSION = `0.3`

// Msg structure for client-server side
type Msg struct {
	Count int    `json:"count"`
	Text  string `json:"text"`
}

// Messages is for list of Msg.
type Messages []Msg

// Handle http request and return response.
func Handle(message interface{}, ws *websocket.Conn) {
	var err error

	switch message.(type) {
	case *twitter.Tweet:

		t := message.(*twitter.Tweet)

		word := WordCount(t.Text)

		if err = websocket.JSON.Send(ws, word); err != nil {
			fmt.Println("Can't send")
			break
		}

		// fmt.Println(WordCount(t.Text))
	}
}

// WordCount is for count word in a string.
func WordCount(s string) Messages {
	strs := strings.Fields(s)
	res := make(map[string]int)

	for _, str := range strs {
		if len(str) > 2 {
			res[strings.ToLower(str)]++
		}
	}

	messages := []Msg{}

	for k, v := range res {
		m := Msg{
			Text:  k,
			Count: v,
		}
		messages = append(messages, m)
	}

	return messages
}

var (
	configFile  = flag.String("config", "./config.json", "Config file for twitter token, porst....")
	version     = flag.Bool("v", false, "display version info and exit")
	help        = flag.Bool("h", false, "display usage")
	twitterConf conf.TwitterKey
)

// main method for application.
func main() {
	log.Println("Application is starting")

	flag.Parse()
	if *version {
		fmt.Printf("twitter-streaming %s\n", VERSION)
		return
	}

	if *help {
		fmt.Println("twitter-streaming Help")
		fmt.Println("twitter-streaming -v => Display version")
		fmt.Println("twitter-streaming -h => Display usage")
		fmt.Println("twitter-streaming -config=config.json => Set the config file")
		fmt.Println("                   config.json")
		fmt.Println("                   {")
		fmt.Println("                       \"Port\":3000,")
		fmt.Println("                       \"Twitter\" : {")
		fmt.Println("                            	\"ConsumerKey\":\"xxx\"")
		fmt.Println("                            	\"ConsumerKeyConsumerSecret\":\"xxx\"")
		fmt.Println("                            	\"ConsumerKeyAccessTokenKey\":\"xxx\"")
		fmt.Println("                            	\"ConsumerKeyAccessTokenSecret\":\"xxx\"")
		fmt.Println("                        }")
		fmt.Println("                   }")
		return
	}

	log.Println("Lets load conf file, ", *configFile)

	c, err := conf.Load(*configFile)
	if err != nil {
		fmt.Println(err)
	}

	if len(strings.TrimSpace(c.Twitter.AccessTokenKey)) == 0 {
		log.Println("You have to define 'AccessTokenKey'. More info use documentation.")
		os.Exit(1)
	}

	if len(strings.TrimSpace(c.Twitter.AccessTokenSecret)) == 0 {
		log.Println("You have to define 'AccessTokenSecret'. More info use documentation.")
		os.Exit(1)
	}

	if len(strings.TrimSpace(c.Twitter.ConsumerKey)) == 0 {
		log.Println("You have to define 'ConsumerKey'. More info use documentation.")
		os.Exit(1)
	}

	if len(strings.TrimSpace(c.Twitter.ConsumerSecret)) == 0 {
		log.Println("You have to define 'ConsumerSecret'. More info use documentation.")
		os.Exit(1)
	}

	log.Println("Created conf. Now create http Handler")

	twitterConf = c.Twitter

	// go initTwitter(c.Twitter)
	initHTTP(c)

	log.Println("App started successfully")
}

// initTwitter object.
func initTwitter(c conf.TwitterKey) {

	twitterClient, err := twitter.NewClient(&c)
	if err != nil {
		fmt.Println("Error creating new client")
	}

	p := &twitter.StreamFilterParams{
		Localtion: []string{"-74,40,-73,41", "-122.75,36.8", "-121.75,37.8"}, // NYC, SF
	}
	filter, err := twitterClient.Streams.Filter(p)
	if err != nil {
		log.Println("Could not create a stream")
		return
	}
	defer filter.Stop()

	for message := range filter.Messages {
		Handle(message, nil)
	}

	filterController, err := api.NewFilterController(&c)
	if err != nil {
		fmt.Printf("Some error while creating controller %s\n", err)
		return
	}

	filterController.FilterByLocation(p.Localtion)
}

// Filter for websocket.
func Filter(ws *websocket.Conn) {
	var err error

	twitterClient, err := twitter.NewClient(&twitterConf)
	if err != nil {
		log.Println("Error creating new client")
	}

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)
		p := &twitter.StreamFilterParams{
			Track: []string{reply},
		}
		filter, err := twitterClient.Streams.Filter(p)
		if err != nil {
			fmt.Println("Could not create a stream")
			return
		}
		defer filter.Stop()

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)
		for message := range filter.Messages {
			Handle(message, ws)
		}
	}

}

// initHTTP for application
func initHTTP(c *conf.Configuration) {
	// Static http serve
	if c.Port == 0 {
		c.Port = 3000
	}

	strPort := ":" + strconv.Itoa(c.Port)
	fs := http.FileServer(http.Dir("/app/public"))
	log.Println("Port ", strPort)

	http.Handle("/api/filter", websocket.Handler(Filter))

	http.Handle("/", fs)

	http.ListenAndServe(strPort, nil)
	log.Println("Http serve on ", strPort)
}
