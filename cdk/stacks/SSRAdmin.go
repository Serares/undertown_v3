package stacks

import (
	"cdk/utils"
	"os"

	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type SSRAdminStackProps struct {
	PUQueue         awssqs.Queue
	RawImagesBucket awss3.Bucket
	awscdk.StackProps
	Env string
}

type SSRAdminStackReturn struct {
	LambdaUrl awslambda.FunctionUrl
	Stack     awscdk.Stack
}

// This might have to be created after the APIStack
// because it needs to import a refference to the endpoint resources so it can call the CRUD operations
func SSRAdmin(scope constructs.Construct, id string, props *SSRAdminStackProps) SSRAdminStackReturn {
	stack := awscdk.NewStack(scope, &id, nil)

	getPropertiesUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", GetProperties.String(), props.Env))
	getPropertyUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", GetProperty.String(), props.Env))
	loginUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", LoginEndpoint.String(), props.Env))
	deletePropertyUrl := awscdk.Fn_ImportValue(jsii.Sprintf("%s-%s", DeleteProperty.String(), props.Env))

	var PUQueueUrl string
	var rawImagesBucketName string
	if props.RawImagesBucket.BucketName() != nil {
		rawImagesBucketName = *props.RawImagesBucket.BucketName()
	}

	if props.PUQueue.QueueUrl() != nil {
		PUQueueUrl = *props.PUQueue.QueueUrl()
	}

	s3BucketAccessRole := utils.CreateLambdaBasicRole(stack, "s3fullaccesslambdarole", props.Env)
	s3BucketAccessRole.AddManagedPolicy(
		awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonS3FullAccess")))

	envVars := map[string]*string{
		env.GET_PROPERTIES_URL:  getPropertiesUrl,
		env.GET_PROPERTY_URL:    getPropertyUrl,
		env.LOGIN_URL:           loginUrl,
		env.DELETE_PROPERTY_URL: deletePropertyUrl,
		env.SQS_PU_QUEUE_URL:    jsii.String(PUQueueUrl),
		env.RAW_IMAGES_BUCKET:   jsii.String(rawImagesBucketName),
		env.JWT_SECRET:          jsii.String(os.Getenv(env.JWT_SECRET)), // this is needed to decode the cookie and add the user id into the sqs message
	}
	// SSR
	// TODO how to import the api root path?
	// this has to be part of the api
	adminSSRLambda := awslambdago.NewGoFunction(stack, jsii.Sprintf("UndertownAdmin-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/ssr/admin/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &envVars,
		Role:         s3BucketAccessRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})
	// Add a Function URL.
	lambdaURL := adminSSRLambda.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})
	awscdk.NewCfnOutput(stack, jsii.String("UndertownAdminLambdaUrl"), &awscdk.CfnOutputProps{
		ExportName: jsii.Sprintf("UndertownAdminLambdaUrl-%s", props.Env),
		Value:      lambdaURL.Url(),
	})

	props.PUQueue.GrantSendMessages(adminSSRLambda)

	return SSRAdminStackReturn{
		LambdaUrl: lambdaURL,
		Stack:     stack,
	}
}
