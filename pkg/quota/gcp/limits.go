package gcp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/googleapi"
	serviceusage "google.golang.org/api/serviceusage/v1beta1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Keep Cap unset. wait.Backoff stops early when the next interval exceeds a
// cap. With five steps, the largest interval used here is already eight seconds.
var quotaLimitRetryBackoff = wait.Backoff{
	Duration: time.Second,
	Factor:   2,
	Jitter:   0.2,
	Steps:    5,
}

// loadLimits loads the consumer quota metric for a given project and the list of services.
// The limits are defined wrt to the regions. If the limit is global, it's location is treated as `global` else the
// region value is used for location.
func loadLimits(ctx context.Context, client *serviceusage.ServicesService, project string, services ...string) ([]record, error) {
	return loadLimitsWithBackoff(ctx, client, quotaLimitRetryBackoff, project, services...)
}

func loadLimitsWithBackoff(ctx context.Context, client *serviceusage.ServicesService, backoff wait.Backoff, project string, services ...string) ([]record, error) {
	// services.consumerQuotaMetrics/list requires each service to be of the
	// form projects/{project id}/services/{service name}
	// see https://cloud.google.com/service-usage/docs/reference/rest/v1beta1/services.consumerQuotaMetrics/list
	parent := fmt.Sprintf("projects/%s/services/", project)

	var limits []record
	for _, service := range services {
		qualifiedService := service
		if !strings.HasPrefix(qualifiedService, parent) {
			qualifiedService = fmt.Sprintf("%s%s", parent, qualifiedService)
		}

		var serviceLimits []record
		var lastErr error
		err := wait.ExponentialBackoffWithContext(ctx, backoff, func(ctx context.Context) (bool, error) {
			var attemptLimits []record
			err := client.ConsumerQuotaMetrics.
				List(qualifiedService).
				Context(ctx).
				Pages(ctx, func(page *serviceusage.ListConsumerQuotaMetricsResponse) error {
					for _, qm := range page.Metrics {
						for _, ql := range qm.ConsumerQuotaLimits {
							for _, qlb := range ql.QuotaBuckets {
								limit := record{
									Service:  qualifiedService[strings.LastIndex(qualifiedService, "/")+1:],
									Name:     ql.Metric,
									Location: "global",
								}
								region, ok := qlb.Dimensions["region"]
								if ok {
									limit.Location = region
								}
								if qlb.EffectiveLimit < 0 {
									continue
								}
								limit.Value = qlb.EffectiveLimit
								attemptLimits = append(attemptLimits, limit)
							}
						}
					}
					return nil
				})
			if err != nil {
				lastErr = err
				if isRetryableQuotaLimitError(err) {
					return false, nil
				}
				return false, err
			}

			serviceLimits = attemptLimits
			lastErr = nil
			return true, nil
		})
		if err != nil {
			if wait.Interrupted(err) && ctx.Err() == nil && lastErr != nil {
				return nil, lastErr
			}
			return nil, err
		}
		limits = append(limits, serviceLimits...)
	}
	return limits, nil
}

func isRetryableQuotaLimitError(err error) bool {
	var gErr *googleapi.Error
	if !errors.As(err, &gErr) {
		return false
	}
	return gErr.Code == http.StatusTooManyRequests ||
		(gErr.Code >= http.StatusInternalServerError && gErr.Code <= 599)
}
