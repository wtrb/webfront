package webfront

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	hitCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "webfront_hits",
			Help: "Cumulative hits since startup.",
		},
		[]string{"host"},
	)
	voltageGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "inverter_voltage_vol",
			Help: "Current voltage of an inverter.",
		},
		[]string{"country", "site", "pv", "iv"},
	)
)

func init() {
	prometheus.MustRegister(hitCounter)
	prometheus.MustRegister(voltageGauge)
}

func (s *Server) MetricsHandler() http.Handler {
	return promhttp.Handler()
}
