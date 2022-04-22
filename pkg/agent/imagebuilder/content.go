package imagebuilder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	template "text/template"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"

	data "github.com/openshift-agent-team/fleeting/data/data/agent"
	"github.com/openshift-agent-team/fleeting/pkg/agent/manifests"

	"github.com/openshift/assisted-service/models"
	"github.com/sirupsen/logrus"
)

const NMCONNECTIONS_DIR = "/etc/assisted/network"

// ConfigBuilder builds an Ignition config
type ConfigBuilder struct {
	pullSecret               string
	serviceBaseURL           string
	pullSecretToken          string
	nodeZeroIP               string
	createClusterParamsJSON  string
	createInfraEnvParamsJSON string
	apiVip                   string
	controlPlaneAgents       int
	workerAgents             int
	staticNetworkConfig      []*models.HostStaticNetworkConfig
}

func New() *ConfigBuilder {
	pullSecret := manifests.GetPullSecret()

	n := manifests.NewNMConfig()
	nodeZeroIP := n.GetNodeZeroIP()

	// TODO: needs appropriate value if AUTH_TYPE != none
	pullSecretToken := getEnv("PULL_SECRET_TOKEN", "")

	serviceBaseURL := fmt.Sprintf("http://%s/",
		net.JoinHostPort(nodeZeroIP, "8090"))

	clusterParams := manifests.CreateClusterParams()
	clusterJSON, err := json.Marshal(clusterParams)
	if err != nil {
		logrus.Errorf("Error marshalling cluster params into json: %v", err)
	}

	infraEnvParams, err := manifests.CreateInfraEnvParams()
	if err != nil {
		logrus.Errorf("Error building infra env params: %v", err)
	}

	infraEnvJSON, err := json.Marshal(infraEnvParams)
	if err != nil {
		logrus.Errorf("Error marshal infra env params into json: %v", err)
	}

	aci := manifests.GetAgentClusterInstall()
	clusterInstall := &aci

	return &ConfigBuilder{
		pullSecret:               pullSecret,
		serviceBaseURL:           serviceBaseURL,
		pullSecretToken:          pullSecretToken,
		nodeZeroIP:               nodeZeroIP,
		createClusterParamsJSON:  string(clusterJSON),
		createInfraEnvParamsJSON: string(infraEnvJSON),
		apiVip:                   clusterInstall.Spec.APIVIP,
		controlPlaneAgents:       clusterInstall.Spec.ProvisionRequirements.ControlPlaneAgents,
		workerAgents:             clusterInstall.Spec.ProvisionRequirements.WorkerAgents,
		staticNetworkConfig:      infraEnvParams.StaticNetworkConfig,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ignitionFileEmbed(path string, mode int, overwrite bool, data []byte) igntypes.File {
	source := ignutil.StrToPtr(dataurl.EncodeBytes(data))

	return igntypes.File{
		Node: igntypes.Node{Path: path, Overwrite: &overwrite},
		FileEmbedded1: igntypes.FileEmbedded1{
			Contents: igntypes.Resource{Source: source},
			Mode:     &mode,
		},
	}
}

// Ignition builds an ignition file and returns the bytes
func (c ConfigBuilder) Ignition() ([]byte, error) {
	var err error

	config := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Passwd: igntypes.Passwd{
			Users: []igntypes.PasswdUser{
				{
					Name:              "core",
					SSHAuthorizedKeys: c.getSSHPubKey(),
				},
			},
		},
	}

	files, err := c.getFiles()
	if err != nil {
		return nil, err
	}

	// pull secret not included in data/data/agent/files because embed.FS
	// does not list directories with name starting with '.'
	if c.pullSecret != "" {
		pullSecret := ignitionFileEmbed("/root/.docker/config.json", 0420, true, []byte(c.pullSecret))
		files = append(files, pullSecret)
	}

	if len(c.staticNetworkConfig) > 0 {
		// Get the static network configuration from nmstate and generate NetworkManager ignition files
		filesList, err := manifests.GetNMIgnitionFiles(c.staticNetworkConfig)
		if err == nil {
			for i := range filesList {
				nmFilePath := path.Join(NMCONNECTIONS_DIR, filesList[i].FilePath)
				nmStateIgnFile := ignitionFileEmbed(nmFilePath, 0600, true, []byte(filesList[i].FileContents))
				files = append(files, nmStateIgnFile)
			}

			nmStateScriptFilePath := "/usr/local/bin/pre-network-manager-config.sh"
			// A local version of the assisted-service internal script is currently used
			nmStateScript := ignitionFileEmbed(nmStateScriptFilePath, 0755, true, []byte(manifests.PreNetworkConfigScript))
			files = append(files, nmStateScript)
		} else {
			// If manifest files are invalid, terminate to avoid networking problems at boot
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	config.Storage.Files = files

	config.Systemd.Units, err = c.getUnits()
	if err != nil {
		return nil, err
	}

	return json.Marshal(config)
}

func (c ConfigBuilder) getSSHPubKey() (keys []igntypes.SSHAuthorizedKey) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	pubkey, err := os.ReadFile(path.Join(home, ".ssh", "id_rsa.pub"))
	if err != nil {
		return
	}
	return append(keys, igntypes.SSHAuthorizedKey(pubkey))
}

func (c ConfigBuilder) getFiles() ([]igntypes.File, error) {
	var readDir func(dirPath string, files []igntypes.File) ([]igntypes.File, error)
	files := make([]igntypes.File, 0)

	readDir = func(dirPath string, files []igntypes.File) ([]igntypes.File, error) {
		entries, err := data.IgnitionData.ReadDir(path.Join("files", dirPath))
		if err != nil {
			return files, fmt.Errorf("failed to open file dir \"%s\": %w", dirPath, err)
		}
		for _, e := range entries {
			fullPath := path.Join(dirPath, e.Name())
			if e.IsDir() {
				files, err = readDir(fullPath, files)
				if err != nil {
					return files, err
				}
			} else {
				contents, err := data.IgnitionData.ReadFile(path.Join("files", fullPath))
				if err != nil {
					return files, fmt.Errorf("failed to read file %s: %w", fullPath, err)
				}
				templated, err := c.templateString(e.Name(), string(contents))
				if err != nil {
					return files, err
				}

				mode := 0600
				if _, dirName := path.Split(dirPath); dirName == "bin" || dirName == "dispatcher.d" {
					mode = 0555
				}
				file := ignitionFileEmbed(strings.TrimSuffix(fullPath, ".template"), mode, true, []byte(templated))
				files = append(files, file)
			}
		}
		return files, nil
	}

	return readDir("/", files)
}

func (c ConfigBuilder) getUnits() ([]igntypes.Unit, error) {
	units := make([]igntypes.Unit, 0)
	basePath := "systemd/units"
	staticNetworkService := "pre-network-manager-config.service"

	entries, err := data.IgnitionData.ReadDir(basePath)
	if err != nil {
		return units, fmt.Errorf("failed to read systemd units: %w", err)
	}

	for _, e := range entries {
		if len(c.staticNetworkConfig) == 0 && e.Name() == staticNetworkService {
			continue
		}

		contents, err := data.IgnitionData.ReadFile(path.Join(basePath, e.Name()))
		if err != nil {
			return units, fmt.Errorf("failed to read unit %s: %w", e.Name(), err)
		}

		templated, err := c.templateString(e.Name(), string(contents))
		if err != nil {
			return units, err
		}

		unit := igntypes.Unit{
			Name:     strings.TrimSuffix(e.Name(), ".template"),
			Enabled:  ignutil.BoolToPtr(true),
			Contents: ignutil.StrToPtr(string(templated)),
		}
		units = append(units, unit)
	}

	return units, nil
}

func (c ConfigBuilder) templateString(name string, text string) (string, error) {
	params := map[string]interface{}{
		"ServiceBaseURL":           c.serviceBaseURL,
		"PullSecretToken":          c.pullSecretToken,
		"NodeZeroIP":               c.nodeZeroIP,
		"ClusterCreateParamsJSON":  c.createClusterParamsJSON,
		"InfraEnvCreateParamsJSON": c.createInfraEnvParamsJSON,
		"APIVIP":                   c.apiVip,
		"ControlPlaneAgents":       c.controlPlaneAgents,
		"WorkerAgents":             c.workerAgents,
	}

	tmpl, err := template.New(name).Parse(string(text))
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}
