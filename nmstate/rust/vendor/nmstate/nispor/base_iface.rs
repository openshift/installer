use crate::{
    nispor::ethtool::np_ethtool_to_nmstate,
    nispor::ip::{np_ipv4_to_nmstate, np_ipv6_to_nmstate},
    BaseInterface, InterfaceState, InterfaceType,
};

fn np_iface_type_to_nmstate(
    np_iface_type: &nispor::IfaceType,
) -> InterfaceType {
    match np_iface_type {
        nispor::IfaceType::Bond => InterfaceType::Bond,
        nispor::IfaceType::Bridge => InterfaceType::LinuxBridge,
        nispor::IfaceType::Dummy => InterfaceType::Dummy,
        nispor::IfaceType::Ethernet => InterfaceType::Ethernet,
        nispor::IfaceType::Loopback => InterfaceType::Loopback,
        nispor::IfaceType::MacVlan => InterfaceType::MacVlan,
        nispor::IfaceType::MacVtap => InterfaceType::MacVtap,
        nispor::IfaceType::OpenvSwitch => InterfaceType::OvsInterface,
        nispor::IfaceType::Veth => InterfaceType::Veth,
        nispor::IfaceType::Vlan => InterfaceType::Vlan,
        nispor::IfaceType::Vrf => InterfaceType::Vrf,
        nispor::IfaceType::Vxlan => InterfaceType::Vxlan,
        nispor::IfaceType::Ipoib => InterfaceType::InfiniBand,
        _ => InterfaceType::Other(format!("{:?}", np_iface_type)),
    }
}

impl From<(&nispor::IfaceState, &[nispor::IfaceFlags])> for InterfaceState {
    fn from(tuple: (&nispor::IfaceState, &[nispor::IfaceFlags])) -> Self {
        let (state, flags) = tuple;
        if *state == nispor::IfaceState::Up
            || flags.contains(&nispor::IfaceFlags::Up)
            || flags.contains(&nispor::IfaceFlags::Running)
        {
            InterfaceState::Up
        } else if *state == nispor::IfaceState::Down {
            InterfaceState::Down
        } else {
            InterfaceState::Unknown
        }
    }
}

pub(crate) fn np_iface_to_base_iface(
    np_iface: &nispor::Iface,
    running_config_only: bool,
) -> BaseInterface {
    let mut base_iface = BaseInterface {
        name: np_iface.name.to_string(),
        state: (&np_iface.state, np_iface.flags.as_slice()).into(),
        iface_type: np_iface_type_to_nmstate(&np_iface.iface_type),
        ipv4: np_ipv4_to_nmstate(np_iface, running_config_only),
        ipv6: np_ipv6_to_nmstate(np_iface, running_config_only),
        mac_address: Some(np_iface.mac_address.to_uppercase()),
        permanent_mac_address: get_permanent_mac_address(np_iface),
        controller: np_iface.controller.as_ref().map(|c| c.to_string()),
        mtu: if np_iface.mtu >= 0 {
            Some(np_iface.mtu as u64)
        } else {
            Some(0u64)
        },
        accept_all_mac_addresses: if np_iface
            .flags
            .contains(&nispor::IfaceFlags::Promisc)
        {
            Some(true)
        } else {
            Some(false)
        },
        ethtool: np_ethtool_to_nmstate(np_iface),
        prop_list: vec![
            "name",
            "state",
            "iface_type",
            "ipv4",
            "ipv6",
            "mac_address",
            "controller",
            "mtu",
            "accept_all_mac_addresses",
            "ethtool",
        ],
        ..Default::default()
    };
    if let InterfaceType::Other(t) = &base_iface.iface_type {
        log::info!(
            "Got unsupported interface type {}: {}, ignoring",
            t,
            &base_iface.name
        );
        base_iface.state = InterfaceState::Ignore;
    }

    base_iface
}

fn get_permanent_mac_address(iface: &nispor::Iface) -> Option<String> {
    if iface.permanent_mac_address.is_empty() {
        // Bond port also hold perm_hwaddr which is the mac address before
        // this interface been assgined to bond as subordinate.
        if let Some(bond_port_info) = &iface.bond_subordinate {
            if bond_port_info.perm_hwaddr.is_empty() {
                None
            } else {
                Some(bond_port_info.perm_hwaddr.clone())
            }
        } else {
            None
        }
    } else {
        Some(iface.permanent_mac_address.clone())
    }
}
