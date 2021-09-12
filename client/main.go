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
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUrlShortenerClient(conn)
	var answer string
	fmt.Print("Введите 'CREATE' для получения сокращенного URL или 'GET' для получения оригинального: ")
	fmt.Scan(&answer)
	if answer == "CREATE" {
		CreatingUrl(c)
	} else if answer == "GET" {
		GettingUrl(c)
	} else {
		fmt.Println("неверный ввод, функция не существует")
	}
}

func CreatingUrl(c pb.UrlShortenerClient) {
	var urlIn string

	fmt.Print("Введите URL: ")
	fmt.Scan(&urlIn)

	if len(os.Args) > 1 {
		urlIn = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Create(ctx, &pb.Url{LongUrl: urlIn})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Short URL: %s", r.GetShortUrl())
}

func GettingUrl(c pb.UrlShortenerClient) {
	var urlIn string

	fmt.Print("Введите URL: ")
	fmt.Scan(&urlIn)

	if len(os.Args) > 1 {
		urlIn = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &pb.ShortUrl{ShortUrl: urlIn})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Original URL: %s", r.GetLongUrl())
}
