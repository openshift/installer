use std::str::FromStr;

use crate::{BaseInterface, VxlanConfig, VxlanInterface};

pub(crate) fn np_vxlan_to_nmstate(
    np_iface: &nispor::Iface,
    base_iface: BaseInterface,
) -> VxlanInterface {
    let vxlan_conf = np_iface.vxlan.as_ref().map(|np_vxlan_info| VxlanConfig {
        id: np_vxlan_info.vxlan_id,
        base_iface: np_vxlan_info.base_iface.clone(),
        remote: std::net::IpAddr::from_str(np_vxlan_info.remote.as_str()).ok(),
        dst_port: Some(np_vxlan_info.dst_port),
    });

    VxlanInterface {
        base: base_iface,
        vxlan: vxlan_conf,
    }
}
