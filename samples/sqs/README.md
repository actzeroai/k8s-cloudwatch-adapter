# Sample Amazon SQS Application

This directory contains a sample Amazon SQS application to test out `k8s-cloudwatch-adapter`. SQS
producer and consumer are provided, together with the YAML files for deploying the consumer, metric
configuration and HPA.

## Prerequisites

Before starting, you need to first create an IAM role which can be assumed by a Kubernetes `ServiceAccount` (named 
`sqstest` in the manifest) and an SQS queue. Attach the following policy to the role you created:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "sqs:*",
            "Resource": "<YOUR_SQS_QUEUE_ARN>"
        }
    ]
}
```

## Deploying the Amazon SQS Consumer

Modify [serviceaccount.yaml](deploy/serviceaccount.yaml) to include the ARN of the role you created and then deploy it.

```bash
$ kubectl apply -f deploy/serviceaccount.yaml
```

Next, add the name of your SQS queue and the name of the region in which it exists to the `env:` section of
[consumer-deployment.yaml](deploy/consumer-deployment.yaml) as `QUEUE_NAME` and `AWS_REGION`. Then, we can deploy it:

```bash
$ kubectl apply -f deploy/consumer-deployment.yaml
```

You can verify the consumer is running by executing this command:

```bash
$ kubectl get deploy
NAME           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
sqs-consumer   1         1         1            0           5s
```

## Setup Amazon CloudWatch Metric and HPA

Next, we will need to create an `ExternalMetric` resource for Amazon CloudWatch metric. This resource will tell the 
adapter how to retrieve metric data from Amazon CloudWatch. Modify [externalmetric.yaml](deploy/externalmetric.yaml) to
include the SQS queue name (see line 15) and deploy it.

```bash
$ kubectl apply -f deploy/extermalmetric.yaml
```

For details about how metric data query works, please refer to 
[CloudWatch GetMetricData API](https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_GetMetricData.html).

Now we can deploy the HorizontalPodAutoscaler for our consumer:

```bash
$ kubectl apply -f deploy/hpa.yaml
```

## Generate Load Using the Producer

Finally, we can start generating messages to the queue. A simple [python script](producer/producer.py) will, using your
current AWS CLI credentials, publish 10,000 messages to the queue.

```bash
$ python3 -m venv .venv
$ source .venv/bin/activate
(.venv) $ pip install boto3
(.venv) $ QUEUE_NAME=<YOUR_SQS_QUEUE_NAME> python producer/producer.py
(.venv) $ deactivate
```

> **NOTE**: The producer can also be built and run in a Docker container, if Python 3.8 or later is not available on your system.

On a separate terminal, you can now watch your HPA retrieving the queue length and start scaling the replicas. Amazon 
SQS now supports metrics at 1-minute interval, so you should start to see the deployment scale pretty quickly.

```bash
$ kubectl get hpa sqs-consumer-scaler --watch
```

## Clean Up

Once you are done with this experiment, you can delete the Kubernetes deployment and respective
resources.

Press `ctrl+c` to terminate the producer if it is still running.

Execute the following commands to remove the consumer, external metric, HPA and Amazon SQS queue.
```bash
$ kubectl delete -f deploy/hpa.yaml
$ kubectl delete -f deploy/externalmetric.yaml
$ kubectl delete -f deploy/consumer-deployment.yaml
$ kubectl delete -f deploy/serviceaccount.yaml

$ aws sqs delete-queue <YOUR_SQS_QUEUE_NAME>
```

Don't forget to clean up your IAM role and policy, too.
