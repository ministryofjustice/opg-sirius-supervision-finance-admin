all: go-lint test build-all scan cypress down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results cypress/screenshots .trivy-cache .go-cache

setup-directories: test-results

go-lint:
	docker compose run --rm go-lint

build:
	docker compose build --no-cache --parallel finance-admin finance-admin-api

build-dev:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build --parallel finance-admin finance-admin-api yarn json-server

build-all:
	docker compose build --parallel finance-admin finance-admin-api yarn cypress

test: setup-directories
	go run gotest.tools/gotestsum@latest --format testname  --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-finance-admin:latest
	docker compose run --rm trivy image --format sarif --output /test-results/hub.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-finance-admin:latest

clean:
	docker compose down
	docker compose run --rm yarn

up: clean compile-assets build-dev
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up finance-admin finance-admin-api yarn

down:
	docker compose down

compile-assets:
	docker compose run --rm yarn build

cypress: setup-directories clean
	docker compose run --build cypress
