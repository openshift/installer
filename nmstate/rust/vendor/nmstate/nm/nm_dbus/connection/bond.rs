use std::collections::HashMap;
use std::convert::TryFrom;

use serde::Deserialize;
use zvariant::Value;

use super::super::{connection::DbusDictionary, NmError};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingBond {
    pub options: HashMap<String, String>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSettingBond {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            options: _from_map!(
                v,
                "options",
                <HashMap<String, String>>::try_from
            )?
            .unwrap_or_default(),
            _other: v,
        })
    }
}

impl NmSettingBond {
    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        ret.insert("options", zvariant::Value::from(self.options.clone()));
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }

    pub(crate) fn to_keyfile(
        &self,
    ) -> Result<HashMap<String, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        for (key, value) in self.options.iter() {
            ret.insert(key.to_string(), Value::new(value));
        }
        Ok(ret)
    }
}
