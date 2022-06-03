use std::collections::{hash_map::Entry, HashMap, HashSet};

use log::{debug, error};
use serde::{Deserialize, Serialize};

use crate::{ip::is_ipv6_addr, ErrorKind, NmstateError};

#[derive(Debug, Clone, PartialEq, Eq, Default, Serialize, Deserialize)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct Routes {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub running: Option<Vec<RouteEntry>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub config: Option<Vec<RouteEntry>>,
}

impl Routes {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn validate(&self) -> Result<(), NmstateError> {
        // All desire non-absent route should have next hop interface
        if let Some(config_routes) = self.config.as_ref() {
            for route in config_routes.iter().filter(|r| !r.is_absent()) {
                if route.next_hop_iface.is_none() {
                    let e = NmstateError::new(
                        ErrorKind::NotImplementedError,
                        format!(
                            "Route with empty next hop interface \
                        is not supported: {:?}",
                            route
                        ),
                    );
                    error!("{}", e);
                    return Err(e);
                }
            }
        }
        Ok(())
    }

    // Kernel might append additional routes. For example, IPv6 default
    // gateway will generate /128 static direct route.
    // Hence, we only check:
    // * desired absent route is removed unless another matching route been
    //   added.
    // * desired static route exists.
    pub(crate) fn verify(
        &self,
        current: &Self,
        ignored_ifaces: &[String],
    ) -> Result<(), NmstateError> {
        if let Some(config_routes) = self.config.as_ref() {
            let cur_config_routes = match current.config.as_ref() {
                Some(c) => c.to_vec(),
                None => Vec::new(),
            };
            for desire_route in config_routes.iter().filter(|r| !r.is_absent())
            {
                if !cur_config_routes.iter().any(|r| desire_route.is_match(r)) {
                    let e = NmstateError::new(
                        ErrorKind::VerificationError,
                        format!(
                            "Desired route {:?} not found after apply",
                            desire_route
                        ),
                    );
                    error!("{}", e);
                    return Err(e);
                }
            }

            for absent_route in config_routes.iter().filter(|r| r.is_absent()) {
                // We ignore absent route if user is replacing old route
                // with new one.
                if config_routes
                    .iter()
                    .any(|r| (!r.is_absent()) && absent_route.is_match(r))
                {
                    continue;
                }

                if let Some(cur_route) =
                    cur_config_routes.iter().find(|r|
                        if let Some(iface) = r.next_hop_iface.as_ref() {
                            !ignored_ifaces.contains(
                                iface
                            )
                        } else {
                            true
                        } && absent_route.is_match(r))
                {
                    let e = NmstateError::new(
                        ErrorKind::VerificationError,
                        format!(
                            "Desired absent route {:?} still found \
                            after apply: {:?}",
                            absent_route, cur_route
                        ),
                    );
                    error!("{}", e);
                    return Err(e);
                }
            }
        }
        Ok(())
    }

    // RouteEntry been added/removed from specific interface, all(including
    // desire and current) its routes will be included in return hash.
    // Steps:
    //  1. Find out all interface with desired add routes.
    //  2. Find out all interface impacted by desired absent routes.
    //  3. Copy all routes from current which are to changed interface.
    //  4. Remove routes base on absent.
    //  5. Add routes in desire.
    //  6. Sort and remove duplicate route.
    pub(crate) fn gen_changed_ifaces_and_routes(
        &self,
        current: &Self,
    ) -> HashMap<String, Vec<RouteEntry>> {
        let mut ret: HashMap<String, Vec<RouteEntry>> = HashMap::new();
        let cur_routes_index = current
            .config
            .as_ref()
            .map(|c| create_route_index_by_iface(c.as_slice()))
            .unwrap_or_default();
        let des_routes_index = self
            .config
            .as_ref()
            .map(|c| create_route_index_by_iface(c.as_slice()))
            .unwrap_or_default();

        let mut iface_names_in_desire: HashSet<&str> =
            des_routes_index.keys().copied().collect();

        // Convert the absent route without iface to multiple routes with
        // iface define.
        let absent_routes = flat_absent_route(
            self.config.as_deref().unwrap_or(&[]),
            current.config.as_deref().unwrap_or(&[]),
        );

        // Include interface which will be impacted by absent routes
        for absent_route in &absent_routes {
            if let Some(i) = absent_route.next_hop_iface.as_ref() {
                debug!(
                    "Interface is impacted by absent route {:?}",
                    absent_route
                );
                iface_names_in_desire.insert(i);
            }
        }

        // Copy current routes next hop to changed interfaces
        for iface_name in &iface_names_in_desire {
            if let Some(cur_routes) = cur_routes_index.get(iface_name) {
                ret.insert(
                    iface_name.to_string(),
                    cur_routes
                        .as_slice()
                        .iter()
                        .map(|r| (*r).clone())
                        .collect::<Vec<RouteEntry>>(),
                );
            }
        }

        // Apply absent routes
        for absent_route in &absent_routes {
            // All absent_route should have interface name here
            if let Some(iface_name) = absent_route.next_hop_iface.as_ref() {
                if let Some(routes) = ret.get_mut(iface_name) {
                    routes.retain(|r| !absent_route.is_match(r));
                }
            }
        }

        // Append desire routes
        for (iface_name, desire_routes) in des_routes_index.iter() {
            let new_routes = desire_routes
                .iter()
                .map(|r| (*r).clone())
                .collect::<Vec<RouteEntry>>();
            match ret.entry(iface_name.to_string()) {
                Entry::Occupied(o) => {
                    o.into_mut().extend(new_routes);
                }
                Entry::Vacant(v) => {
                    v.insert(new_routes);
                }
            };
        }

        // Sort and remove the duplicated routes
        for desire_routes in ret.values_mut() {
            desire_routes.sort_unstable();
            desire_routes.dedup();
        }

        ret
    }

    pub(crate) fn remove_ignored_iface_routes(
        &mut self,
        iface_names: &[String],
    ) {
        if let Some(config_routes) = self.config.as_mut() {
            config_routes.retain(|r| {
                if let Some(i) = r.next_hop_iface.as_ref() {
                    !iface_names.contains(i)
                } else {
                    true
                }
            })
        }
    }

    pub(crate) fn get_config_routes_of_iface(
        &self,
        iface_name: &str,
    ) -> Option<Vec<RouteEntry>> {
        self.config.as_ref().map(|rts| {
            rts.iter()
                .filter(|r| r.next_hop_iface.as_deref() == Some(iface_name))
                .cloned()
                .collect()
        })
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum RouteState {
    Absent,
}

impl Default for RouteState {
    fn default() -> Self {
        Self::Absent
    }
}

#[derive(Debug, Clone, Default, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct RouteEntry {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub state: Option<RouteState>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub destination: Option<String>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        rename = "next-hop-interface"
    )]
    pub next_hop_iface: Option<String>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        rename = "next-hop-address"
    )]
    pub next_hop_addr: Option<String>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_i64_or_string"
    )]
    pub metric: Option<i64>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub table_id: Option<u32>,
}

impl RouteEntry {
    pub const USE_DEFAULT_METRIC: i64 = -1;
    pub const USE_DEFAULT_ROUTE_TABLE: u32 = 0;

    pub fn new() -> Self {
        Self::default()
    }

    fn is_absent(&self) -> bool {
        matches!(self.state, Some(RouteState::Absent))
    }

    fn is_match(&self, other: &Self) -> bool {
        if self.destination.as_ref().is_some()
            && self.destination != other.destination
        {
            return false;
        }
        if self.next_hop_iface.as_ref().is_some()
            && self.next_hop_iface != other.next_hop_iface
        {
            return false;
        }

        if self.next_hop_addr.as_ref().is_some()
            && self.next_hop_addr != other.next_hop_addr
        {
            return false;
        }
        if self.metric.is_some()
            && self.metric != Some(RouteEntry::USE_DEFAULT_METRIC)
            && self.metric != other.metric
        {
            return false;
        }
        if self.table_id.is_some()
            && self.table_id != Some(RouteEntry::USE_DEFAULT_ROUTE_TABLE)
            && self.table_id != other.table_id
        {
            return false;
        }
        true
    }

    // Return tuple of (no_absent, is_ipv4, table_id, next_hop_iface,
    // destination, next_hop_addr, metric)
    // The metric difference is ignored
    fn sort_key(&self) -> (bool, bool, u32, &str, &str, &str, i64) {
        (
            !matches!(self.state, Some(RouteState::Absent)),
            !self
                .destination
                .as_ref()
                .map(|d| is_ipv6_addr(d.as_str()))
                .unwrap_or_default(),
            self.table_id.unwrap_or(RouteEntry::USE_DEFAULT_ROUTE_TABLE),
            self.next_hop_iface.as_deref().unwrap_or(""),
            self.destination.as_deref().unwrap_or(""),
            self.next_hop_addr.as_deref().unwrap_or(""),
            self.metric.unwrap_or(RouteEntry::USE_DEFAULT_METRIC),
        )
    }
}

// For Vec::dedup()
impl PartialEq for RouteEntry {
    fn eq(&self, other: &Self) -> bool {
        self.sort_key() == other.sort_key()
    }
}

// For Vec::sort_unstable()
impl Ord for RouteEntry {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        self.sort_key().cmp(&other.sort_key())
    }
}

// For ord
impl Eq for RouteEntry {}

// For ord
impl PartialOrd for RouteEntry {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}

// Absent route will be ignored
fn create_route_index_by_iface(
    routes: &[RouteEntry],
) -> HashMap<&str, Vec<&RouteEntry>> {
    let mut ret: HashMap<&str, Vec<&RouteEntry>> = HashMap::new();
    for route in routes {
        if route.is_absent() {
            continue;
        }
        let next_hop_iface = route.next_hop_iface.as_deref().unwrap_or("");
        match ret.entry(next_hop_iface) {
            Entry::Occupied(o) => {
                o.into_mut().push(route);
            }
            Entry::Vacant(v) => {
                v.insert(vec![route]);
            }
        };
    }
    ret
}

// All the routes sending to this function has no interface defined.
fn flat_absent_route(
    desire_routes: &[RouteEntry],
    cur_routes: &[RouteEntry],
) -> Vec<RouteEntry> {
    let mut ret: Vec<RouteEntry> = Vec::new();
    for absent_route in desire_routes.iter().filter(|r| r.is_absent()) {
        if absent_route.next_hop_iface.is_none() {
            for cur_route in cur_routes {
                if absent_route.is_match(cur_route) {
                    let mut new_absent_route = absent_route.clone();
                    new_absent_route.next_hop_iface =
                        cur_route.next_hop_iface.as_ref().cloned();
                    ret.push(new_absent_route);
                }
            }
        } else {
            ret.push(absent_route.clone());
        }
    }
    ret
}
