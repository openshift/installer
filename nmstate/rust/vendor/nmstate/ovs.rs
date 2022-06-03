use std::collections::HashMap;

use serde::{Deserialize, Deserializer, Serialize};

use crate::{state::get_json_value_difference, ErrorKind, NmstateError};

#[derive(Debug, Clone, PartialEq, Eq, Default, Serialize)]
#[non_exhaustive]
pub struct OvsDbGlobalConfig {
    #[serde(skip_serializing_if = "Option::is_none")]
    // When the value been set as None, specified key will be removed instead
    // of merging.
    // To remove all settings of external_ids or other_config, use empty
    // HashMap
    pub external_ids: Option<HashMap<String, Option<String>>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub other_config: Option<HashMap<String, Option<String>>>,
}

impl OvsDbGlobalConfig {
    pub fn is_none(&self) -> bool {
        self.external_ids.is_none() && self.other_config.is_none()
    }

    pub(crate) fn verify(&self, current: &Self) -> Result<(), NmstateError> {
        let self_value = serde_json::to_value(self)?;
        let current_value = serde_json::to_value(current)?;

        if let Some((reference, desire, current)) = get_json_value_difference(
            "ovsdb".to_string(),
            &self_value,
            &current_value,
        ) {
            let e = NmstateError::new(
                ErrorKind::VerificationError,
                format!(
                    "Verification failure: {} desire '{}', current '{}'",
                    reference, desire, current
                ),
            );
            log::error!("{}", e);
            Err(e)
        } else {
            Ok(())
        }
    }

    pub(crate) fn get_external_ids(&self) -> HashMap<&str, &str> {
        let mut ret = HashMap::new();
        if let Some(eids) = self.external_ids.as_ref() {
            for (k, v) in eids {
                if let Some(v) = v {
                    ret.insert(k.as_str(), v.as_str());
                }
            }
        }
        ret
    }

    pub(crate) fn get_other_config(&self) -> HashMap<&str, &str> {
        let mut ret = HashMap::new();
        if let Some(ocfg) = self.other_config.as_ref() {
            for (k, v) in ocfg.iter() {
                if let Some(v) = v {
                    ret.insert(k.as_str(), v.as_str());
                }
            }
        }
        ret
    }

    // Currently, we only support full editing on OVSDB.
    //  * If OVSDB setting not mentioned, preserve old configure.
    //  * If OVSDB is set, override.
    pub(crate) fn merge(&mut self, current: &Self) {
        if self.external_ids.is_some() || self.other_config.is_some() {
            if self.external_ids.is_none() {
                self.external_ids = current.external_ids.clone();
            }
            if self.other_config.is_none() {
                self.other_config = current.other_config.clone();
            }
        }
    }
}

impl<'de> Deserialize<'de> for OvsDbGlobalConfig {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let mut ret = Self::default();
        let v = serde_json::Value::deserialize(deserializer)?;
        if let Some(v) = v.as_object() {
            if let Some(v) = v.get("external_ids") {
                ret.external_ids = Some(value_to_hash_map(v));
            }
            if let Some(v) = v.get("other_config") {
                ret.other_config = Some(value_to_hash_map(v));
            }
        } else {
            return Err(serde::de::Error::custom(format!(
                "Expecting dict/HashMap, but got {:?}",
                v
            )));
        }
        Ok(ret)
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default, Serialize)]
#[non_exhaustive]
pub struct OvsDbIfaceConfig {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub external_ids: Option<HashMap<String, Option<String>>>,
}

impl OvsDbIfaceConfig {
    pub(crate) fn get_external_ids(&self) -> HashMap<&str, &str> {
        let mut ret = HashMap::new();
        if let Some(eids) = self.external_ids.as_ref() {
            for (k, v) in eids {
                if let Some(v) = v {
                    ret.insert(k.as_str(), v.as_str());
                }
            }
        }
        ret
    }
}

impl<'de> Deserialize<'de> for OvsDbIfaceConfig {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let mut ret = Self::default();
        let v = serde_json::Value::deserialize(deserializer)?;
        if let Some(v) = v.as_object() {
            if let Some(v) = v.get("external_ids") {
                ret.external_ids = Some(value_to_hash_map(v));
            }
        } else {
            return Err(serde::de::Error::custom(format!(
                "Expecting dict/HashMap, but got {:?}",
                v
            )));
        }
        Ok(ret)
    }
}

fn value_to_hash_map(
    value: &serde_json::Value,
) -> HashMap<String, Option<String>> {
    let mut ret: HashMap<String, Option<String>> = HashMap::new();
    if let Some(value) = value.as_object() {
        for (k, v) in value.iter() {
            let v = match &v {
                serde_json::Value::Number(i) => Some({
                    if let Some(i) = i.as_i64() {
                        format!("{}", i)
                    } else if let Some(i) = i.as_u64() {
                        format!("{}", i)
                    } else if let Some(i) = i.as_f64() {
                        format!("{}", i)
                    } else {
                        continue;
                    }
                }),
                serde_json::Value::String(s) => Some(s.to_string()),
                serde_json::Value::Bool(b) => Some(format!("{}", b)),
                serde_json::Value::Null => None,
                _ => continue,
            };
            ret.insert(k.to_string(), v);
        }
    }
    ret
}
