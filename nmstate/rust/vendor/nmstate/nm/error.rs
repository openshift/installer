use crate::nm::nm_dbus::NmError;

use crate::{ErrorKind, NmstateError};

pub(crate) fn nm_error_to_nmstate(nm_error: NmError) -> NmstateError {
    if nm_error
        .to_string()
        .contains("NetworkManager plugin for 'ovs-bridge' unavailable")
    {
        NmstateError::new(
            ErrorKind::DependencyError,
            "NetworkManager does not have OVS plugin installed for \
            OVS modification"
                .to_string(),
        )
    } else {
        NmstateError::new(
            ErrorKind::Bug,
            format!(
                "{}: {} dbus: {:?}",
                nm_error.kind, nm_error.msg, nm_error.dbus_error
            ),
        )
    }
}
