package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client *s3.Client
}

func NewClient(cfg *aws.Config) *Client {
	client := s3.NewFromConfig(*cfg)

	return &Client{client}
}

func (c *Client) GetObject(ctx context.Context, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return c.client.GetObject(ctx, in)
}

func (c *Client) PutObject(ctx context.Context, in *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return c.client.PutObject(ctx, in)
}