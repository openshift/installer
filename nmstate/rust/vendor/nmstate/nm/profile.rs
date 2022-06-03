use std::collections::{hash_map::Entry, HashMap};
use std::str::FromStr;
use std::time::Instant;

use super::nm_dbus::{NmApi, NmConnection};

use crate::{
    nm::checkpoint::{
        nm_checkpoint_timeout_extend, CHECKPOINT_ROLLBACK_TIMEOUT,
    },
    nm::connection::{
        iface_type_to_nm, NM_SETTING_CONTROLLERS, NM_SETTING_USER_SPACES,
        NM_SETTING_VETH_SETTING_NAME, NM_SETTING_WIRED_SETTING_NAME,
    },
    nm::error::nm_error_to_nmstate,
    nm::ovs::get_ovs_port_name,
    ErrorKind, Interface, InterfaceType, NmstateError,
};

// Found existing profile, prefer the activated one
pub(crate) fn get_exist_profile<'a>(
    exist_nm_conns: &'a [NmConnection],
    iface_name: &str,
    iface_type: &InterfaceType,
    nm_ac_uuids: &[&str],
) -> Option<&'a NmConnection> {
    let mut found_nm_conns: Vec<&NmConnection> = Vec::new();
    for exist_nm_conn in exist_nm_conns {
        let nm_iface_type = if let Ok(t) = iface_type_to_nm(iface_type) {
            // The iface_type will never be veth as top level code
            // `pre_edit_clean()` has confirmed so.
            t
        } else {
            continue;
        };
        if exist_nm_conn.iface_name() == Some(iface_name)
            && (exist_nm_conn.iface_type() == Some(&nm_iface_type)
                || (nm_iface_type == NM_SETTING_WIRED_SETTING_NAME
                    && exist_nm_conn.iface_type()
                        == Some(NM_SETTING_VETH_SETTING_NAME)))
        {
            if let Some(uuid) = exist_nm_conn.uuid() {
                // Prefer activated connection
                if nm_ac_uuids.contains(&uuid) {
                    return Some(exist_nm_conn);
                }
            }
            found_nm_conns.push(exist_nm_conn);
        }
    }
    found_nm_conns.pop()
}

pub(crate) fn delete_exist_profiles(
    nm_api: &NmApi,
    exist_nm_conns: &[NmConnection],
    nm_conns: &[NmConnection],
    checkpoint: &str,
) -> Result<(), NmstateError> {
    let mut now = Instant::now();
    let mut excluded_uuids: Vec<&str> = Vec::new();
    let mut changed_iface_name_types: Vec<(&str, &str)> = Vec::new();
    for nm_conn in nm_conns {
        if let Some(uuid) = nm_conn.uuid() {
            excluded_uuids.push(uuid);
        }
        if let Some(name) = nm_conn.iface_name() {
            if let Some(nm_iface_type) = nm_conn.iface_type() {
                changed_iface_name_types.push((name, nm_iface_type));
            }
        }
    }
    for exist_nm_conn in exist_nm_conns {
        let uuid = if let Some(u) = exist_nm_conn.uuid() {
            u
        } else {
            continue;
        };
        let iface_name = if let Some(i) = exist_nm_conn.iface_name() {
            i
        } else {
            continue;
        };
        let nm_iface_type = if let Some(t) = exist_nm_conn.iface_type() {
            t
        } else {
            continue;
        };
        if !excluded_uuids.contains(&uuid)
            && changed_iface_name_types.contains(&(iface_name, nm_iface_type))
        {
            extend_timeout_if_required(&mut now, checkpoint)?;
            log::info!(
                "Deleting existing connection \
                UUID {:?}, id {:?} type {:?} name {:?}",
                exist_nm_conn.uuid(),
                exist_nm_conn.id(),
                exist_nm_conn.iface_type(),
                exist_nm_conn.iface_name(),
            );
            nm_api
                .connection_delete(uuid)
                .map_err(nm_error_to_nmstate)?;
        }
    }
    Ok(())
}

pub(crate) fn save_nm_profiles(
    nm_api: &NmApi,
    nm_conns: &[NmConnection],
    checkpoint: &str,
    memory_only: bool,
) -> Result<(), NmstateError> {
    let mut now = Instant::now();
    for nm_conn in nm_conns {
        extend_timeout_if_required(&mut now, checkpoint)?;
        log::info!(
            "Creating/Modifying connection \
            UUID {:?}, ID {:?}, type {:?} name {:?}",
            nm_conn.uuid(),
            nm_conn.id(),
            nm_conn.iface_type(),
            nm_conn.iface_name(),
        );
        nm_api
            .connection_add(nm_conn, memory_only)
            .map_err(nm_error_to_nmstate)?;
    }
    Ok(())
}

pub(crate) fn activate_nm_profiles(
    nm_api: &NmApi,
    nm_conns: &[NmConnection],
    nm_ac_uuids: &[&str],
    checkpoint: &str,
) -> Result<(), NmstateError> {
    let mut now = Instant::now();
    let mut new_controllers: Vec<&str> = Vec::new();
    for nm_conn in nm_conns.iter().filter(|c| {
        c.iface_type().map(|t| NM_SETTING_CONTROLLERS.contains(&t))
            == Some(true)
    }) {
        extend_timeout_if_required(&mut now, checkpoint)?;
        if let Some(uuid) = nm_conn.uuid() {
            log::info!(
                "Activating connection {}: {}/{}",
                uuid,
                nm_conn.iface_name().unwrap_or(""),
                nm_conn.iface_type().unwrap_or("")
            );
            if nm_ac_uuids.contains(&uuid) {
                if let Err(e) = nm_api.connection_reapply(nm_conn) {
                    log::info!(
                        "Reapply operation failed trying activation, \
                        reason: {}, retry on normal activation",
                        e
                    );
                    nm_api
                        .connection_activate(uuid)
                        .map_err(nm_error_to_nmstate)?;
                }
            } else {
                new_controllers.push(uuid);
                nm_api
                    .connection_activate(uuid)
                    .map_err(nm_error_to_nmstate)?;
            }
        }
    }
    for nm_conn in nm_conns.iter().filter(|c| {
        c.iface_type().map(|t| NM_SETTING_CONTROLLERS.contains(&t))
            != Some(true)
    }) {
        extend_timeout_if_required(&mut now, checkpoint)?;

        if let Some(uuid) = nm_conn.uuid() {
            if nm_ac_uuids.contains(&uuid) {
                log::info!(
                    "Reapplying connection {}: {}/{}",
                    uuid,
                    nm_conn.iface_name().unwrap_or(""),
                    nm_conn.iface_type().unwrap_or("")
                );
                if let Err(e) = nm_api.connection_reapply(nm_conn) {
                    log::info!(
                        "Reapply operation failed trying activation, \
                        reason: {}, retry on normal activation",
                        e
                    );
                    log::info!(
                        "Activating connection {}: {}/{}",
                        uuid,
                        nm_conn.iface_name().unwrap_or(""),
                        nm_conn.iface_type().unwrap_or("")
                    );
                    nm_api
                        .connection_activate(uuid)
                        .map_err(nm_error_to_nmstate)?;
                }
            } else {
                if let Some(ctrller) = nm_conn.controller() {
                    if nm_conn.iface_type() != Some("ovs-interface") {
                        // OVS port does not do auto port activation.
                        if new_controllers.contains(&ctrller)
                            && nm_conn.controller_type() != Some("ovs-port")
                        {
                            log::info!(
                                "Skip connection activation as its \
                                controller already activated its ports: \
                                {}: {}/{}",
                                uuid,
                                nm_conn.iface_name().unwrap_or(""),
                                nm_conn.iface_type().unwrap_or("")
                            );
                            continue;
                        }
                    }
                }
                log::info!(
                    "Activating connection {}: {}/{}",
                    uuid,
                    nm_conn.iface_name().unwrap_or(""),
                    nm_conn.iface_type().unwrap_or("")
                );
                nm_api
                    .connection_activate(uuid)
                    .map_err(nm_error_to_nmstate)?;
            }
        }
    }
    Ok(())
}

pub(crate) fn deactivate_nm_profiles(
    nm_api: &NmApi,
    nm_conns: &[&NmConnection],
    checkpoint: &str,
) -> Result<(), NmstateError> {
    let mut now = Instant::now();
    for nm_conn in nm_conns {
        extend_timeout_if_required(&mut now, checkpoint)?;
        if let Some(uuid) = nm_conn.uuid() {
            log::info!(
                "Deactivating connection {}: {}/{}",
                uuid,
                nm_conn.iface_name().unwrap_or(""),
                nm_conn.iface_type().unwrap_or("")
            );
            nm_api
                .connection_deactivate(uuid)
                .map_err(nm_error_to_nmstate)?;
        }
    }
    Ok(())
}

pub(crate) fn use_uuid_for_controller_reference(
    nm_conns: &mut [NmConnection],
    des_user_space_ifaces: &HashMap<(String, InterfaceType), Interface>,
    cur_user_space_ifaces: &HashMap<(String, InterfaceType), Interface>,
    exist_nm_conns: &[NmConnection],
) -> Result<(), NmstateError> {
    let mut name_type_2_uuid_index: HashMap<(String, String), String> =
        HashMap::new();

    // This block does not need nm_conn to be mutable, using iter_mut()
    // just to suppress the rust clippy warning message
    for nm_conn in nm_conns.iter_mut() {
        let iface_type = if let Some(i) = nm_conn.iface_type() {
            i
        } else {
            continue;
        };
        if let Some(uuid) = nm_conn.uuid() {
            if let Some(iface_name) = nm_conn.iface_name() {
                name_type_2_uuid_index.insert(
                    (iface_name.to_string(), iface_type.to_string()),
                    uuid.to_string(),
                );
            }
        }
    }

    for nm_conn in exist_nm_conns {
        let iface_type = if let Some(i) = nm_conn.iface_type() {
            i
        } else {
            continue;
        };
        if let Some(uuid) = nm_conn.uuid() {
            if let Some(iface_name) = nm_conn.iface_name() {
                match name_type_2_uuid_index
                    .entry((iface_name.to_string(), iface_type.to_string()))
                {
                    // Prefer newly created NmConnection over existing one
                    Entry::Occupied(_) => {
                        continue;
                    }
                    Entry::Vacant(v) => {
                        v.insert(uuid.to_string());
                    }
                }
            }
        }
    }

    let mut pending_changes: Vec<(&mut NmConnection, String)> = Vec::new();

    for nm_conn in nm_conns.iter_mut() {
        let ctrl_type = if let Some(t) = nm_conn.controller_type() {
            t
        } else {
            continue;
        };
        let mut ctrl_name = if let Some(n) = nm_conn.controller() {
            n.to_string()
        } else {
            continue;
        };

        // Skip if its controller is already a UUID
        if uuid::Uuid::from_str(ctrl_name.as_str()).is_ok() {
            continue;
        }

        if ctrl_type == "ovs-port" {
            if let Some(Interface::OvsBridge(ovs_br_iface)) =
                des_user_space_ifaces
                    .get(&(ctrl_name.to_string(), InterfaceType::OvsBridge))
                    .or_else(|| {
                        cur_user_space_ifaces.get(&(
                            ctrl_name.to_string(),
                            InterfaceType::OvsBridge,
                        ))
                    })
            {
                if let Some(iface_name) = nm_conn.iface_name() {
                    if let Some(ovs_port_name) =
                        get_ovs_port_name(ovs_br_iface, iface_name)
                    {
                        ctrl_name = ovs_port_name.to_string();
                    } else {
                        let e = NmstateError::new(
                            ErrorKind::Bug,
                            format!(
                                "Failed to find OVS port name for \
                                NmConnection {:?}",
                                nm_conn
                            ),
                        );
                        log::error!("{}", e);
                        return Err(e);
                    }
                }
            } else {
                let e = NmstateError::new(
                    ErrorKind::Bug,
                    format!(
                        "Failed to find OVS Bridge interface for \
                        NmConnection {:?}",
                        nm_conn
                    ),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }

        if let Some(uuid) = name_type_2_uuid_index
            .get(&(ctrl_name.clone(), ctrl_type.to_string()))
        {
            pending_changes.push((nm_conn, uuid.to_string()));
        } else {
            let e = NmstateError::new(
                ErrorKind::Bug,
                format!(
                    "BUG: Failed to find UUID of controller connection: \
                {}, {}",
                    ctrl_name, ctrl_type
                ),
            );
            log::error!("{}", e);
            return Err(e);
        }
    }
    for (nm_conn, uuid) in pending_changes {
        if let Some(ref mut nm_conn_set) = &mut nm_conn.connection {
            nm_conn_set.controller = Some(uuid.to_string());
        }
    }
    Ok(())
}

pub(crate) fn extend_timeout_if_required(
    now: &mut Instant,
    checkpoint: &str,
) -> Result<(), NmstateError> {
    // Only extend the timeout when only half of it elapsed
    if now.elapsed().as_secs() >= CHECKPOINT_ROLLBACK_TIMEOUT as u64 / 2 {
        log::debug!("Extending checkpoint timeout");
        nm_checkpoint_timeout_extend(checkpoint, CHECKPOINT_ROLLBACK_TIMEOUT)?;
        *now = Instant::now();
    }
    Ok(())
}

pub(crate) fn use_uuid_for_parent_reference(
    nm_conns: &mut [NmConnection],
    des_kernel_ifaces: &HashMap<String, Interface>,
    exist_nm_conns: &[NmConnection],
) {
    // Pending changes: "child_iface_name: parent_nm_uuid"
    let mut pending_changes: HashMap<String, String> = HashMap::new();

    for iface in des_kernel_ifaces.values() {
        if let Some(parent) = iface.parent() {
            if let Some(parent_uuid) =
                search_uuid_of_kernel_nm_conns(nm_conns, parent).or_else(|| {
                    search_uuid_of_kernel_nm_conns(exist_nm_conns, parent)
                })
            {
                pending_changes
                    .insert(iface.name().to_string(), parent_uuid.to_string());
            }
        }
    }

    for nm_conn in nm_conns {
        if let (Some(iface_name), Some(nm_iface_type)) =
            (nm_conn.iface_name(), nm_conn.iface_type())
        {
            if !NM_SETTING_USER_SPACES.contains(&nm_iface_type) {
                if let Some(parent_uuid) = pending_changes.get(iface_name) {
                    nm_conn.set_parent(parent_uuid);
                }
            }
        }
    }
}

fn search_uuid_of_kernel_nm_conns(
    nm_conns: &[NmConnection],
    iface_name: &str,
) -> Option<String> {
    for nm_conn in nm_conns {
        if let (Some(cur_iface_name), Some(nm_iface_type), Some(uuid)) =
            (nm_conn.iface_name(), nm_conn.iface_type(), nm_conn.uuid())
        {
            if cur_iface_name == iface_name
                && !NM_SETTING_USER_SPACES.contains(&nm_iface_type)
            {
                return Some(uuid.to_string());
            }
        }
    }
    None
}
