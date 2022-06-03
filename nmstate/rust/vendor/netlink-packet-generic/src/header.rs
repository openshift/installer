// SPDX-License-Identifier: MIT

//! header definition of generic netlink packet
use crate::{buffer::GenlBuffer, constants::GENL_HDRLEN};
use netlink_packet_core::DecodeError;
use netlink_packet_utils::{Emitable, Parseable};

/// Generic Netlink header
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct GenlHeader {
    pub cmd: u8,
    pub version: u8,
}

impl Emitable for GenlHeader {
    fn buffer_len(&self) -> usize {
        GENL_HDRLEN
    }

    fn emit(&self, buffer: &mut [u8]) {
        let mut packet = GenlBuffer::new(buffer);
        packet.set_cmd(self.cmd);
        packet.set_version(self.version);
    }
}

impl<T: AsRef<[u8]>> Parseable<GenlBuffer<T>> for GenlHeader {
    fn parse(buf: &GenlBuffer<T>) -> Result<Self, DecodeError> {
        Ok(Self {
            cmd: buf.cmd(),
            version: buf.version(),
        })
    }
}
