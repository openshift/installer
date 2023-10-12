package ipnet

import (
	"encoding/json"
	"net"
	"reflect"
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

func Test_SplitInto(t *testing.T) {
	tests := []struct {
		count    uint
		parent   string
		expected []string
	}{
		{
			count:    8,
			parent:   "1.2.3.0/24",
			expected: []string{"1.2.3.0/27", "1.2.3.32/27", "1.2.3.64/27", "1.2.3.96/27", "1.2.3.128/27", "1.2.3.160/27", "1.2.3.192/27", "1.2.3.224/27"},
		},
		{
			count:    8,
			parent:   "1.2.3.0/27",
			expected: []string{"1.2.3.0/30", "1.2.3.4/30", "1.2.3.8/30", "1.2.3.12/30", "1.2.3.16/30", "1.2.3.20/30", "1.2.3.24/30", "1.2.3.28/30"},
		},
	}
	for _, test := range tests {
		parent, err := ParseCIDR(test.parent)
		if err != nil {
			t.Fatalf("error parsing parent cidr %q: %v", test.parent, err)
		}

		subnets, err := SplitInto(test.count, parent)
		if err != nil {
			t.Fatalf("error splitting parent cidr %q: %v", parent, err)
		}

		var actual []string
		for _, subnet := range subnets {
			actual = append(actual, subnet.String())
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Fatalf("unexpected result of split: actual=%v, expected=%v", actual, test.expected)
		}
	}
}
