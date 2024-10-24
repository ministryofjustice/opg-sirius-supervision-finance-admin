# 3. Downloading reports

Date: 2024-10-18

## Status

Accepted

## Context

Once a report has been generated and stored in S3, the intention is that an email will be sent to the user with a download 
link. Clicking this link should authenticate the user and download the file to their system. This should be done in an 
efficient manner, i.e. not reading the file into memory at every stage.

## Decision

- io.Copy for data streaming
- auth
- JWT
- UX

## Consequences

As this is being done in advance of the work to generate the reports and send the emails, there is a risk some of this 
process may need to be amended to adapt to those implementations.
