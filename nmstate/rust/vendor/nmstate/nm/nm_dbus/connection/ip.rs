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

use std::collections::HashMap;
use std::convert::TryFrom;

use log::warn;

use serde::Deserialize;

use super::super::{
    connection::dns::{
        nm_ip_dns_search_to_value, nm_ip_dns_to_value, parse_nm_dns,
        parse_nm_dns_search,
    },
    connection::route::{
        nm_ip_routes_to_value, parse_nm_ip_route_data, NmIpRoute,
    },
    connection::route_rule::{
        nm_ip_rules_to_value, parse_nm_ip_rule_data, NmIpRouteRule,
    },
    connection::DbusDictionary,
    error::{ErrorKind, NmError},
};

#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(try_from = "zvariant::OwnedValue")]
#[non_exhaustive]
pub enum NmSettingIpMethod {
    Auto,
    Disabled,
    LinkLocal,
    Manual,
    Shared,
    Dhcp,   // IPv6 only,
    Ignore, // Ipv6 only,
}

impl Default for NmSettingIpMethod {
    fn default() -> Self {
        Self::Auto
    }
}

impl std::fmt::Display for NmSettingIpMethod {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Auto => "auto",
                Self::Disabled => "disabled",
                Self::LinkLocal => "link-local",
                Self::Manual => "manual",
                Self::Shared => "shared",
                Self::Dhcp => "dhcp",
                Self::Ignore => "ignore",
            }
        )
    }
}

impl TryFrom<zvariant::OwnedValue> for NmSettingIpMethod {
    type Error = NmError;
    fn try_from(value: zvariant::OwnedValue) -> Result<Self, Self::Error> {
        let str_value = String::try_from(value)?;
        match str_value.as_str() {
            "auto" => Ok(Self::Auto),
            "disabled" => Ok(Self::Disabled),
            "link-local" => Ok(Self::LinkLocal),
            "manual" => Ok(Self::Manual),
            "shared" => Ok(Self::Shared),
            "dhcp" => Ok(Self::Dhcp),
            "ignore" => Ok(Self::Ignore),
            _ => Err(NmError::new(
                ErrorKind::InvalidArgument,
                format!("Invalid IP method {}", str_value),
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingIp {
    pub method: Option<NmSettingIpMethod>,
    pub addresses: Vec<String>,
    pub routes: Vec<NmIpRoute>,
    pub route_rules: Vec<NmIpRouteRule>,
    pub dns_priority: Option<i32>,
    pub dns_search: Option<Vec<String>>,
    pub dns: Option<Vec<String>>,
    pub ignore_auto_dns: Option<bool>,
    pub never_default: Option<bool>,
    pub ignore_auto_routes: Option<bool>,
    pub route_table: Option<u32>,
    pub dhcp_client_id: Option<String>,
    pub dhcp_timeout: Option<i32>,
    pub gateway: Option<String>,
    // IPv6 only
    pub ra_timeout: Option<i32>,
    // IPv6 only
    pub addr_gen_mode: Option<i32>,
    // IPv6 only
    pub dhcp_duid: Option<String>,
    // IPv6 only
    pub dhcp_iaid: Option<String>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSettingIp {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        let mut setting = Self {
            method: _from_map!(v, "method", NmSettingIpMethod::try_from)?,
            addresses: _from_map!(v, "address-data", parse_nm_ip_address_data)?
                .unwrap_or_default(),
            routes: _from_map!(v, "route-data", parse_nm_ip_route_data)?
                .unwrap_or_default(),
            route_rules: _from_map!(v, "routing-rules", parse_nm_ip_rule_data)?
                .unwrap_or_default(),
            dns: _from_map!(v, "dns", parse_nm_dns)?,
            dns_search: _from_map!(v, "dns-search", parse_nm_dns_search)?,
            dns_priority: _from_map!(v, "dns-priority", i32::try_from)?,
            ignore_auto_dns: _from_map!(v, "ignore-auto-dns", bool::try_from)?,
            never_default: _from_map!(v, "never-default", bool::try_from)?,
            ignore_auto_routes: _from_map!(
                v,
                "ignore-auto-routes",
                bool::try_from
            )?,
            dhcp_client_id: _from_map!(v, "dhcp-client-id", String::try_from)?,
            dhcp_timeout: _from_map!(v, "dhcp-timeout", i32::try_from)?,
            ra_timeout: _from_map!(v, "ra-timeout", i32::try_from)?,
            addr_gen_mode: _from_map!(v, "addr-gen-mode", i32::try_from)?,
            dhcp_duid: _from_map!(v, "dhcp-duid", String::try_from)?,
            dhcp_iaid: _from_map!(v, "dhcp-iaid", String::try_from)?,
            route_table: _from_map!(v, "route-table", u32::try_from)?,
            gateway: _from_map!(v, "gateway", String::try_from)?,
            ..Default::default()
        };

        // NM deprecated `addresses` property in the favor of `addresss-data`
        v.remove("addresses");
        // NM deprecated `routes` property in the favor of `routes-data`
        v.remove("routes");
        setting._other = v;
        Ok(setting)
    }
}

impl NmSettingIp {
    pub(crate) fn to_keyfile(
        &self,
    ) -> Result<HashMap<String, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        for (k, v) in self.to_value()?.drain() {
            if !vec!["address-data", "route-data", "dns"].contains(&k) {
                ret.insert(k.to_string(), v);
            }
        }
        for (i, addr) in self.addresses.as_slice().iter().enumerate() {
            ret.insert(format!("address{}", i), zvariant::Value::new(addr));
        }

        for (i, route) in self.routes.as_slice().iter().enumerate() {
            for (k, v) in route.to_keyfile().drain() {
                ret.insert(
                    if k.is_empty() {
                        format!("route{}", i)
                    } else {
                        format!("route{}_{}", i, k)
                    },
                    zvariant::Value::new(v),
                );
            }
        }
        if let Some(dns) = self.dns.as_ref() {
            ret.insert("dns".to_string(), zvariant::Value::new(dns));
        }

        Ok(ret)
    }

    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        if let Some(v) = &self.method {
            ret.insert("method", zvariant::Value::new(format!("{}", v)));
        }
        let mut addresss_data = zvariant::Array::new(
            zvariant::Signature::from_str_unchecked("a{sv}"),
        );
        for addr_str in &self.addresses {
            let addr_str_split: Vec<&str> = addr_str.split('/').collect();
            if addr_str_split.len() != 2 {
                return Err(NmError::new(
                    ErrorKind::InvalidArgument,
                    format!("Invalid IP address {}", addr_str),
                ));
            }
            let prefix = addr_str_split[1].parse::<u32>().map_err(|e| {
                NmError::new(
                    ErrorKind::InvalidArgument,
                    format!(
                        "Invalid IP address prefix {}: {}",
                        addr_str_split[1], e
                    ),
                )
            })?;
            let mut addr_dict = zvariant::Dict::new(
                zvariant::Signature::from_str_unchecked("s"),
                zvariant::Signature::from_str_unchecked("v"),
            );
            addr_dict.append(
                zvariant::Value::new("address".to_string()),
                zvariant::Value::Value(Box::new(zvariant::Value::new(
                    addr_str_split[0].to_string(),
                ))),
            )?;
            addr_dict.append(
                zvariant::Value::new("prefix".to_string()),
                zvariant::Value::Value(Box::new(zvariant::Value::U32(prefix))),
            )?;
            addresss_data.append(zvariant::Value::Dict(addr_dict))?;
        }
        ret.insert("address-data", zvariant::Value::Array(addresss_data));
        ret.insert("route-data", nm_ip_routes_to_value(&self.routes)?);
        ret.insert("routing-rules", nm_ip_rules_to_value(&self.route_rules)?);
        if let Some(dns_servers) = self.dns.as_ref() {
            if !dns_servers.is_empty() {
                ret.insert("dns", nm_ip_dns_to_value(dns_servers)?);
            }
        }
        if let Some(dns_searches) = self.dns_search.as_ref() {
            ret.insert("dns-search", nm_ip_dns_search_to_value(dns_searches)?);
        }
        if let Some(dns_priority) = self.dns_priority {
            ret.insert("dns-priority", zvariant::Value::new(dns_priority));
        }
        if let Some(v) = self.ignore_auto_dns {
            ret.insert("ignore-auto-dns", zvariant::Value::new(v));
        }
        if let Some(v) = self.ignore_auto_routes {
            ret.insert("ignore-auto-routes", zvariant::Value::new(v));
        }
        if let Some(v) = self.never_default {
            ret.insert("never-default", zvariant::Value::new(v));
        }
        if let Some(v) = &self.dhcp_client_id {
            ret.insert("dhcp-client-id", zvariant::Value::new(v));
        }
        if let Some(v) = self.dhcp_timeout {
            ret.insert("dhcp-timeout", zvariant::Value::new(v));
        }
        if let Some(v) = self.ra_timeout {
            ret.insert("ra-timeout", zvariant::Value::new(v));
        }
        if let Some(v) = self.addr_gen_mode {
            ret.insert("addr-gen-mode", zvariant::Value::new(v));
        }
        if let Some(v) = &self.dhcp_duid {
            ret.insert("dhcp-duid", zvariant::Value::new(v));
        }
        if let Some(v) = &self.dhcp_iaid {
            ret.insert("dhcp-iaid", zvariant::Value::new(v));
        }
        if let Some(v) = &self.route_table {
            ret.insert("route-table", zvariant::Value::new(v));
        }
        if let Some(v) = &self.gateway {
            ret.insert("gateway", zvariant::Value::new(v));
        }
        ret.extend(self._other.iter().map(|(key, value)| {
            (key.as_str(), zvariant::Value::from(value.clone()))
        }));
        Ok(ret)
    }
}

fn parse_nm_ip_address_data(
    value: zvariant::OwnedValue,
) -> Result<Vec<String>, NmError> {
    let mut addresses = Vec::new();
    for nm_addr in <Vec<zvariant::OwnedValue>>::try_from(value)? {
        let nm_addr_display = format!("{:?}", nm_addr);
        let mut nm_addr =
            match <HashMap<String, zvariant::OwnedValue>>::try_from(nm_addr) {
                Ok(a) => a,
                Err(e) => {
                    warn!(
                        "Failed to convert {} to HashMap: {}",
                        nm_addr_display, e
                    );
                    continue;
                }
            };
        let address = if let Some(a) = nm_addr
            .remove("address")
            .and_then(|a| String::try_from(a).ok())
        {
            a
        } else {
            warn!("Failed to find address property from {:?}", nm_addr);

            continue;
        };
        let prefix = if let Some(a) =
            nm_addr.remove("prefix").and_then(|a| u32::try_from(a).ok())
        {
            a
        } else {
            warn!("Failed to find address property from {:?}", nm_addr);

            continue;
        };
        addresses.push(format!("{}/{}", address, prefix));
    }
    Ok(addresses)
}
