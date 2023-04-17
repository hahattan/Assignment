# gRPC Client-Server Testing Tool

This repository holds the gRPC client-server example code and provides various deployment option.

## Build

### Native binaries

#### Prerequisites

- [Go 1.20](https://go.dev/doc/install)

Use the `Makefile` in the root directory:
```shell
make build
```

### Build your own Docker Images
In addition to running the services directly, Docker Container can be used.

#### Prerequisites

- [Docker](https://docs.docker.com/engine/install/)

Use the `Makefile` in the root directory:
```shell
make docker
```

## Usage 

### Server

gRPC server, listening on 50051 port, accepts gRPC request from gRPC clients.
```shell
$ make server
go build -o server/server server/main.go

$ ./server/server
2023/04/13 20:50:05 server listening at [::]:50051
```

### Client

gRPC client calls gRPC function of the server, and reports statistics periodically.
```shell
$ make client
go build -o client/client client/main.go

$ ./client/client -h
Usage of ./client/client:
  -addr string
        the address to connect to (default "localhost:50051")
  -freq int
        request frequency in millisecond (default 5000)
  -name string
        Name to greet (default "world")
  -number int
        concurrent request number (default 1)

```

## Deployment

### Docker Compose

#### Prerequisites

- [Docker Compose](https://docs.docker.com/compose/install/)

Run docker compose command from the root directory:
```shell
$ docker compose -f deployment/docker-compose.yml up -d
[+] Running 3/3
 ⠿ Network bluex_default  Created                                                                                                          0.1s
 ⠿ Container grpc-server  Started                                                                                                          0.8s
 ⠿ Container grpc-client  Started                                                                                                          0.8s
```

To clean up:
```shell
docker compose -f deployment/docker-compose.yml down -v 
```

### Kubernetes

#### Prerequisites

- [minikube](https://minikube.sigs.k8s.io/docs/start/)

1. Start a local kubernetes cluster
```shell
minikube start
```

2. Point your shell to minikube's docker-daemon
```shell
eval $(minikube -p minikube docker-env)
```

3. Build Docker images
```shell
make docker
```

4. Create server resources
```shell
kubectl create -f deployment/server.yml
```

5. Cerate client resource
```shell
kubectl create -f deployment/client.yml
```

To clean up:
```shell
minikube delete
```