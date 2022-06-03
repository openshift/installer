// SPDX-License-Identifier: MIT

#[macro_use]
extern crate bitflags;
#[macro_use]
pub(crate) extern crate netlink_packet_utils as utils;
pub(crate) use self::utils::parsers;
pub use self::utils::{traits, DecodeError};

pub use netlink_packet_core::{
    ErrorMessage,
    NetlinkBuffer,
    NetlinkHeader,
    NetlinkMessage,
    NetlinkPayload,
};
pub(crate) use netlink_packet_core::{NetlinkDeserializable, NetlinkSerializable};

pub mod rtnl;
pub use self::rtnl::*;

#[cfg(test)]
#[macro_use]
extern crate lazy_static;

use std::net::IpAddr;

pub(crate) fn ip_len(addr: &IpAddr) -> usize {
    match addr {
        IpAddr::V4(_) => 4,
        IpAddr::V6(_) => 16,
    }
}

pub(crate) fn emit_ip(buf: &mut [u8], addr: &IpAddr) {
    match addr {
        IpAddr::V4(ref ip) => buf.copy_from_slice(&ip.octets()),
        IpAddr::V6(ref ip) => buf.copy_from_slice(&ip.octets()),
    }
}

#[cfg(test)]
#[macro_use]
extern crate pretty_assertions;
