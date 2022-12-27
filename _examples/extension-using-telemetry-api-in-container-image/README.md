Telemetry API Example - The Lambda Extension using Telemetry API in container image
===

# Preparing

1. Login to ECR

    ```bash
    REGION='ap-northeast-1'
    AWS_ACCOUNT_ID=$(
      aws sts get-caller-identity \
      --query 'Account' \
      --output text) \
    && aws ecr get-login-password \
      --region "${REGION}" \
      | docker login \
      --username AWS \
      --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com
    ```

2. Build and zip extension

    ```bash
    make build-ex
    ```

3. Build function image

    ```bash
    make build-func
    ```
    
# Invoke function

1. Run Lambda

    ```bash
    make run
    ```

2. Invoke function

    ```bash
    curl \
    -H 'Content-Type: application/json' \
    http://localhost:9000/2015-03-31/functions/function/invocations
    ```

# Author

[michimani210](https://twitter.com/michimani210)