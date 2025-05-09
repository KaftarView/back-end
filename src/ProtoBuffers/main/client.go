package main

import (
	"context"
	"log"
	"time"

	pb "first-project/src/ProtoBuffers"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "Hello, how can you assist me?"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Chat(ctx, &pb.ChatRequest{Message: defaultName})
	if err != nil {
		log.Fatalf("could not chat: %v", err)
	}
	log.Printf("ChatService replied: %s", r.GetReply())
}
