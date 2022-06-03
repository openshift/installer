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
    convert::mac_str_to_u8_array,
    convert::{own_value_to_bytes_array, u8_array_to_mac_string},
    error::NmError,
};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingWired {
    pub cloned_mac_address: Option<String>,
    pub mtu: Option<u32>,
    pub accept_all_mac_addresses: Option<i32>,
    pub speed: Option<u32>,
    pub duplex: Option<String>,
    pub auto_negotiate: Option<bool>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSettingWired {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            cloned_mac_address: _from_map!(
                v,
                "cloned-mac-address",
                own_value_to_bytes_array
            )?
            .map(u8_array_to_mac_string),
            mtu: _from_map!(v, "mtu", u32::try_from)?,
            accept_all_mac_addresses: _from_map!(
                v,
                "accept-all-mac-addresses",
                i32::try_from
            )?,
            speed: _from_map!(v, "speed", u32::try_from)?,
            duplex: _from_map!(v, "duplex", String::try_from)?,
            auto_negotiate: _from_map!(v, "auto-negotiate", bool::try_from)?,
            _other: v,
        })
    }
}

impl NmSettingWired {
    pub(crate) fn to_keyfile(
        &self,
    ) -> Result<HashMap<String, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        for (k, v) in self.to_value()?.drain() {
            if k != "cloned-mac-address" {
                ret.insert(k.to_string(), v);
            }
        }
        if let Some(v) = &self.cloned_mac_address {
            ret.insert(
                "cloned-mac-address".to_string(),
                zvariant::Value::new(v),
            );
        }
        Ok(ret)
    }

    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        if let Some(v) = &self.cloned_mac_address {
            ret.insert(
                "cloned-mac-address",
                zvariant::Value::new(mac_str_to_u8_array(v)),
            );
        }
        if let Some(v) = &self.mtu {
            ret.insert("mtu", zvariant::Value::new(v));
        }
        if let Some(v) = &self.accept_all_mac_addresses {
            ret.insert("accept-all-mac-addresses", zvariant::Value::new(v));
        }
        if let Some(v) = &self.speed {
            ret.insert("speed", zvariant::Value::new(v));
        }
        if let Some(v) = &self.auto_negotiate {
            ret.insert("auto-negotiate", zvariant::Value::new(v));
        }
        if let Some(v) = &self.duplex {
            ret.insert("duplex", zvariant::Value::new(v));
        }
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }
}
