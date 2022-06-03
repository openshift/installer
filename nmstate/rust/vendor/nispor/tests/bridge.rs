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

const IFACE_NAME: &str = "br0";
const PORT1_NAME: &str = "eth1";
const PORT2_NAME: &str = "eth2";

const EXPECTED_BRIDGE_INFO: &str = r#"---
ports:
  - eth1
  - eth2
ageing_time: 30000
bridge_id: 8000.00234567891c
group_fwd_mask: 0
root_id: 8000.00234567891c
root_port: 0
root_path_cost: 0
topology_change: false
topology_change_detected: false
tcn_timer: 0
topology_change_timer: 0
group_addr: "01:80:c2:00:00:00"
nf_call_iptables: false
nf_call_ip6tables: false
nf_call_arptables: false
vlan_filtering: false
vlan_protocol: 802.1q
default_pvid: 1
vlan_stats_enabled: false
vlan_stats_per_host: false
stp_state: disabled
hello_time: 200
hello_timer: 0
forward_delay: 1500
max_age: 2000
priority: 32768
multicast_router: temp_query
multicast_snooping: true
multicast_query_use_ifaddr: false
multicast_querier: false
multicast_stats_enabled: false
multicast_hash_elasticity: 16
multicast_hash_max: 4096
multicast_last_member_count: 2
multicast_last_member_interval: 100
multicast_startup_query_count: 2
multicast_membership_interval: 26000
multicast_querier_interval: 25500
multicast_query_interval: 12500
multicast_query_response_interval: 1000
multicast_igmp_version: 2
multicast_mld_version: 1"#;

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
    is_egress_untagged: true"#;

#[test]
fn test_get_br_iface_yaml() {
    with_br_iface(|| {
        let mut state = NetState::retrieve().unwrap();
        let iface = state.ifaces.get_mut(IFACE_NAME).unwrap();
        if let Some(ref mut bridge_info) = iface.bridge {
            bridge_info.gc_timer = None;
            // Below value is not supported by RHEL 8 and Ubuntu CI
            bridge_info.multi_bool_opt = None;
            // Below value is different between CI and RHEL/CentOS 8
            // https://blog.grisge.info/posts/br_on_250hz_kernel/
            bridge_info.multicast_startup_query_interval = None;
        }
        let port1 = state.ifaces.get_mut(PORT1_NAME).unwrap();
        if let Some(ref mut port_info) = port1.bridge_port {
            port_info.forward_delay_timer = 0;
            // Below values are not supported by Github CI Ubuntu 20.04
            port_info.mrp_in_open = None;
        }
        let port2 = state.ifaces.get_mut(PORT2_NAME).unwrap();
        if let Some(ref mut port_info) = port2.bridge_port {
            port_info.forward_delay_timer = 0;
            port_info.mrp_in_open = None;
        }

        let iface = &state.ifaces[IFACE_NAME];
        let port1 = &state.ifaces[PORT1_NAME];
        let port2 = &state.ifaces[PORT2_NAME];
        assert_eq!(iface.iface_type, nispor::IfaceType::Bridge);
        assert_eq!(
            serde_yaml::to_string(&iface.bridge).unwrap().trim(),
            EXPECTED_BRIDGE_INFO,
        );
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

fn with_br_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("br");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}

const BRIDGE_CREATE_YML: &str = r#"---
ifaces:
  - name: br0
    type: bridge"#;

const BRIDGE_DELETE_YML: &str = r#"---
ifaces:
  - name: br0
    type: bridge
    state: absent"#;

#[test]
fn test_create_delete_bridge() {
    let net_conf: NetConf = serde_yaml::from_str(BRIDGE_CREATE_YML).unwrap();
    net_conf.apply().unwrap();
    let state = NetState::retrieve().unwrap();
    let iface = &state.ifaces[IFACE_NAME];
    assert_eq!(&iface.iface_type, &nispor::IfaceType::Bridge);

    let net_conf: NetConf = serde_yaml::from_str(BRIDGE_DELETE_YML).unwrap();
    net_conf.apply().unwrap();
    let state = NetState::retrieve().unwrap();
    assert_eq!(None, state.ifaces.get(IFACE_NAME));
}
