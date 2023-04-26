package aws

import (
	sqsclient "github.com/MovAh13h/vidsy/internal/aws/sqs"
	stsclient "github.com/MovAh13h/vidsy/internal/aws/sts"
	dynamoclient "github.com/MovAh13h/vidsy/internal/aws/dynamo"
	s3client "github.com/MovAh13h/vidsy/internal/aws/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
)


type AwsClient struct {
	Sqs *sqsclient.Client
	Sts *stsclient.Client
	Dynamo *dynamoclient.Client
	S3 *s3client.Client
}

func NewAwsClient(cfg *aws.Config) (*AwsClient, error) {
	stsClient := stsclient.NewClient(cfg)
	sqsClient := sqsclient.NewClient(cfg, stsClient)
	dynamoClient, err := dynamoclient.NewClient(cfg)
	s3Client := s3client.NewClient(cfg)

	if err != nil {
		return nil, err
	}

	client := AwsClient {
		Sqs: sqsClient,
		Sts: stsClient,
		Dynamo: dynamoClient,
		S3: s3Client,
	}

	return &client, nil
}

func (c *AwsClient) Close() {
	c.Dynamo.Close()
}