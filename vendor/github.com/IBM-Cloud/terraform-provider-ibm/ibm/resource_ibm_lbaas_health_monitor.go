// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMLbaasHealthMonitor() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMLbaasHealthMonitorCreate,
		Read:     resourceIBMLbaasHealthMonitorRead,
		Delete:   resourceIBMLbaasHealthMonitorDelete,
		Update:   resourceIBMLbaasHealthMonitorUpdate,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS", "TCP"}),
				Description:  "Protocol value",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validatePortRange(1, 65535),
				Description:  "Port number",
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validateInterval,
				Description:  "Interval value",
			},
			"max_retries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validateMaxRetries,
				Description:  "Maximum retry counts",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validateTimeout,
				Description:  "Timeout in seconds",
			},
			"url_path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "/",
				ValidateFunc: validateURLPath,
				Description:  "URL Path",
			},
			"monitor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Monitor ID",
			},
			"lbaas_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "LBAAS id",
			},
		},
	}
}

func resourceIBMLbaasHealthMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	healthMonitorService := services.GetNetworkLBaaSHealthMonitorService(sess.SetRetries(0))

	lbaasID := d.Get("lbaas_id").(string)
	healthMonitors := make([]datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration, 0, 1)
	healthMonitor := datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration{
		BackendPort:       sl.Int(d.Get("port").(int)),
		BackendProtocol:   sl.String(d.Get("protocol").(string)),
		HealthMonitorUuid: sl.String(d.Get("monitor_id").(string)),
		Interval:          sl.Int(d.Get("interval").(int)),
		Timeout:           sl.Int(d.Get("timeout").(int)),
		MaxRetries:        sl.Int(d.Get("max_retries").(int)),
		UrlPath:           sl.String(d.Get("url_path").(string)),
	}

	healthMonitors = append(healthMonitors, healthMonitor)

	_, err := waitForLbaasLBActive(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
	}

	_, err = healthMonitorService.UpdateLoadBalancerHealthMonitors(sl.String(lbaasID), healthMonitors)
	if err != nil {
		return fmt.Errorf("Error adding health monitors: %#v", err)
	}
	_, err = waitForLbaasLBActive(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
	}
	d.SetId(fmt.Sprintf("%s/%s", lbaasID, d.Get("monitor_id").(string)))
	return resourceIBMLbaasHealthMonitorRead(d, meta)
}

func resourceIBMLbaasHealthMonitorRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkLBaaSLoadBalancerService(sess)
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	lbaasID := parts[0]
	monitorID := parts[1]

	result, err := service.Mask("listeners.defaultPool.healthMonitor").GetLoadBalancer(sl.String(lbaasID))
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}
	for _, i := range result.Listeners {
		if monitorID == *i.DefaultPool.HealthMonitor.Uuid {
			d.Set("protocol", *i.Protocol)
			d.Set("port", *i.DefaultPool.ProtocolPort)
			d.Set("interval", *i.DefaultPool.HealthMonitor.Interval)
			d.Set("max_retries", *i.DefaultPool.HealthMonitor.MaxRetries)
			d.Set("timeout", *i.DefaultPool.HealthMonitor.Timeout)
			if i.DefaultPool.HealthMonitor.UrlPath != nil && *i.Protocol == "HTTP" {
				d.Set("url_path", *i.DefaultPool.HealthMonitor.UrlPath)
			}

			break
		}
	}

	return nil
}

func resourceIBMLbaasHealthMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	healthMonitorService := services.GetNetworkLBaaSHealthMonitorService(sess.SetRetries(0))
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	lbaasID := parts[0]
	monitorID := parts[1]

	if d.HasChange("interval") || d.HasChange("timeout") || d.HasChange("max_retries") || d.HasChange("url_path") {
		healthMonitors := make([]datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration, 0, 1)
		healthMonitor := datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration{
			BackendPort:       sl.Int(d.Get("port").(int)),
			BackendProtocol:   sl.String(d.Get("protocol").(string)),
			HealthMonitorUuid: sl.String(monitorID),
			Interval:          sl.Int(d.Get("interval").(int)),
			Timeout:           sl.Int(d.Get("timeout").(int)),
			MaxRetries:        sl.Int(d.Get("max_retries").(int)),
			UrlPath:           sl.String(d.Get("url_path").(string)),
		}

		healthMonitors = append(healthMonitors, healthMonitor)

		_, err = waitForLbaasLBActive(d, meta)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
		}

		_, err := healthMonitorService.UpdateLoadBalancerHealthMonitors(sl.String(lbaasID), healthMonitors)
		if err != nil {
			return fmt.Errorf("Error adding health monitors: %#v", err)
		}
		_, err = waitForLbaasLBActive(d, meta)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
		}
	}
	return resourceIBMLbaasHealthMonitorRead(d, meta)
}

func resourceIBMLbaasHealthMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	fmt.Println("Health monitor is destroyed only when the corresponding protocol is removed")
	d.SetId("")
	return nil
}
