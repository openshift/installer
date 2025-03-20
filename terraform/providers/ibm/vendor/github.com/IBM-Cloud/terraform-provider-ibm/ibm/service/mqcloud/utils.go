package mqcloud

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"strconv"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func loadWaitForQmStatusEnvVar() bool {
	envValue := os.Getenv("IBMCLOUD_MQCLOUD_WAIT_FOR_QM_STATUS")
	return strings.ToLower(envValue) != "false"
}

var waitForQmStatus = loadWaitForQmStatusEnvVar()
var planCache = make(map[string]string)

const qmStatus = "running"
const qmCreating = "initializing"
const isQueueManagerDeleting = "true"
const isQueueManagerDeleteDone = "true"
const reservedDeploymentPlan = "reserved-deployment"
const reservedCapacityPlan = "reserved-capacity"
const reservedCapacitySubscriptionPlan = "reserved-capacity-subscription"
const enforceReservedDeploymentPlan = true // change to false for testing in IKS clusters.

// waitForQmStatusUpdate waits for Queue Manager to be in running state
func waitForQmStatusUpdate(context context.Context, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return "", err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{qmCreating},
		Target:  []string{qmStatus},
		Refresh: func() (interface{}, string, error) {
			getQueueManagerStatusOptions := &mqcloudv1.GetQueueManagerStatusOptions{}
			parts, err := flex.SepIdParts(d.Id(), "/")
			if err != nil {
				return nil, "", err
			}
			getQueueManagerStatusOptions.SetServiceInstanceGuid(parts[0])
			getQueueManagerStatusOptions.SetQueueManagerID(parts[1])
			queueManagerStatus, response, err := mqcloudClient.GetQueueManagerStatusWithContext(context, getQueueManagerStatusOptions)
			if err != nil {
				return "", "", fmt.Errorf("GetQueueManagerWithContext ...... %s err: %s", response, err)
			}
			if queueManagerStatus == nil || queueManagerStatus.Status == nil {
				return nil, "", fmt.Errorf("queueManagerStatus or queueManagerStatus.Status is nil")
			}

			if *queueManagerStatus.Status == "running" {
				return queueManagerStatus, qmStatus, nil
			} else if *queueManagerStatus.Status == "initialization_failed" || *queueManagerStatus.Status == "restore_failed" || *queueManagerStatus.Status == "status_not_available" {
				return queueManagerStatus, qmStatus, fmt.Errorf("%s", err)
			}
			return queueManagerStatus, qmCreating, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}
	return stateConf.WaitForStateContext(context)
}

func waitForQueueManagerToDelete(context context.Context, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{isQueueManagerDeleting},
		Target:  []string{isQueueManagerDeleteDone},
		Refresh: func() (interface{}, string, error) {
			getQueueManagerOptions := &mqcloudv1.GetQueueManagerOptions{}

			parts, err := flex.SepIdParts(d.Id(), "/")
			if err != nil {
				return "", "", err
			}

			getQueueManagerOptions.SetServiceInstanceGuid(parts[0])
			getQueueManagerOptions.SetQueueManagerID(parts[1])

			queueManagerDetails, response, err := mqcloudClient.GetQueueManagerWithContext(context, getQueueManagerOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return queueManagerDetails, isQueueManagerDeleteDone, nil
				}
				return nil, "", err
			}
			return queueManagerDetails, isQueueManagerDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForStateContext(context)
}

func IsVersionDowngrade(oldVersion, newVersion string) bool {
	oldParts := strings.Split(oldVersion, ".")
	newParts := strings.Split(newVersion, ".")

	for i := 0; i < len(oldParts); i++ {
		oldPartNum := strings.Split(oldParts[i], "_")[0]
		newPartNum := strings.Split(newParts[i], "_")[0]

		oldPart, _ := strconv.Atoi(oldPartNum)
		newPart, _ := strconv.Atoi(newPartNum)

		if newPart < oldPart {
			return true
		} else if newPart > oldPart {
			return false
		}
	}

	oldPatchNum := strings.Split(oldParts[len(oldParts)-1], "_")
	newPatchNum := strings.Split(newParts[len(newParts)-1], "_")

	if len(oldPatchNum) > 1 && len(newPatchNum) > 1 {
		oldPatch, _ := strconv.Atoi(oldPatchNum[1])
		newPatch, _ := strconv.Atoi(newPatchNum[1])

		return newPatch < oldPatch
	}

	return false
}

// Function variable to allow mocking in tests
var fetchServicePlanFunc = fetchServicePlan

// checkSIPlan checks whether a Service Instance (SI) is in a Reserved Deployment.
func checkSIPlan(d *schema.ResourceData, meta interface{}) error {
	if !enforceReservedDeploymentPlan {
		return nil
	}

	instanceID := d.Get("service_instance_guid").(string)
	// Check cache first
	if plan, found := planCache[instanceID]; found {
		return handlePlanCheck(plan, instanceID)
	}

	plan, err := fetchServicePlanFunc(meta, instanceID)
	if err != nil {
		return err
	}

	// Update cache
	planCache[instanceID] = plan

	return handlePlanCheck(plan, instanceID)
}

// fetchServicePlan handles all external function calls to fetch the service plan name.
func fetchServicePlan(meta interface{}, instanceID string) (string, error) {
	clientAPIProvider, ok := meta.(conns.ClientSession)
	if !ok {
		return "", fmt.Errorf("meta does not implement ClientAPIProvider")
	}

	// Creating a Resource Controller Client
	rsConClient, err := clientAPIProvider.ResourceControllerV2API()
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to create Resource Controller Client: %s", err)
	}

	// Getting the resource instance
	rsInst := resourcecontrollerv2.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	instance, response, err := rsConClient.GetResourceInstance(&rsInst)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to retrieve Resource Instance: %s, Response: %s", err, response)
	}

	// Creating a Resource Catalog Client
	rsCatClient, err := clientAPIProvider.ResourceCatalogAPI()
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to create Resource Catalog Client: %s", err)
	}
	rsCatRepo := rsCatClient.ResourceCatalog()
	if rsCatRepo == nil {
		return "", fmt.Errorf("[ERROR] Catalog client is nil")
	}
	if instance.ResourcePlanID == nil {
		return "", fmt.Errorf("[ERROR] Instance ResourcePlanID is nil")
	}

	// Checking the service plan
	plan, err := rsCatRepo.GetServicePlanName(*instance.ResourcePlanID)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to retrieve Service Plan: %s", err)
	}

	return plan, nil
}

func handlePlanCheck(plan string, instanceID string) error {
	if !strings.Contains(plan, reservedDeploymentPlan) && !strings.Contains(plan, reservedCapacityPlan) && !strings.Contains(plan, reservedCapacitySubscriptionPlan) {
		return fmt.Errorf("[ERROR] Terraform is only supported for Reserved Deployment, Reserved Capacity, and Reserved Capacity Subscription Plans. Your Service Plan is: %s for the instance %s", plan, instanceID)
	}
	return nil
}
