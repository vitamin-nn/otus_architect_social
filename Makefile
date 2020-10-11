ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

ifneq (,$(wildcard ./deployments/env/front.env))
    include ./deployments/env/front.env
    export
endif

up: build-docker-migrate build-docker-server build-docker-front
	env $(cat ./deployments/env/front.env | xargs) docker-compose -f ./deployments/docker-compose.yml up -d

down:
	docker-compose -f ./deployments/docker-compose.yml down

build-docker-migrate:
	docker build -t social/migrate ./migrate

build-docker-server:
	docker build -t social/server -f server/Dockerfile.server ./server

build-docker-front:
	docker build -t social/front ./front --build-arg API_SERVER_URL=${API_SERVER_URL}

build:
	cd server && \
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./.bin/social_server

run-server:
	source ./configs/.local.env && cd server && go run . server

lint:
	golangci-lint run ./server/...

include local.mk

.PHONY: build migrate
