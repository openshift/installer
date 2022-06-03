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

use crate::error::ErrorKind;
use crate::netlink::parse_as_ipv4;
use crate::netlink::parse_as_ipv6;
use crate::route::AddressFamily;
use crate::route::RouteProtocol;
use crate::NisporError;
use futures::stream::TryStreamExt;
use netlink_packet_route::rtnl::rule::nlas::Nla;
use netlink_packet_route::RuleMessage;
use rtnetlink::new_connection;
use rtnetlink::IpVersion;
use serde::{Deserialize, Serialize};

const FR_ACT_TO_TBL: u8 = 1;
const FR_ACT_GOTO: u8 = 2;
const FR_ACT_NOP: u8 = 4;
const FR_ACT_BLACKHOLE: u8 = 32;
const FR_ACT_UNREACHABLE: u8 = 64;
const FR_ACT_PROHIBIT: u8 = 128;

const RT_TABLE_UNSPEC: u8 = 0;

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[serde(rename_all = "snake_case")]
pub enum RuleAction {
    /* Pass to fixed table or l3mdev */
    Table,
    /* Jump to another rule */
    Goto,
    /* No operation */
    Nop,
    /* Drop without notification */
    Blackhole,
    /* Drop with ENETUNREACH */
    Unreachable,
    /* Drop with EACCES */
    Prohibit,
    Other(u8),
    Unknown,
}

impl From<u8> for RuleAction {
    fn from(d: u8) -> Self {
        match d {
            FR_ACT_TO_TBL => Self::Table,
            FR_ACT_GOTO => Self::Goto,
            FR_ACT_NOP => Self::Nop,
            FR_ACT_BLACKHOLE => Self::Blackhole,
            FR_ACT_UNREACHABLE => Self::Unreachable,
            FR_ACT_PROHIBIT => Self::Prohibit,
            _ => Self::Other(d),
        }
    }
}

impl Default for RuleAction {
    fn default() -> Self {
        Self::Unknown
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct RouteRule {
    pub action: RuleAction,
    pub address_family: AddressFamily,
    pub flags: u32,
    pub tos: u8,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub table: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub dst: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub src: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub iif: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub oif: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub goto: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub priority: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub fw_mark: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub fw_mask: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub mask: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub flow: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tun_id: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub suppress_ifgroup: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub suppress_prefix_len: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub protocol: Option<RouteProtocol>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ip_proto: Option<AddressFamily>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub src_port_range: Option<Vec<u8>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub dst_port_range: Option<Vec<u8>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub l3mdev: Option<bool>,
}

pub(crate) async fn get_route_rules() -> Result<Vec<RouteRule>, NisporError> {
    let mut rules = Vec::new();
    let (connection, handle, _) = new_connection()?;
    tokio::spawn(connection);

    let mut links = handle.rule().get(IpVersion::V6).execute();
    while let Some(rt_msg) = links.try_next().await? {
        rules.push(get_rule(rt_msg)?);
    }
    let mut links = handle.rule().get(IpVersion::V4).execute();
    while let Some(rt_msg) = links.try_next().await? {
        rules.push(get_rule(rt_msg)?);
    }
    Ok(rules)
}

fn get_rule(rule_msg: RuleMessage) -> Result<RouteRule, NisporError> {
    let mut rl = RouteRule::default();
    let header = &rule_msg.header;
    rl.address_family = header.family.into();
    let src_prefix_len = header.src_len;
    let dst_prefix_len = header.dst_len;
    rl.tos = header.tos;
    rl.action = header.action.into();
    if header.table > RT_TABLE_UNSPEC {
        rl.table = Some(header.table.into());
    }
    let family = &rl.address_family;
    for nla in &rule_msg.nlas {
        match nla {
            Nla::Destination(ref d) => {
                rl.dst = Some(format!(
                    "{}/{}",
                    _addr_to_string(d, family)?,
                    dst_prefix_len,
                ));
            }
            Nla::Source(ref d) => {
                rl.src = Some(format!(
                    "{}/{}",
                    _addr_to_string(d, family)?,
                    src_prefix_len,
                ));
            }
            Nla::Iifname(ref d) => {
                rl.iif = Some(d.clone().to_string());
            }
            Nla::OifName(ref d) => {
                rl.oif = Some(d.clone().to_string());
            }
            Nla::Goto(ref d) => {
                rl.goto = Some(*d);
            }
            Nla::Priority(ref d) => {
                rl.priority = Some(*d);
            }
            Nla::FwMark(ref d) => {
                rl.fw_mark = Some(*d);
            }
            Nla::FwMask(ref d) => {
                rl.fw_mask = Some(*d);
            }
            Nla::Flow(ref d) => {
                rl.flow = Some(*d);
            }
            Nla::TunId(ref d) => {
                rl.tun_id = Some(*d);
            }
            Nla::SuppressIfGroup(ref d) => {
                if *d != std::u32::MAX {
                    rl.suppress_ifgroup = Some(*d);
                }
            }
            Nla::SuppressPrefixLen(ref d) => {
                if *d != std::u32::MAX {
                    rl.suppress_prefix_len = Some(*d);
                }
            }
            Nla::Table(ref d) => {
                if *d > RT_TABLE_UNSPEC.into() {
                    rl.table = Some(*d);
                }
            }
            Nla::Protocol(ref d) => {
                rl.protocol = Some((*d).into());
            }
            Nla::IpProto(ref d) => {
                rl.ip_proto = Some((*d).into());
            }
            Nla::L3MDev(ref d) => {
                rl.l3mdev = Some(*d > 0);
            }
            _ => log::warn!("Unknown NLA message for route rule {:?}", nla),
        }
    }

    Ok(rl)
}

fn _addr_to_string(
    data: &[u8],
    family: &AddressFamily,
) -> Result<String, NisporError> {
    match family {
        AddressFamily::IPv4 => Ok(parse_as_ipv4(data)?.to_string()),
        AddressFamily::IPv6 => Ok(parse_as_ipv6(data)?.to_string()),
        _ => Err(NisporError {
            kind: ErrorKind::NisporBug,
            msg: "The rule is not a valid IPv4 and IPv6 rule".to_string(),
        }),
    }
}
