package monitoring

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"abb_ia/internal/logger"
)

// MetricType defines the type of metric being collected
type MetricType string

const (
	CounterMetric MetricType = "counter"
	GaugeMetric   MetricType = "gauge"
	TimerMetric   MetricType = "timer"
)

// Metric represents a single metric with its value and metadata
type Metric struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Type        MetricType  `json:"type"`
	Labels      Labels      `json:"labels"`
	LastUpdated time.Time   `json:"last_updated"`
}

// Labels are key-value pairs that can be attached to metrics
type Labels map[string]string

// MetricsCollector manages application metrics
type MetricsCollector struct {
	mu      sync.RWMutex
	metrics map[string]*Metric
}

// Global metrics collector instance and configuration
var (
	globalCollector *MetricsCollector
	once           sync.Once
	metricsEnabled bool = false // Metrics disabled by default
)

// DisableMetrics disables metrics collection
func DisableMetrics() {
	metricsEnabled = false
	if globalCollector != nil {
		globalCollector.ResetMetrics()
	}
}

// EnableMetrics enables metrics collection
func EnableMetrics() {
	metricsEnabled = true
}

// GetMetricsCollector returns the global metrics collector instance
func GetMetricsCollector() *MetricsCollector {
	if !metricsEnabled {
		return &MetricsCollector{
			metrics: make(map[string]*Metric), // Initialize map even when disabled
		}
	}

	once.Do(func() {
		globalCollector = &MetricsCollector{
			metrics: make(map[string]*Metric),
		}
	})
	return globalCollector
}

// IncrementCounter increments a counter metric
func (mc *MetricsCollector) IncrementCounter(name string, labels Labels) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	metric, exists := mc.metrics[name]
	if !exists {
		metric = &Metric{
			Name:   name,
			Value:  int64(0),
			Type:   CounterMetric,
			Labels: labels,
		}
		mc.metrics[name] = metric
	}

	if counter, ok := metric.Value.(int64); ok {
		metric.Value = counter + 1
		metric.LastUpdated = time.Now()
	}
}

// SetGauge sets a gauge metric value
func (mc *MetricsCollector) SetGauge(name string, value float64, labels Labels) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.metrics[name] = &Metric{
		Name:        name,
		Value:       value,
		Type:        GaugeMetric,
		Labels:      labels,
		LastUpdated: time.Now(),
	}
}

// RecordTimer records a duration for a timer metric
func (mc *MetricsCollector) RecordTimer(name string, duration time.Duration, labels Labels) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	metric, exists := mc.metrics[name]
	if !exists {
		metric = &Metric{
			Name:   name,
			Value:  []time.Duration{},
			Type:   TimerMetric,
			Labels: labels,
		}
		mc.metrics[name] = metric
	}

	if durations, ok := metric.Value.([]time.Duration); ok {
		metric.Value = append(durations, duration)
		metric.LastUpdated = time.Now()
	}
}

// GetMetric returns a specific metric by name
func (mc *MetricsCollector) GetMetric(name string) (*Metric, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metric, exists := mc.metrics[name]
	if !exists {
		return nil, fmt.Errorf("metric not found: %s", name)
	}
	return metric, nil
}

// GetAllMetrics returns all metrics
func (mc *MetricsCollector) GetAllMetrics() map[string]*Metric {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metrics := make(map[string]*Metric)
	for k, v := range mc.metrics {
		metrics[k] = v
	}
	return metrics
}

// ExportMetrics exports all metrics as JSON
func (mc *MetricsCollector) ExportMetrics() ([]byte, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	return json.MarshalIndent(mc.metrics, "", "  ")
}

// StartMetricsReporter starts a goroutine that periodically reports metrics
func (mc *MetricsCollector) StartMetricsReporter(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			metrics := mc.GetAllMetrics()
			for name, metric := range metrics {
				// Format value based on metric type
				var valueStr string
				switch metric.Type {
				case TimerMetric:
					if d, ok := metric.Value.(time.Duration); ok {
						valueStr = d.String()
					}
				case GaugeMetric:
					if f, ok := metric.Value.(float64); ok {
						valueStr = fmt.Sprintf("%.2f", f)
					}
				default:
					valueStr = fmt.Sprintf("%v", metric.Value)
				}

				// Format labels
				labelsStr := ""
				if len(metric.Labels) > 0 {
					labels := make([]string, 0, len(metric.Labels))
					for k, v := range metric.Labels {
						labels = append(labels, fmt.Sprintf("%s=%s", k, v))
					}
					labelsStr = fmt.Sprintf(" [%s]", strings.Join(labels, ", "))
				}

				logger.Debug(fmt.Sprintf("METRIC %s%s = %s (%s) @ %s",
					name,
					labelsStr,
					valueStr,
					metric.Type,
					metric.LastUpdated.Format("15:04:05")))
			}
		}
	}()
}

// ResetMetrics resets all metrics to their zero values
func (mc *MetricsCollector) ResetMetrics() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.metrics = make(map[string]*Metric)
}
