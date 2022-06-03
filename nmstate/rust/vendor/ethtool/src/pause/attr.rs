// SPDX-License-Identifier: MIT

use anyhow::Context;
use byteorder::{ByteOrder, NativeEndian};
use netlink_packet_utils::{
    nla::{DefaultNla, Nla, NlaBuffer, NlasIterator, NLA_F_NESTED},
    parsers::{parse_u64, parse_u8},
    DecodeError,
    Emitable,
    Parseable,
};

use crate::{EthtoolAttr, EthtoolHeader};

const ETHTOOL_A_PAUSE_HEADER: u16 = 1;
const ETHTOOL_A_PAUSE_AUTONEG: u16 = 2;
const ETHTOOL_A_PAUSE_RX: u16 = 3;
const ETHTOOL_A_PAUSE_TX: u16 = 4;
const ETHTOOL_A_PAUSE_STATS: u16 = 5;

const ETHTOOL_A_PAUSE_STAT_TX_FRAMES: u16 = 2;
const ETHTOOL_A_PAUSE_STAT_RX_FRAMES: u16 = 3;

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum EthtoolPauseStatAttr {
    Rx(u64),
    Tx(u64),
    Other(DefaultNla),
}

impl Nla for EthtoolPauseStatAttr {
    fn value_len(&self) -> usize {
        match self {
            Self::Rx(_) | Self::Tx(_) => 8,
            Self::Other(attr) => attr.value_len(),
        }
    }

    fn kind(&self) -> u16 {
        match self {
            Self::Rx(_) => ETHTOOL_A_PAUSE_STAT_RX_FRAMES,
            Self::Tx(_) => ETHTOOL_A_PAUSE_STAT_RX_FRAMES,
            Self::Other(attr) => attr.kind(),
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        match self {
            Self::Rx(value) | Self::Tx(value) => NativeEndian::write_u64(buffer, *value),
            Self::Other(ref attr) => attr.emit_value(buffer),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for EthtoolPauseStatAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            ETHTOOL_A_PAUSE_STAT_TX_FRAMES => Self::Tx(
                parse_u64(payload).context("invalid ETHTOOL_A_PAUSE_STAT_TX_FRAMES value")?,
            ),
            ETHTOOL_A_PAUSE_STAT_RX_FRAMES => Self::Rx(
                parse_u64(payload).context("invalid ETHTOOL_A_PAUSE_STAT_RX_FRAMES value")?,
            ),
            _ => Self::Other(DefaultNla::parse(buf).context("invalid NLA (unknown kind)")?),
        })
    }
}

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum EthtoolPauseAttr {
    Header(Vec<EthtoolHeader>),
    AutoNeg(bool),
    Rx(bool),
    Tx(bool),
    Stats(Vec<EthtoolPauseStatAttr>),
    Other(DefaultNla),
}

impl Nla for EthtoolPauseAttr {
    fn value_len(&self) -> usize {
        match self {
            Self::Header(hdrs) => hdrs.as_slice().buffer_len(),
            Self::AutoNeg(_) | Self::Rx(_) | Self::Tx(_) => 1,
            Self::Stats(ref nlas) => nlas.as_slice().buffer_len(),
            Self::Other(attr) => attr.value_len(),
        }
    }

    fn kind(&self) -> u16 {
        match self {
            Self::Header(_) => ETHTOOL_A_PAUSE_HEADER | NLA_F_NESTED,
            Self::AutoNeg(_) => ETHTOOL_A_PAUSE_AUTONEG,
            Self::Rx(_) => ETHTOOL_A_PAUSE_RX,
            Self::Tx(_) => ETHTOOL_A_PAUSE_TX,
            Self::Stats(_) => ETHTOOL_A_PAUSE_STATS | NLA_F_NESTED,
            Self::Other(attr) => attr.kind(),
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        match self {
            Self::Header(ref nlas) => nlas.as_slice().emit(buffer),
            Self::AutoNeg(value) | Self::Rx(value) | Self::Tx(value) => buffer[0] = *value as u8,
            Self::Stats(ref nlas) => nlas.as_slice().emit(buffer),
            Self::Other(ref attr) => attr.emit(buffer),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for EthtoolPauseAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            ETHTOOL_A_PAUSE_HEADER => {
                let mut nlas = Vec::new();
                let error_msg = "failed to parse pause header attributes";
                for nla in NlasIterator::new(payload) {
                    let nla = &nla.context(error_msg)?;
                    let parsed = EthtoolHeader::parse(nla).context(error_msg)?;
                    nlas.push(parsed);
                }
                Self::Header(nlas)
            }
            ETHTOOL_A_PAUSE_AUTONEG => Self::AutoNeg(
                parse_u8(payload).context("invalid ETHTOOL_A_PAUSE_AUTONEG value")? == 1,
            ),
            ETHTOOL_A_PAUSE_RX => {
                Self::Rx(parse_u8(payload).context("invalid ETHTOOL_A_PAUSE_RX value")? == 1)
            }
            ETHTOOL_A_PAUSE_TX => {
                Self::Tx(parse_u8(payload).context("invalid ETHTOOL_A_PAUSE_TX value")? == 1)
            }
            ETHTOOL_A_PAUSE_STATS => {
                let mut nlas = Vec::new();
                let error_msg = "failed to parse pause stats attributes";
                for nla in NlasIterator::new(payload) {
                    let nla = &nla.context(error_msg)?;
                    let parsed = EthtoolPauseStatAttr::parse(nla).context(error_msg)?;
                    nlas.push(parsed);
                }
                Self::Stats(nlas)
            }
            _ => Self::Other(DefaultNla::parse(buf).context("invalid NLA (unknown kind)")?),
        })
    }
}

pub(crate) fn parse_pause_nlas(buffer: &[u8]) -> Result<Vec<EthtoolAttr>, DecodeError> {
    let mut nlas = Vec::new();
    for nla in NlasIterator::new(buffer) {
        let error_msg = format!("Failed to parse ethtool pause message attribute {:?}", nla);
        let nla = &nla.context(error_msg.clone())?;
        let parsed = EthtoolPauseAttr::parse(nla).context(error_msg)?;
        nlas.push(EthtoolAttr::Pause(parsed));
    }
    Ok(nlas)
}
