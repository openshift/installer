// SPDX-License-Identifier: MIT

use std::ffi::CString;

use anyhow::Context;
use byteorder::{ByteOrder, NativeEndian};
use netlink_packet_utils::{
    nla::{DefaultNla, Nla, NlaBuffer},
    parsers::{parse_string, parse_u32},
    DecodeError,
    Parseable,
};

const ALTIFNAMSIZ: usize = 128;
const ETHTOOL_A_HEADER_DEV_INDEX: u16 = 1;
const ETHTOOL_A_HEADER_DEV_NAME: u16 = 2;
const ETHTOOL_A_HEADER_FLAGS: u16 = 3;

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum EthtoolHeader {
    DevIndex(u32),
    DevName(String),
    Flags(u32),
    Other(DefaultNla),
}

impl Nla for EthtoolHeader {
    fn value_len(&self) -> usize {
        match self {
            Self::DevIndex(_) | Self::Flags(_) => 4,
            Self::DevName(s) => {
                if s.len() + 1 > ALTIFNAMSIZ {
                    ALTIFNAMSIZ
                } else {
                    s.len() + 1
                }
            }
            Self::Other(attr) => attr.value_len(),
        }
    }

    fn kind(&self) -> u16 {
        match self {
            Self::DevIndex(_) => ETHTOOL_A_HEADER_DEV_INDEX,
            Self::DevName(_) => ETHTOOL_A_HEADER_DEV_NAME,
            Self::Flags(_) => ETHTOOL_A_HEADER_FLAGS,
            Self::Other(attr) => attr.kind(),
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        match self {
            Self::DevIndex(value) | Self::Flags(value) => NativeEndian::write_u32(buffer, *value),
            Self::DevName(s) => str_to_zero_ended_u8_array(s, buffer, ALTIFNAMSIZ),
            Self::Other(ref attr) => attr.emit_value(buffer),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for EthtoolHeader {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            ETHTOOL_A_HEADER_DEV_INDEX => Self::DevIndex(
                parse_u32(payload).context("invalid ETHTOOL_A_HEADER_DEV_INDEX value")?,
            ),
            ETHTOOL_A_HEADER_FLAGS => {
                Self::Flags(parse_u32(payload).context("invalid ETHTOOL_A_HEADER_FLAGS value")?)
            }
            ETHTOOL_A_HEADER_DEV_NAME => Self::DevName(
                parse_string(payload).context("invalid ETHTOOL_A_HEADER_DEV_NAME value")?,
            ),
            _ => Self::Other(DefaultNla::parse(buf).context("invalid NLA (unknown kind)")?),
        })
    }
}

fn str_to_zero_ended_u8_array(src_str: &str, buffer: &mut [u8], max_size: usize) {
    if let Ok(src_cstring) = CString::new(src_str.as_bytes()) {
        let src_null_ended_str = src_cstring.into_bytes_with_nul();
        if src_null_ended_str.len() > max_size {
            buffer[..max_size].clone_from_slice(&src_null_ended_str[..max_size])
        } else {
            buffer[..src_null_ended_str.len()].clone_from_slice(&src_null_ended_str)
        }
    }
}
