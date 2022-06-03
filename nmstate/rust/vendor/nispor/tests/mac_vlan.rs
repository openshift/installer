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

const IFACE_NAME: &str = "mac0";

const EXPECTED_MAC_VLAN_STATE: &str = r#"---
base_iface: eth1
mode: source
flags: 0
allowed_mac_addresses:
  - "00:23:45:67:89:1d"
  - "00:23:45:67:89:1c""#;

#[test]
fn test_get_macvlan_iface_yaml() {
    with_macvlan_iface(|| {
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_NAME];
        assert_eq!(iface.iface_type, nispor::IfaceType::MacVlan);
        assert_eq!(
            serde_yaml::to_string(&iface.mac_vlan).unwrap().trim(),
            EXPECTED_MAC_VLAN_STATE
        );
    });
}

fn with_macvlan_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("macvlan");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}
