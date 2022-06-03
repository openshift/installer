use std::collections::HashMap;

use crate::nm::nm_dbus::NmDevice;

pub(crate) fn create_index_for_nm_devs(
    nm_devs: &[NmDevice],
) -> HashMap<(String, String), &NmDevice> {
    let mut ret: HashMap<(String, String), &NmDevice> = HashMap::new();
    for nm_dev in nm_devs {
        ret.insert(
            (nm_dev.name.to_string(), nm_dev.iface_type.to_string()),
            nm_dev,
        );
    }
    ret
}
