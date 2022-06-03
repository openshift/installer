// SPDX-License-Identifier: MIT

//! Message definition and method implementations

use crate::{buffer::GenlBuffer, header::GenlHeader, traits::*};
use netlink_packet_core::{
    DecodeError,
    NetlinkDeserializable,
    NetlinkHeader,
    NetlinkPayload,
    NetlinkSerializable,
};
use netlink_packet_utils::{Emitable, ParseableParametrized};
use std::fmt::Debug;

#[cfg(doc)]
use netlink_packet_core::NetlinkMessage;

/// Represent the generic netlink messages
///
/// This type can wrap data types `F` which represents a generic family payload.
/// The message can be serialize/deserialize if the type `F` implements [`GenlFamily`],
/// [`Emitable`], and [`ParseableParametrized<[u8], GenlHeader>`](ParseableParametrized).
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct GenlMessage<F> {
    pub header: GenlHeader,
    pub payload: F,
    resolved_family_id: u16,
}

impl<F> GenlMessage<F>
where
    F: Debug,
{
    /// Construct the message
    pub fn new(header: GenlHeader, payload: F, family_id: u16) -> Self {
        Self {
            header,
            payload,
            resolved_family_id: family_id,
        }
    }

    /// Construct the message by the given header and payload
    pub fn from_parts(header: GenlHeader, payload: F) -> Self {
        Self {
            header,
            payload,
            resolved_family_id: 0,
        }
    }

    /// Consume this message and return its header and payload
    pub fn into_parts(self) -> (GenlHeader, F) {
        (self.header, self.payload)
    }

    /// Return the previously set resolved family ID in this message.
    ///
    /// This value would be used to serialize the message only if
    /// the ([`GenlFamily::family_id()`]) return 0 in the underlying type.
    pub fn resolved_family_id(&self) -> u16 {
        self.resolved_family_id
    }

    /// Set the resolved dynamic family ID of the message, if the generic family
    /// uses dynamic generated ID by kernel.
    ///
    /// This method is a interface to provide other high level library to
    /// set the resolved family ID before the message is serialized.
    ///
    /// # Usage
    /// Normally, you don't have to call this function directly if you are
    /// using library which helps you handle the dynamic family id.
    ///
    /// If you are the developer of some high level generic netlink library,
    /// you can call this method to set the family id resolved by your resolver.
    /// Without having to modify the `message_type` field of the serialized
    /// netlink packet header before sending it.
    pub fn set_resolved_family_id(&mut self, family_id: u16) {
        self.resolved_family_id = family_id;
    }
}

impl<F> GenlMessage<F>
where
    F: GenlFamily + Debug,
{
    /// Build the message from the payload
    ///
    /// This function would automatically fill the header for you. You can directly emit
    /// the message without having to call [`finalize()`](Self::finalize).
    pub fn from_payload(payload: F) -> Self {
        Self {
            header: GenlHeader {
                cmd: payload.command(),
                version: payload.version(),
            },
            payload,
            resolved_family_id: 0,
        }
    }

    /// Ensure the header ([`GenlHeader`]) is consistent with the payload (`F: GenlFamily`):
    ///
    /// - Fill the command and version number into the header
    ///
    /// If you are not 100% sure the header is correct, this method should be called before calling
    /// [`Emitable::emit()`], as it could get error result if the header is inconsistent with the message.
    pub fn finalize(&mut self) {
        self.header.cmd = self.payload.command();
        self.header.version = self.payload.version();
    }

    /// Return the resolved family ID which should be filled into the `message_type`
    /// field in [`NetlinkHeader`].
    ///
    /// The implementation of [`NetlinkSerializable::message_type()`] would use
    /// this function's result as its the return value. Thus, the family id can
    /// be automatically filled into the `message_type` during the call to
    /// [`NetlinkMessage::finalize()`].
    pub fn family_id(&self) -> u16 {
        let static_id = self.payload.family_id();
        if static_id == 0 {
            self.resolved_family_id
        } else {
            static_id
        }
    }
}

impl<F> Emitable for GenlMessage<F>
where
    F: GenlFamily + Emitable + Debug,
{
    fn buffer_len(&self) -> usize {
        self.header.buffer_len() + self.payload.buffer_len()
    }

    fn emit(&self, buffer: &mut [u8]) {
        self.header.emit(buffer);

        let buffer = &mut buffer[self.header.buffer_len()..];
        self.payload.emit(buffer);
    }
}

impl<F> NetlinkSerializable for GenlMessage<F>
where
    F: GenlFamily + Emitable + Debug,
{
    fn message_type(&self) -> u16 {
        self.family_id()
    }

    fn buffer_len(&self) -> usize {
        <Self as Emitable>::buffer_len(self)
    }

    fn serialize(&self, buffer: &mut [u8]) {
        self.emit(buffer)
    }
}

impl<F> NetlinkDeserializable for GenlMessage<F>
where
    F: ParseableParametrized<[u8], GenlHeader> + Debug,
{
    type Error = DecodeError;
    fn deserialize(header: &NetlinkHeader, payload: &[u8]) -> Result<Self, Self::Error> {
        let buffer = GenlBuffer::new_checked(payload)?;
        GenlMessage::parse_with_param(&buffer, header.message_type)
    }
}

impl<F> From<GenlMessage<F>> for NetlinkPayload<GenlMessage<F>>
where
    F: Debug,
{
    fn from(message: GenlMessage<F>) -> Self {
        NetlinkPayload::InnerMessage(message)
    }
}
