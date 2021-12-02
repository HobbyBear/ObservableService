package main

import (
	"ObservableService/pkg/monitor"
	"ObservableService/services/postservice/impl"
	"ObservableService/services/postservice/pb"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	monitor.Init(monitor.WithMetricsPort(8080), monitor.WithTraceReporterConfig(&jaegercfg.ReporterConfig{
		CollectorEndpoint: "http://172.27.0.6:14268/api/traces",
	}), monitor.WithTracerServiceName("postservice"))
	defer monitor.Close()

	listener, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		log.Println("init listener fail ", err)
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(monitor.UnaryInterceptorChain(monitor.TraceUnaryServerInterceptor(), monitor.UnaryMetricServerInterceptor)))

	pb.RegisterPostServiceServer(server, &impl.Service{})
	reflection.Register(server)
	server.Serve(listener)

}
