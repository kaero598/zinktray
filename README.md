# Fake SMTP Server

This application provides basic SMTP server functionality for email testing.
All emails sent are stored in memory and may be retrieved via API.

## Features

* Anonymous and authenticated email sending.
* API to retrieve all registered mailboxes.
* API to retrieve all stored messages.
* API to retrieve raw message contents.

## Usage

To launch the application use `go run` command:

```shell
$ go run .
```

Or use `go build` and execute `go-fake-smtp` binary.

No configuration options are exposed yet.

SMTP server is exposed on port `2525`. HTTP server and API are exposed on `localhost:8080`.

## API

To retrieve stored messages make an HTTP request to API endpoint `http://localhost:8080/api/messages`. The endpoint returns JSON-encoded list of stored messages, each with a single field containing raw email contents along with headers and body as sent via SMTP session.
