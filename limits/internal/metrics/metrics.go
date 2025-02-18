package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RequestLatency *prometheus.HistogramVec
	HttpCollector  *prometheus.CounterVec
}

type MetricLabelNames string

func CreateMetrics(namespace, subsystem string) *Metrics {
	var (
		requestLatency = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_latency",
				Help:      "The API methods latency info",
			},
			[]string{"service", "method", "success"})
		httpCollector = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_collector",
				Help:      "The number of http requests",
			},
			[]string{"service", "method", "success"},
		)
	)

	prometheus.MustRegister(httpCollector)
	prometheus.MustRegister(requestLatency)

	return &Metrics{
		HttpCollector:  httpCollector,
		RequestLatency: requestLatency,
	}
}
