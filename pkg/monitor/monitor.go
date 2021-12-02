package monitor

import (
	"ObservableService/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	monitor = "monitor"

	apiHttpType = "http"
	apiGrpcType = "grpc"
)

type OptionsFunc func(config *Config)

type metric struct {
	Path       string
	Port       int
	On         bool
	Collectors []prometheus.Collector
}

type Config struct {
	Metrics metric
	Logger  logger.Logger
	Tracer  Tracer
}

func WithLogger(logger logger.Logger) OptionsFunc {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithMetricsPath(path string) OptionsFunc {
	return func(config *Config) {
		config.Metrics.Path = path
	}
}

func WithMetricsPort(port int) OptionsFunc {
	return func(config *Config) {
		config.Metrics.Port = port
	}
}

func WithMetricsCollector(metric ...prometheus.Collector) OptionsFunc {
	return func(config *Config) {
		config.Metrics.Collectors = append(config.Metrics.Collectors, metric...)
	}
}

type Tracer struct {
	ServiceName string
	Reporter    *jaegercfg.ReporterConfig
	Sample      *jaegercfg.SamplerConfig
	On          bool
}

func Init(optionsFuncs ...OptionsFunc) {
	for _, f := range optionsFuncs {
		f(defaultConfig)
	}
	initMetrics()
	initTracer()
}

func WithTracerServiceName(serviceName string) OptionsFunc {
	return func(config *Config) {
		config.Tracer.ServiceName = serviceName
	}
}

func WithTracerSampleConfig(sampleConfig *jaegercfg.SamplerConfig) OptionsFunc {
	return func(config *Config) {
		config.Tracer.Sample = sampleConfig
		config.Tracer.On = true
	}
}

func WithTraceReporterConfig(reporterConfig *jaegercfg.ReporterConfig) OptionsFunc {
	return func(config *Config) {
		config.Tracer.Reporter = reporterConfig
		config.Tracer.On = true
	}
}

var defaultConfig = &Config{
	Metrics: metric{
		Path:       "/metrics",
		Port:       9090,
		On:         true,
		Collectors: metrics,
	},
	Logger: logger.DefaultLogger,
	Tracer: Tracer{
		ServiceName: "",
		On:          false,
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "http://127.0.0.1:14268/api/traces",
		},
		Sample: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
	},
}

func initMetrics() {
	for _, collector := range defaultConfig.Metrics.Collectors {
		prometheus.MustRegister(collector)
	}
	if defaultConfig.Metrics.On {
		http.Handle(defaultConfig.Metrics.Path, promhttp.Handler())
		go func() {
			err := http.ListenAndServe(":"+strconv.Itoa(defaultConfig.Metrics.Port), nil)
			if err != nil {
				defaultConfig.Logger.Error(monitor, zap.String("err", "listen prometheus port fail"), zap.Error(err), zap.Int("port", defaultConfig.Metrics.Port))
			}
		}()
	}
}

func initTracer() {
	if defaultConfig.Tracer.On {
		tracer, closer, err := traceCfg().NewTracer(
			jaegercfg.Logger(decoratorJaegerLog(defaultConfig.Logger)),
		)
		if err != nil {
			defaultConfig.Logger.Error("tracer init", zap.Error(err))
			return
		}
		appendClose(closer)
		opentracing.InitGlobalTracer(tracer)
	}
}

func traceCfg() *jaegercfg.Configuration {
	return &jaegercfg.Configuration{
		ServiceName: defaultConfig.Tracer.ServiceName,
		Sampler:     defaultConfig.Tracer.Sample,
		Reporter:    defaultConfig.Tracer.Reporter,
	}
}
