package spbm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/pbm"
	"github.com/vmware/govmomi/pbm/methods"
	pbmtypes "github.com/vmware/govmomi/pbm/types"
	"github.com/vmware/govmomi/vim25/types"
)

// pbmClientFromGovmomiClient creates a new pbm client from given govmomi client.
// Can we have it in govmomi client as a field similar to tag client?
// We should not create a new pbm client every time we need it.
func pbmClientFromGovmomiClient(ctx context.Context, client *govmomi.Client) (*pbm.Client, error) {
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}

	pc, err := pbm.NewClient(ctx, client.Client)
	return pc, err

}

// PolicyIDByName finds a SPBM storage policy by name and returns its ID.
func PolicyIDByName(client *govmomi.Client, name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	pc, err := pbmClientFromGovmomiClient(ctx, client)
	if err != nil {
		return "", err
	}

	return pc.ProfileIDByName(ctx, name)
}

// policyNameByID returns storage policy name by its ID.
func PolicyNameByID(client *govmomi.Client, id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	pc, err := pbmClientFromGovmomiClient(ctx, client)
	if err != nil {
		return "", provider.ProviderError(id, "policyNameByID", err)
	}

	log.Printf("[DEBUG] Retrieving contents of storage profiles by id: %s.", id)
	profileId := []pbmtypes.PbmProfileId{
		pbmtypes.PbmProfileId{
			UniqueId: id,
		},
	}
	policies, err := pc.RetrieveContent(ctx, profileId)
	if err != nil {
		return "", provider.ProviderError(id, "RetrieveContent", err)
	}

	return policies[0].GetPbmProfile().Name, err
}

// PolicySpecByID creates and returns VirtualMachineDefinedProfileSpec by policy ID.
func PolicySpecByID(id string) []types.BaseVirtualMachineProfileSpec {
	return []types.BaseVirtualMachineProfileSpec{
		&types.VirtualMachineDefinedProfileSpec{
			ProfileId: id,
		},
	}
}

// PolicyIDByVirtualDisk fetches the storage policy associated with a virtual disk.
func PolicyIDByVirtualDisk(client *govmomi.Client, vmMOID string, diskKey int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	pc, err := pbmClientFromGovmomiClient(ctx, client)
	if err != nil {
		return "", provider.ProviderError(vmMOID, "PolicyIDByVirtualDisk", err)
	}

	pbmSOR := pbmtypes.PbmServerObjectRef{
		ObjectType: "virtualDiskId",
		Key:        fmt.Sprintf("%s:%d", vmMOID, diskKey),
	}

	policies, err := queryAssociatedProfile(ctx, pc, pbmSOR)
	if err != nil {
		return "", provider.ProviderError(vmMOID, "PolicyIDByVirtualDisk", err)
	}

	// If no policy returned then virtual disk is not associated with a policy
	if policies == nil || len(policies) == 0 {
		return "", nil
	}

	return policies[0].UniqueId, nil
}

// PolicyIDByVirtualMachine fetches the storage policy associated with a virtual machine.
func PolicyIDByVirtualMachine(client *govmomi.Client, vmMOID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	pc, err := pbmClientFromGovmomiClient(ctx, client)
	if err != nil {
		return "", provider.ProviderError(vmMOID, "PolicyIDByVirtualMachine", err)
	}

	pbmSOR := pbmtypes.PbmServerObjectRef{
		ObjectType: "virtualMachine",
		Key:        vmMOID,
	}

	policies, err := queryAssociatedProfile(ctx, pc, pbmSOR)
	if err != nil {
		return "", provider.ProviderError(vmMOID, "PolicyIDByVirtualMachine", err)
	}

	// If no policy returned then VM is not associated with a policy
	if policies == nil || len(policies) == 0 {
		return "", nil
	}

	return policies[0].UniqueId, nil
}

// queryAssociatedProfile returns the PbmProfileId of the storage policy associated with entity.
func queryAssociatedProfile(ctx context.Context, pc *pbm.Client, ref pbmtypes.PbmServerObjectRef) ([]pbmtypes.PbmProfileId, error) {
	log.Printf("[DEBUG] queryAssociatedProfile: Retrieving storage policy of server object of type [%s] and key [%s].", ref.ObjectType, ref.Key)
	req := pbmtypes.PbmQueryAssociatedProfile{
		This:   pc.ServiceContent.ProfileManager,
		Entity: ref,
	}

	res, err := methods.PbmQueryAssociatedProfile(ctx, pc, &req)
	if err != nil {
		return nil, provider.ProviderError(ref.Key, "queryAssociatedProfile", err)
	}

	return res.Returnval, nil
}
