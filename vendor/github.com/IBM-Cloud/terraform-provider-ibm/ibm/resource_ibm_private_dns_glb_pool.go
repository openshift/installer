// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v3/core"
	dns "github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmDNSGlbPool                         = "ibm_dns_glb_pool"
	pdnsGlbPoolID                         = "pool_id"
	pdnsGlbPoolName                       = "name"
	pdnsGlbPoolDescription                = "description"
	pdnsGlbPoolEnabled                    = "enabled"
	pdnsGlbPoolHealth                     = "health"
	pdnsGlbPoolHealthyOriginsThreshold    = "healthy_origins_threshold"
	pdnsGlbPoolOrigins                    = "origins"
	pdnsGlbPoolOriginsName                = "name"
	pdnsGlbPoolOriginsDescription         = "description"
	pdnsGlbPoolOriginsAddress             = "address"
	pdnsGlbPoolOriginsEnabled             = "enabled"
	pdnsGlbPoolOriginsHealth              = "health"
	pdnsGlbPoolOriginsHealthFailureReason = "health_failure_reason"
	pdnsGlbPoolMonitor                    = "monitor"
	pdnsGlbPoolChannel                    = "notification_channel"
	pdnsGlbPoolRegion                     = "healthcheck_region"
	pdnsGlbPoolSubnet                     = "healthcheck_subnets"
	pdnsGlbPoolCreatedOn                  = "created_on"
	pdnsGlbPoolModifiedOn                 = "modified_on"
	pdnsGlbPoolDeletePending              = "deleting"
	pdnsGlbPoolDeleted                    = "deleted"
)

func resourceIBMPrivateDNSGLBPool() *schema.Resource {
	return &schema.Resource{

		Create:   resourceIBMPrivateDNSGLBPoolCreate,
		Read:     resourceIBMPrivateDNSGLBPoolRead,
		Update:   resourceIBMPrivateDNSGLBPoolUpdate,
		Delete:   resourceIBMPrivateDNSGLBPoolDelete,
		Exists:   resourceIBMPrivateDNSGLBPoolExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance Id",
			},
			pdnsGlbPoolID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pool Id",
			},
			pdnsGlbPoolName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of a service instance.",
			},
			pdnsGlbPoolDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Descriptive text of the load balancer pool",
			},
			pdnsGlbPoolEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the load balancer pool is enabled",
			},
			pdnsGlbPoolHealth: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the load balancer pool is enabled",
			},
			pdnsGlbPoolHealthyOriginsThreshold: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The minimum number of origins that must be healthy for this pool to serve traffic",
			},
			pdnsGlbPoolOrigins: {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Origins info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsGlbPoolOriginsName: {
							Type:        schema.TypeString,
							Description: "The name of the origin server.",
							Required:    true,
						},
						pdnsGlbPoolOriginsAddress: {
							Type:        schema.TypeString,
							Description: "The address of the origin server. It can be a hostname or an IP address.",
							Required:    true,
						},
						pdnsGlbPoolOriginsEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether the origin server is enabled.",
							Required:    true,
						},
						pdnsGlbPoolOriginsDescription: {
							Type:        schema.TypeString,
							Description: "Description of the origin server.",
							Optional:    true,
						},
						pdnsGlbPoolOriginsHealth: {
							Type:        schema.TypeBool,
							Description: "Whether the health is `true` or `false`.",
							Computed:    true,
						},
						pdnsGlbPoolOriginsHealthFailureReason: {
							Type:        schema.TypeString,
							Description: "The Reason for health check failure",
							Computed:    true,
						},
					},
				},
			},
			pdnsGlbPoolMonitor: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the load balancer monitor to be associated to this pool",
			},
			pdnsGlbPoolChannel: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The notification channel,It is a webhook url",
			},
			pdnsGlbPoolRegion: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator(ibmDNSGlbPool, pdnsGlbPoolRegion),
				Description:  "Health check region of VSIs",
			},
			pdnsGlbPoolSubnet: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Health check subnet crn of VSIs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			pdnsGlbPoolCreatedOn: {
				Type:        schema.TypeString,
				Description: "The time when a load balancer pool is created.",
				Computed:    true,
			},
			pdnsGlbPoolModifiedOn: {
				Type:        schema.TypeString,
				Description: "The recent time when a load balancer pool is modified.",
				Computed:    true,
			},
		},
	}
}

func resourceIBMPrivateDNSGLBPoolValidator() *ResourceValidator {
	regions := "us-south,us-east,eu-gb,eu-du,au-syd,jp-tok"

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 pdnsGlbPoolRegion,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              regions})
	dnsPoolValidator := ResourceValidator{ResourceName: ibmDNSGlbPool, Schema: validateSchema}
	return &dnsPoolValidator
}

func resourceIBMPrivateDNSGLBPoolCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	instanceID := d.Get(pdnsInstanceID).(string)
	CreatePoolOptions := sess.NewCreatePoolOptions(instanceID)

	poolname := d.Get(pdnsGlbPoolName).(string)
	CreatePoolOptions.SetName(poolname)

	if description, ok := d.GetOk(pdnsGlbPoolDescription); ok {
		CreatePoolOptions.SetDescription(description.(string))
	}
	if enable, ok := d.GetOk(pdnsGlbPoolEnabled); ok {
		CreatePoolOptions.SetEnabled(enable.(bool))
	}
	if threshold, ok := d.GetOk(pdnsGlbPoolHealthyOriginsThreshold); ok {
		CreatePoolOptions.SetHealthyOriginsThreshold(int64(threshold.(int)))
	}
	if monitor, ok := d.GetOk(pdnsGlbPoolMonitor); ok {
		monitorID, _, _ := convertTftoCisTwoVar(monitor.(string))
		CreatePoolOptions.SetMonitor(monitorID)
	}
	if chanel, ok := d.GetOk(pdnsGlbPoolChannel); ok {
		CreatePoolOptions.SetNotificationChannel(chanel.(string))
	}
	if region, ok := d.GetOk(pdnsGlbPoolRegion); ok {
		CreatePoolOptions.SetHealthcheckRegion(region.(string))
	}
	if subnets, ok := d.GetOk(pdnsGlbPoolSubnet); ok {
		CreatePoolOptions.SetHealthcheckSubnets(expandStringList(subnets.([]interface{})))
	}

	poolorigins := d.Get(pdnsGlbPoolOrigins).(*schema.Set)
	CreatePoolOptions.SetOrigins(expandPDNSGlbPoolOrigins(poolorigins))

	result, resp, err := sess.CreatePool(CreatePoolOptions)
	if err != nil {
		log.Printf("create global load balancer pool failed %s", resp)
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceID, *result.ID))

	return resourceIBMPrivateDNSGLBPoolRead(d, meta)
}

func resourceIBMPrivateDNSGLBPoolRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	idset := strings.Split(d.Id(), "/")

	getPoolOptions := sess.NewGetPoolOptions(idset[0], idset[1])
	presponse, resp, err := sess.GetPool(getPoolOptions)
	if err != nil {
		return fmt.Errorf("Error fetching pdns GLB Pool:%s\n%s", err, resp)
	}

	response := *presponse
	d.Set(pdnsGlbPoolName, response.Name)
	d.Set(pdnsGlbPoolID, response.ID)
	d.Set(pdnsInstanceID, idset[0])
	d.Set(pdnsGlbPoolDescription, response.Description)
	d.Set(pdnsGlbPoolEnabled, response.Enabled)
	d.Set(pdnsGlbPoolHealth, response.Health)
	d.Set(pdnsGlbPoolHealthyOriginsThreshold, response.HealthyOriginsThreshold)
	d.Set(pdnsGlbPoolMonitor, response.Monitor)
	d.Set(pdnsGlbPoolChannel, response.NotificationChannel)
	d.Set(pdnsGlbPoolRegion, response.HealthcheckRegion)
	d.Set(pdnsGlbPoolSubnet, response.HealthcheckSubnets)
	d.Set(pdnsGlbPoolCreatedOn, response.CreatedOn)
	d.Set(pdnsGlbPoolModifiedOn, response.ModifiedOn)
	d.Set(pdnsGlbPoolOrigins, flattenPDNSGlbPoolOrigins(response.Origins))

	return nil
}

func flattenPDNSGlbPoolOrigins(list []dns.Origin) []map[string]interface{} {
	origins := []map[string]interface{}{}
	for _, origin := range list {
		l := map[string]interface{}{
			pdnsGlbPoolOriginsName:                *origin.Name,
			pdnsGlbPoolOriginsAddress:             *origin.Address,
			pdnsGlbPoolOriginsEnabled:             *origin.Enabled,
			pdnsGlbPoolOriginsDescription:         *origin.Description,
			pdnsGlbPoolOriginsHealth:              *origin.Health,
			pdnsGlbPoolOriginsHealthFailureReason: *origin.HealthFailureReason,
		}
		origins = append(origins, l)
	}
	return origins
}

func resourceIBMPrivateDNSGLBPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")
	updatePoolOptions := sess.NewUpdatePoolOptions(idset[0], idset[1])

	if d.HasChange(pdnsGlbPoolName) ||
		d.HasChange(pdnsGlbPoolDescription) ||
		d.HasChange(pdnsGlbPoolEnabled) ||
		d.HasChange(pdnsGlbPoolHealthyOriginsThreshold) ||
		d.HasChange(pdnsGlbPoolMonitor) ||
		d.HasChange(pdnsGlbPoolChannel) ||
		d.HasChange(pdnsGlbPoolRegion) ||
		d.HasChange(pdnsGlbPoolOrigins) ||
		d.HasChange(pdnsGlbPoolSubnet) {
		if Mname, ok := d.GetOk(pdnsGlbPoolName); ok {
			updatePoolOptions.SetName(Mname.(string))
		}
		if description, ok := d.GetOk(pdnsGlbPoolDescription); ok {
			updatePoolOptions.SetDescription(description.(string))
		}
		if enable, ok := d.GetOk(pdnsGlbPoolEnabled); ok {
			updatePoolOptions.SetEnabled(enable.(bool))
		}
		if threshold, ok := d.GetOk(pdnsGlbPoolHealthyOriginsThreshold); ok {
			updatePoolOptions.SetHealthyOriginsThreshold(int64(threshold.(int)))
		}
		if monitor, ok := d.GetOk(pdnsGlbPoolMonitor); ok {
			monitorID, _, _ := convertTftoCisTwoVar(monitor.(string))
			updatePoolOptions.SetMonitor(monitorID)
		}
		if chanel, ok := d.GetOk(pdnsGlbPoolChannel); ok {
			updatePoolOptions.SetNotificationChannel(chanel.(string))
		}
		if region, ok := d.GetOk(pdnsGlbPoolRegion); ok {
			updatePoolOptions.SetHealthcheckRegion(region.(string))
		}
		if _, ok := d.GetOk(pdnsGlbPoolSubnet); ok {
			updatePoolOptions.SetHealthcheckSubnets(expandStringList(d.Get(pdnsGlbPoolSubnet).([]interface{})))
		}
		if _, ok := d.GetOk(pdnsGlbPoolOrigins); ok {
			poolorigins := d.Get(pdnsGlbPoolOrigins).(*schema.Set)
			updatePoolOptions.SetOrigins(expandPDNSGlbPoolOrigins(poolorigins))

		}
		_, detail, err := sess.UpdatePool(updatePoolOptions)
		if err != nil {
			return fmt.Errorf("Error updating pdns GLB Pool:%s\n%s", err, detail)
		}
	}

	return resourceIBMPrivateDNSGLBPoolRead(d, meta)
}

func resourceIBMPrivateDNSGLBPoolDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")
	DeletePoolOptions := sess.NewDeletePoolOptions(idset[0], idset[1])
	response, err := sess.DeletePool(DeletePoolOptions)
	if err != nil {
		return fmt.Errorf("Error deleting pdns GLB Pool:%s\n%s", err, response)
	}
	_, err = waitForPDNSGlbPoolDelete(d, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceIBMPrivateDNSGLBPoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return false, err
	}

	idset := strings.Split(d.Id(), "/")

	getPoolOptions := sess.NewGetPoolOptions(idset[0], idset[1])
	response, detail, err := sess.GetPool(getPoolOptions)
	if err != nil {
		if response != nil && detail != nil && detail.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func expandPDNSGlbPoolOrigins(originsList *schema.Set) (origins []dns.OriginInput) {
	for _, iface := range originsList.List() {
		orig := iface.(map[string]interface{})
		origin := dns.OriginInput{
			Name:        core.StringPtr(orig[pdnsGlbPoolOriginsName].(string)),
			Address:     core.StringPtr(orig[pdnsGlbPoolOriginsAddress].(string)),
			Enabled:     core.BoolPtr(orig[pdnsGlbPoolOriginsEnabled].(bool)),
			Description: core.StringPtr(orig[pdnsGlbPoolOriginsDescription].(string)),
		}
		origins = append(origins, origin)
	}
	return
}

func waitForPDNSGlbPoolDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	cisClient, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return nil, err
	}
	idset := strings.Split(d.Id(), "/")
	getPoolOptions := cisClient.NewGetPoolOptions(idset[0], idset[1])
	stateConf := &resource.StateChangeConf{
		Pending: []string{pdnsGlbPoolDeletePending},
		Target:  []string{pdnsGlbPoolDeleted},
		Refresh: func() (interface{}, string, error) {
			_, detail, err := cisClient.GetPool(getPoolOptions)
			if err != nil {
				if detail != nil && detail.StatusCode == 404 {
					return detail, clusterDeleted, nil
				}
				return nil, "", err
			}
			return detail, clusterDeletePending, nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		MinTimeout:   10 * time.Second,
		PollInterval: 60 * time.Second,
	}

	return stateConf.WaitForState()
}
