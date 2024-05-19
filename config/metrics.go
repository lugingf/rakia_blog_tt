package config

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const AppName = "rakia_blog_tt"

const (
	labelApp    = "app"
	labelPath   = "path"
	labelCode   = "code"
	labelMethod = "method"
	labelName   = "name"
)

type Metrics struct {
	httpDurationSummary         *prometheus.SummaryVec
	storageQueryDurationSummary *prometheus.SummaryVec
}

func InitMetrics() *Metrics {
	var metrics Metrics
	metrics.httpDurationSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_request_duration_summary",
		Help:       "Duration of HTTP requests Summary.",
		Objectives: map[float64]float64{0.5: 0.5, 0.9: 0.9, 1: 1},
		AgeBuckets: 3,
		MaxAge:     120 * time.Second,
	}, []string{labelApp, labelPath, labelCode, labelMethod})

	metrics.storageQueryDurationSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "storage_query_duration_summary",
		Help:       "Duration of query made to storage.",
		Objectives: map[float64]float64{0.5: 0.5, 0.9: 0.9, 1: 1},
		AgeBuckets: 3,
		MaxAge:     120 * time.Second,
	}, []string{labelApp, labelName})

	return &metrics
}

func (m *Metrics) ObserveHTTPDuration(timeSince time.Time, path string, code int, method string) {
	m.httpDurationSummary.With(map[string]string{
		labelApp:    AppName,
		labelPath:   path,
		labelCode:   strconv.Itoa(code),
		labelMethod: method,
	}).Observe(float64(time.Since(timeSince).Seconds()))
}

func (m *Metrics) ObserveQueryDuration(timeSince time.Time, name string) {
	m.storageQueryDurationSummary.With(map[string]string{
		labelApp:  AppName,
		labelName: name,
	}).Observe(float64(time.Since(timeSince).Seconds()))
}
