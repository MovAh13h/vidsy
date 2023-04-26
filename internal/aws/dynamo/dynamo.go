package dynamo

import (
	"context"
	"time"

	"cirello.io/dynamolock/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client struct {
	client *dynamodb.Client
	lock *dynamolock.Client
}

func NewClient(cfg *aws.Config) (*Client, error) {
	client := dynamodb.NewFromConfig(*cfg)

	lock, err := dynamolock.New(client,
		"locks", dynamolock.WithHeartbeatPeriod(time.Second),
		dynamolock.WithPartitionKeyName("key"))

	if err != nil {
		return nil, err
	}

	return &Client{client, lock}, nil
}

func (d *Client) AcquireLock(ctx context.Context, lockName *string) (*dynamolock.Lock, error) {
	return d.lock.AcquireLockWithContext(ctx, *lockName)
}

func (d *Client) ReleaseLock(ctx context.Context, lock *dynamolock.Lock) (bool, error) {
	return d.lock.ReleaseLockWithContext(ctx, lock)
}

func (d *Client) Close() {
	d.lock.Close()
}