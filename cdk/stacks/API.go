package stacks

import (
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AuthorizerType int

const (
	CRUDAuthorizer AuthorizerType = iota
	RegisterAuthorizer
)

func (at AuthorizerType) String() string {
	return [...]string{"CrudAuthorizer", "RegisterAuthorizer"}[at]
}

type IntegrationLambda struct {
	goLambda   *awslambdago.GoFunction
	path       string
	method     []string
	authorizer string
}

type APIStackProps struct {
	awscdk.StackProps
	IntegrationLambdas []IntegrationLambda
	Env                string
}

type ExportedApiPathResources map[string]string

// The API Gateway resources and deployments
func API(scope constructs.Construct, id string, props *APIStackProps) awscdk.Stack {
	var authorizationHeader = "Authorization"
	var JwtSecret = os.Getenv("JWT_SECRET")
	var registratorSecret = os.Getenv("REGISTRATOR_SECRET")
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Define the API Gateway
	spaApi := awsapigateway.NewRestApi(stack, jsii.Sprintf("UndertownAPI-%s", props.Env), &awsapigateway.RestApiProps{
		RestApiName:    jsii.Sprintf("UndertownAPI-%s", props.Env),
		Deploy:         jsii.Bool(false),
		CloudWatchRole: jsii.Bool(true),
	})
	//Register Authorizer
	registerAuthorizerLambda := awslambdago.NewGoFunction(stack, jsii.Sprintf("RegisterAuthorizer-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/registerAuthorizer/lambda"),
		Bundling:     BundlingOptions,
		Environment: &map[string]*string{
			"JWT_SECRET":         &JwtSecret,
			"REGISTRATOR_SECRET": &registratorSecret,
		},
	})
	// API Authorizer lambda
	apiAuthorizerLambda := awslambdago.NewGoFunction(stack, jsii.Sprintf("APIAuthorizer-%s", props.Env), &awslambdago.GoFunctionProps{
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
	crudAuthorizer := awsapigateway.NewRequestAuthorizer(stack, jsii.String(CRUDAuthorizer.String()), &awsapigateway.RequestAuthorizerProps{
		Handler: apiAuthorizerLambda,
		IdentitySources: &[]*string{
			awsapigateway.IdentitySource_Header(jsii.String(authorizationHeader)),
		},
		ResultsCacheTtl: awscdk.Duration_Seconds(jsii.Number(0)),
	})
	// Define the Register Authorizer
	registerAuth := awsapigateway.NewRequestAuthorizer(stack, jsii.String(RegisterAuthorizer.String()), &awsapigateway.RequestAuthorizerProps{
		Handler: registerAuthorizerLambda,
		IdentitySources: &[]*string{
			awsapigateway.IdentitySource_Header(jsii.String(authorizationHeader)),
		},
		ResultsCacheTtl: awscdk.Duration_Seconds(jsii.Number(0)),
	})

	for _, lambda := range props.IntegrationLambdas {
		var methodOptions *awsapigateway.MethodOptions
		switch lambda.authorizer {
		case CRUDAuthorizer.String():
			methodOptions = &awsapigateway.MethodOptions{
				Authorizer: crudAuthorizer,
			}
		case RegisterAuthorizer.String():
			methodOptions = &awsapigateway.MethodOptions{
				Authorizer: registerAuth,
			}
		default:
			// TODO try to initialize with an empty struct and do nothing on the default case
			methodOptions = &awsapigateway.MethodOptions{}
		}

		integration := awsapigateway.NewLambdaIntegration(*lambda.goLambda, &awsapigateway.LambdaIntegrationOptions{})

		resource := spaApi.Root().AddResource(jsii.String(lambda.path), &awsapigateway.ResourceOptions{})
		for _, httpVerb := range lambda.method {
			resource.AddMethod(jsii.String(httpVerb), integration, methodOptions)
		}
	}

	deployment := awsapigateway.NewDeployment(stack, jsii.Sprintf("UndertownAPIDeployment-%s", props.Env), &awsapigateway.DeploymentProps{
		Api: spaApi,
	})

	devLogGroup := awslogs.NewLogGroup(stack, jsii.String("api-logs"), &awslogs.LogGroupProps{})

	stage := awsapigateway.NewStage(stack, jsii.Sprintf("APISTAGE-%s", props.Env), &awsapigateway.StageProps{
		Deployment:           deployment,
		StageName:            jsii.String(props.Env),
		AccessLogDestination: awsapigateway.NewLogGroupLogDestination(devLogGroup),
		AccessLogFormat: awsapigateway.AccessLogFormat_JsonWithStandardFields(&awsapigateway.JsonWithStandardFieldProps{
			Caller:         jsii.Bool(false),
			HttpMethod:     jsii.Bool(true),
			Ip:             jsii.Bool(true),
			Protocol:       jsii.Bool(true),
			RequestTime:    jsii.Bool(true),
			ResourcePath:   jsii.Bool(true),
			ResponseLength: jsii.Bool(true),
			Status:         jsii.Bool(true),
			User:           jsii.Bool(true),
		}),
	})
	spaApi.SetDeploymentStage(stage)
	var exportedPaths = make(ExportedApiPathResources)
	for _, lambda := range props.IntegrationLambdas {
		awscdk.NewCfnOutput(stack, jsii.Sprintf("UndertownAPIPath-%s", lambda.path), &awscdk.CfnOutputProps{
			Value:      jsii.Sprintf("%s", *spaApi.UrlForPath(jsii.String("/" + lambda.path))),
			ExportName: jsii.Sprintf("%s-%s", lambda.path, props.Env),
		})
		exportedPaths[lambda.path] = fmt.Sprintf("%s/%s/%s", *spaApi.Url(), *stage.StageName(), lambda.path)
	}
	return stack
}
