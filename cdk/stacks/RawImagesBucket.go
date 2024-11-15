package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func RawImagesBucket(scope constructs.Construct, id string) *AssetsBucketStack {
	stack := awscdk.NewStack(scope, &id, nil)
	rawImagesBucket := awss3.NewBucket(
		stack,
		jsii.String("RawImagesBucket"),
		&awss3.BucketProps{
			BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
			Encryption:        awss3.BucketEncryption_S3_MANAGED,
			EnforceSSL:        jsii.Bool(true),
			RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
			Versioned:         jsii.Bool(false),
			AutoDeleteObjects: jsii.Bool(true),
		})
	// Allow CloudFront to read from the bucket.
	// resourcePolicy := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{})
	// resourcePolicy.AddActions(jsii.String("s3:GetBucket*"))
	// resourcePolicy.AddActions(jsii.String("s3:GetObject*"))
	// resourcePolicy.AddActions(jsii.String("s3:List*"))
	// resourcePolicy.AddResources(rawImagesBucket.BucketArn())
	// resourcePolicy.AddResources(jsii.String(fmt.Sprintf("%v/*", *rawImagesBucket.BucketArn())))
	// rawImagesBucket.AddToResourcePolicy(resourcePolicy)
	rawImagesBucket.AddCorsRule(
		&awss3.CorsRule{
			AllowedMethods: &[]awss3.HttpMethods{
				awss3.HttpMethods_PUT,
			},
			AllowedOrigins: jsii.Strings(
				*jsii.String("https://undertown.ro"),
				*jsii.String("https://admin.undertown.ro"),
				*jsii.String("http://localhost:4031"),
				*jsii.String("http://127.0.0.1:4031"),
			),
			AllowedHeaders: &[]*string{
				jsii.String("*"),
			},
		},
	)
	return &AssetsBucketStack{
		Stack:  stack,
		Bucket: rawImagesBucket,
	}
}
