FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy
RUN go test ./...
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./out/app ./server/cmd/main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/out/app .
COPY .env .
EXPOSE 8080
CMD ["./app"]
