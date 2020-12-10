FROM golang:1.14-alpine

RUN apk add git
RUN apk add gcc
RUN apk add g++
RUN go get -u github.com/pressly/goose/cmd/goose

RUN mkdir -p /opt/go/app/db
COPY . /opt/go/app/db/

WORKDIR /opt/go/app/db/

RUN chmod u+x /opt/go/app/db/migrate.sh

CMD ["sh", "-c", "/opt/go/app/db/migrate.sh $SERVICE_NAME"]
