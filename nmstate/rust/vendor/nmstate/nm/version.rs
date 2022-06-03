use crate::nm::nm_dbus::NmApi;
use crate::{nm::error::nm_error_to_nmstate, NmstateError};

pub(crate) fn nm_version() -> Result<String, NmstateError> {
    let nm_api = NmApi::new().map_err(nm_error_to_nmstate)?;
    nm_api.version().map_err(nm_error_to_nmstate)
}

// This helper function will help us to avoid introducing new dependencies to
// the project.
pub(crate) fn nm_supports_accept_all_mac_addresses_mode(
) -> Result<bool, NmstateError> {
    let version = nm_version()?;
    let version_split = version.split('.');
    let supported_version = Vec::<u32>::from([1, 32]);
    let mut supported_elem = supported_version.iter();

    for v_elem in version_split {
        if v_elem.chars().all(char::is_numeric) {
            if let Some(supported_v) = supported_elem.next() {
                if v_elem.parse::<u32>().unwrap_or_default() < *supported_v {
                    return Ok(false);
                }
            } else {
                return Ok(true);
            }
        }
    }

    Ok(true)
}
