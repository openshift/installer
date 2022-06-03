// SPDX-License-Identifier: MIT

use anyhow::Context;
use log::warn;
use netlink_packet_utils::{
    nla::{DefaultNla, Nla, NlaBuffer, NlasIterator, NLA_F_NESTED},
    parsers::{parse_string, parse_u32},
    DecodeError,
    Emitable,
    Parseable,
};

use crate::{EthtoolAttr, EthtoolHeader};

const ETHTOOL_A_FEATURES_HEADER: u16 = 1;
const ETHTOOL_A_FEATURES_HW: u16 = 2; // User changable features
const ETHTOOL_A_FEATURES_WANTED: u16 = 3; // User requested fatures
const ETHTOOL_A_FEATURES_ACTIVE: u16 = 4; // Active features
const ETHTOOL_A_FEATURES_NOCHANGE: u16 = 5;

const ETHTOOL_A_BITSET_BITS: u16 = 3;

const ETHTOOL_A_BITSET_BITS_BIT: u16 = 1;

const ETHTOOL_A_BITSET_BIT_INDEX: u16 = 1;
const ETHTOOL_A_BITSET_BIT_NAME: u16 = 2;
const ETHTOOL_A_BITSET_BIT_VALUE: u16 = 3;

#[derive(Debug, PartialEq, Eq, Clone)]
pub struct EthtoolFeatureBit {
    pub index: u32,
    pub name: String,
    pub value: bool,
}

impl EthtoolFeatureBit {
    fn new(has_mask: bool) -> Self {
        Self {
            index: 0,
            name: "".into(),
            value: !has_mask,
        }
    }
}

fn feature_bits_len(_feature_bits: &[EthtoolFeatureBit]) -> usize {
    todo!("Does not support changing ethtool feature yet")
}

fn feature_bits_emit(_feature_bits: &[EthtoolFeatureBit], _buffer: &mut [u8]) {
    todo!("Does not support changing ethtool feature yet")
}

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum EthtoolFeatureAttr {
    Header(Vec<EthtoolHeader>),
    Hw(Vec<EthtoolFeatureBit>),
    Wanted(Vec<EthtoolFeatureBit>),
    Active(Vec<EthtoolFeatureBit>),
    NoChange(Vec<EthtoolFeatureBit>),
    Other(DefaultNla),
}

impl Nla for EthtoolFeatureAttr {
    fn value_len(&self) -> usize {
        match self {
            Self::Header(hdrs) => hdrs.as_slice().buffer_len(),
            Self::Hw(feature_bits)
            | Self::Wanted(feature_bits)
            | Self::Active(feature_bits)
            | Self::NoChange(feature_bits) => feature_bits_len(feature_bits.as_slice()),
            Self::Other(attr) => attr.value_len(),
        }
    }

    fn kind(&self) -> u16 {
        match self {
            Self::Header(_) => ETHTOOL_A_FEATURES_HEADER | NLA_F_NESTED,
            Self::Hw(_) => ETHTOOL_A_FEATURES_HW | NLA_F_NESTED,
            Self::Wanted(_) => ETHTOOL_A_FEATURES_WANTED | NLA_F_NESTED,
            Self::Active(_) => ETHTOOL_A_FEATURES_ACTIVE | NLA_F_NESTED,
            Self::NoChange(_) => ETHTOOL_A_FEATURES_NOCHANGE | NLA_F_NESTED,
            Self::Other(attr) => attr.kind(),
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        match self {
            Self::Header(ref nlas) => nlas.as_slice().emit(buffer),
            Self::Hw(feature_bits)
            | Self::Wanted(feature_bits)
            | Self::Active(feature_bits)
            | Self::NoChange(feature_bits) => feature_bits_emit(feature_bits.as_slice(), buffer),
            Self::Other(ref attr) => attr.emit(buffer),
        }
    }
}

impl<'a, T: AsRef<[u8]> + ?Sized> Parseable<NlaBuffer<&'a T>> for EthtoolFeatureAttr {
    fn parse(buf: &NlaBuffer<&'a T>) -> Result<Self, DecodeError> {
        let payload = buf.value();
        Ok(match buf.kind() {
            ETHTOOL_A_FEATURES_HEADER => {
                let mut nlas = Vec::new();
                let error_msg = "failed to parse feature header attributes";
                for nla in NlasIterator::new(payload) {
                    let nla = &nla.context(error_msg)?;
                    let parsed = EthtoolHeader::parse(nla).context(error_msg)?;
                    nlas.push(parsed);
                }
                Self::Header(nlas)
            }
            ETHTOOL_A_FEATURES_HW => {
                Self::Hw(parse_bitset_bits_nlas(
                    payload, true, /* ETHTOOL_A_FEATURES_HW is using mask */
                )?)
            }
            ETHTOOL_A_FEATURES_WANTED => {
                Self::Wanted(parse_bitset_bits_nlas(
                    payload, false, /* ETHTOOL_A_FEATURES_WANTED does not use mask */
                )?)
            }
            ETHTOOL_A_FEATURES_ACTIVE => {
                Self::Active(parse_bitset_bits_nlas(
                    payload, false, /* ETHTOOL_A_FEATURES_ACTIVE does not use mask */
                )?)
            }
            ETHTOOL_A_FEATURES_NOCHANGE => {
                Self::NoChange(parse_bitset_bits_nlas(
                    payload, false, /* ETHTOOL_A_FEATURES_NOCHANGE does not use mask */
                )?)
            }
            _ => Self::Other(DefaultNla::parse(buf).context("invalid NLA (unknown kind)")?),
        })
    }
}

fn parse_bitset_bits_nlas(
    raw: &[u8],
    has_mask: bool,
) -> Result<Vec<EthtoolFeatureBit>, DecodeError> {
    let error_msg = "failed to parse feature bit sets";
    for nla in NlasIterator::new(raw) {
        let nla = &nla.context(error_msg)?;
        if nla.kind() == ETHTOOL_A_BITSET_BITS {
            return parse_bitset_bits_nla(nla.value(), has_mask);
        }
    }
    Err("No ETHTOOL_A_BITSET_BITS NLA found".into())
}

fn parse_bitset_bits_nla(
    raw: &[u8],
    has_mask: bool,
) -> Result<Vec<EthtoolFeatureBit>, DecodeError> {
    let mut feature_bits = Vec::new();
    let error_msg = "Failed to parse ETHTOOL_A_BITSET_BITS attributes";
    for bit_nla in NlasIterator::new(raw) {
        let bit_nla = &bit_nla.context(error_msg)?;
        match bit_nla.kind() {
            ETHTOOL_A_BITSET_BITS_BIT => {
                let error_msg = "Failed to parse ETHTOOL_A_BITSET_BITS_BIT attributes";
                let nlas = NlasIterator::new(bit_nla.value());
                let mut cur_bit_info = EthtoolFeatureBit::new(has_mask);
                for nla in nlas {
                    let nla = &nla.context(error_msg)?;
                    let payload = nla.value();
                    match nla.kind() {
                        ETHTOOL_A_BITSET_BIT_INDEX => {
                            if cur_bit_info.index != 0 && !&cur_bit_info.name.is_empty() {
                                feature_bits.push(cur_bit_info);
                                cur_bit_info = EthtoolFeatureBit::new(has_mask);
                            }
                            cur_bit_info.index = parse_u32(payload)
                                .context("Invald ETHTOOL_A_BITSET_BIT_INDEX value")?;
                        }
                        ETHTOOL_A_BITSET_BIT_NAME => {
                            cur_bit_info.name = parse_string(payload)
                                .context("Invald ETHTOOL_A_BITSET_BIT_NAME value")?;
                        }
                        ETHTOOL_A_BITSET_BIT_VALUE => {
                            cur_bit_info.value = true;
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
                if cur_bit_info.index != 0 && !&cur_bit_info.name.is_empty() {
                    feature_bits.push(cur_bit_info);
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
    Ok(feature_bits)
}

pub(crate) fn parse_feature_nlas(buffer: &[u8]) -> Result<Vec<EthtoolAttr>, DecodeError> {
    let mut nlas = Vec::new();
    for nla in NlasIterator::new(buffer) {
        let error_msg = format!(
            "Failed to parse ethtool feature message attribute {:?}",
            nla
        );
        let nla = &nla.context(error_msg.clone())?;
        let parsed = EthtoolFeatureAttr::parse(nla).context(error_msg)?;
        nlas.push(EthtoolAttr::Feature(parsed));
    }
    Ok(nlas)
}
