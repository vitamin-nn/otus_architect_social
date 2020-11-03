ifneq (,$(wildcard ./deployments/env/prod.database.env))
    include ./deployments/env/prod.database.env
    export
endif

migrate-heroku:
	goose -dir migrate/migrations mysql "${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_DB_HOST})/${MYSQL_DATABASE}?parseTime=true" up

migrate:
	goose -dir migrate/migrations mysql "otus_social:otus_social_passwd@tcp(0.0.0.0:3306)/otus_social?parseTime=true" up

run-server:
	source ./configs/.local.env && cd server && go run . server

run-db:
	docker-compose -f ./deployments/docker-compose.yml up -d db