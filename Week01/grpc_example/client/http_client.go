package main

import (
	"Go-000-main/Week01/grpc_example/public"
	"context"
	"log"

	pb "Go-000-main/Week01/grpc_example/proto"
	"google.golang.org/grpc"
)

const PORT = "9003"

func main() {
	tlsClient := public.Client{
		ServerName: "grpc_ca",
		CertFile:   "../conf/server/server.pem",
	}
	c, err := tlsClient.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsClient.GetTLSCredentials err: %v", err)
	}
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse())
}
