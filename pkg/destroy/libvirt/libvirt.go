//go:build libvirt
// +build libvirt

package libvirt

import (
	"fmt"
	"net/url"
	"strings"

	libvirt "github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
)

// filterFunc allows filtering based on names.
// returns true, when the name should be handled.
type filterFunc func(name string) bool

// ClusterIDPrefixFilter returns true for names
// that are prefixed with clusterid.
// `clusterid` cannot be empty.
var ClusterIDPrefixFilter = func(clusterid string) filterFunc {
	if clusterid == "" {
		panic("clusterid cannot be empty")
	}
	return func(name string) bool {
		return strings.HasPrefix(name, clusterid)
	}
}

// AlwaysTrueFilter returns true for all
// names except `default`.
var AlwaysTrueFilter = func() filterFunc {
	return func(name string) bool {
		return name != "default"
	}
}

// deleteFunc is the interface a function needs to implement to be delete resources.
type deleteFunc func(conn *libvirt.Libvirt, filter filterFunc, logger logrus.FieldLogger) error

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	LibvirtURI string
	Filter     filterFunc
	Logger     logrus.FieldLogger
}

// New returns libvirt Uninstaller from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		LibvirtURI: metadata.ClusterPlatformMetadata.Libvirt.URI,
		Filter:     ClusterIDPrefixFilter(metadata.InfraID),
		Logger:     logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	uri, err := url.Parse(o.LibvirtURI)
	if err != nil {
		return nil, err
	}

	virt, err := libvirt.ConnectToURI(uri)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := virt.Disconnect(); err != nil {
			if o.Logger != nil {
				o.Logger.Warn("failed to disconnect from libvirt", err)
			}
		}
	}()

	for _, del := range []deleteFunc{
		deleteDomains,
		deleteNetwork,
		deleteStoragePool,
	} {
		err = del(virt, o.Filter, o.Logger)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// deleteDomains calls deleteDomainsSinglePass until it finds no
// matching domains.  This guards against the machine-API launching
// additional nodes after the initial list call.  We continue deleting
// domains until we either hit an error or we have a list call with no
// matching domains.
func deleteDomains(virt *libvirt.Libvirt, filter filterFunc, logger logrus.FieldLogger) error {
	logger.Debug("Deleting libvirt domains")
	var err error
	nothingToDelete := false
	for !nothingToDelete {
		nothingToDelete, err = deleteDomainsSinglePass(virt, filter, logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteDomainsSinglePass(virt *libvirt.Libvirt, filter filterFunc, logger logrus.FieldLogger) (nothingToDelete bool, err error) {
	domains, _, err := virt.ConnectListAllDomains(1, 0)
	if err != nil {
		return false, fmt.Errorf("list domains: %w", err)
	}

	nothingToDelete = true
	for _, domain := range domains {
		if !filter(domain.Name) {
			continue
		}

		nothingToDelete = false
		dState, _, err := virt.DomainGetState(domain, 0)
		if err != nil {
			return false, fmt.Errorf("get domain state %d: %w", domain.Name, err)
		}

		if libvirt.DomainState(dState) != libvirt.DomainShutoff && libvirt.DomainState(dState) != libvirt.DomainShutdown {
			if err := virt.DomainDestroy(domain); err != nil {
				return false, fmt.Errorf("destroy domain %q: %w", domain.Name, err)
			}
		}
		if err := virt.DomainUndefineFlags(domain, libvirt.DomainUndefineNvram); err != nil {
			if e, ok := err.(libvirt.Error); ok && (libvirt.ErrorNumber(e.Code) == libvirt.ErrNoSupport || libvirt.ErrorNumber(e.Code) == libvirt.ErrInvalidArg) {
				logger.WithField("domain", domain.Name).Info("libvirt does not support undefine flags: will try again without flags")
				if err := virt.DomainUndefine(domain); err != nil {
					return false, fmt.Errorf("could not undefine libvirt domain: %q: %w", domain.Name, err)
				}
			} else {
				return false, fmt.Errorf("could not undefine libvirt domain %q with flags: %w", domain.Name, err)
			}
		}
		logger.WithField("domain", domain.Name).Info("Deleted domain")
	}

	return nothingToDelete, nil
}

func deleteStoragePool(virt *libvirt.Libvirt, filter filterFunc, logger logrus.FieldLogger) error {
	logger.Debug("Deleting libvirt volumes")

	pools, _, err := virt.ConnectListAllStoragePools(1, 0)
	if err != nil {
		return fmt.Errorf("list storage pools: %w", err)
	}

	for _, pool := range pools {
		// pool name that returns true from filter
		if !filter(pool.Name) {
			continue
		}

		// delete all vols that return true from filter.
		vols, _, err := virt.StoragePoolListAllVolumes(pool, 1, 0)
		if err != nil {
			return fmt.Errorf("list volumes in %q: %w", pool.Name, err)
		}

		for _, vol := range vols {
			if err := virt.StorageVolDelete(vol, 0); err != nil {
				return fmt.Errorf("delete volume %q from %q: %w", vol.Name, pool.Name, err)
			}
			logger.WithField("volume", vol.Name).Info("Deleted volume")
		}

		// blow away entire pool.
		if err := virt.StoragePoolDestroy(pool); err != nil {
			return fmt.Errorf("destroy pool %q: %w", pool.Name, err)
		}

		if err := virt.StoragePoolDelete(pool, 0); err != nil {
			return fmt.Errorf("delete pool %q: %w", pool.Name, err)
		}

		if err := virt.StoragePoolUndefine(pool); err != nil {
			return fmt.Errorf("undefine pool %q: %w", pool.Name, err)
		}
		logger.WithField("pool", pool.Name).Info("Deleted pool")
	}

	return nil
}

func deleteNetwork(virt *libvirt.Libvirt, filter filterFunc, logger logrus.FieldLogger) error {
	logger.Debug("Deleting libvirt network")

	networks, _, err := virt.ConnectListAllNetworks(1, 0)
	if err != nil {
		return fmt.Errorf("list networks: %w", err)
	}

	for _, network := range networks {
		if !filter(network.Name) {
			continue
		}

		if err := virt.NetworkDestroy(network); err != nil {
			return fmt.Errorf("destroy network %q: %w", network.Name, err)
		}

		if err := virt.NetworkUndefine(network); err != nil {
			return fmt.Errorf("undefine network %q: %w", network.Name, err)
		}
		logger.WithField("network", network.Name).Info("Deleted network")
	}
	return nil
}
