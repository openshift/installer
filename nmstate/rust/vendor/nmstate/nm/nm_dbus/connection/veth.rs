use std::collections::HashMap;
use std::convert::TryFrom;

use serde::Deserialize;

use super::super::{connection::DbusDictionary, NmError};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingVeth {
    pub peer: Option<String>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSettingVeth {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            peer: _from_map!(v, "peer", String::try_from)?,
            _other: v,
        })
    }
}

impl NmSettingVeth {
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
        if let Some(v) = &self.peer {
            ret.insert("peer", zvariant::Value::new(v.clone()));
        }
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }
}
