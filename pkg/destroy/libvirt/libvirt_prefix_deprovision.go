package libvirt

import (
	"os"
	"strings"
	"time"

	libvirt "github.com/libvirt/libvirt-go"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
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
		if name == "default" {
			return false
		}
		return true
	}
}

// deleteFunc type is the interface a function needs to implement to be called as a goroutine.
// The (bool, error) return type mimics wait.ExponentialBackoff where the bool indicates successful
// completion, and the error is for unrecoverable errors.
type deleteFunc func(conn *libvirt.Connect, filter filterFunc, logger log.FieldLogger) (bool, error)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	LibvirtURI string
	Filter     filterFunc
	Logger     log.FieldLogger
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() error {
	deleteFuncs := map[string]deleteFunc{}
	populateDeleteFuncs(deleteFuncs)
	returnChannel := make(chan string)

	conn, err := libvirt.NewConnect(o.LibvirtURI)
	if err != nil {
		return err
	}

	// launch goroutines
	for name, function := range deleteFuncs {
		go deleteRunner(name, function, conn, o.Filter, o.Logger, returnChannel)
	}

	// wait for them to finish
	for i := 0; i < len(deleteFuncs); i++ {
		select {
		case res := <-returnChannel:
			o.Logger.Debugf("goroutine %v complete", res)
		}
	}

	return nil
}

func deleteRunner(deleteFuncName string, dFunction deleteFunc, conn *libvirt.Connect, filter filterFunc, logger log.FieldLogger, channel chan string) {
	backoffSettings := wait.Backoff{
		Duration: time.Second * 10,
		Factor:   1.3,
		Steps:    100,
	}

	err := wait.ExponentialBackoff(backoffSettings, func() (bool, error) {
		return dFunction(conn, filter, logger)
	})

	if err != nil {
		logger.Fatalf("Unrecoverable error/timed out: %v", err)
		os.Exit(1)
	}

	// record that the goroutine has run to completion
	channel <- deleteFuncName
	return
}

// populateDeleteFuncs is the list of functions that will be launched as goroutines
func populateDeleteFuncs(funcs map[string]deleteFunc) {
	funcs["deleteDomains"] = deleteDomains
	funcs["deleteVolumes"] = deleteVolumes
	funcs["deleteNetwork"] = deleteNetwork
}

func deleteDomains(conn *libvirt.Connect, filter filterFunc, logger log.FieldLogger) (bool, error) {
	logger.Debug("Deleting libvirt domains")
	defer logger.Debugf("Exiting deleting libvirt domains")

	domains, err := conn.ListAllDomains(0)
	if err != nil {
		logger.Errorf("Error listing domains: %v", err)
		return false, nil
	}

	for _, domain := range domains {
		defer domain.Free()
		dName, err := domain.GetName()
		if err != nil {
			logger.Errorf("Error getting name for domain: %v", err)
			return false, nil
		}
		if !filter(dName) {
			continue
		}

		if err := domain.Destroy(); err != nil {
			logger.Errorf("Error destroying domain %s: %v", dName, err)
			return false, nil
		}

		if err := domain.Undefine(); err != nil {
			logger.Errorf("Error un-defining domain %s: %v", dName, err)
			return false, nil
		}
		logger.WithField("domain", dName).Info("Deleted domain")
	}
	return true, nil
}

func deleteVolumes(conn *libvirt.Connect, filter filterFunc, logger log.FieldLogger) (bool, error) {
	logger.Debug("Deleting libvirt volumes")
	defer logger.Debugf("Exiting deleting libvirt volumes")

	pools, err := conn.ListStoragePools()
	if err != nil {
		logger.Errorf("Error listing storage pools: %v", err)
		return false, nil
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
		logger.Errorf("Error getting storage pool %s: %v", tpool, err)
		return false, nil
	}
	defer pool.Free()

	switch tpool {
	case "default":
		// delete all vols that return true from filter.
		vols, err := pool.ListAllStorageVolumes(0)
		if err != nil {
			logger.Errorf("Error listing storage volumes in %s: %v", tpool, err)
			return false, nil
		}

		for _, vol := range vols {
			defer vol.Free()
			vName, err := vol.GetName()
			if err != nil {
				logger.Errorf("Error getting name for volume: %v", err)
				return false, nil
			}
			if !filter(vName) {
				continue
			}
			if err := vol.Delete(0); err != nil {
				logger.Errorf("Error deleting volume %s: %v", vName, err)
				return false, nil
			}
			logger.WithField("volume", vName).Info("Deleted volume")
		}
	default:
		// blow away entire pool.
		if err := pool.Destroy(); err != nil {
			logger.Errorf("Error destroying pool %s: %v", tpool, err)
			return false, nil
		}

		if err := pool.Undefine(); err != nil {
			logger.Errorf("Error undefining pool %s: %v", tpool, err)
			return false, nil
		}
		logger.WithField("pool", tpool).Info("Deleted pool")
	}

	return true, nil
}

func deleteNetwork(conn *libvirt.Connect, filter filterFunc, logger log.FieldLogger) (bool, error) {
	logger.Debug("Deleting libvirt network")
	defer logger.Debugf("Exiting deleting libvirt network")

	networks, err := conn.ListNetworks()
	if err != nil {
		logger.Errorf("Error listing networks: %v", err)
		return false, nil
	}

	for _, nName := range networks {
		if !filter(nName) {
			continue
		}
		network, err := conn.LookupNetworkByName(nName)
		if err != nil {
			logger.Errorf("Error getting network %s: %v", nName, err)
			return false, nil
		}
		defer network.Free()

		if err := network.Destroy(); err != nil {
			logger.Errorf("Error destroying network %s: %v", nName, err)
			return false, nil
		}

		if err := network.Undefine(); err != nil {
			logger.Errorf("Error undefining network %s: %v", nName, err)
			return false, nil
		}
		logger.WithField("network", nName).Info("Deleted network")
	}
	return true, nil
}
