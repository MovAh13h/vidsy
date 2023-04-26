package sts

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Client struct {
	client *sts.Client
}

func NewClient(cfg *aws.Config) *Client {
	client := sts.NewFromConfig(*cfg)

	return &Client{client}
}

func (c *Client) GetCallerIdentity(ctx context.Context) (*sts.GetCallerIdentityOutput, error) {
	return c.client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
}