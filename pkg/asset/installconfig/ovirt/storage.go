package ovirt

import (
	"fmt"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types/ovirt"
)

func askStorage(c *ovirtsdk4.Connection, p *ovirt.Platform, clusterName string) error {
	var storageDomainName string
	var domainsForCluster = make(map[string]*ovirtsdk4.StorageDomain)
	var domainNames []string
	systemService := c.SystemService()
	domainsSearch, err := systemService.StorageDomainsService().List().Search(fmt.Sprintf("cluster.name=%s", clusterName)).Send()
	if err != nil {
		return err
	}
	domains, ok := domainsSearch.StorageDomains()
	if !ok {
		return fmt.Errorf("there are no available storage domains for cluster %s", clusterName)
	}

	for _, domain := range domains.Slice() {
		domainsForCluster[domain.MustName()] = domain
		domainNames = append(domainNames, domain.MustName())
	}
	if err := survey.AskOne(&survey.Select{
		Message: "Storage domain",
		Help:    "The storage domain will be used to create the disks of all the cluster nodes.",
		Options: domainNames,
	},
		&storageDomainName,
		survey.WithValidator(func(ans interface{}) error {
			choice := ans.(core.OptionAnswer).Value
			sort.Strings(domainNames)
			i := sort.SearchStrings(domainNames, choice)
			if i == len(domainNames) || domainNames[i] != choice {
				return fmt.Errorf("invalid storage domain %s", choice)
			}
			domain, ok := domainsForCluster[choice]
			if !ok {
				return fmt.Errorf("cannot find storage domain id for the storage domain %s", storageDomainName)
			}
			p.StorageDomainID = domain.MustId()
			return nil
		})); err != nil {
		return errors.Wrap(err, "failed UserInput")
	}
	return nil
}
