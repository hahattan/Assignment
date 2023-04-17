.PHONY: build tidy docker server client test

tidy:
	go mod tidy

build: client server

server:
	go build -o server/server server/main.go

client:
	go build -o client/client client/main.go

docker: docker_server docker_client

dserver: docker_server
docker_server:
	docker build -t grpc-server -f server/Dockerfile .

dclient: docker_client
docker_client:
	docker build -t grpc-client -f client/Dockerfile .

test:
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -n 1 | awk '{print $3}'