package gcp

import (
	"context"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/api/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// loadUsage loads the usage from `consumer_quota` metric type for the project. It pulls metric data for last hour and
// and the uses the latest data point as the usage value for the specific metric type.
func loadUsage(ctx context.Context, client *monitoring.MetricClient, project string) ([]record, error) {
	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   fmt.Sprintf("projects/%s", project),
		Filter: `metric.type = "serviceruntime.googleapis.com/quota/allocation/usage" AND resource.type = "consumer_quota"`,
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: time.Now().Add(-5 * time.Minute).Unix(),
			},
			StartTime: &googlepb.Timestamp{
				Seconds: time.Now().Add(-1 * time.Hour).Unix(),
			},
		},
	}

	var usages []record
	it := client.ListTimeSeries(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return usages, errors.Wrap(err, "failed to list quota/allocation/usage timeseries")
		}
		usage, err := latestRecord(resp)
		if err != nil {
			return usages, errors.Wrap(err, "failed to load usage from timeseries")
		}
		usages = append(usages, usage)
	}
	return usages, nil
}

// latestRecord find the latest data point for the timeseries ans retuns that as the usage for the metric type.
// Based on https://cloud.google.com/monitoring/api/ref_v3/rest/v3/TimeSeries and definition of the points API
// "The data points of this time series. When listing time series, points are returned in reverse time order.",
// The latestRecord returns the first element of the points as the usage value. In case the points list is empty it
/// returns 0 as the usage value.
func latestRecord(ts *monitoringpb.TimeSeries) (record, error) {
	service, ok := ts.GetResource().GetLabels()["service"]
	if !ok {
		return record{}, errors.Errorf("serivce not found for timeseries, actual label %s", ts.GetResource().GetLabels())
	}
	name, ok := ts.GetMetric().Labels["quota_metric"]
	if !ok {
		return record{}, errors.New("no name found for timeseries")
	}
	location, ok := ts.GetResource().GetLabels()["location"]
	if !ok {
		return record{}, errors.Errorf("location not found for timeseries, actual label %s", ts.GetResource().GetLabels())
	}
	if ts.GetValueType() != metric.MetricDescriptor_INT64 {
		return record{}, errors.Errorf("invalid value type for timeseries, was %s", ts.GetValueType())
	}
	value := int64(0)
	if points := ts.GetPoints(); len(points) > 0 {
		value = points[0].GetValue().GetInt64Value()
	}
	return record{
		Service:  service,
		Name:     name,
		Location: location,
		Value:    value,
	}, nil
}
