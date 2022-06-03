use std::collections::HashMap;
use std::convert::TryFrom;

use log::error;
use serde::Deserialize;

use super::super::{connection::DbusDictionary, ErrorKind, NmError};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingVlan {
    pub parent: Option<String>,
    pub id: Option<u32>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSettingVlan {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            parent: _from_map!(v, "parent", String::try_from)?,
            id: _from_map!(v, "id", u32::try_from)?,
            _other: v,
        })
    }
}

impl NmSettingVlan {
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
        if let Some(id) = self.id {
            ret.insert("id", zvariant::Value::new(id));
        }
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }
}

const NM_VLAN_PROTOCOL_802_1Q: &str = "802.1Q";
const NM_VLAN_PROTOCOL_802_1AD: &str = "802.1AD";

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
#[non_exhaustive]
pub enum NmVlanProtocol {
    Dot1Q,
    Dot1Ad,
}

impl Default for NmVlanProtocol {
    fn default() -> Self {
        Self::Dot1Q
    }
}

impl TryFrom<String> for NmVlanProtocol {
    type Error = NmError;
    fn try_from(vlan_protocol: String) -> Result<Self, Self::Error> {
        match vlan_protocol.as_str() {
            NM_VLAN_PROTOCOL_802_1Q => Ok(Self::Dot1Q),
            NM_VLAN_PROTOCOL_802_1AD => Ok(Self::Dot1Ad),
            _ => {
                let e = NmError::new(
                    ErrorKind::InvalidArgument,
                    format!(
                        "Invalid VLAN protocol {}, only support: {} and {}",
                        vlan_protocol,
                        NM_VLAN_PROTOCOL_802_1Q,
                        NM_VLAN_PROTOCOL_802_1AD
                    ),
                );
                error!("{}", e);
                Err(e)
            }
        }
    }
}

impl NmVlanProtocol {
    pub fn to_str(self) -> &'static str {
        match self {
            Self::Dot1Q => NM_VLAN_PROTOCOL_802_1Q,
            Self::Dot1Ad => NM_VLAN_PROTOCOL_802_1AD,
        }
    }
}
