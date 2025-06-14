package server

import (
	pb "AuthService/proto"
	"AuthService/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Start() {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
	authSrv := &service.AuthServiceServer{}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, authSrv)
	reflection.Register(s)
	log.Println("Сервис авторизации запущен на порте :50055")
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Не удалось запустить gRPC: %v", err)
	}
}
