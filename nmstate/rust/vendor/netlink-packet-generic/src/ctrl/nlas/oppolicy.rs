// SPDX-License-Identifier: MIT

use crate::constants::*;
use anyhow::Context;
use byteorder::{ByteOrder, NativeEndian};
use netlink_packet_utils::{
    nla::{Nla, NlaBuffer, NlasIterator},
    parsers::*,
    traits::*,
    DecodeError,
};
use std::mem::size_of_val;

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct OppolicyAttr {
    pub cmd: u8,
    pub policy_idx: Vec<OppolicyIndexAttr>,
}

impl Nla for OppolicyAttr {
    fn value_len(&self) -> usize {
        self.policy_idx.as_slice().buffer_len()
    }

    fn kind(&self) -> u16 {
        self.cmd as u16
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        self.policy_idx.as_slice().emit(buffer);
    }

    fn is_nested(&self) -> bool {
        true
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for OppolicyAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        let policy_idx = NlasIterator::new(payload)
            .map(|nla| nla.and_then(|nla| OppolicyIndexAttr::parse(&nla)))
            .collect::<Result<Vec<_>, _>>()
            .context("failed to parse OppolicyAttr")?;

        Ok(Self {
            cmd: buf.kind() as u8,
            policy_idx,
        })
    }
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum OppolicyIndexAttr {
    Do(u32),
    Dump(u32),
}

impl Nla for OppolicyIndexAttr {
    fn value_len(&self) -> usize {
        use OppolicyIndexAttr::*;
        match self {
            Do(v) => size_of_val(v),
            Dump(v) => size_of_val(v),
        }
    }

    fn kind(&self) -> u16 {
        use OppolicyIndexAttr::*;
        match self {
            Do(_) => CTRL_ATTR_POLICY_DO,
            Dump(_) => CTRL_ATTR_POLICY_DUMP,
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        use OppolicyIndexAttr::*;
        match self {
            Do(v) => NativeEndian::write_u32(buffer, *v),
            Dump(v) => NativeEndian::write_u32(buffer, *v),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for OppolicyIndexAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            CTRL_ATTR_POLICY_DO => {
                Self::Do(parse_u32(payload).context("invalid CTRL_ATTR_POLICY_DO value")?)
            }
            CTRL_ATTR_POLICY_DUMP => {
                Self::Dump(parse_u32(payload).context("invalid CTRL_ATTR_POLICY_DUMP value")?)
            }
            kind => return Err(DecodeError::from(format!("Unknown NLA type: {}", kind))),
        })
    }
}
