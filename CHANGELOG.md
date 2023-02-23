# Changelog

`aws-whoami` uses [monotonic versioning](https://github.com/benkehoe/monotonic-versioning-manifesto) since v1.0 across both [the older Python implementation](https://github.com/benkehoe/aws-whoami) (compatibility number 1) and this Go implementation (compatibility number 2).

## v2.6

* With `--json`, errors are printed as JSON in the form `{"Error": "The error message"}`

## v2.5

* Add `--disable-account-alias` flag.
* Handle paths in user ARNs.
* Disable account alias check by matching SSO Permission Set name.
* Internal revamp for testing.
* Change repo layout to work better with `go install`.
* Add tests.

## v2.4

* Handle root user

## v2.3

* Initial implementation in Go.
* `--debug` has been removed.
* A region is now required (this appears to be an SDK constraint).

## v1.2

Latest version of [the Python implementation](https://github.com/benkehoe/aws-whoami).

