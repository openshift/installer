package virtualdisk

import (
	"context"
	"fmt"
	"log"
	"path"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

// DatastorePathFromString is a convenience method that returns a
// fully-populated DatastorePath from a string containing a datastore path. A
// flag indicating a successful parsing is also returned.
func DatastorePathFromString(p string) (*object.DatastorePath, bool) {
	dp := new(object.DatastorePath)
	success := dp.FromString(p)
	log.Printf("[DEBUG] DatastorePathFromString: success: %t, path: %q", success, p)
	return dp, success
}

// IsVmdkDatastorePath ensures that a string can be parsed as a datastore path
// pointing to a virtual disk. This only checks the validity of the path, not
// whether or not the file exists.
func IsVmdkDatastorePath(p string) bool {
	dp, success := DatastorePathFromString(p)
	if !success {
		log.Printf("[DEBUG] IsVmdkDatastorePath: %q is not a datastore path", p)
		return false
	}
	isVMDK := dp.IsVMDK()
	log.Printf("[DEBUG] IsVmdkDatastorePath: %q %s a datastore path", p, structure.LogCond(isVMDK, "is", "is not"))
	return isVMDK
}

// Move moves a virtual disk from one location to another. The move is not
// forced.
//
// srcPath needs to be a datastore path (ie: "[datastore1] vm/vm.vmdk"),
// however the destination path (dstPath) can be a simple path - if it is, the
// source datastore is used and dstDC is ignored. Further, if dstPath has no
// directory, the directory of srcPath is used. dstDC can be nil if the
// destination datastore is in the same datacenter.
//
// The new datastore path is returned along with any error, to avoid the need
// to re-calculate the path separately.
func Move(client *govmomi.Client, srcPath string, srcDC *object.Datacenter, dstPath string, dstDC *object.Datacenter) (string, error) {
	vdm := object.NewVirtualDiskManager(client.Client)
	if srcDC == nil {
		return "", fmt.Errorf("source datacenter cannot be nil")
	}
	if !IsVmdkDatastorePath(srcPath) {
		return "", fmt.Errorf("path %q is not a datastore path", srcPath)
	}
	if !IsVmdkDatastorePath(dstPath) {
		ddp := dstDataStorePathFromLocalSrc(srcPath, dstPath)
		// One more validation
		if !IsVmdkDatastorePath(ddp) {
			return "", fmt.Errorf("path %q is not a valid destination path", dstPath)
		}
		dstPath = ddp
		dstDC = nil
	}
	log.Printf("[DEBUG] Moving virtual disk from %q in datacenter %s to destination %s%s",
		srcPath,
		srcDC,
		dstPath,
		structure.LogCond(dstDC != nil, fmt.Sprintf("in datacenter %s", dstDC), ""),
	)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := vdm.MoveVirtualDisk(ctx, srcPath, srcDC, dstPath, dstDC, false)
	if err != nil {
		return "", err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer tcancel()
	if err := task.Wait(tctx); err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Virtual disk %q in datacenter %s successfully moved to destination %s%s",
		srcPath,
		srcDC,
		dstPath,
		structure.LogCond(dstDC != nil, fmt.Sprintf("in datacenter %s", dstDC), ""),
	)
	return dstPath, nil
}

// QueryDiskType queries the disk type of the specified virtual disk.
func QueryDiskType(client *govmomi.Client, name string, dc *object.Datacenter) (types.VirtualDiskType, error) {
	di, err := FromPath(client, name, dc)
	if err != nil {
		return types.VirtualDiskType(""), err
	}
	t := di.DiskType
	log.Printf("[DEBUG] QueryDiskType: Disk %q is of type %q", name, t)
	return types.VirtualDiskType(t), nil
}

func dstDataStorePathFromLocalSrc(src, dst string) string {
	ddpTmp, success := DatastorePathFromString(dst)
	if success {
		dst = ddpTmp.Path
	}
	dstDP, success := DatastorePathFromString(src)
	if !success {
		panic(fmt.Errorf("dstDataStorePathFromLocalSrc: Expected valid source, got %q", src))
	}
	if path.Dir(dst) == "." {
		dstDP.Path = path.Join(path.Dir(dstDP.Path), dst)
	} else {
		dstDP.Path = dst
	}
	return dstDP.String()
}

// Delete deletes the virtual disk at the specified datastore path.
func Delete(client *govmomi.Client, name string, dc *object.Datacenter) error {
	if dc == nil {
		return fmt.Errorf("datacenter cannot be nil")
	}
	log.Printf("[DEBUG] Deleting virtual disk %q in datacenter %s", name, dc)
	vdm := object.NewVirtualDiskManager(client.Client)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := vdm.DeleteVirtualDisk(ctx, name, dc)
	if err != nil {
		return err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer tcancel()
	if err := task.Wait(tctx); err != nil {
		return err
	}
	log.Printf("[DEBUG] Virtual disk %q in datacenter %s deleted succesfully", name, dc)
	return nil
}

// FromPath loads a datastore from its path.
func FromPath(client *govmomi.Client, p string, dc *object.Datacenter) (*object.VirtualDiskInfo, error) {
	vdm := object.NewVirtualDiskManager(client.Client)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	di, err := vdm.QueryVirtualDiskInfo(ctx, p, dc, false)
	if err != nil {
		return nil, err
	}
	return &di[0], nil
}
