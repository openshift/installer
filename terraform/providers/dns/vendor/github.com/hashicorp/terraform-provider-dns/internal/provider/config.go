package provider

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/bodgit/tsig"
	"github.com/bodgit/tsig/gss"
	"github.com/miekg/dns"
)

type Config struct {
	server    string
	port      int
	transport string
	timeout   time.Duration
	retries   int
	keyname   string
	keyalgo   string
	keysecret string
	gssapi    bool
	realm     string
	username  string
	password  string
	keytab    string
}

type DNSClient struct {
	c         *dns.Client
	srv_addr  string
	transport string
	retries   int
	keyname   string
	keysecret string
	keyalgo   string
	gssClient *gss.Client
	realm     string
	username  string
	password  string
	keytab    string
}

// Configures and returns a fully initialized DNSClient
func (c *Config) Client() (interface{}, error) {
	log.Println("[INFO] Building DNSClient config structure")

	var client DNSClient
	client.srv_addr = net.JoinHostPort(c.server, strconv.Itoa(c.port))

	// This block is a little unwieldy but there are a few combinations of
	// settings we need to check for
	if c.gssapi && // GSSAPI requested
		!(c.realm != "" && ((c.username == "" && c.password == "" && c.keytab == "") || // Rely on current user session
			(c.username != "" && (c.password != "" || c.keytab != "")))) { // Supplied credentials with either password or keytab
		return nil, fmt.Errorf("Error configuring provider: when using GSSAPI, \"realm\", \"username\" and either \"password\" or \"keytab\" should be non empty")
	} else if !((c.keyname == "" && c.keysecret == "" && c.keyalgo == "") || // No TSIG required
		(c.keyname != "" && c.keysecret != "" && c.keyalgo != "")) { // Supplied key name, secret and algorithm
		return nil, fmt.Errorf("Error configuring provider: when using authentication, \"key_name\", \"key_secret\" and \"key_algorithm\" should be non empty")
	}

	client.c = new(dns.Client)
	client.c.Net = c.transport
	client.transport = c.transport
	client.c.Timeout = c.timeout
	client.retries = c.retries
	client.realm = c.realm
	client.username = c.username
	client.password = c.password
	client.keytab = c.keytab
	if !c.gssapi && c.keyname != "" {
		if !dns.IsFqdn(c.keyname) {
			return nil, fmt.Errorf("Error configuring provider: \"key_name\" should be fully-qualified")
		}
		keyname := strings.ToLower(c.keyname)
		client.keyname = keyname
		client.keysecret = c.keysecret
		keyalgo, err := convertHMACAlgorithm(c.keyalgo)
		if err != nil {
			return nil, fmt.Errorf("Error configuring provider: %s", err)
		}
		client.keyalgo = keyalgo
		client.c.TsigProvider = tsig.HMAC{keyname: c.keysecret}
	} else if c.gssapi {
		g, err := gss.NewClient(client.c)
		if err != nil {
			return nil, fmt.Errorf("Error initializing GSS library: %s", err)
		}

		client.gssClient = g
		client.keyalgo = tsig.GSS
		client.c.TsigProvider = g
	}
	return &client, nil
}

// Validates and converts HMAC algorithm
func convertHMACAlgorithm(name string) (string, error) {
	switch name {
	case "hmac-md5":
		return dns.HmacMD5, nil
	case "hmac-sha1":
		return dns.HmacSHA1, nil
	case "hmac-sha256":
		return dns.HmacSHA256, nil
	case "hmac-sha512":
		return dns.HmacSHA512, nil
	default:
		return "", fmt.Errorf("Unknown HMAC algorithm: %s", name)
	}
}
