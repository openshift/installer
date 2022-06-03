use serde::{Deserialize, Serialize};

use crate::{BaseInterface, InterfaceType};

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[non_exhaustive]
pub struct VrfInterface {
    #[serde(flatten)]
    pub base: BaseInterface,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub vrf: Option<VrfConfig>,
}

impl Default for VrfInterface {
    fn default() -> Self {
        Self {
            base: BaseInterface {
                iface_type: InterfaceType::Vrf,
                ..Default::default()
            },
            vrf: None,
        }
    }
}

impl VrfInterface {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn ports(&self) -> Option<Vec<&str>> {
        self.vrf
            .as_ref()
            .and_then(|vrf_conf| vrf_conf.port.as_ref())
            .map(|ports| ports.as_slice().iter().map(|p| p.as_str()).collect())
    }

    pub(crate) fn update_vrf(&mut self, other: &VrfInterface) {
        // TODO: this should be done by Trait
        if let Some(vrf_conf) = &mut self.vrf {
            vrf_conf.update(other.vrf.as_ref());
        } else {
            self.vrf = other.vrf.clone();
        }
    }

    pub(crate) fn pre_verify_cleanup(&mut self) {
        self.base.mac_address = None;
        if self.base.accept_all_mac_addresses == Some(false) {
            self.base.accept_all_mac_addresses = None;
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct VrfConfig {
    pub port: Option<Vec<String>>,
    #[serde(
        rename = "route-table-id",
        deserialize_with = "crate::deserializer::u32_or_string"
    )]
    pub table_id: u32,
}

impl VrfConfig {
    fn update(&mut self, other: Option<&Self>) {
        if let Some(other) = other {
            self.port = other.port.clone();
            self.table_id = other.table_id;
        }
    }
}
