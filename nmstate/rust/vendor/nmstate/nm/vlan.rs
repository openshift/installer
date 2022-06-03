use crate::nm::nm_dbus::{NmConnection, NmSettingVlan};

use crate::VlanConfig;

impl From<&VlanConfig> for NmSettingVlan {
    fn from(config: &VlanConfig) -> Self {
        let mut settings = NmSettingVlan::default();
        settings.id = Some(config.id.into());
        settings.parent = Some(config.base_iface.clone());
        settings
    }
}

pub(crate) fn is_vlan_id_changed(
    new_nm_conn: &NmConnection,
    cur_nm_conn: &NmConnection,
) -> bool {
    if let (Some(new_vlan_conf), Some(cur_vlan_conf)) =
        (new_nm_conn.vlan.as_ref(), cur_nm_conn.vlan.as_ref())
    {
        new_vlan_conf.id != cur_vlan_conf.id
    } else {
        false
    }
}
