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

This is a mono repo that contains the following services:

### Consumer

The consumer will periodically run to seed a mongo database from the HackerNews API.

To run the consumer:

The command below relies on a running mongo database running

```bash
cp .env.example app.env && go run cmd/consumer/main.go
```

Remember to remove the `app.env` file once you are done:

```bash
rm app.env
```

To run a fully functional consumer use the following command:

```bash
docker-compose up -d
```

```bash
docker-compose down --remove-orphans
```
