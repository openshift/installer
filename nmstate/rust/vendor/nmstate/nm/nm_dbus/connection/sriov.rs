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
use std::fmt::Write;

use serde::Deserialize;

use super::super::{
    connection::DbusDictionary,
    dbus::{NM_TERNARY_FALSE, NM_TERNARY_TRUE},
    error::NmError,
    NmVlanProtocol,
};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingSriov {
    pub autoprobe_drivers: Option<bool>,
    pub total_vfs: Option<u32>,
    pub vfs: Option<Vec<NmSettingSriovVf>>,
    _other: DbusDictionary,
}

impl TryFrom<DbusDictionary> for NmSettingSriov {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            autoprobe_drivers: match _from_map!(
                v,
                "autoprobe-drivers",
                i32::try_from
            )? {
                Some(NM_TERNARY_TRUE) => Some(true),
                Some(NM_TERNARY_FALSE) => Some(false),
                _ => None,
            },
            total_vfs: _from_map!(v, "total-vfs", u32::try_from)?,
            vfs: _from_map!(v, "vfs", own_value_to_vfs)?,
            _other: v,
        })
    }
}

impl NmSettingSriov {
    pub(crate) fn to_keyfile(
        &self,
    ) -> Result<HashMap<String, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        for (k, v) in self.to_value()?.drain() {
            if k != "vfs" {
                ret.insert(k.to_string(), v);
            }
        }
        if let Some(vfs) = self.vfs.as_ref() {
            for vf in vfs {
                if let Some(i) = vf.index {
                    ret.insert(
                        format!("vf.{}", i),
                        zvariant::Value::new(vf.to_keyfile()),
                    );
                }
            }
        }
        Ok(ret)
    }

    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        if let Some(v) = &self.autoprobe_drivers {
            ret.insert(
                "autoprobe-drivers",
                zvariant::Value::new(match v {
                    true => NM_TERNARY_TRUE,
                    false => NM_TERNARY_FALSE,
                }),
            );
        }
        if let Some(v) = &self.total_vfs {
            ret.insert("total-vfs", zvariant::Value::new(v));
        }
        if let Some(vfs) = self.vfs.as_ref() {
            let mut vf_values = zvariant::Array::new(
                zvariant::Signature::from_str_unchecked("a{sv}"),
            );
            for vf in vfs {
                vf_values.append(vf.to_value()?)?;
            }
            ret.insert("vfs", zvariant::Value::Array(vf_values));
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
pub struct NmSettingSriovVf {
    pub index: Option<u32>,
    pub mac: Option<String>,
    pub spoof_check: Option<bool>,
    pub trust: Option<bool>,
    pub min_tx_rate: Option<u32>,
    pub max_tx_rate: Option<u32>,
    pub vlans: Option<Vec<NmSettingSriovVfVlan>>,
    _other: DbusDictionary,
}

impl TryFrom<DbusDictionary> for NmSettingSriovVf {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            index: _from_map!(v, "index", u32::try_from)?,
            mac: _from_map!(v, "mac", String::try_from)?,
            spoof_check: _from_map!(v, "spoof-check", bool::try_from)?,
            trust: _from_map!(v, "trust", bool::try_from)?,
            min_tx_rate: _from_map!(v, "min-tx-rate", u32::try_from)?,
            max_tx_rate: _from_map!(v, "max-tx-rate", u32::try_from)?,
            vlans: _from_map!(v, "vlans", own_value_to_vf_vlans)?,
            _other: v,
        })
    }
}

impl NmSettingSriovVf {
    pub(crate) fn to_keyfile(&self) -> String {
        let mut ret = String::new();
        if let Some(v) = self.mac.as_ref() {
            let _ = write!(ret, "mac={} ", v);
        }
        if let Some(v) = self.spoof_check {
            let _ = write!(ret, "spoof-check={} ", v);
        }
        if let Some(v) = self.trust {
            let _ = write!(ret, "trust={} ", v);
        }
        if let Some(v) = self.min_tx_rate {
            let _ = write!(ret, "min-tx-rate={} ", v);
        }
        if let Some(v) = self.max_tx_rate {
            let _ = write!(ret, "max-tx-rate={} ", v);
        }
        if let Some(vlans) = self.vlans.as_ref() {
            let mut vlans_str = Vec::new();
            for vlan in vlans {
                vlans_str.push(vlan.to_keyfile());
            }
            let _ = write!(ret, "vlans={}", vlans_str.join(";"));
        }
        if ret.ends_with(' ') {
            ret.pop();
        }
        ret
    }

    pub(crate) fn to_value(&self) -> Result<zvariant::Value, NmError> {
        let mut ret = zvariant::Dict::new(
            zvariant::Signature::from_str_unchecked("s"),
            zvariant::Signature::from_str_unchecked("v"),
        );

        if let Some(v) = &self.index {
            ret.append(
                zvariant::Value::new("index"),
                zvariant::Value::new(zvariant::Value::new(v)),
            )?;
        }
        if let Some(v) = &self.mac {
            ret.append(
                zvariant::Value::new("mac"),
                zvariant::Value::new(zvariant::Value::new(v)),
            )?;
        }
        if let Some(v) = &self.spoof_check {
            ret.append(
                zvariant::Value::new("spoof-check"),
                zvariant::Value::new(zvariant::Value::new(v)),
            )?;
        }
        if let Some(v) = &self.trust {
            ret.append(
                zvariant::Value::new("trust"),
                zvariant::Value::new(zvariant::Value::new(v)),
            )?;
        }
        if let Some(v) = &self.min_tx_rate {
            ret.append(
                zvariant::Value::new("min-tx-rate"),
                zvariant::Value::new(zvariant::Value::new(v)),
            )?;
        }
        if let Some(v) = &self.max_tx_rate {
            ret.append(
                zvariant::Value::new("max-tx-rate"),
                zvariant::Value::new(zvariant::Value::new(v)),
            )?;
        }
        if let Some(vlans) = self.vlans.as_ref() {
            let mut vlan_values = zvariant::Array::new(
                zvariant::Signature::from_str_unchecked("a{sv}"),
            );
            for vlan in vlans {
                vlan_values.append(vlan.to_value()?)?;
            }
            ret.append(
                zvariant::Value::new("vlans"),
                zvariant::Value::new(zvariant::Value::Array(vlan_values)),
            )?;
        }
        for (key, value) in self._other.iter() {
            ret.append(
                zvariant::Value::new(key.as_str()),
                zvariant::Value::from(value.clone()),
            )?;
        }

        Ok(zvariant::Value::Dict(ret))
    }
}

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingSriovVfVlan {
    pub id: u32,
    pub qos: u32,
    pub protocol: NmVlanProtocol,
    _other: DbusDictionary,
}

impl TryFrom<DbusDictionary> for NmSettingSriovVfVlan {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            id: _from_map!(v, "id", u32::try_from)?.unwrap_or_default(),
            qos: _from_map!(v, "qos", u32::try_from)?.unwrap_or_default(),
            protocol: match _from_map!(v, "protocol", u32::try_from)? {
                Some(1) => NmVlanProtocol::Dot1Ad,
                _ => NmVlanProtocol::Dot1Q,
            },
            _other: v,
        })
    }
}

impl NmSettingSriovVfVlan {
    pub(crate) fn to_keyfile(&self) -> String {
        match self.protocol {
            NmVlanProtocol::Dot1Q => format!("{}.{}.q", self.id, self.qos),
            NmVlanProtocol::Dot1Ad => format!("{}.{}.ad", self.id, self.qos),
        }
    }

    pub(crate) fn to_value(&self) -> Result<zvariant::Value, NmError> {
        let mut ret = zvariant::Dict::new(
            zvariant::Signature::from_str_unchecked("s"),
            zvariant::Signature::from_str_unchecked("v"),
        );
        ret.append(
            zvariant::Value::new("id"),
            zvariant::Value::new(zvariant::Value::U32(self.id)),
        )?;
        ret.append(
            zvariant::Value::new("qos"),
            zvariant::Value::new(zvariant::Value::U32(self.qos)),
        )?;
        ret.append(
            zvariant::Value::new("protocol"),
            zvariant::Value::new(zvariant::Value::U32(match self.protocol {
                NmVlanProtocol::Dot1Q => 0,
                NmVlanProtocol::Dot1Ad => 1,
            })),
        )?;
        Ok(zvariant::Value::Dict(ret))
    }
}

fn own_value_to_vfs(
    value: zvariant::OwnedValue,
) -> Result<Vec<NmSettingSriovVf>, NmError> {
    let mut ret = Vec::new();
    let raw_vfs = Vec::<DbusDictionary>::try_from(value)?;
    for raw_vf in raw_vfs {
        ret.push(NmSettingSriovVf::try_from(raw_vf)?);
    }
    Ok(ret)
}

fn own_value_to_vf_vlans(
    value: zvariant::OwnedValue,
) -> Result<Vec<NmSettingSriovVfVlan>, NmError> {
    let mut ret = Vec::new();
    let raw_vfs = Vec::<DbusDictionary>::try_from(value)?;
    for raw_vf in raw_vfs {
        ret.push(NmSettingSriovVfVlan::try_from(raw_vf)?);
    }
    Ok(ret)
}
