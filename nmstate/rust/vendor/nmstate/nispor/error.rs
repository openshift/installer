use crate::{ErrorKind, NmstateError};

pub(crate) fn np_error_to_nmstate(
    np_error: nispor::NisporError,
) -> NmstateError {
    NmstateError::new(
        ErrorKind::Bug,
        format!("{}: {}", np_error.kind, np_error.msg),
    )
}
