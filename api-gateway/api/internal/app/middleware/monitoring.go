// Package middleware provides monitoring and observability utilities for HTTP requests.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define Prometheus metrics for tracking HTTP requests.
var (
	// httpRequestsTotal counts the total number of HTTP requests received.
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// httpRequestDuration measures the duration of HTTP requests.
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

// init registers the Prometheus metrics.
func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

// PrometheusMiddleware collects HTTP request metrics for monitoring.
//
// It tracks request count, status codes, and request duration using Prometheus metrics.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start the request timer
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()))

		// Process the request
		c.Next()

		// Record request metrics after execution
		status := c.Writer.Status()
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), http.StatusText(status)).Inc()
		timer.ObserveDuration()
	}
}

// PrometheusHandler serves the Prometheus metrics endpoint.
//
// This handler exposes metrics in the Prometheus format, which can be scraped by Prometheus servers.
//
// Returns:
// - gin.HandlerFunc: A Gin handler function for serving metrics.
func PrometheusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}
