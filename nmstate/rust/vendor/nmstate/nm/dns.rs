use crate::nm::nm_dbus::{NmApi, NmSettingIp};

use crate::{
    nm::error::nm_error_to_nmstate, DnsClientState, DnsState, Interfaces,
    NmstateError,
};

pub(crate) fn nm_dns_to_nmstate(nm_ip_setting: &NmSettingIp) -> DnsClientState {
    DnsClientState {
        server: nm_ip_setting.dns.clone(),
        search: nm_ip_setting.dns_search.clone(),
        priority: nm_ip_setting.dns_priority,
    }
}

pub(crate) fn apply_nm_dns_setting(
    nm_ip_setting: &mut NmSettingIp,
    dns_conf: &DnsClientState,
) {
    nm_ip_setting.dns = dns_conf.server.clone();
    nm_ip_setting.dns_search = dns_conf.search.clone();
    nm_ip_setting.dns_priority = dns_conf.priority;
}

pub(crate) fn retrieve_dns_info(
    nm_api: &NmApi,
    ifaces: &Interfaces,
) -> Result<DnsState, NmstateError> {
    let mut nm_dns_entires = nm_api
        .get_dns_configuration()
        .map_err(nm_error_to_nmstate)?;
    nm_dns_entires.sort_unstable_by_key(|d| d.priority);
    let mut running_srvs: Vec<String> = Vec::new();
    let mut running_schs: Vec<String> = Vec::new();
    for nm_dns_entry in nm_dns_entires {
        running_srvs.extend_from_slice(nm_dns_entry.name_servers.as_slice());
        running_schs.extend_from_slice(nm_dns_entry.domains.as_slice());
    }

    let mut dns_confs: Vec<&DnsClientState> = Vec::new();
    for iface in ifaces.kernel_ifaces.values() {
        if let Some(ip_conf) = iface.base_iface().ipv6.as_ref() {
            if let Some(dns_conf) = ip_conf.dns.as_ref() {
                dns_confs.push(dns_conf);
            }
        }
        if let Some(ip_conf) = iface.base_iface().ipv4.as_ref() {
            if let Some(dns_conf) = ip_conf.dns.as_ref() {
                dns_confs.push(dns_conf);
            }
        }
    }
    dns_confs.sort_unstable_by_key(|d| d.priority.unwrap_or_default());
    let mut config_srvs: Vec<String> = Vec::new();
    let mut config_schs: Vec<String> = Vec::new();
    for dns_conf in dns_confs {
        if let Some(srvs) = dns_conf.server.as_ref() {
            config_srvs.extend_from_slice(srvs);
        }
        if let Some(schs) = dns_conf.search.as_ref() {
            config_schs.extend_from_slice(schs);
        }
    }

    Ok(DnsState {
        running: Some(DnsClientState {
            server: Some(running_srvs),
            search: Some(running_schs),
            ..Default::default()
        }),
        config: Some(DnsClientState {
            server: if config_srvs.is_empty() && config_schs.is_empty() {
                None
            } else {
                Some(config_srvs.clone())
            },
            search: if config_srvs.is_empty() && config_schs.is_empty() {
                None
            } else {
                Some(config_schs)
            },
            ..Default::default()
        }),
    })
}
