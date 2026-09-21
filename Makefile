all: go-lint gosec test build-all scan cypress down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results cypress/screenshots .go-cache

setup-directories: test-results

go-lint:
	docker compose run --rm go-lint

gosec: setup-directories
	docker compose run --rm gosec

build:
	docker compose build --no-cache --parallel finance-admin

build-dev:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build --parallel finance-admin npm json-server finance-hub-api

build-all:
	docker compose build --parallel finance-admin npm cypress

test: setup-directories
	go run gotest.tools/gotestsum@latest --format testname  --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

clean:
	docker compose down
	docker compose run --rm npm

up: clean compile-assets build-dev
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up finance-admin finance-hub-api npm

down:
	docker compose down

compile-assets:
	docker compose run --rm npm run build

cypress: setup-directories
	docker compose up -d localstack
	docker compose run --build cypress

cypress-single: setup-directories
	docker compose up -d localstack
	docker compose run cypress --spec cypress/e2e/$(spec)