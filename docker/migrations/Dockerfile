FROM golang:1.21-alpine

RUN apk --no-cache add curl jq git

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY db/migrations /migrations

ENTRYPOINT ["goose"]
