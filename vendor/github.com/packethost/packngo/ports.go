package packngo

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const portBasePath = "/ports"

// DevicePortService handles operations on a port which belongs to a particular device
type DevicePortService interface {
	Assign(*PortAssignRequest) (*Port, *Response, error)
	Unassign(*PortAssignRequest) (*Port, *Response, error)
	AssignNative(*PortAssignRequest) (*Port, *Response, error)
	UnassignNative(string) (*Port, *Response, error)
	Bond(*Port, bool) (*Port, *Response, error)
	Disbond(*Port, bool) (*Port, *Response, error)
	DeviceToNetworkType(string, string) (*Device, error)
	DeviceNetworkType(string) (string, error)
	PortToLayerTwo(string, string) (*Port, *Response, error)
	PortToLayerThree(string, string) (*Port, *Response, error)
	GetPortByName(string, string) (*Port, error)
	GetOddEthPorts(*Device) (map[string]*Port, error)
	GetAllEthPorts(*Device) (map[string]*Port, error)
	ConvertDevice(*Device, string) error
}

type PortData struct {
	MAC    string `json:"mac"`
	Bonded bool   `json:"bonded"`
}

type BondData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Port struct {
	ID                      string           `json:"id"`
	Type                    string           `json:"type"`
	Name                    string           `json:"name"`
	Data                    PortData         `json:"data"`
	NetworkType             string           `json:"network_type,omitempty"`
	NativeVirtualNetwork    *VirtualNetwork  `json:"native_virtual_network"`
	AttachedVirtualNetworks []VirtualNetwork `json:"virtual_networks"`
	Bond                    *BondData        `json:"bond"`
}

type AddressRequest struct {
	AddressFamily int  `json:"address_family"`
	Public        bool `json:"public"`
}

type BackToL3Request struct {
	RequestIPs []AddressRequest `json:"request_ips"`
}

type DevicePortServiceOp struct {
	client *Client
}

type PortAssignRequest struct {
	PortID           string `json:"id"`
	VirtualNetworkID string `json:"vnid"`
}

type BondRequest struct {
	PortID     string `json:"id"`
	BulkEnable bool   `json:"bulk_enable"`
}

type DisbondRequest struct {
	PortID      string `json:"id"`
	BulkDisable bool   `json:"bulk_disable"`
}

func (i *DevicePortServiceOp) GetPortByName(deviceID, name string) (*Port, error) {
	device, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return nil, err
	}
	return device.GetPortByName(name)
}

func (i *DevicePortServiceOp) Assign(par *PortAssignRequest) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/assign", portBasePath, par.PortID)
	return i.portAction(path, par)
}

func (i *DevicePortServiceOp) AssignNative(par *PortAssignRequest) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/native-vlan", portBasePath, par.PortID)
	return i.portAction(path, par)
}

func (i *DevicePortServiceOp) UnassignNative(portID string) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/native-vlan", portBasePath, portID)
	port := new(Port)

	resp, err := i.client.DoRequest("DELETE", path, nil, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

func (i *DevicePortServiceOp) Unassign(par *PortAssignRequest) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/unassign", portBasePath, par.PortID)
	return i.portAction(path, par)
}

func (i *DevicePortServiceOp) Bond(p *Port, be bool) (*Port, *Response, error) {
	if p.Data.Bonded {
		return p, nil, nil
	}
	br := &BondRequest{PortID: p.ID, BulkEnable: be}
	path := fmt.Sprintf("%s/%s/bond", portBasePath, br.PortID)
	return i.portAction(path, br)
}

func (i *DevicePortServiceOp) Disbond(p *Port, bd bool) (*Port, *Response, error) {
	if !p.Data.Bonded {
		return p, nil, nil
	}
	dr := &DisbondRequest{PortID: p.ID, BulkDisable: bd}
	path := fmt.Sprintf("%s/%s/disbond", portBasePath, dr.PortID)
	return i.portAction(path, dr)
}

func (i *DevicePortServiceOp) portAction(path string, req interface{}) (*Port, *Response, error) {
	port := new(Port)

	resp, err := i.client.DoRequest("POST", path, req, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

func (i *DevicePortServiceOp) PortToLayerTwo(deviceID, portName string) (*Port, *Response, error) {
	p, err := i.client.DevicePorts.GetPortByName(deviceID, portName)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasPrefix(p.NetworkType, "layer2") {
		return p, nil, nil
	}
	path := fmt.Sprintf("%s/%s/convert/layer-2", portBasePath, p.ID)
	port := new(Port)

	resp, err := i.client.DoRequest("POST", path, nil, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

func (i *DevicePortServiceOp) PortToLayerThree(deviceID, portName string) (*Port, *Response, error) {
	p, err := i.client.DevicePorts.GetPortByName(deviceID, portName)
	if err != nil {
		return nil, nil, err
	}
	if (p.NetworkType == NetworkTypeL3) || (p.NetworkType == NetworkTypeHybrid) {
		return p, nil, nil
	}
	path := fmt.Sprintf("%s/%s/convert/layer-3", portBasePath, p.ID)
	port := new(Port)

	req := BackToL3Request{
		RequestIPs: []AddressRequest{
			{AddressFamily: 4, Public: true},
			{AddressFamily: 4, Public: false},
			{AddressFamily: 6, Public: true},
		},
	}

	resp, err := i.client.DoRequest("POST", path, &req, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

func (i *DevicePortServiceOp) DeviceNetworkType(deviceID string) (string, error) {
	d, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return "", err
	}
	return d.GetNetworkType(), nil
}

func (i *DevicePortServiceOp) GetAllEthPorts(d *Device) (map[string]*Port, error) {
	d, _, err := i.client.Devices.Get(d.ID, nil)
	if err != nil {
		return nil, err
	}
	return d.GetPhysicalPorts(), nil
}

func (i *DevicePortServiceOp) GetOddEthPorts(d *Device) (map[string]*Port, error) {
	d, _, err := i.client.Devices.Get(d.ID, nil)
	if err != nil {
		return nil, err
	}
	ret := map[string]*Port{}
	eth1, err := d.GetPortByName("eth1")
	if err != nil {
		return nil, err
	}
	ret["eth1"] = eth1

	eth3, err := d.GetPortByName("eth3")
	if err != nil {
		return ret, nil
	}
	ret["eth3"] = eth3
	return ret, nil

}

func (i *DevicePortServiceOp) ConvertDevice(d *Device, targetType string) error {
	bondPorts := d.GetBondPorts()

	if targetType == NetworkTypeL3 {
		// TODO: remove vlans from all the ports
		for _, p := range bondPorts {
			_, _, err := i.client.DevicePorts.Bond(p, false)
			if err != nil {
				return err
			}
		}
		_, _, err := i.client.DevicePorts.PortToLayerThree(d.ID, "bond0")
		if err != nil {
			return err
		}
		allEthPorts, err := i.client.DevicePorts.GetAllEthPorts(d)
		if err != nil {
			return err
		}
		for _, p := range allEthPorts {
			_, _, err := i.client.DevicePorts.Bond(p, false)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeHybrid {
		for _, p := range bondPorts {
			_, _, err := i.client.DevicePorts.Bond(p, false)
			if err != nil {
				return err
			}
		}

		_, _, err := i.client.DevicePorts.PortToLayerThree(d.ID, "bond0")
		if err != nil {
			return err
		}

		// ports need to be refreshed before bonding/disbonding
		oddEthPorts, err := i.client.DevicePorts.GetOddEthPorts(d)
		if err != nil {
			return err
		}

		for _, p := range oddEthPorts {
			_, _, err := i.client.DevicePorts.Disbond(p, false)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeL2Individual {
		_, _, err := i.client.DevicePorts.PortToLayerTwo(d.ID, "bond0")
		if err != nil {
			return err
		}
		for _, p := range bondPorts {
			_, _, err = i.client.DevicePorts.Disbond(p, true)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeL2Bonded {

		for _, p := range bondPorts {
			_, _, err := i.client.DevicePorts.PortToLayerTwo(d.ID, p.Name)
			if err != nil {
				return err
			}
		}
		allEthPorts, err := i.client.DevicePorts.GetAllEthPorts(d)
		if err != nil {
			return err
		}
		for _, p := range allEthPorts {
			_, _, err := i.client.DevicePorts.Bond(p, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// waitDeviceNetworkType waits for a device's computed network type (as
// determined by GetNetworkType()) to reach the specified state. An error will
// be returned if the device does not attain the desired network type state when
// the timeout is reached without
func waitDeviceNetworkType(id, networkType string, c *Client) (*Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15)*time.Minute)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d, _, err := c.Devices.Get(id, nil)
			if err != nil {
				return nil, err
			}
			if d.GetNetworkType() == networkType {
				return d, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("device %s is still not in state %s after timeout", id, networkType)
		}
	}
}

func (i *DevicePortServiceOp) DeviceToNetworkType(deviceID string, targetType string) (*Device, error) {

	d, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return nil, err
	}

	curType := d.GetNetworkType()

	if curType == targetType {
		return nil, fmt.Errorf("Device already is in state %s", targetType)
	}
	err = i.client.DevicePorts.ConvertDevice(d, targetType)
	if err != nil {
		return nil, err
	}

	d, _, err = i.client.Devices.Get(deviceID, nil)
	//d, err = waitDeviceNetworkType(deviceID, targetType, i.client)
	if err != nil {
		return nil, err
	}

	finalType := d.GetNetworkType()

	if finalType != targetType {
		return nil, fmt.Errorf(
			"Failed to convert device %s from %s to %s. New type was %s",
			deviceID, curType, targetType, finalType)

	}
	return d, err
}
