package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/MovAh13h/vidsy/go/pb"
	"google.golang.org/grpc"
)


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