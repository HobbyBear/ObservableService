package middlewire

import (
	"ObservableService/constants"
	"bytes"
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"net/http"
)

func Trace(next http.HandlerFunc) http.HandlerFunc  {
	return func(writer http.ResponseWriter, request *http.Request) {

		tracer := opentracing.GlobalTracer()
		if tracer == nil {
			// Tracer not found, just skip.
			next(writer,request)
		}

		buf := bytes.NewBuffer(make([]byte,127))
		buf.WriteString("HTTP ")
		buf.WriteString(request.Method)

		// Start span.
		span := opentracing.StartSpan(buf.String())
		rc := opentracing.ContextWithSpan(request.Context(), span)

		// Set request ID for context.
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			rc = context.WithValue(rc, constants.RequestID, sc.TraceID().String())
		}

		next(writer, request.WithContext(rc))

		// Finish span.
		span.Finish()

	}
}
