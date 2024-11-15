package main

import (
	"cdk/stacks"
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
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
	PUqueue := stacks.PUQueue(
		app,
		fmt.Sprintf("PUQeue-%s", theEnv),
		stacks.PUQueueProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env: theEnv,
		},
	)

	PRIQueue := stacks.PRIQueue(
		app,
		fmt.Sprintf("PRIQueue-%s", theEnv),
		stacks.PRIQueueProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env: theEnv,
		},
	)

	// ❗
	// grant send message permissions to persistUpdateProperty and deleteProperty lambdas
	// grant consume permission to deleteProcessedImagesLambda
	deleteProcessedImagesQueue := stacks.DeleteProcessedImagesQueue(
		app,
		fmt.Sprintf("DeleteProcessedImagesQueue-%s", theEnv),
		stacks.DeleteProcessedImagesQueueProps{
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

	rawImagesBucketStack := stacks.RawImagesBucket(
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

	// deleteProperty lambda is going to send messages to deleteImagesQueue
	crudLambdas := stacks.U1Lambda(
		app,
		fmt.Sprintf("U1Lambda-%s", theEnv),
		&stacks.U1LambdaProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:                        theEnv,
			DeleteProcessedImagesQueue: deleteProcessedImagesQueue.Queue,
		},
	)

	// processImagesLambda will poll the PRIQueue
	// needs access to ImagesBucket
	stacks.ProcessImagesLambda(
		app,
		fmt.Sprintf("ProcessImages-%s", theEnv),
		&stacks.ProcessImagesLambdaProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:                   theEnv,
			ProcessedImagesBucket: processedImagesBucket.Bucket,
			RawImagesBucket:       rawImagesBucketStack.Bucket,
			PRIQueue:              PRIQueue.Queue,
		},
	)

	// polls from sqs deleteImagesQueue
	// needs access to S3 ImagesBucket
	stacks.DeleteProcessedImagesLambda(
		app,
		fmt.Sprintf("DeleteProcessedImagesLambda-%s", theEnv),
		&stacks.DeleteProcessedImagesLambdaProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:                   theEnv,
			ProcessedImagesBucket: processedImagesBucket.Bucket,
			DeleteImagesQueue:     deleteProcessedImagesQueue.Queue,
		},
	)

	// Polls from PUQueue
	// sends messages to DeleteImagesQueue
	// sends messages to PRIQueue
	stacks.PersistUpdatePropertyLambda(
		app,
		fmt.Sprintf("PersistUpdateLambda-%s", theEnv),
		&stacks.PersistUpdatePropertyLambdaProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:                        theEnv,
			PUQueue:                    PUqueue.Queue,
			PRIQueue:                   PRIQueue.Queue,
			DeleteProcessedImagesQueue: deleteProcessedImagesQueue.Queue,
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

	adminSSRLambda := stacks.SSRAdmin(
		app,
		fmt.Sprintf("AdminSSR-%s", theEnv),
		&stacks.SSRAdminStackProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			Env:             theEnv,
			PUQueue:         PUqueue.Queue,
			RawImagesBucket: rawImagesBucketStack.Bucket,
		})

	// ❗
	// Attach the S3 event to the process images lambda
	// Deprecated
	// rawImagesBucketStack.Bucket.AddEventNotification(
	// 	awss3.EventType_OBJECT_CREATED_PUT,
	// 	awss3notifications.NewLambdaDestination(processImagesLambda.Lambda),
	// )

	ssrStack.Stack.AddDependency(apiStack, jsii.String("needs the api gateway getProperty and getProperties paths"))
	adminSSRLambda.Stack.AddDependency(apiStack, jsii.String("needs the api gateway crud and login paths"))

	homeDistribution := stacks.HomeDistribution(
		app,
		fmt.Sprintf("HomeDistribution-%s", theEnv),
		&stacks.HomeProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			HomeLambdaUrl: ssrStack.LambdaUrl,
			Env:           theEnv,
		})

	adminDistribution := stacks.AdminDistribution(
		app,
		fmt.Sprintf("AdminDistribution-%s", theEnv),
		&stacks.AdminProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			AdminLambdaUrl: adminSSRLambda.LambdaUrl,
		},
	)

	stacks.SharedDistributionOrigins(
		app,
		fmt.Sprintf("SharedOriginsStack-%s", theEnv),
		stacks.SharedDistribitionProps{
			StackProps: awscdk.StackProps{
				Env: env(),
			},
			HomeDistribution:     homeDistribution,
			AdminDistribution:    adminDistribution,
			AssetsBucketStack:    *assetsBucket,
			ProcessedImagesStack: *processedImagesBucket,
		},
	)

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
