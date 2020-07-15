package packet

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func bgpNeighborSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"address_family": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"customer_as": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"customer_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"md5_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"md5_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"multihop": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"peer_as": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"peer_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"routes_in": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bgpRouteSchema(),
			},
			"routes_out": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bgpRouteSchema(),
			},
		},
	}
}

func bgpRouteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"route": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exact": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePacketDeviceBGPNeighbors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketDeviceBGPNeighborsRead,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bgp_neighbors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bgpNeighborSchema(),
			},
		},
	}
}

func dataSourcePacketDeviceBGPNeighborsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	deviceID := d.Get("device_id").(string)

	bgpNeighborsRaw, _, err := client.Devices.ListBGPNeighbors(deviceID, nil)
	if err != nil {
		return err
	}

	d.Set("bgp_neighbors", getBgpNeighbors(bgpNeighborsRaw))
	d.SetId(deviceID)
	return nil
}

func getRoutesSlice(routes []packngo.BGPRoute) []map[string]interface{} {
	ret := []map[string]interface{}{}
	for _, r := range routes {
		ret = append(ret, map[string]interface{}{
			"route": r.Route, "exact": r.Exact,
		})
	}
	return ret
}

func getBgpNeighbors(ns []packngo.BGPNeighbor) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 1)
	for _, n := range ns {
		neighbor := map[string]interface{}{
			"address_family": n.AddressFamily,
			"customer_as":    n.CustomerAs,
			"customer_ip":    n.CustomerIP,
			"md5_enabled":    n.Md5Enabled,
			"md5_password":   n.Md5Password,
			"multihop":       n.Multihop,
			"peer_as":        n.PeerAs,
			"peer_ips":       n.PeerIps,
			"routes_in":      getRoutesSlice(n.RoutesIn),
			"routes_out":     getRoutesSlice(n.RoutesOut),
		}
		ret = append(ret, neighbor)
	}
	return ret
}
