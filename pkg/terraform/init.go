package terraform

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/data"
	prov "github.com/openshift/installer/pkg/terraform/providers"
)

// unpack unpacks the platform-specific Terraform modules into the
// given directory.
func unpack(dir string, platform string, target string) (err error) {
	err = data.Unpack(dir, filepath.Join(platform, target))
	if err != nil {
		return err
	}

	err = data.Unpack(filepath.Join(dir, "config.tf"), "config.tf")
	if err != nil {
		return err
	}

	platformVarFile := fmt.Sprintf("variables-%s.tf", platform)

	err = data.Unpack(filepath.Join(dir, platformVarFile), filepath.Join(platform, platformVarFile))
	if err != nil {
		return err
	}

	err = data.Unpack(filepath.Join(dir, "terraform.rc"), "terraform.rc")
	if err != nil {
		return err
	}

	return nil
}

// unpackAndInit unpacks the platform-specific Terraform modules into
// the given directory and then runs 'terraform init'.
func unpackAndInit(dir string, platform string, target string, terraformDir string, providers []prov.Provider) (err error) {
	err = unpack(dir, platform, target)
	if err != nil {
		return errors.Wrap(err, "failed to unpack Terraform modules")
	}

	if err := addVersionsFiles(dir, providers); err != nil {
		return errors.Wrap(err, "failed to write versions.tf files")
	}

	tf, err := newTFExec(dir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}

	// Explicitly specify the CLI config file to use so that we control the providers that are used.
	os.Setenv("TF_CLI_CONFIG_FILE", filepath.Join(dir, "terraform.rc"))

	return errors.Wrap(
		tf.Init(context.Background(), tfexec.PluginDir(filepath.Join(terraformDir, "plugins"))),
		"failed doing terraform init",
	)
}

const versionFileTemplate = `terraform {
  required_version = ">= 1.0.0"
  required_providers {
{{- range .}}
    {{.Name}} = {
      source = "{{.Source}}"
    }
{{- end}}
  }
}
`

func addVersionsFiles(dir string, providers []prov.Provider) error {
	tmpl := template.Must(template.New("versions").Parse(versionFileTemplate))
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, providers); err != nil {
		return errors.Wrap(err, "could not create versions.tf from template")
	}
	return addFileToAllDirectories("versions.tf", buf.Bytes(), dir)
}

func addFileToAllDirectories(name string, data []byte, dir string) error {
	if err := os.WriteFile(filepath.Join(dir, name), data, 0666); err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if err := addFileToAllDirectories(name, data, filepath.Join(dir, entry.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

// UnpackTerraform unpacks the terraform binary and the specified provider binaries into the specified directory.
func UnpackTerraform(dir string, stages []Stage) error {
	// Unpack the terraform binary.
	if err := prov.UnpackTerraformBinary(filepath.Join(dir, "bin")); err != nil {
		return err
	}

	// Unpack the providers.
	providers := sets.NewString()
	for _, stage := range stages {
		for _, provider := range stage.Providers() {
			if providers.Has(provider.Name) {
				continue
			}
			if err := provider.Extract(filepath.Join(dir, "plugins")); err != nil {
				return err
			}
			providers.Insert(provider.Name)
		}
	}

	return nil
}
