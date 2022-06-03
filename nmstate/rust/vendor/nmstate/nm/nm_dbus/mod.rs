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

mod active_connection;
mod connection;
mod convert;
mod dbus;
mod dbus_proxy;
mod device;
mod dns;
mod error;
mod keyfile;
mod lldp;
mod nm_api;

pub use self::active_connection::NmActiveConnection;
pub use self::connection::{
    NmConnection, NmIpRoute, NmIpRouteRule, NmSetting8021X, NmSettingBond,
    NmSettingBridge, NmSettingBridgeVlanRange, NmSettingConnection,
    NmSettingEthtool, NmSettingInfiniBand, NmSettingIp, NmSettingIpMethod,
    NmSettingMacVlan, NmSettingOvsBridge, NmSettingOvsDpdk, NmSettingOvsExtIds,
    NmSettingOvsIface, NmSettingOvsPatch, NmSettingOvsPort, NmSettingSriov,
    NmSettingSriovVf, NmSettingSriovVfVlan, NmSettingUser, NmSettingVeth,
    NmSettingVlan, NmSettingVrf, NmSettingVxlan, NmSettingWired,
    NmSettingsConnectionFlag, NmVlanProtocol,
};
pub use self::device::{NmDevice, NmDeviceState, NmDeviceStateReason};
pub use self::dns::NmDnsEntry;
pub use self::error::{ErrorKind, NmError};
pub use self::lldp::{
    NmLldpNeighbor, NmLldpNeighbor8021Ppvid, NmLldpNeighbor8021Vlan,
    NmLldpNeighbor8023MacPhyConf, NmLldpNeighbor8023PowerViaMdi,
    NmLldpNeighborMgmtAddr,
};
pub use self::nm_api::NmApi;
