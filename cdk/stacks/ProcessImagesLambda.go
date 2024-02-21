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

type ProcessImagesLambdaProps struct {
	awscdk.StackProps
	Env                   string
	ProcessedImagesBucket awss3.Bucket
	RawImagesBucket       awss3.Bucket
	PRIQueue              awssqs.Queue
}

type ProcessImagesLambdaReturn struct {
	Lambda awslambdago.GoFunction
	Stack  awscdk.Stack
}

func ProcessImagesLambda(scope constructs.Construct, id string, props *ProcessImagesLambdaProps) ProcessImagesLambdaReturn {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	var processedImagesBucketName string
	var rawImagesBucket string

	if props.ProcessedImagesBucket != nil {
		processedImagesBucketName = *props.ProcessedImagesBucket.BucketName()
	}

	if props.RawImagesBucket != nil {
		rawImagesBucket = *props.RawImagesBucket.BucketName()
	}

	processImagesEnv := map[string]*string{
		env.PROCESSED_IMAGES_BUCKET: jsii.String(processedImagesBucketName),
		env.RAW_IMAGES_BUCKET:       jsii.String(rawImagesBucket),
	}
	s3BucketAccessRole := utils.CreateLambdaBasicRole(stack, "s3fullaccesslambdarole", props.Env)
	s3BucketAccessRole.AddManagedPolicy(
		awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonS3FullAccess")))

	processImages := awslambdago.NewGoFunction(
		stack,
		jsii.Sprintf(
			"ProcessImages-%s",
			props.Env,
		),
		&awslambdago.GoFunctionProps{
			Runtime:      awslambda.Runtime_PROVIDED_AL2(),
			MemorySize:   jsii.Number(1024),
			Architecture: awslambda.Architecture_X86_64(),
			Entry:        jsii.String("../services/events/processImages/lambda"),
			Bundling: &awslambdago.BundlingOptions{
				CgoEnabled:   jsii.Bool(true),
				GoBuildFlags: &[]*string{jsii.String(`-ldflags '-extldflags "-static -s -w"' -tags lambda.norpc`)},
			},
			Environment: &processImagesEnv,
			Role:        s3BucketAccessRole,
			Timeout:     awscdk.Duration_Seconds(jsii.Number(30)),
		},
	)

	// ❗ can't add the event notifications in here
	//  because of some cyclic dependencies bullshit
	// ❗but in the ProcessImagesLambda you can get the bucket name from the dispatched event
	// meaning that there's no need to send the raw-images-bucket as a prop to this lambda

	// props.RawImagesBucket.GrantReadWrite(
	// 	processImages,
	// 	nil,
	// )
	// props.ProcessedImagesBucket.GrantReadWrite(
	// 	processImages,
	// 	nil,
	// )

	// props.RawImagesBucket.AddEventNotification(
	// 	awss3.EventType_OBJECT_CREATED_PUT,
	// 	awss3notifications.NewLambdaDestination(processImages),
	// )

	// props.RawImagesBucket.AddEventNotification(
	// 	awss3.EventType_OBJECT_REMOVED_DELETE,
	// 	awss3notifications.NewLambdaDestination(processImages),
	// 	&awss3.NotificationKeyFilter{},
	// )

	processImages.AddEventSource(
		awslambdaeventsources.NewSqsEventSource(
			props.PRIQueue,
			&awslambdaeventsources.SqsEventSourceProps{},
		),
	)

	return ProcessImagesLambdaReturn{
		Lambda: processImages,
		Stack:  stack,
	}
}
