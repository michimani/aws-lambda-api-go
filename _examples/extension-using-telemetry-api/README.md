Telemetry API Example - The Lambda Extension using Telemetry API
===

# Prepare

## Build Extension

```bash
make build-ex
```

## Build Function

```bash
make build-func
```

## Create IAM Role for Lambda Function

```bash
make role
```

# Deploy

## Publish Extension as Lambda Layer

```bash
aws lambda publish-layer-version \
--layer-name 'telemetry-api-extension' \
--zip-file 'fileb://bin/extension.zip' \
--region ap-northeast-1
```

## Create Lambda Function

```bash
aws lambda create-function \
--function-name 'function-with-telemetry-api-extension' \
--runtime 'python3.9' \
--handler 'main.lambda_handler' \
--role $(
  aws iam get-role \
  --role-name 'telemetry-api-function-role' \
  --query 'Role.Arn' \
  --output text) \
--zip-file fileb://function.zip \
--layers $(
  aws lambda list-layer-versions \
  --layer-name 'telemetry-api-extension' \
  --query 'LayerVersions[0].LayerVersionArn' \
  --output text) \
--region ap-northeast-1
```

# Invoke

```bash
aws lambda invoke \
--function-name 'function-with-telemetry-api-extension' \
--invocation-type 'RequestResponse' \
--cli-binary-format 'raw-in-base64-out' \
--region 'ap-northeast-1' \
--log-type 'Tail' \
/dev/stdout \
| jq -sr '.[1] | .LogResult' \
| base64 -d
```