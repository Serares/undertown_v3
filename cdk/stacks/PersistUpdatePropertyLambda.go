package stacks

import (
	"cdk/utils"
	"os"

	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type PersistUpdatePropertyLambdaProps struct {
	awscdk.StackProps
	Env               string
	PIUQueue          awssqs.Queue
	DeleteImagesQueue awssqs.Queue
}

func PersistUpdatePropertyLambda(scope constructs.Construct, id string, props *PersistUpdatePropertyLambdaProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	var deleteImagesQueueUrl string
	var piuQueueUrl string

	if props.PIUQueue.QueueUrl() != nil {
		piuQueueUrl = *props.PIUQueue.QueueUrl()
	}

	if props.DeleteImagesQueue.QueueUrl() != nil {
		deleteImagesQueueUrl = *props.DeleteImagesQueue.QueueUrl()
	}

	persistUpdatePropertyEnv := map[string]*string{
		// ❗
		// TODO try to obfuscate somehow the values
		// don't store them in plain text
		// store them as an encrypted string?
		// how to decrypt them
		env.DB_HOST:                 jsii.String(os.Getenv(env.DB_HOST)),
		env.DB_NAME:                 jsii.String(os.Getenv(env.DB_NAME)),
		env.DB_PROTOCOL:             jsii.String(os.Getenv(env.DB_PROTOCOL)),
		env.TURSO_DB_TOKEN:          jsii.String(env.TURSO_DB_TOKEN),
		env.PIU_QUEUE_URL:           jsii.String(piuQueueUrl),
		env.DELETE_IMAGES_QUEUE_URL: jsii.String(deleteImagesQueueUrl), // used to dispatch the names of images that needs to be deleted
	}

	s3BucketAccessRole := utils.CreateLambdaBasicRole(stack, "s3fullaccesslambdarole", props.Env)
	s3BucketAccessRole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonS3FullAccess")))

	persistUpdateProperty := awslambdago.NewGoFunction(
		stack,
		jsii.Sprintf("PersistUpdateProperty-%s", props.Env),
		&awslambdago.GoFunctionProps{
			Runtime:      awslambda.Runtime_PROVIDED_AL2(),
			MemorySize:   jsii.Number(1024),
			Architecture: awslambda.Architecture_ARM_64(),
			Entry:        jsii.String("../services/events/persistUpdateProperty/lambda"),
			Bundling:     BundlingOptions,
			Environment:  &persistUpdatePropertyEnv,
			Role:         s3BucketAccessRole,
			Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		},
	)
	// add the event source mapping to the PersistUpdateProperty lambda
	persistUpdateProperty.AddEventSource(
		awslambdaeventsources.NewSqsEventSource(
			props.PIUQueue,
			&awslambdaeventsources.SqsEventSourceProps{},
		),
	)
	// grant piuqueue access to addPropertyLambda
	props.PIUQueue.GrantConsumeMessages(persistUpdateProperty)

	// ❗
	// dispatch delete permissions
	props.DeleteImagesQueue.GrantSendMessages(persistUpdateProperty)

	return stack
}
