// Package metrics constructs the metrics the application will track.
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics represents the set of metrics we gather.
//
// Metrics methods should be used to collect metrics from all
// the different parts of the codebase.
// This will keep this package the central authority
// for metrics and metrics won't get lost.
//
// Metrics is composed by other structs where the actual counters are
// stored and managed. They are embedded in Metrics so that their methods
// can be accessed direcly.
type Metrics struct {
	*api
}

// New initializes and returns the Metrics struct containing all
// the metrics collected by our service. It also returns an instrumented
// http handler to serve the collected metrics.
func New() (*Metrics, http.Handler) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())

	api := newAPI()
	m := &Metrics{
		api: api,
	}

	m.api.register(registry)

	return m, promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
