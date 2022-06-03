// SPDX-License-Identifier: MIT

use crate::{EthtoolHandle, EthtoolRingGetRequest};

pub struct EthtoolRingHandle(EthtoolHandle);

impl EthtoolRingHandle {
    pub fn new(handle: EthtoolHandle) -> Self {
        EthtoolRingHandle(handle)
    }

    /// Retrieve the ethtool rings of a interface (equivalent to `ethtool -g eth1`)
    pub fn get(&mut self, iface_name: Option<&str>) -> EthtoolRingGetRequest {
        EthtoolRingGetRequest::new(self.0.clone(), iface_name)
    }
}
