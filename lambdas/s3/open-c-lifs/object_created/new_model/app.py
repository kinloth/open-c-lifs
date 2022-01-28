import json
import boto3
import os

s3 = boto3.resource('s3')

def get_object(bucket_name: str, key: str):
    bucket = s3.Bucket(bucket_name)
    obj = bucket.Object(key)
    return obj.get()['Body'].read().decode('utf-8')

def lambda_handler(event, context):
    print('####################################')
    print(f'event = {event}')
    print('####################################')

    record = event['Records'][0]['s3']
    data = get_object(record['bucket']['name'], record['object']['key'])
    print(f'data = {data}')

    data = json.loads(data)
    print(f'data = {data}')

    # TODO: do something with data