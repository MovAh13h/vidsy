package common

import "github.com/MovAh13h/vidsy/go/pb"


type QueueJob struct {
	Src string
	Dest string
	OutputFormat *pb.VideoOutputFormat
	Resolution *pb.VideoResolution
}