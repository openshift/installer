package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
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
func createDNSEntries(ctx context.Context, in clusterapi.InfraReadyInput, extLBFQDN string, resourceGroup string, opts *arm.ClientOptions) error {
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
	// useIPv6 := false
	// for _, network := range in.InstallConfig.Config.Networking.ServiceNetwork {
	// 	if network.IP.To4() == nil {
	// 		useIPv6 = true
	// 	}
	// }

	privateRecords := []recordPrivateList{}
	ttl := int64(300)
	recordType := arecord
	// if useIPv6 {
	// 	recordType = aaaarecord
	// }
	privateRecords = append(privateRecords, createPrivateRecordSet("api-int", azureTags, ttl, recordType, ipIlb, ""))
	privateRecords = append(privateRecords, createPrivateRecordSet("api", azureTags, ttl, recordType, ipIlb, ""))

	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	subscriptionID := session.Credentials.SubscriptionID

	recordSetClient, err := armdns.NewRecordSetsClient(subscriptionID, session.TokenCreds, opts)
	if err != nil {
		return fmt.Errorf("failed to create public record client: %w", err)
	}
	privateRecordSetClient, err := armprivatedns.NewRecordSetsClient(subscriptionID, session.TokenCreds, opts)
	if err != nil {
		return fmt.Errorf("failed to create private record client: %w", err)
	}

	// Create the records for api and api-int in the private zone and api.<clustername> for public zone.
	// CAPI currently creates a record called "apiserver" instead of "api" so creating "api" for the installer in the private zone.
	if in.InstallConfig.Config.PublicAPI() {
		cnameRecordName := apiExternalName
		// apiExternalNameV6 := fmt.Sprintf("v6-api.%s", infraID)
		// if useIPv6 {
		// 	cnameRecordName = apiExternalNameV6
		// }
		publicRecords := createRecordSet(cnameRecordName, azureTags, ttl, cname, "", extLBFQDN)
		_, err = recordSetClient.CreateOrUpdate(ctx, baseDomainResourceGroup, zone, publicRecords.Name, publicRecords.RecordType, publicRecords.RecordSet, nil)
		if err != nil {
			return fmt.Errorf("failed to create public record set: %w", err)
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
