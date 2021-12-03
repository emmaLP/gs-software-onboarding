# Gymshark Software Onboarding Project

This project helps new software engineers get familiar with coding standards and help gain a better understanding of how
GS software team work.

## Services

This is a mono repo that contains the following services:

### Consumer

The consumer will periodically run to seed a mongo database from the HackerNews API.

To run the consumer:

```bash
cp .env.example app.env && go run cmd/consumer/main.go
```

Remember to remove the `app.env` file:

```bash
rm app.env
```