// Copyright IBM Corp. 2021, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudant

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func ResourceIBMCloudantDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCloudantDatabaseCreate,
		ReadContext:   resourceIBMCloudantDatabaseRead,
		DeleteContext: resourceIBMCloudantDatabaseDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_crn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloudant Instance CRN.",
			},
			"db": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Path parameter to specify the database name.",
			},
			"partitioned": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Query parameter to specify whether to enable database partitions when creating a database.",
			},
			"shards": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The number of shards in the database. Each shard is a partition of the hash value range. You are encouraged to talk to support about appropriate values before changing this.",
			},
		},
	}
}

func resourceIBMCloudantDatabaseCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceCRN := d.Get("instance_crn").(string)
	cUrl, err := GetCloudantInstanceUrl(instanceCRN, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	cloudantClient, err := GetCloudantClientForUrl(cUrl, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	dbName := d.Get("db").(string)
	putDatabaseOptions := cloudantClient.NewPutDatabaseOptions(dbName)
	if _, ok := d.GetOk("partitioned"); ok {
		putDatabaseOptions.SetPartitioned(d.Get("partitioned").(bool))
	}
	if _, ok := d.GetOk("shards"); ok {
		putDatabaseOptions.SetQ(int64(d.Get("shards").(int)))
	}

	_, response, err := cloudantClient.PutDatabaseWithContext(context, putDatabaseOptions)
	if err != nil {
		log.Printf("[DEBUG] PutDatabaseWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("PutDatabaseWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceCRN, dbName))

	return resourceIBMCloudantDatabaseRead(context, d, meta)
}

func resourceIBMCloudantDatabaseRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	instanceCRN, dbName := strings.Join(parts[:len(parts)-1], "/"), parts[len(parts)-1]
	cUrl, err := GetCloudantInstanceUrl(instanceCRN, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	cloudantClient, err := GetCloudantClientForUrl(cUrl, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getDatabaseInformationOptions := cloudantClient.NewGetDatabaseInformationOptions(dbName)

	databaseInformation, response, err := cloudantClient.GetDatabaseInformationWithContext(context, getDatabaseInformationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDatabaseInformationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetDatabaseInformationWithContext failed %s\n%s", err, response))
	}

	d.Set("instance_crn", instanceCRN)

	if err = d.Set("db", *databaseInformation.DbName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting db: %s", err))
	}

	if err = d.Set("partitioned", databaseInformation.Props.Partitioned != nil); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting partitioned: %s", err))
	}

	if err = d.Set("shards", int(*databaseInformation.Cluster.Q)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting shards: %s", err))
	}

	return nil
}

func resourceIBMCloudantDatabaseDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	instanceCRN, dbName := strings.Join(parts[:len(parts)-1], "/"), parts[len(parts)-1]
	cUrl, err := GetCloudantInstanceUrl(instanceCRN, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	cloudantClient, err := GetCloudantClientForUrl(cUrl, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteDatabaseOptions := cloudantClient.NewDeleteDatabaseOptions(dbName)

	_, response, err := cloudantClient.DeleteDatabaseWithContext(context, deleteDatabaseOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteDatabaseWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteDatabaseWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func GetCloudantInstanceUrl(instanceCRN string, meta interface{}) (string, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return "", err
	}

	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: flex.PtrToString(instanceCRN),
	}

	instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil {
		return "", fmt.Errorf("Error retrieving resource instance: %s with resp code: %s", err, resp)
	}

	if instance.Extensions != nil {
		instanceExtensionMap := flex.Flatten(instance.Extensions)
		if instanceExtensionMap != nil {
			cloudantInstanceUrl := "https://" + instanceExtensionMap["endpoints.public"]
			cloudantInstanceUrl = conns.EnvFallBack([]string{"IBMCLOUD_CLOUDANT_API_ENDPOINT"}, cloudantInstanceUrl)
			return cloudantInstanceUrl, nil
		}
	}

	return "", fmt.Errorf("Unable to get URL for cloudant instance")
}
