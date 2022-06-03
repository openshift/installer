// SPDX-License-Identifier: MIT

use anyhow::Context;
use byteorder::{ByteOrder, NativeEndian};
use netlink_packet_utils::{
    nla::{DefaultNla, Nla, NlaBuffer, NlasIterator, NLA_F_NESTED},
    parsers::{parse_u32, parse_u8},
    DecodeError,
    Emitable,
    Parseable,
};

use crate::{EthtoolAttr, EthtoolHeader};

const ETHTOOL_A_COALESCE_HEADER: u16 = 1;
const ETHTOOL_A_COALESCE_RX_USECS: u16 = 2;
const ETHTOOL_A_COALESCE_RX_MAX_FRAMES: u16 = 3;
const ETHTOOL_A_COALESCE_RX_USECS_IRQ: u16 = 4;
const ETHTOOL_A_COALESCE_RX_MAX_FRAMES_IRQ: u16 = 5;
const ETHTOOL_A_COALESCE_TX_USECS: u16 = 6;
const ETHTOOL_A_COALESCE_TX_MAX_FRAMES: u16 = 7;
const ETHTOOL_A_COALESCE_TX_USECS_IRQ: u16 = 8;
const ETHTOOL_A_COALESCE_TX_MAX_FRAMES_IRQ: u16 = 9;
const ETHTOOL_A_COALESCE_STATS_BLOCK_USECS: u16 = 10;
const ETHTOOL_A_COALESCE_USE_ADAPTIVE_RX: u16 = 11;
const ETHTOOL_A_COALESCE_USE_ADAPTIVE_TX: u16 = 12;
const ETHTOOL_A_COALESCE_PKT_RATE_LOW: u16 = 13;
const ETHTOOL_A_COALESCE_RX_USECS_LOW: u16 = 14;
const ETHTOOL_A_COALESCE_RX_MAX_FRAMES_LOW: u16 = 15;
const ETHTOOL_A_COALESCE_TX_USECS_LOW: u16 = 16;
const ETHTOOL_A_COALESCE_TX_MAX_FRAMES_LOW: u16 = 17;
const ETHTOOL_A_COALESCE_PKT_RATE_HIGH: u16 = 18;
const ETHTOOL_A_COALESCE_RX_USECS_HIGH: u16 = 19;
const ETHTOOL_A_COALESCE_RX_MAX_FRAMES_HIGH: u16 = 20;
const ETHTOOL_A_COALESCE_TX_USECS_HIGH: u16 = 21;
const ETHTOOL_A_COALESCE_TX_MAX_FRAMES_HIGH: u16 = 22;
const ETHTOOL_A_COALESCE_RATE_SAMPLE_INTERVAL: u16 = 23;

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum EthtoolCoalesceAttr {
    Header(Vec<EthtoolHeader>),
    RxUsecs(u32),
    RxMaxFrames(u32),
    RxUsecsIrq(u32),
    RxMaxFramesIrq(u32),
    TxUsecs(u32),
    TxMaxFrames(u32),
    TxUsecsIrq(u32),
    TxMaxFramesIrq(u32),
    StatsBlockUsecs(u32),
    UseAdaptiveRx(bool),
    UseAdaptiveTx(bool),
    PktRateLow(u32),
    RxUsecsLow(u32),
    RxMaxFramesLow(u32),
    TxUsecsLow(u32),
    TxMaxFramesLow(u32),
    PktRateHigh(u32),
    RxUsecsHigh(u32),
    RxMaxFramesHigh(u32),
    TxUsecsHigh(u32),
    TxMaxFramesHigh(u32),
    RateSampleInterval(u32),
    Other(DefaultNla),
}

impl Nla for EthtoolCoalesceAttr {
    fn value_len(&self) -> usize {
        match self {
            Self::Header(hdrs) => hdrs.as_slice().buffer_len(),
            Self::RxUsecs(_)
            | Self::RxMaxFrames(_)
            | Self::RxUsecsIrq(_)
            | Self::RxMaxFramesIrq(_)
            | Self::TxUsecs(_)
            | Self::TxMaxFrames(_)
            | Self::TxUsecsIrq(_)
            | Self::TxMaxFramesIrq(_)
            | Self::StatsBlockUsecs(_)
            | Self::PktRateLow(_)
            | Self::RxUsecsLow(_)
            | Self::RxMaxFramesLow(_)
            | Self::TxUsecsLow(_)
            | Self::TxMaxFramesLow(_)
            | Self::PktRateHigh(_)
            | Self::RxUsecsHigh(_)
            | Self::RxMaxFramesHigh(_)
            | Self::TxUsecsHigh(_)
            | Self::TxMaxFramesHigh(_)
            | Self::RateSampleInterval(_) => 4,
            Self::UseAdaptiveRx(_) | Self::UseAdaptiveTx(_) => 1,
            Self::Other(attr) => attr.value_len(),
        }
    }

    fn kind(&self) -> u16 {
        match self {
            Self::Header(_) => ETHTOOL_A_COALESCE_HEADER | NLA_F_NESTED,
            Self::RxUsecs(_) => ETHTOOL_A_COALESCE_RX_USECS,
            Self::RxMaxFrames(_) => ETHTOOL_A_COALESCE_RX_MAX_FRAMES,
            Self::RxUsecsIrq(_) => ETHTOOL_A_COALESCE_RX_USECS_IRQ,
            Self::RxMaxFramesIrq(_) => ETHTOOL_A_COALESCE_RX_MAX_FRAMES_IRQ,
            Self::TxUsecs(_) => ETHTOOL_A_COALESCE_TX_USECS,
            Self::TxMaxFrames(_) => ETHTOOL_A_COALESCE_TX_MAX_FRAMES,
            Self::TxUsecsIrq(_) => ETHTOOL_A_COALESCE_TX_USECS_IRQ,
            Self::TxMaxFramesIrq(_) => ETHTOOL_A_COALESCE_TX_MAX_FRAMES_IRQ,
            Self::StatsBlockUsecs(_) => ETHTOOL_A_COALESCE_STATS_BLOCK_USECS,
            Self::UseAdaptiveRx(_) => ETHTOOL_A_COALESCE_USE_ADAPTIVE_RX,
            Self::UseAdaptiveTx(_) => ETHTOOL_A_COALESCE_USE_ADAPTIVE_TX,
            Self::PktRateLow(_) => ETHTOOL_A_COALESCE_PKT_RATE_LOW,
            Self::RxUsecsLow(_) => ETHTOOL_A_COALESCE_RX_USECS_LOW,
            Self::RxMaxFramesLow(_) => ETHTOOL_A_COALESCE_RX_MAX_FRAMES_LOW,
            Self::TxUsecsLow(_) => ETHTOOL_A_COALESCE_TX_USECS_LOW,
            Self::TxMaxFramesLow(_) => ETHTOOL_A_COALESCE_TX_MAX_FRAMES_LOW,
            Self::PktRateHigh(_) => ETHTOOL_A_COALESCE_PKT_RATE_HIGH,
            Self::RxUsecsHigh(_) => ETHTOOL_A_COALESCE_RX_USECS_HIGH,
            Self::RxMaxFramesHigh(_) => ETHTOOL_A_COALESCE_RX_MAX_FRAMES_HIGH,
            Self::TxUsecsHigh(_) => ETHTOOL_A_COALESCE_TX_USECS_HIGH,
            Self::TxMaxFramesHigh(_) => ETHTOOL_A_COALESCE_TX_MAX_FRAMES_HIGH,
            Self::RateSampleInterval(_) => ETHTOOL_A_COALESCE_RATE_SAMPLE_INTERVAL,
            Self::Other(attr) => attr.kind(),
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        match self {
            Self::Header(ref nlas) => nlas.as_slice().emit(buffer),
            Self::Other(ref attr) => attr.emit(buffer),
            Self::RxUsecs(d)
            | Self::RxMaxFrames(d)
            | Self::RxUsecsIrq(d)
            | Self::RxMaxFramesIrq(d)
            | Self::TxUsecs(d)
            | Self::TxMaxFrames(d)
            | Self::TxUsecsIrq(d)
            | Self::TxMaxFramesIrq(d)
            | Self::StatsBlockUsecs(d)
            | Self::PktRateLow(d)
            | Self::RxUsecsLow(d)
            | Self::RxMaxFramesLow(d)
            | Self::TxUsecsLow(d)
            | Self::TxMaxFramesLow(d)
            | Self::PktRateHigh(d)
            | Self::RxUsecsHigh(d)
            | Self::RxMaxFramesHigh(d)
            | Self::TxUsecsHigh(d)
            | Self::TxMaxFramesHigh(d)
            | Self::RateSampleInterval(d) => NativeEndian::write_u32(buffer, *d),
            Self::UseAdaptiveRx(d) | Self::UseAdaptiveTx(d) => buffer[0] = if *d { 1 } else { 0 },
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for EthtoolCoalesceAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            ETHTOOL_A_COALESCE_HEADER => {
                let mut nlas = Vec::new();
                let error_msg = "failed to parse coalesce header attributes";
                for nla in NlasIterator::new(payload) {
                    let nla = &nla.context(error_msg)?;
                    let parsed = EthtoolHeader::parse(nla).context(error_msg)?;
                    nlas.push(parsed);
                }
                Self::Header(nlas)
            }
            ETHTOOL_A_COALESCE_RX_USECS => Self::RxUsecs(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_RX_USECS value")?,
            ),
            ETHTOOL_A_COALESCE_RX_MAX_FRAMES => Self::RxMaxFrames(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_RX_MAX_FRAMES value")?,
            ),

            ETHTOOL_A_COALESCE_RX_USECS_IRQ => Self::RxUsecsIrq(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_RX_USECS_IRQ value")?,
            ),

            ETHTOOL_A_COALESCE_RX_MAX_FRAMES_IRQ => Self::RxMaxFramesIrq(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_RX_MAX_FRAMES_IRQ value")?,
            ),

            ETHTOOL_A_COALESCE_TX_USECS => Self::TxUsecs(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_TX_USECS value")?,
            ),

            ETHTOOL_A_COALESCE_TX_MAX_FRAMES => Self::TxMaxFrames(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_TX_MAX_FRAMES value")?,
            ),

            ETHTOOL_A_COALESCE_TX_USECS_IRQ => Self::TxUsecsIrq(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_TX_USECS_IRQ value")?,
            ),

            ETHTOOL_A_COALESCE_TX_MAX_FRAMES_IRQ => Self::TxMaxFramesIrq(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_TX_MAX_FRAMES_IRQ value")?,
            ),

            ETHTOOL_A_COALESCE_STATS_BLOCK_USECS => Self::StatsBlockUsecs(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_STATS_BLOCK_USECS value")?,
            ),

            ETHTOOL_A_COALESCE_USE_ADAPTIVE_RX => Self::UseAdaptiveRx(
                parse_u8(payload).context("Invalid ETHTOOL_A_COALESCE_USE_ADAPTIVE_RX value")? == 1,
            ),

            ETHTOOL_A_COALESCE_USE_ADAPTIVE_TX => Self::UseAdaptiveTx(
                parse_u8(payload).context("Invalid ETHTOOL_A_COALESCE_USE_ADAPTIVE_TX value")? == 1,
            ),

            ETHTOOL_A_COALESCE_PKT_RATE_LOW => Self::PktRateLow(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_PKT_RATE_LOW value")?,
            ),

            ETHTOOL_A_COALESCE_RX_USECS_LOW => Self::RxUsecsLow(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_RX_USECS_LOW value")?,
            ),

            ETHTOOL_A_COALESCE_RX_MAX_FRAMES_LOW => Self::RxMaxFramesLow(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_RX_MAX_FRAMES_LOW value")?,
            ),

            ETHTOOL_A_COALESCE_TX_USECS_LOW => Self::TxUsecsLow(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_TX_USECS_LOW value")?,
            ),

            ETHTOOL_A_COALESCE_TX_MAX_FRAMES_LOW => Self::TxMaxFramesLow(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_TX_MAX_FRAMES_LOW value")?,
            ),

            ETHTOOL_A_COALESCE_PKT_RATE_HIGH => Self::PktRateHigh(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_PKT_RATE_HIGH value")?,
            ),

            ETHTOOL_A_COALESCE_RX_USECS_HIGH => Self::RxUsecsHigh(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_RX_USECS_HIGH value")?,
            ),

            ETHTOOL_A_COALESCE_RX_MAX_FRAMES_HIGH => Self::RxMaxFramesHigh(
                parse_u32(payload)
                    .context("Invalid ETHTOOL_A_COALESCE_RX_MAX_FRAMES_HIGH value")?,
            ),

            ETHTOOL_A_COALESCE_TX_USECS_HIGH => Self::TxUsecsHigh(
                parse_u32(payload).context("Invalid ETHTOOL_A_COALESCE_TX_USECS_HIGH value")?,
            ),

            ETHTOOL_A_COALESCE_TX_MAX_FRAMES_HIGH => Self::TxMaxFramesHigh(
                parse_u32(payload)
                    .context("Invalid ETHTOOL_A_COALESCE_TX_MAX_FRAMES_HIGH value")?,
            ),

            ETHTOOL_A_COALESCE_RATE_SAMPLE_INTERVAL => Self::RateSampleInterval(
                parse_u32(payload)
                    .context("Invalid ETHTOOL_A_COALESCE_RATE_SAMPLE_INTERVAL value")?,
            ),

            _ => Self::Other(DefaultNla::parse(buf).context("invalid NLA (unknown kind)")?),
        })
    }
}

pub(crate) fn parse_coalesce_nlas(buffer: &[u8]) -> Result<Vec<EthtoolAttr>, DecodeError> {
    let mut nlas = Vec::new();
    for nla in NlasIterator::new(buffer) {
        let error_msg = format!(
            "Failed to parse ethtool coalesce message attribute {:?}",
            nla
        );
        let nla = &nla.context(error_msg.clone())?;
        let parsed = EthtoolCoalesceAttr::parse(nla).context(error_msg)?;
        nlas.push(EthtoolAttr::Coalesce(parsed));
    }
    Ok(nlas)
}
