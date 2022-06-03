use std::collections::HashMap;

use log::{debug, info, warn};
use serde::{Deserialize, Deserializer, Serialize};

use crate::{
    dns::{
        get_cur_dns_ifaces, is_dns_changed, purge_dns_config,
        reselect_dns_ifaces,
    },
    ifaces::get_ignored_ifaces,
    nispor::{nispor_apply, nispor_retrieve, set_running_hostname},
    nm::{
        nm_apply, nm_checkpoint_create, nm_checkpoint_destroy,
        nm_checkpoint_rollback, nm_checkpoint_timeout_extend, nm_gen_conf,
        nm_retrieve,
    },
    ovsdb::{ovsdb_apply, ovsdb_is_running, ovsdb_retrieve},
    DnsState, ErrorKind, HostNameState, Interface, InterfaceType, Interfaces,
    NmstateError, OvsDbGlobalConfig, RouteRules, Routes,
};

const DEFAULT_ROLLBACK_TIMEOUT: u32 = 60;
const VERIFY_RETRY_INTERVAL_MILLISECONDS: u64 = 1000;
const VERIFY_RETRY_COUNT: usize = 5;
const VERIFY_RETRY_COUNT_SRIOV: usize = 60;
const VERIFY_RETRY_COUNT_KERNEL_MODE: usize = 5;
const MAX_SUPPORTED_INTERFACES: usize = 1000;

#[derive(Clone, Debug, Serialize, Default, PartialEq, Eq)]
#[serde(deny_unknown_fields)]
#[non_exhaustive]
pub struct NetworkState {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub hostname: Option<HostNameState>,
    #[serde(rename = "dns-resolver", default)]
    pub dns: DnsState,
    #[serde(rename = "route-rules", default)]
    pub rules: RouteRules,
    #[serde(default)]
    pub routes: Routes,
    #[serde(default)]
    pub interfaces: Interfaces,
    #[serde(
        default,
        rename = "ovs-db",
        skip_serializing_if = "OvsDbGlobalConfig::is_none"
    )]
    pub ovsdb: OvsDbGlobalConfig,
    #[serde(skip)]
    // Contain a list of struct member name which is defined explicitly in
    // desire state instead of generated.
    pub prop_list: Vec<&'static str>,
    #[serde(skip)]
    // TODO: Hide user space only info when serialize
    kernel_only: bool,
    #[serde(skip)]
    no_verify: bool,
    #[serde(skip)]
    no_commit: bool,
    #[serde(skip)]
    timeout: Option<u32>,
    #[serde(skip)]
    include_secrets: bool,
    #[serde(skip)]
    include_status_data: bool,
    #[serde(skip)]
    running_config_only: bool,
    #[serde(skip)]
    memory_only: bool,
}

impl<'de> Deserialize<'de> for NetworkState {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let mut net_state = NetworkState::new();
        let v = serde_json::Value::deserialize(deserializer)?;
        if let Some(ifaces_value) = v.get("interfaces") {
            net_state.prop_list.push("interfaces");
            net_state.interfaces = Interfaces::deserialize(ifaces_value)
                .map_err(serde::de::Error::custom)?;
        }
        if let Some(dns_value) = v.get("dns-resolver") {
            net_state.prop_list.push("dns");
            net_state.dns = DnsState::deserialize(dns_value)
                .map_err(serde::de::Error::custom)?;
        }
        if let Some(route_value) = v.get("routes") {
            net_state.prop_list.push("routes");
            net_state.routes = Routes::deserialize(route_value)
                .map_err(serde::de::Error::custom)?;
        }
        if let Some(rule_value) = v.get("route-rules") {
            net_state.prop_list.push("rules");
            net_state.rules = RouteRules::deserialize(rule_value)
                .map_err(serde::de::Error::custom)?;
        }
        if let Some(ovsdb_value) = v.get("ovs-db") {
            net_state.prop_list.push("ovsdb");
            net_state.ovsdb = OvsDbGlobalConfig::deserialize(ovsdb_value)
                .map_err(serde::de::Error::custom)?;
        }
        if let Some(hostname_value) = v.get("hostname") {
            net_state.prop_list.push("hostname");
            net_state.hostname = Some(
                HostNameState::deserialize(hostname_value)
                    .map_err(serde::de::Error::custom)?,
            );
        }
        Ok(net_state)
    }
}

impl NetworkState {
    pub(crate) const PASSWORD_HID_BY_NMSTATE: &'static str =
        "<_password_hid_by_nmstate>";

    pub fn set_kernel_only(&mut self, value: bool) -> &mut Self {
        self.kernel_only = value;
        self
    }

    pub fn set_verify_change(&mut self, value: bool) -> &mut Self {
        self.no_verify = !value;
        self
    }

    pub fn set_commit(&mut self, value: bool) -> &mut Self {
        self.no_commit = !value;
        self
    }

    pub fn set_timeout(&mut self, value: u32) -> &mut Self {
        self.timeout = Some(value);
        self
    }

    pub fn set_include_secrets(&mut self, value: bool) -> &mut Self {
        self.include_secrets = value;
        self
    }

    pub fn set_include_status_data(&mut self, value: bool) -> &mut Self {
        self.include_status_data = value;
        self
    }

    // Query activated/running network configuration excluding:
    // * IP address retrieved by DHCP or IPv6 auto configuration.
    // * DNS client resolver retrieved by DHCP or IPv6 auto configuration.
    // * Routes retrieved by DHCPv4 or IPv6 router advertisement.
    // * LLDP neighbor information.
    pub fn set_running_config_only(&mut self, value: bool) -> &mut Self {
        self.running_config_only = value;
        self
    }

    pub fn set_memory_only(&mut self, value: bool) -> &mut Self {
        self.memory_only = value;
        self
    }

    pub fn new() -> Self {
        Default::default()
    }

    // We provide this instead asking use to do serde_json::from_str(), so that
    // we could provide better error NmstateError instead of serde_json one.
    pub fn new_from_json(net_state_json: &str) -> Result<Self, NmstateError> {
        match serde_json::from_str(net_state_json) {
            Ok(s) => Ok(s),
            Err(e) => Err(NmstateError::new(
                ErrorKind::InvalidArgument,
                format!("Invalid json string: {}", e),
            )),
        }
    }

    pub fn append_interface_data(&mut self, iface: Interface) {
        self.interfaces.push(iface);
    }

    pub fn retrieve(&mut self) -> Result<&mut Self, NmstateError> {
        let state = nispor_retrieve(self.running_config_only)?;
        if state.prop_list.contains(&"hostname") {
            self.hostname = state.hostname;
        }
        if state.prop_list.contains(&"interfaces") {
            self.interfaces = state.interfaces;
        }
        if state.prop_list.contains(&"routes") {
            self.routes = state.routes;
        }
        if state.prop_list.contains(&"rules") {
            self.rules = state.rules;
        }
        if !self.kernel_only {
            let nm_state = nm_retrieve(self.running_config_only)?;
            // TODO: Priority handling
            self.update_state(&nm_state);
            if ovsdb_is_running() {
                match ovsdb_retrieve() {
                    Ok(ovsdb_state) => self.update_state(&ovsdb_state),
                    Err(e) => {
                        log::warn!("Failed to retrieve OVS DB state: {}", e);
                    }
                }
            }
        }
        if !self.include_secrets {
            self.hide_secrets();
        }
        Ok(self)
    }

    pub fn hide_secrets(&mut self) {
        self.interfaces.hide_secrets();
    }

    pub fn apply(&self) -> Result<(), NmstateError> {
        let mut desire_state_to_verify = self.clone();
        let mut desire_state_to_apply = self.clone();
        let mut cur_net_state = NetworkState::new();
        cur_net_state.set_kernel_only(self.kernel_only);
        cur_net_state.set_include_secrets(true);
        cur_net_state.retrieve()?;

        if desire_state_to_apply.interfaces.to_vec().len()
            >= MAX_SUPPORTED_INTERFACES
        {
            log::warn!(
                "Interfaces count exceeds the support limit {} in \
                desired state",
                MAX_SUPPORTED_INTERFACES,
            );
        }

        let (ignored_kernel_ifaces, ignored_user_ifaces) =
            get_ignored_ifaces(&self.interfaces, &cur_net_state.interfaces);

        for iface_name in &ignored_kernel_ifaces {
            log::info!("Ignoring kernel interface {}", iface_name)
        }
        for (iface_name, iface_type) in &ignored_user_ifaces {
            log::info!(
                "Ignoring user space interface {} with type {}",
                iface_name,
                iface_type
            )
        }

        desire_state_to_apply.interfaces.remove_ignored_ifaces(
            &ignored_kernel_ifaces,
            &ignored_user_ifaces,
        );
        cur_net_state.interfaces.remove_ignored_ifaces(
            &ignored_kernel_ifaces,
            &ignored_user_ifaces,
        );

        desire_state_to_apply
            .routes
            .remove_ignored_iface_routes(ignored_kernel_ifaces.as_slice());
        cur_net_state
            .routes
            .remove_ignored_iface_routes(ignored_kernel_ifaces.as_slice());

        desire_state_to_verify
            .interfaces
            .resolve_unknown_ifaces(&cur_net_state.interfaces)?;
        desire_state_to_apply
            .interfaces
            .resolve_unknown_ifaces(&cur_net_state.interfaces)?;

        let (add_net_state, chg_net_state, del_net_state) =
            desire_state_to_apply.gen_state_for_apply(&cur_net_state)?;

        debug!("Adding net state {:?}", &add_net_state);
        debug!("Changing net state {:?}", &chg_net_state);
        debug!("Deleting net state {:?}", &del_net_state);

        if let Some(running_hostname) =
            self.hostname.as_ref().and_then(|c| c.running.as_ref())
        {
            set_running_hostname(running_hostname)?;
        }

        if !self.kernel_only {
            let retry_count =
                if desire_state_to_apply.interfaces.has_sriov_enabled() {
                    VERIFY_RETRY_COUNT_SRIOV
                } else {
                    VERIFY_RETRY_COUNT
                };

            let timeout = self.timeout.unwrap_or(DEFAULT_ROLLBACK_TIMEOUT);
            let checkpoint = nm_checkpoint_create(timeout)?;
            info!("Created checkpoint {}", &checkpoint);

            with_nm_checkpoint(&checkpoint, self.no_commit, || {
                nm_apply(
                    &add_net_state,
                    &chg_net_state,
                    &del_net_state,
                    // TODO: Passing full(desire + current) network state
                    // instead of current,
                    &cur_net_state,
                    self,
                    &checkpoint,
                    self.memory_only,
                )?;
                if ovsdb_is_running() {
                    ovsdb_apply(&desire_state_to_apply, &cur_net_state)?;
                }
                nm_checkpoint_timeout_extend(&checkpoint, timeout)?;
                if !self.no_verify {
                    with_retry(
                        VERIFY_RETRY_INTERVAL_MILLISECONDS,
                        retry_count,
                        || {
                            let mut new_cur_net_state = cur_net_state.clone();
                            new_cur_net_state.set_include_secrets(true);
                            new_cur_net_state.retrieve()?;
                            desire_state_to_verify.verify(&new_cur_net_state)
                        },
                    )
                } else {
                    Ok(())
                }
            })
        } else {
            // TODO: Need checkpoint for kernel only mode
            nispor_apply(
                &add_net_state,
                &chg_net_state,
                &del_net_state,
                &cur_net_state,
            )?;
            if !self.no_verify {
                with_retry(
                    VERIFY_RETRY_INTERVAL_MILLISECONDS,
                    VERIFY_RETRY_COUNT_KERNEL_MODE,
                    || {
                        let mut new_cur_net_state = cur_net_state.clone();
                        new_cur_net_state.retrieve()?;
                        desire_state_to_verify.verify(&new_cur_net_state)
                    },
                )
            } else {
                Ok(())
            }
        }
    }

    fn update_state(&mut self, other: &Self) {
        if other.prop_list.contains(&"hostname") {
            if let Some(h) = self.hostname.as_mut() {
                if let Some(other_h) = other.hostname.as_ref() {
                    h.update(other_h);
                }
            } else {
                self.hostname = other.hostname.clone();
            }
        }
        if other.prop_list.contains(&"interfaces") {
            self.interfaces.update(&other.interfaces);
        }
        if other.prop_list.contains(&"dns") {
            self.dns = other.dns.clone();
        }
        if other.prop_list.contains(&"ovsdb") {
            self.ovsdb = other.ovsdb.clone();
        }
    }

    pub fn gen_conf(
        &self,
    ) -> Result<HashMap<String, Vec<String>>, NmstateError> {
        let mut ret = HashMap::new();
        let mut self_clone = self.clone();
        self_clone.interfaces.set_unknown_iface_to_eth();
        self_clone.interfaces.set_missing_port_to_eth();
        let (add_net_state, _, _) =
            self_clone.gen_state_for_apply(&Self::new())?;
        ret.insert("NetworkManager".to_string(), nm_gen_conf(&add_net_state)?);
        Ok(ret)
    }

    fn verify(&self, current: &Self) -> Result<(), NmstateError> {
        if let Some(desired_hostname) = self.hostname.as_ref() {
            desired_hostname.verify(current.hostname.as_ref())?;
        }
        self.interfaces.verify(&current.interfaces)?;
        let (ignored_kernel_ifaces, _) =
            get_ignored_ifaces(&self.interfaces, &current.interfaces);
        self.routes
            .verify(&current.routes, &ignored_kernel_ifaces)?;
        self.rules.verify(&current.rules)?;
        self.dns.verify(&current.dns)?;
        self.ovsdb.verify(&current.ovsdb)
    }

    // Return three NetworkState:
    //  * State for addition.
    //  * State for change.
    //  * State for deletion.
    // This function is the entry point for decision making which
    // expanding complex desire network layout to flat network layout.
    pub(crate) fn gen_state_for_apply(
        &self,
        current: &Self,
    ) -> Result<(Self, Self, Self), NmstateError> {
        self.routes.validate()?;
        self.rules.validate()?;
        self.dns.validate()?;

        let mut add_net_state = NetworkState::new();
        let mut chg_net_state = NetworkState::new();
        let mut del_net_state = NetworkState::new();

        let mut ifaces = self.interfaces.clone();

        let (add_ifaces, chg_ifaces, del_ifaces) = ifaces
            .gen_state_for_apply(&current.interfaces, self.memory_only)?;

        add_net_state.interfaces = add_ifaces;
        add_net_state.hostname = self.hostname.clone();
        chg_net_state.interfaces = chg_ifaces;
        del_net_state.interfaces = del_ifaces;

        self.include_route_changes(
            &mut add_net_state,
            &mut chg_net_state,
            &del_net_state,
            current,
        );

        self.include_rule_changes(
            &mut add_net_state,
            &mut chg_net_state,
            &del_net_state,
            current,
        )?;

        self.include_dns_changes(
            &mut add_net_state,
            &mut chg_net_state,
            current,
        )?;

        Ok((add_net_state, chg_net_state, del_net_state))
    }

    fn include_route_changes(
        &self,
        add_net_state: &mut Self,
        chg_net_state: &mut Self,
        del_net_state: &Self,
        current: &Self,
    ) {
        let mut changed_iface_routes =
            self.routes.gen_changed_ifaces_and_routes(&current.routes);

        for (iface_name, routes) in changed_iface_routes.drain() {
            let cur_iface = current
                .interfaces
                .get_iface(&iface_name, InterfaceType::Unknown);
            if del_net_state
                .interfaces
                .kernel_ifaces
                .get(&iface_name)
                .is_some()
            {
                // Ignore routes on absent interfaces.
                continue;
            }

            if let Some(iface) =
                add_net_state.interfaces.kernel_ifaces.get_mut(&iface_name)
            {
                // Desire interface might not have IP information defined
                // in that case, we copy from current
                iface.base_iface_mut().routes = Some(routes);
                if let Some(cur_iface) = cur_iface {
                    iface
                        .base_iface_mut()
                        .copy_ip_config_if_none(cur_iface.base_iface());
                }
            } else if let Some(iface) =
                chg_net_state.interfaces.kernel_ifaces.get_mut(&iface_name)
            {
                iface.base_iface_mut().routes = Some(routes);
                if let Some(cur_iface) = cur_iface {
                    iface
                        .base_iface_mut()
                        .copy_ip_config_if_none(cur_iface.base_iface());
                }
            } else if let Some(cur_iface) = cur_iface {
                // Interface not mentioned in desire state but impacted by
                // wildcard absent route
                let mut new_iface = cur_iface.clone_name_type_only();
                new_iface
                    .base_iface_mut()
                    .copy_ip_config_if_none(cur_iface.base_iface());
                new_iface.base_iface_mut().routes = Some(routes);
                chg_net_state.append_interface_data(new_iface);
            } else {
                warn!(
                    "The next hop interface of desired routes {:?} \
                        does not exist",
                    routes
                );
            }
        }
    }

    fn include_rule_changes(
        &self,
        add_net_state: &mut Self,
        chg_net_state: &mut Self,
        del_net_state: &Self,
        current: &Self,
    ) -> Result<(), NmstateError> {
        let mut changed_rules =
            self.rules.gen_rule_changed_table_ids(&current.rules);

        // Convert table id to interface name
        for (table_id, rules) in changed_rules.drain() {
            // We does not differentiate the IPv4 and IPv6 route table ID.
            // The verification process will find out the error.
            // We did not head any use case been limited by this.
            let iface_name =
                self.get_iface_name_for_route_table(current, table_id)?;
            let cur_iface = current
                .interfaces
                .get_iface(&iface_name, InterfaceType::Unknown);
            if del_net_state
                .interfaces
                .kernel_ifaces
                .get(&iface_name)
                .is_some()
            {
                // Ignore rules on absent interfaces.
                continue;
            }
            if let Some(iface) =
                add_net_state.interfaces.kernel_ifaces.get_mut(&iface_name)
            {
                if let Some(ex_rules) = iface.base_iface_mut().rules.as_mut() {
                    ex_rules.extend(rules);
                    ex_rules.sort_unstable();
                    ex_rules.dedup();
                } else {
                    iface.base_iface_mut().rules = Some(rules);
                }
                if let Some(cur_iface) = cur_iface {
                    iface
                        .base_iface_mut()
                        .copy_ip_config_if_none(cur_iface.base_iface());
                }
            } else if let Some(iface) =
                chg_net_state.interfaces.kernel_ifaces.get_mut(&iface_name)
            {
                if let Some(ex_rules) = iface.base_iface_mut().rules.as_mut() {
                    ex_rules.extend(rules);
                    ex_rules.sort_unstable();
                    ex_rules.dedup();
                } else {
                    iface.base_iface_mut().rules = Some(rules);
                }
                if let Some(cur_iface) = cur_iface {
                    iface
                        .base_iface_mut()
                        .copy_ip_config_if_none(cur_iface.base_iface());
                }
            } else if let Some(cur_iface) = cur_iface {
                // Interface not mentioned in desire state but impacted by
                // wildcard absent route rule
                let mut new_iface = cur_iface.clone_name_type_only();
                new_iface
                    .base_iface_mut()
                    .copy_ip_config_if_none(cur_iface.base_iface());
                new_iface.base_iface_mut().rules = Some(rules);
                chg_net_state.append_interface_data(new_iface);
            } else {
                let e = NmstateError::new(
                    ErrorKind::InvalidArgument,
                    format!(
                    "Failed to find a interface for desired routes rules {:?} ",
                    rules
                ),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }
        Ok(())
    }

    // The whole design of this DNS setting is matching NetworkManager
    // design where DNS information is saved into each interface.
    // When different network control backend nmstate supported, we need to
    // move these code to NetworkManager plugin. Currently, we keep it
    // here for consistency with route/route rule.
    fn include_dns_changes(
        &self,
        add_net_state: &mut Self,
        chg_net_state: &mut Self,
        current: &Self,
    ) -> Result<(), NmstateError> {
        let mut self_clone = self.clone();
        self_clone.dns.merge_current(&current.dns);

        if is_dns_changed(&self_clone, current) {
            let (v4_iface_name, v6_iface_name) =
                reselect_dns_ifaces(&self_clone, current);
            let (cur_v4_ifaces, cur_v6_ifaces) =
                get_cur_dns_ifaces(&current.interfaces);
            if let Some(dns_conf) = &self_clone.dns.config {
                if dns_conf.is_purge() {
                    purge_dns_config(
                        false,
                        cur_v4_ifaces.as_slice(),
                        self,
                        chg_net_state,
                        current,
                    );
                    purge_dns_config(
                        true,
                        cur_v6_ifaces.as_slice(),
                        self,
                        chg_net_state,
                        current,
                    );
                } else {
                    purge_dns_config(
                        false,
                        &cur_v4_ifaces,
                        self,
                        chg_net_state,
                        current,
                    );
                    purge_dns_config(
                        true,
                        &cur_v6_ifaces,
                        self,
                        chg_net_state,
                        current,
                    );
                    dns_conf.save_dns_to_iface(
                        &v4_iface_name,
                        &v6_iface_name,
                        add_net_state,
                        chg_net_state,
                        current,
                    )?;
                }
            }
        } else {
            log::debug!("DNS configuration unchanged");
        }
        Ok(())
    }

    fn _get_iface_name_for_route_table(&self, table_id: u32) -> Option<String> {
        if let Some(routes) = self.routes.config.as_ref() {
            for route in routes {
                if route.table_id == Some(table_id) {
                    if let Some(iface_name) = route.next_hop_iface.as_ref() {
                        return Some(iface_name.to_string());
                    }
                }
            }
        }
        // We need to differentiate IPv4 and IPv6 auto route table ID when
        // user case shows up. Currently, we just assume user does not
        // mix up the table number for IPv4 and IPv6 between interfaces.
        for iface in self.interfaces.kernel_ifaces.values() {
            if iface
                .base_iface()
                .ipv6
                .as_ref()
                .and_then(|c| c.auto_table_id)
                .or_else(|| {
                    iface
                        .base_iface()
                        .ipv4
                        .as_ref()
                        .and_then(|c| c.auto_table_id)
                })
                == Some(table_id)
            {
                return Some(iface.name().to_string());
            }
        }
        None
    }

    // * Find desired interface with static route to given table ID.
    // * Find desired interface with dynamic route to given table ID.
    // * Find current interface with static route to given table ID.
    // * Find current interface with dynamic route to given table ID.
    fn get_iface_name_for_route_table(
        &self,
        current: &Self,
        table_id: u32,
    ) -> Result<String, NmstateError> {
        match self
            ._get_iface_name_for_route_table(table_id)
            .or_else(|| current._get_iface_name_for_route_table(table_id))
        {
            Some(iface_name) => Ok(iface_name),
            None => {
                let e = NmstateError::new(
                    ErrorKind::InvalidArgument,
                    format!(
                        "Route table {} for route rule is not defined by \
                        any routes",
                        table_id
                    ),
                );
                log::error!("{}", e);
                Err(e)
            }
        }
    }
    pub fn checkpoint_rollback(checkpoint: &str) -> Result<(), NmstateError> {
        nm_checkpoint_rollback(checkpoint)
    }

    pub fn checkpoint_commit(checkpoint: &str) -> Result<(), NmstateError> {
        nm_checkpoint_destroy(checkpoint)
    }

    pub(crate) fn get_kernel_iface_with_route(
        &self,
        iface_name: &str,
    ) -> Option<Interface> {
        if let Some(iface) = self.interfaces.kernel_ifaces.get(iface_name) {
            let mut ret = iface.clone();
            ret.base_iface_mut().routes =
                self.routes.get_config_routes_of_iface(iface_name);
            Some(ret)
        } else {
            None
        }
    }
}

fn with_nm_checkpoint<T>(
    checkpoint: &str,
    no_commit: bool,
    func: T,
) -> Result<(), NmstateError>
where
    T: FnOnce() -> Result<(), NmstateError>,
{
    match func() {
        Ok(()) => {
            if !no_commit {
                nm_checkpoint_destroy(checkpoint)?;

                info!("Destroyed checkpoint {}", checkpoint);
            } else {
                info!("Skipping commit for checkpoint {}", checkpoint);
            }
            Ok(())
        }
        Err(e) => {
            if let Err(e) = nm_checkpoint_rollback(checkpoint) {
                warn!("nm_checkpoint_rollback() failed: {}", e);
            }
            info!("Rollbacked to checkpoint {}", checkpoint);
            Err(e)
        }
    }
}

fn with_retry<T>(
    interval_ms: u64,
    count: usize,
    func: T,
) -> Result<(), NmstateError>
where
    T: FnOnce() -> Result<(), NmstateError> + Copy,
{
    let mut cur_count = 0usize;
    while cur_count < count {
        if let Err(e) = func() {
            if cur_count == count - 1 {
                return Err(e);
            } else {
                info!("Retrying on verification failure: {}", e);
                std::thread::sleep(std::time::Duration::from_millis(
                    interval_ms,
                ));
                cur_count += 1;
                continue;
            }
        } else {
            return Ok(());
        }
    }
    Ok(())
}
