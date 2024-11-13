# 4. Authentication

Date: 2024-10-29

## Status

Accepted

## Context

Users accessing Finance Admin are able to request and download payments data, and as a result, we need to ensure they are
authenticated in our system and authorised to perform the action. There is a need to ensure that some level of auth is
in place before the report download functionality is in place, as that is the greatest risk to data loss or unauthorised 
access.

## Decision

The simplest solution to authentication is to piggyback on the existing Sirius user session and refresh the session on
each API call by fetching the current user from the Sirius API. In the event that this session is not valid, we redirect
to the login page. This is the same method by which we authenticate in Finance Hub.

## Consequences

This just solves authentication, not authorisation, which means that any user with a Sirius session can access Finance 
Admin. This will be resolved in a future change. 

We also have no auth between the front end and back end services, and while this is restricted by security groups in AWS, 
we should still consider some authentication between the two.

Lastly, this relies on the Sirius API and session handling, which means there is a risk of authentication failing if the 
underlying API is ever changed.
