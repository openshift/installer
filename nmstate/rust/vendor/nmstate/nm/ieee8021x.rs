use crate::nm::nm_dbus::{NmConnection, NmSetting8021X};

use crate::{Ieee8021XConfig, Interface, NetworkState};

pub(crate) fn nm_802_1x_to_nmstate(
    nm_setting: &NmSetting8021X,
) -> Ieee8021XConfig {
    Ieee8021XConfig {
        identity: nm_setting.identity.clone(),
        private_key: nm_setting
            .private_key
            .as_deref()
            .and_then(vec_u8_to_file_path),
        eap: nm_setting.eap.clone(),
        client_cert: nm_setting
            .client_cert
            .as_deref()
            .and_then(vec_u8_to_file_path),
        ca_cert: nm_setting.ca_cert.as_deref().and_then(vec_u8_to_file_path),
        private_key_password: nm_setting.private_key_password.clone(),
    }
}

fn vec_u8_to_file_path(raw: &[u8]) -> Option<String> {
    match NmSetting8021X::glib_bytes_to_file_path(raw) {
        Ok(s) => Some(s),
        Err(e) => {
            log::error!(
                "Unsupported NetworkManager 802.1x glib bytes: {:?}, error: {}",
                raw,
                e
            );
            None
        }
    }
}

pub(crate) fn gen_nm_802_1x_setting(
    iface: &Interface,
    nm_conn: &mut NmConnection,
) {
    if let Some(conf) = iface.base_iface().ieee8021x.as_ref() {
        let mut nm_setting = NmSetting8021X::default();
        nm_setting.identity = conf.identity.clone();
        nm_setting.eap = conf.eap.clone();
        nm_setting.private_key = conf
            .private_key
            .as_deref()
            .map(NmSetting8021X::file_path_to_glib_bytes);
        nm_setting.client_cert = conf
            .client_cert
            .as_deref()
            .map(NmSetting8021X::file_path_to_glib_bytes);
        nm_setting.ca_cert = conf
            .ca_cert
            .as_deref()
            .map(NmSetting8021X::file_path_to_glib_bytes);
        if conf.private_key_password.as_deref()
            == Some(NetworkState::PASSWORD_HID_BY_NMSTATE)
        {
            if let Some(cur_pass) = nm_conn
                .ieee8021x
                .as_ref()
                .and_then(|c| c.private_key_password.as_deref())
            {
                nm_setting.private_key_password = Some(cur_pass.to_string());
            }
        } else {
            nm_setting.private_key_password = conf.private_key_password.clone();
        }
        nm_conn.ieee8021x = Some(nm_setting);
    }
}
