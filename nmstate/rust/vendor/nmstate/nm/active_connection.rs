use std::collections::HashMap;

use crate::nm::nm_dbus::NmActiveConnection;

use crate::nm::connection::{
    NM_SETTING_VETH_SETTING_NAME, NM_SETTING_WIRED_SETTING_NAME,
};

pub(crate) fn create_index_for_nm_acs_by_name_type(
    nm_acs: &[NmActiveConnection],
) -> HashMap<(&str, &str), &NmActiveConnection> {
    let mut ret = HashMap::new();
    for nm_ac in nm_acs {
        let nm_iface_type = match nm_ac.iface_type.as_str() {
            NM_SETTING_VETH_SETTING_NAME => NM_SETTING_WIRED_SETTING_NAME,
            t => t,
        };
        ret.insert((nm_ac.iface_name.as_str(), nm_iface_type), nm_ac);
    }
    ret
}
