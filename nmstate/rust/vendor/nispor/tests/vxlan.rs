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

const IFACE_NAME: &str = "vxlan0";

const EXPECTED_VXLAN_INFO: &str = r#"---
remote: 8.8.8.8
vxlan_id: 101
base_iface: eth1
local: 1.1.1.1
ttl: 0
tos: 0
learning: true
ageing: 300
max_address: 0
src_port_min: 0
src_port_max: 0
proxy: false
rsc: false
l2miss: false
l3miss: false
dst_port: 4789
udp_check_sum: true
udp6_zero_check_sum_tx: false
udp6_zero_check_sum_rx: false
remote_check_sum_tx: false
remote_check_sum_rx: false
gbp: false
remote_check_sum_no_partial: false
collect_metadata: false
label: 0
gpe: false
ttl_inherit: false
df: 0"#;

#[test]
fn test_get_vxlan_iface_yaml() {
    with_vxlan_iface(|| {
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_NAME];
        assert_eq!(iface.iface_type, nispor::IfaceType::Vxlan);
        assert_eq!(
            serde_yaml::to_string(&iface.vxlan).unwrap().trim(),
            EXPECTED_VXLAN_INFO
        );
    });
}

fn with_vxlan_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("vxlan");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}
