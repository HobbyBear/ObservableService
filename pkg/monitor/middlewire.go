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

func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
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

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

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
