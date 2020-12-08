package main

import (
	pb "Go-000-main/Week01/grpc_example/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const PORT = "9001"

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

func main() {
	// 根据服务端输入的证书文件和密钥构造 TLS 凭证
	c, err := credentials.NewServerTLSFromFile("../conf/server.pem", "../conf/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	// 设置连接选项
	server := grpc.NewServer(grpc.Creds(c))
	pb.RegisterSearchServiceServer(server, &SearchService{})
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}
