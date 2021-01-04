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
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// return value of the given metric (or error if not present in metrics)
func getMetric(m *PortierisMetrics, metricName string) (string, error) {
	var req *http.Request
	var err error
	var line []byte
	var result string

	if req, err = http.NewRequest("GET", "/metrics", nil); err != nil {
		return result, err
	}
	rr := httptest.NewRecorder()
	m.GetMetricsHandler().ServeHTTP(rr, req)

	for err == nil {
		line, err = rr.Body.ReadBytes('\n')
		lineString := string(line)
		if strings.HasPrefix(lineString, metricName) {
			result = strings.TrimSpace(strings.Trim(lineString, metricName))
		}
	}
	if err == io.EOF {
		err = nil
	}

	return result, err
}

type MockRegisterer struct {
	numRegistered int
}

func (mr *MockRegisterer) Register(c prometheus.Collector) error {
	mr.numRegistered++
	return nil
}

func (mr *MockRegisterer) MustRegister(c ...prometheus.Collector) {
	mr.numRegistered += len(c)
}

func (mr *MockRegisterer) Unregister(c prometheus.Collector) bool {
	mr.numRegistered--
	return true
}

func TestMetricRegisterUnregister(t *testing.T) {
	r := &MockRegisterer{}
	dr := prometheus.DefaultRegisterer
	defer func() {
		prometheus.DefaultRegisterer = dr
	}()
	prometheus.DefaultRegisterer = r

	m := NewMetrics()
	assert.True(t, r.numRegistered > 0)

	m.UnregisterAll()
	assert.Equal(t, 0, r.numRegistered)
}

func TestAllowDecisionMetric(t *testing.T) {

	pm := NewMetrics()
	defer pm.UnregisterAll()
	pm.AllowDecisionCount.Inc()
	metric, metricErr := getMetric(pm, "portieris_pod_admission_decision_allow_count")

	assert.Nil(t, metricErr)

	assert.NotEqual(t, "", metric)
	value, err := strconv.Atoi(metric)
	assert.Nil(t, err)
	assert.NotZero(t, value)
}

func TestDenyDecisionMetric(t *testing.T) {

	pm := NewMetrics()
	defer pm.UnregisterAll()
	pm.DenyDecisionCount.Inc()
	metric, metricErr := getMetric(pm, "portieris_pod_admission_decision_deny_count")

	assert.Nil(t, metricErr)

	assert.NotEqual(t, "", metric)
	value, err := strconv.Atoi(metric)
	assert.Nil(t, err)
	assert.NotZero(t, value)
}
