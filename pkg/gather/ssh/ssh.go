// Package ssh contains utilities that help gather logs, etc. on failures using ssh.
package ssh

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/net/proxy"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/installer/pkg/lineprinter"
)

type httpProxy struct {
	host     string
	haveAuth bool
	username string
	password string
	forward  proxy.Dialer
}

type httpsDialer struct{}

// Dialer for HTTPS.
var HTTPSDialer = httpsDialer{}

// TLS configuration needed for HTTPS dialer.
var TLSConfig = &tls.Config{}

// NewClient creates a new SSH client which can be used to SSH to address using user and the keys.
// if keys list is empty, it tries to load the keys from the user's environment.
func NewClient(user, address string, keys []string, proxyURL string) (*ssh.Client, error) {
	ag, agentType, err := getAgent(keys)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize the SSH agent")
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Use a callback rather than PublicKeys
			// so we only consult the agent once the remote server
			// wants it.
			ssh.PublicKeysCallback(ag.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var client *ssh.Client
	if proxyURL != "" {
		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse proxy URL")
		}
		proxy.RegisterDialerType("http", newHTTPProxy)
		proxy.RegisterDialerType("https", newHTTPSProxy)
		httpDialer, err := proxy.FromURL(parsedURL, proxy.Direct)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create HTTP dialer")
		}
		httpConn, err := httpDialer.Dial("tcp", address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to dial HTTP proxy")
		}
		ncc, chans, reqs, err := ssh.NewClientConn(httpConn, address, sshConfig)
		if err != nil {
			return nil, handleAgentError(err, agentType)
		}
		client = ssh.NewClient(ncc, chans, reqs)
	} else {
		client, err = ssh.Dial("tcp", address, sshConfig)
		if err != nil {
			return nil, handleAgentError(err, agentType)
		}
	}

	if err := agent.ForwardToAgent(client, ag); err != nil {
		return nil, errors.Wrap(err, "failed to forward agent")
	}
	return client, nil
}

func handleAgentError(err error, agentType string) error {
	if strings.Contains(err.Error(), "ssh: handshake failed: ssh: unable to authenticate") {
		if agentType == "agent" {
			return errors.Wrap(err, "failed to use pre-existing agent, make sure the appropriate keys exist in the agent for authentication")
		}
		return errors.Wrap(err, "failed to use the provided keys for authentication")
	}
	return err
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
	paths, err := os.ReadDir(d)
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
		data, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to read %q", path))
			continue
		}
		key, err := ssh.ParseRawPrivateKey(data)
		if err != nil {
			logrus.Debugf("failed to parse SSH private key from %q", path)
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

// newHTTPProxy creates a new httpProxy populated with credentials if any is available.
func newHTTPProxy(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	proxyDialer := new(httpProxy)
	proxyDialer.host = uri.Host
	proxyDialer.forward = forward
	if uri.User != nil {
		proxyDialer.haveAuth = true
		proxyDialer.username = uri.User.Username()
		proxyDialer.password, _ = uri.User.Password()
	}

	return proxyDialer, nil
}

// newHTTPSProxy creates a new httpProxy populated with credentials using HTTPS dialer.
func newHTTPSProxy(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	return newHTTPProxy(uri, HTTPSDialer)
}

// Dial is a custom dialer for the ssh client that will use HTTP CONNECT method via TLS handshake.
func (d httpsDialer) Dial(network, addr string) (c net.Conn, err error) {
	c, err = tls.Dial(network, addr, TLSConfig)
	if err != nil {
		fmt.Println(err)
		c.Close()
		return nil, errors.Wrap(err, "failed to connect to proxy")
	}
	return c, nil
}

// Dial is a custom dialer for the ssh client that will use HTTP CONNECT method
// to establish a connection to the proxy. This is needed because the net/proxy module
// only supports SOCKS5 proxies.
func (s *httpProxy) Dial(network, addr string) (net.Conn, error) {
	c, err := s.forward.Dial(network, s.host)
	if c == nil {
		return nil, errors.Wrap(err, "failed to connect to proxy")
	}
	if err != nil {
		c.Close()
		return nil, err
	}

	reqURL, err := url.Parse("http://" + addr)
	if err != nil {
		c.Close()
		return nil, err
	}
	reqURL.Scheme = ""

	req, err := http.NewRequest("CONNECT", reqURL.String(), nil)
	if err != nil {
		c.Close()
		return nil, err
	}
	req.Close = false
	if s.haveAuth {
		req.SetBasicAuth(s.username, s.password)
	}

	err = req.Write(c)
	if err != nil {
		c.Close()
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		c.Close()
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		c.Close()
		errMsg := fmt.Sprintf("error when connecting to proxy, StatusCode [%d]", resp.StatusCode)
		err := errors.Wrap(err, errMsg)
		return nil, err
	}

	return c, nil
}
