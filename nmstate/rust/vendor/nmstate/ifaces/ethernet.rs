use serde::{Deserialize, Serialize};

use crate::{
    BaseInterface, ErrorKind, Interface, InterfaceState, InterfaceType,
    Interfaces, NmstateError, SrIovConfig,
};

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[non_exhaustive]
pub struct EthernetInterface {
    #[serde(flatten)]
    pub base: BaseInterface,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ethernet: Option<EthernetConfig>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub veth: Option<VethConfig>,
}

impl Default for EthernetInterface {
    fn default() -> Self {
        let mut base = BaseInterface::new();
        base.iface_type = InterfaceType::Ethernet;
        Self {
            base,
            ethernet: None,
            veth: None,
        }
    }
}

impl EthernetInterface {
    pub(crate) fn update_ethernet(&mut self, other: &EthernetInterface) {
        if let Some(eth_conf) = &mut self.ethernet {
            eth_conf.update(other.ethernet.as_ref())
        } else {
            self.ethernet = other.ethernet.clone()
        }
    }

    pub(crate) fn update_veth(&mut self, other: &EthernetInterface) {
        if let Some(veth_conf) = &mut self.veth {
            veth_conf.update(other.veth.as_ref());
        } else {
            self.veth = other.veth.clone();
        }
    }

    pub(crate) fn pre_verify_cleanup(&mut self) {
        if let Some(eth_conf) = self.ethernet.as_mut() {
            eth_conf.pre_verify_cleanup()
        }
        self.base.iface_type = InterfaceType::Ethernet;
    }

    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn sriov_is_enabled(&self) -> bool {
        self.ethernet
            .as_ref()
            .and_then(|eth_conf| {
                eth_conf.sr_iov.as_ref().map(SrIovConfig::sriov_is_enabled)
            })
            .unwrap_or_default()
    }

    pub(crate) fn verify_sriov(
        &self,
        cur_ifaces: &Interfaces,
    ) -> Result<(), NmstateError> {
        if let Some(eth_conf) = &self.ethernet {
            if let Some(sriov_conf) = &eth_conf.sr_iov {
                sriov_conf.verify_sriov(self.base.name.as_str(), cur_ifaces)?;
            }
        }
        Ok(())
    }

    // veth config is ignored unless iface type is veth
    pub(crate) fn veth_sanitize(&mut self) {
        if self.base.iface_type == InterfaceType::Ethernet
            && self.veth.is_some()
        {
            log::warn!(
                "Veth configuration is ignored, please set interface type \
                    to InterfaceType::Veth (veth) to change veth \
                    configuration"
            );
            self.veth = None;
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum EthernetDuplex {
    Full,
    Half,
}

impl std::fmt::Display for EthernetDuplex {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Full => "full",
                Self::Half => "half",
            }
        )
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "kebab-case", deny_unknown_fields)]
#[non_exhaustive]
pub struct EthernetConfig {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub sr_iov: Option<SrIovConfig>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        rename = "auto-negotiation",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub auto_neg: Option<bool>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub speed: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub duplex: Option<EthernetDuplex>,
}

impl EthernetConfig {
    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn update(&mut self, other: Option<&EthernetConfig>) {
        if let Some(other) = other {
            if let Some(sr_iov_conf) = &mut self.sr_iov {
                sr_iov_conf.update(other.sr_iov.as_ref())
            } else {
                self.sr_iov = other.sr_iov.clone()
            }
        }
    }

    pub(crate) fn pre_verify_cleanup(&mut self) {
        if self.auto_neg == Some(true) {
            self.speed = None;
            self.duplex = None;
        }
        if let Some(sriov_conf) = self.sr_iov.as_mut() {
            sriov_conf.pre_verify_cleanup()
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[non_exhaustive]
pub struct VethConfig {
    pub peer: String,
}

impl VethConfig {
    fn update(&mut self, other: Option<&VethConfig>) {
        if let Some(other) = other {
            self.peer = other.peer.clone();
        }
    }
}

// Raise error if new veth interface has no peer defined.
// Mark old veth peer as absent when veth changed its peer.
// Mark veth peer as absent also when veth is marked as absent.
pub(crate) fn handle_veth_peer_changes(
    add_ifaces: &Interfaces,
    chg_ifaces: &mut Interfaces,
    del_ifaces: &mut Interfaces,
    current: &Interfaces,
) -> Result<(), NmstateError> {
    for iface in add_ifaces
        .kernel_ifaces
        .values()
        .filter(|i| i.iface_type() == InterfaceType::Veth)
    {
        if let Interface::Ethernet(eth_iface) = iface {
            if eth_iface.veth.is_none() {
                let e = NmstateError::new(
                    ErrorKind::InvalidArgument,
                    format!(
                        "Veth interface {} does not exists, \
                        peer name is required for creating it",
                        iface.name()
                    ),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }
    }
    for (iface_name, iface) in chg_ifaces.kernel_ifaces.iter() {
        if let Interface::Ethernet(eth_iface) = iface {
            let cur_eth_iface = if let Some(Interface::Ethernet(i)) =
                current.kernel_ifaces.get(iface_name)
            {
                i
            } else {
                continue;
            };
            if let (Some(veth_conf), Some(cur_veth_conf)) =
                (eth_iface.veth.as_ref(), cur_eth_iface.veth.as_ref())
            {
                if veth_conf.peer != cur_veth_conf.peer {
                    del_ifaces.push(new_absent_eth_iface(
                        cur_veth_conf.peer.as_str(),
                    ));
                }
            }
        }
    }

    let mut del_peers: Vec<&str> = Vec::new();
    for iface in del_ifaces
        .kernel_ifaces
        .values()
        .filter(|i| matches!(i, Interface::Ethernet(_)))
    {
        if let Some(Interface::Ethernet(cur_eth_iface)) =
            current.kernel_ifaces.get(iface.name())
        {
            if let Some(veth_conf) = cur_eth_iface.veth.as_ref() {
                del_peers.push(veth_conf.peer.as_str());
            }
        }
    }
    for del_peer in del_peers {
        if !del_ifaces.kernel_ifaces.contains_key(del_peer) {
            del_ifaces.push(new_absent_eth_iface(del_peer));
        }
    }
    Ok(())
}

fn new_absent_eth_iface(name: &str) -> Interface {
    let mut iface = EthernetInterface::new();
    iface.base = BaseInterface {
        name: name.to_string(),
        iface_type: InterfaceType::Ethernet,
        state: InterfaceState::Absent,
        ..Default::default()
    };
    Interface::Ethernet(iface)
}
