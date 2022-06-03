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

const IFACE_NAME: &str = "br0";
const PORT1_NAME: &str = "eth1";
const PORT2_NAME: &str = "eth2";

const EXPECTED_PORT1_BRIDGE_INFO: &str = r#"---
stp_state: forwarding
stp_priority: 32
stp_path_cost: 2
hairpin_mode: false
bpdu_guard: false
root_block: false
multicast_fast_leave: false
learning: true
unicast_flood: true
proxyarp: false
proxyarp_wifi: false
designated_root: 8000.00234567891c
designated_bridge: 8000.00234567891c
designated_port: 32769
designated_cost: 0
port_id: "0x8001"
port_no: "0x1"
change_ack: false
config_pending: false
message_age_timer: 0
forward_delay_timer: 0
hold_timer: 0
multicast_router: temp_query
multicast_flood: true
multicast_to_unicast: false
vlan_tunnel: false
broadcast_flood: true
group_fwd_mask: 0
neigh_suppress: false
isolated: false
mrp_ring_open: false
vlans:
  - vid: 1
    is_pvid: false
    is_egress_untagged: true
  - vid: 10
    is_pvid: true
    is_egress_untagged: true"#;

const EXPECTED_PORT2_BRIDGE_INFO: &str = r#"---
stp_state: forwarding
stp_priority: 32
stp_path_cost: 2
hairpin_mode: false
bpdu_guard: false
root_block: false
multicast_fast_leave: false
learning: true
unicast_flood: true
proxyarp: false
proxyarp_wifi: false
designated_root: 8000.00234567891c
designated_bridge: 8000.00234567891c
designated_port: 32770
designated_cost: 0
port_id: "0x8002"
port_no: "0x2"
change_ack: false
config_pending: false
message_age_timer: 0
forward_delay_timer: 0
hold_timer: 0
multicast_router: temp_query
multicast_flood: true
multicast_to_unicast: false
vlan_tunnel: false
broadcast_flood: true
group_fwd_mask: 0
neigh_suppress: false
isolated: false
mrp_ring_open: false
vlans:
  - vid: 1
    is_pvid: true
    is_egress_untagged: true
  - vid_range:
      - 2
      - 4094
    is_pvid: false
    is_egress_untagged: false"#;

#[test]
fn test_get_br_vlan_filter_iface_yaml() {
    with_br_with_vlan_filter_iface(|| {
        let mut state = NetState::retrieve().unwrap();
        let port1 = state.ifaces.get_mut(PORT1_NAME).unwrap();
        if let Some(ref mut port_info) = port1.bridge_port {
            port_info.forward_delay_timer = 0;
            // Below values are not supported by Github CI Ubuntu 20.04
            port_info.mrp_in_open = None;
        }
        let port2 = state.ifaces.get_mut(PORT2_NAME).unwrap();
        if let Some(ref mut port_info) = port2.bridge_port {
            port_info.forward_delay_timer = 0;
            // Below values are not supported by Github CI Ubuntu 20.04
            port_info.mrp_in_open = None;
        }
        let iface = state.ifaces.get(IFACE_NAME).unwrap();
        if let Some(bridge_info) = &iface.bridge {
            assert_eq!(bridge_info.vlan_filtering, Some(true))
        }

        let port1 = &state.ifaces[PORT1_NAME];
        let port2 = &state.ifaces[PORT2_NAME];
        assert_eq!(
            serde_yaml::to_string(&port1.bridge_port).unwrap().trim(),
            EXPECTED_PORT1_BRIDGE_INFO,
        );
        assert_eq!(
            serde_yaml::to_string(&port2.bridge_port).unwrap().trim(),
            EXPECTED_PORT2_BRIDGE_INFO,
        );
    });
}

fn with_br_with_vlan_filter_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("brv");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}
