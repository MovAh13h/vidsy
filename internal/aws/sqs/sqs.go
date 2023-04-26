package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	stsclient "github.com/MovAh13h/vidsy/internal/aws/sts"
)

type Client struct {
	client *sqs.Client
	sts *stsclient.Client
}

func NewClient(cfg *aws.Config, sts *stsclient.Client) *Client {
	client := sqs.NewFromConfig(*cfg)

	return &Client{client, sts}
}

func (c *Client) CreateQueue(ctx context.Context, input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	return c.client.CreateQueue(ctx, input)
}

func (c *Client) CreateQueueIfNotExist(ctx context.Context, cr_input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	o, err := c.sts.GetCallerIdentity(ctx)

	if err != nil {
		return nil, err
	}

	q, err := c.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: cr_input.QueueName,
		QueueOwnerAWSAccountId: o.Account,
	})

	if err != nil {
		return nil, err
	}

	if q.QueueUrl == nil {
		return c.CreateQueue(ctx, cr_input)	
	}

	return &sqs.CreateQueueOutput{QueueUrl: q.QueueUrl, ResultMetadata: q.ResultMetadata}, nil
}

func (c *Client) DeleteQueue(ctx context.Context, input *sqs.DeleteQueueInput) (*sqs.DeleteQueueOutput, error) {
	return c.client.DeleteQueue(ctx, input)
}

func (c *Client) DeleteMessage(ctx context.Context, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return c.client.DeleteMessage(ctx, input)
}

func (c *Client) DeleteMessageBatch(ctx context.Context, input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	return c.client.DeleteMessageBatch(ctx, input)
}

func (c *Client) GetQueueAttributes(ctx context.Context, input *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error) {
	return c.client.GetQueueAttributes(ctx, input)
}

func (c *Client) GetQueueUrl(ctx context.Context, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return c.client.GetQueueUrl(ctx, input)
}

func (c *Client) ListQueues(ctx context.Context, input *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	return c.client.ListQueues(ctx, input)
}

func (c *Client) PurgeQueue(ctx context.Context, input *sqs.PurgeQueueInput) (*sqs.PurgeQueueOutput, error) {
	return c.client.PurgeQueue(ctx, input)
}

func (c *Client) ReceiveMessage(ctx context.Context, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return c.client.ReceiveMessage(ctx, input)
}

func (c *Client) SendMessage(ctx context.Context, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return c.client.SendMessage(ctx, input)
}

func (c *Client) SendMessageBatch(ctx context.Context, input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	return c.client.SendMessageBatch(ctx, input)
}

func (c *Client) SetQueueAttributes(ctx context.Context, input *sqs.SetQueueAttributesInput) (*sqs.SetQueueAttributesOutput, error) {
	return c.client.SetQueueAttributes(ctx, input)
}
