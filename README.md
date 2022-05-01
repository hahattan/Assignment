# assignment

## Installation and deployment options 

### Native binary

#### Prerequisites

- Go 1.18
- Redis DB

#### Installation and Execution
```shell
git clone git@github.com:hahattan/assignment.git
cd assignment
make build
make run
```

### Build your own Docker Container

#### Prerequisites

- Docker
- [Docker Compose](https://docs.docker.com/compose/)

#### Build and Execution
```shell
git clone git@github.com:hahattan/assignment.git
cd assignment
make docker
docker-compose up -d
```