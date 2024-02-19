package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type SharedDistribitionProps struct {
	awscdk.StackProps
	HomeDistribution     awscloudfront.Distribution
	AdminDistribution    awscloudfront.Distribution
	AssetsBucketStack    AssetsBucketStack
	ProcessedImagesStack ProcessedImagesBucketStack
}

func SharedDistributionOrigins(scope constructs.Construct, id string, props SharedDistribitionProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, nil)

	// Add /assets* to the distribution backed by S3.
	assetsOrigin := awscloudfrontorigins.NewS3Origin(
		props.AssetsBucketStack.Bucket,
		&awscloudfrontorigins.S3OriginProps{
			// Get content from the / directory in the bucket.
			OriginPath:           jsii.String("/"),
			OriginAccessIdentity: props.AssetsBucketStack.OAI,
		})

	processedImagesOrigin := awscloudfrontorigins.NewS3Origin(
		props.ProcessedImagesStack.Bucket,
		&awscloudfrontorigins.S3OriginProps{
			// Get content from the / directory in the bucket.
			OriginPath:           jsii.String("/"),
			OriginAccessIdentity: props.ProcessedImagesStack.OAI,
		})

	props.HomeDistribution.AddBehavior(
		jsii.String("/assets*"),
		assetsOrigin,
		&awscloudfront.AddBehaviorOptions{
			CachePolicy:         awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
			OriginRequestPolicy: awscloudfront.OriginRequestPolicy_ALL_VIEWER_EXCEPT_HOST_HEADER(),
		})

	// Add /images* to the distribution backed by S3.
	props.HomeDistribution.AddBehavior(
		jsii.String("/images*"),
		processedImagesOrigin,
		&awscloudfront.AddBehaviorOptions{
			CachePolicy:         awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
			OriginRequestPolicy: awscloudfront.OriginRequestPolicy_ALL_VIEWER_EXCEPT_HOST_HEADER(),
		})

	props.AdminDistribution.AddBehavior(
		jsii.String("/assets*"),
		assetsOrigin,
		&awscloudfront.AddBehaviorOptions{
			CachePolicy:         awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
			OriginRequestPolicy: awscloudfront.OriginRequestPolicy_ALL_VIEWER_EXCEPT_HOST_HEADER(),
		})

	// Add /images* to the distribution backed by S3.
	props.AdminDistribution.AddBehavior(
		jsii.String("/images*"),
		processedImagesOrigin,
		&awscloudfront.AddBehaviorOptions{
			CachePolicy:         awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
			OriginRequestPolicy: awscloudfront.OriginRequestPolicy_ALL_VIEWER_EXCEPT_HOST_HEADER(),
		})

	// // Deploy the contents of the ./assets directory to the S3 bucket.
	awss3deployment.NewBucketDeployment(stack, jsii.String("assetsDeployment"), &awss3deployment.BucketDeploymentProps{
		DestinationBucket: props.AssetsBucketStack.Bucket,
		Sources: &[]awss3deployment.ISource{
			awss3deployment.Source_Asset(jsii.String("../services/ssr/assets"), nil),
		},
		DestinationKeyPrefix: jsii.String("assets"),
		Distribution:         props.HomeDistribution,
		DistributionPaths:    jsii.Strings("/assets*"),
	})

	return stack
}
