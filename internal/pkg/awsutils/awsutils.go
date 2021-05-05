package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

var sess *session.Session

const MAX_COUNT = 30

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

//SsmDescribeParameters is returning ssm params staring with specific prefix
func SsmDescribeParameters(paramPrefix *string, nextToken *string, ssmParams []ssm.ParameterMetadata) ([]ssm.ParameterMetadata, *string, error) {

	client := ssm.New(sess)

	filters := []*ssm.ParameterStringFilter{
		{
			Key:    aws.String(ssm.ParametersFilterKeyName),
			Option: aws.String("Contains"),
			Values: []*string{paramPrefix},
		},
	}

	input := &ssm.DescribeParametersInput{
		ParameterFilters: filters,
		MaxResults:       aws.Int64(int64(MAX_COUNT)),
		NextToken:        nextToken,
	}

	output, err := client.DescribeParameters(input)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Error when trying to get ssm params containing '%s'", *paramPrefix))
		return nil, nil, err
	}

	// // Add parameters returned by SSM to ssmParams
	for _, p := range output.Parameters {
		ssmParams = append(ssmParams, *p)
	}

	return ssmParams, output.NextToken, err
}

//GetParameter is
func GetParameter(paramName string) (*ssm.GetParameterOutput, error) {
	client := ssm.New(sess)

	input := &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	}
	
	output, err := client.GetParameter(input)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Error when trying to get ssm param: %s", paramName))
		return nil, err
	}

	return output, err
}

//GetAwsSessionDetails is
func GetAwsSessionDetails() (accountID *string, region *string, err error) {

	client := sts.New(sess)
	output, err := client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		err = errors.Wrap(err, "Error when trying to current aws get caller identity")
		return nil, nil, err
	}

	return output.Account, client.Config.Region, nil
}
