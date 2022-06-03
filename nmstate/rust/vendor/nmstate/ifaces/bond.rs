use crate::{BaseInterface, ErrorKind, Interface, InterfaceType, NmstateError};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq, Eq)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub struct BondInterface {
    #[serde(flatten)]
    pub base: BaseInterface,
    #[serde(
        skip_serializing_if = "Option::is_none",
        rename = "link-aggregation"
    )]
    pub bond: Option<BondConfig>,
}

impl Default for BondInterface {
    fn default() -> Self {
        let mut base = BaseInterface::new();
        base.iface_type = InterfaceType::Bond;
        Self { base, bond: None }
    }
}

impl BondInterface {
    pub(crate) fn update_bond(&mut self, other: &BondInterface) {
        if let Some(bond_conf) = &mut self.bond {
            bond_conf.update(other.bond.as_ref());
        } else {
            self.bond = other.bond.clone();
        }
    }

    // Return None when desire state does not mention ports
    pub(crate) fn ports(&self) -> Option<Vec<&str>> {
        self.bond
            .as_ref()
            .and_then(|bond_conf| bond_conf.port.as_ref())
            .map(|ports| ports.as_slice().iter().map(|p| p.as_str()).collect())
    }

    pub(crate) fn mode(&self) -> Option<BondMode> {
        self.bond.as_ref().and_then(|bond_conf| bond_conf.mode)
    }

    pub(crate) fn pre_verify_cleanup(&mut self) {
        self.drop_empty_arp_ip_target();
        self.sort_ports();
    }

    pub fn new() -> Self {
        Self::default()
    }

    fn is_mac_restricted_mode(&self) -> bool {
        self.bond
            .as_ref()
            .and_then(|bond_conf| {
                if self.mode() == Some(BondMode::ActiveBackup) {
                    bond_conf.options.as_ref()
                } else {
                    None
                }
            })
            .and_then(|bond_opts| bond_opts.fail_over_mac)
            == Some(BondFailOverMac::Active)
    }

    fn is_not_mac_restricted_mode_explicitly(&self) -> bool {
        (self.mode().is_some() && self.mode() != Some(BondMode::ActiveBackup))
            || ![None, Some(BondFailOverMac::Active)].contains(
                &self
                    .bond
                    .as_ref()
                    .and_then(|bond_conf| bond_conf.options.as_ref())
                    .and_then(|bond_opts| bond_opts.fail_over_mac),
            )
    }

    fn sort_ports(&mut self) {
        if let Some(ref mut bond_conf) = self.bond {
            if let Some(ref mut port_conf) = &mut bond_conf.port {
                port_conf.sort_unstable_by_key(|p| p.clone())
            }
        }
    }

    fn drop_empty_arp_ip_target(&mut self) {
        if let Some(ref mut bond_conf) = self.bond {
            if let Some(ref mut bond_opts) = &mut bond_conf.options {
                if let Some(ref mut arp_ip_target) = bond_opts.arp_ip_target {
                    if arp_ip_target.is_empty() {
                        bond_opts.arp_ip_target = None;
                    }
                }
            }
        }
    }

    pub(crate) fn validate(
        &self,
        current: Option<&Interface>,
    ) -> Result<(), NmstateError> {
        self.base.validate()?;
        self.validate_new_iface_with_no_mode(current)?;
        self.validate_mac_restricted_mode(current)?;
        if let Some(bond_conf) = &self.bond {
            bond_conf.validate()?;
        }
        Ok(())
    }

    pub(crate) fn remove_port(&mut self, port_to_remove: &str) {
        if let Some(index) = self.bond.as_ref().and_then(|bond_conf| {
            bond_conf.port.as_ref().and_then(|ports| {
                ports
                    .iter()
                    .position(|port_name| port_name == port_to_remove)
            })
        }) {
            self.bond
                .as_mut()
                .and_then(|bond_conf| bond_conf.port.as_mut())
                .map(|ports| ports.remove(index));
        }
    }

    fn validate_new_iface_with_no_mode(
        &self,
        current: Option<&Interface>,
    ) -> Result<(), NmstateError> {
        if current.is_none() && self.mode().is_none() {
            let e = NmstateError::new(
                ErrorKind::InvalidArgument,
                format!(
                    "Bond mode is mandatory for new bond interface: {}",
                    &self.base.name
                ),
            );
            log::error!("{}", e);
            return Err(e);
        }
        Ok(())
    }

    // Fail on
    // * Desire mac restricted mode with mac defined
    // * Desire mac address with current interface in mac restricted mode with
    //   desired not changing mac restricted mode
    fn validate_mac_restricted_mode(
        &self,
        current: Option<&Interface>,
    ) -> Result<(), NmstateError> {
        let e = NmstateError::new(
            ErrorKind::InvalidArgument,
            "MAC address cannot be specified in bond interface along with \
            fail_over_mac active on active backup mode"
                .to_string(),
        );
        if self.is_mac_restricted_mode() && self.base.mac_address.is_some() {
            log::error!("{}", e);
            return Err(e);
        }

        if let Some(Interface::Bond(current)) = current {
            if current.is_mac_restricted_mode()
                && self.base.mac_address.is_some()
                && !self.is_not_mac_restricted_mode_explicitly()
            {
                log::error!("{}", e);
                return Err(e);
            }
        }
        Ok(())
    }

    pub(crate) fn is_options_reset(&self) -> bool {
        if let Some(bond_opts) = self
            .bond
            .as_ref()
            .and_then(|bond_conf| bond_conf.options.as_ref())
        {
            bond_opts == &BondOptions::default()
        } else {
            false
        }
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone, Copy)]
#[non_exhaustive]
pub enum BondMode {
    #[serde(rename = "balance-rr")]
    RoundRobin,
    #[serde(rename = "active-backup")]
    ActiveBackup,
    #[serde(rename = "balance-xor")]
    XOR,
    #[serde(rename = "broadcast")]
    Broadcast,
    #[serde(rename = "802.3ad")]
    LACP,
    #[serde(rename = "balance-tlb")]
    TLB,
    #[serde(rename = "balance-alb")]
    ALB,
    #[serde(rename = "unknown")]
    Unknown,
}

impl Default for BondMode {
    fn default() -> Self {
        Self::RoundRobin
    }
}

impl std::fmt::Display for BondMode {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                BondMode::RoundRobin => "balance-rr",
                BondMode::ActiveBackup => "active-backup",
                BondMode::XOR => "balance-xor",
                BondMode::Broadcast => "broadcast",
                BondMode::LACP => "802.3ad",
                BondMode::TLB => "balance-tlb",
                BondMode::ALB => "balance-alb",
                BondMode::Unknown => "unknown",
            }
        )
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "kebab-case", deny_unknown_fields)]
#[non_exhaustive]
pub struct BondConfig {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub mode: Option<BondMode>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub options: Option<BondOptions>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub port: Option<Vec<String>>,
}

impl BondConfig {
    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn validate(&self) -> Result<(), NmstateError> {
        if let Some(opts) = &self.options {
            opts.validate()?;
        }
        Ok(())
    }

    pub(crate) fn update(&mut self, other: Option<&BondConfig>) {
        if let Some(other) = other {
            self.mode = other.mode;
            self.options = other.options.clone();
            self.port = other.port.clone();
        }
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone, Copy)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum BondAdSelect {
    #[serde(alias = "0")]
    Stable,
    #[serde(alias = "1")]
    Bandwidth,
    #[serde(alias = "2")]
    Count,
}

impl From<BondAdSelect> for u8 {
    fn from(v: BondAdSelect) -> u8 {
        match v {
            BondAdSelect::Stable => 0,
            BondAdSelect::Bandwidth => 1,
            BondAdSelect::Count => 2,
        }
    }
}

impl std::fmt::Display for BondAdSelect {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Stable => "stable",
                Self::Bandwidth => "bandwidth",
                Self::Count => "count",
            }
        )
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum BondLacpRate {
    #[serde(alias = "0")]
    Slow,
    #[serde(alias = "1")]
    Fast,
}

impl std::fmt::Display for BondLacpRate {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Slow => "slow",
                Self::Fast => "fast",
            }
        )
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone, Copy)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum BondAllPortsActive {
    #[serde(alias = "0")]
    Dropped,
    #[serde(alias = "1")]
    Delivered,
}

impl std::fmt::Display for BondAllPortsActive {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Dropped => "dropped",
                Self::Delivered => "delivered",
            }
        )
    }
}

impl From<BondAllPortsActive> for u8 {
    fn from(v: BondAllPortsActive) -> u8 {
        match v {
            BondAllPortsActive::Dropped => 0,
            BondAllPortsActive::Delivered => 1,
        }
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum BondArpAllTargets {
    #[serde(alias = "0")]
    Any,
    #[serde(alias = "1")]
    All,
}

impl std::fmt::Display for BondArpAllTargets {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Any => "any",
                Self::All => "all",
            }
        )
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum BondArpValidate {
    #[serde(alias = "0")]
    None,
    #[serde(alias = "1")]
    Active,
    #[serde(alias = "2")]
    Backup,
    #[serde(alias = "3")]
    All,
    #[serde(alias = "4")]
    Filter,
    #[serde(rename = "filter_active", alias = "5")]
    FilterActive,
    #[serde(rename = "filter_backup", alias = "6")]
    FilterBackup,
}

impl std::fmt::Display for BondArpValidate {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::None => "none",
                Self::Active => "active",
                Self::Backup => "backup",
                Self::All => "all",
                Self::Filter => "filter",
                Self::FilterActive => "filter_active",
                Self::FilterBackup => "filter_backup",
            }
        )
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone, Copy)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum BondFailOverMac {
    #[serde(alias = "0")]
    None,
    #[serde(alias = "1")]
    Active,
    #[serde(alias = "2")]
    Follow,
}

impl std::fmt::Display for BondFailOverMac {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::None => "none",
                Self::Active => "active",
                Self::Follow => "follow ",
            }
        )
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum BondPrimaryReselect {
    #[serde(alias = "0")]
    Always,
    #[serde(alias = "1")]
    Better,
    #[serde(alias = "2")]
    Failure,
}

impl std::fmt::Display for BondPrimaryReselect {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Always => "always",
                Self::Better => "better",
                Self::Failure => "failure",
            }
        )
    }
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone)]
#[non_exhaustive]
pub enum BondXmitHashPolicy {
    #[serde(rename = "layer2")]
    #[serde(alias = "0")]
    Layer2,
    #[serde(rename = "layer3+4")]
    #[serde(alias = "1")]
    Layer34,
    #[serde(rename = "layer2+3")]
    #[serde(alias = "2")]
    Layer23,
    #[serde(rename = "encap2+3")]
    #[serde(alias = "3")]
    Encap23,
    #[serde(rename = "encap3+4")]
    #[serde(alias = "4")]
    Encap34,
    #[serde(rename = "vlan+srcmac")]
    #[serde(alias = "5")]
    VlanSrcMac,
}

impl BondXmitHashPolicy {
    pub fn to_u8(&self) -> u8 {
        match self {
            Self::Layer2 => 0,
            Self::Layer34 => 1,
            Self::Layer23 => 2,
            Self::Encap23 => 3,
            Self::Encap34 => 4,
            Self::VlanSrcMac => 5,
        }
    }
}

impl std::fmt::Display for BondXmitHashPolicy {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Layer2 => "layer2",
                Self::Layer34 => "layer3+4",
                Self::Layer23 => "layer2+3",
                Self::Encap23 => "encap2+3",
                Self::Encap34 => "encap3+4",
                Self::VlanSrcMac => "vlan+srcmac",
            }
        )
    }
}

#[derive(Debug, Serialize, Deserialize, Default, Clone, PartialEq, Eq)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct BondOptions {
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u16_or_string"
    )]
    pub ad_actor_sys_prio: Option<u16>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ad_actor_system: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ad_select: Option<BondAdSelect>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u16_or_string"
    )]
    pub ad_user_port_key: Option<u16>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub all_slaves_active: Option<BondAllPortsActive>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub arp_all_targets: Option<BondArpAllTargets>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub arp_interval: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub arp_ip_target: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub arp_validate: Option<BondArpValidate>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub downdelay: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub fail_over_mac: Option<BondFailOverMac>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub lacp_rate: Option<BondLacpRate>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub lp_interval: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub miimon: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub min_links: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u8_or_string"
    )]
    pub num_grat_arp: Option<u8>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u8_or_string"
    )]
    pub num_unsol_na: Option<u8>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub packets_per_slave: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub primary: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub primary_reselect: Option<BondPrimaryReselect>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub resend_igmp: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub tlb_dynamic_lb: Option<bool>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub updelay: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub use_carrier: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub xmit_hash_policy: Option<BondXmitHashPolicy>,
}

impl BondOptions {
    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn validate(&self) -> Result<(), NmstateError> {
        self.validate_ad_actor_system_mac_address()?;
        self.validate_miimon_and_arp_interval()?;
        Ok(())
    }

    fn validate_ad_actor_system_mac_address(&self) -> Result<(), NmstateError> {
        if let Some(ad_actor_system) = &self.ad_actor_system {
            if ad_actor_system.to_uppercase().starts_with("01:00:5E") {
                let e = NmstateError::new(
                    ErrorKind::InvalidArgument,
                    "The ad_actor_system bond option cannot be an IANA \
                    multicast address(prefix with 01:00:5E)"
                        .to_string(),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }
        Ok(())
    }

    fn validate_miimon_and_arp_interval(&self) -> Result<(), NmstateError> {
        if let (Some(miimon), Some(arp_interval)) =
            (self.miimon, self.arp_interval)
        {
            if miimon > 0 && arp_interval > 0 {
                let e = NmstateError::new(
                    ErrorKind::InvalidArgument,
                    "Bond miimon and arp interval are not compatible options."
                        .to_string(),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }
        Ok(())
    }
}
