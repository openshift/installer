package vsphere

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"

	"github.com/openshift/installer/pkg/types/vsphere"
)

// NetworkNameMap contains the vCenter cluster name, resource pools and network names.
type NetworkNameMap struct {
	Cluster       string
	ResourcePools map[string]*object.ResourcePool
	NetworkNames  map[string]string
}

// VCenterContext maintains context of known vCenters to be used in CAPI manifest reconciliation.
type VCenterContext struct {
	VCenter           string
	Datacenters       []string
	ClusterNetworkMap map[string]NetworkNameMap
}

// VCenterCredential contains the vCenter username and password.
type VCenterCredential struct {
	Username string
	Password string
}

// Metadata holds vcenter stuff.
type Metadata struct {
	sessions    map[string]*session.Session
	credentials map[string]*session.Params

	VCenterContexts map[string]VCenterContext

	VCenterCredentials map[string]VCenterCredential

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata() *Metadata {
	return &Metadata{
		sessions:           make(map[string]*session.Session),
		credentials:        make(map[string]*session.Params),
		VCenterContexts:    make(map[string]VCenterContext),
		VCenterCredentials: make(map[string]VCenterCredential),
	}
}

// AddCredentials creates a session param from the vCenter server, username and password
// to the Credentials Map.
func (m *Metadata) AddCredentials(server, username, password string) *session.Params {
	if _, ok := m.VCenterCredentials[server]; !ok {
		m.VCenterCredentials[server] = VCenterCredential{
			Username: username,
			Password: password,
		}
	}

	// m.credentials is not stored in the json state file - there is no real reason to do this
	// but upon returning to AddCredentials (create manifest, create cluster) the credentials map is
	// nil, re-make it.
	if m.credentials == nil {
		m.credentials = make(map[string]*session.Params)
	}

	if _, ok := m.credentials[server]; !ok {
		m.credentials[server] = session.NewParams().WithServer(server).WithUserInfo(username, password)
	}
	return m.credentials[server]
}

// Session returns a session from unlockedSession based on the server (vCenter URL).
func (m *Metadata) Session(ctx context.Context, server string) (*session.Session, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// m.sessions is not stored in the json state file - there is no real reason to do this
	// but upon returning to Session (create manifest, create cluster) the sessions map is
	// nil, re-make it.
	if m.sessions == nil {
		m.sessions = make(map[string]*session.Session)
	}

	return m.unlockedSession(ctx, server)
}

func (m *Metadata) unlockedSession(ctx context.Context, server string) (*session.Session, error) {
	var err error
	var ok bool
	var params *session.Params

	if params, ok = m.credentials[server]; !ok {
		if creds, ok := m.VCenterCredentials[server]; ok {
			params = m.AddCredentials(server, creds.Username, creds.Password)
		} else {
			return nil, fmt.Errorf("credentials for %s not found", server)
		}
	}

	// if nil we haven't created a session
	if _, ok := m.sessions[server]; ok {
		// is the session still valid? if not re-run GetOrCreate.
		if !m.sessions[server].Valid() {
			m.sessions[server], err = session.GetOrCreate(ctx, params)
			if err != nil {
				return nil, err
			}
		}
		return m.sessions[server], nil
	}

	// If we have gotten here there is no session for the server name, create.
	m.sessions[server], err = session.GetOrCreate(ctx, params)
	return m.sessions[server], err
}

// unwrapToSoapFault is required because soapErrorFaul is not exported
// and are unable to use errors.As()
// https://github.com/vmware/govmomi/blob/main/vim25/soap/error.go#L38
func unwrapToSoapFault(err error) error {
	if err != nil {
		if soapFault := soap.IsSoapFault(err); !soapFault {
			return unwrapToSoapFault(errors.Unwrap(err))
		}
		return err
	}
	return err
}

// Networks populates VCenterContext and the ClusterNetworkMap based on the vCenter server url and the FailureDomains.
func (m *Metadata) Networks(ctx context.Context, vcenter vsphere.VCenter, failureDomains []vsphere.FailureDomain) error {
	_, err := m.Session(ctx, vcenter.Server)
	if err != nil {
		// Defense against potential issues with assisted installer
		if soapErr := unwrapToSoapFault(err); soapErr != nil {
			soapFault := soap.ToSoapFault(soapErr)
			// The assisted installer provides bogus username and password
			// values. Only return the soap error (fault) if it matches incorrect username or password.
			if strings.Contains(soapFault.String, "Cannot complete login due to an incorrect user name or password") {
				return soapErr
			}
		}
		// if soapErr is nil then this is not a SOAP fault, return err
		return err
	}

	m.VCenterContexts[vcenter.Server] = VCenterContext{
		VCenter:           vcenter.Server,
		Datacenters:       vcenter.Datacenters,
		ClusterNetworkMap: make(map[string]NetworkNameMap),
	}

	for _, fd := range failureDomains {
		if fd.Server != vcenter.Server {
			continue
		}
		if err := m.populateNetworks(ctx, fd); err != nil {
			return fmt.Errorf("unable to retrieve network names: %w", err)
		}
	}

	return nil
}

func (m *Metadata) populateNetworks(ctx context.Context, failureDomain vsphere.FailureDomain) error {
	var ok bool
	var sess *session.Session

	if sess, ok = m.sessions[failureDomain.Server]; !ok {
		return fmt.Errorf("unable to find session")
	}

	finder := sess.Finder
	clusterPath := failureDomain.Topology.ComputeCluster

	clusterRef, err := finder.ClusterComputeResource(ctx, clusterPath)
	if err != nil {
		return fmt.Errorf("unable to retrieve compute cluster: %w", err)
	}

	rpFindPath := path.Join(clusterRef.InventoryPath, "...")
	pools, err := finder.ResourcePoolList(ctx, rpFindPath)
	if err != nil {
		return fmt.Errorf("unable to retrieve resource pools relative to compute cluster: %w", err)
	}

	for _, network := range failureDomain.Topology.Networks {
		clusterMap, present := m.VCenterContexts[failureDomain.Server].ClusterNetworkMap[clusterPath]
		if !present {
			clusterMap = NetworkNameMap{
				Cluster:       clusterPath,
				NetworkNames:  map[string]string{},
				ResourcePools: map[string]*object.ResourcePool{},
			}
			for _, pool := range pools {
				clusterMap.ResourcePools[path.Clean(pool.InventoryPath)] = pool
			}

			m.VCenterContexts[failureDomain.Server].ClusterNetworkMap[clusterPath] = clusterMap
		}

		networkPath := path.Join(clusterRef.InventoryPath, network)

		// Added check to confirm that the path to the network exists.
		networkRef, err := finder.Network(ctx, networkPath)
		if err != nil {
			return err
		}

		clusterMap.NetworkNames[network] = networkRef.GetInventoryPath()
	}
	return nil
}
