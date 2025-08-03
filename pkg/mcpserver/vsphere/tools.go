package vsphere

import (
	"context"
	"encoding/json"
	"path"
	"slices"

	"github.com/sirupsen/logrus"

	icvsphere "github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func GetVSphereTopology(username, password, server string) (string, error) {

	vspherePlatform := &vsphere.Platform{}
	// add a new vcenter to the vsphere platform
	vspherePlatform.VCenters = append(vspherePlatform.VCenters, vsphere.VCenter{
		Server:   server,
		Username: username,
		Password: password,
	})

	vmeta := icvsphere.NewMetadata()

	_ = vmeta.AddCredentials(server, username, password)

	session, err := vmeta.Session(context.Background(), server)
	if err != nil {
		return "", err
	}

	datacenters, err := session.Finder.DatacenterList(context.Background(), "/...")
	if err != nil {
		return "", err
	}

	for _, datacenter := range datacenters {
		logrus.Infof("Found datacenter %s", datacenter.InventoryPath)
		clusters, err := session.Finder.ClusterComputeResourceList(context.Background(), path.Join(datacenter.InventoryPath, "/..."))
		if err != nil {
			return "", err
		}

		// get all datastores in the host
		datastores, err := session.Finder.DatastoreList(context.Background(), path.Join(datacenter.InventoryPath, "datastore", "..."))
		if err != nil {
			logrus.Errorf("Error finding datastores: %v", err)
			return "", err
		}

		for _, cluster := range clusters {
			logrus.Infof("Found cluster %s", cluster.InventoryPath)
			networkObjects, err := session.Finder.NetworkList(context.Background(), path.Join(cluster.InventoryPath, "/*"))
			if err != nil {
				return "", err
			}

			networks := make([]string, len(networkObjects))
			for i, network := range networkObjects {
				logrus.Infof("Network object %s", network.GetInventoryPath())
				networks[i] = network.GetInventoryPath()
			}

			for _, ds := range datastores {

				// check if the datastore is already in the failure domain
				logrus.Infof("Datastore object %s", ds.InventoryPath)
				failureDomain := vsphere.FailureDomain{
					Server: server,
					Topology: vsphere.Topology{
						Datacenter:     datacenter.InventoryPath,
						ComputeCluster: cluster.InventoryPath,
						Networks:       networks,
						Datastore:      ds.InventoryPath,
					},
				}
				vspherePlatform.FailureDomains = append(vspherePlatform.FailureDomains, failureDomain)
			}

		}
	}

	for _, fd := range vspherePlatform.FailureDomains {
		dc := path.Base(fd.Topology.Datacenter)
		if !slices.Contains(vspherePlatform.VCenters[0].Datacenters, dc) {
			vspherePlatform.VCenters[0].Datacenters = append(vspherePlatform.VCenters[0].Datacenters, dc)
		}
	}
	platform := types.Platform{VSphere: vspherePlatform}

	jsonBytes, err := json.Marshal(platform)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil

}
