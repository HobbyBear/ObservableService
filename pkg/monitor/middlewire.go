package monitor

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
	"time"
)

func ApiMetric(c *gin.Context) {
	start := time.Now()
	uri := strings.Split(c.Request.RequestURI, "?")[0]
	method := c.Request.Method
	c.Next()
	var appCode = "200"
	if c.Writer != nil {
		appCode = strconv.Itoa(c.Writer.Status())
	}
	serverHandleCounter.WithLabelValues(method, uri, appCode, apiHttpType).Inc()
	serverHandleHistogram.WithLabelValues(method, uri, apiHttpType).Observe(time.Since(start).Seconds())
}

func UnaryInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	build := func(c grpc.UnaryServerInterceptor, n grpc.UnaryHandler, info *grpc.UnaryServerInfo) grpc.UnaryHandler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			return c(ctx, req, info, n)
		}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = build(interceptors[i], chain, info)
		}
		return chain(ctx, req)
	}
}

func UnaryMetricClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	beg := time.Now()
	err = invoker(ctx, method, req, reply, cc, opts...)
	var code codes.Code
	if err == nil {
		code = codes.OK
	} else {
		code = status.Code(err)
	}
	clientHandleCounter.WithLabelValues(method, cc.Target(), code.String()).Inc()
	clientHandleHistogram.WithLabelValues(method, cc.Target()).Observe(time.Since(beg).Seconds())
	return err
}

func UnaryMetricServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	beg := time.Now()
	_, err = handler(ctx, req)
	var code codes.Code
	if err == nil {
		code = codes.OK
	} else {
		code = status.Code(err)
	}
	serverHandleCounter.WithLabelValues("", info.FullMethod, code.String(), apiGrpcType).Inc()
	serverHandleHistogram.WithLabelValues("", info.FullMethod, apiGrpcType).Observe(time.Since(beg).Seconds())
	return
}
