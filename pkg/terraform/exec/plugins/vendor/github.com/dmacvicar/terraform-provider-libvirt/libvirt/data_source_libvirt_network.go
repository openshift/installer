package libvirt

import (
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

// a libvirt network DNS host template datasource
//
// Datasource example:
//
// data "libvirt_network_dns_host_template" "k8smasters" {
//   count = "${var.master_count}"
//   ip = "${var.master_ips[count.index]}"
//   hostname = "master-${count.index}"
// }
//
// resource "libvirt_network" "k8snet" {
//   ...
//   dns = [{
//     hosts = [ "${flatten(data.libvirt_network_dns_host_template.k8smasters.*.rendered)}" ]
//   }]
//   ...
// }
//
func datasourceLibvirtNetworkDNSHostTemplate() *schema.Resource {
	return &schema.Resource{
		Read: resourceLibvirtNetworkDNSHostRead,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rendered": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func resourceLibvirtNetworkDNSHostRead(d *schema.ResourceData, meta interface{}) error {
	dnsHost := map[string]interface{}{}
	if address, ok := d.GetOk("ip"); ok {
		ip := net.ParseIP(address.(string))
		if ip == nil {
			return fmt.Errorf("Could not parse address '%s'", address)
		}
		dnsHost["ip"] = ip.String()
	}
	if hostname, ok := d.GetOk("hostname"); ok {
		dnsHost["hostname"] = hostname.(string)
	}
	d.Set("rendered", dnsHost)
	d.SetId(strconv.Itoa(hashcode.String(fmt.Sprintf("%v", dnsHost))))

	return nil
}

// a libvirt network DNS SRV template datasource
//
// Datasource example:
//
// data "libvirt_network_dns_srv_template" "etcd_cluster" {
//   count = "${var.etcd_count}"
//   service = "etcd-server"
//   protocol = "tcp"
//   domain = "${discovery_domain}"
//   target = "${var.cluster_name}-etcd-${count.index}.${discovery_domain}"
// }
//
// resource "libvirt_network" "k8snet" {
//   ...
//   dns = [{
//     srvs = [ "${flatten(data.libvirt_network_dns_srv_template.etcd_cluster.*.rendered)}" ]
//   }]
//   ...
// }
//
func datasourceLibvirtNetworkDNSSRVTemplate() *schema.Resource {
	return &schema.Resource{
		Read: resourceLibvirtNetworkDNSSRVRead,
		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rendered": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceLibvirtNetworkDNSSRVRead(d *schema.ResourceData, meta interface{}) error {
	dnsSRV := map[string]interface{}{}
	if service, ok := d.GetOk("service"); ok {
		dnsSRV["service"] = service.(string)
	}
	if protocol, ok := d.GetOk("protocol"); ok {
		dnsSRV["protocol"] = protocol.(string)
	}
	if domain, ok := d.GetOk("domain"); ok {
		dnsSRV["domain"] = domain.(string)
	}
	if target, ok := d.GetOk("target"); ok {
		dnsSRV["target"] = target.(string)
	}
	if port, ok := d.GetOk("port"); ok {
		dnsSRV["port"] = port.(string)
	}
	if weight, ok := d.GetOk("weight"); ok {
		dnsSRV["weight"] = weight.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		dnsSRV["priority"] = priority.(string)
	}
	d.Set("rendered", dnsSRV)
	d.SetId(strconv.Itoa(hashcode.String(fmt.Sprintf("%v", dnsSRV))))

	return nil
}
