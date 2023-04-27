generate-docs:
	swag init -g handler/handler.go


test:
	go test -count=1 -race -coverprofile=coverage.out ./...; \
	go tool cover -html=coverage.out
	rm coverage.out

migrate:
	docker-compose run --rm migrate

rollback:
	docker-compose run --rm rollback

up: migrate
	docker-compose up

dev: migrate
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up