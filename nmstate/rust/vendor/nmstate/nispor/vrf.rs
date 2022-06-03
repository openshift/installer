use crate::{BaseInterface, VrfConfig, VrfInterface};

pub(crate) fn np_vrf_to_nmstate(
    np_iface: &nispor::Iface,
    base_iface: BaseInterface,
) -> VrfInterface {
    let vrf_conf = np_iface.vrf.as_ref().map(|np_vrf_info| VrfConfig {
        table_id: np_vrf_info.table_id,
        port: Some(np_vrf_info.subordinates.clone()),
    });

    VrfInterface {
        base: base_iface,
        vrf: vrf_conf,
    }
}
