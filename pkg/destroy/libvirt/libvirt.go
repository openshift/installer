// +build libvirt_destroy

package libvirt

import (
	"strings"

	libvirt "github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/destroy"
	"github.com/openshift/installer/pkg/types"
)

// filterFunc allows filtering based on names.
// returns true, when the name should be handled.
type filterFunc func(name string) bool

// ClusterNamePrefixFilter returns true for names
// that are prefixed with clustername.
// `clustername` cannot be empty.
var ClusterNamePrefixFilter = func(clustername string) filterFunc {
	if clustername == "" {
		panic("clustername cannot be empty")
	}
	return func(name string) bool {
		return strings.HasPrefix(name, clustername)
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

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() error {
	conn, err := libvirt.NewConnect(o.LibvirtURI)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Libvirt daemon")
	}

	for _, del := range []deleteFunc{
		deleteDomains,
		deleteNetwork,
		deleteVolumes,
	} {
		err = del(conn, o.Filter, o.Logger)
		if err != nil {
			return err
		}
	}

	return nil
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
		if err := domain.Destroy(); err != nil {
			return false, errors.Wrapf(err, "destroy domain %q", dName)
		}

		if err := domain.Undefine(); err != nil {
			return false, errors.Wrapf(err, "undefine domain %q", dName)
		}
		logger.WithField("domain", dName).Info("Deleted domain")
	}

	return nothingToDelete, nil
}

func deleteVolumes(conn *libvirt.Connect, filter filterFunc, logger logrus.FieldLogger) error {
	logger.Debug("Deleting libvirt volumes")

	pools, err := conn.ListStoragePools()
	if err != nil {
		return errors.Wrap(err, "list storage pools")
	}

	tpool := "default"
	for _, pname := range pools {
		// pool name that returns true from filter, override default.
		if filter(pname) {
			tpool = pname
		}
	}
	pool, err := conn.LookupStoragePoolByName(tpool)
	if err != nil {
		return errors.Wrapf(err, "get storage pool %q", tpool)
	}
	defer pool.Free()

	switch tpool {
	case "default":
		// delete all vols that return true from filter.
		vols, err := pool.ListAllStorageVolumes(0)
		if err != nil {
			return errors.Wrapf(err, "list volumes in %q", tpool)
		}

		for _, vol := range vols {
			defer vol.Free()
			vName, err := vol.GetName()
			if err != nil {
				return errors.Wrapf(err, "get volume names in %q", tpool)
			}
			if !filter(vName) {
				continue
			}
			if err := vol.Delete(0); err != nil {
				return errors.Wrapf(err, "delete volume %q from %q", vName, tpool)
			}
			logger.WithField("volume", vName).Info("Deleted volume")
		}
	default:
		// blow away entire pool.
		if err := pool.Destroy(); err != nil {
			return errors.Wrapf(err, "destroy pool %q", tpool)
		}

		if err := pool.Undefine(); err != nil {
			return errors.Wrapf(err, "undefine pool %q", tpool)
		}
		logger.WithField("pool", tpool).Info("Deleted pool")
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

// New returns libvirt Uninstaller from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (destroy.Destroyer, error) {
	return &ClusterUninstaller{
		LibvirtURI: metadata.ClusterPlatformMetadata.Libvirt.URI,
		Filter:     AlwaysTrueFilter(), //TODO: change to ClusterNamePrefixFilter when all resources are prefixed.
		Logger:     logger,
	}, nil
}
