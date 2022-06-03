// SPDX-License-Identifier: MIT

use crate::{EthtoolHandle, EthtoolLinkModeGetRequest};

pub struct EthtoolLinkModeHandle(EthtoolHandle);

impl EthtoolLinkModeHandle {
    pub fn new(handle: EthtoolHandle) -> Self {
        EthtoolLinkModeHandle(handle)
    }

    /// Retrieve the ethtool link_modes(duplex, link speed and etc) of a interface
    pub fn get(&mut self, iface_name: Option<&str>) -> EthtoolLinkModeGetRequest {
        EthtoolLinkModeGetRequest::new(self.0.clone(), iface_name)
    }
}
