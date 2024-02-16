package mock

import (
	"context"
	"crypto/tls"
	"encoding/pem"
	"errors"
	"io/fs"
	"os"
	"strconv"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/simulator/esx"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"

	// required to initialize the REST endpoint.
	_ "github.com/vmware/govmomi/vapi/rest"
	// required to initialize the VAPI endpoint.
	_ "github.com/vmware/govmomi/vapi/simulator"
)

const (
	esxi7U2BuildNumber    int    = 17630552
	esxi7U2Version        string = "7.0.2"
	vcenter7U2BuildNumber int    = 17694817
	vcenter7U2Version     string = "7.0.2"
)

// StartSimulator starts an instance of the simulator which listens on 127.0.0.1.
// Call GetClient to retrieve a vim25.client which will connect to and trust this
// simulator
func StartSimulator(setVersionToSupported bool) (*simulator.Server, error) {
	model := simulator.VPX()

	// Change the simulated vCenter and ESXi hosts
	// to the version and build we support.
	if setVersionToSupported {
		esx.HostSystem.Config.Product.Build = strconv.Itoa(esxi7U2BuildNumber)
		esx.HostSystem.Config.Product.Version = esxi7U2Version
		model.ServiceContent.About.Build = strconv.Itoa(vcenter7U2BuildNumber)
		model.ServiceContent.About.Version = vcenter7U2Version
	}

	model.Folder = 1
	model.Datacenter = 5
	model.OpaqueNetwork = 1
	err := model.Create()
	if err != nil {
		return nil, err
	}

	model.Service.TLS = new(tls.Config)

	model.Service.TLS.ServerName = "127.0.0.1"
	model.Service.RegisterEndpoints = true
	server := model.Service.NewServer()
	return server, nil
}

// GetClient returns a vim25 client which connects to and trusts the simulator
func GetClient(server *simulator.Server) (*vim25.Client, *session.Manager, error) {
	tmpCAdir := "/tmp/vcsimca"
	err := os.Mkdir(tmpCAdir, os.ModePerm)

	if err != nil {
		// If the error is not file existing return err
		if !errors.Is(err, fs.ErrExist) {
			return nil, nil, err
		}
	}
	pemBlock := pem.Block{
		Type:    "CERTIFICATE",
		Headers: nil,
		Bytes:   server.TLS.Certificates[0].Certificate[0],
	}
	tempFile, err := os.CreateTemp(tmpCAdir, "*.pem")
	if err != nil {
		return nil, nil, err
	}
	_, err = tempFile.Write(pem.EncodeToMemory(&pemBlock))
	if err != nil {
		return nil, nil, err
	}

	soapClient := soap.NewClient(server.URL, false)
	err = soapClient.SetRootCAs(tempFile.Name())
	if err != nil {
		return nil, nil, err
	}
	vimClient, err := vim25.NewClient(context.TODO(), soapClient)
	if err != nil {
		return nil, nil, err
	}
	sessionMgr := session.NewManager(vimClient)
	if sessionMgr == nil {
		return nil, nil, errors.New("unable to retrieve session manager")
	}
	if server.URL.User != nil {
		err = sessionMgr.Login(context.TODO(), server.URL.User)
		if err != nil {
			return nil, nil, err
		}
	}
	return vimClient, sessionMgr, err
}

// GetFinder returns an object finder
func GetFinder(server *simulator.Server) (*find.Finder, error) {
	client, _, err := GetClient(server)
	if err != nil {
		return nil, err
	}
	return find.NewFinder(client), nil
}
