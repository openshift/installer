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

use std::convert::TryFrom;

use log::error;

use super::NmError;

pub(crate) fn own_value_to_bytes_array(
    value: zvariant::OwnedValue,
) -> Result<Vec<u8>, NmError> {
    Ok(zvariant::Array::try_from(value)?
        .iter()
        .filter_map(|val| {
            if let zvariant::Value::U8(i) = val {
                return Some(i);
            }
            None
        })
        .copied()
        .collect())
}

pub(crate) fn u8_array_to_mac_string(data: Vec<u8>) -> String {
    // TODO replace collect().join(":") with intersperse(":").collect() once it
    // is stabilized
    data.iter()
        .map(|byte| format!("{:02X}", byte))
        .collect::<Vec<_>>()
        .join(":")
}

pub(crate) fn mac_str_to_u8_array(mac: &str) -> Vec<u8> {
    let result = mac
        .split(':')
        .map(|byte| u8::from_str_radix(byte, 16))
        .collect::<Result<Vec<_>, _>>();
    match result {
        Ok(arr) => arr,
        Err(e) => {
            error!(
                "Failed to convert to MAC address to bytes {:?}: {}",
                mac, e
            );
            Vec::new()
        }
    }
}
