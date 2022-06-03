use serde::{Deserialize, Serialize};

use crate::{ErrorKind, NmstateError};

#[derive(Debug, Clone, PartialEq, Eq, Default, Serialize, Deserialize)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct HostNameState {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub running: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub config: Option<String>,
}

impl HostNameState {
    pub(crate) fn update(&mut self, other: &Self) {
        if other.running.is_some() {
            self.running = other.running.clone();
        }
        if other.config.is_some() {
            self.config = other.config.clone();
        }
    }

    pub(crate) fn verify(
        &self,
        current: Option<&Self>,
    ) -> Result<(), NmstateError> {
        let current = if let Some(c) = current {
            c
        } else {
            // Should never happen
            let e = NmstateError::new(
                ErrorKind::Bug,
                "Got None HostNameState as current".to_string(),
            );
            log::error!("{}", e);
            return Err(e);
        };

        if let Some(running) = self.running.as_ref() {
            if Some(running) != current.running.as_ref() {
                let e = NmstateError::new(
                    ErrorKind::VerificationError,
                    format!(
                        "Verification fail, desire hostname.running: \
                        {}, current: {:?}",
                        running,
                        current.running.as_ref()
                    ),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }
        if let Some(config) = self.config.as_ref() {
            if Some(config) != current.config.as_ref() {
                let e = NmstateError::new(
                    ErrorKind::VerificationError,
                    format!(
                        "Verification fail, desire hostname.config: \
                        {}, current: {:?}",
                        config,
                        current.config.as_ref()
                    ),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }
        Ok(())
    }
}
