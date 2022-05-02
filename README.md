# assignment

## Installation and deployment options 

### Native binary

#### Prerequisites

- Go 1.18
- Redis DB

#### Installation and Execution
```shell
git clone -b PulseiD --single-branch git@github.com:hahattan/assignment.git
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
git clone -b PulseiD --single-branch git@github.com:hahattan/assignment.git
cd assignment
make docker
docker-compose up -d
```

## Usage

### Admin API

Default basic authentication credentials for Admin API is `root:root`.  
Overwrite them with environment variables `ADMIN_USERNAME` and `ADMIN_PASSWORD`

### Database

Web server default connecting to local Redis DB (localhost:6379).  
Overwrite the address with environment variables `DB_HOST` and `DB_PORT`