AWS Lambda API for Go
===

[![codecov](https://codecov.io/gh/michimani/aws-lambda-api-go/branch/main/graph/badge.svg?token=P63U316Y2U)](https://codecov.io/gh/michimani/aws-lambda-api-go)

This is a client library for Go language to use AWS Lambda's Runtime API, Extension API, Telemetry API, and Logs API.

# Supported APIs

## Runtime API

[AWS Lambda runtime API - AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html)

- [x] `GET /runtime/invocation/next`
- [ ] `POST /runtime/invocation/:AwsRequestId/response`
- [ ] `POST /runtime/init/error`
- [ ] `POST /runtime/invocation/:AwsRequestId/error`

## Extension API

- TODO

## Telemetry API

- TODO

# License

[MIT](https://github.com/michimani/aws-lambda-api-go/blob/main/LICENSE)

# Author

[michimani210](https://twitter.com/michimani210)
