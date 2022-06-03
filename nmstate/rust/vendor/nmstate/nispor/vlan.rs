use crate::{BaseInterface, VlanConfig, VlanInterface};

pub(crate) fn np_vlan_to_nmstate(
    np_iface: &nispor::Iface,
    base_iface: BaseInterface,
) -> VlanInterface {
    let vlan_conf = np_iface.vlan.as_ref().map(|np_vlan_info| VlanConfig {
        id: np_vlan_info.vlan_id,
        base_iface: np_vlan_info.base_iface.clone(),
    });

    VlanInterface {
        base: base_iface,
        vlan: vlan_conf,
    }
}

pub(crate) fn nms_vlan_conf_to_np(
    nms_vlan_conf: Option<&VlanConfig>,
) -> Option<nispor::VlanConf> {
    nms_vlan_conf.map(|nms_vlan_conf| {
        let mut np_vlan_conf = nispor::VlanConf::default();
        np_vlan_conf.vlan_id = nms_vlan_conf.id;
        np_vlan_conf.base_iface = nms_vlan_conf.base_iface.clone();
        np_vlan_conf
    })
}
