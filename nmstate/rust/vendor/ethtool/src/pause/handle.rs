// SPDX-License-Identifier: MIT

use crate::{EthtoolHandle, EthtoolPauseGetRequest};

pub struct EthtoolPauseHandle(EthtoolHandle);

impl EthtoolPauseHandle {
    pub fn new(handle: EthtoolHandle) -> Self {
        EthtoolPauseHandle(handle)
    }

    /// Retrieve the pause setting of a interface (equivalent to `ethtool -a eth1`)
    pub fn get(&mut self, iface_name: Option<&str>) -> EthtoolPauseGetRequest {
        EthtoolPauseGetRequest::new(self.0.clone(), iface_name)
    }
}
