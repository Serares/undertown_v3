package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	RegisterAuthorizerName = "GatewayRegisterAuthorizer"
	CRUDAuthorizer         = "GatewayCRUDAuthorizer"
)

// create a type that's a list of structs
type AuthorizerProps struct {
	awscdk.StackProps
}

// each key is a constant defined in this file
type APIAuthorizersMap map[string]APIAuthorizers

type APIAuthorizers struct {
	name    string
	handler *awsapigateway.RequestAuthorizer
}

// Those lambdas are authorizing gateway endpoints
// create both the lambdas and the apigw authorizer resources
func AuthorizersStack(scope constructs.Construct, id string, props *AuthorizerProps) APIAuthorizersMap {
	var authorizationHeader = "Authorization"
	var authorizers APIAuthorizersMap

	var sprops awscdk.StackProps
	stack := awscdk.NewStack(scope, &id, &sprops)

	//Register Authorizer
	registerAuthorizerLambda := awslambdago.NewGoFunction(stack, jsii.String("RegisterAuthorizer"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/registerAuthorizer/lambda"),
		Bundling:     BundlingOptions,
		Environment: &map[string]*string{
			"JWT_SECRET": &JwtSecret,
		},
	})
	// API Authorizer lambda
	apiAuthorizerLambda := awslambdago.NewGoFunction(stack, jsii.String("APIAuthorizer"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/authorizer/lambda"),
		Bundling:     BundlingOptions,
		Environment: &map[string]*string{
			"JWT_SECRET": &JwtSecret,
		},
	})

	// Define the SPA Authorizer
	// this authorizer only checks the validity of a jwt
	crudAuthorizer := awsapigateway.NewRequestAuthorizer(stack, jsii.String(CRUDAuthorizer), &awsapigateway.RequestAuthorizerProps{
		Handler: apiAuthorizerLambda,
		IdentitySources: &[]*string{
			awsapigateway.IdentitySource_Header(jsii.String(authorizationHeader)),
		},
	})
	// Define the Register Authorizer
	registerAuth := awsapigateway.NewRequestAuthorizer(stack, jsii.String(RegisterAuthorizerName), &awsapigateway.RequestAuthorizerProps{
		Handler: registerAuthorizerLambda,
		IdentitySources: &[]*string{
			awsapigateway.IdentitySource_Header(jsii.String(authorizationHeader)),
		},
	})
	authorizers = APIAuthorizersMap{
		RegisterAuthorizerName: APIAuthorizers{
			name:    RegisterAuthorizerName,
			handler: &registerAuth,
		},
		CRUDAuthorizer: APIAuthorizers{
			name:    CRUDAuthorizer,
			handler: &crudAuthorizer,
		},
	}

	return authorizers
}
