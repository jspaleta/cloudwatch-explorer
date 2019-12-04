package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"fmt"
	"log"
	"os"
	"strings"
)

type Example struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

var (
	foo string
	bar string
)

func CreateAwsSessionWithOptions() *session.Session {
	// Create a Session with a custom region
	aws_session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return aws_session
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("You must supply a metric name, namespace, and dimension name")
		os.Exit(1)
	}

	metric := os.Args[1]
	namespace := os.Args[2]
	var dimensions []*cloudwatch.DimensionFilter
	if len(os.Args) > 3 {
		for _, v := range strings.Split(os.Args[3], ",") {
			name := strings.Split(v, "=")[0]
			value := strings.Split(v, "=")[1]
			filter := &cloudwatch.DimensionFilter{
				Name:  &name,
				Value: &value,
			}
			dimensions = append(dimensions, filter)
		}
	}

	session := CreateAwsSessionWithOptions()
	svc := cloudwatch.New(session)

	result, err := svc.ListMetrics(&cloudwatch.ListMetricsInput{
		MetricName: aws.String(metric),
		Namespace:  aws.String(namespace),
		Dimensions: dimensions,
	})
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	fmt.Println("Metrics", result.Metrics)

}
