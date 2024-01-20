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

type CrudEndpoints int

const (
	AddProperty CrudEndpoints = iota
	GetProperties
	GetProperty
)

func (ce CrudEndpoints) String() string {
	return [...]string{"addProperty", "getProperties", "getProperty"}[ce]
}

type U1LambdaProps struct {
	awscdk.StackProps
	Env       string
	BasicRole awsiam.Role
}

// This stack is used for CRUD operations
func U1Lambda(scope constructs.Construct, id string, props *U1LambdaProps) []IntegrationLambda {
	var sprops awscdk.StackProps
	var lambdas []IntegrationLambda
	lambdasEnvVars := map[string]*string{
		// ❗
		// TODO try to obfuscate somehow the values
		// don't store them in plain text
		// store them as an encrypted string?
		// how to decrypt them
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
		Environment:  &lambdasEnvVars,
		Role:         lambdaRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})

	getProperties := awslambdago.NewGoFunction(stack, jsii.Sprintf("GetProperties-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/getProperties/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &lambdasEnvVars,
		Role:         lambdaRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})

	getProperty := awslambdago.NewGoFunction(stack, jsii.Sprintf("GetProperty-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/getProperty/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &lambdasEnvVars,
		Role:         lambdaRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &addProperty,
		path:       AddProperty.String(),
		method:     http.MethodPost,
		authorizer: CRUDAuthorizer.String(),
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &getProperties,
		path:       GetProperties.String(),
		method:     http.MethodGet,
		authorizer: CRUDAuthorizer.String(),
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &getProperty,
		path:       GetProperty.String(),
		method:     http.MethodGet,
		authorizer: CRUDAuthorizer.String(),
	})

	return lambdas
}
