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

mod error;
mod ifaces;
mod ip;
mod mac;
mod net_conf;
mod net_state;
mod netlink;
mod route;
mod route_rule;

pub use crate::error::NisporError;
pub use crate::ifaces::{
    BondAdInfo, BondAdSelect, BondAllSubordinatesActive, BondArpValidate,
    BondFailOverMac, BondInfo, BondLacpRate, BondMiiStatus, BondMode,
    BondModeArpAllTargets, BondPrimaryReselect, BondSubordinateInfo,
    BondSubordinateState, BondXmitHashPolicy, BridgeInfo, BridgePortInfo,
    BridgePortMulticastRouterType, BridgePortStpState, BridgeStpState,
    BridgeVlanEntry, BridgeVlanProtocol, ControllerType, EthtoolCoalesceInfo,
    EthtoolFeatureInfo, EthtoolInfo, EthtoolLinkModeDuplex,
    EthtoolLinkModeInfo, EthtoolPauseInfo, EthtoolRingInfo, Iface, IfaceConf,
    IfaceFlags, IfaceState, IfaceType, IpoibInfo, IpoibMode, MacVlanInfo,
    MacVlanMode, MacVtapInfo, MacVtapMode, SriovInfo, TunInfo, TunMode,
    VethConf, VethInfo, VfInfo, VfLinkState, VfState, VlanConf, VlanInfo,
    VlanProtocol, VrfInfo, VrfSubordinateInfo, VxlanInfo,
};
pub use crate::ip::{
    IpAddrConf, IpConf, IpFamily, Ipv4AddrInfo, Ipv4Info, Ipv6AddrInfo,
    Ipv6Info,
};
pub use crate::net_conf::NetConf;
pub use crate::net_state::NetState;
pub use crate::route::{
    AddressFamily, MultipathRoute, MultipathRouteFlags, Route, RouteConf,
    RouteProtocol, RouteScope, RouteType,
};
pub use crate::route_rule::{RouteRule, RuleAction};
