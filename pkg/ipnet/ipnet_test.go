package ipnet

import (
	"encoding/json"
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

func TestUnmarshal(t *testing.T) {
	for _, ipNetIn := range []*IPNet{
		nil,
		{IPNet: net.IPNet{
			IP:   net.IP{192, 168, 0, 10},
			Mask: net.IPv4Mask(255, 255, 255, 0),
		}},
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
