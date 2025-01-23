// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vmware

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/IBM/vmware-go-sdk/vmwarev1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func loadWaitForVdcStatusEnvVar() bool {
	envValue := os.Getenv("IBM_VMAAS_WAIT_FOR_VDC_STATUS")
	return strings.ToLower(envValue) != "false"
}

var waitForVdcStatus = loadWaitForVdcStatusEnvVar()

const VdcFinalState = "ready_to_use"
const VdcCreatingState = "creating"
const isVdcDeleting = "false"
const isVdcDeleteDone = "true"

// waits for Vdc instance to be in ready state
func waitForVdcStatusUpdate(context context.Context, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	vmwareClient, err := meta.(conns.ClientSession).VmwareV1()
	if err != nil {
		return "", err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{VdcCreatingState},
		Target:  []string{VdcFinalState},
		Refresh: func() (interface{}, string, error) {
			getVdcOptions := &vmwarev1.GetVdcOptions{}

			getVdcOptions.SetID(d.Id())

			vdc, response, err := vmwareClient.GetVdcWithContext(context, getVdcOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil, "", err
				}
				return nil, "", err
			}
			if err = d.Set("status", vdc.Status); err != nil {
				return "", "", fmt.Errorf("Error setting status: %s", err)
			}

			fmt.Println("The vdc is currently in the " + *vdc.Status + " state ....")

			if *vdc.Status == "ready_to_use" {
				return vdc, VdcFinalState, nil
			} else if *vdc.Status == "failed" {
				return vdc, VdcFinalState, fmt.Errorf("%s", err)
			}
			return vdc, VdcCreatingState, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}
	return stateConf.WaitForStateContext(context)
}

func waitForVdcToDelete(context context.Context, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	vmwareClient, err := meta.(conns.ClientSession).VmwareV1()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{isVdcDeleting},
		Target:  []string{isVdcDeleteDone},
		Refresh: func() (interface{}, string, error) {
			getVdcOptions := &vmwarev1.GetVdcOptions{}

			getVdcOptions.SetID(d.Id())

			vdc, response, err := vmwareClient.GetVdcWithContext(context, getVdcOptions)

			if err != nil {
				if response != nil && response.StatusCode == 404 {
					fmt.Println("The vdc is deleted.")
					return vdc, isVdcDeleteDone, nil
				}
				return nil, "", err
			}
			fmt.Println("The vdc is currently in the " + *vdc.Status + " state ....")
			return vdc, isVdcDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForStateContext(context)
}
