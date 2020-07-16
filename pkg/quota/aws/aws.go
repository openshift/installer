package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicequotas"

	"github.com/openshift/installer/pkg/quota"
	"github.com/pkg/errors"
)

// Load load the quota information for a region. It provides information
// about the usage and limit for each resource quota.
func Load(ctx context.Context, sess *session.Session, region string, services ...string) ([]quota.Quota, error) {
	client := servicequotas.New(sess, aws.NewConfig().WithRegion(region))
	records, err := loadLimits(ctx, client, services...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load limits for servicequotas")
	}
	return newQuota(region, records), nil
}

func newQuota(region string, limits []record) []quota.Quota {
	var ret []quota.Quota
	for _, limit := range limits {
		q := quota.Quota{
			Service: limit.Service,
			Name:    fmt.Sprintf("%s/%s", limit.Service, limit.Name),
			Region:  region,
			InUse:   0,
			Limit:   limit.Value,
		}
		if limit.global {
			q.Region = "global"
		}
		ret = append(ret, q)
	}
	return ret
}

// IsUnauthorized checks if the error is un authorized.
func IsUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	var awsErr awserr.Error
	if errors.As(err, &awsErr) {
		// see reference:
		// https://docs.aws.amazon.com/servicequotas/2019-06-24/apireference/API_GetServiceQuota.html
		// https://docs.aws.amazon.com/servicequotas/2019-06-24/apireference/API_GetAWSDefaultServiceQuota.html
		return awsErr.Code() == "AccessDeniedException" || awsErr.Code() == "NoSuchResourceException"
	}
	return false
}
