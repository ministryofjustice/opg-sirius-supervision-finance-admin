# 2. Create a new Finance Admin Front End

Date: 2024-03-14

## Status

Accepted

## Context

As part of their role, the Finance team need to import and export data from Sirius in order to reconcile payments. This 
is currently done via the existing Finance Admin screen. However, many of those functions are now changing as a result of
the source of truth for payments moving to us from SOP. We are therefore replacing the existing page with a new one.

## Decision

This will be a new Front End service following our established Golang FE pattern. The web pages it serves will have two 
separate functions: Uploading files for processing, and downloading reports. Reporting will be broken out into a separate
microservice, as it has a defined scope and will allow us to provide resources for potentially long-running database queries
without impacting other services. The architecture for processing files will depend on the file and how coupled the data
are with existing processes.

## Consequences

Splitting the Front End in this way allows us to work on the product independently and not be blocked by architectural decisions.
However, with all microservices, it adds to the overall system complexity. Our contract testing capability is fairly immature
so care will be required to ensure all services work together as expected.
