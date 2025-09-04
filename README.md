# OPG SIRIUS SUPERVISION FINANCE ADMIN

### Major dependencies

- [Go](https://golang.org/) (>= 1.22)
- [docker compose](https://docs.docker.com/compose/install/) (>= 2.26.0)

#### Installing dependencies locally:
(This is only necessary if running without docker)

- `yarn install`
- `go mod download`
---

## Local development

The application ran through Docker can be accessed on `localhost:8889/finance-admin/downloads`.

To enable debugging and hot-reloading of Go files:

`make up`

Hot-reloading is managed independently for both apps and should happen seamlessly. Hot-reloading for web assets (JS, CSS, etc.)
is also provided via a Yarn watch command.

-----
## Run the unit/integration tests

`make test`

## Run Cypress tests headless

`make cypress`

Finance admin pulls in the finance hub container to run the cypress tests, so if behaviour differs with your tests across your local environment and build pipeline then you might need to run `docker compose pull finance-hub-api` to pull in the latest changes.

## Run Cypress tests with UI

`make up`
`npx cypress open -c baseUrl=https://localhost:8888/supervision/finance-admin`

## Run Trivy scanning

`make scan`

-----
## Architectural Decision Records

The major decisions made on this project are documented as ADRs in `/adrs`. The process for contributing to these is documented
in the first ADR.

-----
## HTMX & JS

This project uses [HTMX](https://htmx.org/) to render partial HTML instead of reloading the whole page on each request. However, this can 
mean that event listeners added on page load may fail to register/get deregistered when a partial is loaded. To avoid this,
you can force event listeners to register on every HTMX load event by putting them within the `htmx.onLoad` function.

HTMX also includes a range of utility functions that can be used in place of more unwieldy native DOM functions.
