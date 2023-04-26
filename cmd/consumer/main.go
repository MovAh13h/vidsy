package main

import (
	"flag"
	"log"
)

func main() {
	name := flag.String("name", "default", "Name of the queue service")
	profile := flag.String("profile", "default", "AWS Profile to run the service")

	flag.Parse()

	consumer, err := NewConsumer(name, profile)

	if err != nil {
		log.Fatalf("%v", err)
	}

	consumer.Start()
}