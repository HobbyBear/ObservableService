package monitor

import (
	"ObservableService/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const monitor = "monitor"

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

func Init(optionsFuncs ...OptionsFunc) {
	for _, f := range optionsFuncs {
		f(defaultConfig)
	}
	initMetrics()
}

var defaultConfig = &Config{
	Metrics: metric{
		Path:       "/metrics",
		Port:       9090,
		On:         true,
		Collectors: metrics,
	},
	Logger: logger.DefaultLogger,
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
