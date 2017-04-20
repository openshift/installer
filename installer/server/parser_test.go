package server

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalNode(t *testing.T) {
	data := []byte(`{"name": "node1.example.com", "mac": "52:54:00:a1:9c:ae"}`)
	node := new(Node)
	err := json.Unmarshal(data, node)
	assert.Nil(t, err)
	assert.Equal(t, "52:54:00:a1:9c:ae", node.MAC.String())
	assert.Equal(t, "node1.example.com", node.Name)
}

func TestUnmarshalMACAddress(t *testing.T) {
	cases := []struct {
		data []byte
		// normalized MAC string
		address string
	}{
		{[]byte("\"52:54:00:a1:9c:ae\""), "52:54:00:a1:9c:ae"},
		{[]byte("\"52:54:00:b2:2f:86\""), "52:54:00:b2:2f:86"},
		{[]byte("\"52:54:00:c3:61:77\""), "52:54:00:c3:61:77"},
		{[]byte("\"52:54:00:d7:99:c7\""), "52:54:00:d7:99:c7"},
		{[]byte("\"52-54-00-a1-9c-ae\""), "52:54:00:a1:9c:ae"},
	}

	for _, c := range cases {
		addr := new(macAddr)
		json.Unmarshal(c.data, addr)
		assert.Equal(t, c.address, addr.String())
	}
}

func TestDashString(t *testing.T) {
	cases := []struct {
		data []byte
		// normalized, dashed MAC string
		address string
	}{
		{[]byte("\"52:54:00:A1:9c:ae\""), "52-54-00-a1-9c-ae"},
		{[]byte("\"52-54-00-a1-9c-ae\""), "52-54-00-a1-9c-ae"},
	}

	for _, c := range cases {
		addr := new(macAddr)
		json.Unmarshal(c.data, addr)
		assert.Equal(t, c.address, addr.DashString())
	}

}
