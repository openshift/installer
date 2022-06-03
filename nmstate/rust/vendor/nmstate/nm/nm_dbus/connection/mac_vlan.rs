use std::collections::HashMap;
use std::convert::TryFrom;

use serde::Deserialize;

use super::super::{connection::DbusDictionary, NmError};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingMacVlan {
    pub parent: Option<String>,
    pub mode: Option<u32>,
    pub accept_all_mac: Option<bool>,
    pub tap: Option<bool>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSettingMacVlan {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            parent: _from_map!(v, "parent", String::try_from)?,
            mode: _from_map!(v, "mode", u32::try_from)?,
            accept_all_mac: _from_map!(v, "promiscuous", bool::try_from)?,
            tap: _from_map!(v, "tap", bool::try_from)?,
            _other: v,
        })
    }
}

impl NmSettingMacVlan {
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
        if let Some(v) = self.mode {
            ret.insert("mode", zvariant::Value::new(v));
        }
        if let Some(v) = self.accept_all_mac {
            ret.insert("promiscuous", zvariant::Value::new(v));
        }
        if let Some(v) = self.tap {
            ret.insert("tap", zvariant::Value::new(v));
        }
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }
}
