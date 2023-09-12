package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MacroPower/osrs_ge_exporter/internal/log"
	"github.com/MacroPower/osrs_ge_exporter/pkg/client"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "osrs"
	subsystem = "ge"
)

type Exporter struct {
	ItemValue          *prometheus.GaugeVec
	ItemHigh5m         *prometheus.GaugeVec
	ItemLow5m          *prometheus.GaugeVec
	ItemHighVolume5m   *prometheus.GaugeVec
	ItemLowVolume5m    *prometheus.GaugeVec
	ItemHigh1h         *prometheus.GaugeVec
	ItemLow1h          *prometheus.GaugeVec
	ItemHighVolume1h   *prometheus.GaugeVec
	ItemLowVolume1h    *prometheus.GaugeVec
	ItemHighLatest     *prometheus.GaugeVec
	ItemHighLatestTime *prometheus.GaugeVec
	ItemLowLatest      *prometheus.GaugeVec
	ItemLowLatestTime  *prometheus.GaugeVec
	ItemHighAlch       *prometheus.GaugeVec
	ItemLowAlch        *prometheus.GaugeVec
	ItemLimit          *prometheus.GaugeVec

	mu            sync.Mutex
	up            prometheus.Gauge
	totalScrapes  prometheus.Counter
	queryFailures prometheus.Counter

	client  *client.PriceClient
	timeout time.Duration
	logger  log.Logger
}

// NewExporter creates an Exporter.
func NewExporter(client *client.PriceClient, timeout time.Duration, logger log.Logger) *Exporter {
	labels := []string{
		"name",
		"id",
		"members",
		"icon",
	}

	return &Exporter{
		ItemValue: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_value",
				Help:      "Current value of an item.",
			},
			labels,
		),
		ItemHigh5m: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_high_5m",
				Help:      "High value of an item (5m avg).",
			},
			labels,
		),
		ItemLow5m: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_low_5m",
				Help:      "Low value of an item (5m avg).",
			},
			labels,
		),
		ItemHighVolume5m: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_high_volume_5m",
				Help:      "Traded volume of an item (5m).",
			},
			labels,
		),
		ItemLowVolume5m: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_low_volume_5m",
				Help:      "Traded volume of an item (5m).",
			},
			labels,
		),
		ItemHigh1h: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_high_1h",
				Help:      "High value of an item (1h avg).",
			},
			labels,
		),
		ItemLow1h: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_low_1h",
				Help:      "Low value of an item (1h avg).",
			},
			labels,
		),
		ItemHighVolume1h: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_high_volume_1h",
				Help:      "Traded volume of an item (1h).",
			},
			labels,
		),
		ItemLowVolume1h: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_low_volume_1h",
				Help:      "Traded volume of an item (1h).",
			},
			labels,
		),
		ItemHighLatest: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_high_latest",
				Help:      "High value of an item (latest).",
			},
			labels,
		),
		ItemHighLatestTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_high_latest_time",
				Help:      "Unix timestamp of the latest transaction.",
			},
			labels,
		),
		ItemLowLatest: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_low_latest",
				Help:      "Low value of an item (latest).",
			},
			labels,
		),
		ItemLowLatestTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_low_latest_time",
				Help:      "Unix timestamp of the latest transaction.",
			},
			labels,
		),
		ItemHighAlch: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_high_alch",
				Help:      "High alch value of an item.",
			},
			labels,
		),
		ItemLowAlch: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_low_alch",
				Help:      "Low alch value of an item.",
			},
			labels,
		),
		ItemLimit: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "item_limit",
				Help:      "Buy limit for an item.",
			},
			labels,
		),
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "up",
			Help:      "Was the last scrape successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "exporter_scrapes_total",
			Help:      "Number of scrapes.",
		}),
		queryFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "exporter_query_failures_total",
			Help:      "Number of errors.",
		}),
		client:  client,
		timeout: timeout,
		logger:  logger,
	}
}

// Describe describes all metrics with constant descriptions.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()
	ch <- e.totalScrapes.Desc()
	ch <- e.queryFailures.Desc()
}

// Collect sets and collects all metrics.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mu.Lock() // To protect metrics from concurrent collects.
	defer e.mu.Unlock()

	e.ItemValue.Reset()
	e.ItemHigh5m.Reset()
	e.ItemLow5m.Reset()
	e.ItemHighVolume5m.Reset()
	e.ItemLowVolume5m.Reset()
	e.ItemHigh1h.Reset()
	e.ItemLow1h.Reset()
	e.ItemHighVolume1h.Reset()
	e.ItemLowVolume1h.Reset()
	e.ItemHighLatest.Reset()
	e.ItemHighLatestTime.Reset()
	e.ItemLowLatest.Reset()
	e.ItemLowLatestTime.Reset()
	e.ItemHighAlch.Reset()
	e.ItemLowAlch.Reset()
	e.ItemLimit.Reset()

	up := float64(1)
	err := e.scrape()
	if err != nil {
		up = float64(0)
		e.queryFailures.Inc()
		log.Error(e.logger).Log("msg", "Collection failed", "err", err)
	}
	e.up.Set(up)
	e.totalScrapes.Inc()

	e.ItemValue.Collect(ch)
	e.ItemHigh5m.Collect(ch)
	e.ItemLow5m.Collect(ch)
	e.ItemHighVolume5m.Collect(ch)
	e.ItemLowVolume5m.Collect(ch)
	e.ItemHigh1h.Collect(ch)
	e.ItemLow1h.Collect(ch)
	e.ItemHighVolume1h.Collect(ch)
	e.ItemLowVolume1h.Collect(ch)
	e.ItemHighLatest.Collect(ch)
	e.ItemHighLatestTime.Collect(ch)
	e.ItemLowLatest.Collect(ch)
	e.ItemLowLatestTime.Collect(ch)
	e.ItemHighAlch.Collect(ch)
	e.ItemLowAlch.Collect(ch)
	e.ItemLimit.Collect(ch)

	ch <- e.up
	ch <- e.totalScrapes
	ch <- e.queryFailures
}

func (e *Exporter) scrape() error {
	ctx := context.Background()
	mapping, err := e.client.GetMapping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get mapping: %w", err)
	}
	avg5m, err := e.client.Get5m(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get 5m avg: %w", err)
	}
	avg1h, err := e.client.Get1h(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get 1h avg: %w", err)
	}
	latest, err := e.client.GetLatest(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get latest: %w", err)
	}

	for _, item := range mapping {
		labels := []string{
			item.Name,
			fmt.Sprint(item.ID),
			boolToString(item.Members),
			item.Icon,
		}
		e.ItemValue.WithLabelValues(labels...).Set(float64(item.Value))
		if item.Highalch != nil {
			e.ItemHighAlch.WithLabelValues(labels...).Set(float64(*item.Highalch))
		}
		if item.Lowalch != nil {
			e.ItemLowAlch.WithLabelValues(labels...).Set(float64(*item.Lowalch))
		}
		if item.Limit != nil {
			e.ItemLimit.WithLabelValues(labels...).Set(float64(*item.Limit))
		}

		if avgItem, ok := avg5m.Data[fmt.Sprint(item.ID)]; ok {
			if avgItem.AvgHighPrice != nil {
				e.ItemHigh5m.WithLabelValues(labels...).Set(float64(*avgItem.AvgHighPrice))
			}
			if avgItem.AvgLowPrice != nil {
				e.ItemLow5m.WithLabelValues(labels...).Set(float64(*avgItem.AvgLowPrice))
			}
			if avgItem.HighPriceVolume != nil {
				e.ItemHighVolume5m.WithLabelValues(labels...).Set(float64(*avgItem.HighPriceVolume))
			}
			if avgItem.LowPriceVolume != nil {
				e.ItemLowVolume5m.WithLabelValues(labels...).Set(float64(*avgItem.LowPriceVolume))
			}
		}

		if avgItem, ok := avg1h.Data[fmt.Sprint(item.ID)]; ok {
			if avgItem.AvgHighPrice != nil {
				e.ItemHigh1h.WithLabelValues(labels...).Set(float64(*avgItem.AvgHighPrice))
			}
			if avgItem.AvgLowPrice != nil {
				e.ItemLow1h.WithLabelValues(labels...).Set(float64(*avgItem.AvgLowPrice))
			}
			if avgItem.HighPriceVolume != nil {
				e.ItemHighVolume1h.WithLabelValues(labels...).Set(float64(*avgItem.HighPriceVolume))
			}
			if avgItem.LowPriceVolume != nil {
				e.ItemLowVolume1h.WithLabelValues(labels...).Set(float64(*avgItem.LowPriceVolume))
			}
		}

		if latestItem, ok := latest.Data[fmt.Sprint(item.ID)]; ok {
			if latestItem.High != nil {
				e.ItemHighLatest.WithLabelValues(labels...).Set(float64(*latestItem.High))
			}
			if latestItem.Low != nil {
				e.ItemLowLatest.WithLabelValues(labels...).Set(float64(*latestItem.Low))
			}
			if latestItem.HighTime != nil {
				e.ItemHighLatestTime.WithLabelValues(labels...).Set(float64(*latestItem.HighTime))
			}
			if latestItem.LowTime != nil {
				e.ItemLowLatestTime.WithLabelValues(labels...).Set(float64(*latestItem.LowTime))
			}
		}
	}

	return nil
}

func boolToString(b bool) string {
	if b {
		return "true"
	}

	return "false"
}
