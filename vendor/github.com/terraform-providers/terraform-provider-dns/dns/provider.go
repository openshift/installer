package dns

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/miekg/dns"
)

const (
	defaultPort      = 53
	defaultRetries   = 3
	defaultTimeout   = "0"
	defaultTransport = "udp"
)

// Provider returns a schema.Provider for DNS dynamic updates.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"update": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("DNS_UPDATE_SERVER", nil),
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							DefaultFunc: func() (interface{}, error) {
								if envPortStr := os.Getenv("DNS_UPDATE_PORT"); envPortStr != "" {
									port, err := strconv.Atoi(envPortStr)
									if err != nil {
										err = fmt.Errorf("invalid DNS_UPDATE_PORT environment variable: %s", err)
									}
									return port, err
								}

								return defaultPort, nil
							},
						},
						"transport": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("DNS_UPDATE_TRANSPORT", defaultTransport),
						},
						"timeout": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("DNS_UPDATE_TIMEOUT", defaultTimeout),
						},
						"retries": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							DefaultFunc: func() (interface{}, error) {
								if env := os.Getenv("DNS_UPDATE_RETRIES"); env != "" {
									retries, err := strconv.Atoi(env)
									if err != nil {
										err = fmt.Errorf("invalid DNS_UPDATE_RETRIES environment variable: %s", err)
									}
									return retries, err
								}

								return defaultRetries, nil
							},
						},
						"key_name": {
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("DNS_UPDATE_KEYNAME", nil),
						},
						"key_algorithm": {
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("DNS_UPDATE_KEYALGORITHM", nil),
						},
						"key_secret": {
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("DNS_UPDATE_KEYSECRET", nil),
						},
					},
				},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"dns_a_record_set":     dataSourceDnsARecordSet(),
			"dns_aaaa_record_set":  dataSourceDnsAAAARecordSet(),
			"dns_cname_record_set": dataSourceDnsCnameRecordSet(),
			"dns_mx_record_set":    dataSourceDnsMXRecordSet(),
			"dns_ns_record_set":    dataSourceDnsNSRecordSet(),
			"dns_ptr_record_set":   dataSourceDnsPtrRecordSet(),
			"dns_srv_record_set":   dataSourceDnsSRVRecordSet(),
			"dns_txt_record_set":   dataSourceDnsTxtRecordSet(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"dns_a_record_set":    resourceDnsARecordSet(),
			"dns_aaaa_record_set": resourceDnsAAAARecordSet(),
			"dns_cname_record":    resourceDnsCnameRecord(),
			"dns_mx_record_set":   resourceDnsMXRecordSet(),
			"dns_ns_record_set":   resourceDnsNSRecordSet(),
			"dns_ptr_record":      resourceDnsPtrRecord(),
			"dns_srv_record_set":  resourceDnsSRVRecordSet(),
			"dns_txt_record_set":  resourceDnsTXTRecordSet(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	var server, transport, timeout, keyname, keyalgo, keysecret string
	var port, retries int
	var duration time.Duration

	// if the update block is missing, schema.EnvDefaultFunc is not called
	if v, ok := d.GetOk("update"); ok {
		update := v.([]interface{})[0].(map[string]interface{})
		if val, ok := update["port"]; ok {
			port = int(val.(int))
		}
		if val, ok := update["server"]; ok {
			server = val.(string)
		}
		if val, ok := update["transport"]; ok {
			transport = val.(string)
		}
		if val, ok := update["timeout"]; ok {
			timeout = val.(string)
		}
		if val, ok := update["retries"]; ok {
			retries = int(val.(int))
		}
		if val, ok := update["key_name"]; ok {
			keyname = val.(string)
		}
		if val, ok := update["key_algorithm"]; ok {
			keyalgo = val.(string)
		}
		if val, ok := update["key_secret"]; ok {
			keysecret = val.(string)
		}
	} else {
		if len(os.Getenv("DNS_UPDATE_SERVER")) > 0 {
			server = os.Getenv("DNS_UPDATE_SERVER")
		} else {
			return nil, nil
		}
		if len(os.Getenv("DNS_UPDATE_PORT")) > 0 {
			var err error
			portStr := os.Getenv("DNS_UPDATE_PORT")
			port, err = strconv.Atoi(portStr)
			if err != nil {
				return nil, fmt.Errorf("invalid DNS_UPDATE_PORT environment variable: %s", err)
			}
		} else {
			port = defaultPort
		}
		if len(os.Getenv("DNS_UPDATE_TRANSPORT")) > 0 {
			transport = os.Getenv("DNS_UPDATE_TRANSPORT")
		} else {
			transport = defaultTransport
		}
		if len(os.Getenv("DNS_UPDATE_TIMEOUT")) > 0 {
			timeout = os.Getenv("DNS_UPDATE_TIMEOUT")
		} else {
			timeout = defaultTimeout
		}
		if len(os.Getenv("DNS_UPDATE_RETRIES")) > 0 {
			var err error
			env := os.Getenv("DNS_UPDATE_RETRIES")
			retries, err = strconv.Atoi(env)
			if err != nil {
				return nil, fmt.Errorf("invalid DNS_UPDATE_RETRIES environment variable: %s", err)
			}
		} else {
			retries = defaultRetries
		}
		if len(os.Getenv("DNS_UPDATE_KEYNAME")) > 0 {
			keyname = os.Getenv("DNS_UPDATE_KEYNAME")
		}
		if len(os.Getenv("DNS_UPDATE_KEYALGORITHM")) > 0 {
			keyalgo = os.Getenv("DNS_UPDATE_KEYALGORITHM")
		}
		if len(os.Getenv("DNS_UPDATE_KEYSECRET")) > 0 {
			keysecret = os.Getenv("DNS_UPDATE_KEYSECRET")
		}
	}

	if timeout != "" {
		var err error
		// Try parsing as a duration
		duration, err = time.ParseDuration(timeout)
		if err != nil {
			// Failing that, convert to an integer and treat as seconds
			seconds, err := strconv.Atoi(timeout)
			if err != nil {
				return nil, fmt.Errorf("invalid timeout: %s", timeout)
			}
			duration = time.Duration(seconds) * time.Second
		}
		if duration < 0 {
			return nil, fmt.Errorf("timeout cannot be negative: %s", duration)
		}
	}

	config := Config{
		server:    server,
		port:      port,
		transport: transport,
		timeout:   duration,
		retries:   retries,
		keyname:   keyname,
		keyalgo:   keyalgo,
		keysecret: keysecret,
	}

	return config.Client()
}

func getAVal(record interface{}) (string, int, error) {

	_, ok := record.(*dns.A)
	if !ok {
		return "", 0, fmt.Errorf("didn't get a A record")
	}

	recstr := record.(*dns.A).String()
	var name, class, typ, addr string
	var ttl int

	_, err := fmt.Sscanf(recstr, "%s\t%d\t%s\t%s\t%s", &name, &ttl, &class, &typ, &addr)
	if err != nil {
		return "", 0, fmt.Errorf("Error parsing record: %s", err)
	}

	return addr, ttl, nil
}

func getNSVal(record interface{}) (string, int, error) {

	_, ok := record.(*dns.NS)
	if !ok {
		return "", 0, fmt.Errorf("didn't get a NS record")
	}

	recstr := record.(*dns.NS).String()
	var name, class, typ, nameserver string
	var ttl int

	_, err := fmt.Sscanf(recstr, "%s\t%d\t%s\t%s\t%s", &name, &ttl, &class, &typ, &nameserver)
	if err != nil {
		return "", 0, fmt.Errorf("Error parsing record: %s", err)
	}

	return nameserver, ttl, nil
}

func getAAAAVal(record interface{}) (string, int, error) {

	_, ok := record.(*dns.AAAA)
	if !ok {
		return "", 0, fmt.Errorf("didn't get a AAAA record")
	}

	recstr := record.(*dns.AAAA).String()
	var name, class, typ, addr string
	var ttl int

	_, err := fmt.Sscanf(recstr, "%s\t%d\t%s\t%s\t%s", &name, &ttl, &class, &typ, &addr)
	if err != nil {
		return "", 0, fmt.Errorf("Error parsing record: %s", err)
	}

	return addr, ttl, nil
}

func getCnameVal(record interface{}) (string, int, error) {

	_, ok := record.(*dns.CNAME)
	if !ok {
		return "", 0, fmt.Errorf("didn't get a CNAME record")
	}

	recstr := record.(*dns.CNAME).String()
	var name, class, typ, cname string
	var ttl int

	_, err := fmt.Sscanf(recstr, "%s\t%d\t%s\t%s\t%s", &name, &ttl, &class, &typ, &cname)
	if err != nil {
		return "", 0, fmt.Errorf("Error parsing record: %s", err)
	}

	return cname, ttl, nil
}

func getPtrVal(record interface{}) (string, int, error) {

	_, ok := record.(*dns.PTR)
	if !ok {
		return "", 0, fmt.Errorf("didn't get a PTR record")
	}

	recstr := record.(*dns.PTR).String()
	var name, class, typ, ptr string
	var ttl int

	_, err := fmt.Sscanf(recstr, "%s\t%d\t%s\t%s\t%s", &name, &ttl, &class, &typ, &ptr)
	if err != nil {
		return "", 0, fmt.Errorf("Error parsing record: %s", err)
	}

	return ptr, ttl, nil
}

func isTimeout(err error) bool {

	timeout, ok := err.(net.Error)
	return ok && timeout.Timeout()
}

func exchange(msg *dns.Msg, tsig bool, meta interface{}) (*dns.Msg, error) {

	c := meta.(*DNSClient).c
	srv_addr := meta.(*DNSClient).srv_addr
	keyname := meta.(*DNSClient).keyname
	keyalgo := meta.(*DNSClient).keyalgo
	c.Net = meta.(*DNSClient).transport
	retries := meta.(*DNSClient).retries
	retry_tcp := false

	msg.RecursionDesired = false

Retry:
	if tsig && keyname != "" {
		msg.SetTsig(keyname, keyalgo, 300, time.Now().Unix())
	}

	r, _, err := c.Exchange(msg, srv_addr)

	switch err {
	case dns.ErrTruncated:
		if retry_tcp {
			switch c.Net {
			case "udp":
				c.Net = "tcp"
			case "udp4":
				c.Net = "tcp4"
			case "udp6":
				c.Net = "tcp6"
			default:
				return nil, fmt.Errorf("Unknown transport: %s", c.Net)
			}
		} else {
			msg.SetEdns0(dns.DefaultMsgSize, false)
			retry_tcp = true
		}

		// Reset retries counter on protocol change
		retries = meta.(*DNSClient).retries
		goto Retry
	case nil:
		if r.Rcode == dns.RcodeServerFailure && retries > 0 {
			retries--
			goto Retry
		}
	default:
		if isTimeout(err) && retries > 0 {
			retries--
			goto Retry
		}
	}

	return r, err
}

func resourceDnsImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	record := d.Id()
	if !dns.IsFqdn(record) {
		return nil, fmt.Errorf("Not a fully-qualified DNS name: %s", record)
	}

	labels := dns.SplitDomainName(record)

	msg := new(dns.Msg)

	var zone *string

Loop:
	for l := range labels {

		msg.SetQuestion(dns.Fqdn(strings.Join(labels[l:], ".")), dns.TypeSOA)

		r, err := exchange(msg, true, meta)
		if err != nil {
			return nil, fmt.Errorf("Error querying DNS record: %s", err)
		}

		switch r.Rcode {
		case dns.RcodeSuccess:

			if len(r.Answer) == 0 {
				continue
			}

			for _, ans := range r.Answer {
				switch t := ans.(type) {
				case *dns.SOA:
					zone = &t.Hdr.Name
				case *dns.CNAME:
					continue Loop
				}
			}

			break Loop
		case dns.RcodeNameError:
			continue
		default:
			return nil, fmt.Errorf("Error querying DNS record: %v (%s)", r.Rcode, dns.RcodeToString[r.Rcode])
		}
	}

	if zone == nil {
		return nil, fmt.Errorf("No SOA record in authority section in response for %s", record)
	}

	common := dns.CompareDomainName(record, *zone)
	if common == 0 {
		return nil, fmt.Errorf("DNS record %s shares no common labels with zone %s", record, *zone)
	}

	d.Set("zone", *zone)
	if name := strings.Join(labels[:len(labels)-common], "."); name != "" {
		d.Set("name", name)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceFQDN(d *schema.ResourceData) string {

	fqdn := d.Get("zone").(string)

	if name, ok := d.GetOk("name"); ok {
		fqdn = fmt.Sprintf("%s.%s", name.(string), fqdn)
	}

	return fqdn
}

func resourceDnsRead(d *schema.ResourceData, meta interface{}, rrType uint16) ([]dns.RR, error) {

	if meta != nil {

		fqdn := resourceFQDN(d)

		msg := new(dns.Msg)
		msg.SetQuestion(fqdn, rrType)

		r, err := exchange(msg, true, meta)
		if err != nil {
			return nil, fmt.Errorf("Error querying DNS record: %s", err)
		}
		switch r.Rcode {
		case dns.RcodeSuccess:
			// NS records are returned slightly differently
			if (rrType == dns.TypeNS && len(r.Ns) > 0) || len(r.Answer) > 0 {
				break
			}
			fallthrough
		case dns.RcodeNameError:
			return nil, nil
		default:
			return nil, fmt.Errorf("Error querying DNS record: %v (%s)", r.Rcode, dns.RcodeToString[r.Rcode])
		}

		if rrType == dns.TypeNS {
			return r.Ns, nil
		}
		return r.Answer, nil
	} else {
		return nil, fmt.Errorf("update server is not set")
	}
}

func resourceDnsDelete(d *schema.ResourceData, meta interface{}, rrType uint16) error {

	if meta != nil {

		fqdn := resourceFQDN(d)

		msg := new(dns.Msg)

		msg.SetUpdate(d.Get("zone").(string))

		rr, _ := dns.NewRR(fmt.Sprintf("%s 0 %s", fqdn, dns.TypeToString[rrType]))
		msg.RemoveRRset([]dns.RR{rr})

		r, err := exchange(msg, true, meta)
		if err != nil {
			return fmt.Errorf("Error deleting DNS record: %s", err)
		}
		if r.Rcode != dns.RcodeSuccess {
			return fmt.Errorf("Error deleting DNS record: %v (%s)", r.Rcode, dns.RcodeToString[r.Rcode])
		}

		return nil
	} else {
		return fmt.Errorf("update server is not set")
	}
}
