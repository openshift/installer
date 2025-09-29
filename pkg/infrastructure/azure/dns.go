package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types/azure"
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

type createDNSEntriesInput struct {
	infra clusterapi.InfraReadyInput
	/*
		extLBFQDNIPv4        string
		extLBFQDNIPv6        string
	*/
	publicIPv4           string
	publicIPv6           string
	resourceGroupName    string
	networkClientFactory *armnetwork.ClientFactory
	opts                 *arm.ClientOptions
}

// Create DNS entries for azure.
func createDNSEntries(ctx context.Context, in *createDNSEntriesInput) error {
	baseDomainResourceGroup := in.infra.InstallConfig.Config.Azure.BaseDomainResourceGroupName
	zone := in.infra.InstallConfig.Config.BaseDomain
	privatezone := in.infra.InstallConfig.Config.ClusterDomain()
	apiExternalName := fmt.Sprintf("api.%s", in.infra.InstallConfig.Config.ObjectMeta.Name)

	if in.infra.InstallConfig.Config.Azure.ResourceGroupName != "" {
		in.resourceGroupName = in.infra.InstallConfig.Config.Azure.ResourceGroupName
	}
	azureTags := make(map[string]*string)
	for k, v := range in.infra.InstallConfig.Config.Azure.UserTags {
		azureTags[k] = ptr.To(v)
	}
	azureTags[fmt.Sprintf("kubernetes.io_cluster.%s", in.infra.InfraID)] = ptr.To("owned")
	lb, err := getLoadBalancer(ctx, &lbInput{
		networkClientFactory: in.networkClientFactory,
		resourceGroupName:    in.resourceGroupName,
		loadBalancerName:     fmt.Sprintf("%s-internal", in.infra.InfraID),
	})
	if err != nil {
		return fmt.Errorf("failed to get Azure internal load balancer: %w", err)
	}
	if len(lb.Properties.FrontendIPConfigurations) == 0 {
		return fmt.Errorf("failed to get Azure cluster LB frontend IPs")
	}

	var ipv4Addresses, ipv6Addresses []*string
	for _, frontendIPConfig := range lb.Properties.FrontendIPConfigurations {
		if *frontendIPConfig.Properties.PrivateIPAddressVersion == armnetwork.IPVersionIPv4 {
			ipv4Addresses = append(ipv4Addresses, frontendIPConfig.Properties.PrivateIPAddress)
			logrus.Debugf("XXX: PrivateIPv4Address=%s", *frontendIPConfig.Properties.PrivateIPAddress)
		} else if *frontendIPConfig.Properties.PrivateIPAddressVersion == armnetwork.IPVersionIPv6 {
			ipv6Addresses = append(ipv6Addresses, frontendIPConfig.Properties.PrivateIPAddress)
			logrus.Debugf("XXX: PrivateIPv6Address=%s", *frontendIPConfig.Properties.PrivateIPAddress)
		}
	}

	// useIPv6 := false
	// for _, network := range in.InstallConfig.Config.Networking.ServiceNetwork {
	// 	if network.IP.To4() == nil {
	// 		useIPv6 = true
	// 	}
	// }

	privateRecords := []recordPrivateList{}
	ttl := int64(300)
	// if useIPv6 {
	// 	recordType = aaaarecord
	// }

	if len(ipv4Addresses) > 0 {
		for _, ipv4Address := range ipv4Addresses {
			privateRecords = append(privateRecords, createPrivateRecordSet("api-int", azureTags, ttl, armdns.RecordTypeA, *ipv4Address, ""))
			privateRecords = append(privateRecords, createPrivateRecordSet("api", azureTags, ttl, armdns.RecordTypeA, *ipv4Address, ""))
		}
	}
	if len(ipv6Addresses) > 0 {
		for _, ipv6Address := range ipv6Addresses {
			privateRecords = append(privateRecords, createPrivateRecordSet("api-int", azureTags, ttl, armdns.RecordTypeAAAA, *ipv6Address, ""))
			privateRecords = append(privateRecords, createPrivateRecordSet("api", azureTags, ttl, armdns.RecordTypeAAAA, *ipv6Address, ""))
		}
	}

	session, err := in.infra.InstallConfig.Azure.Session()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	subscriptionID := session.Credentials.SubscriptionID

	onAzureStack := in.infra.InstallConfig.Azure.CloudName == azure.StackCloud
	if onAzureStack {
		in.opts.APIVersion = stackDNSAPIVersion
	}

	recordSetClient, err := armdns.NewRecordSetsClient(subscriptionID, session.TokenCreds, in.opts)
	if err != nil {
		return fmt.Errorf("failed to create public record client: %w", err)
	}
	privateRecordSetClient, err := armprivatedns.NewRecordSetsClient(subscriptionID, session.TokenCreds, in.opts)
	if err != nil {
		return fmt.Errorf("failed to create private record client: %w", err)
	}

	// Azure Stack only supports "DNS zones"--there is not a private/public zone distinction,
	// so we handle Azure Stack differently and create all records in the single zone.
	if onAzureStack {
		azureCluster := &capz.AzureCluster{}
		key := client.ObjectKey{
			Name:      in.infra.InfraID,
			Namespace: capiutils.Namespace,
		}
		if err := in.infra.Client.Get(ctx, key, azureCluster); err != nil && azureCluster != nil {
			return fmt.Errorf("failed to get Azure cluster: %w", err)
		}

		if len(azureCluster.Spec.NetworkSpec.APIServerLB.FrontendIPs) == 0 {
			return fmt.Errorf("failed to get Azure cluster LB frontend IPs")
		}
		ipIlb := azureCluster.Spec.NetworkSpec.APIServerLB.FrontendIPs[0].PrivateIPAddress
		stackRecords := []recordList{}
		apiInternalName := fmt.Sprintf("api-int.%s", in.infra.InstallConfig.Config.ObjectMeta.Name)
		if in.infra.InstallConfig.Config.PublicAPI() {
			stackRecords = append(stackRecords, createRecordSet(apiExternalName, azureTags, ttl, armdns.RecordTypeA, in.publicIPv4, ""))
		} else {
			stackRecords = append(stackRecords, createRecordSet(apiExternalName, azureTags, ttl, armdns.RecordTypeA, ipIlb, ""))
		}
		stackRecords = append(stackRecords, createRecordSet(apiInternalName, azureTags, ttl, armdns.RecordTypeA, ipIlb, ""))
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
	if in.infra.InstallConfig.Config.PublicAPI() {
		logrus.Debugf("XXX: apiExternalName=%s", apiExternalName)
		/*
			logrus.Debugf("XXX: extLBFQDNIPv4=%s", in.extLBFQDNIPv4)
			logrus.Debugf("XXX: extLBFQDNIPv6=%s", in.extLBFQDNIPv6)
		*/
		// apiExternalNameV6 := fmt.Sprintf("v6-api.%s", infraID)
		// if useIPv6 {
		// 	cnameRecordName = apiExternalNameV6
		// }
		if in.publicIPv4 != "" {
			publicRecords := createRecordSet(apiExternalName, azureTags, ttl, armdns.RecordTypeA, in.publicIPv4, "")
			_, err = recordSetClient.CreateOrUpdate(ctx, baseDomainResourceGroup, zone, publicRecords.Name, publicRecords.RecordType, publicRecords.RecordSet, nil)
			if err != nil {
				return fmt.Errorf("failed to create public IPv4 record set: %w", err)
			}
		}

		if in.publicIPv6 != "" {
			publicRecords := createRecordSet(apiExternalName, azureTags, ttl, armdns.RecordTypeAAAA, in.publicIPv6, "")
			_, err = recordSetClient.CreateOrUpdate(ctx, baseDomainResourceGroup, zone, publicRecords.Name, publicRecords.RecordType, publicRecords.RecordSet, nil)
			if err != nil {
				return fmt.Errorf("failed to create public IPv6 record set: %w", err)
			}
		}
	}

	for _, record := range privateRecords {
		_, err = privateRecordSetClient.CreateOrUpdate(ctx, in.resourceGroupName, privatezone, record.RecordType, record.Name, record.RecordSet, nil)
		if err != nil {
			return fmt.Errorf("failed to create private record set: %w", err)
		}
	}

	return nil
}

func createPrivateRecordSet(lbType string, azureTags map[string]*string, ttl int64, rType armdns.RecordType, ipAddress string, recordName string) (record recordPrivateList) {
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
	case armdns.RecordTypeCNAME:
		record.RecordType = armprivatedns.RecordTypeCNAME
		record.RecordSet.Properties.CnameRecord = &armprivatedns.CnameRecord{
			Cname: &recordName,
		}
	case armdns.RecordTypeA:
		record.RecordType = armprivatedns.RecordTypeA
		record.RecordSet.Properties.ARecords = []*armprivatedns.ARecord{
			{
				IPv4Address: &ipAddress,
			},
		}
	case armdns.RecordTypeAAAA:
		record.RecordType = armprivatedns.RecordTypeAAAA
		record.RecordSet.Properties.AaaaRecords = []*armprivatedns.AaaaRecord{
			{
				IPv6Address: &ipAddress,
			},
		}
	}
	return record
}

func createRecordSet(lbType string, azureTags map[string]*string, ttl int64, rType armdns.RecordType, ipAddress string, recordName string) (record recordList) {
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
	case armdns.RecordTypeCNAME:
		record.RecordType = armdns.RecordTypeCNAME
		record.RecordSet.Properties.CnameRecord = &armdns.CnameRecord{
			Cname: &recordName,
		}
	case armdns.RecordTypeA:
		record.RecordType = armdns.RecordTypeA
		record.RecordSet.Properties.ARecords = []*armdns.ARecord{
			{
				IPv4Address: &ipAddress,
			},
		}
	case armdns.RecordTypeAAAA:
		record.RecordType = armdns.RecordTypeAAAA
		record.RecordSet.Properties.AaaaRecords = []*armdns.AaaaRecord{
			{
				IPv6Address: &ipAddress,
			},
		}
	}
	return record
}
