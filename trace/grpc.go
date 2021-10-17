package trace

import (
	"ObservableService/constants"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

func TraceSpanClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string, req, resp interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) (err error) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "RPC Client "+method)
		defer span.Finish()

		// Save current span context.
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		if err = opentracing.GlobalTracer().Inject(
			span.Context(), opentracing.HTTPHeaders, metadataTextMap(md),
		); err != nil {
			log.Println(ctx, "Failed to inject trace span: %v", err)
		}
		return invoker(metadata.NewOutgoingContext(ctx, md), method, req, resp, cc, opts...)
	}
}


func TraceSpanServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Extract parent trace span.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		parentSpanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders, metadataTextMap(md),
		)
		switch err {
		case nil:
		case opentracing.ErrSpanContextNotFound:
			log.Println(ctx, "Parent span not found, will start new one.")
		default:
			log.Println(ctx, "Failed to extract trace span: %v", err)
		}

		// Start new trace span.
		span := opentracing.StartSpan(
			"RPC Server "+info.FullMethod,
			ext.RPCServerOption(parentSpanContext),
		)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)

		// Set request ID for context.
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			ctx = context.WithValue(ctx, constants.RequestID, sc.TraceID().String())
		}

		return handler(ctx, req)
	}
}


// metadataTextMap extends a metadata.MD to be an opentracing textmap
type metadataTextMap metadata.MD

// Set is a opentracing.TextMapReader interface that extracts values.
func (m metadataTextMap) Set(key, val string) {
	// gRPC allows for complex binary values to be written.
	encodedKey, encodedVal := encodeKeyValue(key, val)
	// The metadata object is a multimap, and previous values may exist, but for opentracing headers, we do not append
	// we just override.
	m[encodedKey] = []string{encodedVal}
}

// ForeachKey is a opentracing.TextMapReader interface that extracts values.
func (m metadataTextMap) ForeachKey(callback func(key, val string) error) error {
	for k, vv := range m {
		for _, v := range vv {
			if decodedKey, decodedVal, err := metadata.DecodeKeyValue(k, v); err == nil {
				if err = callback(decodedKey, decodedVal); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed decoding opentracing from gRPC metadata: %v", err)
			}
		}
	}
	return nil
}

const (
	binHeaderSuffix = "_bin"
)

// encodeKeyValue encodes key and value qualified for transmission via gRPC.
// note: copy pasted from private values of grpc.metadata
func encodeKeyValue(k, v string) (string, string) {
	k = strings.ToLower(k)
	if strings.HasSuffix(k, binHeaderSuffix) {
		val := base64.StdEncoding.EncodeToString([]byte(v))
		v = string(val)
	}
	return k, v
}

