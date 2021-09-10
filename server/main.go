/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "ozonProject/proto"
)

const (
	port = ":50051"
)

// server is used to implement proto.UrlShortenerServer.
type server struct {
	pb.UnimplementedUrlShortenerServer
}

// GenerateShortUrl implements proto.UrlShortenerServer
func (s *server) GenerateShortUrl(ctx context.Context, in *pb.UrlRequest) (*pb.UrlResponse, error) {
	log.Printf("Received: %v", in.GetUrlReq())
	link := "n.ts/"
	result := ""
	rand.Seed(time.Now().UnixNano())
	passwordSymbols := make([]rune, 63, 63)
	passwordSymbols = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_")
	for i := 1; i <= 7; i++ {
		result = result + string(passwordSymbols[rand.Intn(len(passwordSymbols))])
	}
	fmt.Printf("%s%s",link,result)
	return &pb.UrlResponse{UrlResp: "Short URL: " + link + result}, nil
	//return &pb.UrlResponse{UrlResp: "Short URL: " + in.GetUrlReq()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUrlShortenerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
