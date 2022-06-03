use std::collections::HashMap;

use crate::nm::nm_dbus::{
    NmActiveConnection, NmApi, NmConnection, NmDevice, NmDeviceState,
    NmLldpNeighbor,
};
use log::{debug, warn};

use crate::{
    nm::active_connection::create_index_for_nm_acs_by_name_type,
    nm::connection::{
        create_index_for_nm_conns_by_ctrler_type,
        create_index_for_nm_conns_by_name_type, get_port_nm_conns,
        NM_SETTING_BOND_SETTING_NAME, NM_SETTING_BRIDGE_SETTING_NAME,
        NM_SETTING_CONTROLLERS, NM_SETTING_DUMMY_SETTING_NAME,
        NM_SETTING_INFINIBAND_SETTING_NAME, NM_SETTING_MACVLAN_SETTING_NAME,
        NM_SETTING_OVS_BRIDGE_SETTING_NAME, NM_SETTING_OVS_IFACE_SETTING_NAME,
        NM_SETTING_VETH_SETTING_NAME, NM_SETTING_VLAN_SETTING_NAME,
        NM_SETTING_VRF_SETTING_NAME, NM_SETTING_VXLAN_SETTING_NAME,
        NM_SETTING_WIRED_SETTING_NAME,
    },
    nm::dns::retrieve_dns_info,
    nm::error::nm_error_to_nmstate,
    nm::ieee8021x::nm_802_1x_to_nmstate,
    nm::ip::{nm_ip_setting_to_nmstate4, nm_ip_setting_to_nmstate6},
    nm::lldp::{get_lldp, is_lldp_enabled},
    nm::ovs::{
        get_ovs_dpdk_config, get_ovs_patch_config, nm_ovs_bridge_conf_get,
    },
    nm::user::get_description,
    BaseInterface, BondInterface, DummyInterface, EthernetInterface,
    InfiniBandInterface, Interface, InterfaceState, InterfaceType, Interfaces,
    LinuxBridgeInterface, MacVlanInterface, MacVtapInterface, NetworkState,
    NmstateError, OvsBridgeInterface, OvsInterface, UnknownInterface,
    VlanInterface, VrfInterface, VxlanInterface,
};

pub(crate) fn nm_retrieve(
    running_config_only: bool,
) -> Result<NetworkState, NmstateError> {
    let mut net_state = NetworkState::new();
    net_state.prop_list = vec!["interfaces", "dns"];
    let nm_api = NmApi::new().map_err(nm_error_to_nmstate)?;
    let nm_conns = nm_api
        .applied_connections_get()
        .map_err(nm_error_to_nmstate)?;
    let nm_devs = nm_api.devices_get().map_err(nm_error_to_nmstate)?;

    let nm_saved_conns =
        nm_api.connections_get().map_err(nm_error_to_nmstate)?;
    let nm_acs = nm_api
        .active_connections_get()
        .map_err(nm_error_to_nmstate)?;

    let nm_conns_name_type_index =
        create_index_for_nm_conns_by_name_type(nm_conns.as_slice());
    let mut nm_saved_conn_uuid_index: HashMap<&str, &NmConnection> =
        HashMap::new();
    for nm_saved_conn in nm_saved_conns.as_slice() {
        if let Some(uuid) = nm_saved_conn.uuid() {
            nm_saved_conn_uuid_index.insert(uuid, nm_saved_conn);
        }
    }
    let nm_saved_conns_ctrler_type_index =
        create_index_for_nm_conns_by_ctrler_type(nm_saved_conns.as_slice());
    let nm_acs_name_type_index =
        create_index_for_nm_acs_by_name_type(nm_acs.as_slice());

    // Include disconnected interface as state:down
    // This is used for verify on `state: absent`
    for nm_dev in &nm_devs {
        match nm_dev.state {
            NmDeviceState::Unmanaged | NmDeviceState::Disconnected => {
                if let Some(iface) = nm_dev_to_nm_iface(nm_dev) {
                    debug!(
                        "Found unmanaged or disconnected interface {:?}",
                        iface
                    );
                    net_state.append_interface_data(iface);
                }
            }
            _ => {
                let nm_conn = if let Some(c) = get_first_nm_conn(
                    &nm_conns_name_type_index,
                    &nm_dev.name,
                    &nm_dev.iface_type,
                ) {
                    c
                } else {
                    if nm_dev.state == NmDeviceState::Activated {
                        warn!(
                            "Failed to find applied NmConnection for \
                            interface {} {}",
                            nm_dev.name, nm_dev.iface_type
                        );
                    }
                    continue;
                };
                let nm_ac = get_nm_ac(
                    &nm_acs_name_type_index,
                    &nm_dev.name,
                    &nm_dev.iface_type,
                );

                // NM developer confirmed NmActiveConnection UUID is the
                // UUID of NmConnection associated
                let nm_saved_conn = if let Some(nm_ac) = nm_ac {
                    nm_saved_conn_uuid_index.get(nm_ac.uuid.as_str()).copied()
                } else {
                    None
                };
                let port_saved_nm_conns = if NM_SETTING_CONTROLLERS
                    .contains(&nm_dev.iface_type.as_str())
                {
                    Some(get_port_nm_conns(
                        nm_conn,
                        &nm_saved_conns_ctrler_type_index,
                    ))
                } else {
                    None
                };

                let lldp_neighbors = if is_lldp_enabled(nm_conn) {
                    if running_config_only {
                        Some(Vec::new())
                    } else {
                        Some(
                            nm_api
                                .device_lldp_neighbor_get(&nm_dev.obj_path)
                                .map_err(nm_error_to_nmstate)?,
                        )
                    }
                } else {
                    None
                };
                if let Some(iface) = iface_get(
                    nm_dev,
                    nm_conn,
                    nm_saved_conn,
                    port_saved_nm_conns.as_ref().map(Vec::as_ref),
                    lldp_neighbors,
                ) {
                    debug!("Found interface {:?}", iface);
                    net_state.append_interface_data(iface);
                }
            }
        }
    }

    net_state.dns = retrieve_dns_info(&nm_api, &net_state.interfaces)?;
    if running_config_only {
        net_state.dns.running = None;
    }

    set_ovs_iface_controller_info(&mut net_state.interfaces);

    Ok(net_state)
}

fn nm_dev_iface_type_to_nmstate(nm_dev: &NmDevice) -> InterfaceType {
    match nm_dev.iface_type.as_str() {
        NM_SETTING_WIRED_SETTING_NAME => InterfaceType::Ethernet,
        NM_SETTING_VETH_SETTING_NAME => InterfaceType::Ethernet,
        NM_SETTING_BOND_SETTING_NAME => InterfaceType::Bond,
        NM_SETTING_DUMMY_SETTING_NAME => InterfaceType::Dummy,
        NM_SETTING_BRIDGE_SETTING_NAME => InterfaceType::LinuxBridge,
        NM_SETTING_OVS_BRIDGE_SETTING_NAME => InterfaceType::OvsBridge,
        NM_SETTING_OVS_IFACE_SETTING_NAME => InterfaceType::OvsInterface,
        NM_SETTING_VRF_SETTING_NAME => InterfaceType::Vrf,
        NM_SETTING_VLAN_SETTING_NAME => InterfaceType::Vlan,
        NM_SETTING_VXLAN_SETTING_NAME => InterfaceType::Vxlan,
        NM_SETTING_MACVLAN_SETTING_NAME => {
            if nm_dev.is_mac_vtap {
                InterfaceType::MacVtap
            } else {
                InterfaceType::MacVlan
            }
        }
        NM_SETTING_INFINIBAND_SETTING_NAME => InterfaceType::InfiniBand,
        _ => InterfaceType::Other(nm_dev.iface_type.to_string()),
    }
}

fn nm_conn_to_base_iface(
    nm_dev: &NmDevice,
    nm_conn: &NmConnection,
    nm_saved_conn: Option<&NmConnection>,
    lldp_neighbors: Option<Vec<NmLldpNeighbor>>,
) -> Option<BaseInterface> {
    if let Some(iface_name) = nm_conn.iface_name() {
        let ipv4 = nm_conn.ipv4.as_ref().map(nm_ip_setting_to_nmstate4);
        let ipv6 = nm_conn.ipv6.as_ref().map(nm_ip_setting_to_nmstate6);

        let mut base_iface = BaseInterface::new();
        base_iface.name = iface_name.to_string();
        base_iface.prop_list = vec![
            "name",
            "state",
            "iface_type",
            "ipv4",
            "ipv6",
            "ieee8021x",
            "description",
            "lldp",
        ];
        base_iface.state = InterfaceState::Up;
        base_iface.iface_type = nm_dev_iface_type_to_nmstate(nm_dev);
        base_iface.ipv4 = ipv4;
        base_iface.ipv6 = ipv6;
        base_iface.controller = nm_conn.controller().map(|c| c.to_string());
        base_iface.description = get_description(nm_conn);
        base_iface.lldp =
            Some(lldp_neighbors.map(get_lldp).unwrap_or_default());
        if let Some(nm_saved_conn) = nm_saved_conn {
            // 802.1x password is only available in saved connection
            base_iface.ieee8021x =
                nm_saved_conn.ieee8021x.as_ref().map(nm_802_1x_to_nmstate);
        }
        return Some(base_iface);
    }
    None
}

// Applied connection does not hold OVS config, we need the NmConnection
// used by `NmActiveConnection` also.
fn iface_get(
    nm_dev: &NmDevice,
    nm_conn: &NmConnection,
    nm_saved_conn: Option<&NmConnection>,
    port_saved_nm_conns: Option<&[&NmConnection]>,
    lldp_neighbors: Option<Vec<NmLldpNeighbor>>,
) -> Option<Interface> {
    if let Some(base_iface) =
        nm_conn_to_base_iface(nm_dev, nm_conn, nm_saved_conn, lldp_neighbors)
    {
        let iface = match &base_iface.iface_type {
            InterfaceType::LinuxBridge => Interface::LinuxBridge({
                let mut iface = LinuxBridgeInterface::new();
                iface.base = base_iface;
                iface
            }),
            InterfaceType::Ethernet => Interface::Ethernet({
                let mut iface = EthernetInterface::new();
                iface.base = base_iface;
                iface
            }),
            InterfaceType::Bond => Interface::Bond({
                let mut iface = BondInterface::new();
                iface.base = base_iface;
                iface
            }),
            InterfaceType::OvsInterface => Interface::OvsInterface({
                let mut iface = OvsInterface::new();
                iface.base = base_iface;
                iface.patch = get_ovs_patch_config(nm_conn);
                iface.dpdk = get_ovs_dpdk_config(nm_conn);
                iface
            }),
            InterfaceType::Dummy => Interface::Dummy({
                let mut iface = DummyInterface::new();
                iface.base = base_iface;
                iface
            }),
            InterfaceType::Vlan => Interface::Vlan({
                let mut iface = VlanInterface::new();
                iface.base = base_iface;
                iface
            }),
            InterfaceType::Vxlan => Interface::Vxlan({
                let mut iface = VxlanInterface::new();
                iface.base = base_iface;
                iface
            }),
            InterfaceType::MacVlan => Interface::MacVlan({
                let mut iface = MacVlanInterface::new();
                iface.base = base_iface;
                iface
            }),
            InterfaceType::MacVtap => Interface::MacVtap({
                let mut iface = MacVtapInterface::new();
                iface.base = base_iface;
                iface
            }),
            InterfaceType::Vrf => Interface::Vrf({
                let mut iface = VrfInterface::new();
                iface.base = base_iface;
                iface
            }),
            InterfaceType::OvsBridge => {
                // NetworkManager applied connection does not
                // have ovs configure
                if let Some(nm_saved_conn) = nm_saved_conn {
                    let mut br_iface = OvsBridgeInterface::new();
                    br_iface.base = base_iface;
                    br_iface.bridge = nm_ovs_bridge_conf_get(
                        nm_saved_conn,
                        port_saved_nm_conns,
                    )
                    .ok();
                    Interface::OvsBridge(br_iface)
                } else {
                    warn!(
                        "Failed to get active connection of interface \
                        {} {}",
                        base_iface.name, base_iface.iface_type
                    );
                    return None;
                }
            }
            _ => {
                debug!("Skip unsupported interface {:?}", base_iface);
                return None;
            }
        };
        debug!("Found interface {:?}", iface);
        Some(iface)
    } else {
        // NmConnection has no interface name
        None
    }
}

fn get_first_nm_conn<'a>(
    nm_conns_name_type_index: &'a HashMap<
        (&'a str, &'a str),
        Vec<&'a NmConnection>,
    >,
    name: &'a str,
    nm_iface_type: &'a str,
) -> Option<&'a NmConnection> {
    // Treating veth as ethernet
    let nm_iface_type = if nm_iface_type == NM_SETTING_VETH_SETTING_NAME {
        NM_SETTING_WIRED_SETTING_NAME
    } else {
        nm_iface_type
    };
    if let Some(nm_conns) = nm_conns_name_type_index.get(&(name, nm_iface_type))
    {
        if nm_conns.is_empty() {
            None
        } else {
            Some(nm_conns[0])
        }
    } else {
        None
    }
}

fn get_nm_ac<'a>(
    nm_acs_name_type_index: &'a HashMap<
        (&'a str, &'a str),
        &'a NmActiveConnection,
    >,
    name: &'a str,
    nm_iface_type: &'a str,
) -> Option<&'a NmActiveConnection> {
    nm_acs_name_type_index
        .get(&(
            name,
            match nm_iface_type {
                NM_SETTING_VETH_SETTING_NAME => NM_SETTING_WIRED_SETTING_NAME,
                t => t,
            },
        ))
        .copied()
}

fn set_ovs_iface_controller_info(ifaces: &mut Interfaces) {
    let mut pending_changes: Vec<(&str, &str)> = Vec::new();
    for iface in ifaces.user_ifaces.values() {
        if iface.iface_type() == InterfaceType::OvsBridge {
            if let Some(port_names) = iface.ports() {
                for port_name in port_names {
                    pending_changes.push((port_name, iface.name()));
                }
            }
        }
    }
    for (iface_name, ctrl_name) in pending_changes {
        if let Some(ref mut iface) = ifaces.kernel_ifaces.get_mut(iface_name) {
            iface.base_iface_mut().prop_list.push("controller");
            iface.base_iface_mut().prop_list.push("controller_type");
            iface.base_iface_mut().controller = Some(ctrl_name.to_string());
            iface.base_iface_mut().controller_type =
                Some(InterfaceType::OvsBridge);
        }
    }
}

fn nm_dev_to_nm_iface(nm_dev: &NmDevice) -> Option<Interface> {
    let mut base_iface = BaseInterface::new();
    if nm_dev.name.is_empty() {
        return None;
    } else {
        base_iface.name = nm_dev.name.clone();
    }
    base_iface.prop_list = vec!["name", "iface_type", "state"];
    match nm_dev.state {
        NmDeviceState::Unmanaged => {
            if !nm_dev.real {
                return None;
            } else {
                base_iface.state = InterfaceState::Ignore;
            }
        }
        NmDeviceState::Disconnected => base_iface.state = InterfaceState::Down,
        _ => base_iface.state = InterfaceState::Up,
    }
    base_iface.iface_type = nm_dev_iface_type_to_nmstate(nm_dev);
    Some(match &base_iface.iface_type {
        InterfaceType::Ethernet => Interface::Ethernet({
            let mut iface = EthernetInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::Dummy => Interface::Dummy({
            let mut iface = DummyInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::LinuxBridge => Interface::LinuxBridge({
            let mut iface = LinuxBridgeInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::OvsInterface => Interface::OvsInterface({
            let mut iface = OvsInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::Bond => Interface::Bond({
            let mut iface = BondInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::Vlan => Interface::Vlan({
            let mut iface = VlanInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::Vxlan => Interface::Vxlan({
            let mut iface = VxlanInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::MacVlan => Interface::MacVlan({
            let mut iface = MacVlanInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::MacVtap => Interface::MacVtap({
            let mut iface = MacVtapInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::Vrf => Interface::Vrf({
            let mut iface = VrfInterface::new();
            iface.base = base_iface;
            iface
        }),
        InterfaceType::InfiniBand => Interface::InfiniBand({
            InfiniBandInterface {
                base: base_iface,
                ..Default::default()
            }
        }),
        _ => Interface::Unknown({
            log::info!(
                "Got unsupported interface type {}: {}, ignoring",
                base_iface.iface_type,
                base_iface.name
            );
            let mut iface = UnknownInterface::new();
            if base_iface.name == "lo" {
                base_iface.iface_type = InterfaceType::Loopback;
            }
            base_iface.state = InterfaceState::Ignore;
            iface.base = base_iface;
            iface
        }),
    })
}
