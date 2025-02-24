all: go-lint test build-all scan cypress down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results cypress/screenshots .trivy-cache .go-cache

setup-directories: test-results

go-lint:
	docker compose run --rm go-lint

gosec: setup-directories
	docker compose run --rm gosec

build:
	docker compose build --no-cache --parallel finance-admin finance-admin-api

build-dev:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build --parallel finance-admin finance-admin-api yarn json-server finance-hub-api

build-all:
	docker compose build --parallel finance-admin finance-admin-api yarn cypress

test: setup-directories
	go run gotest.tools/gotestsum@latest --format testname  --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

download-trivy-dbs: download-trivy-db download-trivy-java-db
download-trivy-db:
	docker compose run trivy image --download-db-only
download-trivy-java-db:
	docker compose run trivy image --download-java-db-only


scan: scan-hub scan-api
scan-hub: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-finance-admin:latest
	docker compose run --rm trivy image --format sarif --output /test-results/admin-hub.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-finance-admin:latest
scan-api: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-finance-admin-api:latest
	docker compose run --rm trivy image --format sarif --output /test-results/admin-api.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-finance-admin-api:latest

clean:
	docker compose down
	docker compose run --rm yarn

up: clean compile-assets build-dev
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up finance-admin finance-admin-api finance-hub-api yarn

down:
	docker compose down

compile-assets:
	docker compose run --rm yarn build

cypress: setup-directories
	docker compose up -d localstack
	docker compose run --build cypress
