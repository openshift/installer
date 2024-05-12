package installationiso

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/openshift-kni/lifecycle-agent/api/ibiconfig"
	"github.com/openshift-kni/lifecycle-agent/lca-cli/ops"
	"github.com/openshift-kni/lifecycle-agent/utils"
)

type InstallationIso struct {
	log     *logrus.Logger
	ops     ops.Ops
	workDir string
}

type IgnitionData struct {
	AuthFilePath         string
	PullSecretPath       string
	IBIConfigurationPath string
	BackupSecret         string
	PullSecret           string
	SshPublicKey         string
	InstallSeedScript    string
	SeedImage            string
	IBIConfiguration     string
}

//go:embed data/*
var folder embed.FS

func NewInstallationIso(log *logrus.Logger, ops ops.Ops, workDir string) *InstallationIso {
	return &InstallationIso{
		log:     log,
		ops:     ops,
		workDir: workDir,
	}
}

const (
	ibiButaneTemplateFilePath = "data/ibi-butane.template"
	seedInstallScriptFilePath = "data/install-rhcos-and-restore-seed.sh"
	butaneFiles               = "butaneFiles"
	butaneConfigFile          = "config.bu"
	ibiIgnitionFileName       = "ibi-ignition.json"
	rhcosIsoFileName          = "rhcos-live.x86_64.iso"
	ibiIsoFileName            = "rhcos-ibi.iso"
	coreosInstallerImage      = "quay.io/coreos/coreos-installer:latest"
	ibiConfigFileName         = "ibi-configuration.json"
	authIgnitionFilePath      = "/var/tmp/backup-secret.json"
	psIgnitioFilePath         = "/var/tmp/pull-secret.json"
	ibiConfigIgnitionFilePath = "/var/tmp/" + ibiConfigFileName
)

func (r *InstallationIso) Create(ibiConfig *ibiconfig.IBIPrepareConfig) error {
	r.log.Info("Creating IBI installation ISO")
	err := r.validate()
	if err != nil {
		return err
	}
	err = r.createIgnitionFile(ibiConfig)
	if err != nil {
		return err
	}
	if err := r.downloadLiveIso(ibiConfig.RHCOSLiveISO); err != nil {
		return err
	}
	if err := r.embedIgnitionToIso(); err != nil {
		return err
	}
	r.log.Infof("installation ISO created at: %s", path.Join(r.workDir, ibiIsoFileName))

	return nil
}

func (r *InstallationIso) validate() error {
	_, err := os.Stat(r.workDir)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("work dir doesn't exists %w", err)
	}
	return nil
}

func (r *InstallationIso) createIgnitionFile(ibiConfig *ibiconfig.IBIPrepareConfig) error {
	r.log.Info("Generating Ignition Config")
	err := r.renderButaneConfig(ibiConfig)
	if err != nil {
		return err
	}
	return r.renderIgnitionFile()
}

func (r *InstallationIso) renderIgnitionFile() error {
	ibiIsoPath := path.Join(r.workDir, ibiIgnitionFileName)
	if _, err := os.Stat(ibiIsoPath); err == nil {
		r.log.Infof("ignition file exists (%s), deleting it", ibiIsoPath)
		if err = os.Remove(ibiIsoPath); err != nil {
			return fmt.Errorf("failed to delete existing ignition config: %w", err)
		}
	}

	command := "podman"
	args := []string{"run",
		"-v", fmt.Sprintf("%s:/data:rw,Z", r.workDir),
		"--rm",
		"quay.io/coreos/butane:release",
		"--pretty", "--strict",
		"-d", "/data",
		path.Join("/data", butaneConfigFile),
		"-o", path.Join("/data", ibiIgnitionFileName),
	}
	_, err := r.ops.RunInHostNamespace(command, args...)
	if err != nil {
		return fmt.Errorf("failed to render ignition from config: %w", err)
	}

	return nil
}

func (r *InstallationIso) embedIgnitionToIso() error {
	ibiIsoPath := path.Join(r.workDir, ibiIsoFileName)
	if _, err := os.Stat(ibiIsoPath); err == nil {
		r.log.Infof("ibi ISO exists (%s), deleting it", ibiIsoPath)
		if err = os.Remove(ibiIsoPath); err != nil {
			return fmt.Errorf("failed to delete existing ibi ISO: %w", err)
		}
	}

	command := "podman"
	args := []string{"run",
		"-v", fmt.Sprintf("%s:/data:rw,Z", r.workDir),
		coreosInstallerImage,
		"iso", "ignition", "embed",
		"-i", path.Join("/data", ibiIgnitionFileName),
		"-o", path.Join("/data", ibiIsoFileName),
		path.Join("/data", rhcosIsoFileName),
	}

	if _, err := r.ops.RunInHostNamespace(command, args...); err != nil {
		return fmt.Errorf("failed to embed ign with args %s: %w", args, err)
	}
	return nil
}

func (r *InstallationIso) renderButaneConfig(ibiConfig *ibiconfig.IBIPrepareConfig) error {
	r.log.Debug("Generating butane config")
	var sshPublicKey []byte
	var err error
	if ibiConfig.SSHPublicKeyFile == "" {
		r.log.Info("ssh key not provided skipping")
	} else {
		sshPublicKey, err = os.ReadFile(ibiConfig.SSHPublicKeyFile)
		if err != nil {
			return fmt.Errorf("failed to read ssh public key: %w", err)
		}
	}

	butaneDataDir := path.Join(r.workDir, butaneFiles)
	r.log.Debugf("Create %s directory for storing butane config files", butaneDataDir)
	os.Mkdir(butaneDataDir, 0o700)
	// We could apply the template data using the files content (referenced in the butane config as inline)
	// but that might result unmarshal errors while translating the config
	// hence we are copying the files to the butaneDataDir to be referenced as local files
	seedInstallScriptInButane := path.Join(butaneFiles, "seedInstallScript")
	if err := r.copyFileToButaneDir(seedInstallScriptFilePath, path.Join(r.workDir, seedInstallScriptInButane)); err != nil {
		return err
	}
	pullSecretInButane := path.Join(butaneFiles, "pullSecret")
	if err := r.copyFileToButaneDir(ibiConfig.PullSecretFile, path.Join(r.workDir, pullSecretInButane)); err != nil {
		return err
	}
	backupSecretInButane := path.Join(butaneFiles, "backupSecret")
	if err := r.copyFileToButaneDir(ibiConfig.AuthFile, path.Join(r.workDir, backupSecretInButane)); err != nil {
		return err
	}
	// inside ignition auth and pull secret files paths differs from the ones provided in ibi cli
	ibiConfig.AuthFile = authIgnitionFilePath
	ibiConfig.PullSecretFile = psIgnitioFilePath
	ibiConfigurationInButane := path.Join(butaneFiles, ibiConfigFileName)
	if err := utils.MarshalToFile(ibiConfig, path.Join(r.workDir, ibiConfigurationInButane)); err != nil {
		return fmt.Errorf("failed to marshal ibi config to file: %w", err)
	}

	templateData := IgnitionData{
		AuthFilePath:         authIgnitionFilePath,
		PullSecretPath:       psIgnitioFilePath,
		IBIConfigurationPath: ibiConfigIgnitionFilePath,
		BackupSecret:         backupSecretInButane,
		PullSecret:           pullSecretInButane,
		SshPublicKey:         string(sshPublicKey),
		InstallSeedScript:    seedInstallScriptInButane,
		IBIConfiguration:     ibiConfigurationInButane,
		SeedImage:            ibiConfig.SeedImage,
	}

	template, err := folder.ReadFile(ibiButaneTemplateFilePath)
	if err != nil {
		return fmt.Errorf("error occurred while trying to read %s: %w", ibiButaneTemplateFilePath, err)
	}

	if err := utils.RenderTemplateFile(string(template), templateData, path.Join(r.workDir, butaneConfigFile), 0o644); err != nil {
		return fmt.Errorf("failed to render %s: %w", butaneConfigFile, err)
	}
	return nil
}

func (r *InstallationIso) copyFileToButaneDir(sourceFile, target string) error {
	var source fs.File
	var err error
	// this file isn't provided by the user, it's part of the data folder embedded into the go binary at the top of this file
	if sourceFile == seedInstallScriptFilePath {
		source, err = folder.Open(sourceFile)
	} else {
		source, err = os.Open(sourceFile)
	}

	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer source.Close()
	fileForButaneConfig, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("failed to create file under workdir: %w", err)
	}
	defer fileForButaneConfig.Close()
	if _, err = io.Copy(fileForButaneConfig, source); err != nil {
		return fmt.Errorf("failed to copy file to workdir: %w", err)
	}
	return nil
}

func (r *InstallationIso) downloadLiveIso(url string) error {
	r.log.Info("Downloading live ISO")
	rhcosIsoPath := path.Join(r.workDir, rhcosIsoFileName)
	if _, err := os.Stat(rhcosIsoPath); err == nil {
		r.log.Infof("rhcos live ISO (%s) exists, skipping download", rhcosIsoPath)
		return nil
	}

	isoFile, err := os.Create(rhcosIsoPath)
	if err != nil {
		return fmt.Errorf("failed to rhcos iso path in %s: %w", rhcosIsoPath, err)
	}
	defer isoFile.Close()

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to make http get call to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download ISO from URL, status: %s", resp.Status)
	}

	if _, err := io.Copy(isoFile, resp.Body); err != nil {
		return fmt.Errorf("failed iso file from resp: %w", err)
	}

	return nil
}
