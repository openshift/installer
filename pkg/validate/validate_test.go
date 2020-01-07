package validate

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterName(t *testing.T) {
	maxSizeName := strings.Repeat("123456789.", 5) + "1234"

	cases := []struct {
		name        string
		clusterName string
		valid       bool
	}{
		{"empty", "", false},
		{"only whitespace", " ", false},
		{"single lowercase", "a", true},
		{"single uppercase", "A", false},
		{"contains whitespace", "abc D", false},
		{"single number", "1", true},
		{"single dot", ".", false},
		{"ends with dot", "a.", false},
		{"starts with dot", ".a", false},
		{"multiple labels", "a.a", true},
		{"starts with dash", "-a", false},
		{"ends with dash", "a-", false},
		{"label starts with dash", "a.-a", false},
		{"label ends with dash", "a-.a", false},
		{"invalid percent", "a%a", false},
		{"only non-ascii", "日本語", false},
		{"contains non-ascii", "a日本語a", false},
		{"max size", maxSizeName, true},
		{"too long", maxSizeName + "a", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ClusterName(tc.clusterName)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestSubnetCIDR(t *testing.T) {
	cases := []struct {
		cidr   string
		expErr string
	}{
		{"0.0.0.0/32", "address must be specified"},
		{"1.2.3.4/0", "invalid network address. got 1.2.3.4/0, expecting 0.0.0.0/0"},
		{"1.2.3.4/1", "invalid network address. got 1.2.3.4/1, expecting 0.0.0.0/1"},
		{"1.2.3.4/31", ""},
		{"1.2.3.4/32", ""},
		{"0:0:0:0:0:1:102:304/116", "invalid network address. got ::1:102:304/116, expecting ::1:102:0/116"},
		{"0:0:0:0:0:ffff:102:304/116", "invalid network address. got 1.2.3.4/20, expecting 1.2.0.0/20"},
		{"172.17.0.0/20", "overlaps with default Docker Bridge subnet (172.17.0.0/20)"},
		{"172.0.0.0/8", "overlaps with default Docker Bridge subnet (172.0.0.0/8)"},
		{"255.255.255.255/1", "invalid network address. got 255.255.255.255/1, expecting 128.0.0.0/1"},
		{"255.255.255.255/32", ""},
	}
	for _, tc := range cases {
		t.Run(tc.cidr, func(t *testing.T) {
			ip, cidr, err := net.ParseCIDR(tc.cidr)
			if err != nil {
				t.Fatalf("could not parse cidr: %v", err)
			}
			err = SubnetCIDR(&net.IPNet{IP: ip, Mask: cidr.Mask})
			if tc.expErr != "" {
				assert.EqualError(t, err, tc.expErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDomainName_AcceptingTrailingDot(t *testing.T) {
	cases := []struct {
		domain string
		valid  bool
	}{
		{"", false},
		{" ", false},
		{"a", true},
		{".", false},
		{"日本語", false},
		{"日本語.com", false},
		{"abc.日本語.com", false},
		{"a日本語a.com", false},
		{"abc", true},
		{"ABC", false},
		{"ABC123", false},
		{"ABC123.COM123", false},
		{"1", true},
		{"0.0", true},
		{"1.2.3.4", true},
		{"1.2.3.4.", true},
		{"abc.", true},
		{"abc.com", true},
		{"abc.com.", true},
		{"a.b.c.d.e.f", true},
		{".abc", false},
		{".abc.com", false},
		{".abc.com", false},
	}
	for _, tc := range cases {
		t.Run(tc.domain, func(t *testing.T) {
			err := DomainName(tc.domain, true)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDomainName_RejectingTrailingDot(t *testing.T) {
	cases := []struct {
		domain string
		valid  bool
	}{
		{"", false},
		{" ", false},
		{"a", true},
		{".", false},
		{"日本語", false},
		{"日本語.com", false},
		{"abc.日本語.com", false},
		{"a日本語a.com", false},
		{"abc", true},
		{"ABC", false},
		{"ABC123", false},
		{"ABC123.COM123", false},
		{"1", true},
		{"0.0", true},
		{"1.2.3.4", true},
		{"1.2.3.4.", false},
		{"abc.", false},
		{"abc.com", true},
		{"abc.com.", false},
		{"a.b.c.d.e.f", true},
		{".abc", false},
		{".abc.com", false},
		{".abc.com", false},
	}
	for _, tc := range cases {
		t.Run(tc.domain, func(t *testing.T) {
			err := DomainName(tc.domain, false)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestNoProxyDomainName(t *testing.T) {
	cases := []struct {
		domain string
		valid  bool
	}{
		{"", false},
		{" ", false},
		{"a", true},
		{".", false},
		{"日本語", false},
		{"日本語.com", false},
		{"abc.日本語.com", false},
		{"a日本語a.com", false},
		{"abc", true},
		{"ABC", false},
		{"ABC123", false},
		{"ABC123.COM123", false},
		{"1", true},
		{"0.0", true},
		{"1.2.3.4", true},
		{"1.2.3.4.", true},
		{"abc.", true},
		{"abc.com", true},
		{"abc.com.", true},
		{"a.b.c.d.e.f", true},
		{".abc", true},
		{".abc.com", true},
	}
	for _, tc := range cases {
		t.Run(tc.domain, func(t *testing.T) {
			err := NoProxyDomainName(tc.domain)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDoCIDRsOverlap(t *testing.T) {
	cases := []struct {
		a       string
		b       string
		overlap bool
	}{
		{
			a:       "192.168.0.0/30",
			b:       "192.168.0.3/30",
			overlap: true,
		},
		{
			a:       "192.168.0.0/30",
			b:       "192.168.0.4/30",
			overlap: false,
		},
		{
			a:       "192.168.0.0/29",
			b:       "192.168.0.4/30",
			overlap: true,
		},
		{
			a:       "0.0.0.0/0",
			b:       "192.168.0.0/24",
			overlap: true,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %s", tc.a, tc.b), func(t *testing.T) {
			_, a, err := net.ParseCIDR(tc.a)
			if err != nil {
				t.Fatalf("could not parse cidr %q: %v", tc.a, err)
			}
			_, b, err := net.ParseCIDR(tc.b)
			if err != nil {
				t.Fatalf("could not parse cidr %q: %v", tc.b, err)
			}
			actual := DoCIDRsOverlap(a, b)
			assert.Equal(t, tc.overlap, actual)
		})
	}
}

func TestImagePullSecret(t *testing.T) {
	cases := []struct {
		name   string
		secret string
		valid  bool
	}{
		{
			name:   "single entry with auth",
			secret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
			valid:  true,
		},
		{
			name:   "single entry with credsStore",
			secret: `{"auths":{"example.com":{"credsStore":"creds store value"}}}`,
			valid:  true,
		},
		{
			name:   "empty",
			secret: `{}`,
			valid:  false,
		},
		{
			name:   "no auths",
			secret: `{"not-auths":{"example.com":{"auth":"authorization value"}}}`,
			valid:  false,
		},
		{
			name:   "no auth or credsStore",
			secret: `{"auths":{"example.com":{"unrequired-field":"value"}}}`,
			valid:  false,
		},
		{
			name:   "additional fields",
			secret: `{"auths":{"example.com":{"auth":"authorization value","other-field":"other field value"}}}`,
			valid:  true,
		},
		{
			name:   "no entries",
			secret: `{"auths":{}}`,
			valid:  false,
		},
		{
			name:   "multiple valid entries",
			secret: `{"auths":{"example.com":{"auth":"authorization value"},"other-example.com":{"auth":"other auth value"}}}`,
			valid:  true,
		},
		{
			name:   "mix of valid and invalid entries",
			secret: `{"auths":{"example.com":{"auth":"authorization value"},"other-example.com":{"unrequired-field":"value"}}}`,
			valid:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ImagePullSecret(tc.secret)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

const invalidFormatCertificate = `-----INVALID FORMAT-----
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
rLWYdt40z0pjcR3FDVzh
-----INVALID FORMAT-----
`
const validCACertificate = `-----BEGIN CERTIFICATE-----
MIIF2zCCA8OgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwgYExCzAJBgNVBAYTAlVT
MRcwFQYDVQQIDA5Ob3J0aCBDYXJvbGluYTEQMA4GA1UEBwwHUmFsZWlnaDEUMBIG
A1UECgwLUmVkIEhhdCBJbmMxHzAdBgNVBAsMFk9wZW5TaGlmdCBJbnN0YWxsIFRl
c3QxEDAOBgNVBAMMB1Jvb3QgQ0EwHhcNMTkwNzIyMjAwNzUxWhcNMjkwNzE5MjAw
NzUxWjB3MQswCQYDVQQGEwJVUzEXMBUGA1UECAwOTm9ydGggQ2Fyb2xpbmExFDAS
BgNVBAoMC1JlZCBIYXQgSW5jMR8wHQYDVQQLDBZPcGVuU2hpZnQgSW5zdGFsbCBU
ZXN0MRgwFgYDVQQDDA9JbnRlcm1lZGlhdGUgQ0EwggIiMA0GCSqGSIb3DQEBAQUA
A4ICDwAwggIKAoICAQDZhc69vEq9XyG+vcOW4rPx9aYJgn7NFXaE88xrKajFyu2v
kD5Mz7geQV/RQKp1RMvj/1JCW5Npw8QwoPXNGQ8M+d+ajGgSkUZNVBQRXiR/hpfK
ohox9gJRsOVCAvhyE15iZHkEVFFcchiWbsTM9QllLsiiI0qZ/QpkUmJmDyXUV4Hq
hoAGXsojp0xaEQhrl+Hayiwao7qZkbKFCbNIDFU++ZDNT41qqDwcYmbkBJgYoGdS
IAk4Mjf7+rLJPXWNYtYB3g1cuN4pH8FkFT9zocNr0xrsx2itY4gvXgIe/vzts8aw
sHx1h2HcZK7iJEHs25QGrsZhiADeb0i5pN1kaPqpY0qgQUCIaqZAtMMeHXQ0k3PB
xTz8vk0388oFLaJFuI0P9Q6CRf5+4rc9O201aUIuue3Y4IS6zAcd8yL5d5vxvCiN
Dbl7YenBS4C9xSEEiVZwN7AtIdKFq5pGrlptmhVbGFW1CLQNsVWpetCY12Sh9FOq
2IBaAup+XgRgO4kHs3t7euVaS2viH3MplPsOUim8NZPZBdZkTtS3W9SynBDriy1d
KtrYgz0zrgEAa82mq4INaR+7Utct97zhKa1zM47KlHgkauiTPkUcqVhoNWxdM5tI
nSWym/9pPHUmzt8v/F8COA/8Xv+db2QX14S3fStI+8mp084RWuevtbh5WcoypQID
AQABo2YwZDAdBgNVHQ4EFgQUPUqJPYDZeUXbBlR0xXA/F+DYYagwHwYDVR0jBBgw
FoAUjWflPh3KYZ5o3BP3Po4v2ZBshVkwEgYDVR0TAQH/BAgwBgEB/wIBADAOBgNV
HQ8BAf8EBAMCAYYwDQYJKoZIhvcNAQELBQADggIBAH665ntrBhyf+MPFnkY+1VUr
VrfRlP4SccoujdLB/sUKqydYsED+mDJ+V8uFOgoi7PHqwvsRS+yR/bB0bNNYSfKY
slCMQA3sJ7SNDPBsec955ehYPNdquhem+oICzgFaQwL9ULDG87fKZjmaKO25dIYX
ttLqn+0b0GjpfQRuZ3NpAnCTWevodc5A3aYQm6vYeCyeIHGPpmtLE6oPRFib7wtD
n4DFVM57F34ClnnF4m8jq9HoTcM1Y3qOFyslK/4FRyx3HXbEVsm5L289l0AS866U
WEVM9DCqpFNLTwRk0mn4mspNcRxTDUTiHAxMhKxHGgbPcFzCJXqZzkW56bDcAGA5
sQr+MOfa1P/K7pVcFtOAhsBi5ff1G4t1G1+amqXEDalL+qKRGFugGVf+poyb2C3g
sfxkPBp9jPPMgMzXULQglwU4IUm8GtBb9Lh6AFPvt78XAWvNvHLP1Rf8JNZ9prx5
N9RzIKSWKm6CVEjSDvQ42j4OpW0eecHAoluZFMrykVl+KmapWUwQF6v0xz1RJdQ+
q3vGJ6shhiFd6y0ygxPwMaEjhhpbRy4tK9iDBj5yRpo+HE5X+FQSN6NHOYWMeDoZ
uzd86/huEH5qIAL4unM9YFTzJ4CFOC8EJMDW6ul0uKjOwGPP3R1Vss6sC7kR0gXI
rLWYdt40z0pjcR3FDVzh
-----END CERTIFICATE-----
`

const validCertificate = `-----BEGIN CERTIFICATE-----
MIIF0TCCA7mgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwdzELMAkGA1UEBhMCVVMx
FzAVBgNVBAgMDk5vcnRoIENhcm9saW5hMRQwEgYDVQQKDAtSZWQgSGF0IEluYzEf
MB0GA1UECwwWT3BlblNoaWZ0IEluc3RhbGwgVGVzdDEYMBYGA1UEAwwPSW50ZXJt
ZWRpYXRlIENBMB4XDTE5MDcyMjIwMTMwMloXDTIwMDczMTIwMTMwMlowgY4xCzAJ
BgNVBAYTAlVTMRcwFQYDVQQIDA5Ob3J0aCBDYXJvbGluYTEQMA4GA1UEBwwHUmFs
ZWlnaDEUMBIGA1UECgwLUmVkIEhhdCBJbmMxHzAdBgNVBAsMFk9wZW5TaGlmdCBJ
bnN0YWxsIFRlc3QxHTAbBgNVBAMMFHJlZ2lzdHJ5LmV4YW1wbGUuY29tMIIBIjAN
BgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyCT3n3zYL7PkLtnzBU9WUyZBz1Q+
SXUP739DjT+xmRunE1ViD2wfkVIhTHowlw6B7+23tSRQngEu5i4+lglzqouYY5jE
sqWUXaPMa5FeeDstI6LIUxqk9/2yWRBrrdJlVWor57F310aTzkYtkmkCJTDy3k9R
Le8jma8fnchaVpttbHgN/F+CiS+OV8u9PtALGuJ4DfHy2hM4pxhiKMxFpYOaxBuq
41Y0ts8CXyEiVWwZB7+fFrAjog8uJuhAdya1rAWvSDo+GQr2CDY2/PJcAVHO1n9F
h1LkjqIOd4OOqOy9gIYDc5bvZvWGuzeCN8icrdH8KM53witq5yhZHRL4EQIDAQAB
o4IBTTCCAUkwCQYDVR0TBAIwADARBglghkgBhvhCAQEEBAMCBkAwMwYJYIZIAYb4
QgENBCYWJE9wZW5TU0wgR2VuZXJhdGVkIFNlcnZlciBDZXJ0aWZpY2F0ZTAdBgNV
HQ4EFgQUxzJ/lMs83RIoGheipNbag+SZr4wwga8GA1UdIwSBpzCBpIAUPUqJPYDZ
eUXbBlR0xXA/F+DYYaihgYekgYQwgYExCzAJBgNVBAYTAlVTMRcwFQYDVQQIDA5O
b3J0aCBDYXJvbGluYTEQMA4GA1UEBwwHUmFsZWlnaDEUMBIGA1UECgwLUmVkIEhh
dCBJbmMxHzAdBgNVBAsMFk9wZW5TaGlmdCBJbnN0YWxsIFRlc3QxEDAOBgNVBAMM
B1Jvb3QgQ0GCAhAAMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcD
ATANBgkqhkiG9w0BAQsFAAOCAgEALXbZyIAXZPBBe+uWBscUFpmXvfQR9CQRbnYi
X4453dLKZqEsWNVcUQihrfaNyHCrkG7dSgT64JWWJcyXLXe9TfVR9FLGjzt4p0P2
V9T+PFjp3JN+Elh6XDeNisZ7fHzYs2yYnugZELdWkLOcUwkvUHhSQ5aSWYFrngn7
J3mT3GS3WSpLUvVQDn3RBDbS0ossnF1tq9n6Nhs4Xvhdso6gEZU9SeztbnSK9N/k
zWLV5PjgwpevJ17jzpxm7ZIAlcp31i4SIircJtGwgUS3cJZXPPWMdK72qnLQjFMF
BNEc11EBilMK8tn/K4Dn06BBJMRtOCkq0KhopeZX0HmtQE29z6hy9fcTBklCwLXQ
NMSOKemXsOiGTwghosa0xw0H2e9R8z9KTX5xBGgHbHbWu7e/oDVY9+XTnfQ0ZqFi
aBa/U/WWLMQQDNvQQGsllxBHC+pOpDD8YhycPmbpsFfhNo58U9VQ6mqtX3o7j5nP
imNTY4B5RmZUILe+C0XhON6VL5RCa+s6YngIUcfeylTSB8BTeVBxIAInubKzrgZM
4ThJWLbaiTkRqaT/viDfxsmgzJsrDm3ZWzYXwF/a5o6NHK4lqCYf/1nvbjgE5PAm
69R88P32rKeiRJ8AoC4N/5YR++NkB11gsW9ooU2nV90owi6eMhQ6+qGLTrq0mTtv
CNA1OOo=
-----END CERTIFICATE-----
`
const invalidStringWithCertificate = `-----BEGIN CERTIFICATE-----
MIIF0TCCA7mgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwdzELMAkGA1UEBhMCVVMx
FzAVBgNVBAgMDk5vcnRoIENhcm9saW5hMRQwEgYDVQQKDAtSZWQgSGF0IEluYzEf
MB0GA1UECwwWT3BlblNoaWZ0IEluc3RhbGwgVGVzdDEYMBYGA1UEAwwPSW50ZXJt
ZWRpYXRlIENBMB4XDTE5MDcyMjIwMTMwMloXDTIwMDczMTIwMTMwMlowgY4xCzAJ
BgNVBAYTAlVTMRcwFQYDVQQIDA5Ob3J0aCBDYXJvbGluYTEQMA4GA1UEBwwHUmFs
ZWlnaDEUMBIGA1UECgwLUmVkIEhhdCBJbmMxHzAdBgNVBAsMFk9wZW5TaGlmdCBJ
bnN0YWxsIFRlc3QxHTAbBgNVBAMMFHJlZ2lzdHJ5LmV4YW1wbGUuY29tMIIBIjAN
BgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyCT3n3zYL7PkLtnzBU9WUyZBz1Q+
SXUP739DjT+xmRunE1ViD2wfkVIhTHowlw6B7+23tSRQngEu5i4+lglzqouYY5jE
sqWUXaPMa5FeeDstI6LIUxqk9/2yWRBrrdJlVWor57F310aTzkYtkmkCJTDy3k9R
Le8jma8fnchaVpttbHgN/F+CiS+OV8u9PtALGuJ4DfHy2hM4pxhiKMxFpYOaxBuq
41Y0ts8CXyEiVWwZB7+fFrAjog8uJuhAdya1rAWvSDo+GQr2CDY2/PJcAVHO1n9F
h1LkjqIOd4OOqOy9gIYDc5bvZvWGuzeCN8icrdH8KM53witq5yhZHRL4EQIDAQAB
o4IBTTCCAUkwCQYDVR0TBAIwADARBglghkgBhvhCAQEEBAMCBkAwMwYJYIZIAYb4
QgENBCYWJE9wZW5TU0wgR2VuZXJhdGVkIFNlcnZlciBDZXJ0aWZpY2F0ZTAdBgNV
HQ4EFgQUxzJ/lMs83RIoGheipNbag+SZr4wwga8GA1UdIwSBpzCBpIAUPUqJPYDZ
eUXbBlR0xXA/F+DYYaihgYekgYQwgYExCzAJBgNVBAYTAlVTMRcwFQYDVQQIDA5O
b3J0aCBDYXJvbGluYTEQMA4GA1UEBwwHUmFsZWlnaDEUMBIGA1UECgwLUmVkIEhh
dCBJbmMxHzAdBgNVBAsMFk9wZW5TaGlmdCBJbnN0YWxsIFRlc3QxEDAOBgNVBAMM
B1Jvb3QgQ0GCAhAAMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcD
ATANBgkqhkiG9w0BAQsFAAOCAgEALXbZyIAXZPBBe+uWBscUFpmXvfQR9CQRbnYi
X4453dLKZqEsWNVcUQihrfaNyHCrkG7dSgT64JWWJcyXLXe9TfVR9FLGjzt4p0P2
V9T+PFjp3JN+Elh6XDeNisZ7fHzYs2yYnugZELdWkLOcUwkvUHhSQ5aSWYFrngn7
J3mT3GS3WSpLUvVQDn3RBDbS0ossnF1tq9n6Nhs4Xvhdso6gEZU9SeztbnSK9N/k
zWLV5PjgwpevJ17jzpxm7ZIAlcp31i4SIircJtGwgUS3cJZXPPWMdK72qnLQjFMF
BNEc11EBilMK8tn/K4Dn06BBJMRtOCkq0KhopeZX0HmtQE29z6hy9fcTBklCwLXQ
NMSOKemXsOiGTwghosa0xw0H2e9R8z9KTX5xBGgHbHbWu7e/oDVY9+XTnfQ0ZqFi
aBa/U/WWLMQQDNvQQGsllxBHC+pOpDD8YhycPmbpsFfhNo58U9VQ6mqtX3o7j5nP
imNTY4B5RmZUILe+C0XhON6VL5RCa+s6YngIUcfeylTSB8BTeVBxIAInubKzrgZM
4ThJWLbaiTkRqaT/viDfxsmgzJsrDm3ZWzYXwF/a5o6NHK4lqCYf/1nvbjgE5PAm
69R88P32rKeiRJ8AoC4N/5YR++NkB11gsW9ooU2nV90owi6eMhQ6+qGLTrq0mTtv
CNA1OOo=
-----END CERTIFICATE-----
Invalid data here
`

const invalidCertificate = `-----BEGIN CERTIFICATE-----
MIIF2zCCA8OgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwgYExCzAJBgNVBAYTAlVT
MRcwFQYDVQQIDA5Ob3J0aCBDYXJvbGluYTEQMA4GA1UEBwwHUmFsZWlnaDEUMBIG
-----END CERTIFICATE-----
`

func TestAdditionalTrustBundle(t *testing.T) {
	cases := []struct {
		name        string
		certificate string
		valid       bool
	}{
		{
			name:        "valid ca certificate",
			certificate: validCACertificate,
			valid:       true,
		},
		{
			name:        "valid certificate",
			certificate: validCertificate,
			valid:       true,
		},
		{
			name:        "invalid format",
			certificate: invalidFormatCertificate,
			valid:       false,
		},
		{
			name:        "invalid certificate",
			certificate: invalidFormatCertificate,
			valid:       false,
		},
		{
			name:        "valid certificate with extra invalid string",
			certificate: invalidStringWithCertificate,
			valid:       false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := CABundle(tc.certificate)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestSSHPublicKey(t *testing.T) {
	cases := []struct {
		name  string
		key   string
		valid bool
	}{
		{
			name:  "valid",
			key:   "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q==",
			valid: true,
		},
		{
			name:  "valid with email",
			key:   "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q== name@example.com",
			valid: true,
		},
		{
			name:  "invalid format",
			key:   "bad-format AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q==",
			valid: true,
		},
		{
			name:  "invalid key",
			key:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL",
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := SSHPublicKey(tc.key)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestURI(t *testing.T) {
	cases := []struct {
		name  string
		uri   string
		valid bool
	}{
		{
			name:  "valid",
			uri:   "https://example.com",
			valid: true,
		},
		{
			name:  "missing scheme",
			uri:   "example.com",
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := URI(tc.uri)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
