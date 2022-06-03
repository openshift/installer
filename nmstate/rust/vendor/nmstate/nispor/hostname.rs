use std::io::Read;

use crate::{ErrorKind, HostNameState, NmstateError};

const HOST_NAME_MAX: usize = 64;

pub(crate) fn get_hostname_state() -> Option<HostNameState> {
    let mut buffer = [0u8; HOST_NAME_MAX];
    let running = match nix::unistd::gethostname(&mut buffer) {
        Ok(hostname_cstr) => match hostname_cstr.to_str() {
            Ok(h) => Some(h.to_string()),
            Err(e) => {
                log::error!("Failed to convert hostname to String: {}", e);
                None
            }
        },
        Err(e) => {
            log::error!("Failed to get hostname: {}", e);
            None
        }
    };
    if running.is_some() {
        Some(HostNameState {
            running,
            config: get_config_hostname(),
        })
    } else {
        None
    }
}

const HOSTNAME_CONFIG_PATH: &str = "/etc/hostname";

fn get_config_hostname() -> Option<String> {
    if !std::path::Path::new("/etc/hostname").exists() {
        return Some("".to_string());
    }

    let mut fd = match std::fs::File::open(HOSTNAME_CONFIG_PATH) {
        Ok(fd) => fd,
        Err(_) => {
            return None;
        }
    };

    let mut contents = String::new();
    if let Err(e) = fd.read_to_string(&mut contents) {
        log::error!(
            "Failed to read hostname config {}: {}",
            HOSTNAME_CONFIG_PATH,
            e
        );
        None
    } else {
        contents.truncate(contents.as_str().trim().len());
        Some(contents)
    }
}

pub(crate) fn set_running_hostname(hostname: &str) -> Result<(), NmstateError> {
    if hostname.is_empty() {
        let e = NmstateError::new(
            ErrorKind::InvalidArgument,
            "Cannot set empty runtime hostname".to_string(),
        );
        log::error!("{}", e);
        return Err(e);
    }
    if hostname.len() >= HOST_NAME_MAX {
        let e = NmstateError::new(
            ErrorKind::InvalidArgument,
            format!("hostname to long, should be less than {}", HOST_NAME_MAX),
        );
        log::error!("{}", e);
        return Err(e);
    }

    let os_str = std::ffi::OsStr::new(hostname);
    if nix::unistd::sethostname(&os_str).is_err() {
        let e = NmstateError::new(
            ErrorKind::InvalidArgument,
            format!(
                "Failed to set hostname {}, errno {}",
                hostname,
                nix::errno::errno()
            ),
        );
        log::error!("{}", e);
        return Err(e);
    }
    Ok(())
}
