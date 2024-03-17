package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 33
`))
	}))
	defer server.Close()

	metrics, err := FetchMetrics(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(metrics))
	assert.Equal(t, "go_goroutines", metrics[0].Name)
	assert.Equal(t, "GAUGE", metrics[0].Type)
	assert.Equal(t, "Number of goroutines that currently exist.", metrics[0].Description)
}

func TestFetchMetricsError(t *testing.T) {
	_, err := FetchMetrics("http://invalid-url")
	assert.Error(t, err)
}

func TestGenerateDocumentation(t *testing.T) {
	metrics := []*MetricDetail{
		{
			Name:        "go_goroutines",
			Type:        "gauge",
			Description: "Number of goroutines that currently exist.",
		},
	}
	expected := `### go_goroutines
**Type**: ` + "`gauge`" + `
**Description**: Number of goroutines that currently exist.
`
	assert.Equal(t, expected, GenerateDocumentation(metrics))
}

func TestGenerateDocumentationEmpty(t *testing.T) {
	metrics := []*MetricDetail{}
	assert.Equal(t, "", GenerateDocumentation(metrics))
}

func TestIntegration(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
# HELP test_metric A test metric for integration testing.
# TYPE test_metric gauge
test_metric{label="value1"} 123
test_metric{label="value2"} 456
# HELP another_metric Another test metric.
# TYPE another_metric counter
another_metric 789
`)
	}))
	defer mockServer.Close()

	metrics, err := FetchMetrics(mockServer.URL)
	assert.NoError(t, err)

	documentation := GenerateDocumentation(metrics)

	expectedStrings := []string{
		"### test_metric",
		"**Type**: `GAUGE`",
		"**Description**: A test metric for integration testing.",
		"**Labels**: `label`",
		"### another_metric",
		"**Type**: `COUNTER`",
		"**Description**: Another test metric.",
	}

	for _, str := range expectedStrings {
		assert.Contains(t, documentation, str)
	}
}

func TestIntegrationHistogram(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
# HELP test_histogram A test histogram metric.
# TYPE test_histogram histogram
test_histogram_bucket{le="0.05"} 4
test_histogram_bucket{le="0.1"} 7
test_histogram_bucket{le="0.2"} 10
test_histogram_bucket{le="+Inf"} 12
test_histogram_sum 1.23
test_histogram_count 12
# HELP test_gauge A test gauge metric.
# TYPE test_gauge gauge
test_gauge 1.23
`)
	}))
	defer mockServer.Close()

	metrics, err := FetchMetrics(mockServer.URL)
	assert.NoError(t, err)

	documentation := GenerateDocumentation(metrics)

	expectedStrings := []string{
		"### test_histogram",
		"**Type**: `HISTOGRAM`",
		"**Description**: A test histogram metric.",
		"**Buckets**: `0.05`, `0.1`, `0.2`, `+Inf`",
		"### test_gauge",
		"**Type**: `GAUGE`",
		"**Description**: A test gauge metric.",
	}

	for _, str := range expectedStrings {
		assert.Contains(t, documentation, str)
	}
}
