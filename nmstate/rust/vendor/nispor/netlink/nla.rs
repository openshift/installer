// Copyright 2021 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

use crate::NisporError;
use std::convert::TryInto;
use std::net::Ipv4Addr;
use std::net::Ipv6Addr;

pub(crate) fn parse_as_u8(data: &[u8]) -> Result<u8, NisporError> {
    Ok(*data.get(0).ok_or_else(|| {
        NisporError::bug("wrong index when parsing as u8".into())
    })?)
}

pub(crate) fn parse_as_u16(data: &[u8]) -> Result<u16, NisporError> {
    let err_msg = "wrong index when parsing as u16";
    Ok(u16::from_ne_bytes([
        *data
            .get(0)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(1)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
    ]))
}

pub(crate) fn parse_as_i32(data: &[u8]) -> Result<i32, NisporError> {
    let err_msg = "wrong index when parsing as i32";
    Ok(i32::from_ne_bytes([
        *data
            .get(0)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(1)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(2)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(3)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
    ]))
}

pub(crate) fn parse_as_u32(data: &[u8]) -> Result<u32, NisporError> {
    let err_msg = "wrong index when parsing as u32";
    Ok(u32::from_ne_bytes([
        *data
            .get(0)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(1)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(2)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(3)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
    ]))
}

pub(crate) fn parse_as_u64(data: &[u8]) -> Result<u64, NisporError> {
    let err_msg = "wrong index when parsing as u64";
    Ok(u64::from_ne_bytes([
        *data
            .get(0)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(1)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(2)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(3)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(4)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(5)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(6)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
        *data
            .get(7)
            .ok_or_else(|| NisporError::bug(err_msg.into()))?,
    ]))
}

pub(crate) fn parse_as_ipv4(data: &[u8]) -> Result<Ipv4Addr, NisporError> {
    let addr_bytes: [u8; 4] = data.try_into().map_err(|_| {
        NisporError::invalid_argument(
            "Got invalid IPv4 address u8, the length is not 4".into(),
        )
    })?;
    Ok(Ipv4Addr::from(addr_bytes))
}

pub(crate) fn parse_as_ipv6(data: &[u8]) -> Result<Ipv6Addr, NisporError> {
    let addr_bytes: [u8; 16] = data.try_into().map_err(|_| {
        NisporError::invalid_argument(
            "Got invalid IPv6 address u8, the length is not 16".into(),
        )
    })?;
    Ok(Ipv6Addr::from(addr_bytes))
}
