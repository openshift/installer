package imagebuilder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	template "text/template"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"

	data "github.com/openshift-agent-team/fleeting/data/data/agent"
)

// ConfigBuilder builds an Ignition config
type ConfigBuilder struct {
	pullSecret      string
	serviceBaseURL  string
	infraEnvID      string
	pullSecretToken string
	nodeZeroIP      string
}

func New(nodeZeroIP string) *ConfigBuilder {
	pullSecret := getEnv("PULL_SECRET", "")
	// TODO: try setting SERVICE_BASE_URL within agent.service
	serviceBaseURL := getEnv("SERVICE_BASE_URL", "http://"+nodeZeroIP+":8090")
	// TODO: get id either from InfraEnv CR that is included
	// with tool, or query the id from the REST_API
	// curl http://SERVICE_BASE_URL/api/assisted-install/v2/infra-envs
	infraEnvID := getEnv("INFRA_ENV_ID", "infra-env-id-missing")
	// TODO: needs appropriate value if AUTH_TYPE != none
	pullSecretToken := getEnv("PULL_SECRET_TOKEN", "")

	return &ConfigBuilder{
		pullSecret:      pullSecret,
		serviceBaseURL:  serviceBaseURL,
		infraEnvID:      infraEnvID,
		pullSecretToken: pullSecretToken,
		nodeZeroIP:      nodeZeroIP,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
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
		mode := 0420
		pullSecret := igntypes.File{
			Node: igntypes.Node{
				Path:      "/root/.docker/config.json",
				Overwrite: ignutil.BoolToPtr(true),
			},
			FileEmbedded1: igntypes.FileEmbedded1{
				Mode: &mode,
				Contents: igntypes.Resource{
					Source: ignutil.StrToPtr(dataurl.EncodeBytes([]byte(c.pullSecret))),
				},
			},
		}
		files = append(files, pullSecret)
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
				file := igntypes.File{
					Node: igntypes.Node{
						Path:      fullPath,
						Overwrite: ignutil.BoolToPtr(true),
					},
					FileEmbedded1: igntypes.FileEmbedded1{
						Mode: &mode,
						Contents: igntypes.Resource{
							Source: ignutil.StrToPtr(dataurl.EncodeBytes([]byte(templated))),
						},
					},
				}
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

	entries, err := data.IgnitionData.ReadDir(basePath)
	if err != nil {
		return units, fmt.Errorf("failed to read systemd units: %w", err)
	}

	for _, e := range entries {
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
		"ServiceBaseURL":  c.serviceBaseURL,
		"infraEnvId":      c.infraEnvID,
		"PullSecretToken": c.pullSecretToken,
		"NodeZeroIP":      c.nodeZeroIP,
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
