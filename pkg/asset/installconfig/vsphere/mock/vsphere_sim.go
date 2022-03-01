package mock

import (
	"context"
	"crypto/tls"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

// StartSimulator starts an instance of the simulator which listens on 127.0.0.1. Call GetClient
//                to retrieve a vim25.client which will connect to and trust this simulator
func StartSimulator() *simulator.Server {
	model := simulator.VPX()
	model.Create()
	model.Service.TLS = new(tls.Config)
	model.Service.TLS.ServerName = "127.0.0.1"
	server := model.Service.NewServer()
	return server
}

// GetClient returns a vim25 client which connects to and trusts the simulator
func GetClient(server *simulator.Server) (*vim25.Client, *session.Manager, error) {
	tmpCAdir := "/tmp/vcsimca"
	os.Mkdir(tmpCAdir, os.ModePerm)
	pemBlock := pem.Block{
		Type:    "CERTIFICATE",
		Headers: nil,
		Bytes:   server.TLS.Certificates[0].Certificate[0],
	}
	tempFile, err := ioutil.TempFile(tmpCAdir, "*.pem")
	if err != nil {
		return nil, nil, err
	}
	tempFile.Write(pem.EncodeToMemory(&pemBlock))

	soapClient := soap.NewClient(server.URL, false)
	soapClient.SetRootCAs(tempFile.Name())
	vimClient, err := vim25.NewClient(context.TODO(), soapClient)
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
