package main

import (
	"ObservableService/pkg/monitor"
	"ObservableService/services/postservice/impl"
	"ObservableService/services/postservice/pb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		log.Println("init listener fail ", err)
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(monitor.UnaryServerInterceptor))

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	pb.RegisterPostServiceServer(server, &impl.Service{})
	reflection.Register(server)
	server.Serve(listener)

}
