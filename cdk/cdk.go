package main

import (
	"cdk/stacks"
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

func main() {
	defer jsii.Close()
	err := godotenv.Load(".env.dev")
	theEnv := os.Getenv("ENV")
	if err != nil {
		// handle error in the cdk stack
		panic(err)
	}
	app := awscdk.NewApp(nil)
	// TODO
	// try to tidy up the queues and event listening lambdas
	piuqueue := stacks.PIUQueue(
		app,
		fmt.Sprintf("PIUQeue-%s", theEnv),
		stacks.PIUQueueProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env: theEnv,
		},
	)

	deleteImagesQueue := stacks.DeleteImagesQueue(
		app,
		fmt.Sprintf("DeleteImagesQueue-%s", theEnv),
		stacks.DeleteImagesQueueProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env: theEnv,
		},
	)

	assetsBucket := stacks.AssetsBucket(
		app,
		fmt.Sprintf("assets-bucket-%s", theEnv),
	)

	rawImagesBucket := stacks.RawImagesBucket(
		app,
		fmt.Sprintf("raw-images-bucket-%s", theEnv),
	)

	processedImagesBucket := stacks.ProcessedImagesBucket(
		app,
		fmt.Sprintf("processed-images-bucket-%s", theEnv),
	)

	authLambdas := stacks.A1Lambda(
		app,
		fmt.Sprintf("A1Lambda-%s", theEnv),
		&stacks.A1LambdaProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env: theEnv,
		},
	)

	crudLambdas := stacks.U1Lambda(
		app,
		fmt.Sprintf("U1Lambda-%s", theEnv),
		&stacks.U1LambdaProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:               theEnv,
			DeleteImagesQueue: deleteImagesQueue.Queue,
		},
	)

	processImagesLambda := stacks.ProcessImagesLambda(
		app,
		fmt.Sprintf("ProcessImages-%s", theEnv),
		&stacks.ProcessImagesLambdaProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:                   theEnv,
			ProcessedImagesBucket: processedImagesBucket.Bucket,
		},
	)

	stacks.DeleteProcessedImagesLambda(
		app,
		fmt.Sprintf("DeleteImagesLambda-%s", theEnv),
		&stacks.DeleteProcessedImagesLambdaProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:                   theEnv,
			ProcessedImagesBucket: processedImagesBucket.Bucket,
			DeleteImagesQueue:     deleteImagesQueue.Queue,
		},
	)

	stacks.PersistUpdatePropertyLambda(
		app,
		fmt.Sprintf("PersistUpdateLambda-%s", theEnv),
		&stacks.PersistUpdatePropertyLambdaProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:               theEnv,
			PIUQueue:          piuqueue.Queue,
			DeleteImagesQueue: deleteImagesQueue.Queue,
		},
	)

	lambdas := append(authLambdas, crudLambdas.Lambdas...)

	apiStack := stacks.API(
		app,
		fmt.Sprintf("Undertown-API-%s", theEnv),
		&stacks.APIStackProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			IntegrationLambdas: lambdas,
			Env:                theEnv,
		})

	ssrStack := stacks.SSRHomepage(
		app,
		fmt.Sprintf("HomepageSSR-%s", theEnv),
		&stacks.SSRHomepageProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env: theEnv,
		})

	adminStack := stacks.SSRAdmin(
		app,
		fmt.Sprintf("AdminSSR-%s", theEnv),
		&stacks.SSRAdminStackProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:             theEnv,
			PIUQueue:        piuqueue.Queue,
			RawImagesBucket: rawImagesBucket.Bucket,
		})

	rawImagesBucket.Bucket.AddEventNotification(
		awss3.EventType_OBJECT_CREATED_PUT,
		awss3notifications.NewLambdaDestination(processImagesLambda.Lambda),
	)

	ssrStack.Stack.AddDependency(apiStack, jsii.String("needs the api gateway getProperty and getProperties paths"))
	adminStack.Stack.AddDependency(apiStack, jsii.String("needs the api gateway crud and login paths"))
	// crudLambdas.Stack.AddDependency(piuqueue.Queue.Stack(), jsii.String("the addProperty lambda needs the piuqueue"))
	// crudLambdas.Stack.AddDependency(assetsBucket.Stack, jsii.String("CRUD lambdas need the Assets bucket stack deployed first"))
	// adminStack.Stack.AddDependency(piuqueue.Queue.Stack(), jsii.String("admin stack needs the queue to dispatch messages"))
	// processImagesLambda.Stack.AddDependency(rawImagesBucket.Stack, jsii.String("process images lambda needs raw images bucket to be deployed"))
	// processImagesLambda.Stack.AddDependency(processedImagesBucket.Stack, jsii.String("process images lambda needs processed images bucket to be deployed"))
	// deleteImagesLambdaStack.AddDependency(deleteImagesQueue.Queue.Stack(), jsii.String("needs the delete images queue to be deployed first"))
	// persistUpdateLambdaStack.AddDependency(piuqueue.Queue.Stack(), jsii.String("need the piu queue "))
	// persistUpdateLambdaStack.AddDependency(deleteImagesQueue.Queue.Stack(), jsii.String("need the delete images queue"))

	cfStack := stacks.CloudFrontAndBuckets(app, fmt.Sprintf("CloudFrontAndBuckets-%s", theEnv), &stacks.BucketProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		HomeLambdaUrl:        ssrStack.LambdaUrl,
		AdminLambdaUrl:       adminStack.LambdaUrl,
		Env:                  theEnv,
		AssetsBucketStack:    *assetsBucket,
		ProcessedImagesStack: *processedImagesBucket,
	})

	cfStack.AddDependency(assetsBucket.Stack, jsii.String("needs the bucket to be deployed first"))

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
