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
use crate::netlink::parse_as_u32;
use crate::netlink::parse_as_u64;
use crate::NisporError;
use netlink_packet_route::rtnl::nlas::NlasIterator;
use serde::{Deserialize, Serialize};

const IFLA_VF_MAC: u16 = 1;
const IFLA_VF_VLAN: u16 = 2;
const IFLA_VF_TX_RATE: u16 = 3;
const IFLA_VF_SPOOFCHK: u16 = 4;
const IFLA_VF_LINK_STATE: u16 = 5;
const IFLA_VF_RATE: u16 = 6;
const IFLA_VF_RSS_QUERY_EN: u16 = 7;
const IFLA_VF_STATS: u16 = 8;
const IFLA_VF_TRUST: u16 = 9;
const IFLA_VF_IB_NODE_GUID: u16 = 10;
const IFLA_VF_IB_PORT_GUID: u16 = 11;
const IFLA_VF_VLAN_LIST: u16 = 12;
const IFLA_VF_BROADCAST: u16 = 13;

const IFLA_VF_LINK_STATE_AUTO: u32 = 0;
const IFLA_VF_LINK_STATE_ENABLE: u32 = 1;
const IFLA_VF_LINK_STATE_DISABLE: u32 = 2;

const IFLA_VF_STATS_RX_PACKETS: u16 = 0;
const IFLA_VF_STATS_TX_PACKETS: u16 = 1;
const IFLA_VF_STATS_RX_BYTES: u16 = 2;
const IFLA_VF_STATS_TX_BYTES: u16 = 3;
const IFLA_VF_STATS_BROADCAST: u16 = 4;
const IFLA_VF_STATS_MULTICAST: u16 = 5;
// const IFLA_VF_STATS_PAD: u16 = 6;
const IFLA_VF_STATS_RX_DROPPED: u16 = 7;
const IFLA_VF_STATS_TX_DROPPED: u16 = 8;

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[serde(rename_all = "snake_case")]
#[non_exhaustive]
pub enum VfLinkState {
    Auto,
    Enable,
    Disable,
    Other(u32),
    Unknown,
}

impl Default for VfLinkState {
    fn default() -> Self {
        VfLinkState::Unknown
    }
}
impl From<u32> for VfLinkState {
    fn from(d: u32) -> Self {
        match d {
            IFLA_VF_LINK_STATE_AUTO => VfLinkState::Auto,
            IFLA_VF_LINK_STATE_ENABLE => VfLinkState::Enable,
            IFLA_VF_LINK_STATE_DISABLE => VfLinkState::Disable,
            _ => VfLinkState::Other(d),
        }
    }
}

#[derive(Serialize, Deserialize, PartialEq, Debug, Clone, Default)]
#[non_exhaustive]
pub struct VfState {
    pub rx_packets: u64,
    pub tx_packets: u64,
    pub rx_bytes: u64,
    pub tx_bytes: u64,
    pub broadcast: u64,
    pub multicast: u64,
    pub rx_dropped: u64,
    pub tx_dropped: u64,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct SriovInfo {
    pub vfs: Vec<VfInfo>,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct VfInfo {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub iface_name: Option<String>,
    pub id: u32,
    pub mac: String,
    pub broadcast: String,
    // 0 disables VLAN filter
    pub vlan_id: u32,
    pub qos: u32,
    // Max TX bandwidth in Mbps, 0 disables throttling
    pub tx_rate: u32,
    pub spoof_check: bool,
    pub link_state: VfLinkState,
    // Min Bandwidth in Mbps
    pub min_tx_rate: u32,
    // Max Bandwidth in Mbps
    pub max_tx_rate: u32,
    pub query_rss: bool,
    pub state: VfState,
    pub trust: bool,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ib_node_guid: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ib_port_guid: Option<String>,
}

pub(crate) fn get_sriov_info(
    pf_iface_name: &str,
    raw: &[u8],
    mac_len: Option<usize>,
) -> Result<SriovInfo, NisporError> {
    let mut sriov_info = SriovInfo::default();
    let ports = NlasIterator::new(raw);
    for port in ports {
        let mut vf_info = VfInfo::default();
        let port = port?;
        let port_nlas = NlasIterator::new(port.value());
        for nla in port_nlas {
            let nla = nla?;
            match nla.kind() {
                IFLA_VF_MAC => {
                    vf_info.id = parse_as_u32(nla.value())?;
                    vf_info.iface_name =
                        get_vf_iface_name(pf_iface_name, &vf_info.id);
                    vf_info.mac = parse_vf_mac(
                        nla.value().get(4..).ok_or_else(|| {
                            NisporError::bug("invalid index into nla".into())
                        })?,
                        mac_len,
                    )?;
                }
                IFLA_VF_VLAN => {
                    vf_info.vlan_id = parse_as_u32(
                        nla.value().get(4..).ok_or_else(|| {
                            NisporError::bug("invalid index into nla".into())
                        })?,
                    )?;
                    vf_info.qos = parse_as_u32(
                        nla.value().get(8..).ok_or_else(|| {
                            NisporError::bug("invalid index into nla".into())
                        })?,
                    )?;
                }
                IFLA_VF_TX_RATE => {
                    vf_info.tx_rate = parse_as_u32(
                        nla.value().get(4..).ok_or_else(|| {
                            NisporError::bug("invalid index into nla".into())
                        })?,
                    )?;
                }
                IFLA_VF_SPOOFCHK => {
                    let d = parse_as_u32(nla.value().get(4..).ok_or_else(
                        || NisporError::bug("invalid index into nla".into()),
                    )?)?;
                    vf_info.spoof_check = d > 0 && d != std::u32::MAX;
                }
                IFLA_VF_LINK_STATE => {
                    vf_info.link_state = parse_as_u32(
                        nla.value().get(4..).ok_or_else(|| {
                            NisporError::bug("invalid index into nla".into())
                        })?,
                    )?
                    .into();
                }
                IFLA_VF_RATE => {
                    vf_info.min_tx_rate = parse_as_u32(
                        nla.value().get(4..).ok_or_else(|| {
                            NisporError::bug("invalid index into nla".into())
                        })?,
                    )?;
                    vf_info.max_tx_rate = parse_as_u32(
                        nla.value().get(8..).ok_or_else(|| {
                            NisporError::bug("invalid index into nla".into())
                        })?,
                    )?;
                }
                IFLA_VF_RSS_QUERY_EN => {
                    let d = parse_as_u32(nla.value().get(4..).ok_or_else(
                        || NisporError::bug("invalid index into nla".into()),
                    )?)?;
                    vf_info.query_rss = d > 0 && d != std::u32::MAX;
                }
                IFLA_VF_STATS => {
                    vf_info.state = parse_vf_stats(nla.value())?;
                }
                IFLA_VF_TRUST => {
                    let d = parse_as_u32(nla.value().get(4..).ok_or_else(
                        || NisporError::bug("invalid index into nla".into()),
                    )?)?;
                    vf_info.trust = d > 0 && d != std::u32::MAX;
                }
                IFLA_VF_IB_NODE_GUID => {
                    vf_info.ib_node_guid =
                        Some(format!("{:X}", parse_as_u64(nla.value())?));
                }
                IFLA_VF_IB_PORT_GUID => {
                    vf_info.ib_port_guid =
                        Some(format!("{:X}", parse_as_u64(nla.value())?));
                }
                IFLA_VF_VLAN_LIST => {
                    // The kernel just store IFLA_VF_VLAN in a list with single
                    // content. The the vlan protocol is always 0 untile
                    // someone set it via IFLA_VF_VLAN_LIST. The iproute does
                    // not support this, so I doubt anyone is using this.
                }
                IFLA_VF_BROADCAST => {
                    vf_info.broadcast = parse_vf_mac(nla.value(), mac_len)?;
                }
                _ => {
                    log::warn!(
                        "Unhandled SRIOV NLA {} {:?}",
                        nla.kind(),
                        nla.value()
                    );
                }
            }
        }

        sriov_info.vfs.push(vf_info);
    }
    Ok(sriov_info)
}

fn parse_vf_mac(
    raw: &[u8],
    mac_len: Option<usize>,
) -> Result<String, NisporError> {
    match mac_len {
        Some(mac_len) => parse_as_mac(mac_len, raw),
        None => parse_as_mac(32, raw),
    }
}

fn parse_vf_stats(raw: &[u8]) -> Result<VfState, NisporError> {
    let mut state = VfState::default();
    let nlas = NlasIterator::new(raw);
    for nla in nlas {
        let nla = nla?;
        match nla.kind() {
            IFLA_VF_STATS_RX_PACKETS => {
                state.rx_packets = parse_as_u64(nla.value())?;
            }
            IFLA_VF_STATS_TX_PACKETS => {
                state.tx_packets = parse_as_u64(nla.value())?;
            }
            IFLA_VF_STATS_RX_BYTES => {
                state.rx_bytes = parse_as_u64(nla.value())?;
            }
            IFLA_VF_STATS_TX_BYTES => {
                state.tx_bytes = parse_as_u64(nla.value())?;
            }
            IFLA_VF_STATS_BROADCAST => {
                state.broadcast = parse_as_u64(nla.value())?;
            }
            IFLA_VF_STATS_MULTICAST => {
                state.multicast = parse_as_u64(nla.value())?;
            }
            IFLA_VF_STATS_RX_DROPPED => {
                state.rx_dropped = parse_as_u64(nla.value())?;
            }
            IFLA_VF_STATS_TX_DROPPED => {
                state.tx_dropped = parse_as_u64(nla.value())?;
            }
            _ => log::warn!(
                "Unhandled IFLA_VF_STATS {}, {:?}",
                nla.kind(),
                nla.value()
            ),
        }
    }
    Ok(state)
}

// Currently there is no valid netlink way to get information as the kernel code
// is in at PCI level: drivers/pci/iov.c
// We use sysfs content /sys/class/net/<pf_name>/devices/virtfn<sriov_id>/net/
fn get_vf_iface_name(pf_name: &str, sriov_id: &u32) -> Option<String> {
    let sysfs_path =
        format!("/sys/class/net/{}/device/virtfn{}/net/", pf_name, sriov_id);
    read_folder(&sysfs_path).pop()
}

fn read_folder(folder_path: &str) -> Vec<String> {
    let mut folder_contents = Vec::new();
    let fd = match std::fs::read_dir(folder_path) {
        Ok(f) => f,
        Err(e) => {
            log::warn!("Failed to read dir {}: {}", folder_path, e);
            return folder_contents;
        }
    };
    for entry in fd {
        let entry = match entry {
            Ok(e) => e,
            Err(e) => {
                log::warn!("Failed to read dir {}: {}", folder_path, e);
                continue;
            }
        };
        let path = entry.path();
        if let Ok(content) = path.strip_prefix(folder_path) {
            if let Some(content_str) = content.to_str() {
                folder_contents.push(content_str.to_string());
            }
        }
    }
    folder_contents
}
