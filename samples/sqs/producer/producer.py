from os import getenv
from time import sleep

import boto3
from botocore.exceptions import NoRegionError


def check_region():
    checks = [
        getenv('AWS_REGION'),
        getenv('AWS_DEFAULT_REGION'),
        boto3.DEFAULT_SESSION.region_name if boto3.DEFAULT_SESSION else None,
        boto3.Session().region_name,
    ]

    for region in checks:
        if region:
            return region


def main():
    region = check_region()
    print(f'Using AWS Region: {region}')

    queue_name = getenv('QUEUE_NAME')
    if not queue_name:
        print('Environment variable "QUEUE_NAME" must be set and is not.')
        return 1

    sts = boto3.client('sts')
    print(f'Current Role: {sts.get_caller_identity()}')

    try:
        sqs = boto3.client('sqs', region_name=region)
        response = sqs.get_queue_url(QueueName=queue_name)
        queue_url = response['QueueUrl']
    except KeyError as ex:
        print(f'Cannot get the URL for {queue_name}: {ex}')
        return 1
    except NoRegionError:
        print('Could not determine region.')
        return 1

    for num in range(20000):
        sqs.send_message(
            QueueUrl=queue_url,
            MessageBody=f'Sending message {num}'
        )


if __name__ == '__main__':
    main()
