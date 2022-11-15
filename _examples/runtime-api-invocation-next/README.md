Runtime API Example - GET /runtime/invocation/next
===

# Usage

0. Build image

    ```bash
    docker build -t runtime-api-invocation-next:local .
    ```

1. Run the image

    ```bash
    docker run \
    --rm \
    -p 9000:8080 \
    runtime-api-invocation-next:local
    ```

2. Invoke function

    ```bash
    curl -X POST \
    -d '{"type": "test", "message": "test event"}' \
    'http://localhost:9000/2015-03-31/functions/function/invocations'
    ```