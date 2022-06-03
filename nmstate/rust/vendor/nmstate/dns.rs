use std::net::{Ipv4Addr, Ipv6Addr};

use serde::{Deserialize, Serialize};

use crate::{
    ip::is_ipv6_addr, ErrorKind, Interface, Interfaces, NetworkState,
    NmstateError,
};

const DEFAULT_DNS_PRIORITY: i32 = 40;

#[derive(Debug, Clone, PartialEq, Eq, Default, Serialize, Deserialize)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct DnsState {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub running: Option<DnsClientState>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub config: Option<DnsClientState>,
}

impl DnsState {
    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn validate(&self) -> Result<(), NmstateError> {
        if let Some(dns_conf) = self.config.as_ref() {
            if let Some(config_srvs) = dns_conf.server.as_ref() {
                if config_srvs.len() > 2
                    && is_mixed_dns_servers(config_srvs.as_slice())
                {
                    let e = NmstateError::new(
                    ErrorKind::NotImplementedError,
                    "Placing IPv4/IPv6 nameserver in the middle of IPv6/IPv4 \
                    nameservers is not supported yet"
                        .to_string(),
                );
                    log::error!("{}", e);
                    return Err(e);
                }
            }
        }
        Ok(())
    }

    pub(crate) fn verify(&self, current: &Self) -> Result<(), NmstateError> {
        if let Some(conf) = self.config.as_ref() {
            if let Some(srvs) = conf.server.as_ref() {
                let cur_conf = current.config.as_ref().ok_or_else(|| {
                    // Do not log verification error as we have fail-retry
                    NmstateError::new(
                        ErrorKind::VerificationError,
                        format!(
                            "Failed to apply DNS config: desire {:?} got {:?}",
                            self, current
                        ),
                    )
                })?;
                let mut canonicalized_srvs = Vec::new();
                for srv in srvs {
                    if is_ipv6_addr(srv) {
                        if let Ok(ip_addr) = srv.parse::<Ipv6Addr>() {
                            canonicalized_srvs.push(ip_addr.to_string());
                        }
                    } else if let Ok(ip_addr) = srv.parse::<Ipv4Addr>() {
                        canonicalized_srvs.push(ip_addr.to_string());
                    }
                }

                if cur_conf.server != Some(canonicalized_srvs)
                    && !(cur_conf.server.is_none() && srvs.is_empty())
                {
                    return Err(NmstateError::new(
                        ErrorKind::VerificationError,
                        format!(
                            "Failed to apply DNS config: desire name servers \
                            {:?}, got {:?}",
                            srvs,
                            cur_conf.server.as_ref()
                        ),
                    ));
                }
            }
            if let Some(schs) = conf.search.as_ref() {
                let cur_conf = current.config.as_ref().ok_or_else(|| {
                    NmstateError::new(
                        ErrorKind::VerificationError,
                        format!(
                            "Failed to apply DNS config: desire {:?} got {:?}",
                            self, current
                        ),
                    )
                })?;
                if cur_conf.search != Some(schs.to_vec())
                    && !(cur_conf.search.is_none() && schs.is_empty())
                {
                    return Err(NmstateError::new(
                        ErrorKind::VerificationError,
                        format!(
                            "Failed to apply DNS config: desire searches \
                            {:?}, got {:?}",
                            schs,
                            cur_conf.search.as_ref()
                        ),
                    ));
                }
            }
        }
        Ok(())
    }

    pub(crate) fn merge_current(&mut self, current: &Self) {
        if let Some(conf) = self.config.as_mut() {
            if conf.is_purge() {
                self.config = Some(DnsClientState {
                    server: Some(Vec::new()),
                    search: Some(Vec::new()),
                    priority: None,
                });
            } else if let Some(cur_conf) = current.config.as_ref() {
                if conf.server.is_none() {
                    conf.server =
                        Some(cur_conf.server.clone().unwrap_or_default());
                }
                if conf.search.is_none() {
                    conf.search =
                        Some(cur_conf.search.clone().unwrap_or_default());
                }
            }
        } else {
            self.config = current.config.clone();
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default, Serialize, Deserialize)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct DnsClientState {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub server: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub search: Option<Vec<String>>,
    #[serde(skip)]
    // Lower is better
    pub(crate) priority: Option<i32>,
}

impl DnsClientState {
    pub fn new() -> Self {
        Self::default()
    }

    // Whether user want to purge all DNS settings
    pub(crate) fn is_purge(&self) -> bool {
        match (&self.server, &self.search) {
            (Some(srvs), Some(schs)) => srvs.is_empty() && schs.is_empty(),
            (Some(srvs), None) => srvs.is_empty(),
            (None, Some(schs)) => schs.is_empty(),
            (None, None) => true,
        }
    }

    pub(crate) fn is_null(&self) -> bool {
        self.server.as_ref().map(|s| s.len()).unwrap_or_default() == 0
            && self.search.as_ref().map(|s| s.len()).unwrap_or_default() == 0
    }

    pub(crate) fn save_dns_to_iface(
        &self,
        v4_iface_name: &str,
        v6_iface_name: &str,
        add_net_state: &mut NetworkState,
        chg_net_state: &mut NetworkState,
        current: &NetworkState,
    ) -> Result<(), NmstateError> {
        let servers = if let Some(srvs) = self.server.as_ref() {
            srvs.clone()
        } else {
            Vec::new()
        };
        let searches = if let Some(schs) = self.search.as_ref() {
            schs.clone()
        } else {
            Vec::new()
        };
        let mut v4_servers = Vec::new();
        let mut v6_servers = Vec::new();
        let prefer_ipv6_srv = servers
            .get(0)
            .map(|s| is_ipv6_addr(s.as_str()))
            .unwrap_or_default();
        for srv in servers {
            if is_ipv6_addr(&srv) {
                v6_servers.push(srv.to_string())
            } else {
                v4_servers.push(srv.to_string())
            }
        }
        if !v6_servers.is_empty() {
            _save_dns_to_iface(
                true,
                v6_iface_name,
                (v6_servers, searches.clone()),
                add_net_state,
                chg_net_state,
                current,
                prefer_ipv6_srv,
            )?;
        }
        if !v4_servers.is_empty() {
            _save_dns_to_iface(
                false,
                v4_iface_name,
                (v4_servers, searches),
                add_net_state,
                chg_net_state,
                current,
                !prefer_ipv6_srv,
            )?;
        }
        Ok(())
    }
}

pub(crate) fn is_dns_changed(
    desired: &NetworkState,
    current: &NetworkState,
) -> bool {
    match (desired.dns.config.as_ref(), current.dns.config.as_ref()) {
        (None, None) => false,
        (Some(des_config), Some(cur_config)) => {
            if des_config == cur_config {
                // DNS configure not changed, but we need to check whether
                // interface holding DNS still valid to hold DNS config
                !current_dns_ifaces_are_still_valid(desired, current)
            } else {
                true
            }
        }
        (Some(_), None) => true,
        (None, Some(_)) => {
            // DNS configure not changed, but we need to check whether
            // interface holding DNS still valid to hold DNS config
            !current_dns_ifaces_are_still_valid(desired, current)
        }
    }
}

// Return interfaces to hold IPv4 and IPv6 DNS configuration.
pub(crate) fn reselect_dns_ifaces(
    desired: &NetworkState,
    current: &NetworkState,
) -> (String, String) {
    (
        find_ifaces_in_desire(false, &desired.interfaces)
            .or_else(|| {
                find_valid_ifaces_for_dns(
                    false,
                    &desired.interfaces,
                    &current.interfaces,
                )
            })
            .unwrap_or_default(),
        find_ifaces_in_desire(true, &desired.interfaces)
            .or_else(|| {
                find_valid_ifaces_for_dns(
                    true,
                    &desired.interfaces,
                    &current.interfaces,
                )
            })
            .unwrap_or_default(),
    )
}

// Return None if specified interface has IP configuration as None.
fn is_iface_valid_for_dns(is_ipv6: bool, iface: &Interface) -> Option<bool> {
    if is_ipv6 {
        iface.base_iface().ipv6.as_ref().map(|ip_conf| {
            ip_conf.enabled
                && (!ip_conf.is_auto() || (ip_conf.auto_dns == Some(false)))
        })
    } else {
        iface.base_iface().ipv4.as_ref().map(|ip_conf| {
            ip_conf.enabled
                && (!ip_conf.is_auto() || ip_conf.auto_dns == Some(false))
        })
    }
}

// Return false when desired state disabled the IP of dns interface.
fn current_dns_ifaces_are_still_valid(
    desired: &NetworkState,
    current: &NetworkState,
) -> bool {
    for (iface_name, cur_iface) in current.interfaces.kernel_ifaces.iter() {
        if let Some(ipv4) = &cur_iface.base_iface().ipv4 {
            if ipv4.enabled && ipv4.dns.is_some() {
                if let Some(des_iface) =
                    desired.interfaces.kernel_ifaces.get(iface_name)
                {
                    if is_iface_valid_for_dns(false, des_iface) == Some(false) {
                        return false;
                    }
                }
            }
        }
        if let Some(ipv6) = &cur_iface.base_iface().ipv6 {
            if ipv6.enabled && ipv6.dns.is_some() {
                if let Some(des_iface) =
                    desired.interfaces.kernel_ifaces.get(iface_name)
                {
                    if is_iface_valid_for_dns(true, des_iface) == Some(false) {
                        return false;
                    }
                }
            }
        }
    }
    true
}

// Find interface with DHCP disabled and IP enabled from desired interfaces.
fn find_ifaces_in_desire(
    is_ipv6: bool,
    desired: &Interfaces,
) -> Option<String> {
    for (iface_name, iface) in desired.kernel_ifaces.iter() {
        if is_iface_valid_for_dns(is_ipv6, iface) == Some(true) {
            return Some(iface_name.to_string());
        }
    }
    None
}

fn find_valid_ifaces_for_dns(
    is_ipv6: bool,
    desired: &Interfaces,
    current: &Interfaces,
) -> Option<String> {
    for (iface_name, iface) in desired
        .kernel_ifaces
        .iter()
        .chain(current.kernel_ifaces.iter())
    {
        if is_iface_valid_for_dns(is_ipv6, iface) == Some(true) {
            let des_iface = desired.kernel_ifaces.get(iface_name);
            if let Some(des_iface) = des_iface {
                if is_iface_valid_for_dns(is_ipv6, des_iface) != Some(false) {
                    return Some(iface_name.to_string());
                }
            } else {
                return Some(iface_name.to_string());
            }
        }
    }
    None
}

fn is_mixed_dns_servers(srvs: &[String]) -> bool {
    let mut pattern = String::new();
    for srv in srvs {
        let cur_char = if is_ipv6_addr(srv) { '6' } else { '4' };
        if !pattern.ends_with(cur_char) {
            pattern.push(cur_char);
        }
    }
    pattern.contains("464") || pattern.contains("646")
}

// Return a list of interfaces hold DNS configurations
pub(crate) fn get_cur_dns_ifaces(
    current: &Interfaces,
) -> (Vec<String>, Vec<String>) {
    let mut v4_ifaces = Vec::new();
    let mut v6_ifaces = Vec::new();
    for (iface_name, cur_iface) in current.kernel_ifaces.iter() {
        if let Some(ipv4) = &cur_iface.base_iface().ipv4 {
            if ipv4.enabled {
                if let Some(dns_conf) = &ipv4.dns {
                    if !dns_conf.is_null() && !v4_ifaces.contains(iface_name) {
                        v4_ifaces.push(iface_name.to_string())
                    }
                }
            }
        }
        if let Some(ipv6) = &cur_iface.base_iface().ipv6 {
            if ipv6.enabled {
                if let Some(dns_conf) = &ipv6.dns {
                    if !dns_conf.is_null() && !v6_ifaces.contains(iface_name) {
                        v6_ifaces.push(iface_name.to_string())
                    }
                }
            }
        }
    }
    (v4_ifaces, v6_ifaces)
}

fn set_iface_dns_conf(
    is_ipv6: bool,
    iface: &mut Interface,
    servers: Vec<String>,
    searches: Vec<String>,
    priority: Option<i32>,
) {
    let dns_conf = DnsClientState {
        server: Some(servers),
        search: Some(searches),
        priority,
    };
    if is_ipv6 {
        if let Some(ip_conf) = iface.base_iface_mut().ipv6.as_mut() {
            ip_conf.dns = Some(dns_conf);
        } else {
            // Should never happen
            log::error!("BUG: The dns interface is hold None IP {:?}", iface);
        }
    } else if let Some(ip_conf) = iface.base_iface_mut().ipv4.as_mut() {
        ip_conf.dns = Some(dns_conf);
    } else {
        // Should never happen
        log::error!("BUG: The dns interface is hold None IP {:?}", iface);
    }
}

pub(crate) fn purge_dns_config(
    is_ipv6: bool,
    ifaces: &[String],
    desired_net_state: &NetworkState,
    chg_net_state: &mut NetworkState,
    current: &NetworkState,
) {
    for iface_name in ifaces {
        if desired_net_state.interfaces.kernel_ifaces.iter().any(
            |(cur_iface_name, iface)| {
                cur_iface_name == iface_name && iface.is_absent()
            },
        ) {
            continue;
        }

        if let Some(iface) =
            chg_net_state.interfaces.kernel_ifaces.get_mut(iface_name)
        {
            set_iface_dns_conf(is_ipv6, iface, Vec::new(), Vec::new(), None);
        } else {
            // Include interface to chg_net_state
            if let Some(cur_iface) =
                current.interfaces.kernel_ifaces.get(iface_name)
            {
                let mut new_iface = cur_iface.clone_name_type_only();
                new_iface
                    .base_iface_mut()
                    .copy_ip_config_if_none(cur_iface.base_iface());
                set_iface_dns_conf(
                    is_ipv6,
                    &mut new_iface,
                    Vec::new(),
                    Vec::new(),
                    None,
                );
                chg_net_state.append_interface_data(new_iface);
            }
        }
    }
}

// Only preferred: true will save the searches
fn _save_dns_to_iface(
    is_ipv6: bool,
    iface_name: &str,
    dns_conf: (Vec<String>, Vec<String>),
    add_net_state: &mut NetworkState,
    chg_net_state: &mut NetworkState,
    current: &NetworkState,
    preferred: bool,
) -> Result<(), NmstateError> {
    let (servers, searches) = dns_conf;
    if iface_name.is_empty() {
        let e = NmstateError::new(
            ErrorKind::InvalidArgument,
            format!(
                "Failed to find suitable(IP enabled with DHCP off \
                or auto-dns: false) interface for DNS server {:?}",
                servers
            ),
        );
        log::error!("{}", e);
        return Err(e);
    }
    let cur_iface = current.interfaces.kernel_ifaces.get(iface_name);
    if let Some(iface) = add_net_state
        .interfaces
        .kernel_ifaces
        .get_mut(iface_name)
        .or_else(|| chg_net_state.interfaces.kernel_ifaces.get_mut(iface_name))
    {
        if let Some(cur_iface) = cur_iface {
            iface
                .base_iface_mut()
                .copy_ip_config_if_none(cur_iface.base_iface());
        }
        if preferred {
            set_iface_dns_conf(
                is_ipv6,
                iface,
                servers,
                searches,
                Some(DEFAULT_DNS_PRIORITY),
            );
        } else {
            set_iface_dns_conf(
                is_ipv6,
                iface,
                servers,
                Vec::new(),
                Some(DEFAULT_DNS_PRIORITY + 10),
            );
        }
    } else {
        // Copy interface from current
        if let Some(cur_iface) = cur_iface {
            let mut new_iface = cur_iface.clone_name_type_only();
            new_iface
                .base_iface_mut()
                .copy_ip_config_if_none(cur_iface.base_iface());
            // We just append the interface, below unwrap() will never fail
            if preferred {
                set_iface_dns_conf(
                    is_ipv6,
                    &mut new_iface,
                    servers,
                    searches,
                    Some(DEFAULT_DNS_PRIORITY),
                );
            } else {
                set_iface_dns_conf(
                    is_ipv6,
                    &mut new_iface,
                    servers,
                    Vec::new(),
                    Some(DEFAULT_DNS_PRIORITY + 10),
                );
            }
            chg_net_state.append_interface_data(new_iface);
        } else {
            let e = NmstateError::new(
                ErrorKind::Bug,
                format!(
                    "Selected interface {} for dns, but not found it",
                    iface_name
                ),
            );
            log::error!("{}", e);
            return Err(e);
        }
    }

    Ok(())
}
