package main

import (
	"ObservableService/services/postservice/impl"
	"ObservableService/services/postservice/pb"
	"ObservableService/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	trace.Init()
	defer trace.Close()

	listener,err := net.Listen("tcp","127.0.0.1:8082")
	if err != nil{
		log.Println("init listener fail ", err)
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(trace.TraceSpanServerInterceptor()))

	pb.RegisterPostServiceServer(server,&impl.Service{})
	reflection.Register(server)
	server.Serve(listener)

}



