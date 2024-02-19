# Builder stage
FROM golang:1.21.6 as builder

WORKDIR /app

COPY . .

RUN go build -o ./target/newsletter ./cmd/newsletter/main.go

# Runtime stage
FROM golang:1.21.6 as runner

WORKDIR /app

COPY --from=builder /app/target/newsletter  newsletter

COPY configuration configuration

ENV GIN_MODE=release
ENV APP_ENVIRONMENT=production

ENTRYPOINT [ "./newsletter" ]