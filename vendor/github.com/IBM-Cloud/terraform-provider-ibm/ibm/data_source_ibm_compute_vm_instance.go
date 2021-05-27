// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
)

func dataSourceIBMComputeVmInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMComputeVmInstanceRead,

		Schema: map[string]*schema.Schema{

			"hostname": &schema.Schema{
				Description: "The hostname of the virtual guest",
				Type:        schema.TypeString,
				Required:    true,
			},

			"domain": &schema.Schema{
				Description: "The domain of the virtual guest",
				Type:        schema.TypeString,
				Required:    true,
			},

			"datacenter": &schema.Schema{
				Description: "Datacenter in which the virtual guest is deployed",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"cores": &schema.Schema{
				Description: "Number of cpu cores",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"status": &schema.Schema{
				Description: "The VSI status",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"last_known_power_state": &schema.Schema{
				Description: "The last known power state of a virtual guest in the event the guest is turned off outside of IMS or has gone offline.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"public_interface_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_interface_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"power_state": &schema.Schema{
				Description: "The current power state of a virtual guest.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"most_recent": &schema.Schema{
				Description: "If true and multiple entries are found, the most recently created virtual guest is used. " +
					"If false, an error is returned",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"public_subnet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"private_subnet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_address_private": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ip_address_id_private": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv6_address_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"public_ipv6_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ipv6_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"secondary_ip_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMComputeVmInstanceRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetAccountService(sess)

	hostname := d.Get("hostname").(string)
	domain := d.Get("domain").(string)
	mostRecent := d.Get("most_recent").(bool)

	vgs, err := service.
		Filter(filter.Build(filter.Path("virtualGuests.hostname").Eq(hostname),
			filter.Path("virtualGuests.domain").Eq(domain))).Mask(
		"hostname,domain,primaryIpAddress,primaryBackendIpAddress,startCpus,datacenter[id,name,longName],statusId,status,id,powerState,lastKnownPowerState,createDate,primaryNetworkComponent[id, primaryIpAddressRecord[subnet,guestNetworkComponentBinding[ipAddressId]]," +
			"primaryVersion6IpAddressRecord[subnet,guestNetworkComponentBinding[ipAddressId]]]," +
			"primaryBackendNetworkComponent[id],primaryBackendNetworkComponent[networkVlan[id]," +
			"securityGroupBindings[securityGroup]]",
	).GetVirtualGuests()

	if err != nil {
		return fmt.Errorf("Error retrieving virtual guest details for host %s: %s", hostname, err)
	}
	if len(vgs) == 0 {
		return fmt.Errorf("No virtual guest with hostname %s and domain  %s", hostname, domain)
	}

	var vg datatypes.Virtual_Guest
	if len(vgs) > 1 {
		if mostRecent {
			vg = mostRecentVirtualGuest(vgs)
		} else {
			return fmt.Errorf(
				"More than one virtual guest found with host matching [%s] and domain "+
					"matching [%s]. Set 'most_recent' to true in your configuration to force the most recent virtual guest "+
					"to be used", hostname, domain)
		}
	} else {
		vg = vgs[0]
	}

	d.SetId(fmt.Sprintf("%d", *vg.Id))
	d.Set("hostname", vg.Hostname)
	d.Set("domain", vg.Domain)

	if vg.Datacenter != nil {
		d.Set("datacenter", *vg.Datacenter.Name)
	}
	d.Set("cores", *vg.StartCpus)
	if vg.Status != nil {
		d.Set("status", vg.Status.KeyName)
	}
	if vg.PowerState != nil {
		d.Set("power_state", vg.PowerState.KeyName)
	}
	if vg.LastKnownPowerState != nil {
		d.Set("last_known_power_state", vg.LastKnownPowerState.KeyName)
	}
	d.Set("public_interface_id", vg.PrimaryNetworkComponent.Id)
	d.Set("private_interface_id", vg.PrimaryBackendNetworkComponent.Id)
	d.Set("ipv4_address", vg.PrimaryIpAddress)
	d.Set("ipv4_address_private", vg.PrimaryBackendIpAddress)
	if vg.PrimaryNetworkComponent.PrimaryIpAddressRecord != nil {
		d.Set("ip_address_id", *vg.PrimaryNetworkComponent.PrimaryIpAddressRecord.GuestNetworkComponentBinding.IpAddressId)
	}
	if vg.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord != nil {
		d.Set("ip_address_id_private",
			*vg.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord.GuestNetworkComponentBinding.IpAddressId)
	}
	if vg.PrimaryNetworkComponent.PrimaryIpAddressRecord != nil {
		d.Set("public_subnet_id", vg.PrimaryNetworkComponent.PrimaryIpAddressRecord.SubnetId)
	}

	if vg.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord != nil {
		d.Set("private_subnet_id", vg.PrimaryBackendNetworkComponent.PrimaryIpAddressRecord.SubnetId)
	}

	if vg.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord != nil {
		d.Set("ipv6_address", *vg.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.IpAddress)
		d.Set("ipv6_address_id", *vg.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.GuestNetworkComponentBinding.IpAddressId)
		publicSubnet := vg.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.Subnet
		d.Set(
			"public_ipv6_subnet",
			fmt.Sprintf("%s/%d", *publicSubnet.NetworkIdentifier, *publicSubnet.Cidr),
		)
		d.Set("public_ipv6_subnet_id", vg.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.SubnetId)
	}

	err = readSecondaryIPAddresses(d, meta, vg.PrimaryIpAddress)
	if err != nil {
		return fmt.Errorf("Error retrieving virtual guest details for host %s: %s", hostname, err)
	}
	return nil
}

type virtualGuests []datatypes.Virtual_Guest

func (k virtualGuests) Len() int { return len(k) }

func (k virtualGuests) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

func (k virtualGuests) Less(i, j int) bool {
	return k[i].CreateDate.Before(k[j].CreateDate.Time)
}

func mostRecentVirtualGuest(keys virtualGuests) datatypes.Virtual_Guest {
	sortedKeys := keys
	sort.Sort(sortedKeys)
	return sortedKeys[len(sortedKeys)-1]
}
