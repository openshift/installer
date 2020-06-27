package gcp

import (
	"context"
	"fmt"
	"strings"

	serviceusage "google.golang.org/api/serviceusage/v1beta1"
)

// loadLimis loads the comumer quota metric for a given project and the list of services.
// The limits are defined wrt to the regions. If the limit is global, it's location is treated as `global` else the
// region value is used for location.
func loadLimits(ctx context.Context, client *serviceusage.ServicesService, project string, services ...string) ([]record, error) {
	// services.consumerQuotaMetrics/list requires each service to be of the
	// form projects/{project id}/services/{service name}
	// see https://cloud.google.com/service-usage/docs/reference/rest/v1beta1/services.consumerQuotaMetrics/list
	parent := fmt.Sprintf("projects/%s/services/", project)
	for i := range services {
		if !strings.HasPrefix(services[i], parent) {
			services[i] = fmt.Sprintf("%s%s", parent, services[i])
		}
	}

	var limits []record
	for _, service := range services {
		if err := client.ConsumerQuotaMetrics.
			List(service).
			Context(ctx).
			Pages(ctx, func(page *serviceusage.ListConsumerQuotaMetricsResponse) error {
				for _, qm := range page.Metrics {
					for _, ql := range qm.ConsumerQuotaLimits {
						for _, qlb := range ql.QuotaBuckets {
							limit := record{
								Service:  service[strings.LastIndex(service, "/")+1:],
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
							limits = append(limits, limit)
						}
					}
				}
				return nil
			}); err != nil {
			return nil, err
		}
	}
	return limits, nil
}
