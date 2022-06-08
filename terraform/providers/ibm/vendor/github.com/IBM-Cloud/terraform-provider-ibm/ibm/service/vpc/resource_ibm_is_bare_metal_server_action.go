// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerActionAvailable = "available"
	isBareMetalServerActionPending   = "pending"
	isBareMetalServerActionFailed    = "failed"
	isBareMetalServerStopType        = "stop_type"
)

func ResourceIBMIsBareMetalServerAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISBareMetalServerActionCreate,
		ReadContext:   resourceIBMISBareMetalServerActionRead,
		UpdateContext: resourceIBMISBareMetalServerActionUpdate,
		DeleteContext: resourceIBMISBareMetalServerActionDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bare metal server identifier",
			},
			isBareMetalServerStopType: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Default:     "hard",
				Description: "The type of stop operation",
			},
			isBareMetalServerAction: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerAction),
				Description:  "This restart/start/stops a bare metal server.",
			},
			isBareMetalServerStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bare metal server status",
			},

			isBareMetalServerStatusReasons: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerStatusReasonsCode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason",
						},

						isBareMetalServerStatusReasonsMessage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISBareMetalServerActionValidator() *validate.ResourceValidator {
	bareMetalServerStopTypes := "soft, hard"
	bareMetalServerActions := "start, restart, stop"
	validateSchema := make([]validate.ValidateSchema, 1)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBareMetalServerAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              bareMetalServerActions})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBareMetalServerStopType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              bareMetalServerStopTypes})
	ibmISBareMetalServerActionResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_bare_metal_server_action", Schema: validateSchema}
	return &ibmISBareMetalServerActionResourceValidator
}

func resourceIBMISBareMetalServerActionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	bareMetalServerId := ""
	if bmsId, ok := d.GetOk(isBareMetalServerID); ok {
		bareMetalServerId = bmsId.(string)
	}

	bareMetalServerAction := ""
	if bmsAction, ok := d.GetOk(isBareMetalServerAction); ok {
		bareMetalServerAction = bmsAction.(string)
	}
	if bareMetalServerAction == "stop" {
		bareMetalServerStopType := "hard"
		if stopType, ok := d.GetOk(isBareMetalServerStopType); ok {
			bareMetalServerStopType = stopType.(string)
		}

		createBareMetalServerStopOptions := &vpcv1.StopBareMetalServerOptions{
			ID:   &bareMetalServerId,
			Type: &bareMetalServerStopType,
		}

		_, err = sess.StopBareMetalServerWithContext(context, createBareMetalServerStopOptions)
		if err != nil {
			return diag.FromErr(err)
		}
		_, waitErr := isWaitForBareMetalServerActionStop(sess, d.Timeout(schema.TimeoutDelete), bareMetalServerId, d)
		if waitErr != nil {
			return diag.FromErr(waitErr)
		}
	} else if bareMetalServerAction == "start" {

		createBareMetalServerStartOptions := &vpcv1.StartBareMetalServerOptions{
			ID: &bareMetalServerId,
		}

		_, err := sess.StartBareMetalServerWithContext(context, createBareMetalServerStartOptions)
		if err != nil {
			return diag.FromErr(err)
		}
		_, waitErr := isWaitForBareMetalServerActionAvailable(sess, bareMetalServerId, d.Timeout(schema.TimeoutDelete), d)
		if waitErr != nil {
			return diag.FromErr(waitErr)
		}
	} else if bareMetalServerAction == "restart" {
		createBareMetalServerRestartOptions := &vpcv1.RestartBareMetalServerOptions{
			ID: &bareMetalServerId,
		}

		_, err := sess.RestartBareMetalServerWithContext(context, createBareMetalServerRestartOptions)
		if err != nil {
			return diag.FromErr(err)
		}
		_, waitErr := isWaitForBareMetalServerActionAvailable(sess, bareMetalServerId, d.Timeout(schema.TimeoutDelete), d)
		if waitErr != nil {
			return diag.FromErr(waitErr)
		}
	}
	d.SetId(bareMetalServerId)
	err = bareMetalServerActionGet(context, sess, bareMetalServerId, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceIBMISBareMetalServerActionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()
	err = bareMetalServerActionGet(context, sess, id, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func bareMetalServerActionGet(context context.Context, sess *vpcv1.VpcV1, id string, d *schema.ResourceData) error {
	options := &vpcv1.GetBareMetalServerOptions{
		ID: &id,
	}
	bms, response, err := sess.GetBareMetalServerWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s): %s\n%s", id, err, response)
	}
	d.SetId(*bms.ID)
	d.Set(isBareMetalServerStatus, *bms.Status)
	statusReasonsList := make([]map[string]interface{}, 0)
	if bms.StatusReasons != nil {
		for _, sr := range bms.StatusReasons {
			currentSR := map[string]interface{}{}
			if sr.Code != nil && sr.Message != nil {
				currentSR[isBareMetalServerStatusReasonsCode] = *sr.Code
				currentSR[isBareMetalServerStatusReasonsMessage] = *sr.Message
				statusReasonsList = append(statusReasonsList, currentSR)
			}
		}
	}
	d.Set(isBareMetalServerStatusReasons, statusReasonsList)
	return nil
}

func resourceIBMISBareMetalServerActionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange(isBareMetalServerAction) {
		sess, err := vpcClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}
		bareMetalServerId := d.Id()

		bareMetalServerAction := ""
		if bmsAction, ok := d.GetOk(isBareMetalServerAction); ok {
			bareMetalServerAction = bmsAction.(string)
		}

		if bareMetalServerAction == "stop" {
			bareMetalServerStopType := "soft"
			if stopType, ok := d.GetOk(isBareMetalServerStopType); ok {
				bareMetalServerStopType = stopType.(string)
			}

			createBareMetalServerStopOptions := &vpcv1.StopBareMetalServerOptions{
				ID:   &bareMetalServerId,
				Type: &bareMetalServerStopType,
			}

			_, err := sess.StopBareMetalServerWithContext(context, createBareMetalServerStopOptions)
			if err != nil {
				return diag.FromErr(err)
			}
			_, waitErr := isWaitForBareMetalServerActionStop(sess, d.Timeout(schema.TimeoutDelete), bareMetalServerId, d)
			if waitErr != nil {
				return diag.FromErr(waitErr)
			}
		} else if bareMetalServerAction == "start" {
			createBareMetalServerStartOptions := &vpcv1.StartBareMetalServerOptions{
				ID: &bareMetalServerId,
			}

			_, err := sess.StartBareMetalServerWithContext(context, createBareMetalServerStartOptions)
			if err != nil {
				return diag.FromErr(err)
			}
			_, waitErr := isWaitForBareMetalServerActionAvailable(sess, bareMetalServerId, d.Timeout(schema.TimeoutDelete), d)
			if waitErr != nil {
				return diag.FromErr(waitErr)
			}
		} else if bareMetalServerAction == "restart" {
			createBareMetalServerRestartOptions := &vpcv1.RestartBareMetalServerOptions{
				ID: &bareMetalServerId,
			}

			_, err := sess.RestartBareMetalServerWithContext(context, createBareMetalServerRestartOptions)
			if err != nil {
				return diag.FromErr(err)
			}
			_, waitErr := isWaitForBareMetalServerActionAvailable(sess, bareMetalServerId, d.Timeout(schema.TimeoutDelete), d)
			if waitErr != nil {
				return diag.FromErr(waitErr)
			}
		}
		err = bareMetalServerActionGet(context, sess, bareMetalServerId, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceIBMISBareMetalServerActionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func isWaitForBareMetalServerActionAvailable(client *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for Bare Metal Server (%s) to be running.", id)
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isBareMetalServerStatusPending, isBareMetalServerActionStatusStarting},
		Target:     []string{isBareMetalServerStatusRunning, isBareMetalServerStatusFailed},
		Refresh:    isBareMetalServerActionRefreshFunc(client, id, d, communicator),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isBareMetalServerActionRefreshFunc(client *vpcv1.VpcV1, id string, d *schema.ResourceData, communicator chan interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		bmsgetoptions := &vpcv1.GetBareMetalServerOptions{
			ID: &id,
		}
		bms, response, err := client.GetBareMetalServer(bmsgetoptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting Bare Metal Server: %s\n%s", err, response)
		}
		d.Set(isBareMetalServerStatus, *bms.Status)

		select {
		case data := <-communicator:
			return nil, "", data.(error)
		default:
			fmt.Println("no message sent")
		}

		if *bms.Status == "running" {
			// let know the isRestartStartAction() to stop
			close(communicator)
			return bms, *bms.Status, nil

		}
		if *bms.Status == "failed" {
			// let know the isRestartStartAction() to stop
			close(communicator)
			return bms, *bms.Status, fmt.Errorf("[ERROR] Error Bare Metal Server is in failed state")

		}
		return bms, isBareMetalServerStatusPending, nil
	}
}
