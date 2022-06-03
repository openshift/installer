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
use crate::Iface;
use crate::IfaceType;
use crate::NisporError;
use netlink_packet_route::rtnl::link::nlas;
use netlink_packet_route::rtnl::link::nlas::InfoMacVlan;
use netlink_packet_route::rtnl::link::nlas::InfoMacVtap;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

const ETH_ALEN: usize = 6;

const MACVLAN_MODE_PRIVATE: u32 = 1;
const MACVLAN_MODE_VEPA: u32 = 2;
const MACVLAN_MODE_BRIDGE: u32 = 4;
const MACVLAN_MODE_PASSTHRU: u32 = 8;
const MACVLAN_MODE_SOURCE: u32 = 16;

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[serde(rename_all = "lowercase")]
#[non_exhaustive]
pub enum MacVlanMode {
    /* don't talk to other macvlans */
    Private,
    /* talk to other ports through ext bridge */
    Vepa,
    /* talk to bridge ports directly */
    Bridge,
    /* take over the underlying device */
    #[serde(rename = "passthru")]
    PassThrough,
    /* use source MAC address list to assign */
    Source,
    Other(u32),
    Unknown,
}

impl Default for MacVlanMode {
    fn default() -> Self {
        MacVlanMode::Unknown
    }
}

impl From<u32> for MacVlanMode {
    fn from(d: u32) -> Self {
        match d {
            MACVLAN_MODE_PRIVATE => Self::Private,
            MACVLAN_MODE_VEPA => Self::Vepa,
            MACVLAN_MODE_BRIDGE => Self::Bridge,
            MACVLAN_MODE_PASSTHRU => Self::PassThrough,
            MACVLAN_MODE_SOURCE => Self::Source,
            _ => Self::Other(d),
        }
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct MacVlanInfo {
    pub base_iface: String,
    pub mode: MacVlanMode,
    pub flags: u16,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub allowed_mac_addresses: Option<Vec<String>>,
}

pub(crate) fn get_mac_vlan_info(
    data: &nlas::InfoData,
) -> Result<Option<MacVlanInfo>, NisporError> {
    let mut macv_info = MacVlanInfo::default();
    if let nlas::InfoData::MacVlan(infos) = data {
        for info in infos {
            if let InfoMacVlan::Mode(d) = *info {
                macv_info.mode = d.into();
            } else if let InfoMacVlan::Flags(d) = *info {
                macv_info.flags = d;
            } else if let InfoMacVlan::MacAddrData(d) = info {
                let mut addrs = Vec::new();
                for macvlan in d {
                    if let InfoMacVlan::MacAddr(mac_d) = macvlan {
                        addrs.push(parse_as_mac(ETH_ALEN, mac_d)?);
                    }
                }
                macv_info.allowed_mac_addresses = Some(addrs);
            } else {
                log::warn!("Unknown MAC VLAN info {:?}", info)
            }
        }
        Ok(Some(macv_info))
    } else if let nlas::InfoData::MacVtap(infos) = data {
        for info in infos {
            if let InfoMacVtap::Mode(d) = *info {
                macv_info.mode = d.into();
            } else if let InfoMacVtap::Flags(d) = *info {
                macv_info.flags = d;
            } else if let InfoMacVtap::MacAddrData(d) = info {
                let mut addrs = Vec::new();
                for macvtap in d {
                    if let InfoMacVtap::MacAddr(mac_d) = macvtap {
                        addrs.push(parse_as_mac(ETH_ALEN, mac_d)?);
                    }
                }
                macv_info.allowed_mac_addresses = Some(addrs);
            } else {
                log::warn!("Unknown MAC VTAP info {:?}", info)
            }
        }
        Ok(Some(macv_info))
    } else {
        Ok(None)
    }
}

pub(crate) fn mac_vlan_iface_tidy_up(
    iface_states: &mut HashMap<String, Iface>,
) {
    convert_base_iface_index_to_name(iface_states);
}

fn convert_base_iface_index_to_name(iface_states: &mut HashMap<String, Iface>) {
    let mut index_to_name = HashMap::new();
    for iface in iface_states.values() {
        index_to_name.insert(format!("{}", iface.index), iface.name.clone());
    }
    for iface in iface_states.values_mut() {
        if iface.iface_type != IfaceType::MacVlan
            && iface.iface_type != IfaceType::MacVtap
        {
            continue;
        }
        if let Some(ref mut info) = iface.mac_vlan {
            if let Some(base_iface_name) = index_to_name.get(&info.base_iface) {
                info.base_iface = base_iface_name.clone();
            }
        } else if let Some(ref mut info) = iface.mac_vtap {
            if let Some(base_iface_name) = index_to_name.get(&info.base_iface) {
                info.base_iface = base_iface_name.clone();
            }
        }
    }
}
