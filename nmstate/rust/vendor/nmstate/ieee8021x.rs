use serde::{Deserialize, Serialize};

use crate::NetworkState;

#[derive(Debug, Clone, PartialEq, Eq, Default, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub struct Ieee8021XConfig {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub identity: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none", rename = "eap-methods")]
    pub eap: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub private_key: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub client_cert: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ca_cert: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub private_key_password: Option<String>,
}

impl Ieee8021XConfig {
    pub(crate) fn hide_secrets(&mut self) {
        if self.private_key_password.is_some() {
            self.private_key_password =
                Some(NetworkState::PASSWORD_HID_BY_NMSTATE.to_string());
        }
    }
}
