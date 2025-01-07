# 8. Merging Backends

Date: 2025-01-07

## Status

Accepted

## Context

Following ADR-0007, work began on building out data seeding functionality for the reporting feature, using the `finance-migration` 
and `finance-hub-api` services. However, this was getting increasingly complex and difficult to maintain, with very little
benefit from keeping the backends separate.

## Decision

The Finance Admin API will be merged into the Finance Hub API, starting with downloads and reports. This will simplify the
architecture and make it easier to test and maintain. The intention will be for this repository to become frontend-only in
the same fashion as the existing Golang microfrontends (e.g. Deputy Hub).

We will not be merging the frontends at the same time. The use case and user journeys for these two services are very different,
with Finance Hub being client-centric and Finance Admin being management-focused. This decision also leaves open the possibility
of using the Finance Admin frontend for wider management tasks in the future.

## Consequences

This should avoid issues with source of truth for data from a testing perspective, making it easier to maintain and test.
Additional testing may be required with regard to the API contract, but this can be mitigated by using the `shared` package
of JSON struct representations between services.

While downloads and reports are fairly straightforward and more pressing in relation to the testing difficulties, file uploads
will require more thought as it is more complex and event driven. Although moving this to the Finance Hub API will remove 
a lot of that complexity, it is not an immediate concern and can be addressed in a future iteration.
