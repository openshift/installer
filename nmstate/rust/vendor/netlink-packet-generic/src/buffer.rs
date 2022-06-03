// SPDX-License-Identifier: MIT

//! Buffer definition of generic netlink packet
use crate::{constants::GENL_HDRLEN, header::GenlHeader, message::GenlMessage};
use netlink_packet_core::DecodeError;
use netlink_packet_utils::{Parseable, ParseableParametrized};
use std::fmt::Debug;

buffer!(GenlBuffer(GENL_HDRLEN) {
    cmd: (u8, 0),
    version: (u8, 1),
    payload: (slice, GENL_HDRLEN..),
});

impl<F> ParseableParametrized<[u8], u16> for GenlMessage<F>
where
    F: ParseableParametrized<[u8], GenlHeader> + Debug,
{
    fn parse_with_param(buf: &[u8], message_type: u16) -> Result<Self, DecodeError> {
        let buf = GenlBuffer::new_checked(buf)?;
        Self::parse_with_param(&buf, message_type)
    }
}

impl<'a, F, T> ParseableParametrized<GenlBuffer<&'a T>, u16> for GenlMessage<F>
where
    F: ParseableParametrized<[u8], GenlHeader> + Debug,
    T: AsRef<[u8]> + ?Sized,
{
    fn parse_with_param(buf: &GenlBuffer<&'a T>, message_type: u16) -> Result<Self, DecodeError> {
        let header = GenlHeader::parse(buf)?;
        let payload_buf = buf.payload();
        Ok(GenlMessage::new(
            header,
            F::parse_with_param(payload_buf, header)?,
            message_type,
        ))
    }
}
