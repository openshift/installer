// SPDX-License-Identifier: MIT

use crate::{EthtoolCoalesceGetRequest, EthtoolHandle};

pub struct EthtoolCoalesceHandle(EthtoolHandle);

impl EthtoolCoalesceHandle {
    pub fn new(handle: EthtoolHandle) -> Self {
        EthtoolCoalesceHandle(handle)
    }

    /// Retrieve the ethtool coalesces of a interface (equivalent to `ethtool -c eth1`)
    pub fn get(&mut self, iface_name: Option<&str>) -> EthtoolCoalesceGetRequest {
        EthtoolCoalesceGetRequest::new(self.0.clone(), iface_name)
    }
}
