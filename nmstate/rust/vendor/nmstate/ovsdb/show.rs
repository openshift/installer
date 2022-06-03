use std::collections::HashMap;
use std::iter::FromIterator;

use crate::{
    ovsdb::db::OvsDbConnection, Interface, NetworkState, NmstateError,
    OvsBridgeInterface, OvsDbIfaceConfig, UnknownInterface,
};

pub(crate) fn ovsdb_is_running() -> bool {
    if let Ok(mut cli) = OvsDbConnection::new() {
        cli.check_connection()
    } else {
        false
    }
}

pub(crate) fn ovsdb_retrieve() -> Result<NetworkState, NmstateError> {
    let mut ret = NetworkState::new();
    ret.prop_list.push("interfaces");
    ret.prop_list.push("ovsdb");
    let mut cli = OvsDbConnection::new()?;
    for mut ovsdb_iface in cli.get_ovs_ifaces()? {
        let mut iface = Interface::Unknown(UnknownInterface::new());
        iface.base_iface_mut().name = ovsdb_iface.name;
        iface.base_iface_mut().prop_list.push("ovsdb");
        iface.base_iface_mut().ovsdb = Some(OvsDbIfaceConfig {
            external_ids: Some(HashMap::from_iter(
                ovsdb_iface.external_ids.drain().map(|(k, v)| (k, Some(v))),
            )),
        });
        ret.append_interface_data(iface);
    }
    for mut ovsdb_iface in cli.get_ovs_bridges()? {
        let mut iface = Interface::OvsBridge(OvsBridgeInterface::new());
        iface.base_iface_mut().name = ovsdb_iface.name;
        iface.base_iface_mut().prop_list.push("ovsdb");
        iface.base_iface_mut().ovsdb = Some(OvsDbIfaceConfig {
            external_ids: Some(HashMap::from_iter(
                ovsdb_iface.external_ids.drain().map(|(k, v)| (k, Some(v))),
            )),
        });
        ret.append_interface_data(iface);
    }

    ret.ovsdb = cli.get_ovsdb_global_conf()?;

    Ok(ret)
}
