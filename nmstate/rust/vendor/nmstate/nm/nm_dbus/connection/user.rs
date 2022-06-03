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

use super::super::{connection::DbusDictionary, error::NmError};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingUser {
    pub data: Option<HashMap<String, String>>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSettingUser {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            data: _from_map!(v, "data", <HashMap<String, String>>::try_from)?,
            _other: v,
        })
    }
}

impl NmSettingUser {
    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        if let Some(v) = &self.data {
            let mut dict_value = zvariant::Dict::new(
                zvariant::Signature::from_str_unchecked("s"),
                zvariant::Signature::from_str_unchecked("s"),
            );
            for (k, v) in v.iter() {
                dict_value
                    .append(zvariant::Value::new(k), zvariant::Value::new(v))?;
            }
            ret.insert("data", zvariant::Value::Dict(dict_value));
        }
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }

    pub(crate) fn to_keyfile(
        &self,
    ) -> Result<HashMap<String, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        if let Some(data) = self.data.as_ref() {
            for (k, v) in data.iter() {
                ret.insert(k.to_string(), zvariant::Value::new(v));
            }
        }
        Ok(ret)
    }
}
