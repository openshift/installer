package ovirtclient

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
)

// AffinityGroupClient describes the methods required for working with affinity groups.
type AffinityGroupClient interface {
	// CreateAffinityGroup creates an affinity group with the specified parameters.
	CreateAffinityGroup(
		clusterID ClusterID,
		name string,
		params CreateAffinityGroupOptionalParams,
		retries ...RetryStrategy,
	) (
		AffinityGroup,
		error,
	)
	// ListAffinityGroups returns a list of all affinity groups in the oVirt engine.
	ListAffinityGroups(clusterID ClusterID, retries ...RetryStrategy) ([]AffinityGroup, error)
	// GetAffinityGroup returns a specific affinity group based on its ID. An error is returned if the affinity label
	// doesn't exist.
	GetAffinityGroup(clusterID ClusterID, id AffinityGroupID, retries ...RetryStrategy) (AffinityGroup, error)
	// GetAffinityGroupByName returns an affinity group by name.
	GetAffinityGroupByName(clusterID ClusterID, name string, retries ...RetryStrategy) (AffinityGroup, error)
	// RemoveAffinityGroup removes the affinity group specified.
	RemoveAffinityGroup(clusterID ClusterID, id AffinityGroupID, retries ...RetryStrategy) error

	AddVMToAffinityGroup(clusterID ClusterID, vmID VMID, agID AffinityGroupID, retries ...RetryStrategy) error
	RemoveVMFromAffinityGroup(clusterID ClusterID, vmID VMID, agID AffinityGroupID, retries ...RetryStrategy) error
}

// CreateAffinityGroupOptionalParams is a list of optional parameters that can be passed for affinity group creation.
type CreateAffinityGroupOptionalParams interface {
	// Priority returns the affinity group priority that should be applied, or nil if no explicit priority should be
	// applied.
	Priority() *AffinityGroupPriority
	// HostsRule returns a hosts rule that should be applied, or nil if no hosts rule should explicitly be applied.
	HostsRule() AffinityHostsRule
	// VMsRule returns a VMs rule that should be applied, or nil if no VMs rule should explicitly be applied.
	VMsRule() AffinityVMsRule
	// Enforcing returns if the affinity group should be enforced.
	Enforcing() *bool
	// Description returns the description for the affinity group.
	Description() string
}

// BuildableCreateAffinityGroupOptionalParams is a buildable version of CreateAffinityGroupOptionalParams.
type BuildableCreateAffinityGroupOptionalParams interface {
	CreateAffinityGroupOptionalParams
	// WithPriority adds a priority to the affinity group.
	WithPriority(priority AffinityGroupPriority) (BuildableCreateAffinityGroupOptionalParams, error)

	// MustWithPriority is equivalent to WithPriority, but panics instead of returning an error.
	MustWithPriority(priority AffinityGroupPriority) BuildableCreateAffinityGroupOptionalParams

	WithHostsRule(rule AffinityHostsRule) (BuildableCreateAffinityGroupOptionalParams, error)
	MustWithHostsRule(rule AffinityHostsRule) BuildableCreateAffinityGroupOptionalParams

	WithHostsRuleParameters(
		enabled bool,
		affinity Affinity,
		enforcing bool,
	) (BuildableCreateAffinityGroupOptionalParams, error)
	MustWithHostsRuleParameters(
		enabled bool,
		affinity Affinity,
		enforcing bool,
	) BuildableCreateAffinityGroupOptionalParams

	WithVMsRule(rule AffinityVMsRule) (BuildableCreateAffinityGroupOptionalParams, error)
	MustWithVMsRule(rule AffinityVMsRule) BuildableCreateAffinityGroupOptionalParams

	WithVMsRuleParameters(enabled bool, affinity Affinity, enforcing bool) (
		BuildableCreateAffinityGroupOptionalParams,
		error,
	)
	MustWithVMsRuleParameters(
		enabled bool,
		affinity Affinity,
		enforcing bool,
	) BuildableCreateAffinityGroupOptionalParams

	WithEnforcing(enforcing bool) (BuildableCreateAffinityGroupOptionalParams, error)
	MustWithEnforcing(enforcing bool) BuildableCreateAffinityGroupOptionalParams

	WithDescription(description string) (BuildableCreateAffinityGroupOptionalParams, error)
	MustWithDescription(description string) BuildableCreateAffinityGroupOptionalParams
}

// CreateAffinityGroupParams creates a buildable set of parameters for creating an affinity group.
func CreateAffinityGroupParams() BuildableCreateAffinityGroupOptionalParams {
	return &createAffinityGroupParams{}
}

// AffinityGroupPriority is a type alias for the type indicating affinity group priority.
type AffinityGroupPriority float64

// AffinityGroupID is the identifier for affinity groups.
type AffinityGroupID string

// Affinity signals if the affinity is positive (attracting VMs to each other) or negative (pushing VMs from each other
// to different hosts).
type Affinity bool

const (
	// AffinityPositive attracts VMs to each other, they are placed on the same host if enforcing is true, or are
	// attempted to place on the same host if possible in case enforcing is false.
	AffinityPositive Affinity = true
	// AffinityNegative pushes VMs from each other, they are placed on different hosts if enforcing is true, or are
	// attempted to place on different hosts if possible in case enforcing is false.
	AffinityNegative Affinity = false
)

// AffinityGroupData contains the base data for the AffinityGroup.
type AffinityGroupData interface {
	// ID returns the oVirt identifier of the affinity group.
	ID() AffinityGroupID
	// Name is the user-readable oVirt name of the affinity group.
	Name() string
	// Description returns the description of the affinity group.
	Description() string
	// ClusterID is the identifier of the cluster this affinity group belongs to.
	ClusterID() ClusterID
	// Priority indicates in which order the affinity groups should be evaluated.
	Priority() AffinityGroupPriority
	// Enforcing indicates if the deployment should fail if the affinity group cannot be respected.
	Enforcing() bool
	// HostsRule contains the rules for hosts.
	HostsRule() AffinityHostsRule
	// VMsRule contains the rule for the virtual machines.
	VMsRule() AffinityVMsRule
	// VMIDs returns the list of current virtual machine IDs assigned to this affinity group.
	VMIDs() []VMID
}

// AffinityGroup labels virtual machines, so they run / don't run on the same host.
type AffinityGroup interface {
	AffinityGroupData

	// Cluster fetches the cluster this affinity group belongs to.
	Cluster(retries ...RetryStrategy) (Cluster, error)
	// Remove removes the current affinity group.
	Remove(retries ...RetryStrategy) error

	// AddVM adds the specified VM to the current affinity group.
	AddVM(id VMID, retries ...RetryStrategy) error
	// RemoveVM removes the specified VM from the current affinity group.
	RemoveVM(id VMID, retries ...RetryStrategy) error
}

// AffinityRule is a rule for either hosts or virtual machines.
type AffinityRule interface {
	// Enabled indicates if the rule is enabled.
	Enabled() bool
	// Affinity indicates if the affinity is positive (attracting VMs) or negative (pushes VMs from each other).
	Affinity() Affinity
	// Enforcing indicates if the deployment should fail if the affinity group cannot be respected.
	Enforcing() bool
}

// AffinityHostsRule is an alias for hosts rules to avoid mixups.
type AffinityHostsRule AffinityRule

// AffinityVMsRule is an alias for VM rules to avoid mixups.
type AffinityVMsRule AffinityRule

type affinityRule struct {
	enabled   bool
	affinity  Affinity
	enforcing bool
}

func (a affinityRule) Enabled() bool {
	return a.enabled
}

func (a affinityRule) Affinity() Affinity {
	return a.affinity
}

func (a affinityRule) Enforcing() bool {
	return a.enforcing
}

type affinityGroup struct {
	client Client

	id          AffinityGroupID
	name        string
	description string
	clusterID   ClusterID
	priority    AffinityGroupPriority
	enforcing   bool

	hostsRule AffinityRule
	vmsRule   AffinityRule
	vmids     []VMID
}

func (a affinityGroup) Description() string {
	return a.description
}

func (a affinityGroup) hasVM(id VMID) bool {
	for _, vmid := range a.vmids {
		if vmid == id {
			return true
		}
	}
	return false
}

func (a affinityGroup) AddVM(id VMID, retries ...RetryStrategy) error {
	return a.client.AddVMToAffinityGroup(a.clusterID, id, a.id, retries...)
}

func (a affinityGroup) RemoveVM(id VMID, retries ...RetryStrategy) error {
	return a.client.RemoveVMFromAffinityGroup(a.clusterID, id, a.id, retries...)
}

func (a affinityGroup) Remove(retries ...RetryStrategy) error {
	return a.client.RemoveAffinityGroup(a.clusterID, a.id, retries...)
}

func (a affinityGroup) ClusterID() ClusterID {
	return a.clusterID
}

func (a affinityGroup) Priority() AffinityGroupPriority {
	return a.priority
}

func (a affinityGroup) HostsRule() AffinityHostsRule {
	return a.hostsRule
}

func (a affinityGroup) VMsRule() AffinityVMsRule {
	return a.vmsRule
}

func (a affinityGroup) Cluster(retries ...RetryStrategy) (Cluster, error) {
	return a.client.GetCluster(a.clusterID, retries...)
}

func (a affinityGroup) ID() AffinityGroupID {
	return a.id
}

func (a affinityGroup) Name() string {
	return a.name
}

func (a affinityGroup) Enforcing() bool {
	return a.enforcing
}

func (a affinityGroup) VMIDs() []VMID {
	return a.vmids
}

func convertSDKAffinityGroupID(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error {
	result.id = AffinityGroupID(sdkObject.MustId())
	return nil
}

func convertSDKAffinityGroupName(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error {
	name, ok := sdkObject.Name()
	if !ok {
		return newFieldNotFound("affinity group", "name")
	}
	result.name = name
	return nil
}

func convertSDKAffinityGroupDescription(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error {
	if description, ok := sdkObject.Description(); ok {
		result.description = description
	}
	return nil
}

func convertSDKAffinityGroupCluster(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error {
	cluster, ok := sdkObject.Cluster()
	if !ok {
		return newFieldNotFound("affinity group", "cluster")
	}
	clusterID, ok := cluster.Id()
	if !ok {
		return newFieldNotFound("cluster in affinity group", "id")
	}
	result.clusterID = ClusterID(clusterID)
	return nil
}

func convertSDKAffinityGroupEnforcing(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error {
	enforcing, ok := sdkObject.Enforcing()
	if !ok {
		return newFieldNotFound("affinity group", "enforcing")
	}
	result.enforcing = enforcing
	return nil
}

func convertSDKAffinityGroupPriority(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error {
	priority, ok := sdkObject.Priority()
	if !ok {
		return newFieldNotFound("affinity group", "priority")
	}
	result.priority = AffinityGroupPriority(priority)
	return nil
}

func convertSDKAffinityGroupHostsRule(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error {
	hostsRule, ok := sdkObject.HostsRule()
	if !ok {
		return newFieldNotFound("affinity group", "hosts rule")
	}
	convertedSDKHostsRule, err := convertSDKAffinityRule(hostsRule)
	if err != nil {
		return err
	}
	result.hostsRule = convertedSDKHostsRule
	return nil
}

func convertSDKAffinityGroupVMsRule(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error {
	vmsRule, ok := sdkObject.VmsRule()
	if !ok {
		return newFieldNotFound("affinity group", "VMs rule")
	}
	convertedSDKVMSRule, err := convertSDKAffinityRule(vmsRule)
	if err != nil {
		return err
	}
	result.vmsRule = convertedSDKVMSRule
	return nil
}

func convertSDKAffinityGroupVMsList(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error {
	vmsList, ok := sdkObject.Vms()
	if !ok {
		return newFieldNotFound("affinity group", "VMs list")
	}
	convertedVMIDs := make([]VMID, len(vmsList.Slice()))
	for i, vm := range vmsList.Slice() {
		vmid, ok := vm.Id()
		if !ok {
			return newFieldNotFound("VM on affinity group", "id")
		}
		convertedVMIDs[i] = VMID(vmid)
	}
	result.vmids = convertedVMIDs
	return nil
}

func convertSDKAffinityGroup(sdkObject *ovirtsdk.AffinityGroup, o *oVirtClient) (AffinityGroup, error) {
	result := &affinityGroup{
		client: o,
	}
	converters := []func(sdkObject *ovirtsdk.AffinityGroup, result *affinityGroup) error{
		convertSDKAffinityGroupID,
		convertSDKAffinityGroupName,
		convertSDKAffinityGroupDescription,
		convertSDKAffinityGroupCluster,
		convertSDKAffinityGroupEnforcing,
		convertSDKAffinityGroupPriority,
		convertSDKAffinityGroupHostsRule,
		convertSDKAffinityGroupVMsRule,
		convertSDKAffinityGroupVMsList,
	}
	for _, converter := range converters {
		if err := converter(sdkObject, result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func convertSDKAffinityRule(sdk *ovirtsdk.AffinityRule) (*affinityRule, error) {
	enabled, ok := sdk.Enabled()
	if !ok {
		return nil, newFieldNotFound("affinity rule", "enabled")
	}
	affinity, ok := sdk.Positive()
	if !ok {
		return nil, newFieldNotFound("affinity rule", "positive")
	}
	enforcing, ok := sdk.Enforcing()
	if !ok {
		return nil, newFieldNotFound("affinity rule", "enforcing")
	}
	return &affinityRule{
		enabled:   enabled,
		affinity:  Affinity(affinity),
		enforcing: enforcing,
	}, nil
}

type createAffinityGroupParams struct {
	priority    *AffinityGroupPriority
	hostsRule   AffinityHostsRule
	vmsRule     AffinityVMsRule
	enforcing   *bool
	description string
}

func (c *createAffinityGroupParams) Description() string {
	return c.description
}

func (c *createAffinityGroupParams) WithDescription(description string) (
	BuildableCreateAffinityGroupOptionalParams,
	error,
) {
	c.description = description
	return c, nil
}

func (c *createAffinityGroupParams) MustWithDescription(description string) BuildableCreateAffinityGroupOptionalParams {
	builder, err := c.WithDescription(description)
	if err != nil {
		panic(err)
	}
	return builder
}

func (c *createAffinityGroupParams) Enforcing() *bool {
	return c.enforcing
}

func (c *createAffinityGroupParams) WithEnforcing(enforcing bool) (BuildableCreateAffinityGroupOptionalParams, error) {
	c.enforcing = &enforcing
	return c, nil
}

func (c *createAffinityGroupParams) MustWithEnforcing(enforcing bool) BuildableCreateAffinityGroupOptionalParams {
	builder, err := c.WithEnforcing(enforcing)
	if err != nil {
		panic(err)
	}
	return builder
}

func (c *createAffinityGroupParams) WithHostsRule(
	rule AffinityHostsRule,
) (BuildableCreateAffinityGroupOptionalParams, error) {
	c.hostsRule = rule
	return c, nil
}

func (c *createAffinityGroupParams) MustWithHostsRule(
	rule AffinityHostsRule,
) BuildableCreateAffinityGroupOptionalParams {
	builder, err := c.WithHostsRule(rule)
	if err != nil {
		panic(err)
	}
	return builder
}

func (c *createAffinityGroupParams) WithHostsRuleParameters(
	enabled bool,
	affinity Affinity,
	enforcing bool,
) (BuildableCreateAffinityGroupOptionalParams, error) {
	c.hostsRule = &affinityRule{
		enabled:   enabled,
		affinity:  affinity,
		enforcing: enforcing,
	}
	return c, nil
}

func (c *createAffinityGroupParams) MustWithHostsRuleParameters(
	enabled bool,
	affinity Affinity,
	enforcing bool,
) BuildableCreateAffinityGroupOptionalParams {
	builder, err := c.WithHostsRuleParameters(enabled, affinity, enforcing)
	if err != nil {
		panic(err)
	}
	return builder
}

func (c *createAffinityGroupParams) WithVMsRule(
	rule AffinityVMsRule,
) (BuildableCreateAffinityGroupOptionalParams, error) {
	c.vmsRule = rule
	return c, nil
}

func (c *createAffinityGroupParams) MustWithVMsRule(rule AffinityVMsRule) BuildableCreateAffinityGroupOptionalParams {
	builder, err := c.WithVMsRule(rule)
	if err != nil {
		panic(err)
	}
	return builder
}

func (c *createAffinityGroupParams) WithVMsRuleParameters(
	enabled bool,
	affinity Affinity,
	enforcing bool,
) (BuildableCreateAffinityGroupOptionalParams, error) {
	c.vmsRule = &affinityRule{
		enabled:   enabled,
		affinity:  affinity,
		enforcing: enforcing,
	}
	return c, nil
}

func (c *createAffinityGroupParams) MustWithVMsRuleParameters(enabled bool, affinity Affinity, enforcing bool) BuildableCreateAffinityGroupOptionalParams {
	builder, err := c.WithVMsRuleParameters(enabled, affinity, enforcing)
	if err != nil {
		panic(err)
	}
	return builder
}

func (c *createAffinityGroupParams) Priority() *AffinityGroupPriority {
	return c.priority
}

func (c *createAffinityGroupParams) HostsRule() AffinityHostsRule {
	return c.hostsRule
}

func (c *createAffinityGroupParams) VMsRule() AffinityVMsRule {
	return c.vmsRule
}

func (c *createAffinityGroupParams) WithPriority(priority AffinityGroupPriority) (
	BuildableCreateAffinityGroupOptionalParams,
	error,
) {
	c.priority = &priority
	return c, nil
}

func (c *createAffinityGroupParams) MustWithPriority(priority AffinityGroupPriority) BuildableCreateAffinityGroupOptionalParams {
	c.priority = &priority
	return c
}
