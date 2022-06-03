// Copyright 2021 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

use std::convert::TryFrom;
use std::time::{Duration, Instant};

use log::debug;

use super::{
    active_connection::{
        get_nm_ac_by_obj_path, nm_ac_obj_path_uuid_get, NmActiveConnection,
    },
    connection::{nm_con_get_from_obj_path, NmConnection},
    dbus::NmDbus,
    device::{
        nm_dev_delete, nm_dev_from_obj_path, nm_dev_get_llpd, NmDevice,
        NmDeviceState, NmDeviceStateReason,
    },
    dns::NmDnsEntry,
    error::{ErrorKind, NmError},
    lldp::NmLldpNeighbor,
};

pub struct NmApi<'a> {
    pub(crate) dbus: NmDbus<'a>,
}

const RETRY_INTERVAL_MILLISECOND: u64 = 500;
const RETRY_COUNT: usize = 60;

impl<'a> NmApi<'a> {
    pub fn new() -> Result<Self, NmError> {
        Ok(Self {
            dbus: NmDbus::new()?,
        })
    }

    pub fn version(&self) -> Result<String, NmError> {
        self.dbus.version()
    }

    pub fn checkpoint_create(&self, timeout: u32) -> Result<String, NmError> {
        debug!("checkpoint_create");
        let cp = self.dbus.checkpoint_create(timeout)?;
        debug!("checkpoint created: {}", &cp);
        Ok(cp)
    }

    pub fn checkpoint_destroy(&self, checkpoint: &str) -> Result<(), NmError> {
        let mut checkpoint_to_destroy: String = checkpoint.to_string();
        if checkpoint_to_destroy.is_empty() {
            checkpoint_to_destroy = self.last_active_checkpoint()?
        }
        debug!("checkpoint_destroy: {}", checkpoint_to_destroy);
        self.dbus.checkpoint_destroy(checkpoint_to_destroy.as_str())
    }

    pub fn checkpoint_rollback(&self, checkpoint: &str) -> Result<(), NmError> {
        let mut checkpoint_to_rollback: String = checkpoint.to_string();
        if checkpoint_to_rollback.is_empty() {
            checkpoint_to_rollback = self.last_active_checkpoint()?
        }
        debug!("checkpoint_rollback: {}", checkpoint_to_rollback);
        self.dbus
            .checkpoint_rollback(checkpoint_to_rollback.as_str())
    }

    fn last_active_checkpoint(&self) -> Result<String, NmError> {
        debug!("last_active_checkpoint");
        let mut checkpoints = self.dbus.checkpoints()?;
        if !checkpoints.is_empty() {
            Ok(checkpoints.remove(0))
        } else {
            Err(NmError::new(
                ErrorKind::NotFound,
                "Not active checkpoints".to_string(),
            ))
        }
    }

    pub fn connection_activate(&self, uuid: &str) -> Result<(), NmError> {
        debug!("connection_activate: {}", uuid);
        // Race: Connection might just created
        with_retry(RETRY_INTERVAL_MILLISECOND, RETRY_COUNT, || {
            let nm_conn = self.dbus.get_connection_by_uuid(uuid)?;
            self.dbus.connection_activate(&nm_conn)
        })
    }

    pub fn connection_deactivate(&self, uuid: &str) -> Result<(), NmError> {
        debug!("connection_deactivate: {}", uuid);
        // Race: ActiveConnection might be deleted
        with_retry(RETRY_INTERVAL_MILLISECOND, RETRY_COUNT, || {
            if let Ok(nm_ac) = get_nm_ac_obj_path_by_uuid(&self.dbus, uuid) {
                if !nm_ac.is_empty() {
                    self.dbus.connection_deactivate(&nm_ac)?;
                }
            }
            Ok(())
        })
    }

    pub fn connections_get(&self) -> Result<Vec<NmConnection>, NmError> {
        debug!("connections_get");
        let mut nm_conns = Vec::new();
        for nm_conn_obj_path in self.dbus.nm_conn_obj_paths_get()? {
            // Race: Connection might just been deleted, hence we ignore error
            // here
            if let Ok(c) = nm_con_get_from_obj_path(
                &self.dbus.connection,
                &nm_conn_obj_path,
            ) {
                debug!("Got connection {:?}", c);
                nm_conns.push(c);
            }
        }
        Ok(nm_conns)
    }

    pub fn applied_connections_get(
        &self,
    ) -> Result<Vec<NmConnection>, NmError> {
        debug!("applied_connections_get");
        let nm_devs = self.dbus.nm_dev_obj_paths_get()?;
        let nm_conns = nm_devs
            .iter()
            .map(|nm_dev| self.dbus.nm_dev_applied_connection_get(nm_dev))
            // Ignore dbus errors
            .filter(|result| {
                !matches!(result, Err(NmError {
                    dbus_error, ..
                    }) if dbus_error.is_some())
            })
            .collect::<Result<Vec<_>, _>>()?;
        nm_conns
            .iter()
            .for_each(|conn| debug!("Get Applied connection {:?}", conn));
        Ok(nm_conns)
    }

    pub fn connection_add(
        &self,
        nm_conn: &NmConnection,
        memory_only: bool,
    ) -> Result<(), NmError> {
        debug!("connection_add: {:?}", nm_conn);
        if let Some(uuid) = nm_conn.uuid() {
            if let Ok(con_obj_path) = self.dbus.get_connection_by_uuid(uuid) {
                return self.dbus.connection_update(
                    &con_obj_path,
                    nm_conn,
                    memory_only,
                );
            }
        };
        self.dbus.connection_add(nm_conn, memory_only)?;
        Ok(())
    }

    pub fn connection_delete(&self, uuid: &str) -> Result<(), NmError> {
        debug!("connection_delete: {}", uuid);
        with_retry(RETRY_INTERVAL_MILLISECOND, RETRY_COUNT, || {
            if let Ok(con_obj_path) = self.dbus.get_connection_by_uuid(uuid) {
                debug!(
                    "Found nm_connection {} for UUID {}",
                    con_obj_path, uuid
                );
                if !con_obj_path.is_empty() {
                    self.dbus.connection_delete(&con_obj_path)?;
                }
            }
            Ok(())
        })
    }

    pub fn connection_reapply(
        &self,
        nm_conn: &NmConnection,
    ) -> Result<(), NmError> {
        debug!("connection_reapply: {:?}", nm_conn);
        if let Some(iface_name) = nm_conn.iface_name() {
            let nm_dev_obj_path = self.dbus.nm_dev_obj_path_get(iface_name)?;
            self.dbus.nm_dev_reapply(&nm_dev_obj_path, nm_conn)
        } else {
            Err(NmError::new(
                ErrorKind::InvalidArgument,
                format!(
                    "Failed to extract interface name from connection {:?}",
                    nm_conn
                ),
            ))
        }
    }

    pub fn uuid_gen() -> String {
        // Use Linux random number generator (RNG) to generate UUID
        uuid::Uuid::new_v4().to_hyphenated().to_string()
    }

    pub fn active_connections_get(
        &self,
    ) -> Result<Vec<NmActiveConnection>, NmError> {
        debug!("active_connections_get");
        let mut nm_acs = Vec::new();
        let nm_ac_obj_paths = self.dbus.active_connections()?;
        for nm_ac_obj_path in nm_ac_obj_paths {
            // Race condition: Active connection might just been deleted,
            // we ignore error here
            if let Ok(Some(nm_ac)) =
                get_nm_ac_by_obj_path(&self.dbus.connection, &nm_ac_obj_path)
            {
                debug!("Got active connection {:?}", nm_ac);
                nm_acs.push(nm_ac);
            }
        }
        Ok(nm_acs)
    }

    pub fn checkpoint_timeout_extend(
        &self,
        checkpoint: &str,
        added_time_sec: u32,
    ) -> Result<(), NmError> {
        debug!(
            "checkpoint_timeout_extend: {} {}",
            checkpoint, added_time_sec
        );
        self.dbus
            .checkpoint_timeout_extend(checkpoint, added_time_sec)
    }

    pub fn devices_get(&self) -> Result<Vec<NmDevice>, NmError> {
        debug!("devices_get");
        let mut ret = Vec::new();
        for nm_dev_obj_path in &self.dbus.nm_dev_obj_paths_get()? {
            match nm_dev_from_obj_path(&self.dbus.connection, nm_dev_obj_path) {
                Ok(nm_dev) => {
                    debug!("Got Device {:?}", nm_dev);
                    ret.push(nm_dev);
                }
                Err(e) => {
                    // We might have race when relieve device list along with
                    // deleting device
                    debug!(
                        "Failed to retrieve device {} {}",
                        nm_dev_obj_path, e
                    )
                }
            }
        }
        Ok(ret)
    }

    pub fn device_delete(&self, nm_dev_obj_path: &str) -> Result<(), NmError> {
        nm_dev_delete(&self.dbus.connection, nm_dev_obj_path)
    }

    pub fn device_lldp_neighbor_get(
        &self,
        nm_dev_obj_path: &str,
    ) -> Result<Vec<NmLldpNeighbor>, NmError> {
        nm_dev_get_llpd(&self.dbus.connection, nm_dev_obj_path)
    }

    // If any device is with NewActivation or IpConfig state,
    // we wait its activation.
    pub fn wait_checkpoint_rollback(
        &self,
        // TODO: return error when waiting_nm_dev is not changing for given
        // time.
        timeout: u32,
    ) -> Result<(), NmError> {
        debug!("wait_checkpoint_rollback");
        let start = Instant::now();
        while start.elapsed() <= Duration::from_secs(timeout.into()) {
            let mut waiting_nm_dev: Vec<&NmDevice> = Vec::new();
            let nm_devs = self.devices_get()?;
            for nm_dev in &nm_devs {
                if nm_dev.state_reason == NmDeviceStateReason::NewActivation
                    || nm_dev.state == NmDeviceState::IpConfig
                    || nm_dev.state == NmDeviceState::Deactivating
                {
                    waiting_nm_dev.push(nm_dev);
                }
            }
            if waiting_nm_dev.is_empty() {
                return Ok(());
            } else {
                debug!(
                    "Waiting rollback on these devices {:?}",
                    waiting_nm_dev
                );
                std::thread::sleep(Duration::from_millis(500));
            }
        }
        Err(NmError::new(
            ErrorKind::Timeout,
            "Timeout on waiting rollback".to_string(),
        ))
    }

    pub fn get_dns_configuration(&self) -> Result<Vec<NmDnsEntry>, NmError> {
        let mut ret: Vec<NmDnsEntry> = Vec::new();
        for dns_value in self.dbus.get_dns_configuration()? {
            ret.push(NmDnsEntry::try_from(dns_value)?);
        }
        Ok(ret)
    }

    pub fn hostname_set(&self, hostname: &str) -> Result<(), NmError> {
        if hostname.is_empty() {
            // Due to bug https://bugzilla.redhat.com/2090946
            // NetworkManager daemon cannot remove static hostname, hence we
            // just delete the /etc/hostname file
            if std::path::Path::new("/etc/hostname").exists() {
                if let Err(e) = std::fs::remove_file("/etc/hostname") {
                    log::error!("Failed to remove static /etc/hostname: {}", e);
                }
            }
            Ok(())
        } else {
            self.dbus.hostname_set(hostname)
        }
    }
}

fn get_nm_ac_obj_path_by_uuid(
    dbus: &NmDbus,
    uuid: &str,
) -> Result<String, NmError> {
    let nm_ac_obj_paths = dbus.active_connections()?;

    for nm_ac_obj_path in nm_ac_obj_paths {
        if nm_ac_obj_path_uuid_get(&dbus.connection, &nm_ac_obj_path)? == uuid {
            return Ok(nm_ac_obj_path);
        }
    }
    Ok("".into())
}

fn with_retry<T>(interval_ms: u64, count: usize, func: T) -> Result<(), NmError>
where
    T: FnOnce() -> Result<(), NmError> + Copy,
{
    let mut cur_count = 0usize;
    while cur_count < count {
        if let Err(e) = func() {
            if cur_count == count - 1 {
                return Err(e);
            } else {
                if let Some(zbus::Error::MethodError(..)) = e.dbus_error {
                    log::error!("{}", e);
                    return Err(e);
                }
                log::info!("Retrying on NM dbus failure: {}", e);
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
