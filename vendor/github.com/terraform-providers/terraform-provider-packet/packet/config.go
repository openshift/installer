package packet

import (
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/packethost/packngo"
)

const (
	consumerToken = "aZ9GmqHTPtxevvFq9SK3Pi2yr9YCbRzduCSXF2SNem5sjB91mDq7Th3ZwTtRqMWZ"
)

type Config struct {
	AuthToken string
}

// Client returns a new client for accessing Packet's API.
func (c *Config) Client() *packngo.Client {
	client := cleanhttp.DefaultClient()
	client.Transport = logging.NewTransport("Packet", client.Transport)
	return packngo.NewClientWithAuth(consumerToken, c.AuthToken, client)
}
