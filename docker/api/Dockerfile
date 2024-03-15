FROM golang:1.21 as build

WORKDIR /tmp/src

COPY go.mod go.sum ./
RUN go mod download

RUN git config --global --add safe.directory /tmp/src
RUN go install github.com/silenceper/gowatch@latest

COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/api cmd/api/main.go
CMD ["/app/bin/api"]

FROM gcr.io/distroless/static-debian12
COPY --from=build /app/bin/api /
CMD ["/api"]
