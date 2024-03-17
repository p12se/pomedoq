package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

var revision = "unknown"

// MetricDetail represents a single metric's details.
type MetricDetail struct {
	Name        string   // Name of the metric
	Type        string   // Type of the metric
	Description string   // Description of the metric
	Labels      []string // Labels associated with the metric
	Buckets     []string // Buckets associated with the metric (only for histogram type)
}

// FetchMetrics fetches and parses the metrics from a given URL.
// It returns a slice of MetricDetail pointers and any error encountered.
func FetchMetrics(url string) ([]*MetricDetail, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return nil, err
	}

	var metrics []*MetricDetail
	for name, mf := range metricFamilies {
		md := &MetricDetail{
			Name:        name,
			Type:        mf.GetType().String(),
			Description: mf.GetHelp(),
		}

		if mf.GetType() == io_prometheus_client.MetricType_HISTOGRAM {
			md.Buckets = getBuckets(mf)
		}

		md.Labels = getLabels(mf)

		metrics = append(metrics, md)
	}

	return metrics, nil
}

// getBuckets extracts the bucket boundaries from a histogram metric.
// It returns a slice of strings representing the bucket boundaries.
func getBuckets(mf *io_prometheus_client.MetricFamily) []string {
	var buckets []string
	for _, b := range mf.Metric[0].Histogram.Bucket {
		buckets = append(buckets, fmt.Sprintf("%v", b.GetUpperBound()))
	}
	return buckets
}

// getLabels extracts the labels from a metric.
// It returns a slice of strings representing the labels.
func getLabels(mf *io_prometheus_client.MetricFamily) []string {
	var labels []string
	for _, label := range mf.Metric[0].Label {
		labels = append(labels, *label.Name)
	}
	return labels
}

// GenerateDocumentation generates a markdown formatted documentation string for a given slice of metrics.
// It returns the generated documentation string.
func GenerateDocumentation(metrics []*MetricDetail) string {
	var doc strings.Builder
	for _, metric := range metrics {
		doc.WriteString(fmt.Sprintf("### %s\n**Type**: `%s`\n**Description**: %s\n", metric.Name, metric.Type, metric.Description))
		if len(metric.Labels) > 0 {
			doc.WriteString("**Labels**: `" + strings.Join(metric.Labels, ", ") + "`\n\n")
		}
		if len(metric.Buckets) > 0 {
			doc.WriteString("**Buckets**: `" + strings.Join(metric.Buckets, "`, `") + "`\n\n")
		}
	}
	return doc.String()
}

// main is the entry point of the program.
// It fetches the metrics from a given URL and prints the generated documentation.
func main() {
	fmt.Println("pomidoq: revision: ", revision)

	metrics, err := FetchMetrics("http://localhost:8081/metrics")
	if err != nil {
		fmt.Printf("Error fetching metrics: %v\n", err)
		return
	}
	fmt.Println(GenerateDocumentation(metrics))
}
