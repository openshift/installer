use std::collections::HashMap;
use std::convert::TryFrom;

use serde::Deserialize;

use super::super::{connection::DbusDictionary, ErrorKind, NmError};

const GLIB_FILE_PATH_PREFIX: &str = "file://";

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSetting8021X {
    pub identity: Option<String>,
    pub private_key: Option<Vec<u8>>,
    pub eap: Option<Vec<String>>,
    pub client_cert: Option<Vec<u8>>,
    pub ca_cert: Option<Vec<u8>>,
    pub private_key_password: Option<String>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSetting8021X {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            identity: _from_map!(v, "identity", String::try_from)?,
            private_key: _from_map!(v, "private-key", <Vec<u8>>::try_from)?,
            eap: _from_map!(v, "eap", <Vec<String>>::try_from)?,
            client_cert: _from_map!(v, "client-cert", <Vec<u8>>::try_from)?,
            ca_cert: _from_map!(v, "ca-cert", <Vec<u8>>::try_from)?,
            private_key_password: None,
            _other: v,
        })
    }
}

impl NmSetting8021X {
    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        if let Some(v) = &self.identity {
            ret.insert("identity", zvariant::Value::new(v));
        }
        if let Some(v) = &self.private_key {
            ret.insert("private-key", zvariant::Value::new(v));
        }
        if let Some(v) = &self.eap {
            ret.insert("eap", zvariant::Value::new(v));
        }
        if let Some(v) = &self.client_cert {
            ret.insert("client-cert", zvariant::Value::new(v));
        }
        if let Some(v) = &self.ca_cert {
            ret.insert("ca-cert", zvariant::Value::new(v));
        }
        if let Some(v) = &self.private_key_password {
            ret.insert("private-key-password", zvariant::Value::new(v));
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
        if let Some(v) = &self.identity {
            ret.insert("identity".to_string(), zvariant::Value::new(v));
        }
        if let Some(v) = &self.private_key {
            ret.insert(
                "private-key".to_string(),
                if let Ok(path) = Self::glib_bytes_to_file_path(v) {
                    zvariant::Value::new(path)
                } else {
                    zvariant::Value::new(v)
                },
            );
        }
        if let Some(v) = &self.eap {
            // Need NULL append at the end
            let mut new_eaps = v.clone();
            new_eaps.push("".to_string());
            ret.insert("eap".to_string(), zvariant::Value::new(new_eaps));
        }
        if let Some(v) = &self.client_cert {
            ret.insert(
                "client-cert".to_string(),
                if let Ok(path) = Self::glib_bytes_to_file_path(v) {
                    zvariant::Value::new(path)
                } else {
                    zvariant::Value::new(v)
                },
            );
        }
        if let Some(v) = &self.ca_cert {
            ret.insert(
                "ca-cert".to_string(),
                if let Ok(path) = Self::glib_bytes_to_file_path(v) {
                    zvariant::Value::new(path)
                } else {
                    zvariant::Value::new(v)
                },
            );
        }
        if let Some(v) = &self.private_key_password {
            ret.insert(
                "private-key-password".to_string(),
                zvariant::Value::new(v),
            );
        }
        Ok(ret)
    }

    pub(crate) fn fill_secrets(&mut self, secrets: &DbusDictionary) {
        if let Some(v) = secrets.get("private-key-password") {
            match String::try_from(v.clone()) {
                Ok(s) => {
                    self.private_key_password = Some(s);
                }
                Err(e) => {
                    log::warn!(
                        "Filed to convert private_key_password: \
                        {:?} {:?}",
                        v,
                        e
                    );
                }
            }
        }
    }

    pub fn file_path_to_glib_bytes(file_path: &str) -> Vec<u8> {
        format!("{}{}\0", GLIB_FILE_PATH_PREFIX, file_path).into_bytes()
    }

    pub fn glib_bytes_to_file_path(value: &[u8]) -> Result<String, NmError> {
        let mut file_path = match String::from_utf8(value.to_vec()) {
            Ok(f) => f.trim_end_matches(char::from(0)).to_string(),
            Err(e) => {
                let e = NmError::new(
                    ErrorKind::InvalidArgument,
                    format!(
                        "Failed to parse glib bytes to UTF-8 string: \
                        {:?}: {:?}",
                        value, e
                    ),
                );
                log::error!("{}", e);
                return Err(e);
            }
        };
        if file_path.starts_with(GLIB_FILE_PATH_PREFIX) {
            file_path.drain(..GLIB_FILE_PATH_PREFIX.len());
            Ok(file_path)
        } else {
            let e = NmError::new(
                ErrorKind::InvalidArgument,
                format!(
                    "Specified glib bytes is started with {}: {:?}",
                    GLIB_FILE_PATH_PREFIX, value
                ),
            );
            log::error!("{}", e);
            Err(e)
        }
    }
}
