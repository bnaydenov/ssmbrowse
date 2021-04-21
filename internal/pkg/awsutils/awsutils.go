package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var sess *session.Session

const MAX_COUNT = 10

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

//DescribeParameters is returning ssm params staring with specific prefix 
func DescribeParameters(paramPrefix *string, startToken *string, ssmParams []ssm.ParameterMetadata) ([]ssm.ParameterMetadata, *string) {

	client := ssm.New(sess)
	
	filters := []*ssm.ParameterStringFilter{
		&ssm.ParameterStringFilter{
			Key:    aws.String(ssm.ParametersFilterKeyName),
			Option: aws.String("Contains"),
			Values: []*string{paramPrefix},
		},
	}

	input := &ssm.DescribeParametersInput{
		ParameterFilters: filters,
		MaxResults:       aws.Int64(int64(MAX_COUNT)),
		NextToken:        startToken,
	}
    // input := &ssm.GetParametersByPathInput{
    // 	Path: paramPrefix,
	// 	MaxResults:       aws.Int64(int64(MAX_COUNT)),
	// 	NextToken:        startToken,
	// 	Recursive: aws.Bool(true) ,
	// 	WithDecryption: aws.Bool(false),
	// }
	// output, err  := client.GetParametersByPath(input)

	output, err := client.DescribeParameters(input)
	if err != nil {
		fmt.Println("error getting parameters:", err)
	}

	// // Add parameters returned by SSM to ssmParams
	for _, p := range output.Parameters {
		ssmParams = append(ssmParams, *p)
	}

	return ssmParams, output.NextToken
}

//GetParameter is 
func GetParameter(paramName string) *ssm.GetParameterOutput {
    client := ssm.New(sess)

	input := &ssm.GetParameterInput{
    	Name: aws.String(paramName),
		WithDecryption: aws.Bool(true),
	}
	output, err  := client.GetParameter(input)
	if err != nil {
		fmt.Println("error getting parameter:", err)
	}

	return output
}