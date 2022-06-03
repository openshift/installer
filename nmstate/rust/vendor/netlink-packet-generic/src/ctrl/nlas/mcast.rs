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
pub enum McastGrpAttrs {
    Name(String),
    Id(u32),
}

impl Nla for McastGrpAttrs {
    fn value_len(&self) -> usize {
        use McastGrpAttrs::*;
        match self {
            Name(s) => s.as_bytes().len() + 1,
            Id(v) => size_of_val(v),
        }
    }

    fn kind(&self) -> u16 {
        use McastGrpAttrs::*;
        match self {
            Name(_) => CTRL_ATTR_MCAST_GRP_NAME,
            Id(_) => CTRL_ATTR_MCAST_GRP_ID,
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        use McastGrpAttrs::*;
        match self {
            Name(s) => {
                buffer[..s.len()].copy_from_slice(s.as_bytes());
                buffer[s.len()] = 0;
            }
            Id(v) => NativeEndian::write_u32(buffer, *v),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for McastGrpAttrs {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            CTRL_ATTR_MCAST_GRP_NAME => {
                Self::Name(parse_string(payload).context("invalid CTRL_ATTR_MCAST_GRP_NAME value")?)
            }
            CTRL_ATTR_MCAST_GRP_ID => {
                Self::Id(parse_u32(payload).context("invalid CTRL_ATTR_MCAST_GRP_ID value")?)
            }
            kind => return Err(DecodeError::from(format!("Unknown NLA type: {}", kind))),
        })
    }
}
