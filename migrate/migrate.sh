#!/bin/sh

DSN="$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_DB_HOST:$MYSQL_PORT)/$MYSQL_DATABASE?parseTime=true"

for i in $(seq 1 5); do
    goose -dir migrations mysql $DSN up && break
    sleep 1
done
