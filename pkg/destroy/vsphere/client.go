package vsphere

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/pbm"
	pbmtypes "github.com/vmware/govmomi/pbm/types"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/installer/pkg/asset/installconfig/vsphere"
)

//go:generate mockgen -source=./client.go -destination=mock/vsphereclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	Logout()
	ListFolders(ctx context.Context, tagID string) ([]mo.Folder, error)
	ListVirtualMachines(ctx context.Context, tagID string) ([]mo.VirtualMachine, error)
	StopVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine) error
	DeleteFolder(ctx context.Context, f mo.Folder) error
	DeleteVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine) error
	DeleteStoragePolicy(ctx context.Context, policyName string) error
	DeleteTag(ctx context.Context, id string) error
	DeleteTagCategory(ctx context.Context, id string) error
	DeleteHostZoneObjects(ctx context.Context, infraID string) error
}

// Client makes calls to the Azure API.
type Client struct {
	client     *vim25.Client
	restClient *rest.Client
	cleanup    vsphere.ClientLogout
}

const defaultTimeout = time.Minute * 5

// NewClient initializes a client.
// Logout() must be called when you are done with the client.
func NewClient(vCenter, username, password string) (*Client, error) {
	vim25Client, restClient, cleanup, err := vsphere.CreateVSphereClients(
		context.TODO(),
		vCenter,
		username,
		password)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:     vim25Client,
		restClient: restClient,
		cleanup:    cleanup,
	}, nil
}

// Logout logs out from the clients used.
func (c *Client) Logout() {
	c.cleanup()
}

func isNotFound(err error) bool {
	return err != nil && strings.HasSuffix(err.Error(), http.StatusText(http.StatusNotFound))
}

func (c *Client) getAttachedObjectsOnTag(ctx context.Context, tag, objType string) ([]types.ManagedObjectReference, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tagManager := tags.NewManager(c.restClient)
	attached, err := tagManager.GetAttachedObjectsOnTags(ctx, []string{tag})
	if err != nil && !isNotFound(err) {
		return nil, err
	}

	// Separate the objects attached to the tag based on type
	var objectList []types.ManagedObjectReference
	for _, attachedObject := range attached {
		for _, ref := range attachedObject.ObjectIDs {
			if ref.Reference().Type == objType {
				objectList = append(objectList, ref.Reference())
			}
		}
	}

	return objectList, nil
}

func (c *Client) getVirtualMachineManagedObjects(ctx context.Context, moRef []types.ManagedObjectReference) ([]mo.VirtualMachine, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var virtualMachineMoList []mo.VirtualMachine
	if len(moRef) > 0 {
		pc := property.DefaultCollector(c.client)
		err := pc.Retrieve(ctx, moRef, nil, &virtualMachineMoList)
		if err != nil {
			return nil, err
		}
	}
	return virtualMachineMoList, nil
}

func (c *Client) getFolderManagedObjects(ctx context.Context, moRef []types.ManagedObjectReference) ([]mo.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var folderMoList []mo.Folder
	if len(moRef) > 0 {
		pc := property.DefaultCollector(c.client)
		err := pc.Retrieve(ctx, moRef, nil, &folderMoList)
		if err != nil {
			return nil, err
		}
	}
	return folderMoList, nil
}

// ListFolders returns all ManagedObjects of type "Folder".
func (c *Client) ListFolders(ctx context.Context, tagID string) ([]mo.Folder, error) {
	folderList, err := c.getAttachedObjectsOnTag(ctx, tagID, "Folder")
	if err != nil {
		return nil, err
	}

	return c.getFolderManagedObjects(ctx, folderList)
}

// ListVirtualMachines returns ManagedObjects of type "VirtualMachine".
func (c *Client) ListVirtualMachines(ctx context.Context, tagID string) ([]mo.VirtualMachine, error) {
	virtualMachineList, err := c.getAttachedObjectsOnTag(ctx, tagID, "VirtualMachine")
	if err != nil {
		return nil, err
	}

	return c.getVirtualMachineManagedObjects(ctx, virtualMachineList)
}

func isPoweredOff(vmMO mo.VirtualMachine) bool {
	return vmMO.Summary.Runtime.PowerState == "poweredOff"
}

// StopVirtualMachine stops a VM if it's not already powered off.
func (c *Client) StopVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	if !isPoweredOff(vmMO) {
		vm := object.NewVirtualMachine(c.client, vmMO.Reference())
		task, err := vm.PowerOff(ctx)
		if err == nil {
			err = task.Wait(ctx)
		}
		return err
	}
	return nil
}

// DeleteVirtualMachine deletes a VirtualMachine.
func (c *Client) DeleteVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	vm := object.NewVirtualMachine(c.client, vmMO.Reference())
	task, err := vm.Destroy(ctx)
	if err == nil {
		err = task.Wait(ctx)
	}
	return err
}

// DeleteFolder deletes a Folder.
func (c *Client) DeleteFolder(ctx context.Context, f mo.Folder) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	folder := object.NewFolder(c.client, f.Reference())

	task, err := folder.Destroy(ctx)
	if err == nil {
		err = task.Wait(ctx)
	}
	return err
}

// DeleteStoragePolicy deletes a Storage Policy named `policyName`.
func (c *Client) DeleteStoragePolicy(ctx context.Context, policyName string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	rtype := pbmtypes.PbmProfileResourceType{
		ResourceType: string(pbmtypes.PbmProfileResourceTypeEnumSTORAGE),
	}

	category := pbmtypes.PbmProfileCategoryEnumREQUIREMENT

	pbmClient, err := pbm.NewClient(ctx, c.client)
	if err != nil {
		return err
	}

	ids, err := pbmClient.QueryProfile(ctx, rtype, string(category))
	if err != nil {
		return err
	}

	profiles, err := pbmClient.RetrieveContent(ctx, ids)
	if err != nil {
		return err
	}

	matchingProfileIds := []pbmtypes.PbmProfileId{}
	for _, p := range profiles {
		if p.GetPbmProfile().Name == policyName {
			profileID := p.GetPbmProfile().ProfileId
			matchingProfileIds = append(matchingProfileIds, profileID)
		}
	}
	if len(matchingProfileIds) > 0 {
		_, err = pbmClient.DeleteProfile(ctx, matchingProfileIds)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteTag deletes a Tag named `id`.
func (c *Client) DeleteTag(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tagManager := tags.NewManager(c.restClient)
	tag, err := tagManager.GetTag(ctx, id)
	if isNotFound(err) {
		return nil
	}
	if err == nil {
		err = tagManager.DeleteTag(ctx, tag)
	}
	return err
}

// DeleteTagCategory deletes a Tag Category named `categoryName`.
func (c *Client) DeleteTagCategory(ctx context.Context, categoryName string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tagManager := tags.NewManager(c.restClient)
	ids, err := tagManager.ListCategories(ctx)
	if err != nil {
		return err
	}

	var errs []error
	for _, id := range ids {
		category, err := tagManager.GetCategory(ctx, id)
		if err != nil {
			if !isNotFound(err) {
				errs = append(errs, errors.Wrapf(err, "could not get category %q", id))
			}
			continue
		}
		if category.Name == categoryName {
			if err = tagManager.DeleteCategory(ctx, category); err != nil {
				return err
			}
			return nil
		}
	}

	return utilerrors.NewAggregate(errs)
}

// DeleteHostZoneObjects removes from the vCenter cluster the associated OCP cluster's vm-host group (VirtualMachine)
// and the vm-host affinity rule.
func (c *Client) DeleteHostZoneObjects(ctx context.Context, infraID string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	finder := find.NewFinder(c.client, false)

	datacenters, err := finder.DatacenterList(ctx, "/...")
	if err != nil {
		return err
	}

	for _, dc := range datacenters {
		finder = finder.SetDatacenter(dc)
		clusterObjs, err := finder.ClusterComputeResourceList(ctx, "/...")
		if err != nil {
			return err
		}

		for _, ccr := range clusterObjs {
			clusterConfigSpec := &types.ClusterConfigSpecEx{}

			clusterConfig, err := ccr.Configuration(ctx)
			if err != nil {
				return err
			}

			for _, r := range clusterConfig.Rule {
				if rule, ok := r.(*types.ClusterVmHostRuleInfo); ok {
					if strings.Contains(rule.Name, infraID) {
						clusterConfigSpec.RulesSpec = append(clusterConfigSpec.RulesSpec, types.ClusterRuleSpec{
							ArrayUpdateSpec: types.ArrayUpdateSpec{
								Operation: "remove",
								RemoveKey: rule.GetClusterRuleInfo().Key,
							},
							Info: &rule.ClusterRuleInfo,
						})
					}
				}
			}

			for _, g := range clusterConfig.Group {
				if vmg, ok := g.(*types.ClusterVmGroup); ok {
					if strings.Contains(vmg.Name, infraID) {
						clusterConfigSpec.GroupSpec = append(clusterConfigSpec.GroupSpec, types.ClusterGroupSpec{
							ArrayUpdateSpec: types.ArrayUpdateSpec{
								Operation: "remove",
								RemoveKey: vmg.Name,
							},
							Info: &vmg.ClusterGroupInfo,
						})
					}
				}
			}

			// If the rules or group spec are empty there is no need to modify the cluster
			if len(clusterConfigSpec.RulesSpec) != 0 || len(clusterConfigSpec.GroupSpec) != 0 {
				task, err := ccr.Reconfigure(ctx, clusterConfigSpec, true)
				if err != nil {
					return err
				}

				if err := task.Wait(ctx); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
