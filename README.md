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

### Run tests with coverage

Run the following command:

```bash
make test_int_coverage
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

### Publisher

The publisher service will make API calls with HackerNews API to retrieve the stories and jobs. The items retrieved will
then be pushed to a RabbitMQ queue

### Consumer

The consumer will poll a RabbitMQ queue to store hacker news items. Once the message is read off RabbitMQ then the GRPC
server is called to save the item to the database

### API

The api reads data from a GRPC server and returns the necessary information based on the API path.

#### GRPC

The GRPC service support communication between services. This service is responsible for reading items either from a
redis cache or from the database, and saving items to the database

This is the single source to read/write data to data stores.

#### Updating the generated go files

If you update the `.proto` then you need to run the following command:

```bash
make proto
```