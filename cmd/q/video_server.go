package main

import (
	"context"
	"encoding/json"
	"strconv"

	pb "github.com/MovAh13h/vidsy/go/pb"
	awsclient "github.com/MovAh13h/vidsy/internal/aws"
	"github.com/MovAh13h/vidsy/internal/common"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type videoServer struct {
	pb.UnimplementedJobServiceServer
	AwsClient *awsclient.AwsClient
	SqsQueueUrl *string
}

func NewVideoServer(name, profile *string) (*videoServer, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(*profile))

    if err != nil {
        return nil, err
    }

	awsClient, err := awsclient.NewAwsClient(&cfg)

	if err != nil {
		return nil, err
	}

	o, err := awsClient.Sqs.CreateQueueIfNotExist(context.Background(), &sqs.CreateQueueInput{
		QueueName: name,
	})

	if err != nil {
		return nil, err
	}

	return &videoServer{AwsClient: awsClient, SqsQueueUrl: o.QueueUrl}, nil
}

func (s *videoServer) QueueJob(ctx context.Context, in *pb.JobRequest) (*pb.JobResponse, error) {
	var entries []types.SendMessageBatchRequestEntry

	for i, resolution := range in.GetConvResolutions() {
		b, err := json.Marshal(common.QueueJob {
			Src: in.GetSrcPath(),
			Dest: in.GetDestPath(),
			OutputFormat: in.GetOutFormat().Enum(),
			Resolution: &resolution,
		})

		if err != nil {
			return nil, err
		}

		str := string(b)
		id := strconv.Itoa(i)

		entries = append(entries, types.SendMessageBatchRequestEntry{
			MessageBody: &str,
			Id: &id,
		})
	}

	_, err := s.AwsClient.Sqs.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		QueueUrl: s.SqsQueueUrl,
		Entries: entries,
	})

	if err != nil {
		return nil, err
	}

	return &pb.JobResponse{
		Status: 200,
	}, nil
}

func (s *videoServer) Close() {
	s.AwsClient.Close()
}
