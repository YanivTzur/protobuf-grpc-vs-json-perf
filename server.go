package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "03-grpc-json-comparison-go/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const port = ":50051"

type server struct {
	pb.UnimplementedUserManagerServer
}

func (s *server) AddUser(ctx context.Context, in *pb.User) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func parseCertificateFileName() string {
	if os.Args[1] == "--certificate" {
		return os.Args[2]
	}
	return os.Args[4]
}

func parsePrivateKeyFileName() string {
	if os.Args[1] == "--key" {
		return os.Args[2]
	}
	return os.Args[4]
}

func main() {
	if len(os.Args) < 5 || os.Args[1] != "--certificate" && os.Args[1] != "--key" ||
		os.Args[3] != "--certificate" && os.Args[3] != "--key" {
		log.Fatalln("Usage: ./server --certificate [TLS certificate] --key [TLS private key]")
	}

	log.Println("Listening for client gRPC connections on port", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	certificateFileName := parseCertificateFileName()
	privateKeyFileName := parsePrivateKeyFileName()
	creds, _ := credentials.NewServerTLSFromFile(certificateFileName, privateKeyFileName)
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterUserManagerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
