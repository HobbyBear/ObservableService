package monitor

import "github.com/prometheus/client_golang/prometheus"

var metrics = []prometheus.Collector{temp, serverHandleHistogram, serverHandleCounter,
	clientHandleCounter, clientHandleHistogram}

var temp = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace:   monitor,
	Subsystem:   "test",
	Name:        "home_temperature_celsius",
	Help:        "The current temperature in degrees Celsius.",
	ConstLabels: prometheus.Labels{},
})

func TestSetGauge(num int) {
	temp.Add(float64(num))
}

// 设计哪些指标 red 方法的实践

var (
	// qps ,每秒错误请求数 error
	serverHandleCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   monitor,
		Name:        "server_handler_counter",
		ConstLabels: nil,
	}, []string{"method", "api", "code", "type"})

	// 请求耗时
	serverHandleHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: monitor,
		Name:      "server_handle_histogram",
	}, []string{"method", "api", "type"})

	clientHandleCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   monitor,
		Name:        "client_handle_counter",
		ConstLabels: nil,
	}, []string{"method", "target", "code"})

	clientHandleHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: monitor,
		Name:      "client_handle_histogram",
	}, []string{"method", "target"})
)
