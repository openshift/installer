package powervs

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/globalsearchv2"
	"k8s.io/utils/ptr"
)

// TagType The different states deleting a job can take.
type TagType int

const (
	// TagTypeVPC is for Virtual Private Cloud types.
	TagTypeVPC TagType = iota

	// TagTypeLoadBalancer is for Load Balancer types.
	TagTypeLoadBalancer

	// TagTypeCloudInstance is for Virtual Machine instance types.
	TagTypeCloudInstance

	// TagTypePublicGateway is for Public Gateway types.
	TagTypePublicGateway

	// TagTypeFloatingIP is for Floating IP types.
	TagTypeFloatingIP

	// TagTypeNetworkACL is for Network Acces Control List types.
	TagTypeNetworkACL

	// TagTypeSubnet is for Subnet types.
	TagTypeSubnet

	// TagTypeSecurityGroup is for Security Group types.
	TagTypeSecurityGroup

	// TagTypeTransitGateway is for Transit Gateway types.
	TagTypeTransitGateway

	// TagTypeServiceInstance is for Service Instance types.
	TagTypeServiceInstance

	// TagTypeCloudObjectStorage is for Cloud Object Storage types.
	TagTypeCloudObjectStorage
)

// listByTag list IBM Cloud resources by matching tag.
func (o *ClusterUninstaller) listByTag(tagType TagType) ([]string, error) {
	var (
		query               string
		authenticator       core.Authenticator
		globalSearchOptions *globalsearchv2.GlobalSearchV2Options
		searchService       *globalsearchv2.GlobalSearchV2
		moreData                  = true
		perPage             int64 = 100
		searchCursor        string
		searchOptions       *globalsearchv2.SearchOptions
		scanResult          *globalsearchv2.ScanResult
		response            *core.DetailedResponse
		crnStruct           crn.CRN
		result              []string
		err                 error
	)

	switch tagType {
	case TagTypeVPC:
		query = fmt.Sprintf("tags:%s AND family:is AND type:vpc", o.ClusterName)
	case TagTypeLoadBalancer:
		query = fmt.Sprintf("tags:%s AND family:is AND type:load-balancer", o.ClusterName)
	case TagTypeCloudInstance:
		query = fmt.Sprintf("tags:%s AND family:is AND type:instance", o.ClusterName)
	case TagTypePublicGateway:
		query = fmt.Sprintf("tags:%s AND family:is AND type:public-gateway", o.ClusterName)
	case TagTypeFloatingIP:
		query = fmt.Sprintf("tags:%s AND family:is AND type:floating-ip", o.ClusterName)
	case TagTypeNetworkACL:
		query = fmt.Sprintf("tags:%s AND family:is AND type:network-acl", o.ClusterName)
	case TagTypeSubnet:
		query = fmt.Sprintf("tags:%s AND family:is AND type:subnet", o.ClusterName)
	case TagTypeSecurityGroup:
		query = fmt.Sprintf("tags:%s AND family:is AND type:security-group", o.ClusterName)
	case TagTypeTransitGateway:
		query = fmt.Sprintf("tags:%s AND family:resource_controller AND type:gateway", o.ClusterName)
	case TagTypeServiceInstance:
		query = fmt.Sprintf("tags:%s AND family:resource_controller AND type:resource-instance AND crn:crn\\:v1\\:bluemix\\:public\\:power-iaas*", o.ClusterName)
	case TagTypeCloudObjectStorage:
		query = fmt.Sprintf("tags:%s AND family:resource_controller AND type:resource-instance AND crn:crn\\:v1\\:bluemix\\:public\\:cloud-object-storage*", o.ClusterName)
	default:
		return nil, fmt.Errorf("listByTag: tagType %d is unknown", tagType)
	}
	o.Logger.Debugf("listByTag: query = %s", query)

	ctx, cancel := contextWithTimeout()
	defer cancel()

	authenticator, err = o.newAuthenticator(o.APIKey)
	if err != nil {
		return nil, fmt.Errorf("listByTag: newAuthenticator: %w", err)
	}

	globalSearchOptions = &globalsearchv2.GlobalSearchV2Options{
		URL:           globalsearchv2.DefaultServiceURL,
		Authenticator: authenticator,
	}

	searchService, err = globalsearchv2.NewGlobalSearchV2(globalSearchOptions)
	if err != nil {
		return nil, fmt.Errorf("listByTag: globalsearchv2.NewGlobalSearchV2: %w", err)
	}

	result = make([]string, 0)

	for moreData {
		searchOptions = &globalsearchv2.SearchOptions{
			Query: &query,
			Limit: ptr.To(perPage),
			// default Fields: []string{"account_id", "name", "type", "family", "crn"},
			// all     Fields: []string{"*"},
		}
		if searchCursor != "" {
			searchOptions.SetSearchCursor(searchCursor)
		}
		o.Logger.Debugf("listByTag: searchOptions = %+v", searchOptions)

		scanResult, response, err = searchService.SearchWithContext(ctx, searchOptions)
		if err != nil {
			return nil, fmt.Errorf("listByTag: searchService.SearchWithContext: err = %w, response = %v", err, response)
		}
		if scanResult.SearchCursor != nil {
			o.Logger.Debugf("listByTag: scanResult = %+v, scanResult.SearchCursor = %+v, len scanResult.Items = %d", scanResult, *scanResult.SearchCursor, len(scanResult.Items))
		} else {
			o.Logger.Debugf("listByTag: scanResult = %+v, scanResult.SearchCursor = nil, len scanResult.Items = %d", scanResult, len(scanResult.Items))
		}

		for _, item := range scanResult.Items {
			crnStruct, err = crn.Parse(*item.CRN)
			if err != nil {
				o.Logger.Debugf("listByTag: crn = %s", *item.CRN)
				return nil, fmt.Errorf("listByTag: could not parse CRN property")
			}
			o.Logger.Debugf("listByTag: crnStruct = %v, crnStruct.Resource = %v", crnStruct, crnStruct.Resource)

			// Append the ID part of the CRN if it exists
			if crnStruct.Resource == "" {
				result = append(result, *item.CRN)
			} else {
				result = append(result, crnStruct.Resource)
			}
		}

		moreData = int64(len(scanResult.Items)) == perPage
		if moreData {
			if scanResult.SearchCursor != nil {
				searchCursor = *scanResult.SearchCursor
			}
		}
	}

	return result, err
}
