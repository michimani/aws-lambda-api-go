AWS Lambda API for Go
===

[![codecov](https://codecov.io/gh/michimani/aws-lambda-api-go/branch/main/graph/badge.svg?token=P63U316Y2U)](https://codecov.io/gh/michimani/aws-lambda-api-go)

This is a client library for Go language to use AWS Lambda's Runtime API, Extension API and Telemetry API.

# Supported APIs

## Runtime API

[AWS Lambda runtime API - AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html)

- [x] `GET /runtime/invocation/next`
- for custom runtime
  - [x] `POST /runtime/invocation/:AwsRequestId/response`
  - [ ] `POST /runtime/init/error`
  - [ ] `POST /runtime/invocation/:AwsRequestId/error`

## Extension API

[Lambda Extensions API - AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-extensions-api.html)

- [x] `POST /extension/register`
- [x] `GET /extension/event/next`
- [ ] `POST /extension/init/error`
- [ ] `POST /extension/exit/error`

## Telemetry API

[Lambda Telemetry API reference - AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/telemetry-api-reference.html)

- [x] `PUT /telemetry`

# License

[MIT](https://github.com/michimani/aws-lambda-api-go/blob/main/LICENSE)

# Author

[michimani210](https://twitter.com/michimani210)
