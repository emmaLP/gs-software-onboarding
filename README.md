# Gymshark Software Onboarding Project

This project helps new software engineers get familiar with coding standards and help gain a better understanding of how
GS software team work.

## Building/Testing

This project requires docker to be running on the machine as integration tests
leverage [testcontainers](https://github.com/testcontainers/testcontainers-go)
which will start&stop docker containers as part of the tests

### Run all tests

Run the following command:

```bash
go test ./...
```

## Services

This is a mono repo that contains multiple services.

To run all services:

```bash
docker-compose up -d
```

```bash
docker-compose down --remove-orphans
```

### Consumer

The consumer will periodically run to seed a mongo database from the HackerNews API.

### API

The api reads data from a mongo database and returns the necessary information based on the API path.

#### GRPC

The GRPC service support communication between services

#### Updating the generated go files

If you update the `.proto` then you need to run the following command:

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/grpc/proto/hackernews.proto 
```
