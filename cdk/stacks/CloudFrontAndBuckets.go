package stacks

import (
	"github.com/Serares/undertown_v3/utils/constants"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type BucketProps struct {
	awscdk.StackProps
	HomeLambdaUrl        awslambda.FunctionUrl
	AdminLambdaUrl       awslambda.FunctionUrl
	AssetsBucketStack    AssetsBucketStack
	ProcessedImagesStack ProcessedImagesBucketStack
	Env                  string
}

func CloudFrontAndBuckets(scope constructs.Construct, id string, props *BucketProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, nil)
	// ❗The bucket is created in a previous stack
	// because the arn has to be passed to lambdas before creating cloudfront
	// assetsBucket := awss3.NewBucket(stack, jsii.String("assets"), &awss3.BucketProps{
	// 	BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
	// 	Encryption:        awss3.BucketEncryption_S3_MANAGED,
	// 	EnforceSSL:        jsii.Bool(true),
	// 	RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	// 	Versioned:         jsii.Bool(false),
	// })

	// origin request policy for admin
	// to allow query string and cookie to be forwarded

	adminOriginRequestPolicy := awscloudfront.NewOriginRequestPolicy(stack, jsii.String("admin-origin-request-policy"), &awscloudfront.OriginRequestPolicyProps{
		CookieBehavior: awscloudfront.OriginRequestCookieBehavior_AllowList(
			jsii.String(constants.CookieTokenKey),
		),
		QueryStringBehavior: awscloudfront.OriginRequestQueryStringBehavior_AllowList(
			jsii.String(constants.HumanReadableIdQueryKey),
		),
	})

	homepageOriginRequestPolicy := awscloudfront.NewOriginRequestPolicy(stack, jsii.String("homepage-origin-request-policy"), &awscloudfront.OriginRequestPolicyProps{
		QueryStringBehavior: awscloudfront.OriginRequestQueryStringBehavior_AllowList(
			jsii.String(constants.HumanReadableIdQueryKey),
		),
	})

	// Cache Policy for chirii/vanzari
	// bacause the user can sort the properties
	// caching will just return the same page without sorted content
	// Also used for single property because it might get updated
	propertiesCachePolicy := awscloudfront.NewCachePolicy(stack, jsii.String("properties-cache-policy"), &awscloudfront.CachePolicyProps{
		CachePolicyName:     jsii.String("homepageProperties"),
		Comment:             jsii.String("Custom cache policy for properties pages"),
		DefaultTtl:          awscdk.Duration_Seconds(jsii.Number(0)),
		MinTtl:              awscdk.Duration_Seconds(jsii.Number(0)),
		MaxTtl:              awscdk.Duration_Seconds(jsii.Number(0)),
		CookieBehavior:      awscloudfront.CacheCookieBehavior_None(),
		HeaderBehavior:      awscloudfront.CacheHeaderBehavior_None(),
		QueryStringBehavior: awscloudfront.CacheQueryStringBehavior_None(),
	})

	// Cache Policy to forward cookies for admin
	// Admin page should have no caching
	adminCachePolicy := awscloudfront.NewCachePolicy(stack, jsii.String("admin-cache-policy"), &awscloudfront.CachePolicyProps{
		CachePolicyName: jsii.String("adminSSRCachePolicy"),
		Comment:         jsii.String("Custom cache policy for admin ssr, it's not caching anything, cloudfront used only as a proxy"),
		DefaultTtl:      awscdk.Duration_Seconds(jsii.Number(0)),
		MinTtl:          awscdk.Duration_Seconds(jsii.Number(0)),
		MaxTtl:          awscdk.Duration_Seconds(jsii.Number(0)),
	})

	// // Add a CloudFront distribution to route between the public directory and the Lambda function URL.
	homeLambdaUrl := awscdk.Fn_Select(jsii.Number(2), awscdk.Fn_Split(jsii.String("/"), props.HomeLambdaUrl.Url(), nil))
	adminLambdaUrl := awscdk.Fn_Select(jsii.Number(2), awscdk.Fn_Split(jsii.String("/"), props.AdminLambdaUrl.Url(), nil))
	lambdaOrigin := awscloudfrontorigins.NewHttpOrigin(homeLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})

	// the homepage and assets should get the most caching TTL
	// defaults to 24 hrs
	cf := awscloudfront.NewDistribution(stack, jsii.String("cdn-ssr-facing"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_ALL(),
			Origin:               lambdaOrigin,
			CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD(),
			OriginRequestPolicy:  homepageOriginRequestPolicy,
			CachePolicy:          propertiesCachePolicy,
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
		PriceClass: awscloudfront.PriceClass_PRICE_CLASS_100,
	})

	// Add /assets* to the distribution backed by S3.
	assetsOrigin := awscloudfrontorigins.NewS3Origin(
		props.AssetsBucketStack.Bucket,
		&awscloudfrontorigins.S3OriginProps{
			// Get content from the / directory in the bucket.
			OriginPath:           jsii.String("/"),
			OriginAccessIdentity: props.AssetsBucketStack.OAI,
		})
	cf.AddBehavior(jsii.String("/assets*"), assetsOrigin, &awscloudfront.AddBehaviorOptions{
		CachePolicy:         awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
		OriginRequestPolicy: awscloudfront.OriginRequestPolicy_ALL_VIEWER_EXCEPT_HOST_HEADER(),
	})

	// Add /images* to the distribution backed by S3.
	processedImagesOrigin := awscloudfrontorigins.NewS3Origin(
		props.ProcessedImagesStack.Bucket,
		&awscloudfrontorigins.S3OriginProps{
			// Get content from the / directory in the bucket.
			OriginPath:           jsii.String("/"),
			OriginAccessIdentity: props.ProcessedImagesStack.OAI,
		})
	cf.AddBehavior(jsii.String("/images*"), processedImagesOrigin, nil)

	// Add /properties origins chirii|vanzari
	chiriiOrigin := awscloudfrontorigins.NewHttpOrigin(homeLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})
	cf.AddBehavior(jsii.String("/chirii"), chiriiOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         propertiesCachePolicy,
		OriginRequestPolicy: homepageOriginRequestPolicy,
	})
	vanzariOrigin := awscloudfrontorigins.NewHttpOrigin(homeLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})
	cf.AddBehavior(jsii.String("/vanzari"), vanzariOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         propertiesCachePolicy,
		OriginRequestPolicy: homepageOriginRequestPolicy,
	})
	// // Export the domain.
	awscdk.NewCfnOutput(stack, jsii.String("cloudFrontDomain"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("cloudfrontDomain"),
		Value:      cf.DomainName(),
	})

	// ❗
	// TODO try to create the origins and behavoiurs in a loop
	// Add ADMIN routes
	loginOrigin := awscloudfrontorigins.NewHttpOrigin(adminLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})
	cf.AddBehavior(jsii.String("/login"), loginOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	cf.AddBehavior(jsii.String("/login/"), loginOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	submitPropertyOrigin := awscloudfrontorigins.NewHttpOrigin(adminLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})

	cf.AddBehavior(jsii.String("/submit"), submitPropertyOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	cf.AddBehavior(jsii.String("/submit/"), submitPropertyOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	listOrigin := awscloudfrontorigins.NewHttpOrigin(adminLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})

	cf.AddBehavior(jsii.String("/list"), listOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})
	cf.AddBehavior(jsii.String("/list/"), listOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	editOrigin := awscloudfrontorigins.NewHttpOrigin(adminLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})

	cf.AddBehavior(jsii.String("/edit"), editOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})
	cf.AddBehavior(jsii.String("/edit/*"), editOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	deleteOrigin := awscloudfrontorigins.NewHttpOrigin(adminLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})

	cf.AddBehavior(jsii.String("/delete"), deleteOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})
	cf.AddBehavior(jsii.String("/delete/*"), deleteOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	// add the processed images bucket as a origin and a cache behaviour

	// // Deploy the contents of the ./assets directory to the S3 bucket.
	awss3deployment.NewBucketDeployment(stack, jsii.String("assetsDeployment"), &awss3deployment.BucketDeploymentProps{
		DestinationBucket: props.AssetsBucketStack.Bucket,
		Sources: &[]awss3deployment.ISource{
			awss3deployment.Source_Asset(jsii.String("../services/ssr/assets"), nil),
		},
		DestinationKeyPrefix: jsii.String("assets"),
		Distribution:         cf,
		DistributionPaths:    jsii.Strings("/assets*"),
	})

	return stack
}
