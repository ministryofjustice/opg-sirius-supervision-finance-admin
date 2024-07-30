# 2. Adopt Templ for UI components

Date: 2024-08-02

## Status

Accepted

## Context

We have used Golang's native HTML templating in all our microfrontends to date, but there are some frustrations. The
Handlebars syntax does allow for some conditional logic and template composition, but it is incomplete. Templates accept
variables but only as an empty interface, so there is no type safety or intellisense, so a simple misspelled variable 
will result in runtime errors. They are also hard to test, meaning our Cypress tests have become page-level component 
tests, rather than system-level end-to-end tests.

## Decision

[Templ](https://templ.guide/) is a widely adopted HTML component generator where components are described as Go functions.
The syntax uses valid HTML, meaning converting our existing designs and templates should be fairly straight-forward, but
allows full use of everything Go has to offer, including typed function parameters. As components are functions, they can
be passed to other components as function parameters. The `.templ` are then transpiled to `.go` files in a compile step,
and then can be imported into other packages.

This also allows for a better separation of state. Previously, we have run into issues where templates used the same structs
as were used for the API, which coupled the two use cases and made it hard to modify (not to mention being a common cause
of runtime errors). With Templ, we can keep all view state in the `components` package, with the separation enforced by
Go not allowing circular dependencies between packages.

## Consequences

Adopting any library outside of the standard library brings risk. However, Templ is very widely adopted, especially when
paired with HTMX, and in active development. It is also another technology to learn, though is arguably easier to understand
in regard to component composition.
