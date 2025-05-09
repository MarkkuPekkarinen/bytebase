package collector

import (
	"context"

	metricapi "github.com/bytebase/bytebase/backend/metric"
	"github.com/bytebase/bytebase/backend/plugin/metric"
	"github.com/bytebase/bytebase/backend/store"
)

var _ metric.Collector = (*issueCountCollector)(nil)

// issueCountCollector is the metric data collector for issue.
type issueCountCollector struct {
	store *store.Store
}

// NewIssueCountCollector creates a new instance of issueCollector.
func NewIssueCountCollector(store *store.Store) metric.Collector {
	return &issueCountCollector{
		store: store,
	}
}

// Collect will collect the metric for issue.
func (c *issueCountCollector) Collect(ctx context.Context) ([]*metric.Metric, error) {
	count, err := c.store.CountIssues(ctx)
	if err != nil {
		return nil, err
	}
	return []*metric.Metric{
		{
			Name:  metricapi.IssueCountMetricName,
			Value: count,
		},
	}, nil
}
