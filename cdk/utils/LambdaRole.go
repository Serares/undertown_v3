package utils

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateLambdaBasicRole(stack constructs.Construct, env string) awsiam.Role {
	basicLambdaRole := awsiam.NewRole(stack, jsii.Sprintf("UndertownLambda-%s", env), &awsiam.RoleProps{
		// TODO think about if this needs to apply the least priviledge principle
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
	})
	basicLambdaRole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaBasicExecutionRole")))
	return basicLambdaRole
}
