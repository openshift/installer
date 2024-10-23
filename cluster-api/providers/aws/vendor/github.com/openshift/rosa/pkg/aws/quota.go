package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	servicequotastypes "github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
)

type quota struct {
	ServiceCode  string
	QuotaName    string
	QuotaCode    string
	DesiredValue *float64
}

const IAMServiceCode = "iam"

// List of service quotas we verify for cluster installs
// to support 5 x multi zone clusters
var serviceQuotaServices = []quota{
	{
		ServiceCode:  "ec2",
		QuotaCode:    "L-0263D0A3",
		QuotaName:    "Number of EIPs - VPC EIPs",
		DesiredValue: aws.Float64(5.0),
	},
	{
		ServiceCode:  "ec2",
		QuotaCode:    "L-1216C47A",
		QuotaName:    "Running On-Demand Standard (A, C, D, H, I, M, R, T, Z) instances",
		DesiredValue: aws.Float64(100.0),
	},
	{
		ServiceCode:  "vpc",
		QuotaCode:    "L-F678F1CE",
		QuotaName:    "VPCs per Region",
		DesiredValue: aws.Float64(5.0),
	},
	{
		ServiceCode:  "vpc",
		QuotaCode:    "L-A4707A72",
		QuotaName:    "Internet gateways per Region",
		DesiredValue: aws.Float64(5.0),
	},
	{
		ServiceCode:  "vpc",
		QuotaCode:    "L-DF5E4CA3",
		QuotaName:    "Network interfaces per Region",
		DesiredValue: aws.Float64(5000.0),
	},
	{
		ServiceCode:  "ebs",
		QuotaCode:    "L-D18FCD1D",
		QuotaName:    "General Purpose SSD (gp2) volume storage",
		DesiredValue: aws.Float64(50.0),
	},
	{
		ServiceCode:  "ebs",
		QuotaCode:    "L-309BACF6",
		QuotaName:    "Number of EBS snapshots",
		DesiredValue: aws.Float64(300.0),
	},
	{
		ServiceCode:  "ebs",
		QuotaCode:    "L-B3A130E6",
		QuotaName:    "Provisioned IOPS",
		DesiredValue: aws.Float64(300000.0),
	},
	{
		ServiceCode:  "ebs",
		QuotaCode:    "L-FD252861",
		QuotaName:    "Provisioned IOPS SSD (io1) volume storage",
		DesiredValue: aws.Float64(50.0),
	},
	{
		ServiceCode:  "elasticloadbalancing",
		QuotaCode:    "L-53DA6B97",
		QuotaName:    "Application Load Balancers per Region",
		DesiredValue: aws.Float64(50.0),
	},
	{
		ServiceCode:  "elasticloadbalancing",
		QuotaCode:    "L-E9E9831D",
		QuotaName:    "Classic Load Balancers per Region",
		DesiredValue: aws.Float64(20.0),
	},
}

// ValidateQuota
func (c *awsClient) ValidateQuota() (bool, error) {
	var invalidQuotas []string
	for _, quota := range serviceQuotaServices {
		serviceQuotas, err := ListServiceQuotas(c, quota.ServiceCode)
		if err != nil {
			return false, fmt.Errorf("Error listing AWS service quotas: %s %v", quota.ServiceCode, err)
		}

		serviceQuota, err := GetServiceQuota(serviceQuotas, quota.QuotaCode)
		if err != nil || (serviceQuota).Value == nil {
			return false, fmt.Errorf("Error getting AWS service quota: %s %v", quota.ServiceCode, err)
		}

		if *serviceQuota.Value < *quota.DesiredValue {
			invalidQuotas = append(invalidQuotas, fmt.Sprintf(
				"- Service %s quota code %s %s not valid, expected quota of at least %d, but got %d",
				quota.ServiceCode, quota.QuotaCode, quota.QuotaName,
				int(*quota.DesiredValue), int(*serviceQuota.Value)))
		}

		c.logger.Debug(fmt.Sprintf("Service %s quota code %s is ok", quota.ServiceCode, quota.QuotaCode))
	}

	if len(invalidQuotas) > 0 {
		return false, fmt.Errorf("Service quota is insufficient for the following service quota codes:\n%s",
			strings.Join(invalidQuotas, "\n"))
	}

	return true, nil
}

// ListServiceQuotas list available quotas for service
func ListServiceQuotas(client *awsClient, serviceCode string) ([]servicequotastypes.ServiceQuota, error) {
	// Create a ServiceQuotas client
	var serviceQuotas []servicequotastypes.ServiceQuota

	// Paginate through quota results
	paginator := servicequotas.NewListServiceQuotasPaginator(
		client.serviceQuotasClient, &servicequotas.ListServiceQuotasInput{
			ServiceCode: &serviceCode,
		})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		serviceQuotas = append(serviceQuotas, page.Quotas...)
	}

	return serviceQuotas, nil
}

// GetServiceQuota extract service quota for the list of service quotas
func GetServiceQuota(serviceQuotas []servicequotastypes.ServiceQuota,
	quotaCode string) (servicequotastypes.ServiceQuota, error) {
	for _, serviceQuota := range serviceQuotas {
		if *serviceQuota.QuotaCode == quotaCode {
			return serviceQuota, nil
		}
	}
	return servicequotastypes.ServiceQuota{}, fmt.Errorf("Unable to find quota with service code: %s", quotaCode)
}

func (c *awsClient) GetIAMServiceQuota(quotaCode string) (
	*servicequotas.GetServiceQuotaOutput, error) {
	return c.iamQuotaClient.GetServiceQuota(context.Background(), &servicequotas.GetServiceQuotaInput{
		ServiceCode: aws.String(IAMServiceCode),
		QuotaCode:   aws.String(quotaCode),
	})
}
