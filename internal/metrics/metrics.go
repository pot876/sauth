package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MetricsHandlerGin() func(c *gin.Context) {
	return func(c *gin.Context) {
		handler := promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				// EnableOpenMetrics: true,
			})

		handler.ServeHTTP(c.Writer, c.Request)
	}
}
