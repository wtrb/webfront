package webfront

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var hitCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "webfront_hits",
		Help: "Cumulative hits since startup.",
	},
	[]string{"host"},
)

func init() {
	prometheus.MustRegister(hitCounter)
}

func (s *Server) MetricsHandler() http.Handler {
	return promhttp.Handler()
}
