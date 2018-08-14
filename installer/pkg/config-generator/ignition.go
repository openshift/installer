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

var (
	ignVersion   = "2.2.0"
	ignFilesPath = map[string]string{
		"master": config.IgnitionMaster,
		"worker": config.IgnitionWorker,
		"etcd":   config.IgnitionEtcd,
	}
	caPath = "generated/tls/root-ca.crt"
)

func (c *ConfigGenerator) poolToRoleMap() map[string]string {
	poolToRole := make(map[string]string)
	// assume no roles can share pools
	for _, n := range c.Master.NodePools {
		poolToRole[n] = "master"
	}
	for _, n := range c.Worker.NodePools {
		poolToRole[n] = "worker"
	}
	for _, n := range c.Etcd.NodePools {
		poolToRole[n] = "etcd"
	}
	return poolToRole
}

// GenerateIgnConfig generates, if successful, files with the ign config for each role.
func (c *ConfigGenerator) GenerateIgnConfig(clusterDir string) error {
	poolToRole := c.poolToRoleMap()
	for _, p := range c.NodePools {
		ignFile := p.IgnitionFile
		ignCfg, err := parseIgnFile(ignFile)
		if err != nil {
			return fmt.Errorf("failed to GenerateIgnConfig for pool %s and file %s: %v", p.Name, p.IgnitionFile, err)
		}
		role := poolToRole[p.Name]
		// TODO(alberto): Append block need to be different for each etcd node.
		// add loop over count if role is etcd
		c.embedAppendBlock(ignCfg, role)

		ca := filepath.Join(clusterDir, caPath)
		if err = c.appendCertificateAuthority(ignCfg, ca); err != nil {
			return err
		}

		// XXX(crawford): The SSH key should only be added to the bootstrap
		//                node. After that, MCO should be responsible for
		//                distributing SSH keys.
		c.embedUserBlock(ignCfg)

		fileTargetPath := filepath.Join(clusterDir, ignFilesPath[role])
		if err = ignCfgToFile(*ignCfg, fileTargetPath); err != nil {
			return err
		}
	}
	return nil
}

func parseIgnFile(filePath string) (*ignconfigtypes.Config, error) {
	if filePath == "" {
		ignition := &ignconfigtypes.Ignition{
			Version: ignVersion,
		}
		return &ignconfigtypes.Config{Ignition: *ignition}, nil
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

func (c *ConfigGenerator) embedAppendBlock(ignCfg *ignconfigtypes.Config, role string) {
	appendBlock := ignconfigtypes.ConfigReference{
		Source:       c.getTNCURL(role),
		Verification: ignconfigtypes.Verification{Hash: nil},
	}
	ignCfg.Ignition.Config.Append = append(ignCfg.Ignition.Config.Append, appendBlock)
}

func (c *ConfigGenerator) appendCertificateAuthority(ignCfg *ignconfigtypes.Config, caPath string) error {
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		return err
	}

	ignCfg.Ignition.Security.TLS.CertificateAuthorities = append(ignCfg.Ignition.Security.TLS.CertificateAuthorities, ignconfigtypes.CaReference{
		Source: dataurl.EncodeBytes(ca),
	})

	return nil
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

func (c *ConfigGenerator) getTNCURL(role string) string {
	var u string

	// cloud platforms put this behind a load balancer which remaps ports;
	// libvirt doesn't do that - use the tnc port directly
	port := 80
	if c.Platform == config.PlatformLibvirt {
		port = 49500
	}

	// XXX: The bootstrap node on AWS uses a CNAME to redirect TNC-bound
	// traffic to S3. Because of this, HTTPS cannot be used.
	scheme := "https"
	if c.Platform == config.PlatformAWS && role == "master" {
		scheme = "http"
	}

	if role == "master" || role == "worker" {
		u = func() *url.URL {
			return &url.URL{
				Scheme: scheme,
				Host:   fmt.Sprintf("%s-tnc.%s:%d", c.Name, c.BaseDomain, port),
				Path:   fmt.Sprintf("/config/%s", role),
			}
		}().String()
	}
	return u
}

func ignCfgToFile(ignCfg ignconfigtypes.Config, filePath string) error {
	data, err := json.MarshalIndent(&ignCfg, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0666)
}
