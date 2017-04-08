# A Simple Twitter Streaming written in #Golang


## First Step

* You need make ```make```. Be sure you have ```make``` command at ```PATH```variable.
* This project is written in **golang** . So you need to install latest golang. You can find details at https://golang.org/doc/install
* You need nodejs. I strongly recommend to install [nvm](https://github.com/creationix/nvm) not to confuse nodejs's versions. You can find the details in order to download and install nodejs at https://nodejs.org/en/
* Also you should install [docker](https://docs.docker.com/engine/installation/) if you have problem running the app.


## Install Required Dependencies

To install required dependencies, just run the ```make initialize``` command. It fetches all golang dependencies and nodejs packages. Don't forget to add ```make```, ```go```, ```npm``` must be in **PATH** enviropment variables.


## Running

After installing required dependencies, now you can build and run the program. Then run this comman one by one:

* To build backend module run ```make build```
* To build and copy them into the public folder, execute ```make build-copy-frontend-files```

Now you can run the program with ```./twitterStreaming``` command. Execute ```./twitterStreaming -h``` to get help about program.

    twitterStreaming -v => Display version
    twitterStreaming -h => Display usage
    twitterStreaming -config=config.json => Set the config file
                   config.json
                   {
                       "Port":3000,
                       "Twitter" : {
                            	"ConsumerKey":"xxx"
                            	"ConsumerKeyConsumerSecret":"xxx"
                            	"ConsumerKeyAccessTokenKey":"xxx"
                            	"ConsumerKeyAccessTokenSecret":"xxx"
                        }
                   }


### Docker Run

You can run it with this command: ```docker run -it --dns 8.8.8.8 -p 3000:3000  rayyildiz/twitter-streaming```

### Configuration

Application looks for a config.json in current folder (defualt). However you can change the config file by ```./twitterStreaming -config=abc.json```

	{
	    "Port":3000,
	    "Twitter" : {
				"ConsumerKey":"xxx"
		      		"ConsumerKeyConsumerSecret":"xxx"
		      		"ConsumerKeyAccessTokenKey":"xxx"
		      		"ConsumerKeyAccessTokenSecret":"xxx"
		      }
	}



### Twitter Configuration

You need to register an application at (Twitter)[https://dev.twitter.com]. Just go to https://apps.twitter.com/app/new and create an application. After that you need **Consumer Key (API Key)** , **Consumer Secret (API Secret)** , **Access Token** , **Access Token Secret** . Also don't forget the change persmision your account as __Read only__

Change these values with yours in ```config.json```. There is a sample config information by default.

Now just execute ```make run``` and hit the http://localhost:3000 (if port is __3000__ in config.json)

## Alternative Way To Run Application

If you have problem and install the docker, you can execute ```make docker build``` and  ```make docker-run``` . This commans build and run the application at docker.

