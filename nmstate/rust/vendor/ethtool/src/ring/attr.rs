// SPDX-License-Identifier: MIT

use anyhow::Context;
use byteorder::{ByteOrder, NativeEndian};
use netlink_packet_utils::{
    nla::{DefaultNla, Nla, NlaBuffer, NlasIterator, NLA_F_NESTED},
    parsers::parse_u32,
    DecodeError,
    Emitable,
    Parseable,
};

use crate::{EthtoolAttr, EthtoolHeader};

const ETHTOOL_A_RINGS_HEADER: u16 = 1;
const ETHTOOL_A_RINGS_RX_MAX: u16 = 2;
const ETHTOOL_A_RINGS_RX_MINI_MAX: u16 = 3;
const ETHTOOL_A_RINGS_RX_JUMBO_MAX: u16 = 4;
const ETHTOOL_A_RINGS_TX_MAX: u16 = 5;
const ETHTOOL_A_RINGS_RX: u16 = 6;
const ETHTOOL_A_RINGS_RX_MINI: u16 = 7;
const ETHTOOL_A_RINGS_RX_JUMBO: u16 = 8;
const ETHTOOL_A_RINGS_TX: u16 = 9;

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum EthtoolRingAttr {
    Header(Vec<EthtoolHeader>),
    RxMax(u32),
    RxMiniMax(u32),
    RxJumboMax(u32),
    TxMax(u32),
    Rx(u32),
    RxMini(u32),
    RxJumbo(u32),
    Tx(u32),
    Other(DefaultNla),
}

impl Nla for EthtoolRingAttr {
    fn value_len(&self) -> usize {
        match self {
            Self::Header(hdrs) => hdrs.as_slice().buffer_len(),
            Self::RxMax(_)
            | Self::RxMiniMax(_)
            | Self::RxJumboMax(_)
            | Self::TxMax(_)
            | Self::Rx(_)
            | Self::RxMini(_)
            | Self::RxJumbo(_)
            | Self::Tx(_) => 4,
            Self::Other(attr) => attr.value_len(),
        }
    }

    fn kind(&self) -> u16 {
        match self {
            Self::Header(_) => ETHTOOL_A_RINGS_HEADER | NLA_F_NESTED,
            Self::RxMax(_) => ETHTOOL_A_RINGS_RX_MAX,
            Self::RxMiniMax(_) => ETHTOOL_A_RINGS_RX_MINI_MAX,
            Self::RxJumboMax(_) => ETHTOOL_A_RINGS_RX_JUMBO_MAX,
            Self::TxMax(_) => ETHTOOL_A_RINGS_TX_MAX,
            Self::Rx(_) => ETHTOOL_A_RINGS_RX,
            Self::RxMini(_) => ETHTOOL_A_RINGS_RX_MINI,
            Self::RxJumbo(_) => ETHTOOL_A_RINGS_RX_JUMBO,
            Self::Tx(_) => ETHTOOL_A_RINGS_TX,
            Self::Other(attr) => attr.kind(),
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        match self {
            Self::Header(ref nlas) => nlas.as_slice().emit(buffer),
            Self::RxMax(d)
            | Self::RxMiniMax(d)
            | Self::RxJumboMax(d)
            | Self::TxMax(d)
            | Self::Rx(d)
            | Self::RxMini(d)
            | Self::RxJumbo(d)
            | Self::Tx(d) => NativeEndian::write_u32(buffer, *d),
            Self::Other(ref attr) => attr.emit(buffer),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for EthtoolRingAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            ETHTOOL_A_RINGS_HEADER => {
                let mut nlas = Vec::new();
                let error_msg = "failed to parse ring header attributes";
                for nla in NlasIterator::new(payload) {
                    let nla = &nla.context(error_msg)?;
                    let parsed = EthtoolHeader::parse(nla).context(error_msg)?;
                    nlas.push(parsed);
                }
                Self::Header(nlas)
            }
            ETHTOOL_A_RINGS_RX_MAX => {
                Self::RxMax(parse_u32(payload).context("Invalid ETHTOOL_A_RINGS_RX_MAX value")?)
            }

            ETHTOOL_A_RINGS_RX_MINI_MAX => Self::RxMiniMax(
                parse_u32(payload).context("Invalid ETHTOOL_A_RINGS_RX_MINI_MAX value")?,
            ),
            ETHTOOL_A_RINGS_RX_JUMBO_MAX => Self::RxJumboMax(
                parse_u32(payload).context("Invalid ETHTOOL_A_RINGS_RX_JUMBO_MAX value")?,
            ),
            ETHTOOL_A_RINGS_TX_MAX => {
                Self::TxMax(parse_u32(payload).context("Invalid ETHTOOL_A_RINGS_TX_MAX value")?)
            }
            ETHTOOL_A_RINGS_RX => {
                Self::Rx(parse_u32(payload).context("Invalid ETHTOOL_A_RINGS_RX value")?)
            }
            ETHTOOL_A_RINGS_RX_MINI => {
                Self::RxMini(parse_u32(payload).context("Invalid ETHTOOL_A_RINGS_RX_MINI value")?)
            }
            ETHTOOL_A_RINGS_RX_JUMBO => {
                Self::RxJumbo(parse_u32(payload).context("Invalid ETHTOOL_A_RINGS_RX_JUMBO value")?)
            }
            ETHTOOL_A_RINGS_TX => {
                Self::Tx(parse_u32(payload).context("Invalid ETHTOOL_A_RINGS_TX value")?)
            }
            _ => Self::Other(DefaultNla::parse(buf).context("invalid NLA (unknown kind)")?),
        })
    }
}

pub(crate) fn parse_ring_nlas(buffer: &[u8]) -> Result<Vec<EthtoolAttr>, DecodeError> {
    let mut nlas = Vec::new();
    for nla in NlasIterator::new(buffer) {
        let error_msg = format!("Failed to parse ethtool ring message attribute {:?}", nla);
        let nla = &nla.context(error_msg.clone())?;
        let parsed = EthtoolRingAttr::parse(nla).context(error_msg)?;
        nlas.push(EthtoolAttr::Ring(parsed));
    }
    Ok(nlas)
}
