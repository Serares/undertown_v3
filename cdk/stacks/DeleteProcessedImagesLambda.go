package stacks

import (
	"cdk/utils"

	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DeleteProcessedImagesLambdaProps struct {
	awscdk.StackProps
	Env                   string
	DeleteImagesQueue     awssqs.Queue
	ProcessedImagesBucket awss3.Bucket
}

func DeleteProcessedImagesLambda(scope constructs.Construct, id string, props *DeleteProcessedImagesLambdaProps) awscdk.Stack {
	var processedImagesBucket string
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	s3BucketAccessRole := utils.CreateLambdaBasicRole(stack, "s3fullaccesslambdarole", props.Env)
	s3BucketAccessRole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonS3FullAccess")))

	if props.ProcessedImagesBucket.BucketName() != nil {
		processedImagesBucket = *props.ProcessedImagesBucket.BucketName()
	}

	deleteImagesLambdaEnv := map[string]*string{
		env.PROCESSED_IMAGES_BUCKET: jsii.String(processedImagesBucket),
	}

	deleteProcessImages := awslambdago.NewGoFunction(
		stack,
		jsii.Sprintf("DeleteProcessedImages-%s", props.Env),
		&awslambdago.GoFunctionProps{
			Runtime:      awslambda.Runtime_PROVIDED_AL2(),
			MemorySize:   jsii.Number(1024),
			Architecture: awslambda.Architecture_ARM_64(),
			Entry:        jsii.String("../services/events/deleteProcessedImages/lambda"),
			Bundling:     BundlingOptions,
			Environment:  &deleteImagesLambdaEnv,
			Role:         s3BucketAccessRole,
			Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		},
	)

	deleteProcessImages.AddEventSource(
		awslambdaeventsources.NewSqsEventSource(
			props.DeleteImagesQueue,
			&awslambdaeventsources.SqsEventSourceProps{},
		),
	)

	props.DeleteImagesQueue.GrantConsumeMessages(
		deleteProcessImages,
	)

	return stack

}
