package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/MovAh13h/vidsy/go/pb"
	awsclient "github.com/MovAh13h/vidsy/internal/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"google.golang.org/grpc"
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

	awsClient := awsclient.NewAwsClient(&cfg)

	cqi := sqs.CreateQueueInput{
		QueueName: name,
	}

	lqi := sqs.ListQueuesInput{}

	o, err := awsClient.Sqs.CreateQueueIfNotExist(context.Background(), &cqi, &lqi)

	if err != nil {
		return nil, err
	}

	return &videoServer{AwsClient: awsClient, SqsQueueUrl: o.QueueUrl}, nil
}

func (s *videoServer) QueueJob(ctx context.Context, in *pb.JobRequest) (*pb.JobResponse, error) {
	for _, resolution := range in.GetConvResolutions() {
		b, err := json.Marshal(struct{
			Src string
			Dest string
			OutputFormat *pb.VideoOutputFormat
			Resolution *pb.VideoResolution
		}{
			Src: in.GetSrcPath(),
			Dest: in.GetDestPath(),
			OutputFormat: in.GetOutFormat().Enum(),
			Resolution: &resolution,
		})

		if err != nil {
			return nil, err
		}

		str := string(b)

		s.AwsClient.Sqs.SendMessage(ctx, &sqs.SendMessageInput{
			MessageBody: &str,
			QueueUrl: s.SqsQueueUrl,
		})
	}

	return &pb.JobResponse{}, nil
}

func main() {
	port := flag.Int("port", 3000, "Port to run gRPC Server on")
	name := flag.String("name", "default", "Name of the queue service")
	profile := flag.String("profile", "default", "AWS Profile to run the service")
	
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("%v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server, err := NewVideoServer(name, profile)

	if err != nil {
		log.Fatalf("%v", err)
	}

	pb.RegisterJobServiceServer(grpcServer, server)

	log.Fatalf("%v", grpcServer.Serve(lis))
}