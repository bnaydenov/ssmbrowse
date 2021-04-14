package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/bnaydenov/ssmbrowse/internal/pkg/aws"
)


func main() {
    
	var startToken *string
	var params []ssm.ParameterMetadata
	
	params, nextToken := aws.GetParemters([]string{"/"}, startToken, params)
    
	for nextToken != nil {
		params, nextToken = aws.GetParemters([]string{"/"}, nextToken, params)
		fmt.Println("next page.....")
	}
	
	for _, p := range params {
		fmt.Println(*p.Name)
	}
	fmt.Println(len(params))
	// cmd.Entrypoint()
}