package stacks

import (
	"github.com/Serares/undertown_v3/utils/constants"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type HomeProps struct {
	awscdk.StackProps
	HomeLambdaUrl awslambda.FunctionUrl
	Env           string
}

func HomeDistribution(scope constructs.Construct, id string, props *HomeProps) awscloudfront.Distribution {
	stack := awscdk.NewStack(scope, &id, nil)

	homepageOriginRequestPolicy := awscloudfront.NewOriginRequestPolicy(stack, jsii.String("homepage-origin-request-policy"), &awscloudfront.OriginRequestPolicyProps{
		QueryStringBehavior: awscloudfront.OriginRequestQueryStringBehavior_AllowList(
			jsii.String(constants.QUERY_PARAMETER_HUMANREADABLEID),
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

	// // Add a CloudFront distribution to route between the public directory and the Lambda function URL.
	homeLambdaUrl := awscdk.Fn_Select(jsii.Number(2), awscdk.Fn_Split(jsii.String("/"), props.HomeLambdaUrl.Url(), nil))

	homeOrigin := awscloudfrontorigins.NewHttpOrigin(homeLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})

	// the homepage and assets should get the most caching TTL
	// defaults to 24 hrs
	homeDistribution := awscloudfront.NewDistribution(stack, jsii.String("home-ssr-facing"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_ALL(),
			Origin:               homeOrigin,
			CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD(),
			OriginRequestPolicy:  homepageOriginRequestPolicy,
			CachePolicy:          propertiesCachePolicy,
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
		PriceClass: awscloudfront.PriceClass_PRICE_CLASS_100,
	})

	// Add /properties origins chirii|vanzari
	chiriiOrigin := awscloudfrontorigins.NewHttpOrigin(homeLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})
	homeDistribution.AddBehavior(jsii.String("/chirii"), chiriiOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         propertiesCachePolicy,
		OriginRequestPolicy: homepageOriginRequestPolicy,
	})
	vanzariOrigin := awscloudfrontorigins.NewHttpOrigin(homeLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})
	homeDistribution.AddBehavior(jsii.String("/vanzari"), vanzariOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         propertiesCachePolicy,
		OriginRequestPolicy: homepageOriginRequestPolicy,
	})

	// // Export the domain.
	awscdk.NewCfnOutput(stack, jsii.String("homeDistributionDomain"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("homeDistributionDomain"),
		Value:      homeDistribution.DomainName(),
	})

	return homeDistribution
}
