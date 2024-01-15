package main

import (
	"cdk/stacks"
	"fmt"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

// example tests. To run these tests, uncomment this file along with the
// example resource in cdk_test.go
func TestCdkStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)
	theEnv := "testing"
	// WHEN
	authLambdas := stacks.A1Lambda(app, fmt.Sprintf("A1Lambda-%s", theEnv), &stacks.A1LambdaProps{
		Env: theEnv,
	})

	crudLambdas := stacks.U1Lambda(app, fmt.Sprintf("U1Lambda-%s", theEnv), &stacks.U1LambdaProps{
		Env: theEnv,
	})
	lambdas := append(authLambdas, crudLambdas...)
	stack := stacks.API(app, fmt.Sprintf("Undertown-API-%s", theEnv), &stacks.APIStackProps{
		IntegrationLambdas: lambdas,
		Env:                theEnv,
	})
	// THEN
	template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})

	template.HasResourceProperties(jsii.String("AWS::SQS::Queue"), map[string]interface{}{
		"VisibilityTimeout": 300,
	})
}
