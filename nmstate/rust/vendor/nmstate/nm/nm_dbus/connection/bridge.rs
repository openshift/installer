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

use std::collections::HashMap;
use std::convert::TryFrom;

use serde::Deserialize;

use super::super::{
    connection::DbusDictionary,
    convert::{
        mac_str_to_u8_array, own_value_to_bytes_array, u8_array_to_mac_string,
    },
    NmError, NmVlanProtocol,
};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingBridge {
    pub ageing_time: Option<u32>,
    pub forward_delay: Option<u32>,
    pub group_address: Option<String>,
    pub group_forward_mask: Option<u32>,
    pub hello_time: Option<u32>,
    pub max_age: Option<u32>,
    pub multicast_hash_max: Option<u32>,
    pub multicast_last_member_count: Option<u32>,
    pub multicast_last_member_interval: Option<u64>,
    pub multicast_membership_interval: Option<u64>,
    pub multicast_querier: Option<bool>,
    pub multicast_querier_interval: Option<u64>,
    pub multicast_query_interval: Option<u64>,
    pub multicast_query_response_interval: Option<u64>,
    pub multicast_query_use_ifaddr: Option<bool>,
    pub multicast_router: Option<String>,
    pub multicast_snooping: Option<bool>,
    pub multicast_startup_query_count: Option<u32>,
    pub multicast_startup_query_interval: Option<u64>,
    pub priority: Option<u32>,
    pub stp: Option<bool>,
    pub vlan_default_pvid: Option<u32>,
    pub vlan_filtering: Option<bool>,
    pub vlan_protocol: Option<NmVlanProtocol>,
    pub vlan_stats_enabled: Option<bool>,
    pub vlans: Option<Vec<NmSettingBridgeVlanRange>>,
    _other: DbusDictionary,
}

impl TryFrom<DbusDictionary> for NmSettingBridge {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        // The 'interface-name' is deprecated in favor of
        // `connection.interface-name`.
        v.remove("interface-name");

        Ok(Self {
            ageing_time: _from_map!(v, "ageing-time", u32::try_from)?,
            forward_delay: _from_map!(v, "forward-delay", u32::try_from)?,
            group_address: _from_map!(
                v,
                "group-address",
                own_value_to_bytes_array
            )?
            .map(u8_array_to_mac_string),
            group_forward_mask: _from_map!(
                v,
                "group-forward-mask",
                u32::try_from
            )?,
            hello_time: _from_map!(v, "hello-time", u32::try_from)?,
            max_age: _from_map!(v, "max-age", u32::try_from)?,
            multicast_hash_max: _from_map!(
                v,
                "multicast-hash-max",
                u32::try_from
            )?,
            multicast_last_member_count: _from_map!(
                v,
                "multicast-last-member-count",
                u32::try_from
            )?,
            multicast_last_member_interval: _from_map!(
                v,
                "multicast-last-member-interval",
                u64::try_from
            )?,
            multicast_membership_interval: _from_map!(
                v,
                "multicast-membership-interval",
                u64::try_from
            )?,
            multicast_querier: _from_map!(
                v,
                "multicast-querier",
                bool::try_from
            )?,
            multicast_querier_interval: _from_map!(
                v,
                "multicast-querier-interval",
                u64::try_from
            )?,
            multicast_query_interval: _from_map!(
                v,
                "multicast-query-interval",
                u64::try_from
            )?,
            multicast_query_response_interval: _from_map!(
                v,
                "multicast-query-response-interval",
                u64::try_from
            )?,
            multicast_query_use_ifaddr: _from_map!(
                v,
                "multicast-query-use-ifaddr",
                bool::try_from
            )?,
            multicast_router: _from_map!(
                v,
                "multicast-router",
                String::try_from
            )?,
            multicast_snooping: _from_map!(
                v,
                "multicast-snooping",
                bool::try_from
            )?,
            multicast_startup_query_count: _from_map!(
                v,
                "multicast-startup-query-count",
                u32::try_from
            )?,
            multicast_startup_query_interval: _from_map!(
                v,
                "multicast-startup-query-interval",
                u64::try_from
            )?,
            priority: _from_map!(v, "priority", u32::try_from)?,
            // Default value of STP is True
            stp: _from_map!(v, "stp", bool::try_from)?.or(Some(true)),
            vlan_default_pvid: _from_map!(
                v,
                "vlan-default-pvid",
                u32::try_from
            )?,
            vlan_filtering: _from_map!(v, "vlan-filtering", bool::try_from)?,
            vlan_protocol: _from_map!(v, "vlan-protocol", String::try_from)?
                .map(NmVlanProtocol::try_from)
                .transpose()?,
            vlan_stats_enabled: _from_map!(
                v,
                "vlan-stats-enabled",
                bool::try_from
            )?,
            vlans: _from_map!(v, "vlans", own_value_to_vlan_ranges)?,
            _other: v,
        })
    }
}

impl NmSettingBridge {
    pub(crate) fn to_keyfile(
        &self,
    ) -> Result<HashMap<String, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();

        for (k, v) in self.to_value()?.drain() {
            if k != "vlans" {
                ret.insert(k.to_string(), v);
            }
        }
        if let Some(vlans) = self.vlans.as_ref() {
            let mut vlans_clone = vlans.clone();
            vlans_clone.sort_unstable_by_key(|v| v.vid_start);
            let mut vlans_str = Vec::new();
            for vlan in vlans_clone {
                vlans_str.push(vlan.to_keyfile());
            }
            ret.insert(
                "vlans".to_string(),
                zvariant::Value::new(vlans_str.join(",")),
            );
        }

        Ok(ret)
    }

    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        if let Some(v) = &self.stp {
            ret.insert("stp", zvariant::Value::new(v));
        }
        if let Some(v) = &self.ageing_time {
            ret.insert("ageing-time", zvariant::Value::new(v));
        }
        if let Some(v) = &self.forward_delay {
            ret.insert("forward-delay", zvariant::Value::new(v));
        }
        if let Some(v) = &self.group_address {
            ret.insert(
                "group-address",
                zvariant::Value::new(mac_str_to_u8_array(v)),
            );
        }
        if let Some(v) = &self.group_forward_mask {
            ret.insert("group-forward-mask", zvariant::Value::new(v));
        }
        if let Some(v) = &self.hello_time {
            ret.insert("hello-time", zvariant::Value::new(v));
        }
        if let Some(v) = &self.max_age {
            ret.insert("max-age", zvariant::Value::new(v));
        }
        if let Some(v) = &self.multicast_hash_max {
            ret.insert("multicast-hash-max", zvariant::Value::new(v));
        }
        if let Some(v) = &self.multicast_last_member_count {
            ret.insert("multicast-last-member-count", zvariant::Value::new(v));
        }
        if let Some(v) = &self.multicast_last_member_interval {
            ret.insert(
                "multicast-last-member-interval",
                zvariant::Value::new(v),
            );
        }
        if let Some(v) = &self.multicast_membership_interval {
            ret.insert(
                "multicast-membership-interval",
                zvariant::Value::new(v),
            );
        }
        if let Some(v) = &self.multicast_querier {
            ret.insert("multicast-querier", zvariant::Value::new(v));
        }
        if let Some(v) = &self.multicast_querier_interval {
            ret.insert("multicast-querier-interval", zvariant::Value::new(v));
        }
        if let Some(v) = &self.multicast_query_interval {
            ret.insert("multicast-query-interval", zvariant::Value::new(v));
        }
        if let Some(v) = &self.multicast_query_response_interval {
            ret.insert(
                "multicast-query-response-interval",
                zvariant::Value::new(v),
            );
        }
        if let Some(v) = &self.multicast_query_use_ifaddr {
            ret.insert("multicast-query-use-ifaddr", zvariant::Value::new(v));
        }
        if let Some(v) = &self.multicast_router {
            ret.insert("multicast-router", zvariant::Value::new(v));
        }
        if let Some(v) = &self.multicast_snooping {
            ret.insert("multicast-snooping", zvariant::Value::new(v));
        }
        if let Some(v) = &self.multicast_startup_query_count {
            ret.insert(
                "multicast-startup-query-count",
                zvariant::Value::new(v),
            );
        }
        if let Some(v) = &self.multicast_startup_query_interval {
            ret.insert(
                "multicast-startup-query-interval",
                zvariant::Value::new(v),
            );
        }
        if let Some(v) = &self.priority {
            ret.insert("priority", zvariant::Value::new(v));
        }
        // Default value of STP is True
        if let Some(v) = &self.stp {
            ret.insert("stp", zvariant::Value::new(v));
        }
        if let Some(v) = &self.vlan_default_pvid {
            ret.insert("vlan-default-pvid", zvariant::Value::new(v));
        }
        if let Some(v) = &self.vlan_filtering {
            ret.insert("vlan-filtering", zvariant::Value::new(v));
        }
        if let Some(v) = &self.vlan_protocol {
            ret.insert("vlan-protocol", zvariant::Value::new(v.to_str()));
        }
        if let Some(v) = &self.vlan_stats_enabled {
            ret.insert("vlan-stats-enabled", zvariant::Value::new(v));
        }
        if let Some(vlans) = &self.vlans {
            let mut vlan_values = zvariant::Array::new(
                zvariant::Signature::from_str_unchecked("a{sv}"),
            );
            for vlan in vlans {
                vlan_values.append(vlan.to_value()?)?;
            }
            ret.insert("vlans", zvariant::Value::Array(vlan_values));
        }
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }
}

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingBridgeVlanRange {
    pub vid_start: u16,
    pub vid_end: u16,
    pub pvid: bool,
    pub untagged: bool,
    _other: DbusDictionary,
}

impl TryFrom<DbusDictionary> for NmSettingBridgeVlanRange {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            vid_start: _from_map!(v, "vid-start", u16::try_from)?
                .unwrap_or_default(),
            vid_end: _from_map!(v, "vid-end", u16::try_from)?
                .unwrap_or_default(),
            pvid: _from_map!(v, "pvid", bool::try_from)?.unwrap_or_default(),
            untagged: _from_map!(v, "untagged", bool::try_from)?
                .unwrap_or_default(),
            _other: v,
        })
    }
}

impl NmSettingBridgeVlanRange {
    pub(crate) fn to_keyfile(&self) -> String {
        let mut ret = if self.vid_start == self.vid_end {
            self.vid_start.to_string()
        } else {
            format!("{}-{}", self.vid_start, self.vid_end)
        };
        if self.pvid {
            ret += " pvid"
        }
        if self.untagged {
            ret += " untagged"
        }
        ret
    }

    pub fn to_value(&self) -> Result<zvariant::Value, NmError> {
        let mut ret = zvariant::Dict::new(
            zvariant::Signature::from_str_unchecked("s"),
            zvariant::Signature::from_str_unchecked("v"),
        );
        ret.append(
            zvariant::Value::new("vid-start"),
            zvariant::Value::new(zvariant::Value::U16(self.vid_start)),
        )?;
        ret.append(
            zvariant::Value::new("vid-end"),
            zvariant::Value::new(zvariant::Value::U16(self.vid_end)),
        )?;
        ret.append(
            zvariant::Value::new("pvid"),
            zvariant::Value::new(zvariant::Value::Bool(self.pvid)),
        )?;
        ret.append(
            zvariant::Value::new("untagged"),
            zvariant::Value::new(zvariant::Value::Bool(self.untagged)),
        )?;
        Ok(zvariant::Value::Dict(ret))
    }
}

#[derive(Debug, Clone, PartialEq, Default)]
#[non_exhaustive]
pub struct NmSettingBridgePort {
    pub hairpin_mode: Option<bool>,
    pub path_cost: Option<u32>,
    pub priority: Option<u32>,
    pub vlans: Option<Vec<NmSettingBridgeVlanRange>>,
    _other: DbusDictionary,
}

impl TryFrom<DbusDictionary> for NmSettingBridgePort {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            hairpin_mode: _from_map!(v, "hairpin_mode", bool::try_from)?,
            path_cost: _from_map!(v, "path-cost", u32::try_from)?,
            priority: _from_map!(v, "priority", u32::try_from)?,
            vlans: _from_map!(v, "vlans", own_value_to_vlan_ranges)?,
            _other: v,
        })
    }
}

impl NmSettingBridgePort {
    pub(crate) fn to_keyfile(
        &self,
    ) -> Result<HashMap<String, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();

        for (k, v) in self.to_value()?.drain() {
            if k != "vlans" {
                ret.insert(k.to_string(), v);
            }
        }
        if let Some(vlans) = self.vlans.as_ref() {
            let mut vlans_clone = vlans.clone();
            vlans_clone.sort_unstable_by_key(|v| v.vid_start);
            let mut vlans_str = Vec::new();
            for vlan in vlans_clone {
                vlans_str.push(vlan.to_keyfile());
            }
            ret.insert(
                "vlans".to_string(),
                zvariant::Value::new(vlans_str.join(",")),
            );
        }

        Ok(ret)
    }

    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();

        self.hairpin_mode
            .map(|v| ret.insert("hairpin-mode", zvariant::Value::new(v)));
        self.path_cost
            .map(|v| ret.insert("path-cost", zvariant::Value::new(v)));
        self.priority
            .map(|v| ret.insert("priority", zvariant::Value::new(v)));

        if let Some(vlans) = self.vlans.as_ref() {
            let mut vlan_values = zvariant::Array::new(
                zvariant::Signature::from_str_unchecked("a{sv}"),
            );
            for vlan in vlans {
                vlan_values.append(vlan.to_value()?)?;
            }
            ret.insert("vlans", zvariant::Value::Array(vlan_values));
        }
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }
}

fn own_value_to_vlan_ranges(
    value: zvariant::OwnedValue,
) -> Result<Vec<NmSettingBridgeVlanRange>, NmError> {
    let mut ret = Vec::new();
    let raw_vlan_ranges = Vec::<DbusDictionary>::try_from(value)?;
    for raw_vlan_range in raw_vlan_ranges {
        ret.push(NmSettingBridgeVlanRange::try_from(raw_vlan_range)?);
    }
    Ok(ret)
}
