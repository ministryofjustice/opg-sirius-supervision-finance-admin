# 3. Create a new Finance Admin Back End

Date: 2024-09-24

## Status

Accepted

## Context

We have created the Finance Admin front end, however the decision was made to add a service between S3 and the front end
to upload and download files from S3 without making AWS credentials accessible to the user.

## Decision

A new back end will be created in this repo, similarly to the structure of the finance hub. This will communicate with S3 
to perform uploads and downloads using Go's AWS SDK. This back end will also allow us to perform any other logic required for
the finance admin system, such as validation. 

## Consequences

This ensures we can interact with S3 whilst keeping our AWS credentials away from the client, but increases the complexity of the system. 

