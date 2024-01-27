package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type BucketProps struct {
	awscdk.StackProps
	HomeLambdaUrl awslambda.FunctionUrl
	AssetsBucket  awss3.Bucket
	OAI           awscloudfront.OriginAccessIdentity
	Env           string
}

func CloudFrontAndBuckets(scope constructs.Construct, id string, props *BucketProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, nil)
	assetsBucket := awss3.NewBucket(stack, jsii.String("assets"), &awss3.BucketProps{
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		Encryption:        awss3.BucketEncryption_S3_MANAGED,
		EnforceSSL:        jsii.Bool(true),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		Versioned:         jsii.Bool(false),
	})

	// // Add a CloudFront distribution to route between the public directory and the Lambda function URL.
	lambdaURLDomain := awscdk.Fn_Select(jsii.Number(2), awscdk.Fn_Split(jsii.String("/"), props.HomeLambdaUrl.Url(), nil))
	lambdaOrigin := awscloudfrontorigins.NewHttpOrigin(lambdaURLDomain, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})
	cf := awscloudfront.NewDistribution(stack, jsii.String("customerFacing"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_ALL(),
			Origin:               lambdaOrigin,
			CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD(),
			OriginRequestPolicy:  awscloudfront.OriginRequestPolicy_ALL_VIEWER_EXCEPT_HOST_HEADER(),
			CachePolicy:          awscloudfront.CachePolicy_CACHING_DISABLED(),
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
		PriceClass: awscloudfront.PriceClass_PRICE_CLASS_100,
	})

	// Add /assets* to the distribution backed by S3.
	assetsOrigin := awscloudfrontorigins.NewS3Origin(assetsBucket, &awscloudfrontorigins.S3OriginProps{
		// Get content from the / directory in the bucket.
		OriginPath:           jsii.String("/"),
		OriginAccessIdentity: props.OAI,
	})
	cf.AddBehavior(jsii.String("/assets*"), assetsOrigin, nil)

	// Add /properties origins chirii|vanzari
	chiriiOrigin := awscloudfrontorigins.NewHttpOrigin(lambdaURLDomain, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})
	cf.AddBehavior(jsii.String("/chirii"), chiriiOrigin, nil)
	vanzariOrigin := awscloudfrontorigins.NewHttpOrigin(lambdaURLDomain, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})
	cf.AddBehavior(jsii.String("/vanzari"), vanzariOrigin, nil)
	// // Export the domain.
	awscdk.NewCfnOutput(stack, jsii.String("cloudFrontDomain"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("cloudfrontDomain"),
		Value:      cf.DomainName(),
	})

	// // Deploy the contents of the ./assets directory to the S3 bucket.
	awss3deployment.NewBucketDeployment(stack, jsii.String("assetsDeployment"), &awss3deployment.BucketDeploymentProps{
		DestinationBucket: assetsBucket,
		Sources: &[]awss3deployment.ISource{
			awss3deployment.Source_Asset(jsii.String("../services/ssr/assets"), nil),
		},
		DestinationKeyPrefix: jsii.String("assets"),
		Distribution:         cf,
		DistributionPaths:    jsii.Strings("/assets*"),
	})

	return stack
}
