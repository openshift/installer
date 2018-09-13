package configgenerator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"

	ignconfig "github.com/coreos/ignition/config/v2_2"
	ignconfigtypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/openshift/installer/installer/pkg/config"
	"github.com/vincent-petithory/dataurl"
)

const (
	caPath = "generated/tls/root-ca.crt"
)

// GenerateIgnConfig generates Ignition configs for the workers and masters.
func (c *ConfigGenerator) GenerateIgnConfig(clusterDir string) error {
	var masters config.NodePool
	var workers config.NodePool
	for _, pool := range c.NodePools {
		switch pool.Name {
		case "master":
			masters = pool
		case "worker":
			workers = pool
		default:
			return fmt.Errorf("unrecognized role: %s", pool.Name)
		}
	}

	ca, err := ioutil.ReadFile(filepath.Join(clusterDir, caPath))
	if err != nil {
		return err
	}

	workerCfg, err := parseIgnFile(workers.IgnitionFile)
	if err != nil {
		return fmt.Errorf("failed to parse Ignition config for workers: %v", err)
	}

	// XXX(crawford): The SSH key should only be added to the bootstrap
	//                node. After that, MCO should be responsible for
	//                distributing SSH keys.
	c.embedUserBlock(workerCfg)
	c.appendCertificateAuthority(workerCfg, ca)
	c.embedAppendBlock(workerCfg, "worker", "")

	if err = ignCfgToFile(workerCfg, filepath.Join(clusterDir, config.IgnitionPathWorker)); err != nil {
		return err
	}

	masterCfg, err := parseIgnFile(masters.IgnitionFile)
	if err != nil {
		return fmt.Errorf("failed to parse Ignition config for masters: %v", err)
	}

	for i := 0; i < masters.Count; i++ {
		ignCfg := masterCfg

		// XXX(crawford): The SSH key should only be added to the bootstrap
		//                node. After that, MCO should be responsible for
		//                distributing SSH keys.
		c.embedUserBlock(ignCfg)
		c.appendCertificateAuthority(ignCfg, ca)
		c.embedAppendBlock(ignCfg, "master", fmt.Sprintf("etcd_index=%d", i))

		if err = ignCfgToFile(ignCfg, filepath.Join(clusterDir, fmt.Sprintf(config.IgnitionPathMaster, i))); err != nil {
			return err
		}
	}

	return nil
}

func parseIgnFile(filePath string) (*ignconfigtypes.Config, error) {
	if filePath == "" {
		return &ignconfigtypes.Config{
			Ignition: ignconfigtypes.Ignition{
				Version: ignconfigtypes.MaxVersion.String(),
			},
		}, nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	cfg, rpt, _ := ignconfig.Parse(data)
	if len(rpt.Entries) > 0 {
		return nil, fmt.Errorf("failed to parse ignition file %s: %s", filePath, rpt.String())
	}

	return &cfg, nil
}

func (c *ConfigGenerator) embedAppendBlock(ignCfg *ignconfigtypes.Config, role string, query string) {
	appendBlock := ignconfigtypes.ConfigReference{
		Source:       c.getTNCURL(role, query),
		Verification: ignconfigtypes.Verification{Hash: nil},
	}
	ignCfg.Ignition.Config.Append = append(ignCfg.Ignition.Config.Append, appendBlock)
}

func (c *ConfigGenerator) appendCertificateAuthority(ignCfg *ignconfigtypes.Config, ca []byte) {
	ignCfg.Ignition.Security.TLS.CertificateAuthorities = append(ignCfg.Ignition.Security.TLS.CertificateAuthorities, ignconfigtypes.CaReference{
		Source: dataurl.EncodeBytes(ca),
	})
}

func (c *ConfigGenerator) embedUserBlock(ignCfg *ignconfigtypes.Config) {
	userBlock := ignconfigtypes.PasswdUser{
		Name: "core",
		SSHAuthorizedKeys: []ignconfigtypes.SSHAuthorizedKey{
			ignconfigtypes.SSHAuthorizedKey(c.SSHKey),
		},
	}

	ignCfg.Passwd.Users = append(ignCfg.Passwd.Users, userBlock)
}

func (c *ConfigGenerator) getTNCURL(role string, query string) string {
	var u string

	// cloud platforms put this behind a load balancer which remaps ports;
	// libvirt doesn't do that - use the tnc port directly
	port := 80
	if c.Platform == config.PlatformLibvirt {
		port = 49500
	}

	if role == "master" || role == "worker" {
		u = func() *url.URL {
			return &url.URL{
				Scheme:   "https",
				Host:     fmt.Sprintf("%s-tnc.%s:%d", c.Name, c.BaseDomain, port),
				Path:     fmt.Sprintf("/config/%s", role),
				RawQuery: query,
			}
		}().String()
	}
	return u
}

func ignCfgToFile(ignCfg *ignconfigtypes.Config, filePath string) error {
	data, err := json.MarshalIndent(ignCfg, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0666)
}
