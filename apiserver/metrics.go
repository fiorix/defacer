package apiserver

import "github.com/prometheus/client_golang/prometheus"

var defaceCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "defacer_count",
		Help: "Total defaces",
	},
)

var defaceImageCache = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "defacer_image_cache",
		Help: "Total cache activity",
	},
	[]string{"counter"},
)

func init() {
	prometheus.MustRegister(defaceCounter)
	prometheus.MustRegister(defaceImageCache)
}
