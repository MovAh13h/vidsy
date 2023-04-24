package sqs

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type AwsSqsClient struct {
	client *sqs.Client
}

func NewAwsSqsClient(cfg *aws.Config) *AwsSqsClient {
	client := sqs.NewFromConfig(*cfg)

	return &AwsSqsClient{client}
}

func (c *AwsSqsClient) CreateQueue(ctx context.Context, input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	return c.client.CreateQueue(ctx, input)
}

func (c *AwsSqsClient) CreateQueueIfNotExist(ctx context.Context, cr_input *sqs.CreateQueueInput, get_input *sqs.ListQueuesInput) (*sqs.CreateQueueOutput, error) {
	o, err := c.client.ListQueues(ctx, get_input)
	if err != nil {
		return nil, err
	}

	for _, url := range o.QueueUrls {
		if strings.Contains(url, *cr_input.QueueName) {
			return &sqs.CreateQueueOutput{
				QueueUrl: &url,
			}, nil
		}
	}

	
	return c.CreateQueue(ctx, cr_input)
}

func (c *AwsSqsClient) DeleteQueue(ctx context.Context, input *sqs.DeleteQueueInput) (*sqs.DeleteQueueOutput, error) {
	return c.client.DeleteQueue(ctx, input)
}

func (c *AwsSqsClient) DeleteMessage(ctx context.Context, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return c.client.DeleteMessage(ctx, input)
}

func (c *AwsSqsClient) DeleteMessageBatch(ctx context.Context, input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	return c.client.DeleteMessageBatch(ctx, input)
}

func (c *AwsSqsClient) GetQueueAttributes(ctx context.Context, input *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error) {
	return c.client.GetQueueAttributes(ctx, input)
}

func (c *AwsSqsClient) GetQueueUrl(ctx context.Context, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return c.client.GetQueueUrl(ctx, input)
}

func (c *AwsSqsClient) ListQueues(ctx context.Context, input *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	return c.client.ListQueues(ctx, input)
}

func (c *AwsSqsClient) PurgeQueue(ctx context.Context, input *sqs.PurgeQueueInput) (*sqs.PurgeQueueOutput, error) {
	return c.client.PurgeQueue(ctx, input)
}

func (c *AwsSqsClient) ReceiveMessage(ctx context.Context, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return c.client.ReceiveMessage(ctx, input)
}

func (c *AwsSqsClient) SendMessage(ctx context.Context, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return c.client.SendMessage(ctx, input)
}

func (c *AwsSqsClient) SetQueueAttributes(ctx context.Context, input *sqs.SetQueueAttributesInput) (*sqs.SetQueueAttributesOutput, error) {
	return c.client.SetQueueAttributes(ctx, input)
}
