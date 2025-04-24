package image

import (
	"context"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/password"
)

// AgentPassword is an asset that generates the kubadmin-password and hash and stores it to the tls directory.
type AgentPassword struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*AgentPassword)(nil)

// Dependencies returns the assets on which the AgentPassword asset depends.
func (a *AgentPassword) Dependencies() []asset.Asset {
	return []asset.Asset{
		&password.KubeadminPassword{},
	}
}

// Generate generates the password file (and hash) and stores it in the tls directory.
func (a *AgentPassword) Generate(ctx context.Context, dependencies asset.Parents) error {
	pwd := &password.KubeadminPassword{}
	dependencies.Get(pwd)

	agentKubeadminPasswordPath := filepath.Join("tls", "kubeadmin-password")
	agentKubeadminPasswordHashPath := filepath.Join("tls", "kubeadmin-password.hash")

	a.FileList = []*asset.File{
		{
			Filename: agentKubeadminPasswordPath,
			Data:     []byte(pwd.Password),
		},
		{
			Filename: agentKubeadminPasswordHashPath,
			Data:     pwd.PasswordHash,
		},
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *AgentPassword) Name() string {
	return "Agent Installer ISO"
}

// Load returns the password from disk.
func (a *AgentPassword) Load(f asset.FileFetcher) (bool, error) {
	// This is implemented because it is required by WritableAsset
	return false, nil
}

// Files returns the asset's files.
func (a *AgentPassword) Files() []*asset.File {
	return a.FileList
}
