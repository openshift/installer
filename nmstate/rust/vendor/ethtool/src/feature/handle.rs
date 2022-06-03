// SPDX-License-Identifier: MIT

use crate::{EthtoolFeatureGetRequest, EthtoolHandle};

pub struct EthtoolFeatureHandle(EthtoolHandle);

impl EthtoolFeatureHandle {
    pub fn new(handle: EthtoolHandle) -> Self {
        EthtoolFeatureHandle(handle)
    }

    /// Retrieve the ethtool features of a interface (equivalent to `ethtool -k eth1`)
    pub fn get(&mut self, iface_name: Option<&str>) -> EthtoolFeatureGetRequest {
        EthtoolFeatureGetRequest::new(self.0.clone(), iface_name)
    }
}
