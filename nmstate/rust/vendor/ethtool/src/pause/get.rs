// SPDX-License-Identifier: MIT

use futures::TryStream;
use netlink_packet_generic::GenlMessage;

use crate::{ethtool_execute, EthtoolError, EthtoolHandle, EthtoolMessage};

pub struct EthtoolPauseGetRequest {
    handle: EthtoolHandle,
    iface_name: Option<String>,
}

impl EthtoolPauseGetRequest {
    pub(crate) fn new(handle: EthtoolHandle, iface_name: Option<&str>) -> Self {
        EthtoolPauseGetRequest {
            handle,
            iface_name: iface_name.map(|i| i.to_string()),
        }
    }

    pub async fn execute(
        self,
    ) -> impl TryStream<Ok = GenlMessage<EthtoolMessage>, Error = EthtoolError> {
        let EthtoolPauseGetRequest {
            mut handle,
            iface_name,
        } = self;

        let ethtool_msg = EthtoolMessage::new_pause_get(iface_name.as_deref());
        ethtool_execute(&mut handle, iface_name.is_none(), ethtool_msg).await
    }
}
