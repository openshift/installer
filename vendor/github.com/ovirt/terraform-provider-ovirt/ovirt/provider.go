// Copyright (C) 2017 Battelle Memorial Institute
// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

type providerContext struct {
	semaphores *semaphoreProvider
}

func ProviderContext() func() terraform.ResourceProvider {
	c := &providerContext{
		semaphores: newSemaphoreProvider(),
	}
	return c.Provider
}

// Provider returns oVirt provider configuration
func (c *providerContext) Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVIRT_USERNAME", ""),
				Description: "Login username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVIRT_PASSWORD", ""),
				Description: "Login password",
				Sensitive:   true,
			},
			"insecure": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVIRT_INSECURE", false),
				Description: "Skip certificate verification",
				Sensitive:   false,
			},
			"cafile": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVIRT_CAFILE", ""),
				Description: "File containing the CA certificate in PEM format",
				Sensitive:   false,
			},
			"ca_bundle": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVIRT_CA_BUNDLE", ""),
				Description: "CA certificate in PEM format",
				Sensitive:   true,
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVIRT_URL", ""),
				Description: "Ovirt server url",
			},
			"headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Additional headers to be added to each API call",
			},
		},
		ConfigureFunc: ConfigureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"ovirt_affinity_group":  resourceOvirtAffinityGroup(),
			"ovirt_vm":              resourceOvirtVM(c),
			"ovirt_template":        resourceOvirtTemplate(),
			"ovirt_disk":            resourceOvirtDisk(),
			"ovirt_disk_attachment": resourceOvirtDiskAttachment(),
			"ovirt_datacenter":      resourceOvirtDataCenter(),
			"ovirt_network":         resourceOvirtNetwork(),
			"ovirt_vnic":            resourceOvirtVnic(),
			"ovirt_vnic_profile":    resourceOvirtVnicProfile(),
			"ovirt_snapshot":        resourceOvirtSnapshot(),
			"ovirt_storage_domain":  resourceOvirtStorageDomain(),
			"ovirt_tag":             resourceOvirtTag(),
			"ovirt_user":            resourceOvirtUser(),
			"ovirt_cluster":         resourceOvirtCluster(),
			"ovirt_mac_pool":        resourceOvirtMacPool(),
			"ovirt_host":            resourceOvirtHost(),
			"ovirt_image_transfer":  resourceOvirtImageTransfer(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ovirt_disks":          dataSourceOvirtDisks(),
			"ovirt_datacenters":    dataSourceOvirtDataCenters(),
			"ovirt_networks":       dataSourceOvirtNetworks(),
			"ovirt_clusters":       dataSourceOvirtClusters(),
			"ovirt_storagedomains": dataSourceOvirtStorageDomains(),
			"ovirt_vnic_profiles":  dataSourceOvirtVNicProfiles(),
			"ovirt_authzs":         dataSourceOvirtAuthzs(),
			"ovirt_users":          dataSourceOvirtUsers(),
			"ovirt_mac_pools":      dataSourceOvirtMacPools(),
			"ovirt_vms":            dataSourceOvirtVMs(),
			"ovirt_hosts":          dataSourceOvirtHosts(),
			"ovirt_nics":           dataSourceOvirtNics(),
			"ovirt_templates":      dataSourceOvirtTemplates(),
		},
	}
}

// ConfigureProvider initializes the API connection object by config items
func ConfigureProvider(d *schema.ResourceData) (interface{}, error) {
	insecure := d.Get("insecure").(bool)
	caFile := d.Get("cafile").(string)
	caCert := []byte(d.Get("ca_bundle").(string))

	if !insecure && caFile == "" && len(caCert) == 0 {
		return nil, fmt.Errorf("either insecure must be set or one of cafile and ca_bundle must be set")
	}

	connBuilder := ovirtsdk4.NewConnectionBuilder().
		URL(d.Get("url").(string)).
		Username(d.Get("username").(string)).
		Password(d.Get("password").(string)).
		CAFile(caFile).
		CACert(caCert).
		Insecure(insecure)

	// Set headers if needed
	if v, ok := d.GetOk("headers"); ok {
		headers := map[string]string{}
		for k, v := range v.(map[string]interface{}) {
			headers[k] = v.(string)
		}
		connBuilder.Headers(headers)
	}

	return connBuilder.Build()
}

func newSemaphoreProvider() *semaphoreProvider {
	return &semaphoreProvider{
		lock:       &sync.Mutex{},
		semaphores: map[string]chan struct{}{},
	}
}

type semaphoreProvider struct {
	lock       *sync.Mutex
	semaphores map[string]chan struct{}
}

func (s *semaphoreProvider) Lock(semName string, capacity uint) {
	if capacity < 1 {
		panic(fmt.Sprintf("Invalid semaphoreProvider capacity %d for sem %s", capacity, semName))
	}
	s.lock.Lock()
	if _, ok := s.semaphores[semName]; !ok {
		s.semaphores[semName] = make(chan struct{}, capacity)
	}
	s.lock.Unlock()
	s.semaphores[semName] <- struct{}{}
}

func (s *semaphoreProvider) Unlock(semName string) {
	s.lock.Lock()
	if _, ok := s.semaphores[semName]; !ok {
		panic(fmt.Sprintf("semaphoreProvider unlock called before lock: %s", semName))
	}
	s.lock.Unlock()
	<-s.semaphores[semName]
}
