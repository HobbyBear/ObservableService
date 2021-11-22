package main

import (
	"ObservableService/services/postservice/impl"
	"ObservableService/services/postservice/pb"
	"ObservableService/trace"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

func main() {
	trace.Init()
	defer trace.Close()

	listener, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		log.Println("init listener fail ", err)
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(trace.TraceSpanServerInterceptor()))

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)

	pb.RegisterPostServiceServer(server, &impl.Service{})
	reflection.Register(server)
	server.Serve(listener)

}
