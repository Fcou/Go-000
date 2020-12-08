package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"Go-000-main/Week01/grpc_example/public"

	pb "Go-000-main/Week01/grpc_example/proto"
)

const PORT = "9004"

type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return true
}

func main() {
	tlsClient := public.Client{
		ServerName: "grpc_ca",
		CertFile:   "../conf/server/server.pem",
	}
	c, err := tlsClient.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsClient.GetTLSCredentials err: %v", err)
	}

	auth := Auth{
		AppKey:    "fcou",
		AppSecret: "2020",
	}
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c), grpc.WithPerRPCCredentials(&auth))
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
