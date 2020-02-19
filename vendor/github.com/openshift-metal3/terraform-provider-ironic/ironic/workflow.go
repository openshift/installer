package ironic

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"log"
	"time"
)

// provisionStateWorkflow is used to track state through the process of updating's it's provision state
type provisionStateWorkflow struct {
	client *gophercloud.ServiceClient
	node   nodes.Node
	uuid   string
	target nodes.TargetProvisionState
	wait   time.Duration

	configDrive interface{}
}

// ChangeProvisionStateToTarget drives Ironic's state machine through the process to reach our desired end state. This requires multiple
// possibly long-running steps.  If required, we'll build a config drive ISO for deployment.
func ChangeProvisionStateToTarget(client *gophercloud.ServiceClient, uuid string, target nodes.TargetProvisionState, configDrive interface{}) error {

	// Run the provisionStateWorkflow - this could take a while
	wf := provisionStateWorkflow{
		target:      target,
		client:      client,
		wait:        5 * time.Second,
		uuid:        uuid,
		configDrive: configDrive,
	}

	err := wf.run()
	return err
}

// Keep driving the state machine forward
func (workflow *provisionStateWorkflow) run() error {
	log.Printf("[INFO] Beginning provisioning workflow, will try to change node to state '%s'", workflow.target)

	for {
		log.Printf("[DEBUG] Node is in state '%s'", workflow.node.ProvisionState)

		done, err := workflow.next()
		if done || err != nil {
			return err
		}

		time.Sleep(workflow.wait)
	}

	return nil
}

// Do the next thing to get us to our target state
func (workflow *provisionStateWorkflow) next() (done bool, err error) {
	// Refresh the node on each run
	if err := workflow.reloadNode(); err != nil {
		return true, err
	}

	log.Printf("[DEBUG] Node current state is '%s', target is %s", workflow.node.ProvisionState, workflow.target)

	switch target := nodes.TargetProvisionState(workflow.target); target {
	case nodes.TargetManage:
		return workflow.toManageable()
	case nodes.TargetProvide:
		return workflow.toAvailable()
	case nodes.TargetActive:
		return workflow.toActive()
	case nodes.TargetDeleted:
		return workflow.toDeleted()
	case nodes.TargetClean:
		return workflow.toClean()
	case nodes.TargetInspect:
		return workflow.toInspect()
	default:
		return true, fmt.Errorf("unknown target state '%s'", target)
	}
}

// Change a node to "manageable" stable
func (workflow *provisionStateWorkflow) toManageable() (done bool, err error) {
	switch state := workflow.node.ProvisionState; state {
	case "manageable":
		// We're done!
		return true, err
	case "enroll",
		"adopt failed",
		"clean failed",
		"inspect failed",
		"available":
		return workflow.changeProvisionState(nodes.TargetManage)
	case "verifying":
		// Not done, no error - Ironic is working
		return false, nil

	default:
		return true, fmt.Errorf("cannot go from state '%s' to state 'manageable'", state)
	}

	return false, nil
}

// Clean a node
func (workflow *provisionStateWorkflow) toClean() (done bool, err error) {
	// Node must be manageable first
	workflow.reloadNode()
	if workflow.node.ProvisionState != string(nodes.Manageable) {
		if err := ChangeProvisionStateToTarget(workflow.client, workflow.uuid, nodes.TargetManage, nil); err != nil {
			return true, err
		}
	}

	// Set target to clean
	workflow.changeProvisionState(nodes.TargetClean)

	for {
		workflow.reloadNode()
		state := workflow.node.ProvisionState

		switch state {
		case "manageable":
			return true, nil
		case "cleaning",
			"clean wait":
			// Not done, no error - Ironic is working
			continue
		default:
			return true, fmt.Errorf("could not clean node, node is currently '%s', last error was '%s'", state, workflow.node.LastError)
		}
	}

	return true, nil
}

// Inspect a node
func (workflow *provisionStateWorkflow) toInspect() (done bool, err error) {
	// Node must be manageable first
	workflow.reloadNode()
	if workflow.node.ProvisionState != string(nodes.Manageable) {
		if err := ChangeProvisionStateToTarget(workflow.client, workflow.uuid, nodes.TargetManage, nil); err != nil {
			return true, err
		}
	}

	// Set target to inspect
	workflow.changeProvisionState(nodes.TargetInspect)

	for {
		workflow.reloadNode()
		state := workflow.node.ProvisionState

		switch state {
		case "manageable":
			return true, nil
		case "inspecting",
			"inspect wait":
			// Not done, no error - Ironic is working
			continue
		default:
			return true, fmt.Errorf("could not inspect node, node is currently '%s', last error was '%s'", state, workflow.node.LastError)
		}
	}

	return true, nil
}

// Change a node to "available" state
func (workflow *provisionStateWorkflow) toAvailable() (done bool, err error) {
	switch state := workflow.node.ProvisionState; state {
	case "available":
		// We're done!
		return true, nil
	case "cleaning",
		"clean wait":
		// Not done, no error - Ironic is working
		log.Printf("[DEBUG] Node %s is '%s', waiting for Ironic to finish.", workflow.uuid, state)
		return false, nil
	case "manageable":
		// From manageable, we can go to provide
		log.Printf("[DEBUG] Node %s is '%s', going to change to 'available'", workflow.uuid, state)
		return workflow.changeProvisionState(nodes.TargetProvide)
	default:
		// Otherwise we have to get into manageable state first
		log.Printf("[DEBUG] Node %s is '%s', going to change to 'manageable'.", workflow.uuid, state)
		_, err := workflow.toManageable()
		if err != nil {
			return true, err
		}
		return false, nil
	}

	return false, nil
}

// Change a node to "active" state
func (workflow *provisionStateWorkflow) toActive() (bool, error) {

	switch state := workflow.node.ProvisionState; state {
	case "active":
		// We're done!
		log.Printf("[DEBUG] Node %s is 'active', we are done.", workflow.uuid)
		return true, nil
	case "deploying",
		"wait call-back":
		// Not done, no error - Ironic is working
		log.Printf("[DEBUG] Node %s is '%s', waiting for Ironic to finish.", workflow.uuid, state)
		return false, nil
	case "available":
		// From available, we can go to active
		log.Printf("[DEBUG] Node %s is 'available', going to change to 'active'.", workflow.uuid)
		workflow.wait = 30 * time.Second // Deployment takes a while
		return workflow.changeProvisionState(nodes.TargetActive)
	default:
		// Otherwise we have to get into available state first
		log.Printf("[DEBUG] Node %s is '%s', going to change to 'available'.", workflow.uuid, state)
		_, err := workflow.toAvailable()
		if err != nil {
			return true, err
		}
		return false, nil
	}
}

// Change a node to be "deleted," and remove the object from Ironic
func (workflow *provisionStateWorkflow) toDeleted() (bool, error) {
	switch state := workflow.node.ProvisionState; state {
	case "manageable",
		"available",
		"enroll":
		// We're done deleting the node
		return true, nil
	case "cleaning",
		"deleting":
		// Not done, no error - Ironic is working
		log.Printf("[DEBUG] Node %s is '%s', waiting for Ironic to finish.", workflow.uuid, state)
		return false, nil
	case "active",
		"wait call-back",
		"deploy failed",
		"error":
		log.Printf("[DEBUG] Node %s is '%s', going to change to 'deleted'.", workflow.uuid, state)
		return workflow.changeProvisionState(nodes.TargetDeleted)
	case "inspect failed",
		"clean failed":
		// We have to get into manageable state first
		log.Printf("[DEBUG] Node %s is '%s', going to change to 'manageable'.", workflow.uuid, state)
		_, err := workflow.toManageable()
		if err != nil {
			return true, err
		}
		return false, nil
	default:
		return true, fmt.Errorf("cannot delete node in state '%s'", state)
	}

	return false, nil
}

// Builds the ProvisionStateOpts to send to Ironic -- including config drive.
func (workflow *provisionStateWorkflow) buildProvisionStateOpts(target nodes.TargetProvisionState) (*nodes.ProvisionStateOpts, error) {
	opts := nodes.ProvisionStateOpts{
		Target: target,
	}

	// If we're deploying, then build a config drive to send to Ironic
	if target == "active" {
		opts.ConfigDrive = workflow.configDrive
	}

	return &opts, nil
}

// Call Ironic's API and issue the change provision state request.
func (workflow *provisionStateWorkflow) changeProvisionState(target nodes.TargetProvisionState) (done bool, err error) {
	opts, err := workflow.buildProvisionStateOpts(target)
	if err != nil {
		log.Printf("[ERROR] Unable to construct provisioning state options: %s", err.Error())
		return true, err
	}

	interval := 5 * time.Second
	for retries := 0; retries < 5; retries++ {
		err = nodes.ChangeProvisionState(workflow.client, workflow.uuid, *opts).ExtractErr()
		if _, ok := err.(gophercloud.ErrDefault409); ok {
			log.Printf("[DEBUG] Failed to change provision state: ironic is busy, will retry in %s.", interval.String())
			time.Sleep(interval)
			interval *= 2
		} else {
			break
		}
	}

	return false, err
}

// Call Ironic's API and reload the node's current state
func (workflow *provisionStateWorkflow) reloadNode() error {
	return nodes.Get(workflow.client, workflow.uuid).ExtractInto(&workflow.node)
}
