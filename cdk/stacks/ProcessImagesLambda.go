package stacks

import (
	"cdk/utils"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ProcessImagesLambdaProps struct {
	awscdk.StackProps
	Env string
}

type ProcessImagesLambdaReturn struct {
	Lambda awslambdago.GoFunction
	Stack  awscdk.Stack
}

func ProcessImagesLambda(scope constructs.Construct, id string, props *ProcessImagesLambdaProps) ProcessImagesLambdaReturn {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	processImagesEnv := map[string]*string{}
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
			Timeout:     awscdk.Duration_Seconds(jsii.Number(3 * 60)),
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

	return ProcessImagesLambdaReturn{
		Lambda: processImages,
		Stack:  stack,
	}
}
