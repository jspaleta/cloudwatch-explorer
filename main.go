package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"fmt"
	"log"
	"strings"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugins-go-library/sensu"
)

type CheckConfig struct {
	sensu.PluginConfig
	cloudwatchMetricName             string
	cloudwatchMetricNamespace        string
	cloudwatchMetricDimensions       string
	cloudwatchMetricDimensionFilters []*cloudwatch.DimensionFilter
}

var (
	config = CheckConfig{
		PluginConfig: sensu.PluginConfig{
			Name:  "cloudwatch-metric-lister",
			Short: "The Sensu Go Cloudwatch check plugin for listing Cloudwatch Metrics.",
		},
	}

	cloudwatchConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      "metric-name",
			Env:       "CLOUDWATCH_METRIC_NAME",
			Argument:  "metric-name",
			Shorthand: "m",
			Usage:     "The AWS Cloudwatch metric name. Can also be set via the $CLOUDWATCH_METRIC_NAME environment variable.",
			Value:     &config.cloudwatchMetricName,
			Default:   "",
		},
		{
			Path:      "metric-namespace",
			Env:       "CLOUDWATCH_METRIC_NAMESPACE",
			Argument:  "metric-namespace",
			Shorthand: "n",
			Usage:     "The AWS Cloudwatch metric namespace. Can also be set via the $CLOUDWATCH_METRIC_NAMESPACE environment variable.",
			Value:     &config.cloudwatchMetricNamespace,
			Default:   "",
		},
		{
			Path:      "metric-dimensions",
			Env:       "CLOUDWATCH_METRIC_DIMENSION",
			Argument:  "metric-dimensions",
			Shorthand: "d",
			Usage:     "The AWS Cloudwatch metric dimension. Can also be set via the $CLOUDWATCH_METRIC_DIMENSION environment variable.",
			Value:     &config.cloudwatchMetricDimensions,
			Default:   "",
		},
	}
)

func main() {
	check := sensu.InitCheck(&config.PluginConfig, cloudwatchConfigOptions, validateArgs, collectMetrics)
	check.Execute()
}

func validateArgs(event *corev2.Event) error {
	if config.cloudwatchMetricNamespace == "" {
		log.Fatalf("ERROR: no Cloudwatch metric namespace provided.")
		return fmt.Errorf("No Cloudwatch metric namespace provided.")
	}

	if config.cloudwatchMetricName == "" {
		log.Fatalf("ERROR: no Cloudwatch metric name provided.")
		return fmt.Errorf("No Cloudwatch metric name provided.")
	}

	if config.cloudwatchMetricDimensions == "" {
		log.Fatalf("ERROR: no Cloudwatch metric dimension(s) provided.")
		return fmt.Errorf("No Cloudwatch metric dimension(s) provided.")
	}

	return nil
}

func CreateAwsSessionWithOptions() *session.Session {
	// Create a Session with a custom region
	aws_session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return aws_session
}

// Parse Dimension strings, return slice of cloudwatch.DimensionFilter structs.
func parseCloudwatchMetricDimensions(d string) error {
	for _, v := range strings.Split(d, ",") {
		name := strings.Split(v, "=")[0]
		value := strings.Split(v, "=")[1]
		filter := &cloudwatch.DimensionFilter{
			Name:  &name,
			Value: &value,
		}
		config.cloudwatchMetricDimensionFilters = append(config.cloudwatchMetricDimensionFilters, filter)
	}
	return nil
}

func collectMetrics(event *corev2.Event) error {
	session := CreateAwsSessionWithOptions()
	svc := cloudwatch.New(session)

	err := parseCloudwatchMetricDimensions(config.cloudwatchMetricDimensions)
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}

	result, err := svc.ListMetrics(&cloudwatch.ListMetricsInput{
		MetricName: aws.String(config.cloudwatchMetricName),
		Namespace:  aws.String(config.cloudwatchMetricNamespace),
		Dimensions: config.cloudwatchMetricDimensionFilters,
	})
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	fmt.Println("Metrics", result.Metrics)
	return nil
}
