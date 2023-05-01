generate-docs:
	swag init -g handler/handler.go


test-native:
	source $(ENV_FILE); \
	DB_DSN=$${DB_DSN} \
	MIGRATE_DSN=$${MIGRATE_DSN} \
	MIGRATE_DIR=$${MIGRATE_DIR} \
	go test -v -count=1 -race -coverprofile=coverage.out ./...; \
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

migrate:
	docker-compose -f docker-compose.dev.yml --env-file=$(ENV_FILE) run --rm migrate

rollback:
	docker-compose -f docker-compose.dev.yml --env-file=$(ENV_FILE) run --rm rollback

up: migrate
	docker-compose up

dev: migrate
	docker-compose -f docker-compose.dev.yml up

test:
	docker-compose -f docker-compose.dev.yml --env-file=$(ENV_FILE) run --rm test;