use std::collections::HashSet;
use std::time::Instant;

use crate::nm::nm_dbus::{NmApi, NmConnection};
use log::info;

use crate::{
    nm::connection::{
        create_index_for_nm_conns_by_name_type, iface_to_nm_connections,
        iface_type_to_nm, NM_SETTING_OVS_PORT_SETTING_NAME,
    },
    nm::device::create_index_for_nm_devs,
    nm::error::nm_error_to_nmstate,
    nm::profile::{
        activate_nm_profiles, deactivate_nm_profiles, delete_exist_profiles,
        extend_timeout_if_required, get_exist_profile, save_nm_profiles,
        use_uuid_for_controller_reference, use_uuid_for_parent_reference,
    },
    nm::route::is_route_removed,
    nm::veth::{is_veth_peer_changed, is_veth_peer_in_desire},
    nm::vlan::is_vlan_id_changed,
    nm::vrf::is_vrf_table_id_changed,
    nm::vxlan::is_vxlan_id_changed,
    Interface, InterfaceType, NetworkState, NmstateError, OvsBridgeInterface,
    RouteEntry,
};

const ACTIVATION_RETRY_COUNT: usize = 5;
const ACTIVATION_RETRY_INTERVAL: u64 = 1;

pub(crate) fn nm_apply(
    add_net_state: &NetworkState,
    chg_net_state: &NetworkState,
    del_net_state: &NetworkState,
    cur_net_state: &NetworkState,
    des_net_state: &NetworkState,
    checkpoint: &str,
    memory_only: bool,
) -> Result<(), NmstateError> {
    let nm_api = NmApi::new().map_err(nm_error_to_nmstate)?;

    if !memory_only {
        delete_net_state(&nm_api, del_net_state, checkpoint)?;
    }
    apply_single_state(
        &nm_api,
        add_net_state,
        cur_net_state,
        des_net_state,
        checkpoint,
        memory_only,
    )?;
    apply_single_state(
        &nm_api,
        chg_net_state,
        cur_net_state,
        des_net_state,
        checkpoint,
        memory_only,
    )?;

    Ok(())
}

fn delete_net_state(
    nm_api: &NmApi,
    net_state: &NetworkState,
    checkpoint: &str,
) -> Result<(), NmstateError> {
    let all_nm_conns = nm_api.connections_get().map_err(nm_error_to_nmstate)?;

    let nm_conns_name_type_index =
        create_index_for_nm_conns_by_name_type(&all_nm_conns);
    let mut uuids_to_delete: HashSet<&str> = HashSet::new();

    for iface in &(net_state.interfaces.to_vec()) {
        if !iface.is_absent() {
            continue;
        }
        // If interface type not mentioned, we delete all profile with interface
        // name
        let nm_conns_to_delete: Option<Vec<&NmConnection>> =
            if iface.iface_type() == InterfaceType::Unknown {
                Some(
                    all_nm_conns
                        .as_slice()
                        .iter()
                        .filter(|c| c.iface_name() == Some(iface.name()))
                        .collect(),
                )
            } else {
                let nm_iface_type = iface_type_to_nm(&iface.iface_type())?;
                nm_conns_name_type_index
                    .get(&(iface.name(), &nm_iface_type))
                    .cloned()
            };
        // Delete all existing connections for this interface
        if let Some(nm_conns) = nm_conns_to_delete {
            for nm_conn in nm_conns {
                if let Some(uuid) = nm_conn.uuid() {
                    info!(
                        "Deleting NM connection for absent interface \
                        {}/{}: {}",
                        &iface.name(),
                        &iface.iface_type(),
                        uuid
                    );
                    uuids_to_delete.insert(uuid);
                }
                // Delete OVS port profile along with OVS system and internal
                // Interface
                if nm_conn.controller_type() == Some("ovs-port") {
                    // TODO: handle pre-exist OVS config using name instead of
                    // UUID for controller
                    if let Some(uuid) = nm_conn.controller() {
                        info!(
                            "Deleting NM OVS port connection {} \
                             for absent OVS interface {}",
                            uuid,
                            &iface.name(),
                        );
                        uuids_to_delete.insert(uuid);
                    }
                }
            }
        }
    }

    let mut now = Instant::now();
    for uuid in &uuids_to_delete {
        extend_timeout_if_required(&mut now, checkpoint)?;
        nm_api
            .connection_delete(uuid)
            .map_err(nm_error_to_nmstate)?;
    }

    delete_orphan_ports(nm_api, &uuids_to_delete, checkpoint)?;
    delete_remain_virtual_interface_as_desired(nm_api, net_state, checkpoint)?;
    Ok(())
}

fn apply_single_state(
    nm_api: &NmApi,
    net_state: &NetworkState,
    cur_net_state: &NetworkState,
    des_net_state: &NetworkState,
    checkpoint: &str,
    memory_only: bool,
) -> Result<(), NmstateError> {
    if let Some(hostname) =
        net_state.hostname.as_ref().and_then(|c| c.config.as_ref())
    {
        if memory_only {
            log::warn!(
                "Cannot change configure hostname in memory only mode, \
                ignoring"
            );
        } else {
            nm_api.hostname_set(hostname).map_err(nm_error_to_nmstate)?;
        }
    }

    if net_state.interfaces.to_vec().is_empty() {
        return Ok(());
    }
    let mut nm_conns_to_activate: Vec<NmConnection> = Vec::new();

    let exist_nm_conns =
        nm_api.connections_get().map_err(nm_error_to_nmstate)?;
    let nm_acs = nm_api
        .active_connections_get()
        .map_err(nm_error_to_nmstate)?;
    let nm_ac_uuids: Vec<&str> =
        nm_acs.iter().map(|nm_ac| &nm_ac.uuid as &str).collect();
    let activated_nm_conns: Vec<&NmConnection> = exist_nm_conns
        .iter()
        .filter(|c| {
            if let Some(uuid) = c.uuid() {
                nm_ac_uuids.contains(&uuid)
            } else {
                false
            }
        })
        .collect();

    let ifaces = net_state.interfaces.to_vec();

    for iface in ifaces.iter() {
        if iface.iface_type() != InterfaceType::Unknown && iface.is_up() {
            let mut ctrl_iface: Option<&Interface> = None;
            if let Some(ctrl_iface_name) = &iface.base_iface().controller {
                if let Some(ctrl_type) = &iface.base_iface().controller_type {
                    ctrl_iface = des_net_state
                        .interfaces
                        .get_iface(ctrl_iface_name, ctrl_type.clone());
                }
            }
            let mut routes: Vec<&RouteEntry> = Vec::new();
            if let Some(config_routes) = net_state.routes.config.as_ref() {
                for route in config_routes {
                    if let Some(i) = route.next_hop_iface.as_ref() {
                        if i == iface.name() {
                            routes.push(route);
                        }
                    }
                }
            }
            for nm_conn in iface_to_nm_connections(
                iface,
                ctrl_iface,
                &exist_nm_conns,
                &nm_ac_uuids,
                is_veth_peer_in_desire(iface, ifaces.as_slice()),
                cur_net_state,
            )? {
                nm_conns_to_activate.push(nm_conn);
            }
        }
    }
    let nm_conns_to_deactivate = ifaces
        .into_iter()
        .filter(|iface| iface.is_down())
        .filter_map(|iface| {
            get_exist_profile(
                &exist_nm_conns,
                &iface.base_iface().name,
                &iface.base_iface().iface_type,
                &nm_ac_uuids,
            )
        })
        .collect::<Vec<_>>();

    let mut ovs_br_ifaces: Vec<&OvsBridgeInterface> = Vec::new();
    for iface in net_state.interfaces.user_ifaces.values() {
        if let Interface::OvsBridge(ref br_iface) = iface {
            ovs_br_ifaces.push(br_iface);
        }
    }

    use_uuid_for_controller_reference(
        &mut nm_conns_to_activate,
        &des_net_state.interfaces.user_ifaces,
        &cur_net_state.interfaces.user_ifaces,
        &exist_nm_conns,
    )?;
    use_uuid_for_parent_reference(
        &mut nm_conns_to_activate,
        &des_net_state.interfaces.kernel_ifaces,
        &exist_nm_conns,
    );

    let nm_conns_to_deactivate_first = gen_nm_conn_need_to_deactivate_first(
        nm_conns_to_activate.as_slice(),
        activated_nm_conns.as_slice(),
    );
    deactivate_nm_profiles(
        nm_api,
        nm_conns_to_deactivate_first.as_slice(),
        checkpoint,
    )?;
    save_nm_profiles(
        nm_api,
        nm_conns_to_activate.as_slice(),
        checkpoint,
        memory_only,
    )?;
    if !memory_only {
        delete_exist_profiles(
            nm_api,
            &exist_nm_conns,
            &nm_conns_to_activate,
            checkpoint,
        )?;
    }

    for i in 0..ACTIVATION_RETRY_COUNT {
        match activate_nm_profiles(
            nm_api,
            nm_conns_to_activate.as_slice(),
            nm_ac_uuids.as_slice(),
            checkpoint,
        ) {
            Ok(()) => break,
            Err(e) => {
                if i == ACTIVATION_RETRY_COUNT - 1 {
                    return Err(e);
                } else {
                    log::warn!("Activation failure: {}, retrying", e);
                    std::thread::sleep(std::time::Duration::from_secs(
                        ACTIVATION_RETRY_INTERVAL,
                    ));
                }
            }
        }
    }

    deactivate_nm_profiles(
        nm_api,
        nm_conns_to_deactivate.as_slice(),
        checkpoint,
    )?;
    Ok(())
}

fn delete_remain_virtual_interface_as_desired(
    nm_api: &NmApi,
    net_state: &NetworkState,
    checkpoint: &str,
) -> Result<(), NmstateError> {
    let nm_devs = nm_api.devices_get().map_err(nm_error_to_nmstate)?;
    let nm_devs_indexed = create_index_for_nm_devs(&nm_devs);
    let mut now = Instant::now();
    // Interfaces created by non-NM tools will not be deleted by connection
    // deletion, remove manually.
    for iface in &(net_state.interfaces.to_vec()) {
        if !iface.is_absent() {
            continue;
        }
        if iface.is_virtual() {
            if let Some(nm_dev) = nm_devs_indexed.get(&(
                iface.name().to_string(),
                iface.iface_type().to_string(),
            )) {
                info!(
                    "Deleting interface {}/{}: {}",
                    &iface.name(),
                    &iface.iface_type(),
                    &nm_dev.obj_path
                );
                // There might be an race with on-going profile/connection
                // deletion, verification will raise error for it later.
                extend_timeout_if_required(&mut now, checkpoint)?;
                if let Err(e) = nm_api.device_delete(&nm_dev.obj_path) {
                    log::debug!("Failed to delete interface {:?}", e);
                }
            }
        }
    }
    Ok(())
}

// If any connection still referring to deleted UUID, we should delete it also
fn delete_orphan_ports(
    nm_api: &NmApi,
    uuids_deleted: &HashSet<&str>,
    checkpoint: &str,
) -> Result<(), NmstateError> {
    let mut uuids_to_delete = Vec::new();
    let all_nm_conns = nm_api.connections_get().map_err(nm_error_to_nmstate)?;
    for nm_conn in &all_nm_conns {
        if nm_conn.iface_type() != Some(NM_SETTING_OVS_PORT_SETTING_NAME) {
            continue;
        }
        if let Some(ctrl_uuid) = nm_conn.controller() {
            if uuids_deleted.contains(ctrl_uuid) {
                if let Some(uuid) = nm_conn.uuid() {
                    info!(
                        "Deleting NM orphan profile {}/{}: {}",
                        nm_conn.iface_name().unwrap_or(""),
                        nm_conn.iface_type().unwrap_or(""),
                        uuid
                    );
                    uuids_to_delete.push(uuid);
                }
            }
        }
    }
    let mut now = Instant::now();
    for uuid in &uuids_to_delete {
        extend_timeout_if_required(&mut now, checkpoint)?;
        nm_api
            .connection_delete(uuid)
            .map_err(nm_error_to_nmstate)?;
    }
    Ok(())
}

// * NM has problem on remove routes, we need to deactivate it first
//  https://bugzilla.redhat.com/1837254
// * NM cannot change VRF table ID, so we deactivate first
fn gen_nm_conn_need_to_deactivate_first<'a>(
    nm_conns_to_activate: &[NmConnection],
    activated_nm_conns: &[&'a NmConnection],
) -> Vec<&'a NmConnection> {
    let mut ret: Vec<&NmConnection> = Vec::new();
    for nm_conn in nm_conns_to_activate {
        if let Some(uuid) = nm_conn.uuid() {
            if let Some(activated_nm_con) =
                activated_nm_conns.iter().find(|c| {
                    if let Some(cur_uuid) = c.uuid() {
                        cur_uuid == uuid
                    } else {
                        false
                    }
                })
            {
                if is_route_removed(nm_conn, activated_nm_con)
                    || is_vrf_table_id_changed(nm_conn, activated_nm_con)
                    || is_vlan_id_changed(nm_conn, activated_nm_con)
                    || is_vxlan_id_changed(nm_conn, activated_nm_con)
                    || is_veth_peer_changed(nm_conn, activated_nm_con)
                {
                    ret.push(activated_nm_con);
                }
            }
        }
    }
    ret
}
