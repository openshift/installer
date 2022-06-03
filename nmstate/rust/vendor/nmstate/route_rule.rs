use std::collections::{hash_map::Entry, HashMap, HashSet};
use std::convert::TryFrom;

use serde::{Deserialize, Serialize};

use crate::{ip::is_ipv6_addr, ErrorKind, InterfaceIpAddr, NmstateError};

#[derive(Debug, Clone, PartialEq, Eq, Default, Serialize, Deserialize)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct RouteRules {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub config: Option<Vec<RouteRuleEntry>>,
}

impl RouteRules {
    pub fn new() -> Self {
        Self::default()
    }

    // * Neither ip_from nor ip_to should be defined
    pub(crate) fn validate(&self) -> Result<(), NmstateError> {
        if let Some(rules) = self.config.as_ref() {
            for rule in rules.iter().filter(|r| !r.is_absent()) {
                rule.validate()?;
            }
        }
        Ok(())
    }

    // * desired absent route rule is removed unless another matching rule been
    //   added.
    // * desired static rule exists.
    pub fn verify(&self, current: &Self) -> Result<(), NmstateError> {
        if let Some(rules) = self.config.as_ref() {
            let empty_vec: Vec<RouteRuleEntry> = Vec::new();
            let cur_rules = match current.config.as_deref() {
                Some(c) => c,
                None => empty_vec.as_slice(),
            };
            for rule in rules.iter().filter(|r| !r.is_absent()) {
                if !cur_rules.iter().any(|r| rule.is_match(r)) {
                    let e = NmstateError::new(
                        ErrorKind::VerificationError,
                        format!(
                            "Desired route rule {:?} not found after apply",
                            rule
                        ),
                    );
                    log::error!("{}", e);
                    return Err(e);
                }
            }

            for absent_rule in rules.iter().filter(|r| r.is_absent()) {
                // We ignore absent rule if user is replacing old rule
                // with new one.
                if rules
                    .iter()
                    .any(|r| (!r.is_absent()) && absent_rule.is_match(r))
                {
                    continue;
                }

                if let Some(cur_rule) =
                    cur_rules.iter().find(|r| absent_rule.is_match(r))
                {
                    let e = NmstateError::new(
                        ErrorKind::VerificationError,
                        format!(
                            "Desired absent route rule {:?} still found \
                            after apply: {:?}",
                            absent_rule, cur_rule
                        ),
                    );
                    log::error!("{}", e);
                    return Err(e);
                }
            }
        }
        Ok(())
    }

    // RouteRuleEntry been added/removed for specific table id , all(including
    // desire and current) its rules will be included in return hash.
    // Steps:
    //  1. Find out all table id with desired add rules.
    //  2. Find out all table id impacted by desired absent rules.
    //  3. Copy all rules from current which are to changed table id.
    //  4. Remove rules base on absent.
    //  5. Add rules in desire.
    //  6. Sort and remove duplicate rule.
    pub(crate) fn gen_rule_changed_table_ids(
        &self,
        current: &Self,
    ) -> HashMap<u32, Vec<RouteRuleEntry>> {
        let mut ret: HashMap<u32, Vec<RouteRuleEntry>> = HashMap::new();
        let cur_rules_index = current
            .config
            .as_ref()
            .map(|c| create_rule_index_by_table_id(c.as_slice()))
            .unwrap_or_default();
        let des_rules_index = self
            .config
            .as_ref()
            .map(|c| create_rule_index_by_table_id(c.as_slice()))
            .unwrap_or_default();

        let mut table_ids_in_desire: HashSet<u32> =
            des_rules_index.keys().copied().collect();

        // Convert the absent rule without table id to multiple
        // rules with table_id define.
        let absent_rules = flat_absent_rule(
            self.config.as_deref().unwrap_or(&[]),
            current.config.as_deref().unwrap_or(&[]),
        );

        // Include table id which will be impacted by absent rules
        for absent_rule in &absent_rules {
            if let Some(i) = absent_rule.table_id {
                log::debug!(
                    "Route table is impacted by absent rule {:?}",
                    absent_rule
                );
                table_ids_in_desire.insert(i);
            }
        }

        // Copy current rules of desired route table
        for table_id in &table_ids_in_desire {
            if let Some(cur_rules) = cur_rules_index.get(table_id) {
                ret.insert(
                    *table_id,
                    cur_rules
                        .as_slice()
                        .iter()
                        .map(|r| (*r).clone())
                        .collect::<Vec<RouteRuleEntry>>(),
                );
            }
        }

        // Apply absent rules
        for absent_rule in &absent_rules {
            // All absent_rule should have table id here
            if let Some(table_id) = absent_rule.table_id.as_ref() {
                if let Some(rules) = ret.get_mut(table_id) {
                    rules.retain(|r| !absent_rule.is_match(r));
                }
            }
        }

        // Append desire rules
        for (table_id, desire_rules) in des_rules_index.iter() {
            let new_rules = desire_rules
                .iter()
                .map(|r| (*r).clone())
                .collect::<Vec<RouteRuleEntry>>();
            match ret.entry(*table_id) {
                Entry::Occupied(o) => {
                    o.into_mut().extend(new_rules);
                }
                Entry::Vacant(v) => {
                    v.insert(new_rules);
                }
            };
        }

        // Sort and remove the duplicated rules
        for desire_rules in ret.values_mut() {
            desire_rules.sort_unstable();
            desire_rules.dedup();
        }

        ret
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum RouteRuleState {
    Absent,
}

impl Default for RouteRuleState {
    fn default() -> Self {
        Self::Absent
    }
}

#[derive(Debug, Clone, Default, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct RouteRuleEntry {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub state: Option<RouteRuleState>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ip_from: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ip_to: Option<String>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_i64_or_string"
    )]
    pub priority: Option<i64>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        rename = "route-table",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub table_id: Option<u32>,
}

impl RouteRuleEntry {
    pub const USE_DEFAULT_PRIORITY: i64 = -1;
    pub const USE_DEFAULT_ROUTE_TABLE: u32 = 0;
    pub const DEFAULR_ROUTE_TABLE_ID: u32 = 254;

    pub fn new() -> Self {
        Self::default()
    }

    // * Neither ip_from nor ip_to should be defined
    pub(crate) fn validate(&self) -> Result<(), NmstateError> {
        if self.ip_from.is_none() && self.ip_to.is_none() {
            let e = NmstateError::new(
                ErrorKind::InvalidArgument,
                format!(
                    "Neither ip-from or ip-to is defined in route rule {:?}",
                    self
                ),
            );
            log::error!("{}", e);
            return Err(e);
        }
        Ok(())
    }

    fn is_absent(&self) -> bool {
        matches!(self.state, Some(RouteRuleState::Absent))
    }

    fn is_match(&self, other: &Self) -> bool {
        if let Some(ip_from) = self.ip_from.as_deref() {
            let ip_from = if !ip_from.contains('/') {
                match InterfaceIpAddr::try_from(ip_from) {
                    Ok(ref i) => i.into(),
                    Err(e) => {
                        log::error!("{}", e);
                        return false;
                    }
                }
            } else {
                ip_from.to_string()
            };
            if other.ip_from != Some(ip_from) {
                return false;
            }
        }
        if let Some(ip_to) = self.ip_to.as_deref() {
            let ip_to = if !ip_to.contains('/') {
                match InterfaceIpAddr::try_from(ip_to) {
                    Ok(ref i) => i.into(),
                    Err(e) => {
                        log::error!("{}", e);
                        return false;
                    }
                }
            } else {
                ip_to.to_string()
            };
            if other.ip_to != Some(ip_to) {
                return false;
            }
        }
        if self.priority.is_some()
            && self.priority != Some(RouteRuleEntry::USE_DEFAULT_PRIORITY)
            && self.priority != other.priority
        {
            return false;
        }
        if self.table_id.is_some()
            && self.table_id != Some(RouteRuleEntry::USE_DEFAULT_ROUTE_TABLE)
            && self.table_id != other.table_id
        {
            return false;
        }
        true
    }

    // Return tuple of (no_absent, is_ipv4, table_id, ip_from,
    // ip_to, priority)
    fn sort_key(&self) -> (bool, bool, u32, &str, &str, i64) {
        (
            !matches!(self.state, Some(RouteRuleState::Absent)),
            {
                if let Some(ip_from) = self.ip_from.as_ref() {
                    !is_ipv6_addr(ip_from.as_str())
                } else if let Some(ip_to) = self.ip_to.as_ref() {
                    !is_ipv6_addr(ip_to.as_str())
                } else {
                    log::warn!(
                        "Neither ip-from nor ip-to \
                    is defined, treating it a IPv4 route rule"
                    );
                    true
                }
            },
            self.table_id
                .unwrap_or(RouteRuleEntry::USE_DEFAULT_ROUTE_TABLE),
            self.ip_from.as_deref().unwrap_or(""),
            self.ip_to.as_deref().unwrap_or(""),
            self.priority
                .unwrap_or(RouteRuleEntry::USE_DEFAULT_PRIORITY),
        )
    }
}

// For Vec::dedup()
impl PartialEq for RouteRuleEntry {
    fn eq(&self, other: &Self) -> bool {
        self.sort_key() == other.sort_key()
    }
}

// For Vec::sort_unstable()
impl Ord for RouteRuleEntry {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        self.sort_key().cmp(&other.sort_key())
    }
}

// For ord
impl Eq for RouteRuleEntry {}

// For ord
impl PartialOrd for RouteRuleEntry {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}

// Absent rule will be ignored
fn create_rule_index_by_table_id(
    rules: &[RouteRuleEntry],
) -> HashMap<u32, Vec<&RouteRuleEntry>> {
    let mut ret: HashMap<u32, Vec<&RouteRuleEntry>> = HashMap::new();
    for rule in rules {
        if rule.is_absent() {
            continue;
        }
        let table_id = match rule.table_id {
            Some(RouteRuleEntry::USE_DEFAULT_ROUTE_TABLE) | None => {
                RouteRuleEntry::DEFAULR_ROUTE_TABLE_ID
            }
            Some(i) => i,
        };
        match ret.entry(table_id) {
            Entry::Occupied(o) => {
                o.into_mut().push(rule);
            }
            Entry::Vacant(v) => {
                v.insert(vec![rule]);
            }
        };
    }
    ret
}

// All the rules sending to this function has no table id defined.
fn flat_absent_rule(
    desire_rules: &[RouteRuleEntry],
    cur_rules: &[RouteRuleEntry],
) -> Vec<RouteRuleEntry> {
    let mut ret: Vec<RouteRuleEntry> = Vec::new();
    for absent_rule in desire_rules.iter().filter(|r| r.is_absent()) {
        if absent_rule.table_id.is_none() {
            for cur_rule in cur_rules {
                if absent_rule.is_match(cur_rule) {
                    let mut new_absent_rule = absent_rule.clone();
                    new_absent_rule.table_id = cur_rule.table_id;
                    ret.push(new_absent_rule);
                }
            }
        } else {
            ret.push(absent_rule.clone());
        }
    }
    ret
}
