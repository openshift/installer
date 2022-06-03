// SPDX-License-Identifier: MIT

#[macro_export]
macro_rules! try_ethtool {
    ($msg: expr) => {{
        use netlink_packet_core::{NetlinkMessage, NetlinkPayload};
        use $crate::EthtoolError;

        match $msg {
            Ok(msg) => {
                let (header, payload) = msg.into_parts();
                match payload {
                    NetlinkPayload::InnerMessage(msg) => msg,
                    NetlinkPayload::Error(err) => return Err(EthtoolError::NetlinkError(err)),
                    _ => {
                        return Err(EthtoolError::UnexpectedMessage(NetlinkMessage::new(
                            header, payload,
                        )))
                    }
                }
            }
            Err(e) => return Err(EthtoolError::Bug(format!("BUG: decode error {:?}", e))),
        }
    }};
}
