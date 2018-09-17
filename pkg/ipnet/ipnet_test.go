package ipnet

import (
	"encoding/json"
	"net"
	"testing"

	yaml "gopkg.in/yaml.v2"
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

func assertYAML(t *testing.T, data interface{}, expected string) {
	actualBytes, err := yaml.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(actualBytes)

	if actual != expected {
		t.Fatalf("%s != %s", actual, expected)
	}
}

func TestMarshalJSON(t *testing.T) {
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

func TestUnmarshalJSON(t *testing.T) {
	for _, ipNetIn := range []*IPNet{
		nil,
		{IPNet: net.IPNet{
			IP:   net.IP{192, 168, 0, 10},
			Mask: net.IPv4Mask(255, 255, 255, 0),
		}},
	} {
		data, err := json.Marshal(ipNetIn)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(string(data), func(t *testing.T) {
			var ipNetOut *IPNet
			err := json.Unmarshal(data, &ipNetOut)
			if err != nil {
				t.Fatal(err)
			}

			if ipNetIn == nil {
				if ipNetOut != nil {
					t.Fatalf("%v != %v", ipNetOut, ipNetIn)
				}
			} else if ipNetOut.String() != ipNetIn.String() {
				t.Fatalf("%v != %v", ipNetOut, ipNetIn)
			}
		})
	}
}

func TestMarshalYAML(t *testing.T) {
	stdlibIPNet := &net.IPNet{
		IP:   net.IP{192, 168, 0, 10},
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}
	assertYAML(t, stdlibIPNet, `ip: 192.168.0.10
mask:
- 255
- 255
- 255
- 0
`)
	wrappedIPNet := &IPNet{IPNet: *stdlibIPNet}
	assertYAML(t, wrappedIPNet, "192.168.0.10/24\n")
	assertYAML(t, &IPNet{}, "null\n")
	assertYAML(t, nil, "null\n")
}

func TestUnmarshalYAML(t *testing.T) {
	for _, ipNetIn := range []*IPNet{
		nil,
		{IPNet: net.IPNet{
			IP:   net.IP{192, 168, 0, 10},
			Mask: net.IPv4Mask(255, 255, 255, 0),
		}},
	} {
		data, err := yaml.Marshal(ipNetIn)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(string(data), func(t *testing.T) {
			var ipNetOut *IPNet
			err := yaml.Unmarshal(data, &ipNetOut)
			if err != nil {
				t.Fatal(err)
			}

			if ipNetIn == nil {
				if ipNetOut != nil {
					t.Fatalf("%v != %v", ipNetOut, ipNetIn)
				}
			} else if ipNetOut.String() != ipNetIn.String() {
				t.Fatalf("%v != %v", ipNetOut, ipNetIn)
			}
		})
	}
}
