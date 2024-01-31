package stacks

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type SSRStackProps struct {
	awscdk.StackProps
	Env string
}

type SSRStackReturn struct {
	LambdaUrl awslambda.FunctionUrl
	Stack     awscdk.Stack
}

// This might have to be created after the APIStack
// because it needs to import a refference to the endpoint resources so it can call the CRUD operations
func SSR(scope constructs.Construct, id string, props *SSRStackProps) SSRStackReturn {
	stack := awscdk.NewStack(scope, &id, nil)

	getPropertiesUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", GetProperties.String(), props.Env))
	getPropertyUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", GetProperty.String(), props.Env))

	homeSsrEnvVars := map[string]*string{
		"GET_PROPERTIES_URL": getPropertiesUrl,
		"GET_PROPERTY_URL":   getPropertyUrl,
		"JWT_SECRET":         jsii.String(os.Getenv("JWT_SECRET")),
	}
	// SSR
	// TODO how to import the api root path?
	// this has to be part of the api
	homeLambda := awslambdago.NewGoFunction(stack, jsii.Sprintf("ServerSideRender-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/ssr/homepage/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &homeSsrEnvVars,
	})
	// Add a Function URL.
	lambdaURL := homeLambda.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})
	awscdk.NewCfnOutput(stack, jsii.String("ServerSideRenderLambdaUrl"), &awscdk.CfnOutputProps{
		ExportName: jsii.Sprintf("ServerSideRenderLambdaUrl-%s", props.Env),
		Value:      lambdaURL.Url(),
	})

	return SSRStackReturn{
		LambdaUrl: lambdaURL,
		Stack:     stack,
	}
}
