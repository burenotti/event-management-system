generate-docs:
	swag init -g handler/handler.go

generate: generate-docs

test: generate
	docker-compose run db-tests
	docker-compose run --rm migrate-tests
	docker-compose run --rm tests
	docker-compose rm db-tests

migrate:
	docker-compose run --rm migrate

rollback:
	docker-compose run --rm rollback

up: generate migrate
	docker-compose up

dev: generate migrate
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up