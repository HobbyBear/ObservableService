package monitor

import (
	"ObservableService/pkg/logger"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type jaegerLogger struct {
	Logger logger.Logger
}

func (j jaegerLogger) Error(msg string) {
	j.Logger.Error(msg)
}

func (j jaegerLogger) Infof(msg string, args ...interface{}) {
	j.Logger.Info(fmt.Sprintf(msg, args...))
}

func (j jaegerLogger) Debugf(msg string, args ...interface{}) {
	j.Logger.Debug(fmt.Sprintf(msg, args...))
}

func decoratorJaegerLog(logger logger.Logger) jaeger.Logger {
	return jaegerLogger{Logger: logger}
}

func HttpTraceInjection(c *gin.Context) {
	path := strings.Split(c.Request.RequestURI, "?")
	url := ""
	if len(path) != 0 {
		url = path[0]
	}
	span := opentracing.StartSpan(url)
	ext.HTTPUrl.Set(span, c.Request.RequestURI)
	ext.HTTPMethod.Set(span, c.Request.Method)
	defer span.Finish()

	c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))
	c.Next()
	ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
}

// metadataReaderWriter ...
type metadataReaderWriter struct {
	MD map[string][]string
}

// Set ...
func (w metadataReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

// ForeachKey ...
func (w metadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func MetadataInjector(ctx context.Context, md metadata.MD) context.Context {
	span := opentracing.SpanFromContext(ctx)
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, metadataReaderWriter{
		MD: md,
	})
	if err != nil {
		defaultConfig.Logger.Error("inject span ctx", zap.Error(err))
		return ctx
	}
	return metadata.NewOutgoingContext(ctx, md)
}

func TraceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		span, cctx := opentracing.StartSpanFromContext(ctx, method, ext.SpanKindRPCClient)
		defer span.Finish()

		ctx2 := MetadataInjector(cctx, md)
		err := invoker(ctx2, method, req, reply, cc, opts...)
		if err != nil {
			code := codes.Unknown
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			}
			span.SetTag("response_code", code)
			ext.Error.Set(span, true)
		}
		return err
	}
}

func FromIncomingContext(ctx context.Context) opentracing.StartSpanOption {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	sc, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, metadataReaderWriter{MD: md})
	if err != nil {
		return NullStartSpanOption{}
	}
	return ext.RPCServerOption(sc)
}

// NullStartSpanOption ...
type NullStartSpanOption struct{}

// Apply ...
func (sso NullStartSpanOption) Apply(options *opentracing.StartSpanOptions) {}

func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod,
			FromIncomingContext(ctx), ext.SpanKindRPCServer)
		defer span.Finish()

		resp, err = handler(ctx, req)

		if err != nil {
			code := codes.Unknown
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			}
			span.SetTag("code", code)
			ext.Error.Set(span, true)
		}
		return
	}
}
