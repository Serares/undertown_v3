package stacks

import (
	"cdk/utils"
	"net/http"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type U1LambdaProps struct {
	awscdk.StackProps
	Env       string
	BasicRole awsiam.Role
}

// This stack is used for CRUD operations
func U1Lambda(scope constructs.Construct, id string, props *U1LambdaProps) []IntegrationLambda {
	var sprops awscdk.StackProps
	var lambdas []IntegrationLambda
	dbCfg := map[string]*string{
		"DB_HOST":        jsii.String(os.Getenv("DB_HOST")),
		"DB_NAME":        jsii.String(os.Getenv("DB_NAME")),
		"DB_PROTOCOL":    jsii.String(os.Getenv("DB_PROTOCOL")),
		"TURSO_DB_TOKEN": jsii.String(os.Getenv("TURSO_DB_TOKEN")),
	}
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	lambdaRole := utils.CreateLambdaBasicRole(stack, props.Env)

	addProperty := awslambdago.NewGoFunction(stack, jsii.Sprintf("AddProperty-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/addProperty/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &dbCfg,
		Role:         lambdaRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60 * 2)),
		// TODO read this
		// if lambdas are part of a VPC they will need EIP's assigned to the subnets that lambdas are part of
		// https://stackoverflow.com/questions/52992085/why-cant-an-aws-lambda-function-inside-a-public-subnet-in-a-vpc-connect-to-the/52994841#52994841

	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &addProperty,
		path:       "addProperty",
		method:     http.MethodPost,
		authorizer: CRUDAuthorizer,
	})

	return lambdas
}
