package proposal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/goverland-labs/core-storage/internal/metrics"
)

var metricHandleHistogram = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: groupName,
		Name:      "handle_duration_seconds",
		Help:      "Handle proposal event duration seconds",
		Buckets:   []float64{.001, .005, .01, .025, .05, .1, .5, 1, 2.5, 5, 10},
	}, []string{"type", "error"},
)

var metricSendEventGauge = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: metrics.Namespace,
		Subsystem: groupName,
		Name:      "events",
		Help:      "Send events summary",
	}, []string{"group", "event", metrics.ErrLabel},
)
