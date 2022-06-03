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

use crate::netlink::nla::parse_as_ipv4;
use crate::netlink::nla::parse_as_ipv6;
use crate::Iface;
use crate::Ipv4AddrInfo;
use crate::Ipv4Info;
use crate::Ipv6AddrInfo;
use crate::Ipv6Info;
use crate::NisporError;
use netlink_packet_route::address::{CacheInfo, CacheInfoBuffer};
use netlink_packet_route::rtnl::address::nlas::Nla;
use netlink_packet_route::rtnl::AddressMessage;
use netlink_packet_utils::Parseable;
use std::collections::HashMap;

pub(crate) const AF_INET: u8 = 2;
pub(crate) const AF_INET6: u8 = 10;

pub(crate) fn fill_ip_addr(
    iface_states: &mut HashMap<String, Iface>,
    nl_msg: &AddressMessage,
) -> Result<(), NisporError> {
    match nl_msg.header.family {
        AF_INET => {
            let (iface_index, addr) = parse_ipv4_nlas(nl_msg)?;
            let iface_name = get_iface_name_by_index(iface_states, iface_index);
            if !iface_name.is_empty() {
                let new_ip4_info = match &iface_states[&iface_name].ipv4 {
                    Some(ip_info) => {
                        let mut new_ip_info = ip_info.clone();
                        new_ip_info.addresses.push(addr);
                        new_ip_info
                    }
                    None => Ipv4Info {
                        addresses: vec![addr],
                    },
                };
                if let Some(iface) = iface_states.get_mut(&iface_name) {
                    iface.ipv4 = Some(new_ip4_info);
                }
            }
        }
        AF_INET6 => {
            let (iface_index, addr) = parse_ipv6_nlas(nl_msg)?;
            let iface_name = get_iface_name_by_index(iface_states, iface_index);
            if !iface_name.is_empty() {
                let new_ip6_info = match &iface_states[&iface_name].ipv6 {
                    Some(ip_info) => {
                        let mut new_ip_info = ip_info.clone();
                        new_ip_info.addresses.push(addr);
                        new_ip_info
                    }
                    None => Ipv6Info {
                        addresses: vec![addr],
                    },
                };
                if let Some(iface) = iface_states.get_mut(&iface_name) {
                    iface.ipv6 = Some(new_ip6_info);
                }
            }
        }
        _ => {
            log::warn!(
                "unknown address family {} {:?}",
                nl_msg.header.family,
                nl_msg
            );
        }
    };
    Ok(())
}

// TODO: remove the dupcode between parse_ipv4_nlas() and parse_ipv6_nlas()
fn parse_ipv4_nlas(
    nl_msg: &AddressMessage,
) -> Result<(u32, Ipv4AddrInfo), NisporError> {
    let iface_index = nl_msg.header.index;
    let mut addr = Ipv4AddrInfo {
        prefix_len: nl_msg.header.prefix_len,
        ..Default::default()
    };
    let mut peer = String::new();
    for nla in &nl_msg.nlas {
        if let Nla::Local(addr_vec) = nla {
            addr.address = parse_as_ipv4(addr_vec.as_slice())?.to_string();
        } else if let Nla::Address(addr_vec) = nla {
            peer = parse_as_ipv4(addr_vec.as_slice())?.to_string();
        } else if let Nla::CacheInfo(cache_info_vec) = nla {
            let cache_info = CacheInfo::parse(&CacheInfoBuffer::new(
                cache_info_vec.as_slice(),
            ))?;
            addr.preferred_lft = left_time_to_string(cache_info.ifa_preferred);
            addr.valid_lft = left_time_to_string(cache_info.ifa_valid);
        }
    }

    if peer != addr.address {
        addr.peer = Some(peer)
    }

    Ok((iface_index, addr))
}

fn parse_ipv6_nlas(
    nl_msg: &AddressMessage,
) -> Result<(u32, Ipv6AddrInfo), NisporError> {
    let iface_index = nl_msg.header.index;
    let mut addr = Ipv6AddrInfo {
        prefix_len: nl_msg.header.prefix_len,
        ..Default::default()
    };

    for nla in &nl_msg.nlas {
        if let Nla::Address(addr_vec) = nla {
            addr.address = parse_as_ipv6(addr_vec.as_slice())?.to_string();
        } else if let Nla::CacheInfo(cache_info_vec) = nla {
            let cache_info = CacheInfo::parse(&CacheInfoBuffer::new(
                cache_info_vec.as_slice(),
            ))?;
            addr.preferred_lft = left_time_to_string(cache_info.ifa_preferred);
            addr.valid_lft = left_time_to_string(cache_info.ifa_valid);
        }
    }

    Ok((iface_index, addr))
}

fn left_time_to_string(left_time: i32) -> String {
    if left_time == -1 {
        "forever".into()
    } else {
        format!("{}sec", left_time)
    }
}

fn get_iface_name_by_index(
    iface_states: &HashMap<String, Iface>,
    iface_index: u32,
) -> String {
    for (iface_name, iface) in iface_states.iter() {
        if iface.index == iface_index {
            return iface_name.clone();
        }
    }
    "".into()
}
