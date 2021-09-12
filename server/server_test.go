package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "ozonProject/proto"
	"testing"
	"time"
)
const (
	address     = "localhost:50051"
)

func TestServer_Create(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c := pb.NewUrlShortenerClient(conn)

	urlIn := ""
	result, _ := c.Create(ctx, &pb.Url{LongUrl: urlIn})
	if result.GetShortUrl() != "URL is not valid" {
		t.Error("Expected \"URL is not valid\", got:", result.GetShortUrl())
	}

	urlIn = "google.com"
	result, _ = c.Create(ctx, &pb.Url{LongUrl: urlIn})
	if result.GetShortUrl() != "URL is not valid" {
		t.Error("Expected \"URL is not valid\", got:", result.GetShortUrl())
	}

	urlIn = "www.site.com"
	result, _ = c.Create(ctx, &pb.Url{LongUrl: urlIn})
	if result.GetShortUrl() != "URL is not valid" {
		t.Error("Expected \"URL is not valid\", got:", result.GetShortUrl())
	}

	urlIn = "http://google.com"
	result, _ = c.Create(ctx, &pb.Url{LongUrl: urlIn})
	if result.GetShortUrl() == "URL is not valid" {
		t.Error("Expected \"short url request\", got:", result.GetShortUrl())
	}
}
func TestServer_Get(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c := pb.NewUrlShortenerClient(conn)

	urlIn := ""
	result, _ := c.Get(ctx, &pb.ShortUrl{ShortUrl: urlIn})
	if result.GetLongUrl() != "DB has no this short URL" {
		t.Error("Expected \"DB has no this short URL\", got:", result.GetLongUrl())
	}

	urlIn = "www.dadada.yellow.com"
	result, _ = c.Get(ctx, &pb.ShortUrl{ShortUrl: urlIn})
	if result.GetLongUrl() != "DB has no this short URL" {
		t.Error("Expected \"DB has no this short URL\", got:", result.GetLongUrl())
	}
}
func TestDbUrl(t *testing.T) {
	inputUrl := "google.com"
	mode := "Generate"
	result := DbUrl(inputUrl, mode)
	if result != "invalid value for variable \"mode\"" {
		t.Error("invalid value for variable \"mode\", got:", result)
	}

}