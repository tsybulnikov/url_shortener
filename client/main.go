package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "ozonProject/proto"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()


	c := pb.NewUrlShortenerClient(conn)


	var url string

	fmt.Print("Введите URL: ")
	fmt.Scan(&url)

	if len(os.Args) > 1 {
		url = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//r, err := c.Create(ctx, &pb.UrlRequest{UrlReq: url})
	r, err := c.Create(ctx, &pb.Url{LongUrl: url})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetShortUrl())
}
