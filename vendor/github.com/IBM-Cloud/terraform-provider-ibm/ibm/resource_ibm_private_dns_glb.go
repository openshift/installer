// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsGLBName             = "name"
	pdnsGLBID               = "glb_id"
	pdnsGLBDescription      = "description"
	pdnsGLBEnabled          = "enabled"
	pdnsGLBTTL              = "ttl"
	pdnsGLBHealth           = "health"
	pdnsGLBFallbackPool     = "fallback_pool"
	pdnsGLBDefaultPool      = "default_pools"
	pdnsGLBAZPools          = "az_pools"
	pdnsGLBAvailabilityZone = "availability_zone"
	pdnsGLBAZPoolsPools     = "pools"
	pdnsGLBCreatedOn        = "created_on"
	pdnsGLBModifiedOn       = "modified_on"
	pdnsGLBDeleting         = "deleting"
	pdnsGLBDeleted          = "done"
)

func resourceIBMPrivateDNSGLB() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPrivateDNSGLBCreate,
		Read:     resourceIBMPrivateDNSGLBRead,
		Update:   resourceIBMPrivateDNSGLBUpdate,
		Delete:   resourceIBMPrivateDNSGLBDelete,
		Exists:   resourceIBMPrivateDNSGLBExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			pdnsGLBID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer Id",
			},
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The GUID of the private DNS.",
			},
			pdnsZoneID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone Id",
			},
			pdnsGLBName: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of the load balancer",
				DiffSuppressFunc: suppressPDNSGlbNameDiff,
			},
			pdnsGLBDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Descriptive text of the load balancer",
			},
			pdnsGLBEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the load balancer is enabled",
			},
			pdnsGLBTTL: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Time to live in second",
			},
			pdnsGLBHealth: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Healthy state of the load balancer.",
			},
			pdnsGLBFallbackPool: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The pool ID to use when all other pools are detected as unhealthy",
			},
			pdnsGLBDefaultPool: {
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of pool IDs ordered by their failover priority",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			pdnsGLBAZPools: {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Map availability zones to pool ID's.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsGLBAvailabilityZone: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Availability zone.",
						},

						pdnsGLBAZPoolsPools: {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of load balancer pools",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			pdnsGLBCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GLB Load Balancer creation date",
			},
			pdnsGLBModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GLB Load Balancer Modification date",
			},
		},
	}
}

func resourceIBMPrivateDNSGLBCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	zoneID := d.Get(pdnsZoneID).(string)
	createlbOptions := sess.NewCreateLoadBalancerOptions(instanceID, zoneID)

	lbname := d.Get(pdnsGLBName).(string)
	createlbOptions.SetName(lbname)
	createlbOptions.SetFallbackPool(d.Get(pdnsGLBFallbackPool).(string))
	createlbOptions.SetDefaultPools(expandStringList(d.Get(pdnsGLBDefaultPool).([]interface{})))

	if description, ok := d.GetOk(pdnsGLBDescription); ok {
		createlbOptions.SetDescription(description.(string))
	}
	if enable, ok := d.GetOkExists(pdnsGLBEnabled); ok {
		createlbOptions.SetEnabled(enable.(bool))
	}
	if ttl, ok := d.GetOk(pdnsGLBTTL); ok {
		createlbOptions.SetTTL(int64(ttl.(int)))
	}

	if AZpools, ok := d.GetOk(pdnsGLBAZPools); ok {
		expandedAzpools, err := expandPDNSGlbAZPools(AZpools)
		if err != nil {
			return err
		}
		createlbOptions.SetAzPools(expandedAzpools)
	}

	result, resp, err := sess.CreateLoadBalancer(createlbOptions)
	if err != nil {
		log.Printf("create global load balancer failed %s", resp)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, zoneID, *result.ID))
	return resourceIBMPrivateDNSGLBRead(d, meta)
}

func resourceIBMPrivateDNSGLBRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	idset := strings.Split(d.Id(), "/")

	getlbOptions := sess.NewGetLoadBalancerOptions(idset[0], idset[1], idset[2])
	presponse, resp, err := sess.GetLoadBalancer(getlbOptions)
	if err != nil {
		return fmt.Errorf("Error fetching pdns GLB :%s\n%s", err, resp)
	}

	response := *presponse
	d.Set(pdnsInstanceID, idset[0])
	d.Set(pdnsZoneID, idset[1])
	d.Set(pdnsGLBName, response.Name)
	d.Set(pdnsGLBID, response.ID)
	d.Set(pdnsGLBDescription, response.Description)
	d.Set(pdnsGLBEnabled, response.Enabled)
	d.Set(pdnsGLBTTL, response.TTL)
	d.Set(pdnsGLBHealth, response.Health)
	d.Set(pdnsGLBFallbackPool, response.FallbackPool)
	d.Set(pdnsGLBDefaultPool, response.DefaultPools)
	d.Set(pdnsGLBCreatedOn, response.CreatedOn)
	d.Set(pdnsGLBModifiedOn, response.ModifiedOn)
	d.Set(pdnsGLBAZPools, flattenPDNSGlbAZpool(response.AzPools))
	return nil
}

func resourceIBMPrivateDNSGLBUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")

	updatelbOptions := sess.NewUpdateLoadBalancerOptions(idset[0], idset[1], idset[2])

	if d.HasChange(pdnsGLBName) ||
		d.HasChange(pdnsGLBDescription) ||
		d.HasChange(pdnsGLBEnabled) ||
		d.HasChange(pdnsGLBTTL) ||
		d.HasChange(pdnsGLBFallbackPool) ||
		d.HasChange(pdnsGLBDefaultPool) ||
		d.HasChange(pdnsGLBAZPools) {

		updatelbOptions.SetName(d.Get(pdnsGLBName).(string))
		updatelbOptions.SetFallbackPool(d.Get(pdnsGLBFallbackPool).(string))
		updatelbOptions.SetDefaultPools(expandStringList(d.Get(pdnsGLBDefaultPool).([]interface{})))

		if description, ok := d.GetOk(pdnsGLBDescription); ok {
			updatelbOptions.SetDescription(description.(string))
		}
		if enable, ok := d.GetOkExists(pdnsGLBEnabled); ok {
			updatelbOptions.SetEnabled(enable.(bool))
		}
		if ttl, ok := d.GetOk(pdnsGLBTTL); ok {
			updatelbOptions.SetTTL(int64(ttl.(int)))
		}

		if AZpools, ok := d.GetOk(pdnsGLBAZPools); ok {
			expandedAzpools, err := expandPDNSGlbAZPools(AZpools)
			if err != nil {
				return err
			}
			updatelbOptions.SetAzPools(expandedAzpools)
		}

		_, detail, err := sess.UpdateLoadBalancer(updatelbOptions)
		if err != nil {
			return fmt.Errorf("Error updating pdns GLB :%s\n%s", err, detail)
		}
	}

	return resourceIBMPrivateDNSGLBRead(d, meta)
}

func resourceIBMPrivateDNSGLBDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")
	deletelbOptions := sess.NewDeleteLoadBalancerOptions(idset[0], idset[1], idset[2])
	response, err := sess.DeleteLoadBalancer(deletelbOptions)
	if err != nil {
		return fmt.Errorf("Error deleting pdns GLB :%s\n%s", err, response)
	}
	_, err = isWaitForLoadBalancerDeleted(sess, d, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	return nil
}

func resourceIBMPrivateDNSGLBExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return false, err
	}
	idset := strings.Split(d.Id(), "/")
	getlbOptions := sess.NewGetLoadBalancerOptions(idset[0], idset[1], idset[2])
	_, detail, err := sess.GetLoadBalancer(getlbOptions)
	if err != nil {
		if detail != nil && detail.StatusCode == 404 {
			log.Printf("Get GLB failed with status code 404: %v", detail)
			return false, nil
		}
		log.Printf("Get GLB failed: %v", detail)
		return false, err
	}
	return true, nil
}

func expandPDNSGlbAZPools(azpool interface{}) ([]dnssvcsv1.LoadBalancerAzPoolsItem, error) {
	azpools := azpool.(*schema.Set).List()
	expandAZpools := make([]dnssvcsv1.LoadBalancerAzPoolsItem, 0)
	for _, v := range azpools {
		locationConfig := v.(map[string]interface{})
		avzone := locationConfig[pdnsGLBAvailabilityZone].(string)
		pools := expandStringList(locationConfig[pdnsGLBAZPoolsPools].([]interface{}))
		aZItem := dnssvcsv1.LoadBalancerAzPoolsItem{
			AvailabilityZone: &avzone,
			Pools:            pools,
		}
		expandAZpools = append(expandAZpools, aZItem)
	}
	return expandAZpools, nil
}

func flattenPDNSGlbAZpool(azpool []dnssvcsv1.LoadBalancerAzPoolsItem) interface{} {
	flattened := make([]interface{}, 0)
	for _, v := range azpool {
		cfg := map[string]interface{}{
			pdnsGLBAvailabilityZone: *v.AvailabilityZone,
			pdnsGLBAZPoolsPools:     flattenStringList(v.Pools),
		}
		flattened = append(flattened, cfg)
	}
	return flattened
}

func suppressPDNSGlbNameDiff(k, old, new string, d *schema.ResourceData) bool {
	// PDNS GLB concantenates name with domain. So just check name is the same
	if strings.SplitN(old, ".", 2)[0] == strings.SplitN(new, ".", 2)[0] {
		return true
	}
	return false
}

func isWaitForLoadBalancerDeleted(LoadBalancer *dnssvcsv1.DnsSvcsV1, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {
	idset := strings.Split(d.Id(), "/")
	log.Printf("Waiting for PDNS GLB (%s) to be deleted.", idset[2])
	stateConf := &resource.StateChangeConf{
		Pending:    []string{pdnsGLBDeleting},
		Target:     []string{pdnsGLBDeleted},
		Refresh:    isVLoadBalancerDeleteRefreshFunc(LoadBalancer, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isVLoadBalancerDeleteRefreshFunc(LoadBalancer *dnssvcsv1.DnsSvcsV1, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		idset := strings.Split(d.Id(), "/")
		getlbOptions := LoadBalancer.NewGetLoadBalancerOptions(idset[0], idset[1], idset[2])
		_, response, err := LoadBalancer.GetLoadBalancer(getlbOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return "", pdnsGLBDeleted, nil
			}
			return "", "", fmt.Errorf("Error Getting PDNS Load Balancer : %s\n%s", err, response)
		}
		return LoadBalancer, pdnsGLBDeleting, err

	}
}
