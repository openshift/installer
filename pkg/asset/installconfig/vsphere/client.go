package vsphere

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

// Finder interface represents the client that is used to connect to VSphere to get specific
// information from the resources in the VCenter. This interface just describes all the useful
// functions used by the installer from the finder function in vmware govmomi package and is
// mostly used to create a mock client that can be used for testing.
type Finder interface {
	Datacenter(ctx context.Context, path string) (*object.Datacenter, error)
	DatacenterList(ctx context.Context, path string) ([]*object.Datacenter, error)
	DatastoreList(ctx context.Context, path string) ([]*object.Datastore, error)
	ClusterComputeResource(ctx context.Context, path string) (*object.ClusterComputeResource, error)
	ClusterComputeResourceList(ctx context.Context, path string) ([]*object.ClusterComputeResource, error)
	Folder(ctx context.Context, path string) (*object.Folder, error)
	NetworkList(ctx context.Context, path string) ([]object.NetworkReference, error)
	Network(ctx context.Context, path string) (object.NetworkReference, error)
	ResourcePool(ctx context.Context, path string) (*object.ResourcePool, error)
	VirtualMachine(ctx context.Context, path string) (*object.VirtualMachine, error)
	VirtualMachineList(ctx context.Context, path string) ([]*object.VirtualMachine, error)
	HostSystemList(ctx context.Context, path string) ([]*object.HostSystem, error)
	ObjectReference(ctx context.Context, ref types.ManagedObjectReference) (object.Reference, error)
}

// NewFinder creates a new client that conforms with the Finder interface and returns a
// vmware govmomi finder object that can be used to search for resources in vsphere.
func NewFinder(client *vim25.Client, all ...bool) Finder {
	return find.NewFinder(client, all...)
}

// ClientLogout is empty function that logs out of vSphere clients
type ClientLogout func()

// SOAPResponse represents the structure of SOAP responses
type SOAPResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		XMLName xml.Name `xml:"Body"`
		Fault   *struct {
			XMLName xml.Name `xml:"Fault"`
			Code    struct {
				XMLName xml.Name `xml:"faultcode"`
				Value   string   `xml:",chardata"`
			} `xml:"faultcode"`
			Reason struct {
				XMLName xml.Name `xml:"faultstring"`
				Value   string   `xml:",chardata"`
			} `xml:"faultstring"`
			Detail struct {
				XMLName xml.Name `xml:"detail"`
				Content string   `xml:",chardata"`
			} `xml:"detail"`
		} `xml:"Fault,omitempty"`
	} `xml:"Body"`
}

// CustomTransport wraps the default transport to intercept SOAP responses
type CustomTransport struct {
	http.RoundTripper
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Call the original transport
	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	resp.Body.Close()

	// Check if it's a SOAP response
	if strings.Contains(string(body), "<?xml") && strings.Contains(string(body), "Envelope") {
		logrus.Info("=== Intercepted SOAP Response ===")
		logrus.Infof("URL: %s", req.URL.String())
		logrus.Infof("Method: %s", req.Method)
		logrus.Infof("Response Body:\n%s", string(body))

		// Parse SOAP response for privilege errors
		var soapResp SOAPResponse
		if err := xml.Unmarshal(body, &soapResp); err == nil {
			if soapResp.Body.Fault != nil {
				logrus.Error("=== PRIVILEGE ERROR DETECTED ===")
				logrus.Errorf("Fault Code: %s", soapResp.Body.Fault.Code.Value)
				logrus.Errorf("Fault Reason: %s", soapResp.Body.Fault.Reason.Value)
				logrus.Errorf("Fault Detail: %s", soapResp.Body.Fault.Detail.Content)
				logrus.Error("================================")
			}
		}

		// Check for privilege-related error messages in the response
		bodyStr := string(body)
		privilegeKeywords := []string{
			"privilege", "permission", "access denied", "unauthorized", "forbidden",
			"NoPermission", "InvalidLogin", "InvalidPrivilege",
		}
		for _, keyword := range privilegeKeywords {
			if strings.Contains(strings.ToLower(bodyStr), strings.ToLower(keyword)) {
				logrus.Errorf("=== POTENTIAL PRIVILEGE ISSUE DETECTED (keyword: %s) ===", keyword)
				logrus.Error("Response contains privilege-related content")
				logrus.Error("==================================================")
				break
			}
		}

				// Check specifically for missingPrivileges and format the message
		if strings.Contains(bodyStr, "missingPrivileges") {
			logrus.Error("=== MISSING PRIVILEGES DETECTED ===")
			logrus.Error("The following SOAP response contains missingPrivileges information:")
			
			// Try to format the XML for better readability
			var v interface{}
			if err := xml.Unmarshal(body, &v); err == nil {
				// Marshal with indentation for pretty formatting
				prettyXML, err := xml.MarshalIndent(v, "", "  ")
				if err == nil {
					logrus.Errorf("Formatted Response:\n%s", string(prettyXML))
				} else {
					logrus.Errorf("Original Response:\n%s", bodyStr)
				}
			} else {
				// If XML parsing fails, try a simpler approach - just add line breaks between tags
				formatted := strings.ReplaceAll(bodyStr, "><", ">\n<")
				formatted = strings.ReplaceAll(formatted, "<?xml", "<?xml\n")
				logrus.Errorf("Formatted Response (simplified):\n%s", formatted)
			}
			logrus.Error("=== END MISSING PRIVILEGES ===")
		}

		logrus.Info("=== End SOAP Response ===")
	}

	// Create a new response with the body
	resp.Body = io.NopCloser(strings.NewReader(string(body)))
	return resp, nil
}

// createTransport creates a transport that respects the insecure flag
func createTransport(insecure bool) http.RoundTripper {
	if insecure {
		// Create a transport that skips TLS verification
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		return transport
	}
	// Use default transport for secure connections
	return http.DefaultTransport
}

// CreateVSphereClients creates the SOAP and REST client to access
// different portions of the vSphere API
// e.g. tags are only available in REST
func CreateVSphereClients(ctx context.Context, vcenter, username, password string) (*vim25.Client, *rest.Client, ClientLogout, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	u, err := soap.ParseURL(vcenter)
	if err != nil {
		return nil, nil, nil, err
	}
	u.User = url.UserPassword(username, password)

	// Create custom transport with SOAP response logging
	customTransport := &CustomTransport{
		RoundTripper: createTransport(false), // Always use secure connections in installer
	}

	// Create SOAP client with custom transport
	soapClient := soap.NewClient(u, false)
	soapClient.Transport = customTransport

	// Create vim25 client
	vimClient, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create govmomi client
	client := &govmomi.Client{
		Client:         vimClient,
		SessionManager: session.NewManager(vimClient),
	}

	// Login to vSphere
	err = client.Login(ctx, u.User)
	if err != nil {
		// Check if it's a credential-related error
		if strings.Contains(err.Error(), "incorrect user name or password") ||
			strings.Contains(err.Error(), "Cannot complete login") ||
			strings.Contains(err.Error(), "InvalidLogin") {
			return nil, nil, nil, errors.Errorf("vSphere authentication failed - please verify username and password: %w", err)
		}
		return nil, nil, nil, errors.Errorf("unable to login to vCenter: %w", err)
	}

	restClient := rest.NewClient(client.Client)
	err = restClient.Login(ctx, u.User)
	if err != nil {
		logoutErr := client.Logout(context.TODO())
		if logoutErr != nil {
			err = logoutErr
		}
		return nil, nil, nil, err
	}

	return client.Client, restClient, func() {
		client.Logout(context.TODO())
		restClient.Logout(context.TODO())
	}, nil
}

// getNetworks returns a slice of Managed Object references for networks in the given vSphere Cluster.
func getNetworks(ctx context.Context, ccr *object.ClusterComputeResource) ([]types.ManagedObjectReference, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	var ccrMo mo.ClusterComputeResource

	err := ccr.Properties(ctx, ccr.Reference(), []string{"network"}, &ccrMo)
	if err != nil {
		return nil, errors.Wrap(err, "could not get properties of cluster")
	}
	return ccrMo.Network, nil
}

// GetClusterNetworks returns a slice of Managed Object references for vSphere networks in the given Datacenter
// and Cluster.
func GetClusterNetworks(ctx context.Context, finder Finder, datacenter, cluster string) ([]types.ManagedObjectReference, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	ccr, err := finder.ClusterComputeResource(context.TODO(), cluster)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find vSphere cluster at %s", cluster)
	}

	// Get list of Networks inside vSphere Cluster
	networks, err := getNetworks(ctx, ccr)
	if err != nil {
		return nil, err
	}

	return networks, nil
}

// GetNetworkName returns the name of a vSphere network given its Managed Object reference.
func GetNetworkName(ctx context.Context, client *vim25.Client, ref types.ManagedObjectReference) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	netObj := object.NewNetwork(client, ref)
	name, err := netObj.ObjectName(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "could not get network name for %s", ref.String())
	}
	return name, nil
}

// GetNetworkMo returns the unique Managed Object for given network name inside of the given Datacenter
// and Cluster.
func GetNetworkMo(ctx context.Context, client *vim25.Client, finder Finder, datacenter, cluster, network string) (*types.ManagedObjectReference, error) {
	networks, err := GetClusterNetworks(ctx, finder, datacenter, cluster)
	if err != nil {
		return nil, err
	}
	for _, net := range networks {
		name, err := GetNetworkName(ctx, client, net)
		if err != nil {
			return nil, err
		}
		if name == network {
			return &net, nil
		}
	}

	return nil, errors.Errorf("unable to find network provided")
}

// GetNetworkMoID returns the unique Managed Object ID for given network name inside of the given Datacenter
// and Cluster.
func GetNetworkMoID(ctx context.Context, client *vim25.Client, finder Finder, datacenter, cluster, network string) (string, error) {
	mo, err := GetNetworkMo(ctx, client, finder, datacenter, cluster, network)
	if err != nil {
		return "", err
	}
	return mo.Value, nil
}
