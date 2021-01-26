ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

ifneq (,$(wildcard ./deployments/env/front.env))
	include ./deployments/env/front.env
	export
endif


up: build-docker-migrate build-docker-server-profile build-docker-server-feed-consumer build-docker-front
	env $(cat ./deployments/env/front.env | xargs) docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-monitoring.yml -f ./deployments/docker-compose-feed-consumer.yml  up -d

down:
	docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-monitoring.yml -f ./deployments/docker-compose-feed-consumer.yml down

up-api: build-docker-migrate build-docker-server-profile
	docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-monitoring.yml up -d db db-slave1 db-slave2 migration server-profile mysqld-exporter-master mysqld-exporter-slave1 cadvisor node-exporter prometheus grafana

up-messenger: build-docker-migrate build-docker-server-messenger
	docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-messenger.yml up -d db-msg-shard1 db-msg-shard2 migration-messenger1 migration-messenger2 server-messenger

down-messenger:
	docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-messenger.yml down

build-docker-migrate:
	docker build -t social/migrate ./migrate

build-docker-server-profile:
	docker build -t social/server-profile -f server/Dockerfile.profile ./server

build-docker-server-messenger:
	docker build -t social/server-messenger -f server/Dockerfile.messenger ./server

build-docker-server-feed-consumer:
	docker build -t social/feed-consumer -f server/Dockerfile.feed_consumer ./server

build-docker-front:
	docker build -t social/front ./front --build-arg API_SERVER_URL=${API_SERVER_URL}

build:
	cd server && \
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./.bin/social_server

lint:
	golangci-lint run ./server/...

include local.mk

.PHONY: build migrate
