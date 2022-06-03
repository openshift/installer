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

use crate::{
    ifaces::{
        bond::{
            get_bond_info, get_bond_subordinate_info, BondInfo,
            BondSubordinateInfo,
        },
        bridge::{
            get_bridge_info, get_bridge_port_info, parse_bridge_vlan_info,
            BridgeConf, BridgeInfo, BridgePortInfo,
        },
        change_ifaces,
        ethtool::EthtoolInfo,
        ipoib::{get_ipoib_info, IpoibInfo},
        mac_vlan::{get_mac_vlan_info, MacVlanInfo},
        mac_vtap::{get_mac_vtap_info, MacVtapInfo},
        sriov::{get_sriov_info, SriovInfo},
        tun::{get_tun_info, TunInfo},
        veth::{VethConf, VethInfo},
        vlan::{get_vlan_info, VlanConf, VlanInfo},
        vrf::{
            get_vrf_info, get_vrf_subordinate_info, VrfInfo, VrfSubordinateInfo,
        },
        vxlan::{get_vxlan_info, VxlanInfo},
    },
    ip::{IpConf, Ipv4Info, Ipv6Info},
    mac::{mac_str_to_raw, parse_as_mac},
    NisporError,
};

use netlink_packet_route::rtnl::{
    link::nlas, LinkMessage, ARPHRD_ETHER, ARPHRD_LOOPBACK, IFF_ALLMULTI,
    IFF_AUTOMEDIA, IFF_BROADCAST, IFF_DEBUG, IFF_DORMANT, IFF_LOOPBACK,
    IFF_LOWER_UP, IFF_MASTER, IFF_MULTICAST, IFF_NOARP, IFF_POINTOPOINT,
    IFF_PORTSEL, IFF_PROMISC, IFF_RUNNING, IFF_SLAVE, IFF_UP,
};
use rtnetlink::packet::rtnl::link::nlas::Nla;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[serde(rename_all = "snake_case")]
#[non_exhaustive]
pub enum IfaceType {
    Bond,
    Veth,
    Bridge,
    Vlan,
    Dummy,
    Vxlan,
    Loopback,
    Ethernet,
    Vrf,
    Tun,
    MacVlan,
    MacVtap,
    OpenvSwitch,
    Ipoib,
    Unknown,
    Other(String),
}

impl Default for IfaceType {
    fn default() -> Self {
        IfaceType::Unknown
    }
}

impl std::fmt::Display for IfaceType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Bond => "bond",
                Self::Veth => "veth",
                Self::Bridge => "bridge",
                Self::Vlan => "vlan",
                Self::Dummy => "dummy",
                Self::Vxlan => "vxlan",
                Self::Loopback => "loopback",
                Self::Ethernet => "ethernet",
                Self::Vrf => "vrf",
                Self::Tun => "tun",
                Self::MacVlan => "macvlan",
                Self::MacVtap => "macvtap",
                Self::OpenvSwitch => "openvswitch",
                Self::Ipoib => "ipoib",
                Self::Unknown => "unknown",
                Self::Other(s) => s,
            }
        )
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[serde(rename_all = "snake_case")]
#[non_exhaustive]
pub enum IfaceState {
    Up,
    Dormant,
    Down,
    LowerLayerDown,
    Absent, // Only for IfaceConf
    Other(String),
    Unknown,
}

impl Default for IfaceState {
    fn default() -> Self {
        IfaceState::Unknown
    }
}

impl std::fmt::Display for IfaceState {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Up => "up",
                Self::Dormant => "dormant",
                Self::Down => "down",
                Self::LowerLayerDown => "lower_layer_down",
                Self::Absent => "absent",
                Self::Other(s) => s.as_str(),
                Self::Unknown => "unknown",
            }
        )
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[serde(rename_all = "snake_case")]
#[non_exhaustive]
pub enum IfaceFlags {
    AllMulti,
    AutoMedia,
    Broadcast,
    Debug,
    Dormant,
    Loopback,
    LowerUp,
    Controller,
    Multicast,
    NoArp,
    PoinToPoint,
    Portsel,
    Promisc,
    Running,
    Subordinate,
    Up,
    Other(u32),
    Unknown,
}

impl Default for IfaceFlags {
    fn default() -> Self {
        Self::Unknown
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[serde(rename_all = "snake_case")]
#[non_exhaustive]
pub enum ControllerType {
    Bond,
    Bridge,
    Vrf,
    OpenvSwitch,
    Other(String),
    Unknown,
}

impl From<&str> for ControllerType {
    fn from(s: &str) -> Self {
        match s {
            "bond" => ControllerType::Bond,
            "bridge" => ControllerType::Bridge,
            "vrf" => ControllerType::Vrf,
            "openvswitch" => ControllerType::OpenvSwitch,
            _ => ControllerType::Other(s.to_string()),
        }
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct Iface {
    pub name: String,
    #[serde(skip_serializing)]
    pub index: u32,
    pub iface_type: IfaceType,
    pub state: IfaceState,
    pub mtu: i64,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub min_mtu: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub max_mtu: Option<i64>,
    pub flags: Vec<IfaceFlags>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ipv4: Option<Ipv4Info>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ipv6: Option<Ipv6Info>,
    #[serde(skip_serializing_if = "String::is_empty")]
    pub mac_address: String,
    #[serde(skip_serializing_if = "String::is_empty")]
    pub permanent_mac_address: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub controller: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub controller_type: Option<ControllerType>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ethtool: Option<EthtoolInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bond: Option<BondInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bond_subordinate: Option<BondSubordinateInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bridge: Option<BridgeInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bridge_port: Option<BridgePortInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tun: Option<TunInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub vlan: Option<VlanInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub vxlan: Option<VxlanInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub veth: Option<VethInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub vrf: Option<VrfInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub vrf_subordinate: Option<VrfSubordinateInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub mac_vlan: Option<MacVlanInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub mac_vtap: Option<MacVtapInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub sriov: Option<SriovInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ipoib: Option<IpoibInfo>,
}

// TODO: impl From Iface to IfaceConf

pub(crate) fn parse_nl_msg_to_name_and_index(
    nl_msg: &LinkMessage,
) -> Option<(String, u32)> {
    let index = nl_msg.header.index;
    let name = _get_iface_name(nl_msg);
    if name.is_empty() {
        None
    } else {
        Some((name, index))
    }
}

pub(crate) fn parse_nl_msg_to_iface(
    nl_msg: &LinkMessage,
) -> Result<Option<Iface>, NisporError> {
    let name = _get_iface_name(nl_msg);
    if name.is_empty() {
        return Ok(None);
    }
    let mut iface_state = Iface {
        name,
        ..Default::default()
    };
    match nl_msg.header.link_layer_type {
        ARPHRD_ETHER => iface_state.iface_type = IfaceType::Ethernet,
        ARPHRD_LOOPBACK => iface_state.iface_type = IfaceType::Loopback,
        _ => (),
    }
    iface_state.index = nl_msg.header.index;
    let mut link: Option<u32> = None;
    let mut mac_len = None;
    for nla in &nl_msg.nlas {
        if let Nla::Mtu(mtu) = nla {
            iface_state.mtu = *mtu as i64;
        } else if let Nla::MinMtu(mtu) = nla {
            iface_state.min_mtu =
                if *mtu != 0 { Some(*mtu as i64) } else { None };
        } else if let Nla::MaxMtu(mtu) = nla {
            iface_state.max_mtu =
                if *mtu != 0 { Some(*mtu as i64) } else { None };
        } else if let Nla::Address(mac) = nla {
            mac_len = Some(mac.len());
            iface_state.mac_address = parse_as_mac(mac.len(), mac)?;
        } else if let Nla::PermAddress(mac) = nla {
            mac_len = Some(mac.len());
            iface_state.permanent_mac_address = parse_as_mac(mac.len(), mac)?;
        } else if let Nla::OperState(state) = nla {
            iface_state.state = _get_iface_state(state);
        } else if let Nla::Master(controller) = nla {
            iface_state.controller = Some(format!("{}", controller));
        } else if let Nla::Link(l) = nla {
            link = Some(*l);
        } else if let Nla::Info(infos) = nla {
            for info in infos {
                if let nlas::Info::Kind(t) = info {
                    iface_state.iface_type = match t {
                        nlas::InfoKind::Bond => IfaceType::Bond,
                        nlas::InfoKind::Veth => IfaceType::Veth,
                        nlas::InfoKind::Bridge => IfaceType::Bridge,
                        nlas::InfoKind::Vlan => IfaceType::Vlan,
                        nlas::InfoKind::Vxlan => IfaceType::Vxlan,
                        nlas::InfoKind::Dummy => IfaceType::Dummy,
                        nlas::InfoKind::Tun => IfaceType::Tun,
                        nlas::InfoKind::Vrf => IfaceType::Vrf,
                        nlas::InfoKind::MacVlan => IfaceType::MacVlan,
                        nlas::InfoKind::MacVtap => IfaceType::MacVtap,
                        nlas::InfoKind::Ipoib => IfaceType::Ipoib,
                        nlas::InfoKind::Other(s) => match s.as_ref() {
                            "openvswitch" => IfaceType::OpenvSwitch,
                            _ => IfaceType::Other(s.clone()),
                        },
                        _ => IfaceType::Other(format!("{:?}", t)),
                    };
                }
            }
            for info in infos {
                if let nlas::Info::Data(d) = info {
                    match iface_state.iface_type {
                        IfaceType::Bond => iface_state.bond = get_bond_info(d)?,
                        IfaceType::Bridge => {
                            iface_state.bridge = get_bridge_info(d)?
                        }
                        IfaceType::Tun => match get_tun_info(d) {
                            Ok(info) => {
                                iface_state.tun = Some(info);
                            }
                            Err(e) => {
                                log::warn!("Error parsing TUN info: {}", e);
                            }
                        },
                        IfaceType::Vlan => iface_state.vlan = get_vlan_info(d),
                        IfaceType::Vxlan => {
                            iface_state.vxlan = get_vxlan_info(d)?
                        }
                        IfaceType::Vrf => iface_state.vrf = get_vrf_info(d),
                        IfaceType::MacVlan => {
                            iface_state.mac_vlan = get_mac_vlan_info(d)?
                        }
                        IfaceType::MacVtap => {
                            iface_state.mac_vtap = get_mac_vtap_info(d)?
                        }
                        IfaceType::Ipoib => {
                            iface_state.ipoib = get_ipoib_info(d);
                        }
                        _ => log::warn!(
                            "Unhandled IFLA_INFO_DATA for iface type {:?}",
                            iface_state.iface_type
                        ),
                    }
                }
            }
            for info in infos {
                if let nlas::Info::SlaveKind(d) = info {
                    // Remove the tailing \0
                    iface_state.controller_type = Some(
                        std::ffi::CStr::from_bytes_with_nul(d.as_slice())?
                            .to_str()?
                            .into(),
                    )
                }
            }
            if let Some(controller_type) = &iface_state.controller_type {
                for info in infos {
                    if let nlas::Info::SlaveData(d) = info {
                        match controller_type {
                            ControllerType::Bond => {
                                iface_state.bond_subordinate =
                                    get_bond_subordinate_info(d)?;
                            }
                            ControllerType::Bridge => {
                                iface_state.bridge_port =
                                    get_bridge_port_info(d)?;
                            }
                            ControllerType::Vrf => {
                                iface_state.vrf_subordinate =
                                    get_vrf_subordinate_info(d)?;
                            }
                            _ => log::warn!(
                                "Unknown controller type {:?}",
                                controller_type
                            ),
                        }
                    }
                }
            }
        } else if let Nla::VfInfoList(data) = nla {
            if let Ok(info) = get_sriov_info(&iface_state.name, data, mac_len) {
                iface_state.sriov = Some(info);
            }
        } else {
            // println!("{} {:?}", name, nla);
        }
    }
    if let Some(ref mut vlan_info) = iface_state.vlan {
        if let Some(base_iface_index) = link {
            vlan_info.base_iface = format!("{}", base_iface_index);
        }
    }
    if let Some(ref mut ib_info) = iface_state.ipoib {
        if let Some(base_iface_index) = link {
            ib_info.base_iface = Some(format!("{}", base_iface_index));
        }
    }
    if let Some(iface_index) = link {
        match iface_state.iface_type {
            IfaceType::Veth => {
                iface_state.veth = Some(VethInfo {
                    peer: format!("{}", iface_index),
                })
            }
            IfaceType::MacVlan => {
                if let Some(ref mut mac_vlan_info) = iface_state.mac_vlan {
                    mac_vlan_info.base_iface = format!("{}", iface_index);
                }
            }
            IfaceType::MacVtap => {
                if let Some(ref mut mac_vtap_info) = iface_state.mac_vtap {
                    mac_vtap_info.base_iface = format!("{}", iface_index);
                }
            }
            _ => (),
        }
    }
    iface_state.flags = _parse_iface_flags(nl_msg.header.flags);
    Ok(Some(iface_state))
}

fn _get_iface_name(nl_msg: &LinkMessage) -> String {
    for nla in &nl_msg.nlas {
        if let Nla::IfName(name) = nla {
            return name.clone();
        }
    }
    "".into()
}

pub(crate) fn fill_bridge_vlan_info(
    iface_states: &mut HashMap<String, Iface>,
    nl_msg: &LinkMessage,
) -> Result<(), NisporError> {
    let name = _get_iface_name(nl_msg);
    if name.is_empty() {
        return Ok(());
    }
    if let Some(iface_state) = iface_states.get_mut(&name) {
        for nla in &nl_msg.nlas {
            if let Nla::AfSpecBridge(data) = nla {
                parse_bridge_vlan_info(iface_state, data)?;
            }
        }
    }
    Ok(())
}

fn _get_iface_state(state: &nlas::State) -> IfaceState {
    match state {
        nlas::State::Up => IfaceState::Up,
        nlas::State::Dormant => IfaceState::Dormant,
        nlas::State::Down => IfaceState::Down,
        nlas::State::LowerLayerDown => IfaceState::LowerLayerDown,
        nlas::State::Unknown => IfaceState::Unknown,
        _ => IfaceState::Other(format!("{:?}", state)),
    }
}

fn _parse_iface_flags(flags: u32) -> Vec<IfaceFlags> {
    let mut ret = Vec::new();
    if (flags & IFF_ALLMULTI) > 0 {
        ret.push(IfaceFlags::AllMulti)
    }
    if (flags & IFF_AUTOMEDIA) > 0 {
        ret.push(IfaceFlags::AutoMedia)
    }
    if (flags & IFF_BROADCAST) > 0 {
        ret.push(IfaceFlags::Broadcast)
    }
    if (flags & IFF_DEBUG) > 0 {
        ret.push(IfaceFlags::Debug)
    }
    if (flags & IFF_DORMANT) > 0 {
        ret.push(IfaceFlags::Dormant)
    }
    if (flags & IFF_LOOPBACK) > 0 {
        ret.push(IfaceFlags::Loopback)
    }
    if (flags & IFF_LOWER_UP) > 0 {
        ret.push(IfaceFlags::LowerUp)
    }
    if (flags & IFF_MASTER) > 0 {
        ret.push(IfaceFlags::Controller)
    }
    if (flags & IFF_MULTICAST) > 0 {
        ret.push(IfaceFlags::Multicast)
    }
    if (flags & IFF_NOARP) > 0 {
        ret.push(IfaceFlags::NoArp)
    }
    if (flags & IFF_POINTOPOINT) > 0 {
        ret.push(IfaceFlags::PoinToPoint)
    }
    if (flags & IFF_PORTSEL) > 0 {
        ret.push(IfaceFlags::Portsel)
    }
    if (flags & IFF_PROMISC) > 0 {
        ret.push(IfaceFlags::Promisc)
    }
    if (flags & IFF_RUNNING) > 0 {
        ret.push(IfaceFlags::Running)
    }
    if (flags & IFF_SLAVE) > 0 {
        ret.push(IfaceFlags::Subordinate)
    }
    if (flags & IFF_UP) > 0 {
        ret.push(IfaceFlags::Up)
    }

    ret
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct IfaceConf {
    pub name: String,
    #[serde(default = "default_iface_state_in_conf")]
    pub state: IfaceState,
    #[serde(rename = "type")]
    pub iface_type: Option<IfaceType>,
    pub controller: Option<String>,
    pub ipv4: Option<IpConf>,
    pub ipv6: Option<IpConf>,
    pub mac_address: Option<String>,
    pub veth: Option<VethConf>,
    pub bridge: Option<BridgeConf>,
    pub vlan: Option<VlanConf>,
}

impl IfaceConf {
    pub async fn apply(&self, cur_iface: &Iface) -> Result<(), NisporError> {
        log::warn!(
            "WARN: IfaceConf::apply() is deprecated, \
            please use NetConf::apply() instead"
        );
        let ifaces = vec![self];
        let mut cur_ifaces = HashMap::new();
        cur_ifaces.insert(self.name.to_string(), cur_iface.clone());
        change_ifaces(&ifaces, &cur_ifaces).await
    }
}

fn default_iface_state_in_conf() -> IfaceState {
    IfaceState::Up
}

pub(crate) async fn change_iface_state(
    handle: &rtnetlink::Handle,
    index: u32,
    up: bool,
) -> Result<(), NisporError> {
    if up {
        handle.link().set(index).up().execute().await?;
    } else {
        handle.link().set(index).down().execute().await?;
    }
    Ok(())
}

pub(crate) async fn change_iface_mac(
    handle: &rtnetlink::Handle,
    index: u32,
    mac_address: &str,
) -> Result<(), NisporError> {
    change_iface_state(handle, index, false).await?;
    handle
        .link()
        .set(index)
        .address(mac_str_to_raw(mac_address)?)
        .execute()
        .await?;
    Ok(())
}
