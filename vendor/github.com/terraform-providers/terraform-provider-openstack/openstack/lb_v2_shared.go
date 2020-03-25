package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gophercloud/gophercloud"
	octavialisteners "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	octaviamonitors "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
	neutronlisteners "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	neutronmonitors "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/monitors"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
)

// lbPendingStatuses are the valid statuses a LoadBalancer will be in while
// it's updating.
var lbPendingStatuses = []string{"PENDING_CREATE", "PENDING_UPDATE"}

// lbPendingDeleteStatuses are the valid statuses a LoadBalancer will be before delete
var lbPendingDeleteStatuses = []string{"ERROR", "PENDING_UPDATE", "PENDING_DELETE", "ACTIVE"}

var lbSkipLBStatuses = []string{"ERROR", "ACTIVE"}

// chooseLBV2Client will determine which load balacing client to use:
// Either the Octavia/LBaaS client or the Neutron/Networking v2 client.
func chooseLBV2Client(d *schema.ResourceData, config *Config) (*gophercloud.ServiceClient, error) {
	if config.UseOctavia {
		return config.LoadBalancerV2Client(GetRegion(d, config))
	}
	return config.NetworkingV2Client(GetRegion(d, config))
}

// chooseLBV2AccTestClient will determine which load balacing client to use:
// Either the Octavia/LBaaS client or the Neutron/Networking v2 client.
// This is similar to the chooseLBV2Client function but specific for acceptance
// tests.
func chooseLBV2AccTestClient(config *Config, region string) (*gophercloud.ServiceClient, error) {
	if config.UseOctavia {
		return config.LoadBalancerV2Client(region)
	}
	return config.NetworkingV2Client(region)
}

// chooseLBV2ListenerCreateOpts will determine which load balancer listener Create options to use:
// Either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBV2ListenerCreateOpts(d *schema.ResourceData, config *Config) neutronlisteners.CreateOptsBuilder {
	adminStateUp := d.Get("admin_state_up").(bool)

	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_container_refs"); ok {
		for _, v := range raw.([]interface{}) {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}

	var createOpts neutronlisteners.CreateOptsBuilder

	if config.UseOctavia {
		// Use Octavia.
		opts := octavialisteners.CreateOpts{
			Protocol:               octavialisteners.Protocol(d.Get("protocol").(string)),
			ProtocolPort:           d.Get("protocol_port").(int),
			ProjectID:              d.Get("tenant_id").(string),
			LoadbalancerID:         d.Get("loadbalancer_id").(string),
			Name:                   d.Get("name").(string),
			DefaultPoolID:          d.Get("default_pool_id").(string),
			Description:            d.Get("description").(string),
			DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
			SniContainerRefs:       sniContainerRefs,
			AdminStateUp:           &adminStateUp,
		}

		if v, ok := d.GetOk("connection_limit"); ok {
			connectionLimit := v.(int)
			opts.ConnLimit = &connectionLimit
		}

		if v, ok := d.GetOk("timeout_client_data"); ok {
			timeoutClientData := v.(int)
			opts.TimeoutClientData = &timeoutClientData
		}

		if v, ok := d.GetOk("timeout_member_connect"); ok {
			timeoutMemberConnect := v.(int)
			opts.TimeoutMemberConnect = &timeoutMemberConnect
		}

		if v, ok := d.GetOk("timeout_member_data"); ok {
			timeoutMemberData := v.(int)
			opts.TimeoutMemberData = &timeoutMemberData
		}

		if v, ok := d.GetOk("timeout_tcp_inspect"); ok {
			timeoutTCPInspect := v.(int)
			opts.TimeoutTCPInspect = &timeoutTCPInspect
		}

		createOpts = opts
	} else {
		// Use Neutron.
		opts := neutronlisteners.CreateOpts{
			Protocol:               neutronlisteners.Protocol(d.Get("protocol").(string)),
			ProtocolPort:           d.Get("protocol_port").(int),
			TenantID:               d.Get("tenant_id").(string),
			LoadbalancerID:         d.Get("loadbalancer_id").(string),
			Name:                   d.Get("name").(string),
			DefaultPoolID:          d.Get("default_pool_id").(string),
			Description:            d.Get("description").(string),
			DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
			SniContainerRefs:       sniContainerRefs,
			AdminStateUp:           &adminStateUp,
		}

		if v, ok := d.GetOk("connection_limit"); ok {
			connectionLimit := v.(int)
			opts.ConnLimit = &connectionLimit
		}

		createOpts = opts
	}

	return createOpts
}

// chooseLBV2ListenerUpdateOpts will determine which load balancer listener Update options to use:
// Either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBV2ListenerUpdateOpts(d *schema.ResourceData, config *Config) neutronlisteners.UpdateOptsBuilder {
	var hasChange bool

	if config.UseOctavia {
		// Use Octavia.
		var opts octavialisteners.UpdateOpts
		if d.HasChange("name") {
			hasChange = true
			name := d.Get("name").(string)
			opts.Name = &name
		}
		if d.HasChange("description") {
			hasChange = true
			description := d.Get("description").(string)
			opts.Description = &description
		}
		if d.HasChange("connection_limit") {
			hasChange = true
			connLimit := d.Get("connection_limit").(int)
			opts.ConnLimit = &connLimit
		}
		if d.HasChange("timeout_client_data") {
			hasChange = true
			timeoutClientData := d.Get("timeout_client_data").(int)
			opts.TimeoutClientData = &timeoutClientData
		}
		if d.HasChange("timeout_member_connect") {
			hasChange = true
			timeoutMemberConnect := d.Get("timeout_member_connect").(int)
			opts.TimeoutMemberConnect = &timeoutMemberConnect
		}
		if d.HasChange("timeout_member_data") {
			hasChange = true
			timeoutMemberData := d.Get("timeout_member_data").(int)
			opts.TimeoutMemberData = &timeoutMemberData
		}
		if d.HasChange("timeout_tcp_inspect") {
			hasChange = true
			timeoutTCPInspect := d.Get("timeout_tcp_inspect").(int)
			opts.TimeoutTCPInspect = &timeoutTCPInspect
		}
		if d.HasChange("default_pool_id") {
			hasChange = true
			defaultPoolID := d.Get("default_pool_id").(string)
			opts.DefaultPoolID = &defaultPoolID
		}
		if d.HasChange("default_tls_container_ref") {
			hasChange = true
			defaultTlsContainerRef := d.Get("default_tls_container_ref").(string)
			opts.DefaultTlsContainerRef = &defaultTlsContainerRef
		}
		if d.HasChange("sni_container_refs") {
			hasChange = true
			var sniContainerRefs []string
			if raw, ok := d.GetOk("sni_container_refs"); ok {
				for _, v := range raw.([]interface{}) {
					sniContainerRefs = append(sniContainerRefs, v.(string))
				}
			}
			opts.SniContainerRefs = &sniContainerRefs
		}
		if d.HasChange("admin_state_up") {
			hasChange = true
			asu := d.Get("admin_state_up").(bool)
			opts.AdminStateUp = &asu
		}

		if hasChange {
			return opts
		}
	} else {
		// Use Neutron.
		var opts neutronlisteners.UpdateOpts
		if d.HasChange("name") {
			hasChange = true
			name := d.Get("name").(string)
			opts.Name = &name
		}
		if d.HasChange("description") {
			hasChange = true
			description := d.Get("description").(string)
			opts.Description = &description
		}
		if d.HasChange("connection_limit") {
			hasChange = true
			connLimit := d.Get("connection_limit").(int)
			opts.ConnLimit = &connLimit
		}
		if d.HasChange("default_pool_id") {
			hasChange = true
			defaultPoolID := d.Get("default_pool_id").(string)
			opts.DefaultPoolID = &defaultPoolID
		}
		if d.HasChange("default_tls_container_ref") {
			hasChange = true
			defaultTlsContainerRef := d.Get("default_tls_container_ref").(string)
			opts.DefaultTlsContainerRef = &defaultTlsContainerRef
		}
		if d.HasChange("sni_container_refs") {
			hasChange = true
			var sniContainerRefs []string
			if raw, ok := d.GetOk("sni_container_refs"); ok {
				for _, v := range raw.([]interface{}) {
					sniContainerRefs = append(sniContainerRefs, v.(string))
				}
			}
			opts.SniContainerRefs = &sniContainerRefs
		}
		if d.HasChange("admin_state_up") {
			hasChange = true
			asu := d.Get("admin_state_up").(bool)
			opts.AdminStateUp = &asu
		}

		if hasChange {
			return opts
		}
	}

	return nil
}

func waitForLBV2Listener(lbClient *gophercloud.ServiceClient, listener *listeners.Listener, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for openstack_lb_listener_v2 %s to become %s.", listener.ID, target)

	if len(listener.Loadbalancers) == 0 {
		return fmt.Errorf("Failed to detect a openstack_lb_loadbalancer_v2 for the %s openstack_lb_listener_v2", listener.ID)
	}

	lbID := listener.Loadbalancers[0].ID

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2ListenerRefreshFunc(lbClient, lbID, listener),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for openstack_lb_listener_v2 %s to become %s: %s", listener.ID, target, err)
	}

	return nil
}

func resourceLBV2ListenerRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, listener *listeners.Listener) resource.StateRefreshFunc {
	if listener.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(lbSkipLBStatuses, status) {
				return lb, status, nil
			}

			listener, err := listeners.Get(lbClient, listener.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return listener, listener.ProvisioningStatus, nil
		}
	}

	return resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "listener", listener.ID, "")
}

// chooseLBV2MonitorCreateOpts will determine which load balancer monitor Create options to use:
// Either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBV2MonitorCreateOpts(d *schema.ResourceData, config *Config) neutronmonitors.CreateOptsBuilder {
	adminStateUp := d.Get("admin_state_up").(bool)

	var createOpts neutronmonitors.CreateOptsBuilder

	if config.UseOctavia {
		// Use Octavia.
		opts := octaviamonitors.CreateOpts{
			PoolID:         d.Get("pool_id").(string),
			TenantID:       d.Get("tenant_id").(string),
			Type:           d.Get("type").(string),
			Delay:          d.Get("delay").(int),
			Timeout:        d.Get("timeout").(int),
			MaxRetries:     d.Get("max_retries").(int),
			MaxRetriesDown: d.Get("max_retries_down").(int),
			URLPath:        d.Get("url_path").(string),
			HTTPMethod:     d.Get("http_method").(string),
			ExpectedCodes:  d.Get("expected_codes").(string),
			Name:           d.Get("name").(string),
			AdminStateUp:   &adminStateUp,
		}

		createOpts = opts
	} else {
		// Use Neutron.
		opts := neutronmonitors.CreateOpts{
			PoolID:        d.Get("pool_id").(string),
			TenantID:      d.Get("tenant_id").(string),
			Type:          d.Get("type").(string),
			Delay:         d.Get("delay").(int),
			Timeout:       d.Get("timeout").(int),
			MaxRetries:    d.Get("max_retries").(int),
			URLPath:       d.Get("url_path").(string),
			HTTPMethod:    d.Get("http_method").(string),
			ExpectedCodes: d.Get("expected_codes").(string),
			Name:          d.Get("name").(string),
			AdminStateUp:  &adminStateUp,
		}

		createOpts = opts
	}

	return createOpts
}

// chooseLBV2MonitorUpdateOpts will determine which load balancer monitor Update options to use:
// Either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBV2MonitorUpdateOpts(d *schema.ResourceData, config *Config) neutronmonitors.UpdateOptsBuilder {
	var hasChange bool

	if config.UseOctavia {
		// Use Octavia.
		var opts octaviamonitors.UpdateOpts

		if d.HasChange("url_path") {
			hasChange = true
			opts.URLPath = d.Get("url_path").(string)
		}
		if d.HasChange("expected_codes") {
			hasChange = true
			opts.ExpectedCodes = d.Get("expected_codes").(string)
		}
		if d.HasChange("delay") {
			hasChange = true
			opts.Delay = d.Get("delay").(int)
		}
		if d.HasChange("timeout") {
			hasChange = true
			opts.Timeout = d.Get("timeout").(int)
		}
		if d.HasChange("max_retries") {
			hasChange = true
			opts.MaxRetries = d.Get("max_retries").(int)
		}
		if d.HasChange("max_retries_down") {
			hasChange = true
			opts.MaxRetriesDown = d.Get("max_retries_down").(int)
		}
		if d.HasChange("admin_state_up") {
			hasChange = true
			asu := d.Get("admin_state_up").(bool)
			opts.AdminStateUp = &asu
		}
		if d.HasChange("name") {
			hasChange = true
			name := d.Get("name").(string)
			opts.Name = &name
		}
		if d.HasChange("http_method") {
			hasChange = true
			opts.HTTPMethod = d.Get("http_method").(string)
		}

		if hasChange {
			return opts
		}
	} else {
		// Use Neutron.
		var opts neutronmonitors.UpdateOpts

		if d.HasChange("url_path") {
			hasChange = true
			opts.URLPath = d.Get("url_path").(string)
		}
		if d.HasChange("expected_codes") {
			hasChange = true
			opts.ExpectedCodes = d.Get("expected_codes").(string)
		}
		if d.HasChange("delay") {
			hasChange = true
			opts.Delay = d.Get("delay").(int)
		}
		if d.HasChange("timeout") {
			hasChange = true
			opts.Timeout = d.Get("timeout").(int)
		}
		if d.HasChange("max_retries") {
			hasChange = true
			opts.MaxRetries = d.Get("max_retries").(int)
		}
		if d.HasChange("admin_state_up") {
			hasChange = true
			asu := d.Get("admin_state_up").(bool)
			opts.AdminStateUp = &asu
		}
		if d.HasChange("name") {
			hasChange = true
			name := d.Get("name").(string)
			opts.Name = &name
		}
		if d.HasChange("http_method") {
			hasChange = true
			opts.HTTPMethod = d.Get("http_method").(string)
		}

		if hasChange {
			return opts
		}
	}

	return nil
}

func waitForLBV2LoadBalancer(lbClient *gophercloud.ServiceClient, lbID string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for loadbalancer %s to become %s.", lbID, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID),
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: loadbalancer %s not found: %s", lbID, err)
			}
		}
		return fmt.Errorf("Error waiting for loadbalancer %s to become %s: %s", lbID, target, err)
	}

	return nil
}

func resourceLBV2LoadBalancerRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := loadbalancers.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return lb, lb.ProvisioningStatus, nil
	}
}

func waitForLBV2Member(lbClient *gophercloud.ServiceClient, parentPool *pools.Pool, member *pools.Member, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for member %s to become %s.", member.ID, target)

	lbID, err := lbV2FindLBIDviaPool(lbClient, parentPool)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2MemberRefreshFunc(lbClient, lbID, parentPool.ID, member),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for member %s to become %s: %s", member.ID, target, err)
	}

	return nil
}

func resourceLBV2MemberRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, poolID string, member *pools.Member) resource.StateRefreshFunc {
	if member.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(lbSkipLBStatuses, status) {
				return lb, status, nil
			}

			member, err := pools.GetMember(lbClient, poolID, member.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return member, member.ProvisioningStatus, nil
		}
	}

	return resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "member", member.ID, poolID)
}

func waitForLBV2Monitor(lbClient *gophercloud.ServiceClient, parentPool *pools.Pool, monitor *neutronmonitors.Monitor, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for openstack_lb_monitor_v2 %s to become %s.", monitor.ID, target)

	lbID, err := lbV2FindLBIDviaPool(lbClient, parentPool)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2MonitorRefreshFunc(lbClient, lbID, monitor),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}
		return fmt.Errorf("Error waiting for openstack_lb_monitor_v2 %s to become %s: %s", monitor.ID, target, err)
	}

	return nil
}

func resourceLBV2MonitorRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, monitor *neutronmonitors.Monitor) resource.StateRefreshFunc {
	if monitor.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(lbSkipLBStatuses, status) {
				return lb, status, nil
			}

			monitor, err := neutronmonitors.Get(lbClient, monitor.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return monitor, monitor.ProvisioningStatus, nil
		}
	}

	return resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "monitor", monitor.ID, "")
}

func waitForLBV2Pool(lbClient *gophercloud.ServiceClient, pool *pools.Pool, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for pool %s to become %s.", pool.ID, target)

	lbID, err := lbV2FindLBIDviaPool(lbClient, pool)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2PoolRefreshFunc(lbClient, lbID, pool),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for pool %s to become %s: %s", pool.ID, target, err)
	}

	return nil
}

func resourceLBV2PoolRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, pool *pools.Pool) resource.StateRefreshFunc {
	if pool.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(lbSkipLBStatuses, status) {
				return lb, status, nil
			}

			pool, err := pools.Get(lbClient, pool.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return pool, pool.ProvisioningStatus, nil
		}
	}

	return resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "pool", pool.ID, "")
}

func lbV2FindLBIDviaPool(lbClient *gophercloud.ServiceClient, pool *pools.Pool) (string, error) {
	if len(pool.Loadbalancers) > 0 {
		return pool.Loadbalancers[0].ID, nil
	}

	if len(pool.Listeners) > 0 {
		listenerID := pool.Listeners[0].ID
		listener, err := listeners.Get(lbClient, listenerID).Extract()
		if err != nil {
			return "", err
		}

		if len(listener.Loadbalancers) > 0 {
			return listener.Loadbalancers[0].ID, nil
		}
	}

	return "", fmt.Errorf("Unable to determine loadbalancer ID from pool %s", pool.ID)
}

func resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient *gophercloud.ServiceClient, lbID, resourceType, resourceID string, parentID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		statuses, err := loadbalancers.GetStatuses(lbClient, lbID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil, "", gophercloud.ErrDefault404{
					gophercloud.ErrUnexpectedResponseCode{BaseError: gophercloud.BaseError{
						DefaultErrString: fmt.Sprintf("Unable to get statuses from the Load Balancer %s statuses tree: %s", lbID, err),
					}},
				}
			}
			return nil, "", fmt.Errorf("Unable to get statuses from the Load Balancer %s statuses tree: %s", lbID, err)
		}

		// Don't fail, when statuses returns "null"
		if statuses == nil || statuses.Loadbalancer == nil {
			statuses = new(loadbalancers.StatusTree)
			statuses.Loadbalancer = new(loadbalancers.LoadBalancer)
		} else if !strSliceContains(lbSkipLBStatuses, statuses.Loadbalancer.ProvisioningStatus) {
			return statuses.Loadbalancer, statuses.Loadbalancer.ProvisioningStatus, nil
		}

		switch resourceType {
		case "listener":
			for _, listener := range statuses.Loadbalancer.Listeners {
				if listener.ID == resourceID {
					if listener.ProvisioningStatus != "" {
						return listener, listener.ProvisioningStatus, nil
					}
				}
			}
			listener, err := listeners.Get(lbClient, resourceID).Extract()
			return listener, "ACTIVE", err

		case "pool":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.ID == resourceID {
					if pool.ProvisioningStatus != "" {
						return pool, pool.ProvisioningStatus, nil
					}
				}
			}
			pool, err := pools.Get(lbClient, resourceID).Extract()
			return pool, "ACTIVE", err

		case "monitor":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.Monitor.ID == resourceID {
					if pool.Monitor.ProvisioningStatus != "" {
						return pool.Monitor, pool.Monitor.ProvisioningStatus, nil
					}
				}
			}
			monitor, err := neutronmonitors.Get(lbClient, resourceID).Extract()
			return monitor, "ACTIVE", err

		case "member":
			for _, pool := range statuses.Loadbalancer.Pools {
				for _, member := range pool.Members {
					if member.ID == resourceID {
						if member.ProvisioningStatus != "" {
							return member, member.ProvisioningStatus, nil
						}
					}
				}
			}
			member, err := pools.GetMember(lbClient, parentID, resourceID).Extract()
			return member, "ACTIVE", err

		case "l7policy":
			for _, listener := range statuses.Loadbalancer.Listeners {
				for _, l7policy := range listener.L7Policies {
					if l7policy.ID == resourceID {
						if l7policy.ProvisioningStatus != "" {
							return l7policy, l7policy.ProvisioningStatus, nil
						}
					}
				}
			}
			l7policy, err := l7policies.Get(lbClient, resourceID).Extract()
			return l7policy, "ACTIVE", err

		case "l7rule":
			for _, listener := range statuses.Loadbalancer.Listeners {
				for _, l7policy := range listener.L7Policies {
					for _, l7rule := range l7policy.Rules {
						if l7rule.ID == resourceID {
							if l7rule.ProvisioningStatus != "" {
								return l7rule, l7rule.ProvisioningStatus, nil
							}
						}
					}
				}
			}
			l7Rule, err := l7policies.GetRule(lbClient, parentID, resourceID).Extract()
			return l7Rule, "ACTIVE", err
		}

		return nil, "", fmt.Errorf("An unexpected error occurred querying the status of %s %s by loadbalancer %s", resourceType, resourceID, lbID)
	}
}

func resourceLBV2L7PolicyRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, l7policy *l7policies.L7Policy) resource.StateRefreshFunc {
	if l7policy.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(lbSkipLBStatuses, status) {
				return lb, status, nil
			}

			l7policy, err := l7policies.Get(lbClient, l7policy.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return l7policy, l7policy.ProvisioningStatus, nil
		}
	}

	return resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7policy", l7policy.ID, "")
}

func waitForLBV2L7Policy(lbClient *gophercloud.ServiceClient, parentListener *listeners.Listener, l7policy *l7policies.L7Policy, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7policy %s to become %s.", l7policy.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("Unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2L7PolicyRefreshFunc(lbClient, lbID, l7policy),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for l7policy %s to become %s: %s", l7policy.ID, target, err)
	}

	return nil
}

func getListenerIDForL7Policy(lbClient *gophercloud.ServiceClient, id string) (string, error) {
	log.Printf("[DEBUG] Trying to get Listener ID associated with the %s L7 Policy ID", id)
	lbsPages, err := loadbalancers.List(lbClient, loadbalancers.ListOpts{}).AllPages()
	if err != nil {
		return "", fmt.Errorf("No Load Balancers were found: %s", err)
	}

	lbs, err := loadbalancers.ExtractLoadBalancers(lbsPages)
	if err != nil {
		return "", fmt.Errorf("Unable to extract Load Balancers list: %s", err)
	}

	for _, lb := range lbs {
		statuses, err := loadbalancers.GetStatuses(lbClient, lb.ID).Extract()
		if err != nil {
			return "", fmt.Errorf("Failed to get Load Balancer statuses: %s", err)
		}
		for _, listener := range statuses.Loadbalancer.Listeners {
			for _, l7policy := range listener.L7Policies {
				if l7policy.ID == id {
					return listener.ID, nil
				}
			}
		}
	}

	return "", fmt.Errorf("Unable to find Listener ID associated with the %s L7 Policy ID", id)
}

func resourceLBV2L7RuleRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, l7policyID string, l7rule *l7policies.Rule) resource.StateRefreshFunc {
	if l7rule.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(lbSkipLBStatuses, status) {
				return lb, status, nil
			}

			l7rule, err := l7policies.GetRule(lbClient, l7policyID, l7rule.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return l7rule, l7rule.ProvisioningStatus, nil
		}
	}

	return resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7rule", l7rule.ID, l7policyID)
}

func waitForLBV2L7Rule(lbClient *gophercloud.ServiceClient, parentListener *listeners.Listener, parentL7policy *l7policies.L7Policy, l7rule *l7policies.Rule, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7rule %s to become %s.", l7rule.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("Unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2L7RuleRefreshFunc(lbClient, lbID, parentL7policy.ID, l7rule),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for l7rule %s to become %s: %s", l7rule.ID, target, err)
	}

	return nil
}

func flattenLBPoolPersistenceV2(p pools.SessionPersistence) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":        p.Type,
			"cookie_name": p.CookieName,
		},
	}
}
