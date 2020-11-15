ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

up: build-docker-migrate build-docker-server build-docker-front
	#ifneq (,$(wildcard ./deployments/env/front.env))
	#	include ./deployments/env/front.env
	#	export
	#endif
	env $(cat ./deployments/env/front.env | xargs) docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-monitoring.yml  up -d

down:
	docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-monitoring.yml down

up-api: build-docker-migrate build-docker-server
	docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-monitoring.yml up -d db db-slave1 db-slave2 migration server mysqld-exporter-master mysqld-exporter-slave1 cadvisor node-exporter prometheus grafana

build-docker-migrate:
	docker build -t social/migrate ./migrate

build-docker-server:
	docker build -t social/server -f server/Dockerfile.server ./server

build-docker-front:
	docker build -t social/front ./front --build-arg API_SERVER_URL=${API_SERVER_URL}

build:
	cd server && \
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./.bin/social_server

lint:
	golangci-lint run ./server/...

include local.mk

.PHONY: build migrate
