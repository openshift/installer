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

use crate::ifaces::mac_vlan::get_mac_vlan_info;
use crate::ifaces::mac_vlan::MacVlanInfo;
use crate::ifaces::mac_vlan::MacVlanMode;
use crate::NisporError;
use netlink_packet_route::rtnl::link::nlas;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[serde(rename_all = "lowercase")]
#[non_exhaustive]
pub enum MacVtapMode {
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

impl Default for MacVtapMode {
    fn default() -> Self {
        MacVtapMode::Unknown
    }
}

impl From<MacVlanMode> for MacVtapMode {
    fn from(d: MacVlanMode) -> Self {
        match d {
            MacVlanMode::Private => Self::Private,
            MacVlanMode::Vepa => Self::Vepa,
            MacVlanMode::Bridge => Self::Bridge,
            MacVlanMode::PassThrough => Self::PassThrough,
            MacVlanMode::Source => Self::Source,
            MacVlanMode::Unknown => Self::Unknown,
            MacVlanMode::Other(u32) => Self::Other(u32),
        }
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct MacVtapInfo {
    pub base_iface: String,
    pub mode: MacVtapMode,
    pub flags: u16,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub allowed_mac_addresses: Option<Vec<String>>,
}

impl From<MacVlanInfo> for MacVtapInfo {
    fn from(d: MacVlanInfo) -> Self {
        Self {
            base_iface: d.base_iface,
            mode: MacVtapMode::from(d.mode),
            flags: d.flags,
            allowed_mac_addresses: d.allowed_mac_addresses,
        }
    }
}

pub(crate) fn get_mac_vtap_info(
    data: &nlas::InfoData,
) -> Result<Option<MacVtapInfo>, NisporError> {
    if let Some(info) = get_mac_vlan_info(data)? {
        Ok(Some(info.into()))
    } else {
        Ok(None)
    }
}
