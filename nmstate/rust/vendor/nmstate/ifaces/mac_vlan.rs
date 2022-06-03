use serde::{Deserialize, Serialize};

use crate::{BaseInterface, ErrorKind, InterfaceType, NmstateError};

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[non_exhaustive]
pub struct MacVlanInterface {
    #[serde(flatten)]
    pub base: BaseInterface,
    #[serde(skip_serializing_if = "Option::is_none", rename = "mac-vlan")]
    pub mac_vlan: Option<MacVlanConfig>,
}

impl Default for MacVlanInterface {
    fn default() -> Self {
        Self {
            base: BaseInterface {
                iface_type: InterfaceType::MacVlan,
                ..Default::default()
            },
            mac_vlan: None,
        }
    }
}

impl MacVlanInterface {
    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn validate(&self) -> Result<(), NmstateError> {
        if let Some(conf) = &self.mac_vlan {
            if conf.accept_all_mac == Some(false)
                && conf.mode != MacVlanMode::Passthru
            {
                let e = NmstateError::new(
                    ErrorKind::InvalidArgument,
                    "Disable accept-all-mac-addresses(promiscuous) \
                    is only allowed on passthru mode"
                        .to_string(),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }
        Ok(())
    }

    pub(crate) fn parent(&self) -> Option<&str> {
        self.mac_vlan.as_ref().map(|cfg| cfg.base_iface.as_str())
    }

    pub(crate) fn update_mac_vlan(&mut self, other: &MacVlanInterface) {
        // TODO: this should be done by Trait
        if let Some(conf) = &mut self.mac_vlan {
            conf.update(other.mac_vlan.as_ref());
        } else {
            self.mac_vlan = other.mac_vlan.clone();
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "kebab-case", deny_unknown_fields)]
#[non_exhaustive]
pub struct MacVlanConfig {
    pub base_iface: String,
    pub mode: MacVlanMode,
    #[serde(
        skip_serializing_if = "Option::is_none",
        rename = "promiscuous",
        alias = "accept-all-mac",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub accept_all_mac: Option<bool>,
}

impl MacVlanConfig {
    fn update(&mut self, other: Option<&Self>) {
        if let Some(other) = other {
            self.base_iface = other.base_iface.clone();
            self.mode = other.mode;
            self.accept_all_mac = other.accept_all_mac;
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum MacVlanMode {
    Vepa,
    Bridge,
    Private,
    Passthru,
    Source,
    Unknown,
}

impl From<MacVlanMode> for u32 {
    fn from(v: MacVlanMode) -> u32 {
        match v {
            MacVlanMode::Unknown => 0,
            MacVlanMode::Vepa => 1,
            MacVlanMode::Bridge => 2,
            MacVlanMode::Private => 3,
            MacVlanMode::Passthru => 4,
            MacVlanMode::Source => 5,
        }
    }
}

impl Default for MacVlanMode {
    fn default() -> Self {
        Self::Unknown
    }
}
