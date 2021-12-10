ARG go_version=1.17

# User
FROM alpine:3.13.1 as user
ARG uid=10001
ARG gid=10001
RUN echo "scratchuser:x:${uid}:${gid}::/home/scratchuser:/bin/sh" > /scratchpasswd

# Certs
FROM alpine:3.13.1 as certs
RUN apk add -U --no-cache ca-certificates

# Build
FROM golang:${go_version} as build
WORKDIR /code/
ENV GO111MODULE=on
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY ./pkg/ ./pkg/

RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/consumer ./cmd/consumer

FROM scratch as consumer
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=user /scratchpasswd /etc/passwd
COPY --from=build /code/bin/consumer .
ENTRYPOINT ["./consumer"]