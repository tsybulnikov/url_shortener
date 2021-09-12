package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"net/url"
	pb "ozonProject/proto"
	"time"
)

const (
	port = ":50051"
)

// server is used to implement proto.UrlShortenerServer.
type server struct {
	pb.UnimplementedUrlShortenerServer
}

// Create implements proto.UrlShortenerServer
func (s *server) Create(ctx context.Context, in *pb.Url) (*pb.ShortUrl, error) {
	log.Printf("Received: %v", in.GetLongUrl())

	_, err := url.ParseRequestURI(in.GetLongUrl())
	if err != nil {
		return &pb.ShortUrl{ShortUrl: "URL is not valid"}, nil
	}

	mode := "create"

	methodResult := DbUrl(in.GetLongUrl(), mode)
	methodResult = "http://" + methodResult
	return &pb.ShortUrl{ShortUrl: methodResult}, nil

}

func (s *server) Get(ctx context.Context, in *pb.ShortUrl) (*pb.Url, error) {
	//log.Printf("Received: %v", in.GetShortUrl())

	mode := "get"
	methodResult := DbUrl(in.GetShortUrl(), mode)

	return &pb.Url{LongUrl: methodResult}, nil

}

func DbUrl(inputUrl string, mode string) string {
	type url struct {
		id int
		longUrl string
		shortUrl string
	}

	connStr := "user=tsybulnikov password=31415 dbname=OzonProjectDb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil{
		log.Fatalf("failed to open DB: %v", err)
	}

	if mode == "get" {
		rows, err := db.Query("select * from urls where shortUrl = $1", inputUrl)
		if err != nil {
			log.Fatalf("failed to query to DB: %v", err)
		}
		defer rows.Close()

		var urls []url

		for rows.Next() {
			u := url{}
			err := rows.Scan(&u.id, &u.longUrl, &u.shortUrl)
			if err != nil {
				log.Fatalf("failed to scan rows: %v", err)
			}
			urls = append(urls, u)
		}

		for _, u := range urls {
			return u.longUrl
		}
		return "DB has no this short URL"
	} else {
		link := "n.ts/"
		code := ""
		rand.Seed(time.Now().UnixNano())
		passwordSymbols := make([]rune, 63, 63)
		passwordSymbols = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_")
		for i := 1; i <= 5; i++ {
			code = code + string(passwordSymbols[rand.Intn(len(passwordSymbols))])
		}
		result := link + code

		rows, err := db.Query("select * from urls where longUrl = $1", inputUrl)
		if err != nil {
			log.Fatalf("failed to query to DB: %v", err)
		}
		defer rows.Close()

		var urls []url

		for rows.Next() {
			u := url{}
			err := rows.Scan(&u.id, &u.longUrl, &u.shortUrl)
			if err != nil {
				log.Fatalf("failed to scan rows: %v", err)
			}
			urls = append(urls, u)
		}

		for _, u := range urls {
			return u.shortUrl
		}

		_, err = db.Exec("insert into urls (longUrl, shortUrl) values ($1, $2)", inputUrl,
			result)
		if err != nil {
			log.Fatalf("failed to exec DB: %v", err)
		}
		defer db.Close()
		return result
	}
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
