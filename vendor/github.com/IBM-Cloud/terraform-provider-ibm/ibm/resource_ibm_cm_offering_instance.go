// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	inProgress = "in progress"
	failed     = "failed"
	success    = "succeeded"

	waitUntilInterval = 10 * time.Second
)

func resourceIBMCmOfferingInstance() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCmOfferingInstanceCreate,
		Read:     resourceIBMCmOfferingInstanceRead,
		Update:   resourceIBMCmOfferingInstanceUpdate,
		Delete:   resourceIBMCmOfferingInstanceDelete,
		Exists:   resourceIBMCmOfferingInstanceExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "url reference to this object.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "platform CRN for this instance.",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "the label for this instance.",
			},
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Catalog ID this instance was created from.",
			},
			"offering_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Offering ID this instance was created from.",
			},
			"kind_format": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "the format this instance has (helm, operator, ova...).",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version this instance was installed from (not version id).",
			},
			"cluster_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID.",
			},
			"cluster_region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster region (e.g., us-south).",
			},
			"cluster_namespaces": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of target namespaces to install into.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"cluster_all_namespaces": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "designate to install into all namespaces.",
			},
			"schematics_workspace_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "id of the schematics workspace, for offerings installed through schematics",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "id of the resource group",
			},
			"install_plan": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "install plan for the subscription of the operator- can be either automatic or manual. Required for operator bundles",
			},
			"channel": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "channel to target for the operator subscription. Required for operator bundles",
			},
			"wait_until_successful": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Default:          true,
				Description:      "Whether to wait until the offering instance successfully provisions, or to return when accepted",
			},
		},
	}
}

func resourceIBMCmOfferingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	rsConClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	createOfferingInstanceOptions := &catalogmanagementv1.CreateOfferingInstanceOptions{}

	schemID, isfound := os.LookupEnv("IC_SCHEMATICS_WORKSPACE_ID")
	if isfound {
		createOfferingInstanceOptions.SetSchematicsWorkspaceID(schemID)
	}
	createOfferingInstanceOptions.SetXAuthRefreshToken(rsConClient.Config.IAMRefreshToken)
	if _, ok := d.GetOk("label"); ok {
		createOfferingInstanceOptions.SetLabel(d.Get("label").(string))
	}
	if _, ok := d.GetOk("catalog_id"); ok {
		createOfferingInstanceOptions.SetCatalogID(d.Get("catalog_id").(string))
	}
	if _, ok := d.GetOk("offering_id"); ok {
		createOfferingInstanceOptions.SetOfferingID(d.Get("offering_id").(string))
	}
	if _, ok := d.GetOk("kind_format"); ok {
		createOfferingInstanceOptions.SetKindFormat(d.Get("kind_format").(string))
	}
	if _, ok := d.GetOk("version"); ok {
		createOfferingInstanceOptions.SetVersion(d.Get("version").(string))
	}
	if _, ok := d.GetOk("cluster_id"); ok {
		createOfferingInstanceOptions.SetClusterID(d.Get("cluster_id").(string))
	}
	if _, ok := d.GetOk("cluster_region"); ok {
		createOfferingInstanceOptions.SetClusterRegion(d.Get("cluster_region").(string))
	}
	if ns, ok := d.GetOk("cluster_namespaces"); ok {
		list := expandStringList(ns.([]interface{}))
		createOfferingInstanceOptions.SetClusterNamespaces(list)
	}
	if _, ok := d.GetOk("cluster_all_namespaces"); ok {
		createOfferingInstanceOptions.SetClusterAllNamespaces(d.Get("cluster_all_namespaces").(bool))
	}
	if _, ok := d.GetOk("resource_group_id"); ok {
		createOfferingInstanceOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	}
	if _, ok := d.GetOk("install_plan"); ok {
		createOfferingInstanceOptions.SetInstallPlan(d.Get("install_plan").(string))
	}
	if _, ok := d.GetOk("channel"); ok {
		createOfferingInstanceOptions.SetChannel(d.Get("channel").(string))
	}

	offeringInstance, response, err := catalogManagementClient.CreateOfferingInstance(createOfferingInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateOfferingInstance failed %s\n%s", err, response)
		return err
	}

	d.SetId(*offeringInstance.ID)

	if d.Get("wait_until_successful").(bool) {
		if _, err = waitUntilSuccess(d, meta); err != nil {
			log.Print(err)
			return err
		}
	}

	log.Printf("LOG2 Service version instance of type %q was created on cluster %q", *createOfferingInstanceOptions.KindFormat, *createOfferingInstanceOptions.ClusterID)

	return resourceIBMCmOfferingInstanceRead(d, meta)
}

func waitUntilSuccess(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return nil, err
	}
	getOfferingInstanceOptions := &catalogmanagementv1.GetOfferingInstanceOptions{}

	getOfferingInstanceOptions.SetInstanceIdentifier(d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{inProgress},
		Target:  []string{success},
		Refresh: func() (interface{}, string, error) {
			offeringInstance, _, err := catalogManagementClient.GetOfferingInstance(getOfferingInstanceOptions)
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving offering instance: %s", err)
			}

			return offeringInstance, *offeringInstance.LastOperation.State, nil
		},
		Delay:      waitUntilInterval * 2,
		MinTimeout: waitUntilInterval,
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}

	return stateConf.WaitForState()
}

func resourceIBMCmOfferingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	getOfferingInstanceOptions := &catalogmanagementv1.GetOfferingInstanceOptions{}

	getOfferingInstanceOptions.SetInstanceIdentifier(d.Id())

	offeringInstance, response, err := catalogManagementClient.GetOfferingInstance(getOfferingInstanceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetOfferingInstance failed %s\n%s", err, response)
		return err
	}

	if err = d.Set("url", offeringInstance.URL); err != nil {
		return fmt.Errorf("Error setting url: %s", err)
	}
	if err = d.Set("crn", offeringInstance.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("label", offeringInstance.Label); err != nil {
		return fmt.Errorf("Error setting label: %s", err)
	}
	if err = d.Set("catalog_id", offeringInstance.CatalogID); err != nil {
		return fmt.Errorf("Error setting catalog_id: %s", err)
	}
	if err = d.Set("offering_id", offeringInstance.OfferingID); err != nil {
		return fmt.Errorf("Error setting offering_id: %s", err)
	}
	if err = d.Set("kind_format", offeringInstance.KindFormat); err != nil {
		return fmt.Errorf("Error setting kind_format: %s", err)
	}
	if err = d.Set("version", offeringInstance.Version); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}
	if err = d.Set("cluster_id", offeringInstance.ClusterID); err != nil {
		return fmt.Errorf("Error setting cluster_id: %s", err)
	}
	if err = d.Set("cluster_region", offeringInstance.ClusterRegion); err != nil {
		return fmt.Errorf("Error setting cluster_region: %s", err)
	}
	if offeringInstance.ClusterNamespaces != nil {
		if err = d.Set("cluster_namespaces", offeringInstance.ClusterNamespaces); err != nil {
			return fmt.Errorf("Error setting cluster_namespaces: %s", err)
		}
	}
	if err = d.Set("cluster_all_namespaces", offeringInstance.ClusterAllNamespaces); err != nil {
		return fmt.Errorf("Error setting cluster_all_namespaces: %s", err)
	}
	if err = d.Set("schematics_workspace_id", offeringInstance.SchematicsWorkspaceID); err != nil {
		return fmt.Errorf("Error setting schematics_workspace_id: %s", err)
	}
	if err = d.Set("install_plan", offeringInstance.InstallPlan); err != nil {
		return fmt.Errorf("Error setting install_plan: %s", err)
	}
	if err = d.Set("channel", offeringInstance.Channel); err != nil {
		return fmt.Errorf("Error setting channel: %s", err)
	}

	return nil
}

func resourceIBMCmOfferingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	getOfferingInstanceOptions := &catalogmanagementv1.GetOfferingInstanceOptions{}

	getOfferingInstanceOptions.SetInstanceIdentifier(d.Id())

	offeringInstance, response, err := catalogManagementClient.GetOfferingInstance(getOfferingInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] Failed to retrieve rev %s\n%s", err, response)
		return err
	}

	rsConClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	putOfferingInstanceOptions := &catalogmanagementv1.PutOfferingInstanceOptions{}

	schemID, isfound := os.LookupEnv("IC_SCHEMATICS_WORKSPACE_ID")
	if isfound {
		putOfferingInstanceOptions.SetSchematicsWorkspaceID(schemID)
	}

	putOfferingInstanceOptions.SetInstanceIdentifier(d.Id())
	putOfferingInstanceOptions.SetID(d.Id())
	putOfferingInstanceOptions.SetXAuthRefreshToken(rsConClient.Config.IAMRefreshToken)
	putOfferingInstanceOptions.SetRev(*offeringInstance.Rev)
	if _, ok := d.GetOk("label"); ok {
		putOfferingInstanceOptions.SetLabel(d.Get("label").(string))
	}
	if _, ok := d.GetOk("catalog_id"); ok {
		putOfferingInstanceOptions.SetCatalogID(d.Get("catalog_id").(string))
	}
	if _, ok := d.GetOk("offering_id"); ok {
		putOfferingInstanceOptions.SetOfferingID(d.Get("offering_id").(string))
	}
	if _, ok := d.GetOk("kind_format"); ok {
		putOfferingInstanceOptions.SetKindFormat(d.Get("kind_format").(string))
	}
	if _, ok := d.GetOk("version"); ok {
		putOfferingInstanceOptions.SetVersion(d.Get("version").(string))
	}
	if _, ok := d.GetOk("cluster_id"); ok {
		putOfferingInstanceOptions.SetClusterID(d.Get("cluster_id").(string))
	}
	if _, ok := d.GetOk("cluster_region"); ok {
		putOfferingInstanceOptions.SetClusterRegion(d.Get("cluster_region").(string))
	}
	if ns, ok := d.GetOk("cluster_namespaces"); ok {
		list := expandStringList(ns.([]interface{}))
		putOfferingInstanceOptions.SetClusterNamespaces(list)
	}
	if _, ok := d.GetOk("cluster_all_namespaces"); ok {
		putOfferingInstanceOptions.SetClusterAllNamespaces(d.Get("cluster_all_namespaces").(bool))
	}
	if _, ok := d.GetOk("resource_group_id"); ok {
		putOfferingInstanceOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	}
	if _, ok := d.GetOk("install_plan"); ok {
		putOfferingInstanceOptions.SetInstallPlan(d.Get("install_plan").(string))
	}
	if _, ok := d.GetOk("channel"); ok {
		putOfferingInstanceOptions.SetChannel(d.Get("channel").(string))
	}

	_, response, err = catalogManagementClient.PutOfferingInstance(putOfferingInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] PutOfferingInstance failed %s\n%s", err, response)
		return err
	}

	return resourceIBMCmOfferingInstanceRead(d, meta)
}

func resourceIBMCmOfferingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}
	rsConClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	deleteOfferingInstanceOptions := &catalogmanagementv1.DeleteOfferingInstanceOptions{}

	deleteOfferingInstanceOptions.SetInstanceIdentifier(d.Id())
	deleteOfferingInstanceOptions.SetXAuthRefreshToken(rsConClient.Config.IAMRefreshToken)

	response, err := catalogManagementClient.DeleteOfferingInstance(deleteOfferingInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteOfferingInstance failed %s\n%s", err, response)
		return err
	}

	d.SetId("")

	return nil
}

func resourceIBMCmOfferingInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return false, err
	}

	getOfferingInstanceOptions := &catalogmanagementv1.GetOfferingInstanceOptions{}

	getOfferingInstanceOptions.SetInstanceIdentifier(d.Id())

	offeringInstance, response, err := catalogManagementClient.GetOfferingInstance(getOfferingInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] GetOfferingInstance failed %s\n%s", err, response)
		return false, err
	}

	return *offeringInstance.ID == d.Id(), nil
}
