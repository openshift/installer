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

use ethtool::EthtoolError;
use libc::{EEXIST, EPERM};
use netlink_packet_utils::DecodeError;
use serde::Serialize;

#[derive(Debug, Clone, Serialize)]
#[serde(rename_all = "snake_case")]
pub enum ErrorKind {
    InvalidArgument,
    NetlinkError,
    NisporBug,
    PermissionDeny,
}

impl std::fmt::Display for ErrorKind {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{:?}", self)
    }
}

#[derive(Debug, Clone, Serialize)]
#[non_exhaustive]
pub struct NisporError {
    pub kind: ErrorKind,
    pub msg: String,
}

impl NisporError {
    pub(crate) fn bug(message: String) -> NisporError {
        NisporError {
            kind: ErrorKind::NisporBug,
            msg: message,
        }
    }
    pub(crate) fn permission_deny(message: String) -> NisporError {
        NisporError {
            kind: ErrorKind::PermissionDeny,
            msg: message,
        }
    }
    pub(crate) fn invalid_argument(message: String) -> NisporError {
        NisporError {
            kind: ErrorKind::InvalidArgument,
            msg: message,
        }
    }
}

impl std::fmt::Display for NisporError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.msg)
    }
}

impl std::error::Error for NisporError {
    /* TODO */
}

impl std::convert::From<rtnetlink::Error> for NisporError {
    fn from(e: rtnetlink::Error) -> Self {
        match e {
            rtnetlink::Error::NetlinkError(netlink_err) => {
                match netlink_err.code.abs() {
                    EEXIST => NisporError::bug(format!(
                        "Got netlink EEXIST error: {}",
                        netlink_err
                    )),
                    EPERM => NisporError::permission_deny(format!(
                        "{}",
                        netlink_err,
                    )),
                    _ => NisporError::bug(format!(
                        "Got netlink unknown error: code {}, msg: {}",
                        netlink_err.code, netlink_err,
                    )),
                }
            }
            _ => NisporError {
                kind: ErrorKind::NetlinkError,
                msg: e.to_string(),
            },
        }
    }
}

impl std::convert::From<EthtoolError> for NisporError {
    fn from(e: EthtoolError) -> Self {
        NisporError {
            kind: ErrorKind::NetlinkError,
            msg: e.to_string(),
        }
    }
}

impl std::convert::From<std::ffi::FromBytesWithNulError> for NisporError {
    fn from(e: std::ffi::FromBytesWithNulError) -> Self {
        NisporError {
            kind: ErrorKind::NisporBug,
            msg: format!("FromBytesWithNulError: {}", e),
        }
    }
}

impl std::convert::From<std::str::Utf8Error> for NisporError {
    fn from(e: std::str::Utf8Error) -> Self {
        NisporError {
            kind: ErrorKind::NisporBug,
            msg: format!("Utf8Error: {}", e),
        }
    }
}

impl std::convert::From<DecodeError> for NisporError {
    fn from(e: DecodeError) -> Self {
        NisporError {
            kind: ErrorKind::NetlinkError,
            msg: e.to_string(),
        }
    }
}

impl std::convert::From<std::io::Error> for NisporError {
    fn from(e: std::io::Error) -> Self {
        NisporError {
            kind: ErrorKind::NisporBug,
            msg: e.to_string(),
        }
    }
}

impl std::convert::From<std::net::AddrParseError> for NisporError {
    fn from(e: std::net::AddrParseError) -> Self {
        NisporError {
            kind: ErrorKind::InvalidArgument,
            msg: e.to_string(),
        }
    }
}
