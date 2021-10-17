package trace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"io"
	"log"
)

var (
	tracer opentracing.Tracer
	closer io.Closer
)


func Init() {
	// todo 搞成读取配置文件
	cfg := jaegercfg.Configuration{
		ServiceName: "client test", // 对其发起请求的的调用链，叫什么服务
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "http://127.0.0.1:14268/api/traces",
		},
	}

	var err error

	jLogger := jaegerlog.StdLogger
	tracer, closer, err = cfg.NewTracer(
		jaegercfg.Logger(jLogger),
	)
	if err != nil{
		log.Fatal(err)
	}
	opentracing.SetGlobalTracer(tracer)
}

func Close() error {
	return closer.Close()
}