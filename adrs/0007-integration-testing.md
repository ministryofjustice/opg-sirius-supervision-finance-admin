# 7. Integration Testing

Date: 2024-12-12

## Status

Accepted

## Context

Finance Admin queries reports against the `supervision_finance` database schema. However, is not the source of truth for this
data and does not manage its migrations. It is also not the entrypoint for that data. This means our existing approach of
seeding data - writing it with SQL queries - is liable to produce false positives, as it would be unaware of any changes to
the data or business logic used to create it.

## Decision

Testcontainers is used to stand up a Docker Compose stack consisting of the database (`supervision_finance` and a small subset
of `public`), the `finance-migration` container, and `finance-hub-api`. Data is seeded by calling the API, using the JSON
data structures imported from the `shared` package of Finance Hub, ensuring data is created in the same way as in production.

## Consequences

An unintended, but welcome, consequence of this approach is the integration tests also assert on the API contract, in the 
form of the `shared` package. This interdependence does make it more important to ensure dependencies are up-to-date, but
as the services are developed by the same team, this is not a significant concern.

Another potential issue is that although these tests will catch breaking changes to both the API and the database schema, they
won't be able to do so until the tests have been run, which is hard to do currently when the breaking changes will be made
in another project. This will be discussed in the future, along with more general end-to-end testing strategies.