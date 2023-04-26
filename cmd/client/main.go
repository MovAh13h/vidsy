package main

import (
	"context"

	"github.com/MovAh13h/vidsy/go/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := pb.NewJobServiceClient(conn)

	// convert to all smaller resolutions
	resolutions := []pb.VideoResolution { 0, 1, 2, 3, 4, 5, 6, 7 } 

	client.QueueJob(context.Background(), &pb.JobRequest{
		OutFormat: pb.VideoOutputFormat_HLS,
		ConvResolutions: resolutions,
		SrcPath: "s3://vidsy-store/mountain600.mp4",
		DestPath: "",
	})
}