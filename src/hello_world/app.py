import json
import boto3

s3 = boto3.resource('s3')


def lambda_handler(event, context):
    bucket = s3.Bucket('open-c-lifs-models')
    bucket.put_object(Key='input.json', Body=event.body)
    return {
        "statusCode": 200,
        "body": json.dumps({
            "message": event,
        }),
    }
