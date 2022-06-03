use std::collections::HashMap;
use std::convert::TryFrom;

use serde::Deserialize;

use super::super::{connection::DbusDictionary, NmError};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingInfiniBand {
    pub parent: Option<String>,
    pub pkey: Option<i32>,
    pub mode: Option<String>,
    pub mtu: Option<u32>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSettingInfiniBand {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            parent: _from_map!(v, "parent", String::try_from)?,
            pkey: _from_map!(v, "p-key", i32::try_from)?,
            mode: _from_map!(v, "transport-mode", String::try_from)?,
            mtu: _from_map!(v, "mtu", u32::try_from)?,
            _other: v,
        })
    }
}

impl NmSettingInfiniBand {
    pub(crate) fn to_keyfile(
        &self,
    ) -> Result<HashMap<String, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        for (k, v) in self.to_value()?.drain() {
            ret.insert(k.to_string(), v);
        }
        Ok(ret)
    }

    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        if let Some(v) = &self.parent {
            ret.insert("parent", zvariant::Value::new(v.clone()));
        }
        if let Some(v) = &self.mode {
            ret.insert("transport-mode", zvariant::Value::new(v));
        }
        if let Some(v) = &self.pkey {
            ret.insert("p-key", zvariant::Value::new(v));
        }
        if let Some(v) = &self.mtu {
            ret.insert("mtu", zvariant::Value::new(v));
        }
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }
}
