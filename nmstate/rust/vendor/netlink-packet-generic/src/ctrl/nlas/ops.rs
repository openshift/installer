// SPDX-License-Identifier: MIT

use crate::constants::*;
use anyhow::Context;
use byteorder::{ByteOrder, NativeEndian};
use netlink_packet_utils::{
    nla::{Nla, NlaBuffer},
    parsers::*,
    traits::*,
    DecodeError,
};
use std::mem::size_of_val;

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum OpAttrs {
    Id(u32),
    Flags(u32),
}

impl Nla for OpAttrs {
    fn value_len(&self) -> usize {
        use OpAttrs::*;
        match self {
            Id(v) => size_of_val(v),
            Flags(v) => size_of_val(v),
        }
    }

    fn kind(&self) -> u16 {
        use OpAttrs::*;
        match self {
            Id(_) => CTRL_ATTR_OP_ID,
            Flags(_) => CTRL_ATTR_OP_FLAGS,
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        use OpAttrs::*;
        match self {
            Id(v) => NativeEndian::write_u32(buffer, *v),
            Flags(v) => NativeEndian::write_u32(buffer, *v),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for OpAttrs {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            CTRL_ATTR_OP_ID => {
                Self::Id(parse_u32(payload).context("invalid CTRL_ATTR_OP_ID value")?)
            }
            CTRL_ATTR_OP_FLAGS => {
                Self::Flags(parse_u32(payload).context("invalid CTRL_ATTR_OP_FLAGS value")?)
            }
            kind => return Err(DecodeError::from(format!("Unknown NLA type: {}", kind))),
        })
    }
}
