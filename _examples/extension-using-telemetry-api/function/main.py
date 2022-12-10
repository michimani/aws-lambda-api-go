import json
import logging


def lambda_handler(event, context):
    logger = logging.getLogger()
    logger.setLevel(logging.INFO)

    logger.info('[%s] Hello Telemetry API!', context.aws_request_id)

    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Lambda!')
    }
