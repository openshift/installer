use crate::{
    BaseInterface, InfiniBandConfig, InfiniBandInterface, InfiniBandMode,
};

impl From<nispor::IpoibMode> for InfiniBandMode {
    fn from(m: nispor::IpoibMode) -> Self {
        match m {
            nispor::IpoibMode::Datagram => Self::Datagram,
            nispor::IpoibMode::Connected => Self::Connected,
            _ => {
                log::warn!("Unknown IP over IB mode {:?}", m);
                Self::default()
            }
        }
    }
}

pub(crate) fn np_ib_to_nmstate(
    np_iface: &nispor::Iface,
    base_iface: BaseInterface,
) -> InfiniBandInterface {
    let ib_conf = np_iface.ipoib.as_ref().map(|np_ib_info| InfiniBandConfig {
        mode: np_ib_info.mode.into(),
        base_iface: np_ib_info.base_iface.clone(),
        pkey: Some(np_ib_info.pkey),
    });

    InfiniBandInterface {
        base: base_iface,
        ib: ib_conf,
    }
}
