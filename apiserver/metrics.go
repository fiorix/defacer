package apiserver

import "github.com/prometheus/client_golang/prometheus"

var defaceCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "defacer_count",
		Help: "Total defaces",
	},
)

func init() {
	prometheus.MustRegister(defaceCounter)
}
