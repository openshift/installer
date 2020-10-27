// Package ssh contains utilities that help gather logs, etc. on failures using ssh.
package ssh

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/installer/pkg/lineprinter"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// NewClient creates a new SSH client which can be used to SSH to address using user and the keys.
//
// if keys list is empty, it tries to load the keys from the user's environment.
func NewClient(user, address string, keys []string) (*ssh.Client, error) {
	ag, agentType, err := getAgent(keys)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize the SSH agent")
	}

	client, err := ssh.Dial("tcp", address, &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Use a callback rather than PublicKeys
			// so we only consult the agent once the remote server
			// wants it.
			ssh.PublicKeysCallback(ag.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "ssh: handshake failed: ssh: unable to authenticate") {
			if agentType == "agent" {
				return nil, errors.Wrap(err, "failed to use pre-existing agent, make sure the appropriate keys exist in the agent for authentication")
			}
			return nil, errors.Wrap(err, "failed to use the provided keys for authentication")
		}
		return nil, err
	}
	if err := agent.ForwardToAgent(client, ag); err != nil {
		return nil, errors.Wrap(err, "failed to forward agent")
	}
	return client, nil
}

// Run uses an SSH client to execute commands.
func Run(client *ssh.Client, command string) error {
	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	if err := agent.RequestAgentForwarding(sess); err != nil {
		return errors.Wrap(err, "failed to setup request agent forwarding")
	}

	debugW := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Debug}).Print}
	defer debugW.Close()
	sess.Stdout = debugW
	sess.Stderr = debugW
	return sess.Run(command)
}

// PullFileTo downloads the file from remote server using SSH connection and writes to localPath.
func PullFileTo(client *ssh.Client, remotePath, localPath string) error {
	sc, err := sftp.NewClient(client)
	if err != nil {
		return errors.Wrap(err, "failed to initialize the sftp client")
	}
	defer sc.Close()

	// Open the source file
	rFile, err := sc.Open(remotePath)
	if err != nil {
		return errors.Wrap(err, "failed to open remote file")
	}
	defer rFile.Close()

	lFile, err := os.Create(localPath)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer lFile.Close()

	if _, err := rFile.WriteTo(lFile); err != nil {
		return err
	}
	return nil
}

// defaultPrivateSSHKeys returns a list of all the PRIVATE SSH keys from user's home directory.
// It does not return any intermediate errors if at least one private key was loaded.
func defaultPrivateSSHKeys() (map[string]interface{}, error) {
	d := filepath.Join(os.Getenv("HOME"), ".ssh")
	paths, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read directory %q", d)
	}

	var files []string
	for _, path := range paths {
		if path.IsDir() {
			continue
		}
		files = append(files, filepath.Join(d, path.Name()))
	}
	keys, err := LoadPrivateSSHKeys(files)
	if len(keys) > 0 {
		return keys, nil
	}
	return nil, err
}

// LoadPrivateSSHKeys try to optimistically load PRIVATE SSH keys from the all paths.
func LoadPrivateSSHKeys(paths []string) (map[string]interface{}, error) {
	var errs []error
	keys := make(map[string]interface{})
	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to read %q", path))
			continue
		}
		key, err := ssh.ParseRawPrivateKey(data)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to parse SSH private key from %q", path))
			continue
		}
		keys[path] = key
	}
	if err := utilerrors.NewAggregate(errs); err != nil {
		return keys, err
	}
	return keys, nil
}
