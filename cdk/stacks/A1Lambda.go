package stacks

import (
	"cdk/utils"
	"net/http"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AuthEndpoints int

const (
	LoginEndpoint AuthEndpoints = iota
	RegisterEndpoint
)

func (ae AuthEndpoints) String() string {
	return [...]string{"login", "register"}[ae]
}

type A1LambdaProps struct {
	awscdk.StackProps
	// Vpc         awsec2.Vpc Deprecated
	Env string
}

// This stack is used for authentication/registration
func A1Lambda(scope constructs.Construct, id string, props *A1LambdaProps) []IntegrationLambda {
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
	lambdaRole := utils.CreateLambdaBasicRole(stack, "lambdaBasicRoleA1", props.Env)

	registerLambda := awslambdago.NewGoFunction(stack, jsii.Sprintf("Register-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/register/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &dbCfg,
		Role:         lambdaRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})

	// Login Lambda
	loginLambda := awslambdago.NewGoFunction(stack, jsii.Sprintf("Login-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/login/lambda"),
		Bundling:     BundlingOptions,
		Environment: &map[string]*string{
			"DB_HOST":        jsii.String(os.Getenv("DB_HOST")),
			"DB_NAME":        jsii.String(os.Getenv("DB_NAME")),
			"DB_PROTOCOL":    jsii.String(os.Getenv("DB_PROTOCOL")),
			"TURSO_DB_TOKEN": jsii.String(os.Getenv("TURSO_DB_TOKEN")),
			"JWT_SECRET":     jsii.String(os.Getenv("JWT_SECRET")),
		},
		// todo have to add the JWT_SECRET IN HERE
		Role:    lambdaRole,
		Timeout: awscdk.Duration_Seconds(jsii.Number(30)),
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &registerLambda,
		path:       RegisterEndpoint.String(),
		method:     http.MethodPost,
		authorizer: RegisterAuthorizer.String(),
	})
	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &loginLambda,
		path:       LoginEndpoint.String(),
		method:     http.MethodPost,
		authorizer: "",
	})

	return lambdas
}
