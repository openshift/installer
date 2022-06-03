use crate::nm::nm_dbus::NmApi;
use log::warn;

use crate::{nm::error::nm_error_to_nmstate, NmstateError};

// Wait maximum 60 seconds for rollback
pub(crate) const CHECKPOINT_ROLLBACK_TIMEOUT: u32 = 60;

pub(crate) fn nm_checkpoint_create(
    timeout: u32,
) -> Result<String, NmstateError> {
    let nm_api = NmApi::new().map_err(nm_error_to_nmstate)?;
    nm_api
        .checkpoint_create(timeout)
        .map_err(nm_error_to_nmstate)
}

pub(crate) fn nm_checkpoint_rollback(
    checkpoint: &str,
) -> Result<(), NmstateError> {
    let nm_api = NmApi::new().map_err(nm_error_to_nmstate)?;
    nm_api
        .checkpoint_rollback(checkpoint)
        .map_err(nm_error_to_nmstate)?;
    if let Err(e) = nm_api.wait_checkpoint_rollback(CHECKPOINT_ROLLBACK_TIMEOUT)
    {
        warn!("{}", e);
    }
    Ok(())
}

pub(crate) fn nm_checkpoint_destroy(
    checkpoint: &str,
) -> Result<(), NmstateError> {
    let nm_api = NmApi::new().map_err(nm_error_to_nmstate)?;
    nm_api
        .checkpoint_destroy(checkpoint)
        .map_err(nm_error_to_nmstate)
}

pub(crate) fn nm_checkpoint_timeout_extend(
    checkpoint: &str,
    added_time_sec: u32,
) -> Result<(), NmstateError> {
    let nm_api = NmApi::new().map_err(nm_error_to_nmstate)?;
    nm_api
        .checkpoint_timeout_extend(checkpoint, added_time_sec)
        .map_err(nm_error_to_nmstate)
}
