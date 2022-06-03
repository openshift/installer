// Copyright 2021 Red Hat, Inc.
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

use nispor::NetState;
use pretty_assertions::assert_eq;
use serde_yaml;
use std::panic;

mod utils;

const IFACE_NAME0: &str = "sim0";
const IFACE_NAME1: &str = "sim1";

const EXPECTED_PAUSE_INFO: &str = r#"---
rx: true
tx: true
auto_negotiate: false"#;

const EXPECTED_FEATURE_INFO: &str = r#"---
fixed:
  esp-hw-offload: false
  esp-tx-csum-hw-offload: false
  fcoe-mtu: false
  highdma: true
  hw-tc-offload: false
  l2-fwd-offload: false
  loopback: true
  macsec-hw-offload: false
  netns-local: true
  rx-all: false
  rx-checksum: true
  rx-fcs: false
  rx-gro-hw: false
  rx-hashing: false
  rx-lro: false
  rx-ntuple-filter: false
  rx-udp_tunnel-port-offload: false
  rx-vlan-filter: false
  rx-vlan-hw-parse: false
  rx-vlan-stag-filter: false
  rx-vlan-stag-hw-parse: false
  tls-hw-record: false
  tls-hw-rx-offload: false
  tls-hw-tx-offload: false
  tx-checksum-fcoe-crc: false
  tx-checksum-ip-generic: true
  tx-checksum-ipv4: false
  tx-checksum-ipv6: false
  tx-checksum-sctp: true
  tx-esp-segmentation: false
  tx-fcoe-segmentation: false
  tx-gre-csum-segmentation: false
  tx-gre-segmentation: false
  tx-gso-partial: false
  tx-gso-robust: false
  tx-ipxip4-segmentation: false
  tx-ipxip6-segmentation: false
  tx-lockless: true
  tx-nocache-copy: false
  tx-scatter-gather-fraglist: true
  tx-tunnel-remcsum-segmentation: false
  tx-udp_tnl-csum-segmentation: false
  tx-udp_tnl-segmentation: false
  tx-vlan-hw-insert: false
  tx-vlan-stag-hw-insert: false
  vlan-challenged: true
changeable:
  rx-gro: true
  rx-gro-list: false
  tx-generic-segmentation: true
  tx-sctp-segmentation: true
  tx-tcp-ecn-segmentation: true
  tx-tcp-mangleid-segmentation: true
  tx-tcp-segmentation: true
  tx-tcp6-segmentation: true"#;

#[test]
#[ignore] // Only new version of netdevsim support pause
fn test_get_ethtool_pause_yaml() {
    with_netdevsim_iface(|| {
        let state = NetState::retrieve().unwrap();
        let iface0 = &state.ifaces[IFACE_NAME0];
        let iface1 = &state.ifaces[IFACE_NAME1];
        assert_eq!(&iface0.iface_type, &nispor::IfaceType::Ethernet);
        assert_eq!(&iface1.iface_type, &nispor::IfaceType::Ethernet);
        assert_eq!(
            serde_yaml::to_string(&iface0.ethtool.as_ref().unwrap().pause)
                .unwrap()
                .trim(),
            EXPECTED_PAUSE_INFO
        );
        assert_eq!(
            serde_yaml::to_string(&iface1.ethtool.as_ref().unwrap().pause)
                .unwrap()
                .trim(),
            EXPECTED_PAUSE_INFO
        );
    });
}

#[test]
fn test_get_ethtool_feature_yaml_of_loopback() {
    let mut state = NetState::retrieve().unwrap();
    let iface = state.ifaces.get_mut("lo").unwrap();
    // These property value is different between Github CI and my Archlinux
    iface
        .ethtool
        .as_mut()
        .unwrap()
        .features
        .as_mut()
        .map(|features| features.fixed.remove("tx-gso-list"));
    iface
        .ethtool
        .as_mut()
        .unwrap()
        .features
        .as_mut()
        .map(|features| features.changeable.remove("tx-gso-list"));
    iface
        .ethtool
        .as_mut()
        .unwrap()
        .features
        .as_mut()
        .map(|features| features.fixed.remove("tx-udp-segmentation"));
    iface
        .ethtool
        .as_mut()
        .unwrap()
        .features
        .as_mut()
        .map(|features| features.changeable.remove("tx-udp-segmentation"));
    assert_eq!(&iface.iface_type, &nispor::IfaceType::Loopback);
    assert_eq!(
        serde_yaml::to_string(&iface.ethtool.as_ref().unwrap().features)
            .unwrap()
            .trim(),
        EXPECTED_FEATURE_INFO
    );
}

// TODO: There is no way to test the ethtool ring.

fn with_netdevsim_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("sim");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}

const IFACE_TUN_NAME: &str = "tun1";
const EXPECTED_ETHTOOL_COALESCE: &str = r#"---
rx_max_frames: 60"#;
const EXPECTED_ETHTOOL_LINK_MODE: &str = r#"---
auto_negotiate: false
ours: []
speed: 10
duplex: full"#;

#[test]
fn test_get_ethtool_coalesce_yaml() {
    with_tun_iface(|| {
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_TUN_NAME];
        assert_eq!(
            serde_yaml::to_string(&iface.ethtool.as_ref().unwrap().coalesce)
                .unwrap()
                .trim(),
            EXPECTED_ETHTOOL_COALESCE
        );
    });
}

#[test]
fn test_get_ethtool_link_mode_yaml() {
    with_tun_iface(|| {
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_TUN_NAME];
        assert_eq!(
            serde_yaml::to_string(&iface.ethtool.as_ref().unwrap().link_mode)
                .unwrap()
                .trim(),
            EXPECTED_ETHTOOL_LINK_MODE
        );
    });
}

fn with_tun_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("tun");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}
