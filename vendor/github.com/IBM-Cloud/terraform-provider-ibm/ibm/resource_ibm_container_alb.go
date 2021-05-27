// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
)

func resourceIBMContainerALB() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerALBCreate,
		Read:     resourceIBMContainerALBRead,
		Update:   resourceIBMContainerALBUpdate,
		Delete:   resourceIBMContainerALBDelete,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"alb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ALB ID",
			},
			"alb_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ALB type",
			},
			"cluster": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster id",
			},
			"user_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "IP assigned by the user",
			},
			"enable": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"disable_deployment"},
				Description:   "set to true if ALB needs to be enabled",
			},
			"disable_deployment": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"enable"},
				Description:   "Set to true if ALB needs to be disabled",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ALB name",
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ALB zone",
			},
			"region": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This field is deprecated",
			},
		},
	}
}

func resourceIBMContainerALBCreate(d *schema.ResourceData, meta interface{}) error {
	albClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	var userIP string
	var enable, disableDeployment bool
	albID := d.Get("alb_id").(string)
	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
	} else if v, ok := d.GetOkExists("disable_deployment"); ok {
		disableDeployment = v.(bool)
	} else {
		return fmt.Errorf("Provide either `enable` or `disable_deployment`")
	}

	numOfInstances := "2"
	if v, ok := d.GetOk("user_ip"); ok {
		userIP = v.(string)
	}
	params := v1.ALBConfig{
		ALBID:          albID,
		Enable:         enable,
		NumOfInstances: numOfInstances,
	}
	if userIP != "" {
		params.ALBIP = userIP
	}

	_, err = waitForClusterAvailable(d, meta, albID)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for cluster resources availabilty (%s) : %s", d.Id(), err)
	}

	albAPI := albClient.Albs()
	targetEnv, err := getAlbTargetHeader(d, meta)
	if err != nil {
		return err
	}
	err = albAPI.ConfigureALB(albID, params, disableDeployment, targetEnv)
	if err != nil {
		return err
	}
	d.SetId(albID)
	_, err = waitForContainerALB(d, meta, albID, schema.TimeoutCreate, enable, disableDeployment)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for create resource alb (%s) : %s", d.Id(), err)
	}

	return resourceIBMContainerALBRead(d, meta)
}

func resourceIBMContainerALBRead(d *schema.ResourceData, meta interface{}) error {
	albClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	albID := d.Id()

	albAPI := albClient.Albs()
	targetEnv, err := getAlbTargetHeader(d, meta)
	if err != nil {
		return err
	}
	albConfig, err := albAPI.GetALB(albID, targetEnv)
	if err != nil {
		return err
	}

	d.Set("alb_type", albConfig.ALBType)
	d.Set("cluster", albConfig.ClusterID)
	d.Set("name", albConfig.Name)
	d.Set("enable", albConfig.Enable)
	d.Set("disable_deployment", albConfig.DisableDeployment)
	d.Set("replicas", albConfig.NumOfInstances)
	d.Set("resize", albConfig.Resize)
	d.Set("user_ip", albConfig.ALBIP)
	d.Set("zone", albConfig.Zone)
	return nil
}

func resourceIBMContainerALBUpdate(d *schema.ResourceData, meta interface{}) error {
	albClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	albAPI := albClient.Albs()

	if d.HasChange("enable") {
		enable := d.Get("enable").(bool)
		disableDeployment := d.Get("disable_deployment").(bool)
		albID := d.Id()
		params := v1.ALBConfig{
			ALBID:  albID,
			Enable: enable,
		}

		targetEnv, err := getAlbTargetHeader(d, meta)
		if err != nil {
			return err
		}

		_, err = waitForClusterAvailable(d, meta, albID)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for cluster resources availabilty (%s) : %s", d.Id(), err)
		}

		err = albAPI.ConfigureALB(albID, params, disableDeployment, targetEnv)
		if err != nil {
			return err
		}
		_, err = waitForContainerALB(d, meta, albID, schema.TimeoutUpdate, enable, disableDeployment)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for updating resource alb (%s) : %s", d.Id(), err)
		}

	}
	return resourceIBMContainerALBRead(d, meta)
}

func waitForContainerALB(d *schema.ResourceData, meta interface{}, albID, timeout string, enable, disableDeployment bool) (interface{}, error) {
	albClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"active"},
		Refresh: func() (interface{}, string, error) {
			targetEnv, err := getAlbTargetHeader(d, meta)
			if err != nil {
				return nil, "", err
			}
			alb, err := albClient.Albs().GetALB(albID, targetEnv)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The resource alb %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", err
			}
			if enable {
				if alb.Enable == false {
					return alb, "pending", nil
				}
			} else if disableDeployment {
				if alb.Enable == true {
					return alb, "pending", nil
				}
			}
			return alb, "active", nil
		},
		Timeout:    d.Timeout(timeout),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIBMContainerALBDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}

// WaitForWorkerAvailable Waits for worker creation
func waitForClusterAvailable(d *schema.ResourceData, meta interface{}, albID string) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}

	target, err := getAlbTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}

	albConfig, err := csClient.Albs().GetALB(albID, target)
	if err != nil {
		return nil, err
	}

	ClusterID := albConfig.ClusterID

	log.Printf("Waiting for worker of the cluster (%s) wokers to be available.", ClusterID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", workerProvisioning},
		Target:     []string{workerNormal},
		Refresh:    workerStateRefreshFunc(csClient.Workers(), ClusterID, target),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func getAlbTargetHeader(d *schema.ResourceData, meta interface{}) (v1.ClusterTargetHeader, error) {
	var region string
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return v1.ClusterTargetHeader{}, err
	}

	if region == "" {
		region = sess.Config.Region
	}

	targetEnv := v1.ClusterTargetHeader{
		Region: region,
	}

	return targetEnv, nil
}
