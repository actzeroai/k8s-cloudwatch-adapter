package aws

import (
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/awslabs/k8s-cloudwatch-adapter/pkg/apis/metrics/v1alpha1"

	"k8s.io/klog"
)

func GetIMDSv2Token() (string, error) {
    tokenURL := "http://169.254.169.254/latest/api/token"
    req, err := http.NewRequest("PUT", tokenURL, nil)
    if err != nil {
        return "", err
    }

    req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "60")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", nil
    }

    defer resp.Body.Close()
    token, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(token), nil
}

// GetLocalRegion gets the region ID from the instance metadata.
func GetLocalRegion() string {
    metadataURL := "http://169.254.169.254/latest/meta-data/placement/availability-zone/"
	resp, err := http.Get(metadataURL)
	if err != nil {
		klog.Errorf("unable to get current region information, %v", err)
		return ""
	}
	if resp.StatusCode == 401 {
	    klog.Info("Token required to obtain instance metadata. Trying that method.")
	    token, err := GetIMDSv2Token()
	    if err != nil {
	        klog.Fatalf("Unable to get metadata v2 token, %v", err)
	    }

	    req, err := http.NewRequest("GET", metadataURL, nil)
	    if err != nil {
	        klog.Errorf("Unable to get metadata v2 token, %v", err)
	    }
        req.Header.Set("X-aws-ec2-metadata-token", token)
        client := &http.Client{}
        resp, err = client.Do(req)
        if err != nil {
	        klog.Errorf("Unable to get metadata v2 token, %v", err)
	    }
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("cannot read response from instance metadata, %v", err)
	}

	// strip the last character from AZ to get region ID
	return string(body[0 : len(body)-1])
}

func toCloudWatchQuery(externalMetric *v1alpha1.ExternalMetric) cloudwatch.GetMetricDataInput {
	queries := externalMetric.Spec.Queries

	cwMetricQueries := make([]*cloudwatch.MetricDataQuery, len(queries))
	for i, q := range queries {
		q := q
		returnData := &q.ReturnData
		mdq := &cloudwatch.MetricDataQuery{
			Id:         &q.ID,
			Label:      &q.Label,
			ReturnData: *returnData,
		}

		if len(q.Expression) == 0 {
			dimensions := make([]*cloudwatch.Dimension, len(q.MetricStat.Metric.Dimensions))
			for j := range q.MetricStat.Metric.Dimensions {
				dimensions[j] = &cloudwatch.Dimension{
					Name:  &q.MetricStat.Metric.Dimensions[j].Name,
					Value: &q.MetricStat.Metric.Dimensions[j].Value,
				}
			}

			metric := &cloudwatch.Metric{
				Dimensions: dimensions,
				MetricName: &q.MetricStat.Metric.MetricName,
				Namespace:  &q.MetricStat.Metric.Namespace,
			}

			mdq.MetricStat = &cloudwatch.MetricStat{
				Metric: metric,
				Period: &q.MetricStat.Period,
				Stat:   &q.MetricStat.Stat,
				Unit:   aws.String(q.MetricStat.Unit),
			}
		} else {
			mdq.Expression = &q.Expression
		}

		cwMetricQueries[i] = mdq
	}
	cwQuery := cloudwatch.GetMetricDataInput{
		MetricDataQueries: cwMetricQueries,
	}

	return cwQuery
}
