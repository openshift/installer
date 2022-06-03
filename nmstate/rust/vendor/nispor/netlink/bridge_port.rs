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

use crate::netlink::nla::parse_as_u16;
use crate::netlink::nla::parse_as_u32;
use crate::netlink::nla::parse_as_u64;
use crate::netlink::nla::parse_as_u8;
use crate::BridgePortInfo;
use crate::NisporError;
use netlink_packet_route::rtnl::nlas::NlasIterator;

fn parse_void_port_info(
    _data: &[u8],
    _port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    Ok(())
}

fn parse_brport_state(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.stp_state = parse_as_u8(data)?.into();
    Ok(())
}

fn parse_brport_priority(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.stp_priority = parse_as_u16(data)?;
    Ok(())
}

fn parse_brport_cost(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.stp_path_cost = parse_as_u32(data)?;
    Ok(())
}

fn parse_brport_mode(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.hairpin_mode = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_guard(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.bpdu_guard = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_protect(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.root_block = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_fast_leave(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.multicast_fast_leave = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_learning(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.learning = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_unicast_flood(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.unicast_flood = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_proxyarp(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.proxyarp = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_learning_sync(
    _data: &[u8],
    _port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    Ok(()) // Ther kernel 5.7-rc6 never update fill value in br_port_fill_attrs
}

fn parse_brport_proxyarp_wifi(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.proxyarp_wifi = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_root_id(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.designated_root = parse_as_bridge_id(data)?;
    Ok(())
}

fn parse_brport_bridge_id(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.designated_bridge = parse_as_bridge_id(data)?;
    Ok(())
}

fn parse_brport_designated_port(
    data: &[u8],
    port_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    port_info.designated_port = parse_as_u16(data)?;
    Ok(())
}

fn parse_brport_designated_cost(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.designated_cost = parse_as_u16(data)?;
    Ok(())
}

fn parse_brport_id(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.port_id = format!("0x{:04x}", parse_as_u16(data)?);
    Ok(())
}

fn parse_brport_no(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.port_no = format!("0x{:x}", parse_as_u16(data)?);
    Ok(())
}

fn parse_brport_topology_change_ack(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.change_ack = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_config_pending(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.config_pending = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_message_age_timer(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.message_age_timer = parse_as_u64(data)?;
    Ok(())
}

fn parse_brport_forward_delay_timer(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.forward_delay_timer = parse_as_u64(data)?;
    Ok(())
}

fn parse_brport_hold_timer(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.hold_timer = parse_as_u64(data)?;
    Ok(())
}

fn parse_brport_multicast_router(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.multicast_router = parse_as_u8(data)?.into();
    Ok(())
}

fn parse_brport_mcast_flood(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.multicast_flood = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_mcast_to_ucast(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.multicast_to_unicast = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_vlan_tunnel(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.vlan_tunnel = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_bast_flood(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.broadcast_flood = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_group_fwd_mask(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.group_fwd_mask = parse_as_u16(data)?;
    Ok(())
}

fn parse_brport_neigh_suppress(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.neigh_suppress = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_isolated(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.isolated = parse_as_u8(data)? > 0;
    Ok(())
}

fn parse_brport_backup_port(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.backup_port = format!("{}", parse_as_u32(data)?);
    Ok(())
}

fn parse_brport_mrp_ring_open(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.mrp_ring_open = Some(parse_as_u8(data)? > 0);
    Ok(())
}

fn parse_brport_mrp_in_open(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.mrp_in_open = Some(parse_as_u8(data)? > 0);
    Ok(())
}

fn parse_brport_mcast_eht_hosts_limit(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.mcast_eht_hosts_limit = Some(parse_as_u32(data)?);
    Ok(())
}

fn parse_brport_mcast_eht_hosts_cnt(
    data: &[u8],
    cost_info: &mut BridgePortInfo,
) -> Result<(), NisporError> {
    cost_info.mcast_eht_hosts_cnt = Some(parse_as_u32(data)?);
    Ok(())
}

type BridgePortParseFunc =
    fn(&[u8], &mut BridgePortInfo) -> Result<(), NisporError>;

const NLA_PORT_PARSE_FUNS: &[BridgePortParseFunc] = &[
    parse_void_port_info, // IFLA_BRPORT_UNSPEC
    parse_brport_state,
    parse_brport_priority,
    parse_brport_cost,
    parse_brport_mode,
    parse_brport_guard,
    parse_brport_protect,
    parse_brport_fast_leave,
    parse_brport_learning,
    parse_brport_unicast_flood,
    parse_brport_proxyarp,
    parse_brport_learning_sync,
    parse_brport_proxyarp_wifi,
    parse_brport_root_id,
    parse_brport_bridge_id,
    parse_brport_designated_port,
    parse_brport_designated_cost,
    parse_brport_id,
    parse_brport_no,
    parse_brport_topology_change_ack,
    parse_brport_config_pending,
    parse_brport_message_age_timer,
    parse_brport_forward_delay_timer,
    parse_brport_hold_timer,
    parse_void_port_info, // IFLA_BRPORT_FLUSH
    parse_brport_multicast_router,
    parse_void_port_info, // IFLA_BRPORT_PAD
    parse_brport_mcast_flood,
    parse_brport_mcast_to_ucast,
    parse_brport_vlan_tunnel,
    parse_brport_bast_flood,
    parse_brport_group_fwd_mask,
    parse_brport_neigh_suppress,
    parse_brport_isolated,
    parse_brport_backup_port,
    parse_brport_mrp_ring_open,
    parse_brport_mrp_in_open,
    parse_brport_mcast_eht_hosts_limit,
    parse_brport_mcast_eht_hosts_cnt,
];

pub(crate) fn parse_bridge_port_info(
    raw: &[u8],
) -> Result<BridgePortInfo, NisporError> {
    let nlas = NlasIterator::new(raw);
    let mut port_info = BridgePortInfo::default();
    // TODO: Dup with parse_bond_info
    for nla in nlas {
        match nla {
            Ok(nla) => {
                if let Some(func) =
                    NLA_PORT_PARSE_FUNS.get::<usize>(nla.kind().into())
                {
                    func(nla.value(), &mut port_info)?;
                } else {
                    log::warn!(
                        "Unhandled BRIDGE_PORT_INFO {} {:?}",
                        nla.kind(),
                        nla.value()
                    );
                }
            }
            Err(e) => {
                log::warn!("{}", e);
            }
        }
    }
    Ok(port_info)
}

fn parse_as_bridge_id(data: &[u8]) -> Result<String, NisporError> {
    let err_msg = "wrong index at bridge_id parsing";
    Ok(format!(
        "{:02x}{:02x}.{:02x}{:02x}{:02x}{:02x}{:02x}{:02x}",
        data.get(0)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        data.get(1)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        data.get(2)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        data.get(3)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        data.get(4)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        data.get(5)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        data.get(6)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        data.get(7)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
    ))
}
