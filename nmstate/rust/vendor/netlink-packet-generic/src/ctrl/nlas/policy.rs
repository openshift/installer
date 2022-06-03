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
use std::{
    convert::TryFrom,
    mem::{size_of, size_of_val},
};

// PolicyAttr

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct PolicyAttr {
    pub index: u16,
    pub attr_policy: AttributePolicyAttr,
}

impl Nla for PolicyAttr {
    fn value_len(&self) -> usize {
        self.attr_policy.buffer_len()
    }

    fn kind(&self) -> u16 {
        self.index
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        self.attr_policy.emit(buffer);
    }

    fn is_nested(&self) -> bool {
        true
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for PolicyAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();

        Ok(Self {
            index: buf.kind(),
            attr_policy: AttributePolicyAttr::parse(&NlaBuffer::new(payload))
                .context("failed to parse PolicyAttr")?,
        })
    }
}

// AttributePolicyAttr

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct AttributePolicyAttr {
    pub index: u16,
    pub policies: Vec<NlPolicyTypeAttrs>,
}

impl Nla for AttributePolicyAttr {
    fn value_len(&self) -> usize {
        self.policies.as_slice().buffer_len()
    }

    fn kind(&self) -> u16 {
        self.index
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        self.policies.as_slice().emit(buffer);
    }

    fn is_nested(&self) -> bool {
        true
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for AttributePolicyAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        let policies = NlasIterator::new(payload)
            .map(|nla| nla.and_then(|nla| NlPolicyTypeAttrs::parse(&nla)))
            .collect::<Result<Vec<_>, _>>()
            .context("failed to parse AttributePolicyAttr")?;

        Ok(Self {
            index: buf.kind(),
            policies,
        })
    }
}

// PolicyTypeAttrs

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum NlPolicyTypeAttrs {
    Type(NlaType),
    MinValueSigned(i64),
    MaxValueSigned(i64),
    MaxValueUnsigned(u64),
    MinValueUnsigned(u64),
    MinLength(u32),
    MaxLength(u32),
    PolicyIdx(u32),
    PolicyMaxType(u32),
    Bitfield32Mask(u32),
    Mask(u64),
}

impl Nla for NlPolicyTypeAttrs {
    fn value_len(&self) -> usize {
        use NlPolicyTypeAttrs::*;
        match self {
            Type(v) => size_of_val(v),
            MinValueSigned(v) => size_of_val(v),
            MaxValueSigned(v) => size_of_val(v),
            MaxValueUnsigned(v) => size_of_val(v),
            MinValueUnsigned(v) => size_of_val(v),
            MinLength(v) => size_of_val(v),
            MaxLength(v) => size_of_val(v),
            PolicyIdx(v) => size_of_val(v),
            PolicyMaxType(v) => size_of_val(v),
            Bitfield32Mask(v) => size_of_val(v),
            Mask(v) => size_of_val(v),
        }
    }

    fn kind(&self) -> u16 {
        use NlPolicyTypeAttrs::*;
        match self {
            Type(_) => NL_POLICY_TYPE_ATTR_TYPE,
            MinValueSigned(_) => NL_POLICY_TYPE_ATTR_MIN_VALUE_S,
            MaxValueSigned(_) => NL_POLICY_TYPE_ATTR_MAX_VALUE_S,
            MaxValueUnsigned(_) => NL_POLICY_TYPE_ATTR_MIN_VALUE_U,
            MinValueUnsigned(_) => NL_POLICY_TYPE_ATTR_MAX_VALUE_U,
            MinLength(_) => NL_POLICY_TYPE_ATTR_MIN_LENGTH,
            MaxLength(_) => NL_POLICY_TYPE_ATTR_MAX_LENGTH,
            PolicyIdx(_) => NL_POLICY_TYPE_ATTR_POLICY_IDX,
            PolicyMaxType(_) => NL_POLICY_TYPE_ATTR_POLICY_MAXTYPE,
            Bitfield32Mask(_) => NL_POLICY_TYPE_ATTR_BITFIELD32_MASK,
            Mask(_) => NL_POLICY_TYPE_ATTR_MASK,
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        use NlPolicyTypeAttrs::*;
        match self {
            Type(v) => NativeEndian::write_u32(buffer, u32::from(*v)),
            MinValueSigned(v) => NativeEndian::write_i64(buffer, *v),
            MaxValueSigned(v) => NativeEndian::write_i64(buffer, *v),
            MaxValueUnsigned(v) => NativeEndian::write_u64(buffer, *v),
            MinValueUnsigned(v) => NativeEndian::write_u64(buffer, *v),
            MinLength(v) => NativeEndian::write_u32(buffer, *v),
            MaxLength(v) => NativeEndian::write_u32(buffer, *v),
            PolicyIdx(v) => NativeEndian::write_u32(buffer, *v),
            PolicyMaxType(v) => NativeEndian::write_u32(buffer, *v),
            Bitfield32Mask(v) => NativeEndian::write_u32(buffer, *v),
            Mask(v) => NativeEndian::write_u64(buffer, *v),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for NlPolicyTypeAttrs {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            NL_POLICY_TYPE_ATTR_TYPE => {
                let value = parse_u32(payload).context("invalid NL_POLICY_TYPE_ATTR_TYPE value")?;
                Self::Type(NlaType::try_from(value)?)
            }
            NL_POLICY_TYPE_ATTR_MIN_VALUE_S => Self::MinValueSigned(
                parse_i64(payload).context("invalid NL_POLICY_TYPE_ATTR_MIN_VALUE_S value")?,
            ),
            NL_POLICY_TYPE_ATTR_MAX_VALUE_S => Self::MaxValueSigned(
                parse_i64(payload).context("invalid NL_POLICY_TYPE_ATTR_MAX_VALUE_S value")?,
            ),
            NL_POLICY_TYPE_ATTR_MIN_VALUE_U => Self::MinValueUnsigned(
                parse_u64(payload).context("invalid NL_POLICY_TYPE_ATTR_MIN_VALUE_U value")?,
            ),
            NL_POLICY_TYPE_ATTR_MAX_VALUE_U => Self::MaxValueUnsigned(
                parse_u64(payload).context("invalid NL_POLICY_TYPE_ATTR_MAX_VALUE_U value")?,
            ),
            NL_POLICY_TYPE_ATTR_MIN_LENGTH => Self::MinLength(
                parse_u32(payload).context("invalid NL_POLICY_TYPE_ATTR_MIN_LENGTH value")?,
            ),
            NL_POLICY_TYPE_ATTR_MAX_LENGTH => Self::MaxLength(
                parse_u32(payload).context("invalid NL_POLICY_TYPE_ATTR_MAX_LENGTH value")?,
            ),
            NL_POLICY_TYPE_ATTR_POLICY_IDX => Self::PolicyIdx(
                parse_u32(payload).context("invalid NL_POLICY_TYPE_ATTR_POLICY_IDX value")?,
            ),
            NL_POLICY_TYPE_ATTR_POLICY_MAXTYPE => Self::PolicyMaxType(
                parse_u32(payload).context("invalid NL_POLICY_TYPE_ATTR_POLICY_MAXTYPE value")?,
            ),
            NL_POLICY_TYPE_ATTR_BITFIELD32_MASK => Self::Bitfield32Mask(
                parse_u32(payload).context("invalid NL_POLICY_TYPE_ATTR_BITFIELD32_MASK value")?,
            ),
            NL_POLICY_TYPE_ATTR_MASK => {
                Self::Mask(parse_u64(payload).context("invalid NL_POLICY_TYPE_ATTR_MASK value")?)
            }
            kind => return Err(DecodeError::from(format!("Unknown NLA type: {}", kind))),
        })
    }
}

#[derive(Copy, Clone, Debug, PartialEq, Eq)]
pub enum NlaType {
    Flag,
    U8,
    U16,
    U32,
    U64,
    S8,
    S16,
    S32,
    S64,
    Binary,
    String,
    NulString,
    Nested,
    NestedArray,
    Bitfield32,
}

impl From<NlaType> for u32 {
    fn from(nlatype: NlaType) -> u32 {
        match nlatype {
            NlaType::Flag => NL_ATTR_TYPE_FLAG,
            NlaType::U8 => NL_ATTR_TYPE_U8,
            NlaType::U16 => NL_ATTR_TYPE_U16,
            NlaType::U32 => NL_ATTR_TYPE_U32,
            NlaType::U64 => NL_ATTR_TYPE_U64,
            NlaType::S8 => NL_ATTR_TYPE_S8,
            NlaType::S16 => NL_ATTR_TYPE_S16,
            NlaType::S32 => NL_ATTR_TYPE_S32,
            NlaType::S64 => NL_ATTR_TYPE_S64,
            NlaType::Binary => NL_ATTR_TYPE_BINARY,
            NlaType::String => NL_ATTR_TYPE_STRING,
            NlaType::NulString => NL_ATTR_TYPE_NUL_STRING,
            NlaType::Nested => NL_ATTR_TYPE_NESTED,
            NlaType::NestedArray => NL_ATTR_TYPE_NESTED_ARRAY,
            NlaType::Bitfield32 => NL_ATTR_TYPE_BITFIELD32,
        }
    }
}

impl TryFrom<u32> for NlaType {
    type Error = DecodeError;

    fn try_from(value: u32) -> Result<Self, Self::Error> {
        Ok(match value {
            NL_ATTR_TYPE_FLAG => NlaType::Flag,
            NL_ATTR_TYPE_U8 => NlaType::U8,
            NL_ATTR_TYPE_U16 => NlaType::U16,
            NL_ATTR_TYPE_U32 => NlaType::U32,
            NL_ATTR_TYPE_U64 => NlaType::U64,
            NL_ATTR_TYPE_S8 => NlaType::S8,
            NL_ATTR_TYPE_S16 => NlaType::S16,
            NL_ATTR_TYPE_S32 => NlaType::S32,
            NL_ATTR_TYPE_S64 => NlaType::S64,
            NL_ATTR_TYPE_BINARY => NlaType::Binary,
            NL_ATTR_TYPE_STRING => NlaType::String,
            NL_ATTR_TYPE_NUL_STRING => NlaType::NulString,
            NL_ATTR_TYPE_NESTED => NlaType::Nested,
            NL_ATTR_TYPE_NESTED_ARRAY => NlaType::NestedArray,
            NL_ATTR_TYPE_BITFIELD32 => NlaType::Bitfield32,
            _ => return Err(DecodeError::from(format!("invalid NLA type: {}", value))),
        })
    }
}

// FIXME: Add this into netlink_packet_utils::parser
fn parse_i64(payload: &[u8]) -> Result<i64, DecodeError> {
    if payload.len() != size_of::<i64>() {
        return Err(format!("invalid i64: {:?}", payload).into());
    }
    Ok(NativeEndian::read_i64(payload))
}
