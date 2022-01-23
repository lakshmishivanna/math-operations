package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	pb "example.com/go-techmgmt-grpc/techmgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50055"
)

func main() {
	// Connection to server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Client: did not connect: %v", err)
	}
	fmt.Println("Client: The status code we got is:", conn.GetState())
	defer conn.Close()
	c := pb.NewTechMangementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		r, err := c.AddNmbers(ctx, &pb.Numbers{Number1: 10, Number2: 20})
		//r, err := c.AddNmbers(ctx, &pb.Numbers{})
		if err != nil {
			log.Fatalf("Client: could not calculate: %v", err)
		}

		s := r.GetHashsum()
		a, _ := base64.StdEncoding.DecodeString(s)

		log.Printf("Client: HashSum: %x", a[:])
		log.Printf("Client: ErrorReponse: %v", r.GetErrorResponse())
	}
}
