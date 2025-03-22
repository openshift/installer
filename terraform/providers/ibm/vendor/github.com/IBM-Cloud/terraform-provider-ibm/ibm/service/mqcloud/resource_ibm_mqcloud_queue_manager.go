// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

package mqcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func ResourceIbmMqcloudQueueManager() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmMqcloudQueueManagerCreate,
		ReadContext:   resourceIbmMqcloudQueueManagerRead,
		UpdateContext: resourceIbmMqcloudQueueManagerUpdate,
		DeleteContext: resourceIbmMqcloudQueueManagerDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_queue_manager", "service_instance_guid"),
				Description:  "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_queue_manager", "name"),
				Description:  "A queue manager name conforming to MQ restrictions.",
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_queue_manager", "display_name"),
				Description:  "A displayable name for the queue manager - limited only in length.",
			},
			"location": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_queue_manager", "location"),
				Description:  "The locations in which the queue manager could be deployed.",
			},
			"size": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_queue_manager", "size"),
				Description:  "The queue manager sizes of deployment available.",
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_queue_manager", "version"),
				Description:  "The MQ version of the queue manager.",
			},
			"status_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A reference uri to get deployment status of the queue manager.",
			},
			"web_console_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url through which to access the web console for this queue manager.",
			},
			"rest_api_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url through which to access REST APIs for this queue manager.",
			},
			"administrator_api_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url through which to access the Admin REST APIs for this queue manager.",
			},
			"connection_info_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uri through which the CDDT for this queue manager can be obtained.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RFC3339 formatted UTC date for when the queue manager was created.",
			},
			"upgrade_available": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Describes whether an upgrade is available for this queue manager.",
			},
			"available_upgrade_versions_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uri through which the available versions to upgrade to can be found for this queue manager.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this queue manager.",
			},
			"queue_manager_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the queue manager which was allocated on creation, and can be used for delete calls.",
			},
		},
	}
}

func ResourceIbmMqcloudQueueManagerValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service_instance_guid",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9_.]*$`,
			MinValueLength:             1,
			MaxValueLength:             48,
		},
		validate.ValidateSchema{
			Identifier:                 "display_name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^.*$`,
			MinValueLength:             0,
			MaxValueLength:             150,
		},
		validate.ValidateSchema{
			Identifier:                 "location",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([^[:ascii:]]|[a-zA-Z0-9-._: ])+$`,
			MinValueLength:             2,
			MaxValueLength:             150,
		},
		validate.ValidateSchema{
			Identifier:                 "size",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "large, medium, small, xsmall",
		},
		validate.ValidateSchema{
			Identifier:                 "version",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[0-9]+.[0-9]+.[0-9]+_[0-9]+$`,
			MinValueLength:             7,
			MaxValueLength:             15,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_mqcloud_queue_manager", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmMqcloudQueueManagerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Create Queue Manager failed: %s", err.Error()), "ibm_mqcloud_queue_manager", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createQueueManagerOptions := &mqcloudv1.CreateQueueManagerOptions{}

	createQueueManagerOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	createQueueManagerOptions.SetName(d.Get("name").(string))
	createQueueManagerOptions.SetLocation(d.Get("location").(string))
	createQueueManagerOptions.SetSize(d.Get("size").(string))
	if _, ok := d.GetOk("display_name"); ok {
		createQueueManagerOptions.SetDisplayName(d.Get("display_name").(string))
	}
	if _, ok := d.GetOk("version"); ok {
		createQueueManagerOptions.SetVersion(d.Get("version").(string))
	}

	queueManagerTaskStatus, _, err := mqcloudClient.CreateQueueManagerWithContext(context, createQueueManagerOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateQueueManagerWithContext failed: %s", err.Error()), "ibm_mqcloud_queue_manager", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createQueueManagerOptions.ServiceInstanceGuid, *queueManagerTaskStatus.QueueManagerID))

	if waitForQmStatus {
		_, err = waitForQmStatusUpdate(context, d, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for Queue Manager (%s) to be in running state: %s", *queueManagerTaskStatus.QueueManagerID, err))
		}
	}

	return resourceIbmMqcloudQueueManagerRead(context, d, meta)
}

func resourceIbmMqcloudQueueManagerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getQueueManagerOptions := &mqcloudv1.GetQueueManagerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "sep-id-parts").GetDiag()
	}

	getQueueManagerOptions.SetServiceInstanceGuid(parts[0])
	getQueueManagerOptions.SetQueueManagerID(parts[1])

	var queueManagerDetails *mqcloudv1.QueueManagerDetails
	var response *core.DetailedResponse

	err = resource.RetryContext(context, 150*time.Second, func() *resource.RetryError {
		queueManagerDetails, response, err = mqcloudClient.GetQueueManagerWithContext(context, getQueueManagerOptions)
		if err != nil || queueManagerDetails == nil {
			if response != nil && response.StatusCode == 404 {
				return resource.RetryableError(fmt.Errorf("Queue Manager not found, retrying"))
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetQueueManagerWithContext failed: %s", err.Error()), "ibm_mqcloud_queue_manager", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("service_instance_guid", parts[0]); err != nil {
		err = fmt.Errorf("Error setting service_instance_guid: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-service_instance_guid").GetDiag()
	}
	if err = d.Set("name", queueManagerDetails.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-name").GetDiag()
	}
	if !core.IsNil(queueManagerDetails.DisplayName) {
		if err = d.Set("display_name", queueManagerDetails.DisplayName); err != nil {
			err = fmt.Errorf("Error setting display_name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-display_name").GetDiag()
		}
	}
	if err = d.Set("location", queueManagerDetails.Location); err != nil {
		err = fmt.Errorf("Error setting location: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-location").GetDiag()
	}
	if err = d.Set("size", queueManagerDetails.Size); err != nil {
		err = fmt.Errorf("Error setting size: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-size").GetDiag()
	}
	if !core.IsNil(queueManagerDetails.Version) {
		if err = d.Set("version", queueManagerDetails.Version); err != nil {
			err = fmt.Errorf("Error setting version: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-version").GetDiag()
		}
	}
	if err = d.Set("status_uri", queueManagerDetails.StatusURI); err != nil {
		err = fmt.Errorf("Error setting status_uri: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-status_uri").GetDiag()
	}
	if err = d.Set("web_console_url", queueManagerDetails.WebConsoleURL); err != nil {
		err = fmt.Errorf("Error setting web_console_url: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-web_console_url").GetDiag()
	}
	if err = d.Set("rest_api_endpoint_url", queueManagerDetails.RestApiEndpointURL); err != nil {
		err = fmt.Errorf("Error setting rest_api_endpoint_url: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-rest_api_endpoint_url").GetDiag()
	}
	if err = d.Set("administrator_api_endpoint_url", queueManagerDetails.AdministratorApiEndpointURL); err != nil {
		err = fmt.Errorf("Error setting administrator_api_endpoint_url: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-administrator_api_endpoint_url").GetDiag()
	}
	if err = d.Set("connection_info_uri", queueManagerDetails.ConnectionInfoURI); err != nil {
		err = fmt.Errorf("Error setting connection_info_uri: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-connection_info_uri").GetDiag()
	}
	if err = d.Set("date_created", flex.DateTimeToString(queueManagerDetails.DateCreated)); err != nil {
		err = fmt.Errorf("Error setting date_created: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-date_created").GetDiag()
	}
	if err = d.Set("upgrade_available", queueManagerDetails.UpgradeAvailable); err != nil {
		err = fmt.Errorf("Error setting upgrade_available: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-upgrade_available").GetDiag()
	}
	if err = d.Set("available_upgrade_versions_uri", queueManagerDetails.AvailableUpgradeVersionsURI); err != nil {
		err = fmt.Errorf("Error setting available_upgrade_versions_uri: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-available_upgrade_versions_uri").GetDiag()
	}
	if err = d.Set("href", queueManagerDetails.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-href").GetDiag()
	}
	if err = d.Set("queue_manager_id", queueManagerDetails.ID); err != nil {
		err = fmt.Errorf("Error setting queue_manager_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "read", "set-queue_manager_id").GetDiag()
	}

	return nil
}

func resourceIbmMqcloudQueueManagerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Update Queue Manager failed: %s", err.Error()), "ibm_mqcloud_queue_manager", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	setQueueManagerVersionOptions := &mqcloudv1.SetQueueManagerVersionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "update", "sep-id-parts").GetDiag()
	}

	setQueueManagerVersionOptions.SetServiceInstanceGuid(parts[0])
	setQueueManagerVersionOptions.SetQueueManagerID(parts[1])

	hasChange := false

	if d.HasChange("service_instance_guid") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "service_instance_guid")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_mqcloud_queue_manager", "update", "service_instance_guid-forces-new").GetDiag()
	}
	if d.HasChange("version") {
		oldVersion, newVersion := d.GetChange("version")
		if IsVersionDowngrade(oldVersion.(string), newVersion.(string)) {
			return diag.FromErr(fmt.Errorf("Version downgrade is not allowed"))
		}
		setQueueManagerVersionOptions.SetVersion(newVersion.(string))
		hasChange = true
	}

	if hasChange {
		_, _, err = mqcloudClient.SetQueueManagerVersionWithContext(context, setQueueManagerVersionOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("SetQueueManagerVersionWithContext failed: %s", err.Error()), "ibm_mqcloud_queue_manager", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if waitForQmStatus {
			_, err = waitForQmStatusUpdate(context, d, meta)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for Queue Manager (%s) to be in running state: %s", *setQueueManagerVersionOptions.QueueManagerID, err))
			}
		}
	}

	return resourceIbmMqcloudQueueManagerRead(context, d, meta)
}

func resourceIbmMqcloudQueueManagerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Delete Queue Manager failed: %s", err.Error()), "ibm_mqcloud_queue_manager", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteQueueManagerOptions := &mqcloudv1.DeleteQueueManagerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_queue_manager", "delete", "sep-id-parts").GetDiag()
	}

	deleteQueueManagerOptions.SetServiceInstanceGuid(parts[0])
	deleteQueueManagerOptions.SetQueueManagerID(parts[1])

	_, _, err = mqcloudClient.DeleteQueueManagerWithContext(context, deleteQueueManagerOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteQueueManagerWithContext failed: %s", err.Error()), "ibm_mqcloud_queue_manager", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if waitForQmStatus {
		_, err = waitForQueueManagerToDelete(context, d, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for Queue Manager (%s) to be in running state: %s", *deleteQueueManagerOptions.QueueManagerID, err))
		}
	}
	d.SetId("")

	return nil
}
