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

//GetParemtersByPrefix is returning ssm params staring with specific prefix 
func GetParemtersByPrefix(paramPrefix *string, startToken *string, ssmParams []ssm.Parameter) ([]ssm.Parameter, *string) {

	client := ssm.New(sess)
	
	// filters := []*ssm.ParameterStringFilter{
	// 	&ssm.ParameterStringFilter{
	// 		Key:    aws.String(ssm.ParametersFilterKeyName),
	// 		Option: aws.String("BeginsWith"),
	// 		Values: aws.StringSlice(paramPrefix),
	// 	},
	// }

	// input := &ssm.DescribeParametersInput{
	// 	ParameterFilters: filters,
	// 	MaxResults:       aws.Int64(int64(MAX_COUNT)),
	// 	NextToken:        startToken,
	// }
    input := &ssm.GetParametersByPathInput{
    	Path: paramPrefix,
		MaxResults:       aws.Int64(int64(MAX_COUNT)),
		NextToken:        startToken,
		Recursive: aws.Bool(true) ,
		WithDecryption: aws.Bool(false),
	}
	output, err  := client.GetParametersByPath(input)

	// output, err := client.DescribeParameters(input)
	if err != nil {
		fmt.Println("error getting parameters:", err)
	}

	// // Add parameters returned by SSM to ssmParams
	for _, p := range output.Parameters {
		ssmParams = append(ssmParams, *p)
	}

	return ssmParams, output.NextToken
}
