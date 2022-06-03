// SPDX-License-Identifier: MIT

use anyhow::Context;
use log::warn;
use netlink_packet_utils::{
    nla::{DefaultNla, Nla, NlaBuffer, NlasIterator, NLA_F_NESTED},
    parsers::{parse_string, parse_u32, parse_u8},
    DecodeError,
    Emitable,
    Parseable,
};

use crate::{EthtoolAttr, EthtoolHeader};

const ETHTOOL_A_LINKMODES_HEADER: u16 = 1;
const ETHTOOL_A_LINKMODES_AUTONEG: u16 = 2;
const ETHTOOL_A_LINKMODES_OURS: u16 = 3;
const ETHTOOL_A_LINKMODES_PEER: u16 = 4;
const ETHTOOL_A_LINKMODES_SPEED: u16 = 5;
const ETHTOOL_A_LINKMODES_DUPLEX: u16 = 6;
const ETHTOOL_A_LINKMODES_SUBORDINATE_CFG: u16 = 7;
const ETHTOOL_A_LINKMODES_SUBORDINATE_STATE: u16 = 8;
const ETHTOOL_A_LINKMODES_LANES: u16 = 9;

const ETHTOOL_A_BITSET_BITS: u16 = 3;

const ETHTOOL_A_BITSET_BITS_BIT: u16 = 1;

const ETHTOOL_A_BITSET_BIT_INDEX: u16 = 1;
const ETHTOOL_A_BITSET_BIT_NAME: u16 = 2;
const ETHTOOL_A_BITSET_BIT_VALUE: u16 = 3;

const DUPLEX_HALF: u8 = 0x00;
const DUPLEX_FULL: u8 = 0x01;
const DUPLEX_UNKNOWN: u8 = 0xff;

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum EthtoolLinkModeDuplex {
    Half,
    Full,
    Unknown,
    Other(u8),
}

impl From<u8> for EthtoolLinkModeDuplex {
    fn from(d: u8) -> Self {
        match d {
            DUPLEX_HALF => Self::Half,
            DUPLEX_FULL => Self::Full,
            DUPLEX_UNKNOWN => Self::Unknown,
            _ => Self::Other(d),
        }
    }
}

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum EthtoolLinkModeAttr {
    Header(Vec<EthtoolHeader>),
    Autoneg(bool),
    Ours(Vec<String>),
    Peer(Vec<String>),
    Speed(u32),
    Duplex(EthtoolLinkModeDuplex),
    ControllerSubordinateCfg(u8),
    ControllerSubordinateState(u8),
    Lanes(u32),
    Other(DefaultNla),
}

impl Nla for EthtoolLinkModeAttr {
    fn value_len(&self) -> usize {
        match self {
            Self::Header(hdrs) => hdrs.as_slice().buffer_len(),
            Self::Autoneg(_)
            | Self::Duplex(_)
            | Self::ControllerSubordinateCfg(_)
            | Self::ControllerSubordinateState(_) => 1,
            Self::Ours(_) => {
                todo!("Does not support changing ethtool link mode yet")
            }
            Self::Peer(_) => {
                todo!("Does not support changing ethtool link mode yet")
            }
            Self::Speed(_) | Self::Lanes(_) => 4,
            Self::Other(attr) => attr.value_len(),
        }
    }

    fn kind(&self) -> u16 {
        match self {
            Self::Header(_) => ETHTOOL_A_LINKMODES_HEADER | NLA_F_NESTED,
            Self::Autoneg(_) => ETHTOOL_A_LINKMODES_AUTONEG,
            Self::Ours(_) => ETHTOOL_A_LINKMODES_OURS,
            Self::Peer(_) => ETHTOOL_A_LINKMODES_PEER,
            Self::Speed(_) => ETHTOOL_A_LINKMODES_SPEED,
            Self::Duplex(_) => ETHTOOL_A_LINKMODES_DUPLEX,
            Self::ControllerSubordinateCfg(_) => ETHTOOL_A_LINKMODES_SUBORDINATE_CFG,
            Self::ControllerSubordinateState(_) => ETHTOOL_A_LINKMODES_SUBORDINATE_STATE,
            Self::Lanes(_) => ETHTOOL_A_LINKMODES_LANES,
            Self::Other(attr) => attr.kind(),
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        match self {
            Self::Header(ref nlas) => nlas.as_slice().emit(buffer),
            Self::Other(ref attr) => attr.emit(buffer),
            _ => todo!("Does not support changing ethtool link mode yet"),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for EthtoolLinkModeAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            ETHTOOL_A_LINKMODES_HEADER => {
                let mut nlas = Vec::new();
                let error_msg = "failed to parse link_mode header attributes";
                for nla in NlasIterator::new(payload) {
                    let nla = &nla.context(error_msg)?;
                    let parsed = EthtoolHeader::parse(nla).context(error_msg)?;
                    nlas.push(parsed);
                }
                Self::Header(nlas)
            }
            ETHTOOL_A_LINKMODES_AUTONEG => Self::Autoneg(
                parse_u8(payload).context("Invalid ETHTOOL_A_LINKMODES_AUTONEG value")? == 1,
            ),

            ETHTOOL_A_LINKMODES_OURS => Self::Ours(parse_bitset_bits_nlas(payload)?),
            ETHTOOL_A_LINKMODES_PEER => Self::Peer(parse_bitset_bits_nlas(payload)?),
            ETHTOOL_A_LINKMODES_SPEED => {
                Self::Speed(parse_u32(payload).context("Invalid ETHTOOL_A_LINKMODES_SPEED value")?)
            }
            ETHTOOL_A_LINKMODES_DUPLEX => Self::Duplex(
                parse_u8(payload)
                    .context("Invalid ETHTOOL_A_LINKMODES_DUPLEX value")?
                    .into(),
            ),
            ETHTOOL_A_LINKMODES_SUBORDINATE_CFG => Self::ControllerSubordinateCfg(
                parse_u8(payload).context("Invalid ETHTOOL_A_LINKMODES_SUBORDINATE_CFG value")?,
            ),
            ETHTOOL_A_LINKMODES_SUBORDINATE_STATE => Self::ControllerSubordinateState(
                parse_u8(payload).context("Invalid ETHTOOL_A_LINKMODES_SUBORDINATE_STATE value")?,
            ),
            ETHTOOL_A_LINKMODES_LANES => {
                Self::Lanes(parse_u32(payload).context("Invalid ETHTOOL_A_LINKMODES_LANES value")?)
            }
            _ => Self::Other(DefaultNla::parse(buf).context("invalid NLA (unknown kind)")?),
        })
    }
}

fn parse_bitset_bits_nlas(raw: &[u8]) -> Result<Vec<String>, DecodeError> {
    let error_msg = "failed to parse mode bit sets";
    for nla in NlasIterator::new(raw) {
        let nla = &nla.context(error_msg)?;
        if nla.kind() == ETHTOOL_A_BITSET_BITS {
            return parse_bitset_bits_nla(nla.value());
        }
    }
    Err("No ETHTOOL_A_BITSET_BITS NLA found".into())
}

fn parse_bitset_bits_nla(raw: &[u8]) -> Result<Vec<String>, DecodeError> {
    let mut modes = Vec::new();
    let error_msg = "Failed to parse ETHTOOL_A_BITSET_BITS attributes";
    for bit_nla in NlasIterator::new(raw) {
        let bit_nla = &bit_nla.context(error_msg)?;
        match bit_nla.kind() {
            ETHTOOL_A_BITSET_BITS_BIT => {
                let error_msg = "Failed to parse ETHTOOL_A_BITSET_BITS_BIT attributes";
                let nlas = NlasIterator::new(bit_nla.value());
                for nla in nlas {
                    let nla = &nla.context(error_msg)?;
                    let payload = nla.value();
                    match nla.kind() {
                        ETHTOOL_A_BITSET_BIT_INDEX | ETHTOOL_A_BITSET_BIT_VALUE => {
                            // ignored
                        }
                        ETHTOOL_A_BITSET_BIT_NAME => {
                            modes.push(
                                parse_string(payload)
                                    .context("Invald ETHTOOL_A_BITSET_BIT_NAME value")?,
                            );
                        }
                        _ => {
                            warn!(
                                "Unknown ETHTOOL_A_BITSET_BITS_BIT {} {:?}",
                                nla.kind(),
                                nla.value(),
                            );
                        }
                    }
                }
            }
            _ => {
                warn!(
                    "Unknown ETHTOOL_A_BITSET_BITS kind {}, {:?}",
                    bit_nla.kind(),
                    bit_nla.value()
                );
            }
        };
    }
    Ok(modes)
}

pub(crate) fn parse_link_mode_nlas(buffer: &[u8]) -> Result<Vec<EthtoolAttr>, DecodeError> {
    let mut nlas = Vec::new();
    for nla in NlasIterator::new(buffer) {
        let error_msg = format!(
            "Failed to parse ethtool link_mode message attribute {:?}",
            nla
        );
        let nla = &nla.context(error_msg.clone())?;
        let parsed = EthtoolLinkModeAttr::parse(nla).context(error_msg)?;
        nlas.push(EthtoolAttr::LinkMode(parsed));
    }
    Ok(nlas)
}
