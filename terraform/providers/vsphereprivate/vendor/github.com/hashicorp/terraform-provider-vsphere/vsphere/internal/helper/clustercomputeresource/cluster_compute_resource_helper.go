package clustercomputeresource

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/computeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// FromID locates a cluster by its managed object reference ID.
func FromID(client *govmomi.Client, id string) (*object.ClusterComputeResource, error) {
	log.Printf("[DEBUG] Locating compute cluster with ID %q", id)
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "ClusterComputeResource",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	r, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	cluster := r.(*object.ClusterComputeResource)
	log.Printf("[DEBUG] Compute cluster with ID %q found (%s)", cluster.Reference().Value, cluster.InventoryPath)
	return cluster, nil
}

func List(client *govmomi.Client) ([]*object.ClusterComputeResource, error) {
	return getComputeClusters(client, "/*")
}

func getComputeClusters(client *govmomi.Client, path string) ([]*object.ClusterComputeResource, error) {
	ctx := context.TODO()
	var dss []*object.ClusterComputeResource
	finder := find.NewFinder(client.Client, false)
	es, err := finder.ManagedObjectListChildren(ctx, path+"/*", "folder", "storagepod", "clustercompute")
	if err != nil {
		return nil, err
	}
	for _, id := range es {
		switch {
		case id.Object.Reference().Type == "ClusterComputeResource":
			ds, err := FromID(client, id.Object.Reference().Value)
			if err != nil {
				return nil, err
			}
			dss = append(dss, ds)
		case id.Object.Reference().Type == "Folder":
			newDSs, err := getComputeClusters(client, id.Path)
			if err != nil {
				return nil, err
			}
			dss = append(dss, newDSs...)
		default:
			continue
		}
	}
	return dss, nil
}

// FromPath loads a ClusterComputeResource from its path. The datacenter is
// optional if the path is specific enough to not require it.
func FromPath(client *govmomi.Client, name string, dc *object.Datacenter) (*object.ClusterComputeResource, error) {
	finder := find.NewFinder(client.Client, false)
	if dc != nil {
		log.Printf("[DEBUG] Attempting to locate compute cluster %q in datacenter %q", name, dc.InventoryPath)
		finder.SetDatacenter(dc)
	} else {
		log.Printf("[DEBUG] Attempting to locate compute cluster at absolute path %q", name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	return finder.ClusterComputeResource(ctx, name)
}

// Properties is a convenience method that wraps fetching the
// ClusterComputeResource MO from its higher-level object.
func Properties(cluster *object.ClusterComputeResource) (*mo.ClusterComputeResource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	var props mo.ClusterComputeResource
	if err := cluster.Properties(ctx, cluster.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

// Create creates a ClusterComputeResource in a supplied folder. The resulting
// ClusterComputeResource is returned.
func Create(f *object.Folder, name string, spec types.ClusterConfigSpecEx) (*object.ClusterComputeResource, error) {
	log.Printf("[DEBUG] Creating compute cluster %q", fmt.Sprintf("%s/%s", f.InventoryPath, name))
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	cluster, err := f.CreateCluster(ctx, name, spec)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// Rename renames a ClusterComputeResource.
func Rename(cluster *object.ClusterComputeResource, name string) error {
	log.Printf("[DEBUG] Renaming compute cluster %q to %s", cluster.InventoryPath, name)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := cluster.Rename(ctx, name)
	if err != nil {
		return err
	}
	return task.Wait(ctx)
}

// MoveToFolder is a complex method that moves a ClusterComputeResource to a given relative
// compute folder path. "Relative" here means relative to a datacenter, which
// is discovered from the current ClusterComputeResource path.
func MoveToFolder(client *govmomi.Client, cluster *object.ClusterComputeResource, relative string) error {
	f, err := folder.HostFolderFromObject(client, cluster, relative)
	if err != nil {
		return err
	}
	return folder.MoveObjectTo(cluster.Reference(), f)
}

// HasChildren checks to see if a compute cluster has any child items (hosts
// and virtual machines) and returns true if that is the case. This is useful
// when checking to see if a compute cluster is safe to delete - destroying a
// compute cluster in vSphere destroys *all* children if at all possible
// (including removing hosts and virtual machines), so extra verification is
// necessary to prevent accidental removal.
func HasChildren(cluster *object.ClusterComputeResource) (bool, error) {
	return computeresource.HasChildren(cluster)
}

// Reconfigure reconfigures a cluster. This just gets dispatched to
// computeresource as both methods are the same.
func Reconfigure(cluster *object.ClusterComputeResource, spec *types.ClusterConfigSpecEx) error {
	return computeresource.Reconfigure(cluster, spec)
}

// Delete destroys a ClusterComputeResource.
func Delete(cluster *object.ClusterComputeResource) error {
	log.Printf("[DEBUG] Deleting compute cluster %q", cluster.InventoryPath)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := cluster.Destroy(ctx)
	if err != nil {
		return err
	}
	return task.Wait(ctx)
}

func Hosts(cluster *object.ClusterComputeResource) ([]*object.HostSystem, error) {
	ctx := context.TODO()
	return cluster.Hosts(ctx)
}

// MoveHostsInto moves all of the supplied hosts into the cluster. All virtual
// machines are moved to the cluster's root resource pool and any resource
// pools on the host itself are deleted.
func MoveHostsInto(client *govmomi.Client, cluster *object.ClusterComputeResource, hosts []*object.HostSystem) error {
	var hsNames []string
	var hsRefs []types.ManagedObjectReference

	seenClusters := map[string]int{}

	for _, hs := range hosts {
		hsNames = append(hsNames, hs.Name())
		hsRefs = append(hsRefs, hs.Reference())
		hsProps, err := hostsystem.Properties(hs)
		if err != nil {
			return fmt.Errorf("while fetching properties for host %q: %s", hs.Reference().Value, err)
		}

		if hsProps.Parent.Type == "ClusterComputeResource" {
			cRef := hsProps.Parent.Value
			parentCluster, err := computeresource.BaseFromReference(client, hsProps.Parent.Reference())
			if err != nil {
				return fmt.Errorf("while retrieving parent cluster (%q) object for host %q: %s", cluster.Reference().Value, hs.Reference().Value, err)
			}
			c, err := computeresource.BaseProperties(parentCluster)
			if err != nil {
				return fmt.Errorf("while retrieving parent cluster (%q) properties for host %q: %s", cluster.Reference().Value, hs.Reference().Value, err)
			}

			var evacuate bool
			hostsLeft, ok := seenClusters[cRef]
			if !ok {
				seenClusters[cRef] = len(c.Host)
				if hostsLeft > 1 {
					evacuate = true
				}
			} else {
				evacuate = false
				if hostsLeft > 1 {
					evacuate = true
				}
			}

			totalVMTimeout := provider.DefaultAPITimeout * time.Duration(len(hsProps.Vm)+1)
			err = hostsystem.EnterMaintenanceMode(hs, totalVMTimeout, evacuate)
			if err != nil {
				return fmt.Errorf("while putting host %q in maintenance mode: %s", hs.Reference().Value, err)
			}
		}
	}
	log.Printf("[DEBUG] Adding hosts into cluster %q: %s", cluster.Name(), strings.Join(hsNames, ", "))

	req := types.MoveInto_Task{
		This: cluster.Reference(),
		Host: hsRefs,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	resp, err := methods.MoveInto_Task(ctx, cluster.Client(), &req)
	if err != nil {
		return err
	}

	task := object.NewTask(cluster.Client(), resp.Returnval)
	return task.Wait(ctx)
}

// MoveHostsOutOf moves a supplied list of hosts out of the specified cluster.
// The host is moved to the root host folder for the datacenter that the
// cluster is in.
//
// The host is placed into maintenance mode with evacuate flagged on, ensuring
// that as many VMs as possible are moved out of the host before removing it
// from the cluster. The effectiveness of this operation is dictated by the
// cluster's DRS settings, which also affects if this means that the task will
// block and require manual intervention. The supplied timeout is passed to the
// maintenance mode operations, and represents the timeout in seconds.
//
// Individual hosts are taken out of maintenance mode after its operation is
// complete.
func MoveHostsOutOf(cluster *object.ClusterComputeResource, hosts []*object.HostSystem, timeout int) error {
	for _, host := range hosts {
		if err := moveHostOutOf(cluster, host, timeout); err != nil {
			return err
		}
	}
	return nil
}

func moveHostOutOf(cluster *object.ClusterComputeResource, host *object.HostSystem, timeout int) error {
	// Place the host into maintenance mode. This blocks until the host is ready.
	timeoutDuration := time.Duration(timeout) * time.Second
	if err := hostsystem.EnterMaintenanceMode(host, timeoutDuration, true); err != nil {
		return fmt.Errorf("error putting host %q into maintenance mode: %s", host.Name(), err)
	}

	// Host should be ready to move out of the cluster now.
	f, err := folder.HostFolderFromObject(&govmomi.Client{Client: cluster.Client()}, host, "/")
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Moving host %q out of cluster %q and to folder %q", host.Name(), cluster.Name(), f.InventoryPath)
	if err := folder.MoveObjectTo(host.Reference(), f); err != nil {
		return fmt.Errorf("error moving host %q out of cluster %q: %s", host.Name(), cluster.Name(), err)
	}

	// Move the host out of maintenance mode now that it's out of the cluster.
	if err := hostsystem.ExitMaintenanceMode(host, timeoutDuration); err != nil {
		return fmt.Errorf("error taking host %q out of maintenance mode: %s", host.Name(), err)
	}

	log.Printf("[DEBUG] Host %q moved out of cluster %q successfully", host.Name(), cluster.Name())
	return nil
}
