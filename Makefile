.PHONY: dev build

dev:
	docker-compose up --build --force-recreate --rm

build:
	docker-compose build

