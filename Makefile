docker-build:
	@make --no-print-directory codegen
	docker build -t pineapple217/mb:latest --build-arg GIT_COMMIT=$(shell git log -1 --format=%h) . 

docker-push:
	docker push pineapple217/mb:latest

docker-update:
	@make --no-print-directory docker-build
	@make --no-print-directory docker-push

codegen:
	templ generate
	sqlc generate
	go run ./cmd/gen/

build:
	@make --no-print-directory codegen
	go build -o ./tmp/main.exe ./cmd/server

start:
	@./tmp/main.exe

run: 
	@make --no-print-directory build
	@make --no-print-directory start