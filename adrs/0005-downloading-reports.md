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
  a `uid` query parameter, being a JSON object containing the filename and version ID, base64 encoded.
- The service checks the file is present in the S3 bucket.
- If it exists, a "Download" button is rendered on the page, which fetches the file on click.
- The file is then streamed to the user without being read using `io.Copy`.
- If it does not exist, an error message is displayed.

The filename and version ID are used as a composite key as we cannot guarantee the filename (which is the key in S3 buckets)
will be unique.

## Consequences

As this is being done in advance of the work to generate the reports and send the emails, there is a risk some of this 
process may need to be amended to adapt to those implementations.

The `uid` in the link is easily decoded, making it a potential vector for malicious use as the filename could be iterated
on (e.g. manipulating dates). However, this is mitigated by the version ID, which is an AWS-generated unique string. The
expiry (TTL) on the bucket is also set to seven days, adding an extra level of security for any leaked email links.
