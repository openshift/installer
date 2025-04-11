package agent

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	certificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	timeout = 90 * time.Minute
)

const (
	firstCSRSignerName  = "kubernetes.io/kube-apiserver-client-kubelet"
	secondCSRSignerName = "kubernetes.io/kubelet-serving"
)

type addNodeStatusHistory struct {
	RestAPISeen            bool
	KubeletIsRunningOnNode bool
	FirstCSRSeen           bool
	SecondCSRSeen          bool
	NodeJoinedCluster      bool
	NodeIsReady            bool
}

type addNodeMonitor struct {
	nodeIPAddress string
	hostnames     []string
	cluster       *Cluster
	status        addNodeStatusHistory
	ch            chan logEntry
}

func newAddNodeMonitor(nodeIP string, cluster *Cluster, ch chan logEntry) (*addNodeMonitor, error) {
	parsedIPAddress := net.ParseIP(nodeIP)
	if parsedIPAddress == nil {
		return nil, fmt.Errorf("%s is not valid IP Address", nodeIP)
	}
	mon := addNodeMonitor{
		nodeIPAddress: parsedIPAddress.String(),
		cluster:       cluster,
		status: addNodeStatusHistory{
			RestAPISeen:            false,
			KubeletIsRunningOnNode: false,
			FirstCSRSeen:           false,
			SecondCSRSeen:          false,
			NodeJoinedCluster:      false,
			NodeIsReady:            false,
		},
		ch: ch,
	}
	hostnames, err := net.LookupAddr(nodeIP)
	if err != nil {
		mon.logStatus(Info, fmt.Sprintf("Cannot resolve IP address %v to a hostname. Skipping checks for pending CSRs.", nodeIP))
	} else {
		mon.hostnames = hostnames
	}
	return &mon, nil
}

// MonitorAddNodes display the progress of one or more nodes being
// added to a cluster. ipAddresses is an array of IP addresses to be
// monitored. clusters is an array of their corresponding initialized Cluster
// struct used to interact with the assisted-service and k8s APIs.
func MonitorAddNodes(ctx context.Context, clusters []*Cluster, ipAddresses []string) {
	waitContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var wg sync.WaitGroup
	ipChanMap := make(map[string]chan logEntry)

	for i, ipAddress := range ipAddresses {
		wg.Add(1)
		ch := make(chan logEntry, 100)
		ipChanMap[ipAddress] = ch
		go MonitorSingleNode(waitContext, clusters[i], ipAddress, &wg, ch)
	}

	wg.Add(1)
	go printLogs(&wg, ipChanMap)

	wg.Wait()
}

// MonitorSingleNode waits for the a node to be added to the cluster
// and reports its status until it becomes Ready.
func MonitorSingleNode(waitContext context.Context, cluster *Cluster, nodeIPAddress string, wg *sync.WaitGroup, ch chan logEntry) {
	defer wg.Done()
	defer close(ch)

	mon, err := newAddNodeMonitor(nodeIPAddress, cluster, ch)
	if err != nil {
		mon.logStatus(Error, fmt.Sprintf("Could not initialize node monitor: %v", err))
		return
	}

	for {
		if !mon.status.RestAPISeen &&
			mon.cluster.API.Rest.IsRestAPILive() {
			mon.status.RestAPISeen = true
			mon.logStatus(Info, "Assisted Service API is available")
		}

		if !mon.status.KubeletIsRunningOnNode &&
			mon.isKubeletRunningOnNode() {
			mon.status.KubeletIsRunningOnNode = true
			mon.logStatus(Info, "Kubelet is running")
		}

		if mon.status.KubeletIsRunningOnNode &&
			!mon.status.FirstCSRSeen &&
			mon.clusterHasFirstCSRPending() {
			mon.status.FirstCSRSeen = true
			mon.logStatus(Info, "First CSR Pending approval")
			mon.logCSRsPendingApproval(firstCSRSignerName)
		}

		if mon.status.KubeletIsRunningOnNode &&
			!mon.status.SecondCSRSeen &&
			mon.clusterHasSecondCSRPending() {
			mon.status.SecondCSRSeen = true
			mon.logStatus(Info, "Second CSR Pending approval")
			mon.logCSRsPendingApproval(secondCSRSignerName)
		}

		hasJoined, isReady, err := mon.nodeHasJoinedClusterAndIsReady()
		if err != nil {
			mon.logStatus(Debug, fmt.Sprintf("Node joined cluster and is ready check returned err: %v", err))
		}

		if !mon.status.NodeJoinedCluster && hasJoined {
			mon.status.NodeJoinedCluster = true
			mon.logStatus(Info, "Node joined cluster")
		}

		if !mon.status.NodeIsReady && isReady {
			mon.status.NodeIsReady = true
			if !mon.clusterHasSecondCSRPending() {
				mon.logStatus(Info, "Node is Ready")
			} else {
				// The node becomes Ready before second CSR is approved. Log Pending CSRs
				// so users are aware there are still some waiting their approval even
				// though the node status is Ready.
				mon.logStatus(Info, "Node is Ready but has CSRs pending approval. Until all CSRs are approved, the node may not be fully functional.")
				mon.logCSRsPendingApproval(secondCSRSignerName)
			}
			break
		}

		if mon.cluster.API.Rest.IsRestAPILive() {
			_, err = cluster.MonitorStatusFromAssistedService(mon.ch)
			if err != nil {
				mon.logStatus(Warning, fmt.Sprintf("Error fetching status from assisted-service for node %s: %s", nodeIPAddress, err))
			}
		}

		waitErr := waitContext.Err()
		if waitErr != nil {
			if errors.Is(waitErr, context.Canceled) {
				mon.logStatus(Info, fmt.Sprintf("Node monitoring cancelled: %v", waitErr))
				break
			}
			if errors.Is(waitErr, context.DeadlineExceeded) {
				mon.logStatus(Info, fmt.Sprintf("Node monitoring timed out after %v minutes", timeout))
				break
			}
		}

		time.Sleep(logInterval * time.Second)
	}

	// print out logs before channel is closed
	printChannelLogs(nodeIPAddress, ch)
}

func (mon *addNodeMonitor) logStatus(level, message string) {
	mon.ch <- logEntry{level: level, message: message}
}

func (mon *addNodeMonitor) nodeHasJoinedClusterAndIsReady() (bool, bool, error) {
	nodes, err := mon.cluster.API.Kube.ListNodes()
	if err != nil {
		mon.logStatus(Debug, fmt.Sprintf("Error getting node list %v", err))
		return false, false, nil
	}

	var joinedNode corev1.Node
	hasJoined := false
	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP {
				if address.Address == mon.nodeIPAddress {
					joinedNode = node
					hasJoined = true
				}
			}
		}
	}

	isReady := false
	if hasJoined {
		mon.logStatus(Debug, fmt.Sprintf("%s has joined cluster", joinedNode.Name))
		for _, cond := range joinedNode.Status.Conditions {
			if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
				isReady = true
			}
		}
		if isReady {
			mon.logStatus(Debug, fmt.Sprintf("%s is Ready", joinedNode.Name))
		} else {
			mon.logStatus(Debug, fmt.Sprintf("%s is not Ready", joinedNode.Name))
		}
	} else {
		mon.logStatus(Debug, "Node has not joined cluster")
	}

	return hasJoined, isReady, nil
}

func (mon *addNodeMonitor) logCSRsPendingApproval(signerName string) {
	csrsPendingApproval := mon.getCSRsPendingApproval(signerName)

	for _, csr := range csrsPendingApproval {
		mon.logStatus(Info, fmt.Sprintf("CSR %s with signerName %s and username %s is Pending and awaiting approval",
			csr.Name, csr.Spec.SignerName, csr.Spec.Username))
	}
}

func (mon *addNodeMonitor) clusterHasFirstCSRPending() bool {
	return len(mon.getCSRsPendingApproval(firstCSRSignerName)) > 0
}

func (mon *addNodeMonitor) clusterHasSecondCSRPending() bool {
	return len(mon.getCSRsPendingApproval(secondCSRSignerName)) > 0
}

func (mon *addNodeMonitor) getCSRsPendingApproval(signerName string) []certificatesv1.CertificateSigningRequest {
	if mon.hostnames == nil {
		return []certificatesv1.CertificateSigningRequest{}
	}

	csrs, err := mon.cluster.API.Kube.ListCSRs()
	if err != nil {
		mon.logStatus(Debug, fmt.Sprintf("Error calling listCSRs(): %v", err))
		mon.logStatus(Info, "Cannot retrieve CSRs from Kube API. Skipping checks for pending CSRs")
		return []certificatesv1.CertificateSigningRequest{}
	}

	return filterCSRsMatchingHostname(signerName, csrs, mon.hostnames, mon.ch)
}

func filterCSRsMatchingHostname(signerName string, csrs *certificatesv1.CertificateSigningRequestList, hostnames []string, ch chan logEntry) []certificatesv1.CertificateSigningRequest {
	matchedCSRs := []certificatesv1.CertificateSigningRequest{}
	for _, csr := range csrs.Items {
		if len(csr.Status.Conditions) > 0 {
			// CSR is not Pending and not awaiting approval
			continue
		}
		if signerName == firstCSRSignerName && csr.Spec.SignerName == firstCSRSignerName &&
			containsHostname(decodedFirstCSRSubject(csr.Spec.Request, ch), hostnames) {
			matchedCSRs = append(matchedCSRs, csr)
		}
		if signerName == secondCSRSignerName && csr.Spec.SignerName == secondCSRSignerName &&
			containsHostname(csr.Spec.Username, hostnames) {
			matchedCSRs = append(matchedCSRs, csr)
		}
	}
	return matchedCSRs
}

// containsHostname checks if the searchString contains one of the node's
// hostnames. Only the first element of the hostname is checked.
// For example if the hostname is "extraworker-0.ostest.test.metalkube.org",
// "extraworker-0" is used to check if it exists in the searchString.
func containsHostname(searchString string, hostnames []string) bool {
	for _, hostname := range hostnames {
		parts := strings.Split(hostname, ".")
		if strings.Contains(searchString, parts[0]) {
			return true
		}
	}
	return false
}

// isKubeletRunningOnNode checks if kubelet responds
// to http. Even if kubelet responds with error like
// TLS errors, kubelet is considered running.
func (mon *addNodeMonitor) isKubeletRunningOnNode() bool {
	url := fmt.Sprintf("https://%s:10250/metrics", mon.nodeIPAddress)
	// http get without authentication
	resp, err := http.Get(url) //nolint mon.nodeIPAddress is prevalidated to be IP address
	if err != nil {
		mon.logStatus(Debug, fmt.Sprintf("kubelet http err: %v", err))
		if strings.Contains(err.Error(), "remote error: tls: internal error") {
			// nodes being added will return this error
			return true
		}
		if strings.Contains(err.Error(), "tls: failed to verify certificate: x509: certificate signed by unknown authority") {
			// existing control plane nodes returns this error
			return true
		}
		if strings.Contains(err.Error(), "connect: no route to host") {
			return false
		}
	} else {
		mon.logStatus(Debug, fmt.Sprintf("kubelet http status code: %v", resp.StatusCode))
	}
	return false
}

// decodedFirstCSRSubject decodes the CSR.Spec.Request PEM block
// into readable output and returns the subject as string.
//
// Example of decoded request:
// Certificate Request:
// Data:
// Version: 1 (0x0)
// Subject: O = system:nodes, CN = system:node:extraworker-1
// Subject Public Key Info:
//
//	Public Key Algorithm: id-ecPublicKey
//		Public-Key: (256 bit)
//		pub:
//			*snip*
//		ASN1 OID: prime256v1
//		NIST CURVE: P-256
//
// Attributes:
//
//	a0:00
//
// Signature Algorithm: ecdsa-with-SHA256
//
//	*snip*
func decodedFirstCSRSubject(request []byte, ch chan logEntry) string {
	block, _ := pem.Decode(request)
	if block == nil {
		return ""
	}
	csrDER := block.Bytes
	decodedRequest, err := x509.ParseCertificateRequest(csrDER)
	if err != nil {
		if ch != nil {
			log(Warning, "error in x509.ParseCertificateRequest(csrDER)", nil, ch)
		}
		return ""
	}
	return decodedRequest.Subject.String()
}
