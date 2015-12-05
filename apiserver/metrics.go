package apiserver

import "github.com/prometheus/client_golang/prometheus"

var defacerImageDefaceSum = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "defacer_image_deface_sum",
		Help: "Total defaces",
	},
)

var defacerImageCacheHitsSum = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "defacer_image_cache_hits_sum",
		Help: "Total cache hits",
	},
)

var defacerImageCacheMissSum = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "defacer_image_cache_miss_sum",
		Help: "Total cache miss",
	},
)

var defacerImageCacheItemsCount = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "defacer_image_cache_items_count",
		Help: "Total items cached",
	},
)

var defacerImageResizeCoalesceSum = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "defacer_image_resize_coalesce_sum",
		Help: "Total coalesced image resizes",
	},
)

func init() {
	prometheus.MustRegister(defacerImageDefaceSum)
	prometheus.MustRegister(defacerImageCacheHitsSum)
	prometheus.MustRegister(defacerImageCacheMissSum)
	prometheus.MustRegister(defacerImageCacheItemsCount)
	prometheus.MustRegister(defacerImageResizeCoalesceSum)
}
