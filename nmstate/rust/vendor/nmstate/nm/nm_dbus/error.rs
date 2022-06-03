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

#[derive(Debug, Clone, PartialEq, Eq)]
#[non_exhaustive]
pub enum ErrorKind {
    DbusConnectionError,
    CheckpointConflict,
    InvalidArgument,
    NotFound,
    IncompatibleReapply,
    Bug,
    Timeout,
}

impl std::fmt::Display for ErrorKind {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{:?}", self)
    }
}

#[derive(Debug)]
pub struct NmError {
    pub kind: ErrorKind,
    pub msg: String,
    pub dbus_error: Option<zbus::Error>,
}

impl NmError {
    pub fn new(kind: ErrorKind, msg: String) -> Self {
        Self {
            kind,
            msg,
            dbus_error: None,
        }
    }
}

impl std::fmt::Display for NmError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.msg)
    }
}

impl From<zbus::Error> for NmError {
    fn from(e: zbus::Error) -> Self {
        Self {
            kind: ErrorKind::DbusConnectionError,
            msg: format!("{}", e),
            dbus_error: Some(e),
        }
    }
}

impl From<zvariant::Error> for NmError {
    fn from(e: zvariant::Error) -> Self {
        Self {
            kind: ErrorKind::Bug,
            msg: format!("{}", e),
            dbus_error: None,
        }
    }
}

impl From<std::io::Error> for NmError {
    fn from(e: std::io::Error) -> Self {
        Self::new(ErrorKind::Bug, format!("failed to write file: {}", e))
    }
}
