package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
)

// SupportedRegions is a list of AWS regions that support the servicequota APIs.
// see https://docs.aws.amazon.com/general/latest/gr/servicequotas.html
var SupportedRegions = sets.NewString(
	"us-east-2",
	"us-east-1",
	"us-west-1",
	"us-west-2",
	"ap-south-1",
	"ap-northeast-3",
	"ap-northeast-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-northeast-1",
	"ca-central-1",
	"eu-central-1",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"sa-east-1",
)

// record stores the data from quota limits and usages.
type record struct {
	Service string
	Name    string
	global  bool

	Value int64
}

func loadLimits(ctx context.Context, client *servicequotas.ServiceQuotas, services ...string) ([]record, error) {
	records := map[string]record{}
	key := func(q *servicequotas.ServiceQuota) string {
		return fmt.Sprintf("%s/%s", aws.StringValue(q.ServiceCode), aws.StringValue(q.QuotaCode))
	}

	for _, service := range services {
		if err := client.ListAWSDefaultServiceQuotasPagesWithContext(ctx,
			&servicequotas.ListAWSDefaultServiceQuotasInput{ServiceCode: aws.String(service)},
			func(page *servicequotas.ListAWSDefaultServiceQuotasOutput, lastPage bool) bool {
				for _, sq := range page.Quotas {
					records[key(sq)] = record{
						Service: service,
						Name:    aws.StringValue(sq.QuotaCode),
						global:  aws.BoolValue(sq.GlobalQuota),
						Value:   int64(aws.Float64Value(sq.Value)),
					}
				}
				return !lastPage
			}); err != nil {
			return nil, errors.Wrapf(err, "failed to list default serviceqquotas for %s", service)

		}

		if err := client.ListServiceQuotasPagesWithContext(ctx,
			&servicequotas.ListServiceQuotasInput{ServiceCode: aws.String(service)},
			func(page *servicequotas.ListServiceQuotasOutput, lastPage bool) bool {
				for _, sq := range page.Quotas {
					records[key(sq)] = record{
						Service: service,
						Name:    aws.StringValue(sq.QuotaCode),
						global:  aws.BoolValue(sq.GlobalQuota),
						Value:   int64(aws.Float64Value(sq.Value)),
					}
				}
				return !lastPage
			}); err != nil {
			return nil, errors.Wrapf(err, "failed to list serviceqquotas for %s", service)
		}

	}
	var ret []record
	for _, r := range records {
		ret = append(ret, r)
	}
	return ret, nil
}
