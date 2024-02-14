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

type ProcessedImagesBucketStack struct {
	Bucket awss3.Bucket
	OAI    awscloudfront.OriginAccessIdentity
	Stack  awscdk.Stack
}

func ProcessedImagesBucket(scope constructs.Construct, id string) *ProcessedImagesBucketStack {
	stack := awscdk.NewStack(scope, &id, nil)
	processImagesBucket := awss3.NewBucket(
		stack,
		jsii.String("ProcessedImagesBucket"),
		&awss3.BucketProps{
			BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
			Encryption:        awss3.BucketEncryption_S3_MANAGED,
			EnforceSSL:        jsii.Bool(true),
			RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
			Versioned:         jsii.Bool(false),
		})
	// Allow CloudFront to read from the bucket.
	cfOAI := awscloudfront.NewOriginAccessIdentity(
		stack,
		jsii.String("cfnOriginAccessIdentity"),
		&awscloudfront.OriginAccessIdentityProps{},
	)
	resourcePolicy := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{})
	resourcePolicy.AddActions(jsii.String("s3:GetBucket*"))
	resourcePolicy.AddActions(jsii.String("s3:GetObject*"))
	resourcePolicy.AddActions(jsii.String("s3:List*"))
	resourcePolicy.AddResources(processImagesBucket.BucketArn())
	resourcePolicy.AddResources(jsii.String(fmt.Sprintf("%v/*", *processImagesBucket.BucketArn())))
	resourcePolicy.AddCanonicalUserPrincipal(cfOAI.CloudFrontOriginAccessIdentityS3CanonicalUserId())
	processImagesBucket.AddToResourcePolicy(resourcePolicy)
	return &ProcessedImagesBucketStack{
		Stack:  stack,
		Bucket: processImagesBucket,
		OAI:    cfOAI,
	}
}
