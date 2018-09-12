// Package ipnet wraps net.IPNet to get CIDR serialization.
package ipnet

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

var nullString = "null"
var nullBytes = []byte(nullString)
var emptyIPNet = net.IPNet{}

// IPNet wraps net.IPNet to get CIDR serialization.
type IPNet struct {
	net.IPNet
}

// MarshalJSON interface for an IPNet.
func (ipnet IPNet) MarshalJSON() (data []byte, err error) {
	if reflect.DeepEqual(ipnet.IPNet, emptyIPNet) {
		return nullBytes, nil
	}

	return json.Marshal(ipnet.String())
}

// UnmarshalJSON interface for an IPNet.
func (ipnet *IPNet) UnmarshalJSON(b []byte) (err error) {
	if string(b) == nullString {
		ipnet.IP = net.IP{}
		ipnet.Mask = net.IPMask{}
		return nil
	}

	var cidr string
	err = yaml.Unmarshal(b, &cidr)
	if err != nil {
		return err
	}

	return ipnet.parseCIDR(cidr)
}

func (ipnet *IPNet) parseCIDR(cidr string) (err error) {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	ipnet.IP = ip
	ipnet.Mask = net.Mask
	return nil
}

// MarshalYAML interface for an IPNet.
func (ipnet *IPNet) MarshalYAML() (replacement interface{}, err error) {
	if ipnet == nil || reflect.DeepEqual(ipnet.IPNet, emptyIPNet) {
		return nil, nil
	}

	return ipnet.String(), nil
}

// UnmarshalYAML interface for an IPNet.
func (ipnet *IPNet) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var data interface{}
	err = unmarshal(&data)
	if err != nil {
		return nil
	}

	switch data.(type) {
	case nil:
		return nil
	case string:
		return ipnet.parseCIDR(data.(string))
	default:
		return fmt.Errorf("cannot unmarshal %v into an IPNet", data)
	}
}
