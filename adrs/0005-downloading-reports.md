# 5. Downloading reports

Date: 2024-10-18

## Status

Accepted

## Context

Once a report has been generated and stored in S3, the intention is that an email will be sent to the user with a download 
link. Clicking this link should authenticate the user and download the file to their system. This should be done in an 
efficient manner, i.e. not reading the file into memory at every stage.

## Decision

The download flow is as follows:
- The URL in the email navigates to Finance Admin, which now requires a valid user session (see ADR 004). The URL contains
  a `uid` query parameter, being the filename in base64.
- The user is presented with a Download button, which on click sends an HTMX request to the server.
- This returns with an `HX-Redirect` header, that sends a further request to fetch the file.
- The file is then streamed to the user without being read using `io.Copy`.

## Consequences

As this is being done in advance of the work to generate the reports and send the emails, there is a risk some of this 
process may need to be amended to adapt to those implementations.

The UID in the link is also easily decoded, making it a potential vector for malicious use. While the scope of attack is
very limited, seeing as files will have a short time-to-live and users will need a valid session to download, but we can 
still iterate on this solution to mitigate this, either by properly encrypting the filename or storing generated UIDs in
a datastore.
