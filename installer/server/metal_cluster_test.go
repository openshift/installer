package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	someMAC = mustMAC(parseMACAddr("52:54:00:a1:9c:ae"))
)

func TestMetalInitialize(t *testing.T) {
	cluster := validTectonicMetalCluster()
	assert.Nil(t, cluster.Initialize())
}

func TestMetalInitialize_MissingControllerDomain(t *testing.T) {
	cluster := validTectonicMetalCluster()
	cluster.ControllerDomain = ""
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errMissingControllerDomain, err)
	}
}

func TestMetalInitialize_TooFewControllers(t *testing.T) {
	cluster := validTectonicMetalCluster()
	// Remove controller node
	cluster.Controllers = []Node{}
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errTooFewControllers, err)
	}
}

func TestMetalInitialize_ClusterTooSmall(t *testing.T) {
	cluster := validTectonicMetalCluster()
	// Remove two nodes
	cluster.Workers = cluster.Workers[2:]
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errClusterTooSmall, err)
	}
}

func TestMetalInitialize_MissingMatchboxEndpoint(t *testing.T) {
	cluster := validTectonicMetalCluster()
	cluster.MatchboxHTTP = ""
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errMissingMatchboxEndpoint, err)
	}
}

func TestMetalInitialize_MissingMAC(t *testing.T) {
	// nil master MAC
	cluster := validTectonicMetalCluster()
	cluster.Controllers[0].MAC = nil
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errMissingMACAddress, err)
	}
	// nil worker MAC
	cluster = validTectonicMetalCluster()
	cluster.Workers[0].MAC = nil
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errMissingMACAddress, err)
	}
}

func TestMetalInitialize_MissingChannel(t *testing.T) {
	cluster := validTectonicMetalCluster()
	cluster.Channel = ""
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errMissingChannel, err)
	}
}

func TestMetalInitialize_MissingVersion(t *testing.T) {
	cluster := validTectonicMetalCluster()
	cluster.Version = ""
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errMissingVersion, err)
	}
}

func TestMetalInitialize_MissingTectonicConfig(t *testing.T) {
	cluster := validTectonicMetalCluster()
	cluster.Tectonic = nil
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errMissingTectonicConfig, err)
	}
}

func TestMetalInitialize_MissingTectonicDomain(t *testing.T) {
	cluster := validTectonicMetalCluster()
	cluster.TectonicDomain = ""
	if err := cluster.Initialize(); assert.Error(t, err) {
		assert.Equal(t, errMissingTectonicDomain, err)
	}
}

// mustMAC is a helper which wraps a call to a function returning a (macAddr,
// error) and panics if the error is non-nil. It should be used in tests only.
//     var mac = mustMac(parseMACAddr("52:54:00:a1:9c:ae"))
func mustMAC(addr macAddr, err error) macAddr {
	if err != nil {
		panic(err)
	}
	return addr
}

// validTectonicMetalCluster returns a TectonicMetalCluster which can be
// modified and used for testing.
func validTectonicMetalCluster() *TectonicMetalCluster {
	someMAC := mustMAC(parseMACAddr("52:54:00:a1:9c:ae"))
	return &TectonicMetalCluster{
		MatchboxHTTP:       "172.18.0.2:8081",
		MatchboxRPC:        "matchbox.example.com:8081",
		MatchboxCA:         "-----BEGIN CERTIFICATE-----\nMIIFDTCCAvWgAwIBAgIJAMTO0L/ekS8KMA0GCSqGSIb3DQEBCwUAMBIxEDAOBgNV\nBAMMB2Zha2UtY2EwHhcNMTYwNjEzMTk1MTI4WhcNMjYwNjExMTk1MTI4WjASMRAw\nDgYDVQQDDAdmYWtlLWNhMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA\nyN01Bd7DwfPgs+9SqoMLdKCe+QcJTxxj14FCDRn4ZwSZKVoSvTf2JBQc+PeUzPNs\nktc++1RzMG5TLYUm+ER7w7zsBiC7SQATlwKgpDybjuKehxMx42GO4IY5+dP3FIyw\nEBvgJOBqeU0Pjhl1sh4fx8RXCmQlHQkNswQ47WEbZvaiTtIEndeE4M8nYOkooy74\nRhxexN794qigHMycqFR2odE26675uDg17ieD1haVtJUfl4aR7xnCuf6UeL/WFOrF\nTxnObrYpvIr/75GCrmXuPrpxaqQZIh3I3Fc16C/7aHGj/2TsaiQl1tqlDbT/liKf\nXu1f7SuYhc+jYb2hwIfD+w6PKRwxlM6GY2wDPYtBwadIECd0Ow7b+bDM2MucWckK\nkc+5AuBBbJ4i9YYmyJqDWI3vaiQwXgbUwR1RgdOnfknU2gKLBWqMAf60w0jC0cKQ\nh/scgGEXv+HHHgYXcZ4L0rutGyjiBhY420p64q2o1e0AlqeuDqCp0pPskfvecE5u\nvv9fpf6VyD2xi70XEeDpy07f3bTL8CkNbGvj0XvH6xlvhuOY7tgsqg7JvZNkZJZ1\n06zYQ6hfhObd9A9sBwXHtegF+kYKU2nUeM0gQBylqqRJ4PAYDQdAsoD6Yc536S5E\nwDSVwI+8Rya12o9KElXvKjDje7OHrkO28cxi/S8pg5kCAwEAAaNmMGQwHQYDVR0O\nBBYEFO48wKIe/blfHFC/59wQgFhhJz5yMB8GA1UdIwQYMBaAFO48wKIe/blfHFC/\n59wQgFhhJz5yMBIGA1UdEwEB/wQIMAYBAf8CAQAwDgYDVR0PAQH/BAQDAgGGMA0G\nCSqGSIb3DQEBCwUAA4ICAQAPi/8RCwwF9TLRqkpVU+YccYjs6ZhJPJCLZ89g7Koo\nTVe6eEW8iVsa2TWs80wB3SjmJy9sqHbjzTmGQCNAEPv/TR0PRNwDE7yutFoXD/Nl\nEVerLpgHj1DkGQ47xE7YrTu3qi65hSRgXL3HmkvZa+ykrtsP+B6HSJDqsqOFgJ/4\n52kZIhhZgp4XHyaMLi2vIIGXlddRBy/xudic5mD6yiKqZwr6xU4vRrFgkiXg6V0X\nDuuZTl0WXQ4jMWL0yb3mAGaj+X0TR/BMkce9Wn9xxB9MkgKh+Ca8WF74dhVTNwIR\n4EUO++N00hS1dJUczrJLsKxaHiLxBI0/aQTU4fb1pfxhqOzSEWqjgS4IMGoGi0PM\naEhExRmPVYxYHKJMxVitL4/hKBhSngScSaRBJslvzWhCWEV2cHnk/4HoA/rpfR6v\nzHSFmMvzwugzSg39HvSaHQ3DPr2lAt67bDw0uGS6FpuAl0ARsXCcEZpLduWBXM1U\nCQ16jaEuV0PUQyul6Jp3tsg2r3OXB5iiWcwLkKkFh1Iz7Vxt7xRreeHdJogFM5x8\nbY9ZM+JzXIhWelf0R0dP/qsLlaqDHj51xZRQGUmjasoJSpXgLEjkvtCH5wT6neUB\nl6JYmUq58fPg9pjlxdaOZFjf7BBpd1r8znj7JzylQS0JfevJXBNdgpMk93qVu6u1\nVA==\n-----END CERTIFICATE-----\n",
		MatchboxClientCert: "-----BEGIN CERTIFICATE-----\nMIIEYDCCAkigAwIBAgICEAEwDQYJKoZIhvcNAQELBQAwEjEQMA4GA1UEAwwHZmFr\nZS1jYTAeFw0xNjA2MTMxOTUxMzBaFw0xNzA2MTMxOTUxMzBaMBYxFDASBgNVBAMM\nC2Zha2UtY2xpZW50MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Txa\nMnZbC/CyZWe93cbrK5ajyaAA753/srMNQusxXRpBEHzjkOo1MPAwo7FPu/08EvAY\nMYzTbLZdHui5xl/sWiO7V3yU6Ve1XmkWG56fIGwT74s1RlnPZ7HEf/893wvf4aax\nE9pMB6WaN3PuQShTDPTGk18bkDPY+eRzsunnYAUwY5nWAD/vhJgsCyALITtZJhBo\nU2yWMctlD9EKAX0uKnDn+3m3RAxUIcOaxVQEkyPSe3zFRWCPsphU4Xrbwmn28dcI\nkwW6Q/9Gmqs1vrHy+NtCNOJphaHADZn0OuH3p5919onZolF5hKd2uNKBOWeZHf+x\n0T1ZekE54ORcPhTH5wIDAQABo4G7MIG4MAkGA1UdEwQCMAAwEQYJYIZIAYb4QgEB\nBAQDAgeAMDMGCWCGSAGG+EIBDQQmFiRPcGVuU1NMIEdlbmVyYXRlZCBDbGllbnQg\nQ2VydGlmaWNhdGUwHQYDVR0OBBYEFOvgy7ppEipt4pA5S15T692ESPf3MB8GA1Ud\nIwQYMBaAFO48wKIe/blfHFC/59wQgFhhJz5yMA4GA1UdDwEB/wQEAwIF4DATBgNV\nHSUEDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQsFAAOCAgEAi84Eciw2QvYxGb1f\nnXrK7VYgwDceCqfjrDzauJrd6JuQKnAvNlMPce/f+UdyE1F4P+I80IgdgAkekWEE\ngyb/gWnpfWH6fFbRLjE1qTf2QiIuSxuWp2UPpfJlFTtHYoZmbBMFxsXSUCwLeUt0\ncNXq0o0rpc/dvvrBjbGZN3ghjO/GSCavw1rbzbSa5xTCDAVIszm7gcKF92eCFDv0\ntUjkq3yhjBv4O5dybcVm+CtIDFZQPJIhWKNA5v1AroFbGn5sfqif2+JqSRKTB8cI\nEL8JNKVwJZlH3HgU7eX7/NE5lM4RJ9E7Gf5bfbzR2qVRhHyVdxn7pExKh8tmgrQL\nbjVySjMTx+mbMszkr047DNKr/YVfUP01RJ0bDl+j3K1f+SDJJUlgp+oHVnObWrkB\nA473pWe5usW3hfYIaQDLa85QEXMj0CNcCvCD2Qe5w2VBaYQAPxtRNw2BaexBqHA7\nEq4DjkWSo+U0yASnRNs5Vem1UkmCplbCWDrTGpXd4uTmlxOEWLEfuDfmy5FRvXWy\nw9m7bLYVUl/JG2NAhvCh3kcuZn634vaWZndiLHNYxGzV37w3fJC/kztBUVS8fKI9\nX5ApB/m3u71HxhiAp14fqE3N4ZJdb2aiPhzMtS4Q9ELQmVKsauiMDgJ2RGPNVx0m\njVHUytMGqE64pLUnJbAq8icHEyY=\n-----END CERTIFICATE-----\n",
		MatchboxClientKey:  "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA0TxaMnZbC/CyZWe93cbrK5ajyaAA753/srMNQusxXRpBEHzj\nkOo1MPAwo7FPu/08EvAYMYzTbLZdHui5xl/sWiO7V3yU6Ve1XmkWG56fIGwT74s1\nRlnPZ7HEf/893wvf4aaxE9pMB6WaN3PuQShTDPTGk18bkDPY+eRzsunnYAUwY5nW\nAD/vhJgsCyALITtZJhBoU2yWMctlD9EKAX0uKnDn+3m3RAxUIcOaxVQEkyPSe3zF\nRWCPsphU4Xrbwmn28dcIkwW6Q/9Gmqs1vrHy+NtCNOJphaHADZn0OuH3p5919onZ\nolF5hKd2uNKBOWeZHf+x0T1ZekE54ORcPhTH5wIDAQABAoIBAQDQUPc4YI/7TlQ/\nE898G80OI7fESTJFXxgx8YSViJYoLLh00vK62APHSowSnblV4CVMcZCU3LGu/c2u\ntWQotl4ZzJN74gRcYb+oVJX3P7EMVa5vgziyZz5Q7jNGgHg4NplbD1wj+OZTgrQM\n27ZtVtmA/78pALrvkj9HZQUwLyL2e8bm4MP01cNOwff3HjYEOXxbJJDr9Jz8OPa8\nXmDuQKwuTRmohCD3FKnJrTeSH9X4DKi1CoMUB7kh6Pt1jCzPsSKKNAquwby0jHvF\nLJOw71l3fect80ankF356nKAv9KgpsGq5qQxPW4OSCvv7iTIcbh/Pyso12NN3GV1\naPf3PB1ZAoGBAPBrWZ5FUccKc89mme1HNV68hhrBGAJlSj8+lbSKjKpB/mTT1D/J\nqhSbuz3/ktxNWMPAx2NYi7iHDRzsUz7CKL55KcioZCCf4jE9enG+EGHNwxN2ac4/\nIguBJwB7qiANEAdJJqsQ7naYzrVHXWfXjXGvb7+5F2cBsoLGpSUkQVojAoGBAN7L\np42+wrumnfI6Iqs6FSrcOfRt3ys7cJFFXWpfv+n86i1Z1fJY7Bs0+nz4gqcbGZD9\nCodDJc7HgPMTZuOUg3yHo3EI5hz9FKKS15EjGCm5SDfs6TbqJSV7xtZfFr8fpD6w\n7tkNoq8CD26FZF3pF0QAoAe4YgmxV60UNAMqZe1tAoGBAN0GM/uXStkruNBhSP2k\ny2Hu+3K5NjNtn1aJWOQDw9H6nb9gJu8FnQEZMoiK3x79VK+SGTwx+TGJpvqCIP2/\nTeneRhWdCYAcvLv8Awdybmkb202XPSpJTCk7cPm2tu6EU8n+7De0dyY80TxDAZIn\nzndHi/q8VNFz9ALaUJTWweX1AoGAXXeGtXp/64V84a/t93OIidCWJ6soYtSu5uL4\ny7Wbp6hI/fmgPel8M/XH2EHRXhWKZj8h+Zj79YHQ4SkUkwktGEM3GCapkyPBUmrU\nMLlOW8K1P3EObdFRACarRifiPRAjMYG80iZcR5tPqgggER3GeurgOBzsVDCoHZ5K\nK8HPvQkCgYBc/+RjHK3HFhfcSWA3/m8/kNu0C41/vAlB8CYr11ojm9CuMgZp1/Sx\nQbD0YwMV1IEtKzTqHCWHkE9KlR2g13VGuIMvPQCVL4hBPwEnGUJRUOJfq7QkI2dm\niixDtZk3iFKLWGj5vjeH71YJsIQ8Z5+EOoOyHYnkhY8nqvhEz6Zmeg==\n-----END RSA PRIVATE KEY-----\n",
		ControllerDomain:   "cluster.example.com",
		TectonicDomain:     "tectonic.example.com",
		Controllers: []Node{
			Node{
				MAC:  &someMAC,
				Name: "node1.example.com",
			},
		},
		Workers: []Node{
			Node{
				MAC:  &someMAC,
				Name: "node2.example.com",
			},
			Node{
				MAC:  &someMAC,
				Name: "node3.example.com",
			},
		},
		SSHAuthorizedKeys: []string{"ssh-rsa pubkey"},
		Tectonic: &TectonicConfig{
			ControllerDomain: "cluster.example.com",
			TectonicDomain:   "cluster.example.com",
		},
		Version: "1081.5.0",
		Channel: "beta",
	}
}
