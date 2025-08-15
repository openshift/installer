package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	sqtypes "github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
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
	"ca-west-1",
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

func loadLimits(ctx context.Context, client *servicequotas.Client, services ...string) ([]record, error) {
	records := map[string]record{}
	key := func(q sqtypes.ServiceQuota) string {
		return fmt.Sprintf("%s/%s", aws.ToString(q.ServiceCode), aws.ToString(q.QuotaCode))
	}

	for _, service := range services {
		defaultSQInput := &servicequotas.ListAWSDefaultServiceQuotasInput{ServiceCode: aws.String(service)}
		defaultSQPaginator := servicequotas.NewListAWSDefaultServiceQuotasPaginator(client, defaultSQInput)

		for defaultSQPaginator.HasMorePages() {
			page, err := defaultSQPaginator.NextPage(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to list default service quotas for %s: %w", service, err)
			}

			for _, sq := range page.Quotas {
				records[key(sq)] = record{
					Service: service,
					Name:    aws.ToString(sq.QuotaCode),
					global:  sq.GlobalQuota,
					Value:   int64(aws.ToFloat64(sq.Value)),
				}
			}
		}

		sqInput := &servicequotas.ListServiceQuotasInput{ServiceCode: aws.String(service)}
		sqPaginator := servicequotas.NewListServiceQuotasPaginator(client, sqInput)
		for sqPaginator.HasMorePages() {
			page, err := sqPaginator.NextPage(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to list service quotas for %s: %w", service, err)
			}

			for _, sq := range page.Quotas {
				records[key(sq)] = record{
					Service: service,
					Name:    aws.ToString(sq.QuotaCode),
					global:  sq.GlobalQuota,
					Value:   int64(aws.ToFloat64(sq.Value)),
				}
			}
		}
	}

	var ret []record
	for _, r := range records {
		ret = append(ret, r)
	}
	return ret, nil
}
