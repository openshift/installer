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

use std::collections::HashMap;
use std::net::IpAddr;
use std::str::FromStr;

use netlink_packet_route::rtnl::{
    address::nlas::{CacheInfo, Nla, ADDRESSS_CACHE_INFO_LEN},
    AddressMessage,
};
use netlink_packet_utils::Emitable;
use serde::{Deserialize, Serialize};

use crate::{Iface, IfaceConf, NisporError};

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct Ipv4Info {
    pub addresses: Vec<Ipv4AddrInfo>,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct Ipv4AddrInfo {
    pub address: String,
    pub prefix_len: u8,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub peer: Option<String>,
    // The renaming seonds for this address be valid
    pub valid_lft: String,
    // The renaming seonds for this address be preferred
    pub preferred_lft: String,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct Ipv6Info {
    pub addresses: Vec<Ipv6AddrInfo>,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct Ipv6AddrInfo {
    pub address: String,
    pub prefix_len: u8,
    // The renaming seonds for this address be valid
    pub valid_lft: String,
    // The renaming seonds for this address be preferred
    pub preferred_lft: String,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct IpConf {
    pub addresses: Vec<IpAddrConf>,
}

impl From<&Ipv4Info> for IpConf {
    fn from(info: &Ipv4Info) -> Self {
        let mut addrs = Vec::new();
        for addr_info in &info.addresses {
            if addr_info.valid_lft == "forever" {
                addrs.push(IpAddrConf {
                    remove: false,
                    address: addr_info.address.clone(),
                    prefix_len: addr_info.prefix_len,
                    preferred_lft: addr_info.preferred_lft.clone(),
                    valid_lft: addr_info.valid_lft.clone(),
                });
            }
        }
        Self { addresses: addrs }
    }
}

impl From<&Ipv6Info> for IpConf {
    fn from(info: &Ipv6Info) -> Self {
        let mut addrs = Vec::new();
        for addr_info in &info.addresses {
            if addr_info.valid_lft == "forever" {
                addrs.push(IpAddrConf {
                    remove: false,
                    address: addr_info.address.clone(),
                    prefix_len: addr_info.prefix_len,
                    preferred_lft: addr_info.preferred_lft.clone(),
                    valid_lft: addr_info.valid_lft.clone(),
                });
            }
        }
        Self { addresses: addrs }
    }
}

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum IpFamily {
    Ipv4,
    Ipv6,
}

#[derive(
    Serialize, Deserialize, Debug, PartialEq, Eq, Hash, Clone, Default,
)]
#[non_exhaustive]
pub struct IpAddrConf {
    #[serde(default)]
    pub remove: bool,
    pub address: String,
    pub prefix_len: u8,
    #[serde(default)]
    pub valid_lft: String,
    #[serde(default)]
    pub preferred_lft: String,
}

impl IpConf {
    pub async fn apply(
        &self,
        handle: &rtnetlink::Handle,
        cur_iface: &Iface,
        family: IpFamily,
    ) -> Result<(), NisporError> {
        log::warn!("WARN: Deprecated, please use NetConf::apply() instead");
        let iface = match family {
            IpFamily::Ipv4 => IfaceConf {
                ipv4: Some(self.clone()),
                ..Default::default()
            },
            IpFamily::Ipv6 => IfaceConf {
                ipv6: Some(self.clone()),
                ..Default::default()
            },
        };
        let ifaces = vec![&iface];
        let mut cur_ifaces = HashMap::new();
        cur_ifaces.insert(cur_iface.name.clone(), cur_iface.clone());
        change_ips(handle, &ifaces, &cur_ifaces).await
    }
}

pub(crate) fn is_ipv6_addr(addr: &str) -> bool {
    addr.contains(':')
}

pub(crate) async fn change_ips(
    handle: &rtnetlink::Handle,
    ifaces: &[&IfaceConf],
    cur_ifaces: &HashMap<String, Iface>,
) -> Result<(), NisporError> {
    for iface in ifaces {
        if let Some(cur_iface) = cur_ifaces.get(&iface.name) {
            if let Some(ip_conf) = iface.ipv4.as_ref() {
                apply_ip_conf(handle, cur_iface.index, ip_conf, IpFamily::Ipv4)
                    .await?;
            }
            if let Some(ip_conf) = iface.ipv6.as_ref() {
                apply_ip_conf(handle, cur_iface.index, ip_conf, IpFamily::Ipv6)
                    .await?;
            }
        }
    }

    Ok(())
}

async fn apply_ip_conf(
    handle: &rtnetlink::Handle,
    iface_index: u32,
    ip_conf: &IpConf,
    ip_family: IpFamily,
) -> Result<(), NisporError> {
    for addr_conf in &ip_conf.addresses {
        if addr_conf.remove {
            let mut nl_msg = AddressMessage::default();
            nl_msg.header.index = iface_index;
            nl_msg.header.prefix_len = addr_conf.prefix_len;
            nl_msg.header.family = match ip_family {
                IpFamily::Ipv4 => libc::AF_INET as u8,
                IpFamily::Ipv6 => libc::AF_INET6 as u8,
            };
            nl_msg.nlas.push(Nla::Address(
                match ip_addr_str_to_enum(&addr_conf.address)? {
                    IpAddr::V4(i) => i.octets().to_vec(),
                    IpAddr::V6(i) => i.octets().to_vec(),
                },
            ));
            if let Err(e) = handle.address().del(nl_msg).execute().await {
                if let rtnetlink::Error::NetlinkError(ref e) = e {
                    if e.code == -libc::EADDRNOTAVAIL {
                        return Ok(());
                    }
                }
                return Err(e.into());
            }
        } else {
            let mut req = handle
                .address()
                .add(
                    iface_index,
                    ip_addr_str_to_enum(&addr_conf.address)?,
                    addr_conf.prefix_len,
                )
                .replace();
            if is_dynamic_ip(&addr_conf.preferred_lft, &addr_conf.valid_lft) {
                handle_dynamic_ip(
                    req.message_mut(),
                    &addr_conf.preferred_lft,
                    &addr_conf.valid_lft,
                )?;
            }
            req.execute().await?;
        }
    }
    Ok(())
}

fn ip_addr_str_to_enum(address: &str) -> Result<IpAddr, NisporError> {
    Ok(if is_ipv6_addr(address) {
        IpAddr::V6(std::net::Ipv6Addr::from_str(address)?)
    } else {
        IpAddr::V4(std::net::Ipv4Addr::from_str(address)?)
    })
}

fn gen_cache_info_u8(
    preferred_lft: &str,
    valid_lft: &str,
) -> Result<[u8; ADDRESSS_CACHE_INFO_LEN], NisporError> {
    let cache_info = CacheInfo {
        ifa_preferred: parse_lft_sec("preferred_lft", preferred_lft)?,
        ifa_valid: parse_lft_sec("valid_lft", valid_lft)?,
        ..Default::default()
    };
    let mut buff = [0u8; ADDRESSS_CACHE_INFO_LEN];
    cache_info.emit(&mut buff);
    Ok(buff)
}

fn handle_dynamic_ip(
    nl_msg: &mut AddressMessage,
    preferred_lft: &str,
    valid_lft: &str,
) -> Result<(), NisporError> {
    nl_msg.nlas.push(Nla::CacheInfo(
        gen_cache_info_u8(preferred_lft, valid_lft)?.to_vec(),
    ));
    Ok(())
}

fn is_dynamic_ip(preferred_lft: &str, valid_lft: &str) -> bool {
    (preferred_lft != "forever" && !preferred_lft.is_empty())
        || (valid_lft != "forever" && !valid_lft.is_empty())
}

fn parse_lft_sec(name: &str, lft_str: &str) -> Result<i32, NisporError> {
    let e = NisporError::invalid_argument(format!(
        "Invalid {} format: expect format 50sec, got {}",
        name, lft_str
    ));
    match lft_str.strip_suffix("sec") {
        Some(a) => a.parse().map_err(|_| {
            log::error!("{}", e);
            e
        }),
        None => {
            log::error!("{}", e);
            Err(e)
        }
    }
}

pub(crate) fn parse_ip_addr_str(
    ip_addr_str: &str,
) -> Result<IpAddr, NisporError> {
    IpAddr::from_str(ip_addr_str).map_err(|e| {
        let e = NisporError::invalid_argument(format!(
            "Invalid IP address {}: {}",
            ip_addr_str, e
        ));
        log::error!("{}", e);
        e
    })
}

pub(crate) fn parse_ip_net_addr_str(
    ip_net_str: &str,
) -> Result<(IpAddr, u8), NisporError> {
    let splits: Vec<&str> = ip_net_str.split('/').collect();
    if splits.len() > 2 || splits.is_empty() {
        let e = NisporError::invalid_argument(format!(
            "Invalid IP network address {}",
            ip_net_str,
        ));
        log::error!("{}", e);
        return Err(e);
    }
    let addr_str = splits[0];
    let prefix_len = if let Some(prefix_len_str) = splits.get(1) {
        prefix_len_str.parse::<u8>().map_err(|e| {
            let e = NisporError::invalid_argument(format!(
                "Invalid IP network prefix {}: {}",
                ip_net_str, e
            ));
            log::error!("{}", e);
            e
        })?
    } else if is_ipv6_addr(addr_str) {
        128
    } else {
        32
    };
    Ok((parse_ip_addr_str(addr_str)?, prefix_len))
}
