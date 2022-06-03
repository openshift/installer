// SPDX-License-Identifier: MIT

//! Generic netlink controller implementation
//!
//! This module provides the definition of the controller packet.
//! It also serves as an example for creating a generic family.

use self::nlas::*;
use crate::{constants::*, traits::*, GenlHeader};
use anyhow::Context;
use netlink_packet_utils::{nla::NlasIterator, traits::*, DecodeError};
use std::convert::{TryFrom, TryInto};

/// Netlink attributes for this family
pub mod nlas;

/// Command code definition of Netlink controller (nlctrl) family
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum GenlCtrlCmd {
    /// Notify from event
    NewFamily,
    /// Notify from event
    DelFamily,
    /// Request to get family info
    GetFamily,
    /// Currently unused
    NewOps,
    /// Currently unused
    DelOps,
    /// Currently unused
    GetOps,
    /// Notify from event
    NewMcastGrp,
    /// Notify from event
    DelMcastGrp,
    /// Currently unused
    GetMcastGrp,
    /// Request to get family policy
    GetPolicy,
}

impl From<GenlCtrlCmd> for u8 {
    fn from(cmd: GenlCtrlCmd) -> u8 {
        use GenlCtrlCmd::*;
        match cmd {
            NewFamily => CTRL_CMD_NEWFAMILY,
            DelFamily => CTRL_CMD_DELFAMILY,
            GetFamily => CTRL_CMD_GETFAMILY,
            NewOps => CTRL_CMD_NEWOPS,
            DelOps => CTRL_CMD_DELOPS,
            GetOps => CTRL_CMD_GETOPS,
            NewMcastGrp => CTRL_CMD_NEWMCAST_GRP,
            DelMcastGrp => CTRL_CMD_DELMCAST_GRP,
            GetMcastGrp => CTRL_CMD_GETMCAST_GRP,
            GetPolicy => CTRL_CMD_GETPOLICY,
        }
    }
}

impl TryFrom<u8> for GenlCtrlCmd {
    type Error = DecodeError;

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        use GenlCtrlCmd::*;
        Ok(match value {
            CTRL_CMD_NEWFAMILY => NewFamily,
            CTRL_CMD_DELFAMILY => DelFamily,
            CTRL_CMD_GETFAMILY => GetFamily,
            CTRL_CMD_NEWOPS => NewOps,
            CTRL_CMD_DELOPS => DelOps,
            CTRL_CMD_GETOPS => GetOps,
            CTRL_CMD_NEWMCAST_GRP => NewMcastGrp,
            CTRL_CMD_DELMCAST_GRP => DelMcastGrp,
            CTRL_CMD_GETMCAST_GRP => GetMcastGrp,
            CTRL_CMD_GETPOLICY => GetPolicy,
            cmd => {
                return Err(DecodeError::from(format!(
                    "Unknown control command: {}",
                    cmd
                )))
            }
        })
    }
}

/// Payload of generic netlink controller
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct GenlCtrl {
    /// Command code of this message
    pub cmd: GenlCtrlCmd,
    /// Netlink attributes in this message
    pub nlas: Vec<GenlCtrlAttrs>,
}

impl GenlFamily for GenlCtrl {
    fn family_name() -> &'static str {
        "nlctrl"
    }

    fn family_id(&self) -> u16 {
        GENL_ID_CTRL
    }

    fn command(&self) -> u8 {
        self.cmd.into()
    }

    fn version(&self) -> u8 {
        2
    }
}

impl Emitable for GenlCtrl {
    fn emit(&self, buffer: &mut [u8]) {
        self.nlas.as_slice().emit(buffer)
    }

    fn buffer_len(&self) -> usize {
        self.nlas.as_slice().buffer_len()
    }
}

impl ParseableParametrized<[u8], GenlHeader> for GenlCtrl {
    fn parse_with_param(buf: &[u8], header: GenlHeader) -> Result<Self, DecodeError> {
        Ok(Self {
            cmd: header.cmd.try_into()?,
            nlas: parse_ctrlnlas(buf)?,
        })
    }
}

fn parse_ctrlnlas(buf: &[u8]) -> Result<Vec<GenlCtrlAttrs>, DecodeError> {
    let nlas = NlasIterator::new(buf)
        .map(|nla| nla.and_then(|nla| GenlCtrlAttrs::parse(&nla)))
        .collect::<Result<Vec<_>, _>>()
        .context("failed to parse control message attributes")?;

    Ok(nlas)
}
