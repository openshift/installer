package terraform

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster/tfvars"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/lineprinter"
	"github.com/openshift/installer/pkg/metrics/timer"
	"github.com/openshift/installer/pkg/types"
)

const (
	tfVarsFileName         = "terraform.tfvars.json"
	tfPlatformVarsFileName = "terraform.platform.auto.tfvars.json"
)

// Provider implements the infrastructure.Provider interface.
type Provider struct {
	stages []Stage
}

// InitializeProvider creates a concrete infrastructure.Provider for the given platform.
func InitializeProvider(stages []Stage) infrastructure.Provider {
	return &Provider{stages}
}

// Provision implements pkg/infrastructure/provider.Provision. Provision iterates
// through each of the stages and applies the Terraform config for the stage.
func (p *Provider) Provision(_ context.Context, dir string, parents asset.Parents) ([]*asset.File, error) {
	tfVars := &tfvars.TerraformVariables{}
	parents.Get(tfVars)
	vars := tfVars.Files()

	fileList := []*asset.File{}
	terraformDir := filepath.Join(dir, "terraform")
	if err := os.Mkdir(terraformDir, 0777); err != nil {
		return nil, fmt.Errorf("could not create the terraform directory: %w", err)
	}

	terraformDirPath, err := filepath.Abs(terraformDir)
	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path of terraform directory: %w", err)
	}

	defer os.RemoveAll(terraformDir)
	if err = UnpackTerraform(terraformDirPath, p.stages); err != nil {
		return nil, fmt.Errorf("error unpacking terraform: %w", err)
	}

	for _, stage := range p.stages {
		outputs, stateFile, err := applyStage(stage.Platform(), stage, terraformDirPath, vars)
		if err != nil {
			// Write the state file to the install directory even if the apply failed.
			if stateFile != nil {
				fileList = append(fileList, stateFile)
			}
			return fileList, fmt.Errorf("failure applying terraform for %q stage: %w", stage.Name(), err)
		}
		vars = append(vars, outputs)
		fileList = append(fileList, outputs)
		fileList = append(fileList, stateFile)

		_, extErr := stage.ExtractLBConfig(dir, terraformDirPath, outputs, vars[0])
		if extErr != nil {
			return fileList, fmt.Errorf("failed to extract load balancer information: %w", extErr)
		}
	}
	return fileList, nil
}

// DestroyBootstrap implements pkg/infrastructure/provider.DestroyBootstrap.
// DestroyBootstrap iterates through each stage, and will run the destroy
// command when defined on a stage.
func (p *Provider) DestroyBootstrap(ctx context.Context, dir string) error {
	varFiles := []string{tfVarsFileName, tfPlatformVarsFileName}
	for _, stage := range p.stages {
		varFiles = append(varFiles, stage.OutputsFilename())
	}

	terraformDir := filepath.Join(dir, "terraform")
	if err := os.Mkdir(terraformDir, 0777); err != nil {
		return fmt.Errorf("could not create the terraform directory: %w", err)
	}

	terraformDirPath, err := filepath.Abs(terraformDir)
	if err != nil {
		return fmt.Errorf("could not get absolute path of terraform directory: %w", err)
	}

	defer os.RemoveAll(terraformDirPath)
	if err = UnpackTerraform(terraformDirPath, p.stages); err != nil {
		return fmt.Errorf("error unpacking terraform: %w", err)
	}

	for i := len(p.stages) - 1; i >= 0; i-- {
		stage := p.stages[i]

		if !stage.DestroyWithBootstrap() {
			continue
		}

		tempDir, err := os.MkdirTemp("", fmt.Sprintf("openshift-install-%s-", stage.Name()))
		if err != nil {
			return fmt.Errorf("failed to create temporary directory for Terraform execution: %w", err)
		}
		defer os.RemoveAll(tempDir)

		stateFilePathInInstallDir := filepath.Join(dir, stage.StateFilename())
		stateFilePathInTempDir := filepath.Join(tempDir, StateFilename)
		if err := copyFile(stateFilePathInInstallDir, stateFilePathInTempDir); err != nil {
			return fmt.Errorf("failed to copy state file to the temporary directory: %w", err)
		}

		targetVarFiles := make([]string, 0, len(varFiles))
		for _, filename := range varFiles {
			sourcePath := filepath.Join(dir, filename)
			targetPath := filepath.Join(tempDir, filename)
			if err := copyFile(sourcePath, targetPath); err != nil {
				// platform may not need platform-specific Terraform variables
				if filename == tfPlatformVarsFileName {
					var pErr *os.PathError
					if errors.As(err, &pErr) && pErr.Path == sourcePath {
						continue
					}
				}
				return fmt.Errorf("failed to copy %s to the temporary directory: %w", filename, err)
			}
			targetVarFiles = append(targetVarFiles, targetPath)
		}

		if err := stage.Destroy(tempDir, terraformDirPath, targetVarFiles); err != nil {
			return err
		}

		if err := copyFile(stateFilePathInTempDir, stateFilePathInInstallDir); err != nil {
			return fmt.Errorf("failed to copy state file from the temporary directory: %w", err)
		}
	}
	return nil
}

// ExtractHostAddresses implements pkg/infrastructure/provider.ExtractHostAddresses. Extracts the addresses to be used
// for gathering debug logs by inspecting the Terraform output files.
func (p *Provider) ExtractHostAddresses(dir string, config *types.InstallConfig, ha *infrastructure.HostAddresses) error {
	for _, stage := range p.stages {
		stageBootstrap, stagePort, stageMasters, err := stage.ExtractHostAddresses(dir, config)
		if err != nil {
			logrus.Warnf("Failed to extract host addresses: %s", err.Error())
		} else {
			if stageBootstrap != "" {
				ha.Bootstrap = stageBootstrap
			}
			if stagePort != 0 {
				ha.Port = stagePort
			}
			if len(stageMasters) > 0 {
				ha.Masters = stageMasters
			}
		}
	}
	return nil
}

// newTFExec creates a tfexec.Terraform for executing Terraform CLI commands.
// The `datadir` is the location to which the terraform plan (tf files, etc) has been unpacked.
// The `terraformDir` is the location to which Terraform, provider binaries, & .terraform data dir have been unpacked.
// The stdout and stderr will be sent to the logger at the debug and error levels,
// respectively.
func newTFExec(datadir string, terraformDir string) (*tfexec.Terraform, error) {
	tfPath := filepath.Join(terraformDir, "bin", "terraform")
	tf, err := tfexec.NewTerraform(datadir, tfPath)
	if err != nil {
		return nil, err
	}

	// terraform-exec will not accept debug logs unless a log file path has
	// been specified. And it makes sense since the logging is very verbose.
	if path, ok := os.LookupEnv("TF_LOG_PATH"); ok {
		// These might fail if tf cli does not have a compatible version. Since
		// the exact same check is repeated, we just have to verify error once
		// for all calls
		if err := tf.SetLog(os.Getenv("TF_LOG")); err != nil {
			// We want to skip setting the log path since tf-exec lib will
			// default to TRACE log levels which can risk leaking sensitive
			// data
			logrus.Infof("Skipping setting terraform log levels: %v", err)
		} else {
			tf.SetLogCore(os.Getenv("TF_LOG_CORE"))         //nolint:errcheck
			tf.SetLogProvider(os.Getenv("TF_LOG_PROVIDER")) //nolint:errcheck
			// This never returns any errors despite its signature
			tf.SetLogPath(path) //nolint:errcheck
		}
	}

	// Add terraform info logs to the installer log
	lpDebug := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Debug}).Print}
	lpError := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Error}).Print}
	defer lpDebug.Close()
	defer lpError.Close()

	tf.SetStdout(lpDebug)
	tf.SetStderr(lpError)
	tf.SetLogger(newPrintfer())

	// Set the Terraform data dir to be the same as the terraformDir so that
	// files we unpack are contained and, more importantly, we can ensure the
	// provider binaries unpacked in the Terraform data dir have the same permission
	// levels as the Terraform binary.
	dd := path.Join(terraformDir, ".terraform")
	os.Setenv("TF_DATA_DIR", dd)

	return tf, nil
}

// Apply unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// apply'.
func Apply(dir string, platform string, stage Stage, terraformDir string, extraOpts ...tfexec.ApplyOption) error {
	if err := unpackAndInit(dir, platform, stage.Name(), terraformDir, stage.Providers()); err != nil {
		return err
	}

	tf, err := newTFExec(dir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}
	err = tf.Apply(context.Background(), extraOpts...)
	return errors.Wrap(diagnoseApplyError(err), "failed to apply Terraform")
}

// Destroy unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// destroy'.
func Destroy(dir string, platform string, stage Stage, terraformDir string, extraOpts ...tfexec.DestroyOption) error {
	if err := unpackAndInit(dir, platform, stage.Name(), terraformDir, stage.Providers()); err != nil {
		return err
	}

	tf, err := newTFExec(dir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}
	return errors.Wrap(
		tf.Destroy(context.Background(), extraOpts...),
		"failed doing terraform destroy",
	)
}

func applyStage(platform string, stage Stage, terraformDir string, tfvarsFiles []*asset.File) (*asset.File, *asset.File, error) {
	// Copy the terraform.tfvars to a temp directory which will contain the terraform plan.
	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("openshift-install-%s-", stage.Name()))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create temp dir for terraform execution")
	}
	defer os.RemoveAll(tmpDir)

	extraOpts := []tfexec.ApplyOption{}
	for _, file := range tfvarsFiles {
		if err := os.WriteFile(filepath.Join(tmpDir, file.Filename), file.Data, 0o600); err != nil {
			return nil, nil, err
		}
		extraOpts = append(extraOpts, tfexec.VarFile(filepath.Join(tmpDir, file.Filename)))
	}

	return applyTerraform(tmpDir, platform, stage, terraformDir, extraOpts...)
}

func applyTerraform(tmpDir string, platform string, stage Stage, terraformDir string, opts ...tfexec.ApplyOption) (outputsFile, stateFile *asset.File, err error) {
	timer.StartTimer(stage.Name())
	defer timer.StopTimer(stage.Name())

	applyErr := Apply(tmpDir, platform, stage, terraformDir, opts...)

	if data, err := os.ReadFile(filepath.Join(tmpDir, StateFilename)); err == nil {
		stateFile = &asset.File{
			Filename: stage.StateFilename(),
			Data:     data,
		}
	} else if !os.IsNotExist(err) {
		logrus.Errorf("Failed to read tfstate: %v", err)
		return nil, nil, errors.Wrap(err, "failed to read tfstate")
	}

	if applyErr != nil {
		return nil, stateFile, fmt.Errorf("error applying Terraform configs: %w", applyErr)
	}

	outputs, err := Outputs(tmpDir, terraformDir)
	if err != nil {
		return nil, stateFile, errors.Wrapf(err, "could not get outputs from stage %q", stage.Name())
	}

	outputsFile = &asset.File{
		Filename: stage.OutputsFilename(),
		Data:     outputs,
	}
	return outputsFile, stateFile, nil
}

func copyFile(from string, to string) error {
	data, err := os.ReadFile(from)
	if err != nil {
		return err
	}

	return os.WriteFile(to, data, 0o666) //nolint:gosec // state file doesn't need to be 0600
}
