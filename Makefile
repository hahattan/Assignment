.PHONY: build test clean

build:
	go build -o app main.go

test:
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -n 1 | awk '{print $3}'

clean:
	rm -f app