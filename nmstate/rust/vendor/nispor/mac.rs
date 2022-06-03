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

pub(crate) fn parse_as_mac(
    mac_len: usize,
    data: &[u8],
) -> Result<String, NisporError> {
    if data.len() != mac_len {
        return Err(NisporError::bug("wrong size at mac parsing".into()));
    }
    let mut rt = String::new();
    for (i, &val) in data.iter().enumerate() {
        rt.push_str(&format!("{:02x}", val));
        if i != mac_len - 1 {
            rt.push(':');
        }
    }
    Ok(rt)
}

pub(crate) fn mac_str_to_raw(mac_addr: &str) -> Result<Vec<u8>, NisporError> {
    let mac_addr = mac_addr.to_string().replace(':', "").replace('-', "");

    let mut mac_raw: Vec<u8> = Vec::new();

    let mac_addr = mac_addr.replace(':', "");
    let mut chars = mac_addr.chars().peekable();

    while chars.peek().is_some() {
        let chunk: String = chars.by_ref().take(2).collect();
        match u8::from_str_radix(&chunk, 16) {
            Ok(i) => mac_raw.push(i),
            Err(e) => {
                return Err(NisporError::invalid_argument(format!(
                    "Invalid hex string for MAC address {}: {}",
                    mac_addr, e
                )));
            }
        }
    }

    Ok(mac_raw)
}
