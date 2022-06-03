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

mod bond;
mod bridge;
mod ethtool;
mod iface;
mod ipoib;
mod mac_vlan;
mod mac_vtap;
mod sriov;
mod tun;
mod veth;
mod vlan;
mod vrf;
mod vxlan;

pub use crate::ifaces::bond::*;
pub use crate::ifaces::bridge::*;
pub use crate::ifaces::ethtool::*;
pub use crate::ifaces::iface::*;
pub use crate::ifaces::ipoib::{IpoibInfo, IpoibMode};
pub use crate::ifaces::mac_vlan::*;
pub use crate::ifaces::mac_vtap::*;
pub use crate::ifaces::sriov::*;
pub use crate::ifaces::tun::*;
pub use crate::ifaces::veth::*;
pub use crate::ifaces::vlan::*;
pub use crate::ifaces::vrf::*;
pub use crate::ifaces::vxlan::*;

use crate::ifaces::bond::bond_iface_tidy_up;
use crate::ifaces::bridge::bridge_iface_tidy_up;
use crate::ifaces::ethtool::get_ethtool_infos;
use crate::ifaces::iface::{
    fill_bridge_vlan_info, parse_nl_msg_to_iface,
    parse_nl_msg_to_name_and_index,
};
use crate::ifaces::ipoib::ipoib_iface_tidy_up;
use crate::ifaces::mac_vlan::mac_vlan_iface_tidy_up;
use crate::ifaces::veth::veth_iface_tidy_up;
use crate::ifaces::vlan::vlan_iface_tidy_up;
use crate::ifaces::vrf::vrf_iface_tidy_up;
use crate::ifaces::vxlan::vxlan_iface_tidy_up;
use crate::ip::change_ips;
use crate::netlink::fill_ip_addr;
use crate::NisporError;

use futures::stream::TryStreamExt;
use netlink_packet_route::rtnl::constants::AF_BRIDGE;
use netlink_packet_route::rtnl::constants::AF_UNSPEC;
use netlink_packet_route::rtnl::constants::RTEXT_FILTER_BRVLAN_COMPRESSED;
use netlink_packet_route::rtnl::constants::RTEXT_FILTER_VF;
use rtnetlink::new_connection;
use std::collections::HashMap;

pub(crate) async fn get_ifaces() -> Result<HashMap<String, Iface>, NisporError>
{
    let mut iface_states: HashMap<String, Iface> = HashMap::new();
    let (connection, handle, _) = new_connection()?;
    tokio::spawn(connection);

    let mut links = handle
        .link()
        .get()
        .set_filter_mask(AF_UNSPEC as u8, RTEXT_FILTER_VF)
        .execute();
    while let Some(nl_msg) = links.try_next().await? {
        if let Some(iface_state) = parse_nl_msg_to_iface(&nl_msg)? {
            iface_states.insert(iface_state.name.clone(), iface_state);
        }
    }
    let mut addrs = handle.address().get().execute();
    while let Some(nl_msg) = addrs.try_next().await? {
        fill_ip_addr(&mut iface_states, &nl_msg)?;
    }
    let mut br_vlan_links = handle
        .link()
        .get()
        .set_filter_mask(AF_BRIDGE as u8, RTEXT_FILTER_BRVLAN_COMPRESSED)
        .execute();
    while let Some(nl_msg) = br_vlan_links.try_next().await? {
        fill_bridge_vlan_info(&mut iface_states, &nl_msg)?;
    }

    match get_ethtool_infos().await {
        Ok(mut ethtool_infos) => {
            ifaces_merge_ethool_infos(&mut iface_states, &mut ethtool_infos);
        }
        Err(e) => {
            // Ethtool is considered as optional
            log::warn!("Failed to query ethtool info: {}", e);
        }
    };

    tidy_up(&mut iface_states);
    Ok(iface_states)
}

fn tidy_up(iface_states: &mut HashMap<String, Iface>) {
    controller_iface_index_to_name(iface_states);
    bond_iface_tidy_up(iface_states);
    bridge_iface_tidy_up(iface_states);
    vlan_iface_tidy_up(iface_states);
    vxlan_iface_tidy_up(iface_states);
    veth_iface_tidy_up(iface_states);
    vrf_iface_tidy_up(iface_states);
    mac_vlan_iface_tidy_up(iface_states);
    ipoib_iface_tidy_up(iface_states);
}

fn controller_iface_index_to_name(iface_states: &mut HashMap<String, Iface>) {
    let mut index_to_name = HashMap::new();
    for iface in iface_states.values() {
        index_to_name.insert(format!("{}", iface.index), iface.name.clone());
    }
    for iface in iface_states.values_mut() {
        if let Some(controller) = &iface.controller {
            if let Some(name) = index_to_name.get(controller) {
                iface.controller = Some(name.to_string());
            }
        }
    }
}

fn ifaces_merge_ethool_infos(
    iface_states: &mut HashMap<String, Iface>,
    ethtool_infos: &mut HashMap<String, EthtoolInfo>,
) {
    for iface in iface_states.values_mut() {
        if let Some(ethtool_info) = ethtool_infos.remove(&iface.name) {
            iface.ethtool = Some(ethtool_info)
        }
    }
}

pub(crate) async fn get_iface_name2index(
) -> Result<HashMap<String, u32>, NisporError> {
    let mut name2index: HashMap<String, u32> = HashMap::new();
    let (connection, handle, _) = new_connection()?;
    tokio::spawn(connection);

    let mut links = handle.link().get().execute();
    while let Some(nl_msg) = links.try_next().await? {
        if let Some((iface_name, iface_index)) =
            parse_nl_msg_to_name_and_index(&nl_msg)
        {
            name2index.insert(iface_name, iface_index);
        }
    }
    Ok(name2index)
}

pub(crate) async fn delete_ifaces(
    ifaces: &[(&str, u32)],
) -> Result<(), NisporError> {
    let (connection, handle, _) = new_connection()?;
    tokio::spawn(connection);
    for (iface_name, iface_index) in ifaces {
        if let Err(e) = handle.link().del(*iface_index).execute().await {
            return Err(NisporError::bug(format!(
                "Failed to delete interface {} with index {}: {}",
                iface_name, iface_index, e
            )));
        }
    }

    Ok(())
}

pub(crate) async fn create_ifaces(
    ifaces: &[&IfaceConf],
    cur_iface_name_2_index: &HashMap<String, u32>,
) -> Result<(), NisporError> {
    let (connection, handle, _) = new_connection()?;
    tokio::spawn(connection);
    for iface in ifaces {
        match iface.iface_type {
            Some(IfaceType::Bridge) => {
                BridgeConf::create(&handle, &iface.name).await?;
            }
            Some(IfaceType::Veth) => {
                if let Some(veth_conf) = &iface.veth {
                    veth_conf.create(&handle, &iface.name).await?;
                }
            }
            Some(IfaceType::Bond) => {
                BondConf::create(&handle, &iface.name).await?;
            }
            Some(IfaceType::Vlan) => {
                if let Some(vlan_conf) = &iface.vlan {
                    if let Some(base_iface_index) =
                        cur_iface_name_2_index.get(&vlan_conf.base_iface)
                    {
                        VlanConf::create(
                            &handle,
                            &iface.name,
                            vlan_conf.vlan_id,
                            *base_iface_index,
                        )
                        .await?;
                    } else {
                        return Err(NisporError::invalid_argument(format!(
                            "Base interface {} for VLAN {} not found",
                            &vlan_conf.base_iface, iface.name
                        )));
                    }
                }
            }
            Some(_) => {
                return Err(NisporError::invalid_argument(format!(
                    "Cannot create unsupported interface {:?}",
                    &iface
                )));
            }
            None => {
                return Err(NisporError::invalid_argument(format!(
                    "No interface type defined for new interface {:?}",
                    &iface
                )));
            }
        }
    }

    Ok(())
}

pub(crate) async fn change_ifaces(
    ifaces: &[&IfaceConf],
    cur_ifaces: &HashMap<String, Iface>,
) -> Result<(), NisporError> {
    let (connection, handle, _) = new_connection()?;
    tokio::spawn(connection);
    change_ifaces_mac(&handle, ifaces, cur_ifaces).await?;
    change_ifaces_controller(&handle, ifaces, cur_ifaces).await?;
    change_ifaces_state(&handle, ifaces, cur_ifaces).await?;
    change_ips(&handle, ifaces, cur_ifaces).await?;
    Ok(())
}

async fn change_ifaces_state(
    handle: &rtnetlink::Handle,
    ifaces: &[&IfaceConf],
    cur_ifaces: &HashMap<String, Iface>,
) -> Result<(), NisporError> {
    for iface in ifaces {
        if let Some(cur_iface) = cur_ifaces.get(&iface.name) {
            if cur_iface.state != iface.state {
                if iface.state == IfaceState::Up {
                    change_iface_state(handle, cur_iface.index, true).await?;
                } else if iface.state == IfaceState::Down {
                    change_iface_state(handle, cur_iface.index, false).await?;
                } else {
                    return Err(NisporError::invalid_argument(format!(
                        "Unsupported interface state in NetConf: {}",
                        iface.state
                    )));
                }
            }
        }
    }

    Ok(())
}

async fn change_ifaces_controller(
    handle: &rtnetlink::Handle,
    ifaces: &[&IfaceConf],
    cur_ifaces: &HashMap<String, Iface>,
) -> Result<(), NisporError> {
    for iface in ifaces {
        if let Some(cur_iface) = cur_ifaces.get(&iface.name) {
            if cur_iface.controller != iface.controller {
                match &iface.controller {
                    Some(ref ctrl_name) => match cur_ifaces.get(ctrl_name) {
                        None => {
                            return Err(NisporError::invalid_argument(
                                format!(
                                    "Controller interface {} not found",
                                    &ctrl_name
                                ),
                            ));
                        }
                        Some(ctrl_iface) => {
                            handle
                                .link()
                                .set(cur_iface.index)
                                .master(ctrl_iface.index)
                                .execute()
                                .await?;
                        }
                    },
                    None => {
                        handle
                            .link()
                            .set(cur_iface.index)
                            .nomaster()
                            .execute()
                            .await?;
                    }
                }
            }
        } else {
            return Err(NisporError::invalid_argument(format!(
                "Interface {} not found",
                iface.name
            )));
        }
    }

    Ok(())
}

async fn change_ifaces_mac(
    handle: &rtnetlink::Handle,
    ifaces: &[&IfaceConf],
    cur_ifaces: &HashMap<String, Iface>,
) -> Result<(), NisporError> {
    for iface in ifaces {
        if let Some(mac_addr) = &iface.mac_address {
            if let Some(cur_iface) = cur_ifaces.get(&iface.name) {
                if cur_iface.state != IfaceState::Down {
                    // We can only change MAC address when link down
                    change_iface_state(handle, cur_iface.index, false).await?;
                }
                change_iface_mac(handle, cur_iface.index, mac_addr).await?;
                if cur_iface.state == IfaceState::Up
                    && iface.state == IfaceState::Up
                {
                    // Restore the interface state
                    change_iface_state(handle, cur_iface.index, true).await?;
                }
            }
        }
    }
    Ok(())
}
