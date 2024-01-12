package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type SSRStackProps struct {
	awscdk.StackProps
}

func SSRStack(scope constructs.Construct, id string, props *SSRStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, nil)
	return stack
	// createPropertyLambda := awslambdago.NewGoFunction(stack, )

	// SSR
	// homeLambda := awslambdago.NewGoFunction(stack, jsii.String("ServerSideRender"), &awslambdago.GoFunctionProps{
	// 	Runtime:      awslambda.Runtime_PROVIDED_AL2(),
	// 	MemorySize:   jsii.Number(1024),
	// 	Architecture: awslambda.Architecture_ARM_64(),
	// 	Entry:        jsii.String("../services/ssr/homepage/lambda"),
	// 	Bundling:     bundlingOptions,
	// 	// Environment: &map[string]*string{
	// 	// 	"TABLE_NAME": db.TableName(),
	// 	// },
	// })
	// // Add a Function URL.
	// lambdaURL := homeLambda.AddFunctionUrl(&awslambda.FunctionUrlOptions{
	// 	AuthType: awslambda.FunctionUrlAuthType_NONE,
	// })
	// awscdk.NewCfnOutput(stack, jsii.String("homeLambdaURL"), &awscdk.CfnOutputProps{
	// 	ExportName: jsii.String("homeLambdaURL"),
	// 	Value:      lambdaURL.Url(),
	// })

	// assetsBucket := awss3.NewBucket(stack, jsii.String("assets"), &awss3.BucketProps{
	// 	BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
	// 	Encryption:        awss3.BucketEncryption_S3_MANAGED,
	// 	EnforceSSL:        jsii.Bool(true),
	// 	RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	// 	Versioned:         jsii.Bool(false),
	// })

	// // Allow CloudFront to read from the bucket.
	// cfOAI := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("cfnOriginAccessIdentity"), &awscloudfront.OriginAccessIdentityProps{})
	// cfs := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{})
	// cfs.AddActions(jsii.String("s3:GetBucket*"))
	// cfs.AddActions(jsii.String("s3:GetObject*"))
	// cfs.AddActions(jsii.String("s3:List*"))
	// cfs.AddResources(assetsBucket.BucketArn())
	// cfs.AddResources(jsii.String(fmt.Sprintf("%v/*", *assetsBucket.BucketArn())))
	// cfs.AddCanonicalUserPrincipal(cfOAI.CloudFrontOriginAccessIdentityS3CanonicalUserId())
	// assetsBucket.AddToResourcePolicy(cfs)

	// // Add a CloudFront distribution to route between the public directory and the Lambda function URL.
	// lambdaURLDomain := awscdk.Fn_Select(jsii.Number(2), awscdk.Fn_Split(jsii.String("/"), lambdaURL.Url(), nil))
	// lambdaOrigin := awscloudfrontorigins.NewHttpOrigin(lambdaURLDomain, &awscloudfrontorigins.HttpOriginProps{
	// 	ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	// })
	// cf := awscloudfront.NewDistribution(stack, jsii.String("customerFacing"), &awscloudfront.DistributionProps{
	// 	DefaultBehavior: &awscloudfront.BehaviorOptions{
	// 		AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_ALL(),
	// 		Origin:               lambdaOrigin,
	// 		CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD(),
	// 		OriginRequestPolicy:  awscloudfront.OriginRequestPolicy_ALL_VIEWER_EXCEPT_HOST_HEADER(),
	// 		CachePolicy:          awscloudfront.CachePolicy_CACHING_DISABLED(),
	// 		ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
	// 	},
	// 	PriceClass: awscloudfront.PriceClass_PRICE_CLASS_100,
	// })

	// // Add /assets* to the distribution backed by S3.
	// assetsOrigin := awscloudfrontorigins.NewS3Origin(assetsBucket, &awscloudfrontorigins.S3OriginProps{
	// 	// Get content from the / directory in the bucket.
	// 	OriginPath:           jsii.String("/"),
	// 	OriginAccessIdentity: cfOAI,
	// })
	// cf.AddBehavior(jsii.String("/assets*"), assetsOrigin, nil)

	// // Export the domain.
	// awscdk.NewCfnOutput(stack, jsii.String("cloudFrontDomain"), &awscdk.CfnOutputProps{
	// 	ExportName: jsii.String("cloudfrontDomain"),
	// 	Value:      cf.DomainName(),
	// })

	// // Deploy the contents of the ./assets directory to the S3 bucket.
	// awss3deployment.NewBucketDeployment(stack, jsii.String("assetsDeployment"), &awss3deployment.BucketDeploymentProps{
	// 	DestinationBucket: assetsBucket,
	// 	Sources: &[]awss3deployment.ISource{
	// 		awss3deployment.Source_Asset(jsii.String("../services/ssr/homepage/assets"), nil),
	// 	},
	// 	DestinationKeyPrefix: jsii.String("assets"),
	// 	Distribution:         cf,
	// 	DistributionPaths:    jsii.Strings("/assets*"),
	// })
}
