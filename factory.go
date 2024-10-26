// factory.go
package my

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/consumer"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

// Config defines the configuration for the custom processor, including metrics to convert.
type Config struct {
	Metrics []string `mapstructure:"metrics"` // List of delta sum metrics to convert to rates
}

// Validate checks if the configuration is valid.
func (config *Config) Validate() error {
	if len(config.Metrics) == 0 {
		return fmt.Errorf("metric names are missing")
	}
	return nil
}

// CreateDefaultConfig returns a default configuration instance for the processor.
func CreateDefaultConfig() component.Config {
	return &Config{
		Metrics: []string{}, // Default to an empty list
	}
}

var processorCapabilities = consumer.Capabilities{MutatesData: true}

// CreateMetricsProcessor creates the processor and ensures configuration is valid.
func CreateMetricsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (processor.Metrics, error) {

	// Type assert and validate configuration
	pConfig := cfg.(*Config)
	if err := pConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config for exampleprocessor: %w", err)
	}
	return processorhelper.NewMetrics(
		ctx,
		set,
		cfg,
		nextConsumer,
		newMetricsProcessor(pConfig),
		processorhelper.WithCapabilities(processorCapabilities))
}

func NewProcessorFactory() processor.Factory {
	typeStr, _ := component.NewType("wireguardprocessor")
	return processor.NewFactory(
		typeStr,
		func() component.Config { return &struct{}{} },
		processor.WithMetrics(CreateMetricsProcessor, component.StabilityLevelBeta),
	)
}
