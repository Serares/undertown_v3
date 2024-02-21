package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type PRIQueueProps struct {
	awscdk.StackProps
	Env string
}

type PRIQueueReturn struct {
	Queue awssqs.Queue
}

// 🪪
// P -> Process
// U -> Raw
// I -> Images
func PRIQueue(scope constructs.Construct, id string, props PRIQueueProps) *PRIQueueReturn {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	dlq := awssqs.NewQueue(
		stack,
		jsii.Sprintf("%s-%s-dlq", id, props.Env),
		&awssqs.QueueProps{},
	)
	queue := awssqs.NewQueue(
		stack,
		jsii.Sprintf("%s-%s", id, props.Env),
		&awssqs.QueueProps{
			DeadLetterQueue: &awssqs.DeadLetterQueue{
				MaxReceiveCount: jsii.Number(3),
				Queue:           dlq,
			},
		},
	)

	return &PRIQueueReturn{
		Queue: queue,
	}
}
