use crate::nm::nm_dbus::{NmConnection, NmSettingVxlan};

use crate::VxlanConfig;

impl From<&VxlanConfig> for NmSettingVxlan {
    fn from(config: &VxlanConfig) -> Self {
        let mut setting = NmSettingVxlan::default();
        setting.id = Some(config.id);
        setting.parent = Some(config.base_iface.clone());
        if let Some(v) = config.remote.as_ref() {
            setting.remote = Some(v.to_string());
        }
        if let Some(v) = config.dst_port {
            setting.dst_port = Some(v.into())
        }
        setting
    }
}

pub(crate) fn is_vxlan_id_changed(
    new_nm_conn: &NmConnection,
    cur_nm_conn: &NmConnection,
) -> bool {
    if let (Some(new_vxlan_conf), Some(cur_vxlan_conf)) =
        (new_nm_conn.vxlan.as_ref(), cur_nm_conn.vxlan.as_ref())
    {
        new_vxlan_conf.id != cur_vxlan_conf.id
    } else {
        false
    }
}
