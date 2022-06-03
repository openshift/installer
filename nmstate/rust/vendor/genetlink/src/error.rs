// SPDX-License-Identifier: MIT

use crate::message::RawGenlMessage;

/// Error type of genetlink
#[derive(Debug, Error)]
pub enum GenetlinkError {
    #[error("Failed to send netlink request")]
    ProtocolError(#[from] netlink_proto::Error<RawGenlMessage>),
    #[error("Failed to decode generic packet")]
    DecodeError(#[from] netlink_packet_utils::DecodeError),
    #[error("Netlink error message: {0}")]
    NetlinkError(std::io::Error),
    #[error("Cannot find specified netlink attribute: {0}")]
    AttributeNotFound(String),
    #[error("Desire netlink message type not received")]
    NoMessageReceived,
}

// Since `netlink_packet_core::error::ErrorMessage` doesn't impl `Error` trait,
// it need to convert to `std::io::Error` here
impl From<netlink_packet_core::error::ErrorMessage> for GenetlinkError {
    fn from(err_msg: netlink_packet_core::error::ErrorMessage) -> Self {
        Self::NetlinkError(err_msg.to_io())
    }
}
