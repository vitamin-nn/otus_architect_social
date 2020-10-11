# This dockerfile was created for heroku deploy.
# For other reasons use Dockerfile.server.

FROM golang:1.14-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o social-server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/social-server .
CMD ["./social-server", "server"]
