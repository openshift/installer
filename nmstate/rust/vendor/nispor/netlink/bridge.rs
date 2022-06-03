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

use crate::mac::parse_as_mac;
use crate::BridgeInfo;
use crate::NisporError;
use netlink_packet_route::rtnl::link::nlas::InfoBridge;

const ETH_ALEN: usize = 6;

pub(crate) fn parse_bridge_info(
    infos: &[InfoBridge],
) -> Result<BridgeInfo, NisporError> {
    let mut bridge_info = BridgeInfo::default();

    for info in infos {
        if let InfoBridge::ForwardDelay(d) = info {
            bridge_info.forward_delay = Some(*d);
        } else if let InfoBridge::HelloTime(d) = info {
            bridge_info.hello_time = Some(*d);
        } else if let InfoBridge::MaxAge(d) = info {
            bridge_info.max_age = Some(*d);
        } else if let InfoBridge::AgeingTime(d) = info {
            bridge_info.ageing_time = Some(*d);
        } else if let InfoBridge::StpState(d) = info {
            bridge_info.stp_state = Some((*d).into());
        } else if let InfoBridge::Priority(d) = info {
            bridge_info.priority = Some(*d);
        } else if let InfoBridge::VlanFiltering(d) = info {
            bridge_info.vlan_filtering = Some(*d > 0);
        } else if let InfoBridge::VlanProtocol(d) = info {
            bridge_info.vlan_protocol = Some((*d).into());
        } else if let InfoBridge::GroupFwdMask(d) = info {
            bridge_info.group_fwd_mask = Some(*d);
        } else if let InfoBridge::RootId((priority, mac)) = info {
            bridge_info.root_id = Some(parse_bridge_id(*priority, mac)?);
        } else if let InfoBridge::BridgeId((priority, mac)) = info {
            bridge_info.bridge_id = Some(parse_bridge_id(*priority, mac)?);
        } else if let InfoBridge::RootPort(d) = info {
            bridge_info.root_port = Some(*d);
        } else if let InfoBridge::RootPathCost(d) = info {
            bridge_info.root_path_cost = Some(*d);
        } else if let InfoBridge::TopologyChange(d) = info {
            bridge_info.topology_change = Some(*d > 0);
        } else if let InfoBridge::TopologyChangeDetected(d) = info {
            bridge_info.topology_change_detected = Some(*d > 0);
        } else if let InfoBridge::HelloTimer(d) = info {
            bridge_info.hello_timer = Some(*d);
        } else if let InfoBridge::TcnTimer(d) = info {
            bridge_info.tcn_timer = Some(*d);
        } else if let InfoBridge::TopologyChangeTimer(d) = info {
            bridge_info.topology_change_timer = Some(*d);
        } else if let InfoBridge::GcTimer(d) = info {
            bridge_info.gc_timer = Some(*d);
        } else if let InfoBridge::GroupAddr(d) = info {
            bridge_info.group_addr = Some(parse_as_mac(ETH_ALEN, d)?);
        // InfoBridge::FdbFlush is only used for changing bridge
        } else if let InfoBridge::MulticastRouter(d) = info {
            bridge_info.multicast_router = Some((*d).into());
        } else if let InfoBridge::MulticastSnooping(d) = info {
            bridge_info.multicast_snooping = Some((*d) > 0);
        } else if let InfoBridge::MulticastQueryUseIfaddr(d) = info {
            bridge_info.multicast_query_use_ifaddr = Some((*d) > 0);
        } else if let InfoBridge::MulticastQuerier(d) = info {
            bridge_info.multicast_querier = Some((*d) > 0);
        } else if let InfoBridge::MulticastHashElasticity(d) = info {
            bridge_info.multicast_hash_elasticity = Some(*d);
        } else if let InfoBridge::MulticastHashMax(d) = info {
            bridge_info.multicast_hash_max = Some(*d);
        } else if let InfoBridge::MulticastLastMemberCount(d) = info {
            bridge_info.multicast_last_member_count = Some(*d);
        } else if let InfoBridge::MulticastStartupQueryCount(d) = info {
            bridge_info.multicast_startup_query_count = Some(*d);
        } else if let InfoBridge::MulticastLastMemberInterval(d) = info {
            bridge_info.multicast_last_member_interval = Some(*d);
        } else if let InfoBridge::MulticastMembershipInterval(d) = info {
            bridge_info.multicast_membership_interval = Some(*d);
        } else if let InfoBridge::MulticastQuerierInterval(d) = info {
            bridge_info.multicast_querier_interval = Some(*d);
        } else if let InfoBridge::MulticastQueryInterval(d) = info {
            bridge_info.multicast_query_interval = Some(*d);
        } else if let InfoBridge::MulticastQueryResponseInterval(d) = info {
            bridge_info.multicast_query_response_interval = Some(*d);
        } else if let InfoBridge::MulticastStartupQueryInterval(d) = info {
            bridge_info.multicast_startup_query_interval = Some(*d);
        } else if let InfoBridge::NfCallIpTables(d) = info {
            bridge_info.nf_call_iptables = Some(*d > 0);
        } else if let InfoBridge::NfCallIp6Tables(d) = info {
            bridge_info.nf_call_ip6tables = Some(*d > 0);
        } else if let InfoBridge::NfCallArpTables(d) = info {
            bridge_info.nf_call_arptables = Some(*d > 0);
        } else if let InfoBridge::VlanDefaultPvid(d) = info {
            bridge_info.default_pvid = Some(*d);
        } else if let InfoBridge::VlanStatsEnabled(d) = info {
            bridge_info.vlan_stats_enabled = Some(*d > 0);
        } else if let InfoBridge::MulticastStatsEnabled(d) = info {
            bridge_info.multicast_stats_enabled = Some(*d > 0);
        } else if let InfoBridge::MulticastIgmpVersion(d) = info {
            bridge_info.multicast_igmp_version = Some(*d);
        } else if let InfoBridge::MulticastMldVersion(d) = info {
            bridge_info.multicast_mld_version = Some(*d);
        } else if let InfoBridge::VlanStatsPerHost(d) = info {
            bridge_info.vlan_stats_per_host = Some(*d > 0);
        } else if let InfoBridge::MultiBoolOpt(d) = info {
            bridge_info.multi_bool_opt = Some(*d);
        } else {
            log::warn!("Unknown NLA {:?}", &info);
        }
    }
    Ok(bridge_info)
}

fn parse_bridge_id(
    priority: u16,
    mac: &[u8; 6],
) -> Result<String, NisporError> {
    //Following the format of sysfs
    let priority_bytes = priority.to_ne_bytes();
    let mac = parse_as_mac(ETH_ALEN, mac)
        .map_err(|_| {
            NisporError::invalid_argument(
                "invalid mac address in bridge_id".into(),
            )
        })?
        .to_lowercase()
        .replace(':', "");

    Ok(format!(
        "{:02x}{:02x}.{}",
        priority_bytes[0], priority_bytes[1], mac
    ))
}
