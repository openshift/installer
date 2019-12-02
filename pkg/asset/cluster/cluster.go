package cluster

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster/aws"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/terraform"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Cluster uses the terraform executable to launch a cluster
// with the given terraform tfvar and generated templates.
type Cluster struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*Cluster)(nil)

// Name returns the human-friendly name of the asset.
func (c *Cluster) Name() string {
	return "Cluster"
}

// Dependencies returns the direct dependency for launching
// the cluster.
func (c *Cluster) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		&installconfig.InstallConfig{},
		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
		&TerraformVariables{},
		&password.KubeadminPassword{},
	}
}

func uploadFile(sftpClient *sftp.Client, src, dest string) error {
	// create destination file
	dstFile, err := sftpClient.Create(dest)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// create source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	// copy source file to destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func sshAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

// Generate launches the cluster and generates the terraform state file on disk.
func (c *Cluster) Generate(parents asset.Parents) (err error) {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	terraformVariables := &TerraformVariables{}
	parents.Get(clusterID, installConfig, terraformVariables)

	if installConfig.Config.Platform.None != nil {
		return errors.New("cluster cannot be created with platform set to 'none'")
	}

	// Copy the terraform.tfvars to a temp directory where the terraform will be invoked within.
	tmpDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return errors.Wrap(err, "failed to create temp dir for terraform execution")
	}
	defer os.RemoveAll(tmpDir)

	extraArgs := []string{}
	for _, file := range terraformVariables.Files() {
		if err := ioutil.WriteFile(filepath.Join(tmpDir, file.Filename), file.Data, 0600); err != nil {
			return err
		}
		extraArgs = append(extraArgs, fmt.Sprintf("-var-file=%s", filepath.Join(tmpDir, file.Filename)))
	}

	logrus.Infof("Creating infrastructure resources...")
	if installConfig.Config.Platform.AWS != nil {
		if err := aws.PreTerraform(context.TODO(), clusterID.InfraID, installConfig); err != nil {
			return err
		}
	}

	stateFile, err := terraform.Apply(tmpDir, installConfig.Config.Platform.Name(), extraArgs...)
	if err != nil {
		err = errors.Wrap(err, "failed to create cluster")
		if stateFile == "" {
			return err
		}
		// Store the error from the apply, but continue with the
		// generation so that the Terraform state file is recovered from
		// the temporary directory.
	}

	if installConfig.Config.Platform.OpenStack != nil {
		// Get public IP of the bootstrap machine
		state, err := terraform.ReadState(stateFile)
		if err != nil {
			return err
		}

		resource, err := terraform.LookupResource(state, "module.bootstrap", "openstack_networking_floatingip_v2", "bootstrap_fip")
		if err != nil {
			return err
		}
		publicIP := resource.Instances[0].Attributes["address"].(string)

		// Create the bootstrapping script
		// TODO(mfedosin): move it to the ignition file
		f, err := os.Create("/tmp/bootstrap_files.py")
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(`import json
import os
with open("bootstrap_files.json", "r") as json_data:
	data = json.load(json_data)
	for path, content in data.items():
		os.makedirs(os.path.dirname(path), exist_ok=True)
		with open(path, 'w') as f:
			if path.endswith('.sh'):
				os.chmod(path, 0o755)
			f.write(content)`)
		if err != nil {
			return err
		}

		time.Sleep(30 * time.Second)
		client, err := ssh.Dial("tcp", publicIP+":22", &ssh.ClientConfig{
			User: "core",
			Auth: []ssh.AuthMethod{
				sshAgent(),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			return err
		}

		// create new SFTP client
		sftpClient, err := sftp.NewClient(client)
		if err != nil {
			return err
		}
		defer sftpClient.Close()

		// Upload bootstrapping data and the script to the remote machine
		err = uploadFile(sftpClient, "/tmp/bootstrap_files.json", "bootstrap_files.json")
		if err != nil {
			return err
		}

		err = uploadFile(sftpClient, "/tmp/bootstrap_files.py", "bootstrap_files.py")
		if err != nil {
			return err
		}

		sess, err := client.NewSession()
		if err != nil {
			return err
		}
		defer sess.Close()

		// Run the script that creates necessary files on the remote machine.
		err = sess.Run("sudo /usr/libexec/platform-python /var/home/core/bootstrap_files.py")
		if err != nil {
			return err
		}
	}

	data, err2 := ioutil.ReadFile(stateFile)
	if err2 == nil {
		c.FileList = append(c.FileList, &asset.File{
			Filename: terraform.StateFileName,
			Data:     data,
		})
	} else if err == nil {
		err = err2
	} else {
		logrus.Errorf("Failed to read tfstate: %v", err2)
	}

	return err
}

// Files returns the FileList generated by the asset.
func (c *Cluster) Files() []*asset.File {
	return c.FileList
}

// Load returns error if the tfstate file is already on-disk, because we want to
// prevent user from accidentally re-launching the cluster.
func (c *Cluster) Load(f asset.FileFetcher) (found bool, err error) {
	_, err = f.FetchByName(terraform.StateFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, errors.Errorf("%q already exists.  There may already be a running cluster", terraform.StateFileName)
}
