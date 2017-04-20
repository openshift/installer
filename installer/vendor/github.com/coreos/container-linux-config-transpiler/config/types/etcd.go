// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"errors"
	"fmt"

	"github.com/coreos/go-semver/semver"
	ignTypes "github.com/coreos/ignition/config/v2_0/types"
	"github.com/coreos/ignition/config/validate/report"
)

var (
	EtcdVersionTooOld      = errors.New("Etcd version specified is not valid (too old)")
	EtcdMinorVersionTooNew = errors.New("Etcd minor version specified is too new, only options available in the previous minor version will be accepted")
	EtcdMajorVersionTooNew = errors.New("Etcd version is not valid (too new)")
	OldestEtcd             = *semver.New("2.3.0")
)

// Options can be the options for any Etcd version
type Options interface{}

type etcdCommon Etcd

type EtcdVersion semver.Version

func (e *EtcdVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	t := semver.Version(*e)
	if err := unmarshal(&t); err != nil {
		return err
	}
	*e = EtcdVersion(t)
	return nil
}

func (e EtcdVersion) Validate() report.Report {
	v := semver.Version(e)
	switch {
	case v.LessThan(OldestEtcd):
		return report.ReportFromError(EtcdVersionTooOld, report.EntryError)
	case v.Major == 2 && v.Minor > 3:
		fallthrough
	case v.Major == 3 && v.Minor > 1:
		return report.ReportFromError(EtcdMinorVersionTooNew, report.EntryWarning)
	case v.Major > 3:
		return report.ReportFromError(EtcdMajorVersionTooNew, report.EntryError)
	}
	return report.Report{}
}

func (e EtcdVersion) String() string {
	return semver.Version(e).String()
}

// Etcd is a stub for yaml unmarshalling that figures out which
// of the other Etcd structs to use and unmarshals to that. Options needs
// to be an embedded type so that the structure of the yaml tree matches the
// structure of the go config tree
type Etcd struct {
	Version EtcdVersion `yaml:"version"`
	Options
}

func (etcd *Etcd) UnmarshalYAML(unmarshal func(interface{}) error) error {
	t := etcdCommon(*etcd)
	if err := unmarshal(&t); err != nil {
		return err
	}
	*etcd = Etcd(t)

	version := semver.Version(etcd.Version)
	if version.Major == 2 && version.Minor >= 3 {
		o := Etcd2{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		etcd.Options = o
	} else if version.Major == 3 && version.Minor == 0 {
		o := Etcd3_0{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		etcd.Options = o
	} else if version.Major == 3 && version.Minor >= 1 {
		o := Etcd3_1{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		etcd.Options = o
	}
	return nil
}

func init() {
	register2_0(func(in Config, out ignTypes.Config, platform string) (ignTypes.Config, report.Report) {
		if in.Etcd != nil {
			contents, err := etcdContents(*in.Etcd, platform)
			if err != nil {
				return ignTypes.Config{}, report.ReportFromError(err, report.EntryError)
			}
			out.Systemd.Units = append(out.Systemd.Units, ignTypes.SystemdUnit{
				Name:   "etcd-member.service",
				Enable: true,
				DropIns: []ignTypes.SystemdUnitDropIn{{
					Name:     "20-clct-etcd-member.conf",
					Contents: contents,
				}},
			})
		}
		return out, report.Report{}
	})
}

// etcdContents creates the string containing the systemd drop in for etcd-member
func etcdContents(etcd Etcd, platform string) (string, error) {
	args := getCliArgs(etcd.Options)
	vars := []string{fmt.Sprintf("ETCD_IMAGE_TAG=v%s", etcd.Version)}

	return assembleUnit("/usr/lib/coreos/etcd-wrapper $ETCD_OPTS", args, vars, platform)
}

type Etcd3_0 struct {
	Name                     string `yaml:"name"                        cli:"name"`
	DataDir                  string `yaml:"data_dir"                    cli:"data-dir"`
	WalDir                   string `yaml:"wal_dir"                     cli:"wal-dir"`
	SnapshotCount            int    `yaml:"snapshot_count"              cli:"snapshot-count"`
	HeartbeatInterval        int    `yaml:"heartbeat_interval"          cli:"heartbeat-interval"`
	ElectionTimeout          int    `yaml:"election_timeout"            cli:"election-timeout"`
	ListenPeerUrls           string `yaml:"listen_peer_urls"            cli:"listen-peer-urls"`
	ListenClientUrls         string `yaml:"listen_client_urls"          cli:"listen-client-urls"`
	MaxSnapshots             int    `yaml:"max_snapshots"               cli:"max-snapshots"`
	MaxWals                  int    `yaml:"max_wals"                    cli:"max-wals"`
	Cors                     string `yaml:"cors"                        cli:"cors"`
	InitialAdvertisePeerUrls string `yaml:"initial_advertise_peer_urls" cli:"initial-advertise-peer-urls"`
	InitialCluster           string `yaml:"initial_cluster"             cli:"initial-cluster"`
	InitialClusterState      string `yaml:"initial_cluster_state"       cli:"initial-cluster-state"`
	InitialClusterToken      string `yaml:"initial_cluster_token"       cli:"initial-cluster-token"`
	AdvertiseClientUrls      string `yaml:"advertise_client_urls"       cli:"advertise-client-urls"`
	Discovery                string `yaml:"discovery"                   cli:"discovery"`
	DiscoverySrv             string `yaml:"discovery_srv"               cli:"discovery-srv"`
	DiscoveryFallback        string `yaml:"discovery_fallback"          cli:"discovery-fallback"`
	DiscoveryProxy           string `yaml:"discovery_proxy"             cli:"discovery-proxy"`
	StrictReconfigCheck      bool   `yaml:"strict_reconfig_check"       cli:"strict-reconfig-check"`
	AutoCompactionRetention  int    `yaml:"auto_compaction_retention"   cli:"auto-compaction-retention"`
	Proxy                    string `yaml:"proxy"                       cli:"proxy"`
	ProxyFailureWait         int    `yaml:"proxy_failure_wait"          cli:"proxy-failure-wait"`
	ProxyRefreshInterval     int    `yaml:"proxy_refresh_interval"      cli:"proxy-refresh-interval"`
	ProxyDialTimeout         int    `yaml:"proxy_dial_timeout"          cli:"proxy-dial-timeout"`
	ProxyWriteTimeout        int    `yaml:"proxy_write_timeout"         cli:"proxy-write-timeout"`
	ProxyReadTimeout         int    `yaml:"proxy_read_timeout"          cli:"proxy-read-timeout"`
	CaFile                   string `yaml:"ca_file"                     cli:"ca-file"                     deprecated:"ca_file obsoleted by trusted_ca_file and client_cert_auth"`
	CertFile                 string `yaml:"cert_file"                   cli:"cert-file"`
	KeyFile                  string `yaml:"key_file"                    cli:"key-file"`
	ClientCertAuth           bool   `yaml:"client_cert_auth"            cli:"client-cert-auth"`
	TrustedCaFile            string `yaml:"trusted_ca_file"             cli:"trusted-ca-file"`
	AutoTls                  bool   `yaml:"auto_tls"                    cli:"auto-tls"`
	PeerCaFile               string `yaml:"peer_ca_file"                cli:"peer-ca-file"                deprecated:"peer_ca_file obsoleted peer_trusted_ca_file and peer_client_cert_auth"`
	PeerCertFile             string `yaml:"peer_cert_file"              cli:"peer-cert-file"`
	PeerKeyFile              string `yaml:"peer_key_file"               cli:"peer-key-file"`
	PeerClientCertAuth       bool   `yaml:"peer_client_cert_auth"       cli:"peer-client-cert-auth"`
	PeerTrustedCaFile        string `yaml:"peer_trusted_ca_file"        cli:"peer-trusted-ca-file"`
	PeerAutoTls              bool   `yaml:"peer_auto_tls"               cli:"peer-auto-tls"`
	Debug                    bool   `yaml:"debug"                       cli:"debug"`
	LogPackageLevels         string `yaml:"log_package_levels"          cli:"log-package-levels"`
	ForceNewCluster          bool   `yaml:"force_new_cluster"           cli:"force-new-cluster"`
}

type Etcd3_1 struct {
	Name                     string `yaml:"name"                        cli:"name"`
	DataDir                  string `yaml:"data_dir"                    cli:"data-dir"`
	WalDir                   string `yaml:"wal_dir"                     cli:"wal-dir"`
	SnapshotCount            int    `yaml:"snapshot_count"              cli:"snapshot-count"`
	HeartbeatInterval        int    `yaml:"heartbeat_interval"          cli:"heartbeat-interval"`
	ElectionTimeout          int    `yaml:"election_timeout"            cli:"election-timeout"`
	ListenPeerUrls           string `yaml:"listen_peer_urls"            cli:"listen-peer-urls"`
	ListenClientUrls         string `yaml:"listen_client_urls"          cli:"listen-client-urls"`
	MaxSnapshots             int    `yaml:"max_snapshots"               cli:"max-snapshots"`
	MaxWals                  int    `yaml:"max_wals"                    cli:"max-wals"`
	Cors                     string `yaml:"cors"                        cli:"cors"`
	InitialAdvertisePeerUrls string `yaml:"initial_advertise_peer_urls" cli:"initial-advertise-peer-urls"`
	InitialCluster           string `yaml:"initial_cluster"             cli:"initial-cluster"`
	InitialClusterState      string `yaml:"initial_cluster_state"       cli:"initial-cluster-state"`
	InitialClusterToken      string `yaml:"initial_cluster_token"       cli:"initial-cluster-token"`
	AdvertiseClientUrls      string `yaml:"advertise_client_urls"       cli:"advertise-client-urls"`
	Discovery                string `yaml:"discovery"                   cli:"discovery"`
	DiscoverySrv             string `yaml:"discovery_srv"               cli:"discovery-srv"`
	DiscoveryFallback        string `yaml:"discovery_fallback"          cli:"discovery-fallback"`
	DiscoveryProxy           string `yaml:"discovery_proxy"             cli:"discovery-proxy"`
	StrictReconfigCheck      bool   `yaml:"strict_reconfig_check"       cli:"strict-reconfig-check"`
	AutoCompactionRetention  int    `yaml:"auto_compaction_retention"   cli:"auto-compaction-retention"`
	Proxy                    string `yaml:"proxy"                       cli:"proxy"`
	ProxyFailureWait         int    `yaml:"proxy_failure_wait"          cli:"proxy-failure-wait"`
	ProxyRefreshInterval     int    `yaml:"proxy_refresh_interval"      cli:"proxy-refresh-interval"`
	ProxyDialTimeout         int    `yaml:"proxy_dial_timeout"          cli:"proxy-dial-timeout"`
	ProxyWriteTimeout        int    `yaml:"proxy_write_timeout"         cli:"proxy-write-timeout"`
	ProxyReadTimeout         int    `yaml:"proxy_read_timeout"          cli:"proxy-read-timeout"`
	CaFile                   string `yaml:"ca_file"                     cli:"ca-file"                     deprecated:"ca_file obsoleted by trusted_ca_file and client_cert_auth"`
	CertFile                 string `yaml:"cert_file"                   cli:"cert-file"`
	KeyFile                  string `yaml:"key_file"                    cli:"key-file"`
	ClientCertAuth           bool   `yaml:"client_cert_auth"            cli:"client-cert-auth"`
	TrustedCaFile            string `yaml:"trusted_ca_file"             cli:"trusted-ca-file"`
	AutoTls                  bool   `yaml:"auto_tls"                    cli:"auto-tls"`
	PeerCaFile               string `yaml:"peer_ca_file"                cli:"peer-ca-file"                deprecated:"peer_ca_file obsoleted peer_trusted_ca_file and peer_client_cert_auth"`
	PeerCertFile             string `yaml:"peer_cert_file"              cli:"peer-cert-file"`
	PeerKeyFile              string `yaml:"peer_key_file"               cli:"peer-key-file"`
	PeerClientCertAuth       bool   `yaml:"peer_client_cert_auth"       cli:"peer-client-cert-auth"`
	PeerTrustedCaFile        string `yaml:"peer_trusted_ca_file"        cli:"peer-trusted-ca-file"`
	PeerAutoTls              bool   `yaml:"peer_auto_tls"               cli:"peer-auto-tls"`
	Debug                    bool   `yaml:"debug"                       cli:"debug"`
	LogPackageLevels         string `yaml:"log_package_levels"          cli:"log-package-levels"`
	ForceNewCluster          bool   `yaml:"force_new_cluster"           cli:"force-new-cluster"`
	Metrics                  string `yaml:"metrics"                     cli:"metrics"`
	LogOutput                string `yaml:"log_output"                  cli:"log-output"`
}

type Etcd2 struct {
	AdvertiseClientURLs      string `yaml:"advertise_client_urls"         cli:"advertise-client-urls"`
	CAFile                   string `yaml:"ca_file"                       cli:"ca-file"                     deprecated:"ca_file obsoleted by trusted_ca_file and client_cert_auth"`
	CertFile                 string `yaml:"cert_file"                     cli:"cert-file"`
	ClientCertAuth           bool   `yaml:"client_cert_auth"              cli:"client-cert-auth"`
	CorsOrigins              string `yaml:"cors"                          cli:"cors"`
	DataDir                  string `yaml:"data_dir"                      cli:"data-dir"`
	Debug                    bool   `yaml:"debug"                         cli:"debug"`
	Discovery                string `yaml:"discovery"                     cli:"discovery"`
	DiscoveryFallback        string `yaml:"discovery_fallback"            cli:"discovery-fallback"`
	DiscoverySRV             string `yaml:"discovery_srv"                 cli:"discovery-srv"`
	DiscoveryProxy           string `yaml:"discovery_proxy"               cli:"discovery-proxy"`
	ElectionTimeout          int    `yaml:"election_timeout"              cli:"election-timeout"`
	EnablePprof              bool   `yaml:"enable_pprof"                  cli:"enable-pprof"`
	ForceNewCluster          bool   `yaml:"force_new_cluster"             cli:"force-new-cluster"`
	HeartbeatInterval        int    `yaml:"heartbeat_interval"            cli:"heartbeat-interval"`
	InitialAdvertisePeerURLs string `yaml:"initial_advertise_peer_urls"   cli:"initial-advertise-peer-urls"`
	InitialCluster           string `yaml:"initial_cluster"               cli:"initial-cluster"`
	InitialClusterState      string `yaml:"initial_cluster_state"         cli:"initial-cluster-state"`
	InitialClusterToken      string `yaml:"initial_cluster_token"         cli:"initial-cluster-token"`
	KeyFile                  string `yaml:"key_file"                      cli:"key-file"`
	ListenClientURLs         string `yaml:"listen_client_urls"            cli:"listen-client-urls"`
	ListenPeerURLs           string `yaml:"listen_peer_urls"              cli:"listen-peer-urls"`
	LogPackageLevels         string `yaml:"log_package_levels"            cli:"log-package-levels"`
	MaxSnapshots             int    `yaml:"max_snapshots"                 cli:"max-snapshots"`
	MaxWALs                  int    `yaml:"max_wals"                      cli:"max-wals"`
	Name                     string `yaml:"name"                          cli:"name"`
	PeerCAFile               string `yaml:"peer_ca_file"                  cli:"peer-ca-file"                deprecated:"peer_ca_file obsoleted peer_trusted_ca_file and peer_client_cert_auth"`
	PeerCertFile             string `yaml:"peer_cert_file"                cli:"peer-cert-file"`
	PeerKeyFile              string `yaml:"peer_key_file"                 cli:"peer-key-file"`
	PeerClientCertAuth       bool   `yaml:"peer_client_cert_auth"         cli:"peer-client-cert-auth"`
	PeerTrustedCAFile        string `yaml:"peer_trusted_ca_file"          cli:"peer-trusted-ca-file"`
	Proxy                    string `yaml:"proxy"                         cli:"proxy"                       valid:"^(on|off|readonly)$"`
	ProxyDialTimeout         int    `yaml:"proxy_dial_timeout"            cli:"proxy-dial-timeout"`
	ProxyFailureWait         int    `yaml:"proxy_failure_wait"            cli:"proxy-failure-wait"`
	ProxyReadTimeout         int    `yaml:"proxy_read_timeout"            cli:"proxy-read-timeout"`
	ProxyRefreshInterval     int    `yaml:"proxy_refresh_interval"        cli:"proxy-refresh-interval"`
	ProxyWriteTimeout        int    `yaml:"proxy_write_timeout"           cli:"proxy-write-timeout"`
	SnapshotCount            int    `yaml:"snapshot_count"                cli:"snapshot-count"`
	StrictReconfigCheck      bool   `yaml:"strict_reconfig_check"         cli:"strict-reconfig-check"`
	TrustedCAFile            string `yaml:"trusted_ca_file"               cli:"trusted-ca-file"`
	WalDir                   string `yaml:"wal_dir"                       cli:"wal-dir"`
}
