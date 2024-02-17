FROM golang:1.21.6

ENV GIN_MODE=release

WORKDIR /app

COPY . .

RUN go build -o ./target/newsletter ./cmd/newsletter/main.go

ENTRYPOINT [ "./target/newsletter" ]