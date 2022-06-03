use crate::nm::nm_dbus::{NmConnection, NmSettingVrf};

use crate::VrfConfig;

impl From<&VrfConfig> for NmSettingVrf {
    fn from(config: &VrfConfig) -> Self {
        let mut settings = NmSettingVrf::default();
        settings.table = Some(config.table_id);
        settings
    }
}

pub(crate) fn is_vrf_table_id_changed(
    new_nm_conn: &NmConnection,
    cur_nm_conn: &NmConnection,
) -> bool {
    if let (Some(new_vrf_conf), Some(cur_vrf_conf)) =
        (new_nm_conn.vrf.as_ref(), cur_nm_conn.vrf.as_ref())
    {
        new_vrf_conf.table != cur_vrf_conf.table
    } else {
        false
    }
}
