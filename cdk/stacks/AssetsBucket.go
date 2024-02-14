package stacks

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AssetsBucketStack struct {
	Bucket awss3.Bucket
	OAI    awscloudfront.OriginAccessIdentity
	Stack  awscdk.Stack
}

func AssetsBucket(scope constructs.Construct, id string) *AssetsBucketStack {
	stack := awscdk.NewStack(scope, &id, nil)
	assetsBucket := awss3.NewBucket(stack, jsii.String("assets"), &awss3.BucketProps{
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		Encryption:        awss3.BucketEncryption_S3_MANAGED,
		EnforceSSL:        jsii.Bool(true),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		Versioned:         jsii.Bool(false),
		AutoDeleteObjects: jsii.Bool(true),
	})
	// // Allow CloudFront to read from the bucket.
	cfOAI := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("cfnOriginAccessIdentity"), &awscloudfront.OriginAccessIdentityProps{})
	cfs := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{})
	cfs.AddActions(jsii.String("s3:GetBucket*"))
	cfs.AddActions(jsii.String("s3:GetObject*"))
	cfs.AddActions(jsii.String("s3:List*"))
	cfs.AddResources(assetsBucket.BucketArn())
	cfs.AddResources(jsii.String(fmt.Sprintf("%v/*", *assetsBucket.BucketArn())))
	cfs.AddCanonicalUserPrincipal(cfOAI.CloudFrontOriginAccessIdentityS3CanonicalUserId())
	assetsBucket.AddToResourcePolicy(cfs)
	return &AssetsBucketStack{
		Stack:  stack,
		Bucket: assetsBucket,
		OAI:    cfOAI,
	}
}
