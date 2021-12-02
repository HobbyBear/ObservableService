package main

import (
	"ObservableService/pkg/logger"
	"ObservableService/pkg/monitor"
	"ObservableService/servers/postserver/api"
	"context"
	"github.com/gin-gonic/gin"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"net"
	"net/http"
)

func main() {

	monitor.Init(monitor.WithMetricsPort(8080), monitor.WithTraceReporterConfig(&jaegercfg.ReporterConfig{
		CollectorEndpoint: "http://172.27.0.6:14268/api/traces",
	}), monitor.WithTracerServiceName("postserver"))
	defer monitor.Close()

	engine := gin.New()
	engine.Use(monitor.HttpTraceInjection)
	engine.Use(monitor.ApiMetric)

	engine.Handle(http.MethodGet, "postserver/get_post", api.GetPostHandler)

	server := &http.Server{
		Handler: engine,
	}
	// 本机上所有网卡的9090端口都会监听到
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		logger.Error("init listener fail ", zap.Error(err))
		return
	}
	err = server.Serve(listener)
	if err != nil {
		logger.Error("init server fail ", zap.Error(err))
		return
	}
	defer server.Shutdown(context.TODO())

}
