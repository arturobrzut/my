// processor.go
package wireguardprocessor

import (
	"context"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	"time"
)

// metricsProcessor defines the structure for the custom metrics processor.
type metricsProcessor struct {
	logger  *zap.Logger
	metrics map[string]struct{} // Set of metric names to process, for faster lookup
}

// newMetricsProcessor creates an instance of metricsProcessor, initializing any necessary dependencies.
func newMetricsProcessor(cfg *Config) func(ctx context.Context, metrics pmetric.Metrics) (pmetric.Metrics, error) {
	// Convert metrics slice into a set for efficient lookup
	metricsSet := make(map[string]struct{}, len(cfg.Metrics))
	for _, metric := range cfg.Metrics {
		metricsSet[metric] = struct{}{}
	}

	return func(ctx context.Context, metrics pmetric.Metrics) (pmetric.Metrics, error) {
		return processMetrics(metrics, metricsSet)
	}
}

// processMetrics contains the core processing logic for the specified metrics.
func processMetrics(metrics pmetric.Metrics, metricsSet map[string]struct{}) (pmetric.Metrics, error) {
	// Iterate over ResourceMetrics
	rmSlice := metrics.ResourceMetrics()
	for i := 0; i < rmSlice.Len(); i++ {
		ilmSlice := rmSlice.At(i).ScopeMetrics()
		for j := 0; j < ilmSlice.Len(); j++ {
			metricSlice := ilmSlice.At(j).Metrics()

			for k := 0; k < metricSlice.Len(); k++ {
				metric := metricSlice.At(k)

				// Check if this is the wireguard_latest_handshake_seconds metric
				if metric.Name() == "wireguard_latest_handshake_seconds" {
					// Create a new metric for wireguard_handshake_seconds
					newMetric := metricSlice.AppendEmpty()
					newMetric.SetName("wireguard_handshake_seconds")
					newMetric.SetUnit("seconds")
					newMetric.SetEmptyGauge()

					// Copy data points with calculated values
					dataPoints := metric.Gauge().DataPoints()
					for l := 0; l < dataPoints.Len(); l++ {
						originalDP := dataPoints.At(l)
						newDP := newMetric.Gauge().DataPoints().AppendEmpty()
						originalDP.CopyTo(newDP)

						// Calculate duration as current time - original timestamp
						timestamp := originalDP.Timestamp().AsTime().Unix()
						duration := time.Now().Unix() - timestamp
						newDP.SetDoubleValue(float64(duration))

						// Copy labels
						originalDP.Attributes().CopyTo(newDP.Attributes())
					}
				}
			}
		}
	}

	return metrics, nil
}
