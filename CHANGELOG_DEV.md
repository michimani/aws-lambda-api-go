CHANGELOG 
===
This is a version of CHANGELOG less than v1.0.0

## [Unreleased]

* Runtime API
  * `POST /runtime/init/error`
  * `POST /runtime/invocation/:AwsRequestId/error`
* Extension API
  * `POST /extension/init/error`
  * `POST /extension/exit/error`

v0.2.0 (2022-12-10)
===

* Support `PUT /telemetry`

v0.1.2 (2022-12-04)
===

* Fix request header for `POST /extension/register` 

v0.1.1 (2022-12-04)
===

* Support `POST /extension/register`
* Support `GET /extension/event/next`

v0.1.0 (2022-11-19)
====

* Support `POST /runtime/invocation/:AwsRequestId/response`

v0.0.0 (2022-11-15)
====

* pre release