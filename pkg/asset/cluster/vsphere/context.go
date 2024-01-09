package vsphere

import (
	"context"
	"fmt"
	"path"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types/vsphere"
	vsphereplatform "github.com/openshift/installer/pkg/types/vsphere"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vapi/tags"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

func (a *VCenterContexts) cacheNetworkNames(ctx context.Context, failureDomain vsphere.FailureDomain, session *session.Session, server string) error {

	finder := session.Finder
	clusterPath := failureDomain.Topology.ComputeCluster

	clusterRef, err := finder.ClusterComputeResource(ctx, clusterPath)
	if err != nil {
		return fmt.Errorf("unable to retrieve compute cluster: %v", err)
	}

	pools, err := finder.ResourcePoolList(ctx, clusterRef.InventoryPath)
	if err != nil {
		return fmt.Errorf("unable to retrieve resource pools relative to compute cluster: %v", err)
	}

	for _, network := range failureDomain.Topology.Networks {
		clusterMap, present := a.VCenters[server].ClusterNetworkMap[clusterPath]
		if !present {
			clusterMap = NetworkNameMap{
				Cluster:       clusterPath,
				NetworkNames:  map[string]string{},
				ResourcePools: map[string]*object.ResourcePool{},
			}
			for _, pool := range pools {
				clusterMap.ResourcePools[path.Clean(pool.InventoryPath)] = pool
			}

			a.VCenters[server].ClusterNetworkMap[clusterPath] = clusterMap
		}

		networkName := path.Join(clusterRef.InventoryPath, network)
		clusterMap.NetworkNames[network] = networkName

	}

	return nil
}

func (a *VCenterContexts) createClusterTagID(ctx context.Context, session *session.Session, clusterId string, server string) error {
	tagManager := session.TagManager
	categories, err := tagManager.GetCategories(ctx)
	if err != nil {
		return fmt.Errorf("unable to get tag categories: %v", err)
	}

	var clusterTagCategory *tags.Category
	clusterTagCategoryName := fmt.Sprintf("openshift-%s", clusterId)
	tagCategoryId := ""

	for _, category := range categories {
		if category.Name == clusterTagCategoryName {
			clusterTagCategory = &category
			tagCategoryId = category.ID
			break
		}
	}

	if clusterTagCategory == nil {
		clusterTagCategory = &tags.Category{
			Name:        clusterTagCategoryName,
			Description: "Added by openshift-install do not remove",
			Cardinality: "SINGLE",
			AssociableTypes: []string{
				"urn:vim25:VirtualMachine",
				"urn:vim25:ResourcePool",
				"urn:vim25:Folder",
				"urn:vim25:Datastore",
				"urn:vim25:StoragePod",
			},
		}
		tagCategoryId, err = tagManager.CreateCategory(ctx, clusterTagCategory)
		if err != nil {
			return fmt.Errorf("unable to create tag category: %v", err)
		}
	}

	var categoryTag *tags.Tag
	tagId := ""

	categoryTags, err := tagManager.GetTagsForCategory(ctx, tagCategoryId)
	if err != nil {
		return fmt.Errorf("unable to get tags for category: %v", err)
	}
	for _, tag := range categoryTags {
		if tag.Name == clusterId {
			categoryTag = &tag
			tagId = tag.ID
			break
		}
	}

	if categoryTag == nil {
		categoryTag = &tags.Tag{
			Description: "Added by openshift-install do not remove",
			Name:        clusterId,
			CategoryID:  tagCategoryId,
		}
		tagId, err = tagManager.CreateTag(ctx, categoryTag)
		if err != nil {
			return fmt.Errorf("unable to create tag: %v", err)
		}
	}

	vCenterContext := a.VCenters[server]
	vCenterContext.TagID = tagId
	a.VCenters[server] = vCenterContext

	return nil
}

type NetworkNameMap struct {
	Cluster       string
	ResourcePools map[string]*object.ResourcePool
	NetworkNames  map[string]string
}

// VCenterContext maintains context of known vCenters to be used in CAPI manifest reconciliation.
type VCenterContext struct {
	VCenter           string
	TagID             string
	Datacenters       []string
	ClusterNetworkMap map[string]NetworkNameMap
}

type VCenterContexts struct {
	VCenters map[string]VCenterContext
}

var (
	_ asset.Asset = (*VCenterContexts)(nil)
)

func (a *VCenterContexts) Generate(parents asset.Parents) error {
	ctx := context.TODO()

	a.VCenters = map[string]VCenterContext{}

	ic := &installconfig.InstallConfig{}
	clusterID := &installconfig.ClusterID{}
	parents.Get(
		ic,
		clusterID,
	)

	if ic.Config.Platform.Name() != vsphereplatform.Name {
		return nil
	}

	installConfig := ic.Config

	for _, vcenter := range installConfig.VSphere.VCenters {
		server := vcenter.Server
		params := session.NewParams().WithServer(server).WithUserInfo(vcenter.Username, vcenter.Password)
		tempConnection, err := session.GetOrCreate(ctx, params)
		if err != nil {
			return fmt.Errorf("unable to create session: %v", err)
		}

		defer tempConnection.CloseIdleConnections()

		a.VCenters[server] = VCenterContext{
			VCenter:           server,
			Datacenters:       vcenter.Datacenters,
			ClusterNetworkMap: map[string]NetworkNameMap{},
		}

		if err = a.createClusterTagID(ctx, tempConnection, clusterID.InfraID, server); err != nil {
			return fmt.Errorf("unable to create cluster tag ID: %v", err)
		}

		for _, failureDomain := range installConfig.VSphere.FailureDomains {
			if failureDomain.Server != server {
				continue
			}
			if err = a.cacheNetworkNames(ctx, failureDomain, tempConnection, server); err != nil {
				return fmt.Errorf("unable to retrieve network names: %v", err)
			}
		}
	}
	return nil
}

func (a *VCenterContexts) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&installconfig.ClusterID{},
	}
}

func (a *VCenterContexts) Name() string {
	return "vCenter Context"
}
