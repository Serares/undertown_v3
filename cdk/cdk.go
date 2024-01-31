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
	assetsBucket := stacks.AssetsBucket(app, "assets-bucket")

	authLambdas := stacks.A1Lambda(app, fmt.Sprintf("A1Lambda-%s", theEnv), &stacks.A1LambdaProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		Env: theEnv,
	})

	crudLambdas := stacks.U1Lambda(app, fmt.Sprintf("U1Lambda-%s", theEnv), &stacks.U1LambdaProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		Env:          theEnv,
		AssetsBucket: assetsBucket.Bucket,
	})
	lambdas := append(authLambdas, crudLambdas.Lambdas...)
	apiStack := stacks.API(app, fmt.Sprintf("Undertown-API-%s", theEnv), &stacks.APIStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		IntegrationLambdas: lambdas,
		Env:                theEnv,
	})

	ssrStack := stacks.SSR(app, fmt.Sprintf("SSRLambda-%s", theEnv), &stacks.SSRStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		Env: theEnv,
	})

	adminStack := stacks.AdminSSR(app, "UndertownAdmin-Stack", &stacks.AdminSSRStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		Env: theEnv,
	})

	ssrStack.Stack.AddDependency(apiStack, jsii.String("needs the api gateway getProperty and getProperties paths"))
	adminStack.Stack.AddDependency(apiStack, jsii.String("needs the api gateway crud and login paths"))

	cfStack := stacks.CloudFrontAndBuckets(app, fmt.Sprintf("CloudFrontAndBuckets-%s", theEnv), &stacks.BucketProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		HomeLambdaUrl:  ssrStack.LambdaUrl,
		AdminLambdaUrl: adminStack.LambdaUrl,
		Env:            theEnv,
		AssetsBucket:   assetsBucket.Bucket,
		OAI:            assetsBucket.OAI,
	})
	crudLambdas.Stack.AddDependency(assetsBucket.Stack, jsii.String("CRUD lambdas need the Assets bucket stack deployed first"))
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
