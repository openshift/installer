use serde::{Deserialize, Serialize};

use crate::{BaseInterface, ErrorKind, InterfaceType, NmstateError};

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[non_exhaustive]
pub struct MacVtapInterface {
    #[serde(flatten)]
    pub base: BaseInterface,
    #[serde(skip_serializing_if = "Option::is_none", rename = "mac-vtap")]
    pub mac_vtap: Option<MacVtapConfig>,
}

impl Default for MacVtapInterface {
    fn default() -> Self {
        Self {
            base: BaseInterface {
                iface_type: InterfaceType::MacVtap,
                ..Default::default()
            },
            mac_vtap: None,
        }
    }
}

impl MacVtapInterface {
    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn validate(&self) -> Result<(), NmstateError> {
        if let Some(conf) = &self.mac_vtap {
            if conf.accept_all_mac == Some(false)
                && conf.mode != MacVtapMode::Passthru
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
        self.mac_vtap.as_ref().map(|cfg| cfg.base_iface.as_str())
    }

    pub(crate) fn update_mac_vtap(&mut self, other: &MacVtapInterface) {
        // TODO: this should be done by Trait
        if let Some(conf) = &mut self.mac_vtap {
            conf.update(other.mac_vtap.as_ref());
        } else {
            self.mac_vtap = other.mac_vtap.clone();
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "kebab-case", deny_unknown_fields)]
#[non_exhaustive]
pub struct MacVtapConfig {
    pub base_iface: String,
    pub mode: MacVtapMode,
    #[serde(
        skip_serializing_if = "Option::is_none",
        rename = "promiscuous",
        alias = "accept-all-mac",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub accept_all_mac: Option<bool>,
}

impl MacVtapConfig {
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
pub enum MacVtapMode {
    Vepa,
    Bridge,
    Private,
    Passthru,
    Source,
    Unknown,
}

impl From<MacVtapMode> for u32 {
    fn from(v: MacVtapMode) -> u32 {
        match v {
            MacVtapMode::Unknown => 0,
            MacVtapMode::Vepa => 1,
            MacVtapMode::Bridge => 2,
            MacVtapMode::Private => 3,
            MacVtapMode::Passthru => 4,
            MacVtapMode::Source => 5,
        }
    }
}

impl Default for MacVtapMode {
    fn default() -> Self {
        Self::Unknown
    }
}
