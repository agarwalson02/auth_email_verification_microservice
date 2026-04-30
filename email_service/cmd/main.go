package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	delivery "email_service/internal/delivery/grpc"
	"email_service/internal/services"
	pb "email_service/proto"
	"google.golang.org/grpc/reflection"
)

func main() {

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	emailService := services.NewEmailService()
	handler := delivery.NewHandler(emailService)

	pb.RegisterEmailServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Println("Email service running on port 5001")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
