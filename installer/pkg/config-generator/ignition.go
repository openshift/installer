package configgenerator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"

	ignconfig "github.com/coreos/ignition/config/v2_2"
	ignconfigtypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/openshift/installer/pkg/types/config"
	"github.com/vincent-petithory/dataurl"
)

const (
	caPath = "generated/tls/root-ca.crt"
)

// GenerateIgnConfig generates Ignition configs for the workers and masters.
// It returns the content of the ignition files.
func (c *ConfigGenerator) GenerateIgnConfig(clusterDir string) (masterIgns []string, workerIgn string, err error) {
	var masters config.NodePool
	var workers config.NodePool
	for _, pool := range c.NodePools {
		switch pool.Name {
		case "master":
			masters = pool
		case "worker":
			workers = pool
		default:
			return nil, "", fmt.Errorf("unrecognized role: %s", pool.Name)
		}
	}

	ca, err := ioutil.ReadFile(filepath.Join(clusterDir, caPath))
	if err != nil {
		return nil, "", err
	}

	workerCfg, err := parseIgnFile(workers.IgnitionFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse Ignition config for workers: %v", err)
	}

	// XXX(crawford): The SSH key should only be added to the bootstrap
	//                node. After that, MCO should be responsible for
	//                distributing SSH keys.
	c.embedUserBlock(&workerCfg)
	c.appendCertificateAuthority(&workerCfg, ca)
	c.embedAppendBlock(&workerCfg, "worker", "")

	ign, err := json.Marshal(&workerCfg)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal worker ignition: %v", err)
	}
	workerIgn = string(ign)

	masterCfg, err := parseIgnFile(masters.IgnitionFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse Ignition config for masters: %v", err)
	}

	for i := 0; i < masters.Count; i++ {
		ignCfg := masterCfg

		// XXX(crawford): The SSH key should only be added to the bootstrap
		//                node. After that, MCO should be responsible for
		//                distributing SSH keys.
		c.embedUserBlock(&ignCfg)
		c.appendCertificateAuthority(&ignCfg, ca)
		c.embedAppendBlock(&ignCfg, "master", fmt.Sprintf("etcd_index=%d", i))

		masterIgn, err := json.Marshal(&ignCfg)
		if err != nil {
			return nil, "", fmt.Errorf("failed to marshal master ignition: %v", err)
		}
		masterIgns = append(masterIgns, string(masterIgn))
	}

	return masterIgns, workerIgn, nil
}

func parseIgnFile(filePath string) (ignconfigtypes.Config, error) {
	if filePath == "" {
		return ignconfigtypes.Config{
			Ignition: ignconfigtypes.Ignition{
				Version: ignconfigtypes.MaxVersion.String(),
			},
		}, nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ignconfigtypes.Config{}, err
	}

	cfg, rpt, _ := ignconfig.Parse(data)
	if len(rpt.Entries) > 0 {
		return ignconfigtypes.Config{}, fmt.Errorf("failed to parse ignition file %s: %s", filePath, rpt.String())
	}

	return cfg, nil
}

func (c *ConfigGenerator) embedAppendBlock(ignCfg *ignconfigtypes.Config, role string, query string) {
	appendBlock := ignconfigtypes.ConfigReference{
		Source:       c.getMCSURL(role, query),
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

func (c *ConfigGenerator) getMCSURL(role string, query string) string {
	var u string
	port := 49500

	if role == "master" || role == "worker" {
		u = func() *url.URL {
			return &url.URL{
				Scheme:   "https",
				Host:     fmt.Sprintf("%s-api.%s:%d", c.Name, c.BaseDomain, port),
				Path:     fmt.Sprintf("/config/%s", role),
				RawQuery: query,
			}
		}().String()
	}
	return u
}
