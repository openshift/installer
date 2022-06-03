use crate::{BaseInterface, EthernetInterface, VethConfig};

pub(crate) fn np_veth_to_nmstate(
    np_iface: &nispor::Iface,
    base_iface: BaseInterface,
) -> EthernetInterface {
    let veth_conf = np_iface.veth.as_ref().and_then(|np_veth_info| {
        if np_veth_info.peer.as_str().parse::<u32>().is_ok() {
            // If veth peer is interface index, it means its veth peer is in
            // another network namespace, we should treat this interface
            // as ethernet
            None
        } else {
            Some(VethConfig {
                peer: np_veth_info.peer.clone(),
            })
        }
    });

    EthernetInterface {
        base: base_iface,
        veth: veth_conf,
        // TODO: Filling the ethernet section
        ..Default::default()
    }
}

pub(crate) fn nms_veth_conf_to_np(
    nms_veth_conf: Option<&VethConfig>,
) -> Option<nispor::VethConf> {
    nms_veth_conf.map(|nms_veth_conf| {
        let mut veth_conf = nispor::VethConf::default();
        veth_conf.peer = nms_veth_conf.peer.to_string();
        veth_conf
    })
}
