[![Build Status](https://travis-ci.org/awslabs/k8s-cloudwatch-adapter.svg?branch=master)](https://travis-ci.org/awslabs/k8s-cloudwatch-adapter)
[![GitHub
release](https://img.shields.io/github/release/actzeroai/k8s-cloudwatch-adapter/all.svg)](https://github.com/actzeroai/k8s-cloudwatch-adapter/releases)
[![GitHub issues](https://img.shields.io/github/issues/actzeroai/k8s-cloudwatch-adapter)](https://github.com/actzeroai/k8s-cloudwatch-adapter/issues)
[![GitHub PRs](https://img.shields.io/github/issues-pr/actzeroai/k8s-cloudwatch-adapter)](https://github.com/actzeroai/k8s-cloudwatch-adapter/pulls)

# Kubernetes Custom Metrics Adapter for Kubernetes

An implementation of the Kubernetes [Custom Metrics API and External Metrics
API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis)
for AWS CloudWatch metrics.

This adapter allows you to scale your Kubernetes deployment using the [Horizontal Pod
Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) with
metrics from AWS CloudWatch.

## Prerequisites

This adapter requires your node groups to have the following permissions to access metric data from Amazon CloudWatch.
- `cloudwatch:GetMetricData`

You can create an IAM policy using this template and attach it to your node group role.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "cloudwatch:GetMetricData"
            ],
            "Resource": "*"
        }
    ]
}
```

## Deploy

Requires a Kubernetes cluster with 
[Metric Server deployed](https://docs.aws.amazon.com/eks/latest/userguide/metrics-server.html).

Now, deploy the [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) and adapter
using the Helm charts in the `/charts` directory:

```bash
$ helm install k8s-cloudwatch-adapter-crd ./charts/k8s-cloudwatch-adapter-crd \
>   --namespace custom-metrics \
>   --create-namespace
NAME: k8s-cloudwatch-adapter-crd
LAST DEPLOYED: Mon May 8 11:36:53 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
$ helm install k8s-cloudwatch-adapter ./charts/k8s-cloudwatch-adapter --namespace custom-metrics
NAME: k8s-cloudwatch-adapter
LAST DEPLOYED: Mon May 8 13:20:17 2023
NAMESPACE: custom-metrics
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

> **NOTE:** A Helm chart repository will soon be available, eliminating the need to clone the Git repository.

### Verifying the deployment

Next you can query the APIs to see if the adapter is deployed correctly by running:

```bash
$ kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq .
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "external.metrics.k8s.io/v1beta1",
  "resources": [
  ]
}
```

## Deploying the Sample Application

There is a sample SQS application provided in this repository for you to test how the adapter works. Refer to this 
[guide](samples/sqs/README.md).

## Further Reading

- [Configuring cross account metric example](docs/cross-account.md)
- [ExternalMetric CRD schema](docs/schema.md)

## License

This library is licensed under the Apache 2.0 License. [Contributions](CONTRIBUTING.md) are welcomed!

## Issues

Report any issues in the [Github Issues](https://github.com/actzeroai/k8s-cloudwatch-adapter/issues)
