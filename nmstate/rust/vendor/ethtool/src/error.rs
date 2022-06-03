// SPDX-License-Identifier: MIT

use thiserror::Error;

use netlink_packet_core::{ErrorMessage, NetlinkMessage};
use netlink_packet_generic::GenlMessage;

use crate::EthtoolMessage;

#[derive(Clone, Eq, PartialEq, Debug, Error)]
pub enum EthtoolError {
    #[error("Received an unexpected message {0:?}")]
    UnexpectedMessage(NetlinkMessage<GenlMessage<EthtoolMessage>>),

    #[error("Received a netlink error message {0}")]
    NetlinkError(ErrorMessage),

    #[error("A netlink request failed")]
    RequestFailed(String),

    #[error("A bug in this crate")]
    Bug(String),
}
