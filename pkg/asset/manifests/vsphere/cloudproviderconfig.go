package vsphere

import (
	"bytes"
	"fmt"
	"net"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func printIfNotEmpty(buf *bytes.Buffer, k, v string) {
	if v != "" {
		fmt.Fprintf(buf, "%s = %q\n", k, v)
	}
}

// CloudProviderConfig generates the cloud provider config for the vSphere platform.
// We check to see if a port value is specified for vcenter("url:port") and send them as
//separate fields in cloud config.
// folderPath is the absolute path to the VM folder that will be used for installation.
// p is the vSphere platform struct.
func CloudProviderConfig(folderPath string, p *vspheretypes.Platform) (string, error) {
	buf := new(bytes.Buffer)

	host, port, err := net.SplitHostPort(p.VCenter)
	if err != nil {
		host = p.VCenter
	}

	fmt.Fprintln(buf, "[Global]")
	printIfNotEmpty(buf, "secret-name", "vsphere-creds")
	printIfNotEmpty(buf, "secret-namespace", "kube-system")
	printIfNotEmpty(buf, "insecure-flag", "1")
	printIfNotEmpty(buf, "port", port)
	fmt.Fprintln(buf, "")

	fmt.Fprintln(buf, "[Workspace]")

	printIfNotEmpty(buf, "server", host)
	printIfNotEmpty(buf, "datacenter", p.Datacenter)
	printIfNotEmpty(buf, "default-datastore", p.DefaultDatastore)
	printIfNotEmpty(buf, "folder", folderPath)
	fmt.Fprintln(buf, "")

	fmt.Fprintf(buf, "[VirtualCenter %q]\n", host)
	printIfNotEmpty(buf, "datacenters", p.Datacenter)

	return buf.String(), nil
}
