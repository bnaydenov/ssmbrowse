package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var sess *session.Session

const MAX_COUNT = 50

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

func GetParemters(paramPrefix []string, startToken *string, ssmParams  []ssm.ParameterMetadata)([]ssm.ParameterMetadata, *string)  {

	client := ssm.New(sess)
	
    filters := []*ssm.ParameterStringFilter{
		&ssm.ParameterStringFilter{
			Key:    aws.String(ssm.ParametersFilterKeyName),
			Option: aws.String("BeginsWith"),
			Values: aws.StringSlice(paramPrefix),
		},
	}

	input := &ssm.DescribeParametersInput{
		ParameterFilters: filters,
		MaxResults: aws.Int64(int64(MAX_COUNT)),
		NextToken: startToken,
	}

	output, err := client.DescribeParameters(input)
	if err != nil {
		fmt.Println("error describing parameters:", err)
	}

		// Add parameters returned by SSM to ssmParams
		for _, p := range output.Parameters {
			ssmParams = append(ssmParams, *p)
		}
	
	return ssmParams, output.NextToken
}