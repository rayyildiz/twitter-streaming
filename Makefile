# Initialize project, fetch required dependencies for golang and nodejs
initialize:
	go get golang.org/x/net/websocket
	go get gopkg.in/redis.v4
	go get github.com/cenkalti/backoff
	go get github.com/dghubble/oauth1
	go get github.com/google/go-querystring/query
	go get golang.org/x/net/websocket

	npm install -g create-react-app
	cd  ./frontend && npm update

# Deprecated (not need anymore, because it requires frontend)
# Install backend module
install:
	go install .

# build frontend project and copy them into the public folder
build-copy-frontend-files:
	rm -rf ./public
	cd ./frontend && npm run build
	cp -rf  "./frontend/build/" "./public"
	cd ..

# build backend module
build:
	go build .

# run whole project . First build and copy frontend module
run: build-copy-frontend-files
	go run main.go

# test backend module
test: build
	go test -v . ./twitter ./conf ./util ./api #./db  : db disable right now
	cd ./frontend && ./node_modules/.bin/flow check

# code format for backend module
fmt:
	gofmt -w *.go */*.go

tags:
	find ./ -name '*.go' -print0 | xargs -0 gotags > TAGS

push:
	git push origin master
	$(MAKE) stat

stat:
	git status

commit: stat
	git add .
	git commit -m "$(comment)"

docker-images:
	echo "Docker Images"
	docker images

docker-ps:
	echo "Docker All Container Passive/Active"
	docker ps -a

# clean docker images and delete container
docker-clean:
	$(MAKE) docker-ps
	docker ps -a | grep 'rayyildiz/twitter-streaming' | awk '{print $3}' | xargs docker rm -f
	docker images --no-trunc | grep none | awk '{print $3}' | xargs docker rmi -f
	$(MAKE) docker-ps

docker-images-clean:
	docker images --no-trunc | grep none | awk '{print $3}' | xargs docker rmi -f

# build for docker. It requires a clean build
docker-build: docker-images docker-ps build-copy-frontend-files
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o twitterStreaming .
	docker build . -t rayyildiz/twitter-streaming
	$(MAKE) docker-images
	$(MAKE) docker-ps

# deprecated ( not need right now. we can store words in redis later)
docker-run-redis:
	docker run --name twitter-streaming-db -p 6379:6379 -d redis

# run the docker image
docker-run:
	docker run -it --dns 8.8.8.8 -p 3000:3000  rayyildiz/twitter-streaming
