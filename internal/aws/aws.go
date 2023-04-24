package aws

import (
	"github.com/MovAh13h/vidsy/internal/aws/sqs"
	"github.com/aws/aws-sdk-go-v2/aws"
)


type AwsClient struct {
	Sqs *sqs.AwsSqsClient
}

func NewAwsClient(cfg *aws.Config) *AwsClient {
	sqsClient := sqs.NewAwsSqsClient(cfg)

	client := AwsClient {
		Sqs: sqsClient,
	}

	return &client
}