ifneq (,$(wildcard ./deployments/env/prod.database.env))
    include ./deployments/env/prod.database.env
    export
endif

migrate-heroku:
	goose -dir migrate/migrations mysql "${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_DB_HOST})/${MYSQL_DATABASE}?parseTime=true" up

migrate:
	goose -dir migrate/profile mysql "otus_social:otus_social_passwd@tcp(0.0.0.0:3306)/otus_social?parseTime=true" up

run-profile-server:
	source ./configs/.local.env && source ./configs/.local.profile.env && cd server && go run . profile

run-feed-consumer:
	source ./configs/.local.env && source ./configs/.local.profile.env && cd server && go run . feed-consumer

run-db:
	docker-compose -f ./deployments/docker-compose.yml up -d db db-slave1 db-slave2

run-messenger-server:
	source ./configs/.local.env && source ./configs/.local.messenger.env && cd server && go run . messenger

run-db-messenger:
	docker-compose -f ./deployments/docker-compose-messenger.yml up -d db-msg-shard1 db-msg-shard2

migrate-messenger:
	goose -dir migrate/migrations mysql "otus_social:otus_social_passwd@tcp(0.0.0.0:6306)/otus_social?parseTime=true" up
	goose -dir migrate/migrations mysql "otus_social:otus_social_passwd@tcp(0.0.0.0:7306)/otus_social?parseTime=true" up
