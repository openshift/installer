package ovirt

import (
	"fmt"
	"sort"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/ovirt"
)

func askCluster(c *ovirtsdk4.Connection, p *ovirt.Platform) (string, error) {
	var clusterName string
	var clusterByNames = make(map[string]*ovirtsdk4.Cluster)
	var clusterNames []string
	systemService := c.SystemService()

	dcResp, err := datacentersAvailable(c, "")
	if err != nil {
		return "", err
	}

	datacenters := dcResp.MustDataCenters()
	for _, dc := range datacenters.Slice() {
		dcService := systemService.DataCentersService().DataCenterService(dc.MustId())
		logrus.Debug("Datacenter:", dc.MustName())
		clusters, err := dcService.ClustersService().List().Send()
		if err != nil {
			return "", errors.Wrap(err, "failed to list clusters")
		}
		clusterSlice := clusters.MustClusters()
		for _, cluster := range clusterSlice.Slice() {
			logrus.Debug("\tcluster:", cluster.MustName())
			clusterByNames[cluster.MustName()] = cluster
			clusterNames = append(clusterNames, cluster.MustName())
		}
	}
	if err := survey.AskOne(&survey.Select{
		Message: "Cluster",
		Help:    "The Cluster where the VMs will be created.",
		Options: clusterNames,
	},
		&clusterName,
		func(ans interface{}) error {
			choice := ans.(string)
			sort.Strings(clusterNames)
			i := sort.SearchStrings(clusterNames, choice)
			if i == len(clusterNames) || clusterNames[i] != choice {
				return fmt.Errorf("invalid cluster %s", choice)
			}
			cl, ok := clusterByNames[choice]
			if !ok {
				return fmt.Errorf("cannot find a cluster id for the cluster name %s", clusterName)
			}
			p.ClusterID = cl.MustId()
			return nil
		}); err != nil {
		return clusterName, errors.Wrap(err, "failed UserInput")
	}
	return clusterName, nil
}
