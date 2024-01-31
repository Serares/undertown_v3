package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AdminSSRStackProps struct {
	awscdk.StackProps
	Env string
}

type AdminSSRStackReturn struct {
	LambdaUrl awslambda.FunctionUrl
	Stack     awscdk.Stack
}

// This might have to be created after the APIStack
// because it needs to import a refference to the endpoint resources so it can call the CRUD operations
func AdminSSR(scope constructs.Construct, id string, props *AdminSSRStackProps) AdminSSRStackReturn {
	stack := awscdk.NewStack(scope, &id, nil)

	getPropertiesUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", GetProperties.String(), props.Env))
	getPropertyUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", GetProperty.String(), props.Env))
	loginUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", LoginEndpoint.String(), props.Env))
	addPropertyUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", AddProperty.String(), props.Env))

	envVars := map[string]*string{
		"GET_PROPERTIES_URL":  getPropertiesUrl,
		"GET_PROPERTY_URL":    getPropertyUrl,
		"LOGIN_URL":           loginUrl,
		"SUBMIT_PROPERTY_URL": addPropertyUrl,
		// "JWT_SECRET":          jsii.String(os.Getenv("JWT_SECRET")), // THIS is not needed because the token is in the cookie after login
	}
	// SSR
	// TODO how to import the api root path?
	// this has to be part of the api
	homeLambda := awslambdago.NewGoFunction(stack, jsii.Sprintf("UndertownAdmin-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/ssr/admin/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &envVars,
	})
	// Add a Function URL.
	lambdaURL := homeLambda.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})
	awscdk.NewCfnOutput(stack, jsii.String("UndertownAdminLambdaUrl"), &awscdk.CfnOutputProps{
		ExportName: jsii.Sprintf("UndertownAdminLambdaUrl-%s", props.Env),
		Value:      lambdaURL.Url(),
	})

	return AdminSSRStackReturn{
		LambdaUrl: lambdaURL,
		Stack:     stack,
	}
}
