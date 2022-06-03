use crate::{
    BaseInterface, MacVlanConfig, MacVlanInterface, MacVlanMode, MacVtapConfig,
    MacVtapInterface, MacVtapMode,
};

const MACVLAN_FLAG_NOPROMISC: u16 = 1;
const MACVTAP_FLAG_NOPROMISC: u16 = 1;

pub(crate) fn np_mac_vlan_to_nmstate(
    np_iface: &nispor::Iface,
    base_iface: BaseInterface,
) -> MacVlanInterface {
    let vlan_conf =
        np_iface
            .mac_vlan
            .as_ref()
            .map(|np_vlan_info| MacVlanConfig {
                mode: match &np_vlan_info.mode {
                    nispor::MacVlanMode::Private => MacVlanMode::Private,
                    nispor::MacVlanMode::Vepa => MacVlanMode::Vepa,
                    nispor::MacVlanMode::Bridge => MacVlanMode::Bridge,
                    nispor::MacVlanMode::PassThrough => MacVlanMode::Passthru,
                    nispor::MacVlanMode::Source => MacVlanMode::Source,
                    _ => {
                        log::warn!(
                            "Unknown supported MacVlan mode {:?}",
                            np_vlan_info.mode
                        );
                        MacVlanMode::Unknown
                    }
                },
                accept_all_mac: Some(
                    np_vlan_info.flags & MACVLAN_FLAG_NOPROMISC == 0,
                ),
                base_iface: np_vlan_info.base_iface.clone(),
            });

    MacVlanInterface {
        base: base_iface,
        mac_vlan: vlan_conf,
    }
}

// TODO: remove below duplicate code
pub(crate) fn np_mac_vtap_to_nmstate(
    np_iface: &nispor::Iface,
    base_iface: BaseInterface,
) -> MacVtapInterface {
    let vtap_conf =
        np_iface
            .mac_vtap
            .as_ref()
            .map(|np_vtap_info| MacVtapConfig {
                mode: match &np_vtap_info.mode {
                    nispor::MacVtapMode::Private => MacVtapMode::Private,
                    nispor::MacVtapMode::Vepa => MacVtapMode::Vepa,
                    nispor::MacVtapMode::Bridge => MacVtapMode::Bridge,
                    nispor::MacVtapMode::PassThrough => MacVtapMode::Passthru,
                    nispor::MacVtapMode::Source => MacVtapMode::Source,
                    _ => {
                        log::warn!(
                            "Unknown supported MacVtap mode {:?}",
                            np_vtap_info.mode
                        );
                        MacVtapMode::Unknown
                    }
                },
                accept_all_mac: Some(
                    np_vtap_info.flags & MACVTAP_FLAG_NOPROMISC == 0,
                ),
                base_iface: np_vtap_info.base_iface.clone(),
            });

    MacVtapInterface {
        base: base_iface,
        mac_vtap: vtap_conf,
    }
}
