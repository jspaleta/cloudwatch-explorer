package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	v2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/aws"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	//Base Sensu plugin configs
	sensu.PluginConfig
	//AWS specific Sensu plugin configs
	aws.AWSPluginConfig
	//Additional configs for this check command
	Example  string
	Verbose  bool
	MaxPages int
}

var (
	//initialize Sensu plugin Config object
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "Cloudwatch Explorer",
			Short:    "Cloudwatch Exploror",
			Keyspace: "sensu.io/plugins/cloudwatch-explorer/config",
		},
	}
	//initialize options list with custom options
	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "verbose",
			Argument:  "verbose",
			Shorthand: "v",
			Default:   false,
			Usage:     "Enable verbose output",
			Value:     &plugin.Verbose,
		},
		&sensu.PluginConfigOption{
			Path:      "max-pages",
			Argument:  "max-pages",
			Shorthand: "m",
			Default:   1,
			Usage:     "Maximum number of result pages",
			Value:     &plugin.MaxPages,
		},
	}
)

func init() {
	//append common AWS options to options list
	options = append(options, plugin.GetAWSOpts()...)
}

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *v2.Event) (int, error) {
	// Check for valid AWS credentials
	if plugin.Verbose {
		fmt.Println("Checking AWS Creds")
	}
	if state, err := plugin.CheckAWSCreds(); err != nil {
		return state, err
	}

	// Specific Argument Checking for this command
	if plugin.Verbose {
		fmt.Println("Checking Arguments")
	}

	return sensu.CheckStateOK, nil
}

func executeCheck(event *v2.Event) (int, error) {
	//Make sure plugin.CheckAwsCreds() worked as expected
	if plugin.AWSConfig == nil {
		return sensu.CheckStateCritical, fmt.Errorf("AWS Config undefined, something went wrong in processing AWS configuration information")
	}
	//Start AWS Service specific client
	client := cloudwatch.NewFromConfig(*plugin.AWSConfig)
	//Run business logic for check
	state, err := checkFunction(client)
	return state, err
}

//Create service interface to help with mock testing
type ServiceAPI interface {
	ListMetrics(ctx context.Context,
		params *cloudwatch.ListMetricsInput,
		optFns ...func(*cloudwatch.Options)) (*cloudwatch.ListMetricsOutput, error)
}

func GetMetrics(c context.Context, api ServiceAPI, input *cloudwatch.ListMetricsInput) (*cloudwatch.ListMetricsOutput, error) {
	return api.ListMetrics(c, input)
}

// Note: Use ServiceAPI interface definition to make function testable with mock API testing pattern
func checkFunction(client ServiceAPI) (int, error) {

	numPages := 0
	/*  Output format from aws documented examples
	fmt.Println("Metrics:")
	numMetrics := 0
	for _, m := range result.Metrics {
		fmt.Println("   Metric Name: " + *m.MetricName)
		fmt.Println("   Namespace:   " + *m.Namespace)
		fmt.Println("   Dimensions:")
		for _, d := range m.Dimensions {
			fmt.Println("      " + *d.Name + ": " + *d.Value)
		}
		fmt.Println("")

	}
	fmt.Println("Found " + strconv.Itoa(numMetrics) + " metrics")
	*/
	for getList := true; getList && numPages < plugin.MaxPages; {
		getList = false
		input := &cloudwatch.ListMetricsInput{}
		result, err := GetMetrics(context.TODO(), client, input)

		if err != nil {
			fmt.Println("Could not get metrics list")
			return sensu.CheckStateCritical, nil
		}
		if result.NextToken != nil {
			getList = true
			numPages++
			input.NextToken = result.NextToken
		}

		for _, m := range result.Metrics {
			namespace := *m.Namespace
			name := *m.MetricName
			var dimensions string
			for _, d := range m.Dimensions {
				k := *d.Name
				v := *d.Value
				dimensions = dimensions + fmt.Sprintf("%s=%s, ", k, v)
			}
			dimensions = strings.TrimRight(dimensions, ", ")
			fmt.Printf("%s/%s (%s)\n", namespace, name, dimensions)
		}
	}
	if numPages > plugin.MaxPages {
		fmt.Println("Warning: max allowed ListMetrics result pages exceeded, increase --max-pages value")
		return sensu.CheckStateWarning, nil
	}

	return sensu.CheckStateOK, nil
}
