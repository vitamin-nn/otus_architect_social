FROM golang:1.14-alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/social-server

FROM scratch

COPY --from=builder /go/bin/social-server /go/bin/social-server
ENTRYPOINT ["/go/bin/social-server"]
CMD ["profile"]
EXPOSE 8090
