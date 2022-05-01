.PHONY: tidy build run test docker

SERVER=server

build: $(SERVER)

tidy:
	go mod tidy

server: tidy
	go build -o $@ main.go

run:
	./server

test:
	go test -coverprofile=coverage.out ./...

docker:
	docker build -t local/assignment:0.0.0 .