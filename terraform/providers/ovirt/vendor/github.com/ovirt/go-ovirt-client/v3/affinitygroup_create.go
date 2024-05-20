package ovirtclient

import (
	"fmt"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CreateAffinityGroup(
	clusterID ClusterID,
	name string,
	params CreateAffinityGroupOptionalParams,
	retries ...RetryStrategy,
) (result AffinityGroup, err error) {
	if params == nil {
		params = CreateAffinityGroupParams()
	}
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("creating affinity group in cluster %s", clusterID),
		o.logger,
		retries,
		func() error {
			agBuilder := ovirtsdk4.NewAffinityGroupBuilder().
				Name(name)
			if enforcing := params.Enforcing(); enforcing != nil {
				agBuilder.Enforcing(*enforcing)
			}
			if description := params.Description(); description != "" {
				agBuilder.Description(description)
			}
			if vmsRule := params.VMsRule(); vmsRule != nil {
				rule := ovirtsdk4.NewAffinityRuleBuilder()
				rule.Enabled(vmsRule.Enabled())
				rule.Positive(bool(vmsRule.Affinity()))
				rule.Enforcing(vmsRule.Enforcing())
				agBuilder.VmsRule(rule.MustBuild())
			}
			if hostsRule := params.HostsRule(); hostsRule != nil {
				rule := ovirtsdk4.NewAffinityRuleBuilder()
				rule.Enabled(hostsRule.Enabled())
				rule.Positive(bool(hostsRule.Affinity()))
				rule.Enforcing(hostsRule.Enforcing())
				agBuilder.HostsRule(rule.MustBuild())
			}
			addRequest := o.conn.
				SystemService().
				ClustersService().
				ClusterService(string(clusterID)).
				AffinityGroupsService().
				Add()
			addRequest.Group(
				agBuilder.MustBuild(),
			)
			response, err := addRequest.Send()
			if err != nil {
				return err
			}
			group, ok := response.Group()
			if !ok {
				return newFieldNotFound("add affinity group response", "group")
			}
			result, err = convertSDKAffinityGroup(group, o)
			return err
		},
	)
	return result, err
}

func (m *mockClient) CreateAffinityGroup(
	clusterID ClusterID,
	name string,
	params CreateAffinityGroupOptionalParams,
	_ ...RetryStrategy,
) (AffinityGroup, error) {
	if params == nil {
		params = CreateAffinityGroupParams()
	}

	priority := AffinityGroupPriority(1)
	if params.Priority() != nil {
		priority = *params.Priority()
	}
	enforcing := false
	if params.Enforcing() != nil {
		enforcing = *params.Enforcing()
	}
	hostsRule := params.HostsRule()
	if hostsRule == nil {
		hostsRule = &affinityRule{}
	}
	vmsRule := params.VMsRule()
	if vmsRule == nil {
		vmsRule = &affinityRule{}
	}
	ag := &affinityGroup{
		client:      m,
		id:          AffinityGroupID(m.GenerateUUID()),
		name:        name,
		description: params.Description(),
		clusterID:   clusterID,
		priority:    priority,
		enforcing:   enforcing,
		hostsRule:   hostsRule,
		vmsRule:     vmsRule,
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.affinityGroups[ag.ClusterID()]; !ok {
		return nil, newError(ENotFound, "Cluster %s not found.", ag.ClusterID())
	}

	m.affinityGroups[ag.ClusterID()][ag.id] = ag

	return ag, nil
}
