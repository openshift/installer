use std::collections::HashMap;

use log::warn;

use crate::{
    nispor::{
        ip::{nmstate_ipv4_to_np, nmstate_ipv6_to_np},
        veth::nms_veth_conf_to_np,
        vlan::nms_vlan_conf_to_np,
    },
    ErrorKind, Interface, InterfaceType, Interfaces, NetworkState,
    NmstateError,
};

pub(crate) fn nispor_apply(
    add_net_state: &NetworkState,
    chg_net_state: &NetworkState,
    del_net_state: &NetworkState,
    cur_net_state: &NetworkState,
) -> Result<(), NmstateError> {
    let mut del_net_state = del_net_state.clone();
    only_delete_one_end_of_veth_peer(
        &mut del_net_state.interfaces,
        &cur_net_state.interfaces,
    );
    apply_single_state(&del_net_state)?;
    apply_single_state(add_net_state)?;
    apply_single_state(chg_net_state)?;
    Ok(())
}

fn net_state_to_nispor(
    net_state: &NetworkState,
) -> Result<nispor::NetConf, NmstateError> {
    let mut np_ifaces: Vec<nispor::IfaceConf> = Vec::new();

    for iface in net_state.interfaces.to_vec() {
        if iface.is_up() {
            let np_iface_type = nmstate_iface_type_to_np(&iface.iface_type());
            if np_iface_type == nispor::IfaceType::Unknown {
                warn!(
                    "Unknown interface type {} for interface {}",
                    iface.iface_type(),
                    iface.name()
                );
                continue;
            }
            np_ifaces.push(nmstate_iface_to_np(iface, np_iface_type)?);
        } else if iface.is_absent() {
            println!("del {:?} {:?}", iface.name(), iface.iface_type());
            let mut iface_conf = nispor::IfaceConf::default();
            iface_conf.name = iface.name().to_string();
            iface_conf.iface_type =
                Some(nmstate_iface_type_to_np(&iface.iface_type()));
            iface_conf.state = nispor::IfaceState::Absent;
            np_ifaces.push(iface_conf);
        }
    }

    let mut net_conf = nispor::NetConf::default();
    net_conf.ifaces = Some(np_ifaces);

    Ok(net_conf)
}

fn nmstate_iface_type_to_np(
    nms_iface_type: &InterfaceType,
) -> nispor::IfaceType {
    match nms_iface_type {
        InterfaceType::LinuxBridge => nispor::IfaceType::Bridge,
        InterfaceType::Bond => nispor::IfaceType::Bond,
        InterfaceType::Ethernet => nispor::IfaceType::Ethernet,
        InterfaceType::Veth => nispor::IfaceType::Veth,
        InterfaceType::Vlan => nispor::IfaceType::Vlan,
        _ => nispor::IfaceType::Unknown,
    }
}

fn nmstate_iface_to_np(
    nms_iface: &Interface,
    np_iface_type: nispor::IfaceType,
) -> Result<nispor::IfaceConf, NmstateError> {
    let mut np_iface = nispor::IfaceConf::default();
    np_iface.name = nms_iface.name().to_string();
    np_iface.iface_type = Some(np_iface_type);
    np_iface.state = nispor::IfaceState::Up;
    let base_iface = &nms_iface.base_iface();
    if let Some(ctrl_name) = &base_iface.controller {
        np_iface.controller = Some(ctrl_name.to_string())
    }
    if base_iface.can_have_ip() {
        np_iface.ipv4 = Some(nmstate_ipv4_to_np(base_iface.ipv4.as_ref()));
        np_iface.ipv6 = Some(nmstate_ipv6_to_np(base_iface.ipv6.as_ref()));
    }

    np_iface.mac_address = base_iface.mac_address.clone();

    if let Interface::Ethernet(eth_iface) = nms_iface {
        np_iface.veth = nms_veth_conf_to_np(eth_iface.veth.as_ref());
    } else if let Interface::Vlan(vlan_iface) = nms_iface {
        np_iface.vlan = nms_vlan_conf_to_np(vlan_iface.vlan.as_ref());
    }

    Ok(np_iface)
}

fn apply_single_state(net_state: &NetworkState) -> Result<(), NmstateError> {
    let np_net_conf = net_state_to_nispor(net_state)?;
    if let Err(e) = np_net_conf.apply() {
        return Err(NmstateError::new(
            ErrorKind::PluginFailure,
            format!("Unknown error from nipsor plugin: {}, {}", e.kind, e.msg),
        ));
    } else {
        Ok(())
    }
}

// Deleting one end of veth peer is enough, remove other end from desire state
// TODO: Fix nispor to ignore ENODEV(19) error when deleting interface which is
// already gone.
fn only_delete_one_end_of_veth_peer(
    desired: &mut Interfaces,
    current: &Interfaces,
) {
    let mut veth_pairs: HashMap<String, String> = HashMap::new();
    for iface in current.kernel_ifaces.values() {
        if let Interface::Ethernet(eth_iface) = iface {
            if let Some(peer_name) = eth_iface
                .veth
                .as_ref()
                .map(|veth_conf| veth_conf.peer.as_str())
            {
                veth_pairs
                    .insert(iface.name().to_string(), peer_name.to_string());
            }
        }
    }

    let mut veth_to_ignore = Vec::new();
    for iface in desired.kernel_ifaces.values().filter(|i| i.is_absent()) {
        if let Some(veth_peer) = veth_pairs.get(iface.name()) {
            if let Some(veth_peer_iface) = desired.kernel_ifaces.get(veth_peer)
            {
                if iface.name() < veth_peer.as_str()
                    && veth_peer_iface.is_absent()
                {
                    veth_to_ignore.push(veth_peer);
                }
            }
        }
    }

    for iface_name in veth_to_ignore {
        desired.kernel_ifaces.remove(iface_name);
    }
}
