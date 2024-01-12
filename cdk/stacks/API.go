package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type IntegrationLambda struct {
	goLambda   *awslambdago.GoFunction
	path       string
	method     string
	authorizer *awsapigateway.RequestAuthorizer
}

type APIStackProps struct {
	awscdk.StackProps
	integrationLambdas []IntegrationLambda
}

// The API Gateway resources and deployments
func API(scope constructs.Construct, id string, props *APIStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Define the API Gateway
	spaApi := awsapigateway.NewRestApi(stack, jsii.String("SPAUndertownAPI"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("SPAUndertownAPI"),
		DeployOptions: &awsapigateway.StageOptions{
			TracingEnabled: jsii.Bool(true),
		},
		// CloudWatchRole: jsii.Bool(true),
	})

	for _, lambda := range props.integrationLambdas {
		integration := awsapigateway.NewLambdaIntegration(*lambda.goLambda, &awsapigateway.LambdaIntegrationOptions{})
		spaApi.Root().AddResource(jsii.String(lambda.path), &awsapigateway.ResourceOptions{}).AddMethod(jsii.String(lambda.method), integration, &awsapigateway.MethodOptions{
			Authorizer: *lambda.authorizer,
		})
	}

	deployment := awsapigateway.NewDeployment(stack, jsii.String("SPAAPIDeployment"), &awsapigateway.DeploymentProps{
		Api: spaApi,
	})

	// devLogGroup := awslogs.NewLogGroup(stack, jsii.String("devlogs"), &awslogs.LogGroupProps{})

	stage := awsapigateway.NewStage(stack, jsii.String("DevStage"), &awsapigateway.StageProps{
		Deployment: deployment,
		StageName:  jsii.String("dev"),
		// AccessLogDestination: awsapigateway.NewLogGroupLogDestination(devLogGroup),
		// AccessLogFormat: awsapigateway.AccessLogFormat_JsonWithStandardFields(&awsapigateway.JsonWithStandardFieldProps{
		// 	Caller:         jsii.Bool(false),
		// 	HttpMethod:     jsii.Bool(true),
		// 	Ip:             jsii.Bool(true),
		// 	Protocol:       jsii.Bool(true),
		// 	RequestTime:    jsii.Bool(true),
		// 	ResourcePath:   jsii.Bool(true),
		// 	ResponseLength: jsii.Bool(true),
		// 	Status:         jsii.Bool(true),
		// 	User:           jsii.Bool(true),
		// }),
	})
	spaApi.SetDeploymentStage(stage)
	return stack
}
