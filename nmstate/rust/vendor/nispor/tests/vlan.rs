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

use nispor::{NetConf, NetState};
use pretty_assertions::assert_eq;
use serde_yaml;
use std::panic;

mod utils;

const IFACE_NAME: &str = "eth1.101";

const EXPECTED_VLAN_INFO: &str = r#"---
vlan_id: 101
protocol: 802.1q
base_iface: eth1
is_reorder_hdr: true
is_gvrp: false
is_loose_binding: false
is_mvrp: false
is_bridge_binding: false"#;

#[test]
fn test_get_vlan_iface_yaml() {
    with_vlan_iface(|| {
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_NAME];
        assert_eq!(iface.iface_type, nispor::IfaceType::Vlan);
        assert_eq!(
            serde_yaml::to_string(&iface.vlan).unwrap().trim(),
            EXPECTED_VLAN_INFO
        );
    });
}

fn with_vlan_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("vlan");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}

const VETH_CREATE_YML: &str = r#"---
ifaces:
  - name: veth1
    type: veth
    veth:
      peer: veth1.ep
  - name: veth1.ep
    type: veth"#;

const VETH_DELETE_YML: &str = r#"---
ifaces:
  - name: veth1
    type: veth
    state: absent"#;

const VLAN_CREATE_YML: &str = r#"---
ifaces:
  - name: veth1.99
    type: vlan
    vlan:
      base_iface: veth1
      vlan_id: 99"#;

const VLAN_DELETE_YML: &str = r#"---
ifaces:
  - name: veth1.99
    type: vlan
    state: absent"#;

#[test]
fn test_create_delete_vlan() {
    let net_conf: NetConf = serde_yaml::from_str(VETH_CREATE_YML).unwrap();
    net_conf.apply().unwrap();

    let net_conf: NetConf = serde_yaml::from_str(VLAN_CREATE_YML).unwrap();
    net_conf.apply().unwrap();
    let state = NetState::retrieve().unwrap();
    let iface = &state.ifaces["veth1.99"];
    assert_eq!(&iface.iface_type, &nispor::IfaceType::Vlan);
    assert_eq!(iface.vlan.as_ref().unwrap().vlan_id, 99);
    assert_eq!(iface.vlan.as_ref().unwrap().base_iface.as_str(), "veth1");

    let net_conf: NetConf = serde_yaml::from_str(VLAN_DELETE_YML).unwrap();
    net_conf.apply().unwrap();
    let state = NetState::retrieve().unwrap();
    assert_eq!(None, state.ifaces.get("veth1.99"));

    let net_conf: NetConf = serde_yaml::from_str(VETH_DELETE_YML).unwrap();
    net_conf.apply().unwrap();
}
