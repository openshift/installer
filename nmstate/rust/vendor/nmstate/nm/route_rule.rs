use std::convert::TryFrom;

use crate::nm::nm_dbus::NmIpRouteRule;

use crate::{ip::is_ipv6_addr, InterfaceIpAddr, NmstateError, RouteRuleEntry};

// NM require route rule priority been set explicitly, use 30,000 when
// desire state instruct to use USE_DEFAULT_PRIORITY
const ROUTE_RULE_DEFAULT_PRIORIRY: u32 = 30000;

pub(crate) fn gen_nm_ip_rules(
    rules: &[RouteRuleEntry],
    is_ipv6: bool,
) -> Result<Vec<NmIpRouteRule>, NmstateError> {
    let mut ret = Vec::new();
    for rule in rules {
        let mut nm_rule = NmIpRouteRule::default();
        nm_rule.family = Some(if is_ipv6 {
            libc::AF_INET6
        } else {
            libc::AF_INET
        });
        if let Some(addr) = rule.ip_from.as_deref() {
            match (is_ipv6, is_ipv6_addr(addr)) {
                (true, true) | (false, false) => {
                    let ip_addr = InterfaceIpAddr::try_from(addr)?;
                    nm_rule.from_len = Some(ip_addr.prefix_length);
                    nm_rule.from = Some(ip_addr.ip.to_string());
                }
                _ => continue,
            }
        }
        if let Some(addr) = rule.ip_to.as_deref() {
            match (is_ipv6, is_ipv6_addr(addr)) {
                (true, true) | (false, false) => {
                    let ip_addr = InterfaceIpAddr::try_from(addr)?;
                    nm_rule.to_len = Some(ip_addr.prefix_length);
                    nm_rule.to = Some(ip_addr.ip.to_string());
                }
                _ => continue,
            }
        }
        nm_rule.priority = match rule.priority {
            Some(RouteRuleEntry::USE_DEFAULT_PRIORITY) | None => {
                Some(ROUTE_RULE_DEFAULT_PRIORIRY)
            }
            Some(i) => Some(i as u32),
        };
        nm_rule.table = match rule.table_id {
            Some(RouteRuleEntry::USE_DEFAULT_ROUTE_TABLE) | None => {
                Some(RouteRuleEntry::DEFAULR_ROUTE_TABLE_ID)
            }
            Some(i) => Some(i),
        };

        ret.push(nm_rule);
    }
    Ok(ret)
}
