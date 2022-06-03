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
//

use std::convert::TryFrom;
use std::net::{Ipv4Addr, Ipv6Addr};
use std::str::FromStr;

use super::super::{ErrorKind, NmError};

const IPV6_ADDR_LEN: usize = 16;

pub(crate) fn parse_nm_dns(
    value: zvariant::OwnedValue,
) -> Result<Vec<String>, NmError> {
    let mut dns_srvs = Vec::new();
    for nm_dns_srv in Vec::<zvariant::OwnedValue>::try_from(value)? {
        match nm_dns_srv.value_signature().as_str() {
            "u" => match u32::try_from(nm_dns_srv) {
                Ok(i) => {
                    dns_srvs.push(Ipv4Addr::from(u32::from_be(i)).to_string())
                }
                Err(e) => {
                    let e = NmError::new(
                        ErrorKind::InvalidArgument,
                        format!("Failed to convert to IP address: {}", e),
                    );
                    log::error!("{}", e);
                    return Err(e);
                }
            },
            "ay" => match Vec::<u8>::try_from(nm_dns_srv) {
                Ok(b) => {
                    if b.len() == IPV6_ADDR_LEN {
                        let mut bytes = [0u8; IPV6_ADDR_LEN];
                        bytes.copy_from_slice(&b[..IPV6_ADDR_LEN]);
                        dns_srvs.push(
                            Ipv6Addr::from(u128::from_be_bytes(bytes))
                                .to_string(),
                        );
                    } else {
                        let e = NmError::new(
                            ErrorKind::InvalidArgument,
                            format!("Failed to convert {:?} to IP address", b),
                        );
                        log::error!("{}", e);
                        return Err(e);
                    }
                }
                Err(e) => {
                    let e = NmError::new(
                        ErrorKind::InvalidArgument,
                        format!("Failed to convert to IP address: {}", e),
                    );
                    log::error!("{}", e);
                    return Err(e);
                }
            },
            s => {
                let e = NmError::new(
                    ErrorKind::InvalidArgument,
                    format!(
                        "Failed to convert to IP address: \
                        invalid signature {:?}",
                        s
                    ),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }
    }
    Ok(dns_srvs)
}

pub(crate) fn parse_nm_dns_search(
    value: zvariant::OwnedValue,
) -> Result<Vec<String>, NmError> {
    Vec::<String>::try_from(value).map_err(|e| {
        let e = NmError::new(
            ErrorKind::InvalidArgument,
            format!("In valid DNS search: {}", e),
        );
        log::error!("{}", e);
        e
    })
}

pub(crate) fn nm_ip_dns_to_value(
    dns_srvs: &[String],
) -> Result<zvariant::Value, NmError> {
    let mut is_ipv6 = false;
    let mut dns_values = if let Some(dns_srv) = dns_srvs.get(0) {
        if dns_srv.contains(':') {
            // is IPv6
            is_ipv6 = true;
            zvariant::Array::new(zvariant::Signature::from_str_unchecked("ay"))
        } else {
            zvariant::Array::new(zvariant::Signature::from_str_unchecked("u"))
        }
    } else {
        let e = NmError::new(
            ErrorKind::Bug,
            "nm_ip_dns_to_value got unexpected empty dns_srvs".to_string(),
        );
        log::error!("{}", e);
        return Err(e);
    };
    for dns_srv in dns_srvs {
        if is_ipv6 {
            let ip_addr = Ipv6Addr::from_str(dns_srv).map_err(|e| {
                let e = NmError::new(
                    ErrorKind::InvalidArgument,
                    format!("Invalid IPv6 address: {}: {}", dns_srv, e),
                );
                log::error!("{}", e);
                e
            })?;
            let mut bytes = [0u8; IPV6_ADDR_LEN];
            bytes.copy_from_slice(&ip_addr.octets()[..IPV6_ADDR_LEN]);
            dns_values.append(zvariant::Value::new(bytes.to_vec()))?;
        } else {
            let ip_addr = Ipv4Addr::from_str(dns_srv).map_err(|e| {
                let e = NmError::new(
                    ErrorKind::InvalidArgument,
                    format!("Invalid IPv4 address: {}: {}", dns_srv, e),
                );
                log::error!("{}", e);
                e
            })?;
            let ip_addr_u32 = u32::from_be_bytes(ip_addr.octets()).to_be();
            dns_values.append(zvariant::Value::new(ip_addr_u32))?;
        }
    }
    Ok(zvariant::Value::Array(dns_values))
}

pub(crate) fn nm_ip_dns_search_to_value(
    dns_searches: &[String],
) -> Result<zvariant::Value, NmError> {
    let mut values =
        zvariant::Array::new(zvariant::Signature::from_str_unchecked("s"));
    for search in dns_searches {
        values.append(zvariant::Value::new(search))?;
    }
    Ok(zvariant::Value::Array(values))
}
