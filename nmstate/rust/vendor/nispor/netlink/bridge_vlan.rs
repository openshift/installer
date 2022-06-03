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

use crate::BridgeVlanEntry;
use crate::NisporError;
use netlink_packet_route::rtnl::nlas::NlasIterator;

const IFLA_BRIDGE_VLAN_INFO: u16 = 2;

// VLAN is PVID, ingress untagged;
const BRIDGE_VLAN_INFO_PVID: u16 = 1 << 1;
// VLAN egresses untagged;
const BRIDGE_VLAN_INFO_UNTAGGED: u16 = 1 << 2;
// VLAN is start of vlan range;
const BRIDGE_VLAN_INFO_RANGE_BEGIN: u16 = 1 << 3;
// VLAN is end of vlan range;
const BRIDGE_VLAN_INFO_RANGE_END: u16 = 1 << 4;

// TODO: Dup with parse_bond_info
pub(crate) fn parse_af_spec_bridge_info(
    raw: &[u8],
) -> Result<Option<Vec<BridgeVlanEntry>>, NisporError> {
    let nlas = NlasIterator::new(raw);
    let mut vlans = Vec::new();

    // TODO: Dup with parse_bond_info
    for nla in nlas {
        match nla {
            Ok(nla) => match nla.kind() {
                IFLA_BRIDGE_VLAN_INFO => {
                    if let Some(v) = parse_vlan_info(nla.value())? {
                        vlans.push(v);
                    }
                }
                _ => {
                    log::warn!(
                        "Unhandled AF_SPEC_BRIDGE_INFO: {} {:?}",
                        nla.kind(),
                        nla.value()
                    );
                }
            },
            Err(e) => log::warn!("{}", e),
        }
    }
    if !vlans.is_empty() {
        Ok(Some(merge_vlan_range(&vlans)))
    } else {
        Ok(None)
    }
}

#[derive(Debug, PartialEq, Clone, Default)]
struct KernelBridgeVlanEntry {
    vid: u16,
    is_pvid: bool, // is PVID and ingress untagged
    is_egress_untagged: bool,
    is_range_start: bool,
    is_range_end: bool,
}

fn parse_vlan_info(
    data: &[u8],
) -> Result<Option<KernelBridgeVlanEntry>, NisporError> {
    if data.len() == 4 {
        let flags = u16::from_ne_bytes([
            *data.get(0).ok_or_else(|| {
                NisporError::bug("wrong index at vlan flags".into())
            })?,
            *data.get(1).ok_or_else(|| {
                NisporError::bug("wrong index at vlan flags".into())
            })?,
        ]);
        let vid = u16::from_ne_bytes([
            *data.get(2).ok_or_else(|| {
                NisporError::bug("wrong index at vlan id".into())
            })?,
            *data.get(3).ok_or_else(|| {
                NisporError::bug("wrong index at vlan id".into())
            })?,
        ]);
        let mut entry = KernelBridgeVlanEntry {
            vid,
            ..Default::default()
        };
        entry.is_pvid = (flags & BRIDGE_VLAN_INFO_PVID) > 0;
        entry.is_egress_untagged = (flags & BRIDGE_VLAN_INFO_UNTAGGED) > 0;
        entry.is_range_start = (flags & BRIDGE_VLAN_INFO_RANGE_BEGIN) > 0;
        entry.is_range_end = (flags & BRIDGE_VLAN_INFO_RANGE_END) > 0;
        Ok(Some(entry))
    } else {
        log::warn!(
            "Invalid kernel bridge vlan info: {:?}, should be [u8;4]",
            data
        );
        Ok(None)
    }
}

fn merge_vlan_range(
    kernel_vlans: &[KernelBridgeVlanEntry],
) -> Vec<BridgeVlanEntry> {
    let mut vlans = Vec::new();
    let mut vlan_start = None;
    for k_vlan in kernel_vlans {
        match (k_vlan.is_range_start, k_vlan.is_range_end) {
            (true, false) => {
                vlan_start = Some(k_vlan.vid);
                continue;
            }
            (false, true) => {
                if let Some(start) = vlan_start {
                    vlans.push(BridgeVlanEntry {
                        vid: None,
                        vid_range: Some((start, k_vlan.vid)),
                        is_pvid: k_vlan.is_pvid,
                        is_egress_untagged: k_vlan.is_egress_untagged,
                    })
                } else {
                    log::warn!(
                        "Invalid kernel bridge vlan information: \
                        missing start VLAN for {}",
                        k_vlan.vid
                    );
                }
                vlan_start = None;
            }
            (false, false) | (true, true) => {
                vlans.push(BridgeVlanEntry {
                    vid: Some(k_vlan.vid),
                    vid_range: None,
                    is_pvid: k_vlan.is_pvid,
                    is_egress_untagged: k_vlan.is_egress_untagged,
                });
                vlan_start = None;
            }
        };
    }
    vlans
}
