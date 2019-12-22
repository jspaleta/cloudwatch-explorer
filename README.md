# AWS CloudWatch Metrics Explorer

## Overview

The `cloudwatch-explorer` utility can be used to lookup available Cloudwatch
metrics.

```
An AWS Cloudwatch metric explorer.

Usage:
  cloudwatch-explorer [flags]

Flags:
  -h, --help                       help for cloudwatch-explorer
  -d, --metric-dimensions string   The AWS Cloudwatch metric dimension. Can also
                                   be set via the $CLOUDWATCH_METRIC_DIMENSION
                                   environment variable. OPTIONAL.
  -m, --metric-name string         The AWS Cloudwatch metric name. Can also be
                                   set via the $CLOUDWATCH_METRIC_NAME
                                   environment variable. OPTIONAL.
  -n, --metric-namespace string    The AWS Cloudwatch metric namespace. Can also
                                   be set via the $CLOUDWATCH_METRIC_NAMESPACE
                                   environment variable. (default "AWS/EC2")
```

Example output:

```shell
$ ./cloudwatch-explorer -n AWS/EC2 -d InstanceId=i-xxxxxxxxxxxxxxxxx
AWS/EC2/MetadataNoToken (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/EBSIOBalance% (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/EBSByteBalance% (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/EBSReadOps (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/EBSReadBytes (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/EBSWriteOps (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/EBSWriteBytes (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/NetworkIn (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/NetworkOut (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/NetworkPacketsIn (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/NetworkPacketsOut (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/CPUUtilization (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/StatusCheckFailed_System (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/StatusCheckFailed_Instance (InstanceId=i-xxxxxxxxxxxxxxxxx)
AWS/EC2/StatusCheckFailed (InstanceId=i-xxxxxxxxxxxxxxxxx)

./cloudwatch-explorer -n AWS/Kinesis
AWS/Kinesis/WriteProvisionedThroughputExceeded (StreamName=example-stream)
AWS/Kinesis/PutRecords.Success (StreamName=example-stream)
AWS/Kinesis/PutRecords.Bytes (StreamName=example-stream)
AWS/Kinesis/IncomingBytes (StreamName=example-stream)
AWS/Kinesis/IncomingRecords (StreamName=example-stream)
AWS/Kinesis/PutRecords.Latency (StreamName=example-stream)
AWS/Kinesis/PutRecords.Records (StreamName=example-stream)
```

## Roadmap

- [x] Authenticate to the AWS API (via `AWS_ACCESS_KEY` and `AWS_SECRET_ACCESS_KEY`)
- [x] Support AWS IAM environment variables
  - [x] `$AWS_ACCESS_KEY`
  - [x] `AWS_SECRET_ACCESS_KEY`
  - [x] `AWS_REGION`
- [x] Use the [CloudWatch ListMetrics API][1] to list available metrics per
      region, CloudWatch Namespace, and CloudWatch Metric
- [x] Optionally filter available metrics by one or more CloudWatch Dimensions
- [x] Add command flags and `--help` usage instructions
- [ ] Fetch CloudWatch metrics
- [ ] Output collected metrics in Graphite plaintext format
- [ ] Output collected metrics in InfluxDB line protocol format
- [ ] Output collected metrics in Prometheus exposition format

[1]: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_ListMetrics.html

## References

- [Cloudwatch Golang SDK examples](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/cw-example-getting-metrics.html)
- [Cloudwatch Namespaces Reference](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/aws-services-cloudwatch-metrics.html)
- [AWS/EC2 Cloudwatch Metrics](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/viewing_metrics_with_cloudwatch.html#ec2-cloudwatch-dimensions)
