ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

up: build-docker-migrate build-docker-server build-docker-front
	docker-compose -f ./deployments/docker-compose.yml up -d

down:
	docker-compose -f ./deployments/docker-compose.yml down

build-docker-migrate:
	docker build -t social/migrate ./migrate

build-docker-server:
	docker build -t social/server ./server

build-docker-front:
	docker build -t social/front ./front

build:
	cd server && \
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./.bin/social_server

run-server:
	source ./configs/.local.env && cd server && go run . server

migrate:
	goose -dir migrate/migrations mysql "otus_social:otus_social_passwd@tcp(0.0.0.0:3306)/otus_social?parseTime=true" up

lint:
	golangci-lint run ./server/...

.PHONY: build migrate
