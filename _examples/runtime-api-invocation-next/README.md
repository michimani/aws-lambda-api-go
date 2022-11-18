Runtime API Example - GET /runtime/invocation/next
===

# Usage

1. Start Lambda

    ```bash
    make build-run
    ```
    
    If you have already built the image, you can start Lambda with the following command
    
    ```bash
    make run
    ```

2. Invoke function

    ```bash
    curl -X POST \
    -d '{"type": "test", "message": "test event"}' \
    'http://localhost:9000/2015-03-31/functions/function/invocations'
    ```