package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types/azure"
)

type recordListType string

const (
	cname      recordListType = "Cname"
	arecord    recordListType = "ARecord"
	aaaarecord recordListType = "AaaaRecord"
)

type recordList struct {
	Name       string
	RecordType armdns.RecordType
	RecordSet  armdns.RecordSet
}

type recordPrivateList struct {
	Name       string
	RecordType armprivatedns.RecordType
	RecordSet  armprivatedns.RecordSet
}

// Create DNS entries for azure.
func createDNSEntries(ctx context.Context, in clusterapi.InfraReadyInput, extLBFQDN, publicIP, intIPv6IP, resourceGroup string, opts *arm.ClientOptions) error {
	baseDomainResourceGroup := in.InstallConfig.Config.Azure.BaseDomainResourceGroupName
	zone := in.InstallConfig.Config.BaseDomain
	privatezone := in.InstallConfig.Config.ClusterDomain()
	apiExternalName := fmt.Sprintf("api.%s", in.InstallConfig.Config.ObjectMeta.Name)

	if in.InstallConfig.Config.Azure.ResourceGroupName != "" {
		resourceGroup = in.InstallConfig.Config.Azure.ResourceGroupName
	}
	azureTags := make(map[string]*string)
	for k, v := range in.InstallConfig.Config.Azure.UserTags {
		azureTags[k] = ptr.To(v)
	}
	azureTags[fmt.Sprintf("kubernetes.io_cluster.%s", in.InfraID)] = ptr.To("owned")
	azureCluster := &capz.AzureCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, azureCluster); err != nil && azureCluster != nil {
		return fmt.Errorf("failed to get Azure cluster: %w", err)
	}

	if len(azureCluster.Spec.NetworkSpec.APIServerLB.FrontendIPs) == 0 {
		return fmt.Errorf("failed to get Azure cluster LB frontend IPs")
	}
	ipIlb := azureCluster.Spec.NetworkSpec.APIServerLB.FrontendIPs[0].PrivateIPAddress

	privateRecords := []recordPrivateList{}
	ttl := int64(300)
	recordType := arecord
	privateRecords = append(privateRecords, createPrivateRecordSet("api-int", azureTags, ttl, recordType, ipIlb, ""))
	privateRecords = append(privateRecords, createPrivateRecordSet("api", azureTags, ttl, recordType, ipIlb, ""))

	if in.InstallConfig.Config.Azure.IPFamily.DualStackEnabled() {
		// Get the internal LB's IPv6 address from Azure (dynamically assigned)
		ipIlbV6, err := getInternalLBIPv6(ctx, in, resourceGroup, opts)
		if err == nil && ipIlbV6 != "" {
			privateRecords = append(privateRecords, createPrivateRecordSet("api-int", azureTags, ttl, aaaarecord, ipIlbV6, ""))
			privateRecords = append(privateRecords, createPrivateRecordSet("api", azureTags, ttl, aaaarecord, ipIlbV6, ""))
		}
	}
	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	subscriptionID := session.Credentials.SubscriptionID

	onAzureStack := in.InstallConfig.Azure.CloudName == azure.StackCloud
	if onAzureStack {
		opts.APIVersion = stackDNSAPIVersion
	}

	recordSetClient, err := armdns.NewRecordSetsClient(subscriptionID, session.TokenCreds, opts)
	if err != nil {
		return fmt.Errorf("failed to create public record client: %w", err)
	}
	privateRecordSetClient, err := armprivatedns.NewRecordSetsClient(subscriptionID, session.TokenCreds, opts)
	if err != nil {
		return fmt.Errorf("failed to create private record client: %w", err)
	}

	// Azure Stack only supports "DNS zones"--there is not a private/public zone distinction,
	// so we handle Azure Stack differently and create all records in the single zone.
	if onAzureStack {
		stackRecords := []recordList{}
		apiInternalName := fmt.Sprintf("api-int.%s", in.InstallConfig.Config.ObjectMeta.Name)
		if in.InstallConfig.Config.PublicAPI() {
			stackRecords = append(stackRecords, createRecordSet(apiExternalName, azureTags, ttl, arecord, publicIP, ""))
		} else {
			stackRecords = append(stackRecords, createRecordSet(apiExternalName, azureTags, ttl, arecord, ipIlb, ""))
		}
		stackRecords = append(stackRecords, createRecordSet(apiInternalName, azureTags, ttl, arecord, ipIlb, ""))
		for _, record := range stackRecords {
			_, err = recordSetClient.CreateOrUpdate(ctx, baseDomainResourceGroup, zone, record.Name, record.RecordType, record.RecordSet, nil)
			if err != nil {
				return fmt.Errorf("failed to create public record set: %w", err)
			}
		}
		return nil
	}

	// Create the records for api and api-int in the private zone and api.<clustername> for public zone.
	// CAPI currently creates a record called "apiserver" instead of "api" so creating "api" for the installer in the private zone.
	if in.InstallConfig.Config.PublicAPI() {
		cnameRecordName := apiExternalName
		// if useIPv6 {
		// 	cnameRecordName = apiExternalNameV6
		// }
		publicRecords := createRecordSet(cnameRecordName, azureTags, ttl, cname, "", extLBFQDN)
		_, err = recordSetClient.CreateOrUpdate(ctx, baseDomainResourceGroup, zone, publicRecords.Name, publicRecords.RecordType, publicRecords.RecordSet, nil)
		if err != nil {
			return fmt.Errorf("failed to create public record set: %w", err)
		}
		if in.InstallConfig.Config.Azure.IPFamily.DualStackEnabled() {
			apiExternalNameV6 := fmt.Sprintf("api.%s", in.InstallConfig.Config.ObjectMeta.Name)
			publicRecords := createRecordSet(apiExternalNameV6, azureTags, ttl, cname, "", extLBFQDN)
			_, err = recordSetClient.CreateOrUpdate(ctx, baseDomainResourceGroup, zone, publicRecords.Name, publicRecords.RecordType, publicRecords.RecordSet, nil)
			if err != nil {
				return fmt.Errorf("failed to create public record set: %w", err)
			}
		}
	}

	for _, record := range privateRecords {
		_, err = privateRecordSetClient.CreateOrUpdate(ctx, resourceGroup, privatezone, record.RecordType, record.Name, record.RecordSet, nil)
		if err != nil {
			return fmt.Errorf("failed to create private record set: %w", err)
		}
	}

	return nil
}

func createPrivateRecordSet(lbType string, azureTags map[string]*string, ttl int64, rType recordListType, ipAddress string, recordName string) (record recordPrivateList) {
	record = recordPrivateList{
		Name: lbType,
		RecordSet: armprivatedns.RecordSet{
			Properties: &armprivatedns.RecordSetProperties{
				TTL:      &ttl,
				Metadata: azureTags,
			},
		},
	}

	switch rType {
	case cname:
		record.RecordType = armprivatedns.RecordTypeCNAME
		record.RecordSet.Properties.CnameRecord = &armprivatedns.CnameRecord{
			Cname: &recordName,
		}
	case arecord:
		record.RecordType = armprivatedns.RecordTypeA
		record.RecordSet.Properties.ARecords = []*armprivatedns.ARecord{
			{
				IPv4Address: &ipAddress,
			},
		}
	case aaaarecord:
		record.RecordType = armprivatedns.RecordTypeAAAA
		record.RecordSet.Properties.AaaaRecords = []*armprivatedns.AaaaRecord{
			{
				IPv6Address: &ipAddress,
			},
		}
	}
	return record
}

func createRecordSet(lbType string, azureTags map[string]*string, ttl int64, rType recordListType, ipAddress string, recordName string) (record recordList) {
	record = recordList{
		Name: lbType,
		RecordSet: armdns.RecordSet{
			Properties: &armdns.RecordSetProperties{
				TTL:      &ttl,
				Metadata: azureTags,
			},
		},
	}

	switch rType {
	case cname:
		record.RecordType = armdns.RecordTypeCNAME
		record.RecordSet.Properties.CnameRecord = &armdns.CnameRecord{
			Cname: &recordName,
		}
	case arecord:
		record.RecordType = armdns.RecordTypeA
		record.RecordSet.Properties.ARecords = []*armdns.ARecord{
			{
				IPv4Address: &ipAddress,
			},
		}
	case aaaarecord:
		record.RecordType = armdns.RecordTypeAAAA
		record.RecordSet.Properties.AaaaRecords = []*armdns.AaaaRecord{
			{
				IPv6Address: &ipAddress,
			},
		}
	}
	return record
}

// getInternalLBIPv6 retrieves the IPv6 private IP address of the internal load balancer.
func getInternalLBIPv6(ctx context.Context, in clusterapi.InfraReadyInput, resourceGroup string, opts *arm.ClientOptions) (string, error) {
	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return "", fmt.Errorf("failed to create azure session: %w", err)
	}

	subscriptionID := session.Credentials.SubscriptionID
	tokenCredential := session.TokenCreds

	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionID, tokenCredential, opts)
	if err != nil {
		return "", fmt.Errorf("failed to create azure network factory: %w", err)
	}

	lbClient := networkClientFactory.NewLoadBalancersClient()
	internalLBName := fmt.Sprintf("%s-internal", in.InfraID)

	lb, err := lbClient.Get(ctx, resourceGroup, internalLBName, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get internal load balancer: %w", err)
	}

	if lb.Properties == nil || lb.Properties.FrontendIPConfigurations == nil {
		return "", fmt.Errorf("internal load balancer has no frontend IP configurations")
	}

	// Find the IPv6 frontend IP configuration (should be the second one for dual-stack)
	for _, frontendIP := range lb.Properties.FrontendIPConfigurations {
		if frontendIP.Properties != nil &&
			frontendIP.Properties.PrivateIPAddress != nil &&
			frontendIP.Properties.PrivateIPAddressVersion != nil &&
			*frontendIP.Properties.PrivateIPAddressVersion == armnetwork.IPVersionIPv6 {
			return *frontendIP.Properties.PrivateIPAddress, nil
		}
	}

	return "", fmt.Errorf("no IPv6 frontend IP found on internal load balancer")
}
