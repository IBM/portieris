// Copyright 2020  Portieris Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PortierisMetrics implements the metrics for Portieris
type PortierisMetrics struct {
	AllowDecisionCount prometheus.Counter
	DenyDecisionCount  prometheus.Counter

	allMetrics []prometheus.Collector
}

// NewMetrics instantiates PortierisMetrics
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

// UnregisterAll is used by the unit tests to clean up all metrics,
// registered with Prometheus, at the end of a test run.
func (p *PortierisMetrics) UnregisterAll() {
	for _, met := range p.allMetrics {
		prometheus.Unregister(met)
	}
	p.allMetrics = p.allMetrics[:0]
}

// GetMetricsHandler is used by the unit tests to retrieve the http.Handler
// that Prometheus provides for retrieving metrics
func (p *PortierisMetrics) GetMetricsHandler() http.Handler {
	return promhttp.Handler()
}
