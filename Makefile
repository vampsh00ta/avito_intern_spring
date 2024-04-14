DATABASE="postgresql://avito_intern:avito_intern@localhost:5432/avito_intern?sslmode=disable"
migrate:
	migrate create -ext sql -dir ./migrations/ -seq $(name)
migration:
	migrate -path ./migrations -database  $(DATABASE)  up
start:
	docker-compose build &&  docker-compose up

test:
	go test ./internal/transport/http/tests
test-docker:
	docker exec -it vampshoota_app_avito make test


