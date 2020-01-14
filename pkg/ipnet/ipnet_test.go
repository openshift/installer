package ipnet

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func assertJSON(t *testing.T, data interface{}, expected string) {
	actualBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(actualBytes)

	if actual != expected {
		t.Fatalf("%s != %s", actual, expected)
	}
}

func TestMarshal(t *testing.T) {
	stdlibIPNet := &net.IPNet{
		IP:   net.IP{192, 168, 0, 10},
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}

	assertJSON(t, stdlibIPNet, "{\"IP\":\"192.168.0.10\",\"Mask\":\"////AA==\"}")
	wrappedIPNet := &IPNet{IPNet: *stdlibIPNet}
	assertJSON(t, wrappedIPNet, "\"192.168.0.10/24\"")
	assertJSON(t, &IPNet{}, "null")
	assertJSON(t, nil, "null")
}

func TestMarshalIPv6(t *testing.T) {
	ipv6 := MustParseCIDR("fd2e:6f44:5dd8:b856::2/64")
	assertJSON(t, ipv6, "\"fd2e:6f44:5dd8:b856::2/64\"")
}

func TestUnmarshal(t *testing.T) {
	for _, ipNetIn := range []*IPNet{
		nil,
		{
			IPNet: net.IPNet{
				IP:   net.IP{192, 168, 0, 10},
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
		},
		{
			IPNet: net.IPNet{
				IP:   net.IP{253, 46, 111, 68, 93, 216, 184, 86, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	} {
		t.Run(ipNetIn.String(), func(t *testing.T) {
			data, err := json.Marshal(ipNetIn)
			if err != nil {
				t.Fatal(err)
			}

			var ipNetOut *IPNet
			err = json.Unmarshal(data, &ipNetOut)
			if err != nil {
				t.Fatal(err)
			}

			if ipNetOut.String() != ipNetIn.String() {
				t.Fatalf("%v != %v", ipNetOut, ipNetIn)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	ipv6 := MustParseCIDR("fd2e:6f44:5dd8:b856::0/64")
	assert.Equal(t, 6, ipv6.Version())

	ipv4 := MustParseCIDR("172.22.0.0/24")
	assert.Equal(t, 4, ipv4.Version())
}

func TestCIDR(t *testing.T) {
	ipv6 := MustParseCIDR("fd2e:6f44:5dd8:b856::0/64")
	assert.Equal(t, 64, ipv6.CIDR())

	ipv4 := MustParseCIDR("172.22.0.0/24")
	assert.Equal(t, 24, ipv4.CIDR())
}
