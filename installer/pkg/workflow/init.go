package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"
	"text/template"
)

// InitWorkflow creates new instances of the 'init' workflow,
// responsible for initializing a new cluster.
func InitWorkflow(domain, name, licensePath, pullSecretPath string) Workflow {
	return Workflow{
		metadata: metadata{
			domain:         domain,
			name:           name,
			licensePath:    licensePath,
			pullSecretPath: pullSecretPath,
		},
		steps: []Step{
			initWorspaceStep,
			refreshConfigStep,
		},
	}
}

func initWorspaceStep(m *metadata) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	if m.name == "" {
		m.name = strings.Replace(m.domain, ".", "", -1)
	}

	// generate clusterDir folder
	clusterDir := filepath.Join(dir, m.name)
	m.clusterDir = clusterDir
	if stat, err := os.Stat(clusterDir); err == nil && stat.IsDir() {
		return fmt.Errorf("cluster directory already exists at %q", clusterDir)
	}

	if err := os.MkdirAll(clusterDir, os.ModeDir|0755); err != nil {
		return fmt.Errorf("failed to create cluster directory at %q", clusterDir)
	}

	// generate cluster config
	configFilePath := filepath.Join(clusterDir, configFileName)
	f, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to create cluster config at %q: %v", clusterDir, err)
	}

	ctd := configTemplateData{
		BaseDomain:     m.domain,
		LicensePath:    m.licensePath,
		Name:           m.name,
		PullSecretPath: m.pullSecretPath,
	}
	tmpl := template.Must(template.New("config").Parse(configTemplate))
	if err := tmpl.Execute(f, ctd); err != nil {
		return err
	}

	// generate the internal config file under the clusterDir folder
	return buildInternalConfig(clusterDir)
}
