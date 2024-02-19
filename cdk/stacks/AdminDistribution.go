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

type AdminProps struct {
	awscdk.StackProps
	AdminLambdaUrl awslambda.FunctionUrl
	Env            string
}

func AdminDistribution(scope constructs.Construct, id string, props *AdminProps) awscloudfront.Distribution {
	stack := awscdk.NewStack(scope, &id, nil)

	adminOriginRequestPolicy := awscloudfront.NewOriginRequestPolicy(stack, jsii.String("admin-origin-request-policy"), &awscloudfront.OriginRequestPolicyProps{
		CookieBehavior: awscloudfront.OriginRequestCookieBehavior_AllowList(
			jsii.String(constants.CookieTokenKey),
		),
		QueryStringBehavior: awscloudfront.OriginRequestQueryStringBehavior_AllowList(
			jsii.String(constants.QUERY_PARAMETER_HUMANREADABLEID),
		),
	})

	adminLambdaUrl := awscdk.Fn_Select(jsii.Number(2), awscdk.Fn_Split(jsii.String("/"), props.AdminLambdaUrl.Url(), nil))

	// Cache Policy to forward cookies for admin
	// Admin page should have no caching
	adminCachePolicy := awscloudfront.NewCachePolicy(stack, jsii.String("admin-cache-policy"), &awscloudfront.CachePolicyProps{
		CachePolicyName: jsii.String("adminSSRCachePolicy"),
		Comment:         jsii.String("Custom cache policy for admin ssr, it's not caching anything, cloudfront used only as a proxy"),
		DefaultTtl:      awscdk.Duration_Seconds(jsii.Number(0)),
		MinTtl:          awscdk.Duration_Seconds(jsii.Number(0)),
		MaxTtl:          awscdk.Duration_Seconds(jsii.Number(0)),
	})

	admninOrigin := awscloudfrontorigins.NewHttpOrigin(adminLambdaUrl, &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	})

	// the homepage and assets should get the most caching TTL
	// defaults to 24 hrs
	adminDistribution := awscloudfront.NewDistribution(stack, jsii.String("admin-ssr-facing"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_ALL(),
			Origin:               admninOrigin,
			CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD(),
			OriginRequestPolicy:  adminOriginRequestPolicy,
			CachePolicy:          adminCachePolicy,
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
		PriceClass: awscloudfront.PriceClass_PRICE_CLASS_100,
	})

	adminDistribution.AddBehavior(jsii.String("/login"), admninOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	adminDistribution.AddBehavior(jsii.String("/login/"), admninOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	adminDistribution.AddBehavior(jsii.String("/submit"), admninOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	adminDistribution.AddBehavior(jsii.String("/submit/"), admninOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	adminDistribution.AddBehavior(jsii.String("/edit"), admninOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	adminDistribution.AddBehavior(jsii.String("/edit/*"), admninOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	adminDistribution.AddBehavior(jsii.String("/delete"), admninOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	adminDistribution.AddBehavior(jsii.String("/delete/*"), admninOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	adminDistribution.AddBehavior(jsii.String("/"), admninOrigin, &awscloudfront.AddBehaviorOptions{
		AllowedMethods:      awscloudfront.AllowedMethods_ALLOW_ALL(),
		CachedMethods:       awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		CachePolicy:         adminCachePolicy,
		OriginRequestPolicy: adminOriginRequestPolicy,
	})

	// // Export the domain.
	awscdk.NewCfnOutput(stack, jsii.String("adminDistributionDomain"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("adminDistributionDomain"),
		Value:      adminDistribution.DomainName(),
	})

	return adminDistribution
}
