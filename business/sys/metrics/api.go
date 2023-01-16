package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	actionLabel = "action"
	apiLabel    = "api"
)

type api struct {
	call    *prometheus.SummaryVec
	extCall *prometheus.SummaryVec
}

func newAPI() *api {
	a := &api{}
	a.call = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_seconds",
			Help:       "The requests latencies in seconds of api calls",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{actionLabel},
	)
	a.extCall = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "ex_http_request_duration_seconds",
			Help:       "The requests latencies in seconds of external api calls",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{actionLabel, apiLabel},
	)
	return a
}

func (a *api) register(registry *prometheus.Registry) {
	registry.MustRegister(a.call, a.extCall)
}

// Trace can be used to register the duration of an event.
type Trace struct {
	tm *prometheus.Timer
}

// Mark registers the time passed since the trace has been created.
func (t *Trace) Mark() {
	t.tm.ObserveDuration()
}

// TraceAPI allows to register the duration of the http requests.
// It returns a Trace that should be Marked (by calling its Mark method)
// after the http request is completed.
func (a *api) TraceAPI(endpoint string) *Trace {
	t := prometheus.NewTimer(a.call.With(prometheus.Labels{
		actionLabel: endpoint,
	}))
	return &Trace{tm: t}
}
