PS_USER=admin
PS_PASSWORD=admin
PS_HOST=localhost
PS_PORT=1234
PS_DB=todo-list
PS_SSLMODE=disable

build:
	go build -o ./bin/todo-api

run: build
	./bin/todo-api

runapp:
	./bin/todo-api

migrate-up:
	goose -dir ./migrations postgres "postgres://${PS_USER}:${PS_PASSWORD}@${PS_HOST}:${PS_PORT}/${PS_DB}?sslmode=${PS_SSLMODE}" up

migrate-down:
	goose -dir ./migrations postgres "postgres://${PS_USER}:${PS_PASSWORD}@${PS_HOST}:${PS_PORT}/${PS_DB}?sslmode=${PS_SSLMODE}" down