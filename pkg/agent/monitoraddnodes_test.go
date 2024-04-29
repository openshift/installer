package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	certificatesv1 "k8s.io/api/certificates/v1"
)

func TestDecodedFirstCSRSubjectContainsHostname(t *testing.T) {
	firstCSRRequestForExtraworker0 := "-----BEGIN CERTIFICATE REQUEST-----\nMIH3MIGdAgEAMDsxFTATBgNVBAoTDHN5c3RlbTpub2RlczEiMCAGA1UEAxMZc3lz\ndGVtOm5vZGU6ZXh0cmF3b3JrZXItMDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBGaK3U+3X3lM6tdgjD2b/y7Kysws8xgFW1rNd/wvKEvXzP5+A1K1M38zJiAWqKXP\n5AL2IDklO4GaO7PcRDNPabigADAKBggqhkjOPQQDAgNJADBGAiEA7C33Nym0Go73\nCZY+XOmyqE/IhaBMSwign+fgbPX1ibkCIQDHIfF7QpZReF93IW0v864/yLoXKyXy\nTGygkuR4KtXTDw==\n-----END CERTIFICATE REQUEST-----\n"
	tests := []struct {
		name           string
		hostnames      []string
		request        string
		expectedResult bool
	}{
		{
			name:           "request contains hostname",
			hostnames:      []string{"extraworker-0"},
			request:        firstCSRRequestForExtraworker0,
			expectedResult: true,
		},
		{
			name:           "request contains hostname using FQDN",
			hostnames:      []string{"extraworker-0.ostest.test.metalkube.org"},
			request:        firstCSRRequestForExtraworker0,
			expectedResult: true,
		},
		{
			name:           "request contains hostname when multiple names are resolved",
			hostnames:      []string{"somename", "extraworker-0.ostest.test.metalkube.org"},
			request:        firstCSRRequestForExtraworker0,
			expectedResult: true,
		},
		{
			name:           "request does not contain hostname",
			hostnames:      []string{"extraworker-1"},
			request:        firstCSRRequestForExtraworker0,
			expectedResult: false,
		},
		{
			name:           "request is empty string",
			hostnames:      []string{"hostname-not-specified"},
			request:        "",
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containsHostname := containsHostname(decodedFirstCSRSubject([]byte(tt.request)), tt.hostnames)
			assert.Equal(t, tt.expectedResult, containsHostname)
		})
	}
}

func TestFilterCSRsMatchingHostnames(t *testing.T) {
	firstCSRRequestForExtraworker0 := "-----BEGIN CERTIFICATE REQUEST-----\nMIH3MIGdAgEAMDsxFTATBgNVBAoTDHN5c3RlbTpub2RlczEiMCAGA1UEAxMZc3lz\ndGVtOm5vZGU6ZXh0cmF3b3JrZXItMDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA\nBGaK3U+3X3lM6tdgjD2b/y7Kysws8xgFW1rNd/wvKEvXzP5+A1K1M38zJiAWqKXP\n5AL2IDklO4GaO7PcRDNPabigADAKBggqhkjOPQQDAgNJADBGAiEA7C33Nym0Go73\nCZY+XOmyqE/IhaBMSwign+fgbPX1ibkCIQDHIfF7QpZReF93IW0v864/yLoXKyXy\nTGygkuR4KtXTDw==\n-----END CERTIFICATE REQUEST-----\n"

	tests := []struct {
		name           string
		csrs           *certificatesv1.CertificateSigningRequestList
		hostnames      []string
		signerName     string
		expectedResult []certificatesv1.CertificateSigningRequest
	}{
		{
			name: "first CSR filtering",
			csrs: &certificatesv1.CertificateSigningRequestList{
				Items: []certificatesv1.CertificateSigningRequest{
					{
						// should match only this one
						Spec: certificatesv1.CertificateSigningRequestSpec{
							SignerName: firstCSRSignerName,
							Request:    []byte(firstCSRRequestForExtraworker0),
						},
					},
					{
						Spec: certificatesv1.CertificateSigningRequestSpec{
							SignerName: "other-request",
							Request:    []byte("other-request"),
						},
					},
				},
			},
			hostnames:  []string{"extraworker-0.ostest.test.metalkube.org"},
			signerName: "kubernetes.io/kube-apiserver-client-kubelet",
			expectedResult: []certificatesv1.CertificateSigningRequest{
				{
					Spec: certificatesv1.CertificateSigningRequestSpec{
						SignerName: "kubernetes.io/kube-apiserver-client-kubelet",
						Request:    []byte(firstCSRRequestForExtraworker0),
					},
				},
			},
		},
		{
			name: "second CSR filtering",
			csrs: &certificatesv1.CertificateSigningRequestList{
				Items: []certificatesv1.CertificateSigningRequest{
					{
						// should match only this one
						Spec: certificatesv1.CertificateSigningRequestSpec{
							SignerName: secondCSRSignerName,
							Username:   "system:node:extraworker-0",
							Request:    []byte("something"),
						},
					},
					{
						Spec: certificatesv1.CertificateSigningRequestSpec{
							SignerName: secondCSRSignerName,
							Username:   "system:node:extraworker-1",
							Request:    []byte("something"),
						},
					},
					{
						Spec: certificatesv1.CertificateSigningRequestSpec{
							SignerName: "other-request",
							Request:    []byte("other-request"),
						},
					},
				},
			},
			hostnames:  []string{"extraworker-0.ostest.test.metalkube.org"},
			signerName: secondCSRSignerName,
			expectedResult: []certificatesv1.CertificateSigningRequest{
				{
					Spec: certificatesv1.CertificateSigningRequestSpec{
						SignerName: "kubernetes.io/kubelet-serving",
						Username:   "system:node:extraworker-0",
						Request:    []byte("something"),
					},
				},
			},
		},
		{
			name: "no CSRs should not result in error",
			csrs: &certificatesv1.CertificateSigningRequestList{
				Items: []certificatesv1.CertificateSigningRequest{},
			},
			hostnames:      []string{"extraworker-0.ostest.test.metalkube.org"},
			signerName:     secondCSRSignerName,
			expectedResult: []certificatesv1.CertificateSigningRequest{},
		},
		{
			name: "no hostnames should not result in error",
			csrs: &certificatesv1.CertificateSigningRequestList{
				Items: []certificatesv1.CertificateSigningRequest{},
			},
			hostnames:      []string{},
			signerName:     secondCSRSignerName,
			expectedResult: []certificatesv1.CertificateSigningRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filteredCSRs := filterCSRsMatchingHostname(tt.signerName, tt.csrs, tt.hostnames)
			assert.Equal(t, tt.expectedResult, filteredCSRs)
		})
	}
}
