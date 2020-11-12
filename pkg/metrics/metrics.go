package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PortierisMetrics struct {
	AllowDecisionCount prometheus.Counter
	DenyDecisionCount  prometheus.Counter

	allMetrics []prometheus.Collector
}

func NewMetrics() *PortierisMetrics {
	p := &PortierisMetrics{}
	p.AllowDecisionCount = p.counter("allow_count", "Allow")
	p.DenyDecisionCount = p.counter("deny_count", "Deny")
	prometheus.MustRegister(p.allMetrics...)
	return p
}

func (p *PortierisMetrics) counter(name, help string) prometheus.Counter {
	result := prometheus.NewCounter(prometheus.CounterOpts{
		Name: metricName(name),
		Help: metricHelp(help),
	})

	p.allMetrics = append(p.allMetrics, result)
	return result
}

func metricName(suffix string) string {
	return fmt.Sprintf("portieris_pod_admission_decision_%s", suffix)
}

func metricHelp(desc string) string {
	return fmt.Sprintf("Portieris count of decision outcomes of %s", desc)
}

// unregisterAll is used by the unit tests to clean up all metrics,
// registered with Prometheus, at the end of a test run.
func (p *PortierisMetrics) UnregisterAll() {
	for _, met := range p.allMetrics {
		prometheus.Unregister(met)
	}
	p.allMetrics = p.allMetrics[:0]
}

// getMetricsHandler is used by the unit tests to retrieve the http.Handler
// that Prometheus provides for retrieving metrics
func (p *PortierisMetrics) GetMetricsHandler() http.Handler {
	return promhttp.Handler()
}
