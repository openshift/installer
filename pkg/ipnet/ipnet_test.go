package ipnet

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
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
	tests := []struct {
		input string

		exp    *IPNet
		expErr string
	}{{
		input:  `""`,
		expErr: "failed to Parse cidr string to net.IPNet: invalid CIDR address: ",
	}, {
		input:  ``,
		expErr: "unexpected end of JSON input",
	}, {
		input: `null`,
		exp:   &IPNet{net.IPNet{IP: net.IP{}, Mask: net.IPMask{}}},
	}, {
		input:  `"null"`,
		expErr: "failed to Parse cidr string to net.IPNet: invalid CIDR address: null",
	}, {
		input: `"192.168.0.10/24"`,
		exp:   &IPNet{net.IPNet{IP: net.IP{0xc0, 0xa8, 0x0, 0x0}, Mask: net.IPMask{0xff, 0xff, 0xff, 0x0}}},
	}, {
		input: `"fe80::c0a8:a/120"`,
		exp:   &IPNet{net.IPNet{IP: net.IP{0xfe, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0xa8, 0x0, 0x0}, Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0}}},
	}}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var ipNetOut IPNet
			err := json.Unmarshal([]byte(test.input), &ipNetOut)
			if test.expErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.exp, &ipNetOut)
			} else {
				assert.EqualError(t, err, test.expErr)
			}
		})
	}
}
