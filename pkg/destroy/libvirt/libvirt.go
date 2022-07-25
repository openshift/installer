//go:build libvirt
// +build libvirt

package libvirt

import (
	"strings"

	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
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
type deleteFunc func(conn *libvirt.Connect, filter filterFunc, logger logrus.FieldLogger) error

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
	conn, err := libvirt.NewConnect(o.LibvirtURI)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to Libvirt daemon")
	}

	for _, del := range []deleteFunc{
		deleteDomains,
		deleteNetwork,
		deleteStoragePool,
	} {
		err = del(conn, o.Filter, o.Logger)
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
func deleteDomains(conn *libvirt.Connect, filter filterFunc, logger logrus.FieldLogger) error {
	logger.Debug("Deleting libvirt domains")
	var err error
	nothingToDelete := false
	for !nothingToDelete {
		nothingToDelete, err = deleteDomainsSinglePass(conn, filter, logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteDomainsSinglePass(conn *libvirt.Connect, filter filterFunc, logger logrus.FieldLogger) (nothingToDelete bool, err error) {
	domains, err := conn.ListAllDomains(0)
	if err != nil {
		return false, errors.Wrap(err, "list domains")
	}

	nothingToDelete = true
	for _, domain := range domains {
		defer domain.Free()
		dName, err := domain.GetName()
		if err != nil {
			return false, errors.Wrap(err, "get domain name")
		}
		if !filter(dName) {
			continue
		}

		nothingToDelete = false
		dState, _, err := domain.GetState()
		if err != nil {
			return false, errors.Wrapf(err, "get domain state %d", dName)
		}

		if dState != libvirt.DOMAIN_SHUTOFF && dState != libvirt.DOMAIN_SHUTDOWN {
			if err := domain.Destroy(); err != nil {
				return false, errors.Wrapf(err, "destroy domain %q", dName)
			}
		}
		if err := domain.UndefineFlags(libvirt.DOMAIN_UNDEFINE_NVRAM); err != nil {
			if e := err.(libvirt.Error); e.Code == libvirt.ERR_NO_SUPPORT || e.Code == libvirt.ERR_INVALID_ARG {
				logger.WithField("domain", dName).Info("libvirt does not support undefine flags: will try again without flags")
				if err := domain.Undefine(); err != nil {
					return false, errors.Wrapf(err, "could not undefine libvirt domain: %q", dName)
				}
			} else {
				return false, errors.Wrapf(err, "could not undefine libvirt domain %q with flags", dName)
			}
		}
		logger.WithField("domain", dName).Info("Deleted domain")
	}

	return nothingToDelete, nil
}

func deleteStoragePool(conn *libvirt.Connect, filter filterFunc, logger logrus.FieldLogger) error {
	logger.Debug("Deleting libvirt volumes")

	pools, err := conn.ListStoragePools()
	if err != nil {
		return errors.Wrap(err, "list storage pools")
	}

	for _, pname := range pools {
		// pool name that returns true from filter
		if !filter(pname) {
			continue
		}

		pool, err := conn.LookupStoragePoolByName(pname)
		if err != nil {
			return errors.Wrapf(err, "get storage pool %q", pname)
		}
		defer pool.Free()

		// delete all vols that return true from filter.
		vols, err := pool.ListAllStorageVolumes(0)
		if err != nil {
			return errors.Wrapf(err, "list volumes in %q", pname)
		}

		for _, vol := range vols {
			defer vol.Free()
			vName, err := vol.GetName()
			if err != nil {
				return errors.Wrapf(err, "get volume names in %q", pname)
			}
			if err := vol.Delete(0); err != nil {
				return errors.Wrapf(err, "delete volume %q from %q", vName, pname)
			}
			logger.WithField("volume", vName).Info("Deleted volume")
		}

		// blow away entire pool.
		if err := pool.Destroy(); err != nil {
			return errors.Wrapf(err, "destroy pool %q", pname)
		}

		if err := pool.Delete(0); err != nil {
			return errors.Wrapf(err, "delete pool %q", pname)
		}

		if err := pool.Undefine(); err != nil {
			return errors.Wrapf(err, "undefine pool %q", pname)
		}
		logger.WithField("pool", pname).Info("Deleted pool")
	}

	return nil
}

func deleteNetwork(conn *libvirt.Connect, filter filterFunc, logger logrus.FieldLogger) error {
	logger.Debug("Deleting libvirt network")

	networks, err := conn.ListNetworks()
	if err != nil {
		return errors.Wrap(err, "list networks")
	}

	for _, nName := range networks {
		if !filter(nName) {
			continue
		}
		network, err := conn.LookupNetworkByName(nName)
		if err != nil {
			return errors.Wrapf(err, "get network %q", nName)
		}
		defer network.Free()

		if err := network.Destroy(); err != nil {
			return errors.Wrapf(err, "destroy network %q", nName)
		}

		if err := network.Undefine(); err != nil {
			return errors.Wrapf(err, "undefine network %q", nName)
		}
		logger.WithField("network", nName).Info("Deleted network")
	}
	return nil
}
