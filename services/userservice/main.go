package main

import (
	"ObservableService/services/userservice/impl"
	"ObservableService/services/userservice/pb"
	"ObservableService/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	trace.Init()
	defer trace.Close()

	listener, err := net.Listen("tcp", "127.0.0.1:8083")
	if err != nil {
		log.Println("init listener fail ", err)
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(trace.TraceSpanServerInterceptor()))

	pb.RegisterUserServiceServer(server, &impl.Service{})
	reflection.Register(server)
	server.Serve(listener)

}
