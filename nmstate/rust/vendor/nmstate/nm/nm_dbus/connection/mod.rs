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
//

#[macro_use]
mod macros;

mod bond;
mod bridge;
mod conn;
mod dns;
mod ethtool;
mod ieee8021x;
mod infiniband;
mod ip;
mod mac_vlan;
mod ovs;
mod route;
mod route_rule;
mod sriov;
mod user;
mod veth;
mod vlan;
mod vrf;
mod vxlan;
mod wired;

pub use self::bond::NmSettingBond;
pub use self::bridge::{
    NmSettingBridge, NmSettingBridgePort, NmSettingBridgeVlanRange,
};
pub use self::conn::{
    NmConnection, NmSettingConnection, NmSettingsConnectionFlag,
};
pub use self::ethtool::NmSettingEthtool;
pub use self::ieee8021x::NmSetting8021X;
pub use self::infiniband::NmSettingInfiniBand;
pub use self::ip::{NmSettingIp, NmSettingIpMethod};
pub use self::mac_vlan::NmSettingMacVlan;
pub use self::ovs::{
    NmSettingOvsBridge, NmSettingOvsDpdk, NmSettingOvsExtIds,
    NmSettingOvsIface, NmSettingOvsPatch, NmSettingOvsPort,
};
pub use self::route::NmIpRoute;
pub use self::route_rule::NmIpRouteRule;
pub use self::sriov::{NmSettingSriov, NmSettingSriovVf, NmSettingSriovVfVlan};
pub use self::user::NmSettingUser;
pub use self::veth::NmSettingVeth;
pub use self::vlan::{NmSettingVlan, NmVlanProtocol};
pub use self::vrf::NmSettingVrf;
pub use self::vxlan::NmSettingVxlan;
pub use self::wired::NmSettingWired;

pub(crate) use self::conn::{
    nm_con_get_from_obj_path, DbusDictionary, NmConnectionDbusValue,
};
pub(crate) use self::macros::_from_map;
