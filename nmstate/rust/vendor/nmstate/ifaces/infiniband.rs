use serde::{Deserialize, Serialize, Serializer};

use crate::{BaseInterface, InterfaceType};

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq, Eq)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub struct InfiniBandInterface {
    #[serde(flatten)]
    pub base: BaseInterface,
    #[serde(skip_serializing_if = "Option::is_none", rename = "infiniband")]
    pub ib: Option<InfiniBandConfig>,
}

impl Default for InfiniBandInterface {
    fn default() -> Self {
        let mut base = BaseInterface::new();
        base.iface_type = InterfaceType::InfiniBand;
        Self { base, ib: None }
    }
}

impl InfiniBandInterface {
    pub(crate) fn parent(&self) -> Option<&str> {
        self.ib.as_ref().and_then(|cfg| cfg.base_iface.as_deref())
    }

    pub(crate) fn update_ib(&mut self, other: &InfiniBandInterface) {
        if let Some(ib_conf) = &mut self.ib {
            ib_conf.update(other.ib.as_ref());
        } else {
            self.ib = other.ib.clone();
        }
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone, Copy)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum InfiniBandMode {
    Datagram,
    Connected,
}

impl Default for InfiniBandMode {
    fn default() -> Self {
        Self::Datagram
    }
}

impl std::fmt::Display for InfiniBandMode {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                InfiniBandMode::Datagram => "datagram",
                InfiniBandMode::Connected => "connected",
            }
        )
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "kebab-case", deny_unknown_fields)]
#[non_exhaustive]
pub struct InfiniBandConfig {
    pub mode: InfiniBandMode,
    #[serde(skip_serializing_if = "crate::serializer::is_option_string_empty")]
    pub base_iface: Option<String>,
    #[serde(
        serialize_with = "show_as_hex",
        deserialize_with = "crate::deserializer::option_u16_or_string"
    )]
    pub pkey: Option<u16>,
}

impl InfiniBandConfig {
    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn update(&mut self, other: Option<&InfiniBandConfig>) {
        if let Some(other) = other {
            self.mode = other.mode;
            self.pkey = other.pkey;
            self.base_iface = other.base_iface.clone();
        }
    }
}

fn show_as_hex<S>(v: &Option<u16>, s: S) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    if let Some(v) = v {
        s.serialize_str(&format!("{:#02x}", v))
    } else {
        s.serialize_none()
    }
}
