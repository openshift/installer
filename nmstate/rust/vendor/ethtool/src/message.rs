// SPDX-License-Identifier: MIT

use netlink_packet_core::DecodeError;
use netlink_packet_generic::{GenlFamily, GenlHeader};
use netlink_packet_utils::{nla::Nla, Emitable, ParseableParametrized};

use crate::{
    coalesce::{parse_coalesce_nlas, EthtoolCoalesceAttr},
    feature::{parse_feature_nlas, EthtoolFeatureAttr},
    link_mode::{parse_link_mode_nlas, EthtoolLinkModeAttr},
    pause::{parse_pause_nlas, EthtoolPauseAttr},
    ring::{parse_ring_nlas, EthtoolRingAttr},
    EthtoolHeader,
};

const ETHTOOL_MSG_PAUSE_GET: u8 = 21;
const ETHTOOL_MSG_PAUSE_GET_REPLY: u8 = 22;
const ETHTOOL_MSG_FEATURES_GET: u8 = 11;
const ETHTOOL_MSG_FEATURES_GET_REPLY: u8 = 11;
const ETHTOOL_MSG_LINKMODES_GET: u8 = 4;
const ETHTOOL_MSG_LINKMODES_GET_REPLY: u8 = 4;
const ETHTOOL_MSG_RINGS_GET: u8 = 15;
const ETHTOOL_MSG_RINGS_GET_REPLY: u8 = 16;
const ETHTOOL_MSG_COALESCE_GET: u8 = 19;
const ETHTOOL_MSG_COALESCE_GET_REPLY: u8 = 20;

#[derive(Debug, PartialEq, Eq, Clone, Copy)]
pub enum EthtoolCmd {
    PauseGet,
    PauseGetReply,
    FeatureGet,
    FeatureGetReply,
    LinkModeGet,
    LinkModeGetReply,
    RingGet,
    RingGetReply,
    CoalesceGet,
    CoalesceGetReply,
}

impl From<EthtoolCmd> for u8 {
    fn from(cmd: EthtoolCmd) -> Self {
        match cmd {
            EthtoolCmd::PauseGet => ETHTOOL_MSG_PAUSE_GET,
            EthtoolCmd::PauseGetReply => ETHTOOL_MSG_PAUSE_GET_REPLY,
            EthtoolCmd::FeatureGet => ETHTOOL_MSG_FEATURES_GET,
            EthtoolCmd::FeatureGetReply => ETHTOOL_MSG_FEATURES_GET_REPLY,
            EthtoolCmd::LinkModeGet => ETHTOOL_MSG_LINKMODES_GET,
            EthtoolCmd::LinkModeGetReply => ETHTOOL_MSG_LINKMODES_GET_REPLY,
            EthtoolCmd::RingGet => ETHTOOL_MSG_RINGS_GET,
            EthtoolCmd::RingGetReply => ETHTOOL_MSG_RINGS_GET_REPLY,
            EthtoolCmd::CoalesceGet => ETHTOOL_MSG_COALESCE_GET,
            EthtoolCmd::CoalesceGetReply => ETHTOOL_MSG_COALESCE_GET_REPLY,
        }
    }
}

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum EthtoolAttr {
    Pause(EthtoolPauseAttr),
    Feature(EthtoolFeatureAttr),
    LinkMode(EthtoolLinkModeAttr),
    Ring(EthtoolRingAttr),
    Coalesce(EthtoolCoalesceAttr),
}

impl Nla for EthtoolAttr {
    fn value_len(&self) -> usize {
        match self {
            Self::Pause(attr) => attr.value_len(),
            Self::Feature(attr) => attr.value_len(),
            Self::LinkMode(attr) => attr.value_len(),
            Self::Ring(attr) => attr.value_len(),
            Self::Coalesce(attr) => attr.value_len(),
        }
    }

    fn kind(&self) -> u16 {
        match self {
            Self::Pause(attr) => attr.kind(),
            Self::Feature(attr) => attr.kind(),
            Self::LinkMode(attr) => attr.kind(),
            Self::Ring(attr) => attr.kind(),
            Self::Coalesce(attr) => attr.kind(),
        }
    }

    fn emit_value(&self, buffer: &mut [u8]) {
        match self {
            Self::Pause(attr) => attr.emit_value(buffer),
            Self::Feature(attr) => attr.emit_value(buffer),
            Self::LinkMode(attr) => attr.emit_value(buffer),
            Self::Ring(attr) => attr.emit_value(buffer),
            Self::Coalesce(attr) => attr.emit_value(buffer),
        }
    }
}

#[derive(Debug, PartialEq, Eq, Clone)]
pub struct EthtoolMessage {
    pub cmd: EthtoolCmd,
    pub nlas: Vec<EthtoolAttr>,
}

impl GenlFamily for EthtoolMessage {
    fn family_name() -> &'static str {
        "ethtool"
    }

    fn version(&self) -> u8 {
        1
    }

    fn command(&self) -> u8 {
        self.cmd.into()
    }
}

impl EthtoolMessage {
    pub fn new_pause_get(iface_name: Option<&str>) -> Self {
        let nlas = match iface_name {
            Some(s) => vec![EthtoolAttr::Pause(EthtoolPauseAttr::Header(vec![
                EthtoolHeader::DevName(s.to_string()),
            ]))],
            None => vec![EthtoolAttr::Pause(EthtoolPauseAttr::Header(vec![]))],
        };
        EthtoolMessage {
            cmd: EthtoolCmd::PauseGet,
            nlas,
        }
    }

    pub fn new_feature_get(iface_name: Option<&str>) -> Self {
        let nlas = match iface_name {
            Some(s) => vec![EthtoolAttr::Feature(EthtoolFeatureAttr::Header(vec![
                EthtoolHeader::DevName(s.to_string()),
            ]))],
            None => vec![EthtoolAttr::Feature(EthtoolFeatureAttr::Header(vec![]))],
        };
        EthtoolMessage {
            cmd: EthtoolCmd::FeatureGet,
            nlas,
        }
    }

    pub fn new_link_mode_get(iface_name: Option<&str>) -> Self {
        let nlas = match iface_name {
            Some(s) => vec![EthtoolAttr::LinkMode(EthtoolLinkModeAttr::Header(vec![
                EthtoolHeader::DevName(s.to_string()),
            ]))],
            None => vec![EthtoolAttr::LinkMode(EthtoolLinkModeAttr::Header(vec![]))],
        };
        EthtoolMessage {
            cmd: EthtoolCmd::LinkModeGet,
            nlas,
        }
    }

    pub fn new_ring_get(iface_name: Option<&str>) -> Self {
        let nlas = match iface_name {
            Some(s) => vec![EthtoolAttr::Ring(EthtoolRingAttr::Header(vec![
                EthtoolHeader::DevName(s.to_string()),
            ]))],
            None => vec![EthtoolAttr::Ring(EthtoolRingAttr::Header(vec![]))],
        };
        EthtoolMessage {
            cmd: EthtoolCmd::RingGet,
            nlas,
        }
    }

    pub fn new_coalesce_get(iface_name: Option<&str>) -> Self {
        let nlas = match iface_name {
            Some(s) => vec![EthtoolAttr::Coalesce(EthtoolCoalesceAttr::Header(vec![
                EthtoolHeader::DevName(s.to_string()),
            ]))],
            None => vec![EthtoolAttr::Coalesce(EthtoolCoalesceAttr::Header(vec![]))],
        };
        EthtoolMessage {
            cmd: EthtoolCmd::CoalesceGet,
            nlas,
        }
    }
}

impl Emitable for EthtoolMessage {
    fn buffer_len(&self) -> usize {
        self.nlas.as_slice().buffer_len()
    }

    fn emit(&self, buffer: &mut [u8]) {
        self.nlas.as_slice().emit(buffer)
    }
}

impl ParseableParametrized<[u8], GenlHeader> for EthtoolMessage {
    fn parse_with_param(buffer: &[u8], header: GenlHeader) -> Result<Self, DecodeError> {
        Ok(match header.cmd {
            ETHTOOL_MSG_PAUSE_GET_REPLY => Self {
                cmd: EthtoolCmd::PauseGetReply,
                nlas: parse_pause_nlas(buffer)?,
            },
            ETHTOOL_MSG_FEATURES_GET_REPLY => Self {
                cmd: EthtoolCmd::FeatureGetReply,
                nlas: parse_feature_nlas(buffer)?,
            },
            ETHTOOL_MSG_LINKMODES_GET_REPLY => Self {
                cmd: EthtoolCmd::LinkModeGetReply,
                nlas: parse_link_mode_nlas(buffer)?,
            },
            ETHTOOL_MSG_RINGS_GET_REPLY => Self {
                cmd: EthtoolCmd::RingGetReply,
                nlas: parse_ring_nlas(buffer)?,
            },
            ETHTOOL_MSG_COALESCE_GET_REPLY => Self {
                cmd: EthtoolCmd::CoalesceGetReply,
                nlas: parse_coalesce_nlas(buffer)?,
            },
            cmd => {
                return Err(DecodeError::from(format!(
                    "Unsupported ethtool reply command: {}",
                    cmd
                )))
            }
        })
    }
}
