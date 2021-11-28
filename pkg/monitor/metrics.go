package monitor

import "github.com/prometheus/client_golang/prometheus"

var metrics = []prometheus.Collector{temp}

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
