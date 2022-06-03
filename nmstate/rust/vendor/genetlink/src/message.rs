// SPDX-License-Identifier: MIT

//! Raw generic netlink payload message
//!
//! # Design
//! Since we use generic type to represent different generic family's message type,
//! and it is not easy to create a underlying [`netlink_proto::new_connection()`]
//! with trait object to multiplex different generic netlink family's message.
//!
//! Therefore, I decided to serialize the generic type payload into bytes before
//! sending to the underlying connection. The [`RawGenlMessage`] is meant for this.
//!
//! This special message doesn't use generic type and its payload is `Vec<u8>`.
//! Therefore, its type is easier to use.
//!
//! Another advantage is that it can let users know when the generic netlink payload
//! fails to decode instead of just dropping the messages.
//! (`netlink_proto` would drop messages if they fails to decode.)
//! I think this can help developers debug their deserializing implementation.
use netlink_packet_core::{
    DecodeError,
    NetlinkDeserializable,
    NetlinkHeader,
    NetlinkMessage,
    NetlinkPayload,
    NetlinkSerializable,
};
use netlink_packet_generic::{GenlBuffer, GenlFamily, GenlHeader, GenlMessage};
use netlink_packet_utils::{Emitable, Parseable, ParseableParametrized};
use std::fmt::Debug;

/// Message type to hold serialized generic netlink payload
///
/// **Note** This message type is not intend to be used by normal users, unless
/// you need to use the `UnboundedReceiver<(NetlinkMessage<RawGenlMessage>, SocketAddr)>`
/// return by [`new_connection()`](crate::new_connection)
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct RawGenlMessage {
    pub header: GenlHeader,
    pub payload: Vec<u8>,
    pub family_id: u16,
}

impl RawGenlMessage {
    /// Construct the message
    pub fn new(header: GenlHeader, payload: Vec<u8>, family_id: u16) -> Self {
        Self {
            header,
            payload,
            family_id,
        }
    }

    /// Consume this message and return its header and payload
    pub fn into_parts(self) -> (GenlHeader, Vec<u8>) {
        (self.header, self.payload)
    }

    /// Serialize the generic netlink payload into raw bytes
    pub fn from_genlmsg<F>(genlmsg: GenlMessage<F>) -> Self
    where
        F: GenlFamily + Emitable + Debug,
    {
        let mut payload_buf = vec![0u8; genlmsg.payload.buffer_len()];
        genlmsg.payload.emit(&mut payload_buf);

        Self {
            header: genlmsg.header,
            payload: payload_buf,
            family_id: genlmsg.family_id(),
        }
    }

    /// Try to deserialize the generic netlink payload from raw bytes
    pub fn parse_into_genlmsg<F>(&self) -> Result<GenlMessage<F>, DecodeError>
    where
        F: GenlFamily + ParseableParametrized<[u8], GenlHeader> + Debug,
    {
        let inner = F::parse_with_param(&self.payload, self.header)?;
        Ok(GenlMessage::new(self.header, inner, self.family_id))
    }
}

impl Emitable for RawGenlMessage {
    fn buffer_len(&self) -> usize {
        self.header.buffer_len() + self.payload.len()
    }

    fn emit(&self, buffer: &mut [u8]) {
        self.header.emit(buffer);

        let buffer = &mut buffer[self.header.buffer_len()..];
        buffer.copy_from_slice(&self.payload);
    }
}

impl<'a, T> ParseableParametrized<GenlBuffer<&'a T>, u16> for RawGenlMessage
where
    T: AsRef<[u8]> + ?Sized,
{
    fn parse_with_param(buf: &GenlBuffer<&'a T>, message_type: u16) -> Result<Self, DecodeError> {
        let header = GenlHeader::parse(buf)?;
        let payload_buf = buf.payload();
        Ok(RawGenlMessage::new(
            header,
            payload_buf.to_vec(),
            message_type,
        ))
    }
}

impl NetlinkSerializable for RawGenlMessage {
    fn message_type(&self) -> u16 {
        self.family_id
    }

    fn buffer_len(&self) -> usize {
        <Self as Emitable>::buffer_len(self)
    }

    fn serialize(&self, buffer: &mut [u8]) {
        self.emit(buffer)
    }
}

impl NetlinkDeserializable for RawGenlMessage {
    type Error = DecodeError;
    fn deserialize(header: &NetlinkHeader, payload: &[u8]) -> Result<Self, Self::Error> {
        let buffer = GenlBuffer::new_checked(payload)?;
        RawGenlMessage::parse_with_param(&buffer, header.message_type)
    }
}

impl From<RawGenlMessage> for NetlinkPayload<RawGenlMessage> {
    fn from(message: RawGenlMessage) -> Self {
        NetlinkPayload::InnerMessage(message)
    }
}

/// Helper function to map the [`NetlinkPayload`] types in [`NetlinkMessage`]
/// and serialize the generic netlink payload into raw bytes.
pub fn map_to_rawgenlmsg<F>(
    message: NetlinkMessage<GenlMessage<F>>,
) -> NetlinkMessage<RawGenlMessage>
where
    F: GenlFamily + Emitable + Debug,
{
    let raw_payload = match message.payload {
        NetlinkPayload::InnerMessage(genlmsg) => {
            NetlinkPayload::InnerMessage(RawGenlMessage::from_genlmsg(genlmsg))
        }
        NetlinkPayload::Done => NetlinkPayload::Done,
        NetlinkPayload::Error(i) => NetlinkPayload::Error(i),
        NetlinkPayload::Ack(i) => NetlinkPayload::Ack(i),
        NetlinkPayload::Noop => NetlinkPayload::Noop,
        NetlinkPayload::Overrun(i) => NetlinkPayload::Overrun(i),
    };
    NetlinkMessage::new(message.header, raw_payload)
}

/// Helper function to map the [`NetlinkPayload`] types in [`NetlinkMessage`]
/// and try to deserialize the generic netlink payload from raw bytes.
pub fn map_from_rawgenlmsg<F>(
    raw_msg: NetlinkMessage<RawGenlMessage>,
) -> Result<NetlinkMessage<GenlMessage<F>>, DecodeError>
where
    F: GenlFamily + ParseableParametrized<[u8], GenlHeader> + Debug,
{
    let payload = match raw_msg.payload {
        NetlinkPayload::InnerMessage(raw_genlmsg) => {
            NetlinkPayload::InnerMessage(raw_genlmsg.parse_into_genlmsg()?)
        }
        NetlinkPayload::Done => NetlinkPayload::Done,
        NetlinkPayload::Error(i) => NetlinkPayload::Error(i),
        NetlinkPayload::Ack(i) => NetlinkPayload::Ack(i),
        NetlinkPayload::Noop => NetlinkPayload::Noop,
        NetlinkPayload::Overrun(i) => NetlinkPayload::Overrun(i),
    };
    Ok(NetlinkMessage::new(raw_msg.header, payload))
}
