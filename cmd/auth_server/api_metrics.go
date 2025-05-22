package main

import "github.com/prometheus/client_golang/prometheus"

type ApiMetrics struct {
	metricsCounters *prometheus.CounterVec
	metricsBuckets  *prometheus.HistogramVec

	loginCounter2XX   prometheus.Counter
	loginCounter4XX   prometheus.Counter
	loginCounter5XX   prometheus.Counter
	refreshCounter2XX prometheus.Counter
	refreshCounter4XX prometheus.Counter
	refreshCounter5XX prometheus.Counter

	loginDurations2XX   prometheus.Observer
	refreshDurations2XX prometheus.Observer
	loginDurationsXXX   prometheus.Observer
	refreshDurationsXXX prometheus.Observer
}

func (a *ApiMetrics) setupMetrics() {
	a.metricsCounters = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_counter",
			Help: "",
		},
		[]string{"prefix", "status"},
	)
	a.loginCounter2XX = a.metricsCounters.WithLabelValues("login", "2XX")
	a.loginCounter4XX = a.metricsCounters.WithLabelValues("login", "4XX")
	a.loginCounter5XX = a.metricsCounters.WithLabelValues("login", "5XX")

	a.refreshCounter2XX = a.metricsCounters.WithLabelValues("refresh", "2XX")
	a.refreshCounter4XX = a.metricsCounters.WithLabelValues("refresh", "4XX")
	a.refreshCounter5XX = a.metricsCounters.WithLabelValues("refresh", "5XX")

	a.metricsBuckets = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_histogram",
			Help:    "",
			Buckets: []float64{.005, .015, .030, .1, .3},
		},
		[]string{"prefix", "status"},
	)
	a.loginDurations2XX = a.metricsBuckets.WithLabelValues("login", "2XX")
	a.loginDurationsXXX = a.metricsBuckets.WithLabelValues("login", "XXX")

	a.refreshDurations2XX = a.metricsBuckets.WithLabelValues("refresh", "2XX")
	a.refreshDurationsXXX = a.metricsBuckets.WithLabelValues("refresh", "XXX")
}

func (a *ApiMetrics) registerMetrics(r prometheus.Registerer) {
	a.setupMetrics()

	r.MustRegister(a.metricsCounters)
	r.MustRegister(a.metricsBuckets)
}
