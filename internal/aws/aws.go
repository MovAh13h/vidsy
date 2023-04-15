package aws

import (
	"github.com/MovAh13h/vidsy-aws/sqs"
	"github.com/aws/aws-sdk-go-v2/aws"
)


type AwsClient struct {
	sqs *sqs.AwsSqsClient
}

func NewAwsClient(cfg *aws.Config) *AwsClient {
	sqsClient := sqs.NewAwsSqsClient(cfg)

	client := AwsClient {
		sqs: sqsClient,
	}

	return &client
}