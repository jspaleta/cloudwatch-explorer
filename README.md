# AWS CloudWatch Metrics Collection Prototype

Something something collect metrics from AWS CloudWatch something.

## Overview

Coming soon.

## Roadmap

- [x] Authenticate to the AWS API (via `AWS_ACCESS_KEY` and `AWS_SECRET_ACCESS_KEY`)
- [x] Support AWS IAM environment variables
  - [x] `$AWS_ACCESS_KEY`
  - [x] `AWS_SECRET_ACCESS_KEY`
  - [x] `AWS_REGION`
- [x] Use the [CloudWatch ListMetrics API][1] to list available metrics per
      region, CloudWatch Namespace, and CloudWatch Metric
- [x] Optionally filter available metrics by one or more CloudWatch Dimensions
- [ ] Fetch CloudWatch metrics
- [ ] Output collected metrics in Graphite plaintext format
- [ ] Output collected metrics in InfluxDB line protocol format
- [ ] Output collected metrics in Prometheus exposition format
- [ ] Add command flags and `--help` usage instructions 

[1]: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_ListMetrics.html

## References

- [Cloudwatch Golang SDK examples](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/cw-example-getting-metrics.html)
- [Cloudwatch Namespaces Reference](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/aws-services-cloudwatch-metrics.html)
- [AWS/EC2 Cloudwatch Metrics](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/viewing_metrics_with_cloudwatch.html#ec2-cloudwatch-dimensions)
