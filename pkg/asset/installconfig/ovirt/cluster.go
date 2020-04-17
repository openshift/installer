package ovirt

import (
	"fmt"
	"sort"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/ovirt"
)

func askCluster(c *ovirtsdk4.Connection, p *ovirt.Platform) (string, error) {
	var clusterName string
	var clusterByNames = make(map[string]*ovirtsdk4.Cluster)
	var clusterNames []string
	systemService := c.SystemService()
	response, err := systemService.ClustersService().List().Send()
	if err != nil {
		return "", err
	}
	clusters, ok := response.Clusters()
	if !ok {
		return "", fmt.Errorf("there are no available clusters")
	}

	for _, cluster := range clusters.Slice() {
		clusterByNames[cluster.MustName()] = cluster
		clusterNames = append(clusterNames, cluster.MustName())
	}
	err = survey.AskOne(&survey.Select{
		Message: "oVirt cluster",
		Help:    "The oVirt cluster where the VMs will be created.",
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
		})
	return clusterName, err
}
