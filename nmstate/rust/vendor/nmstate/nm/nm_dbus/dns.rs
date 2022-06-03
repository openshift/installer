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

use std::collections::HashMap;
use std::convert::TryFrom;

use super::{connection::_from_map, NmError};

#[derive(Debug, Clone, PartialEq, Default)]
pub struct NmDnsEntry {
    pub priority: i32,
    pub domains: Vec<String>,
    pub name_servers: Vec<String>,
    pub interface: String,
    pub is_vpn: bool,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<HashMap<String, zvariant::OwnedValue>> for NmDnsEntry {
    type Error = NmError;
    fn try_from(
        mut v: HashMap<String, zvariant::OwnedValue>,
    ) -> Result<Self, Self::Error> {
        Ok(Self {
            priority: _from_map!(v, "priority", i32::try_from)?
                .unwrap_or_default(),
            domains: _from_map!(v, "domains", Vec::<String>::try_from)?
                .unwrap_or_default(),
            name_servers: _from_map!(
                v,
                "nameservers",
                Vec::<String>::try_from
            )?
            .unwrap_or_default(),
            interface: _from_map!(v, "interface", String::try_from)?
                .unwrap_or_default(),
            is_vpn: _from_map!(v, "vpn", bool::try_from)?.unwrap_or_default(),
            _other: v,
        })
    }
}
