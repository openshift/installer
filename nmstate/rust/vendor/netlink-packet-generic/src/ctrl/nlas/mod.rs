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

mod mcast;
mod oppolicy;
mod ops;
mod policy;

pub use mcast::*;
pub use oppolicy::*;
pub use ops::*;
pub use policy::*;

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum GenlCtrlAttrs {
    FamilyId(u16),
    FamilyName(String),
    Version(u32),
    HdrSize(u32),
    MaxAttr(u32),
    Ops(Vec<Vec<OpAttrs>>),
    McastGroups(Vec<Vec<McastGrpAttrs>>),
    Policy(PolicyAttr),
    OpPolicy(OppolicyAttr),
    Op(u32),
}

impl Nla for GenlCtrlAttrs {
    fn value_len(&self) -> usize {
        use GenlCtrlAttrs::*;
        match self {
            FamilyId(v) => size_of_val(v),
            FamilyName(s) => s.len() + 1,
            Version(v) => size_of_val(v),
            HdrSize(v) => size_of_val(v),
            MaxAttr(v) => size_of_val(v),
            Ops(nlas) => nlas.iter().map(|op| op.as_slice().buffer_len()).sum(),
            McastGroups(nlas) => nlas.iter().map(|op| op.as_slice().buffer_len()).sum(),
            Policy(nla) => nla.buffer_len(),
            OpPolicy(nla) => nla.buffer_len(),
            Op(v) => size_of_val(v),
        }
    }

    fn kind(&self) -> u16 {
        use GenlCtrlAttrs::*;
        match self {
            FamilyId(_) => CTRL_ATTR_FAMILY_ID,
            FamilyName(_) => CTRL_ATTR_FAMILY_NAME,
            Version(_) => CTRL_ATTR_VERSION,
            HdrSize(_) => CTRL_ATTR_HDRSIZE,
            MaxAttr(_) => CTRL_ATTR_MAXATTR,
            Ops(_) => CTRL_ATTR_OPS,
            McastGroups(_) => CTRL_ATTR_MCAST_GROUPS,
            Policy(_) => CTRL_ATTR_POLICY,
            OpPolicy(_) => CTRL_ATTR_OP_POLICY,
            Op(_) => CTRL_ATTR_OP,
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        use GenlCtrlAttrs::*;
        match self {
            FamilyId(v) => NativeEndian::write_u16(buffer, *v),
            FamilyName(s) => {
                buffer[..s.len()].copy_from_slice(s.as_bytes());
                buffer[s.len()] = 0;
            }
            Version(v) => NativeEndian::write_u32(buffer, *v),
            HdrSize(v) => NativeEndian::write_u32(buffer, *v),
            MaxAttr(v) => NativeEndian::write_u32(buffer, *v),
            Ops(nlas) => {
                let mut len = 0;
                for op in nlas {
                    op.as_slice().emit(&mut buffer[len..]);
                    len += op.as_slice().buffer_len();
                }
            }
            McastGroups(nlas) => {
                let mut len = 0;
                for op in nlas {
                    op.as_slice().emit(&mut buffer[len..]);
                    len += op.as_slice().buffer_len();
                }
            }
            Policy(nla) => nla.emit_value(buffer),
            OpPolicy(nla) => nla.emit_value(buffer),
            Op(v) => NativeEndian::write_u32(buffer, *v),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for GenlCtrlAttrs {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            CTRL_ATTR_FAMILY_ID => {
                Self::FamilyId(parse_u16(payload).context("invalid CTRL_ATTR_FAMILY_ID value")?)
            }
            CTRL_ATTR_FAMILY_NAME => Self::FamilyName(
                parse_string(payload).context("invalid CTRL_ATTR_FAMILY_NAME value")?,
            ),
            CTRL_ATTR_VERSION => {
                Self::Version(parse_u32(payload).context("invalid CTRL_ATTR_VERSION value")?)
            }
            CTRL_ATTR_HDRSIZE => {
                Self::HdrSize(parse_u32(payload).context("invalid CTRL_ATTR_HDRSIZE value")?)
            }
            CTRL_ATTR_MAXATTR => {
                Self::MaxAttr(parse_u32(payload).context("invalid CTRL_ATTR_MAXATTR value")?)
            }
            CTRL_ATTR_OPS => {
                let ops = NlasIterator::new(payload)
                    .map(|nlas| {
                        nlas.and_then(|nlas| {
                            NlasIterator::new(nlas.value())
                                .map(|nla| nla.and_then(|nla| OpAttrs::parse(&nla)))
                                .collect::<Result<Vec<_>, _>>()
                        })
                    })
                    .collect::<Result<Vec<Vec<_>>, _>>()
                    .context("failed to parse CTRL_ATTR_OPS")?;

                Self::Ops(ops)
            }
            CTRL_ATTR_MCAST_GROUPS => {
                let groups = NlasIterator::new(payload)
                    .map(|nlas| {
                        nlas.and_then(|nlas| {
                            NlasIterator::new(nlas.value())
                                .map(|nla| nla.and_then(|nla| McastGrpAttrs::parse(&nla)))
                                .collect::<Result<Vec<_>, _>>()
                        })
                    })
                    .collect::<Result<Vec<Vec<_>>, _>>()
                    .context("failed to parse CTRL_ATTR_MCAST_GROUPS")?;

                Self::McastGroups(groups)
            }
            CTRL_ATTR_POLICY => Self::Policy(
                PolicyAttr::parse(&NlaBuffer::new(payload))
                    .context("failed to parse CTRL_ATTR_POLICY")?,
            ),
            CTRL_ATTR_OP_POLICY => Self::OpPolicy(
                OppolicyAttr::parse(&NlaBuffer::new(payload))
                    .context("failed to parse CTRL_ATTR_OP_POLICY")?,
            ),
            CTRL_ATTR_OP => Self::Op(parse_u32(payload)?),
            kind => return Err(DecodeError::from(format!("Unknown NLA type: {}", kind))),
        })
    }
}
