use crate::nm::nm_dbus::NmConnection;

use crate::{
    nm::version::nm_supports_accept_all_mac_addresses_mode, Interface,
};

pub(crate) fn gen_nm_wired_setting(
    iface: &Interface,
    nm_conn: &mut NmConnection,
) {
    let mut nm_wired_set = nm_conn.wired.as_ref().cloned().unwrap_or_default();

    let mut flag_need_wired = false;

    let base_iface = iface.base_iface();

    if let Some(mac) = &base_iface.mac_address {
        nm_wired_set.cloned_mac_address = Some(mac.to_string());
        flag_need_wired = true;
    }
    if let Some(mtu) = &base_iface.mtu {
        nm_wired_set.mtu = Some(*mtu as u32);
        flag_need_wired = true;
    }

    if let Interface::Ethernet(eth_iface) = iface {
        if let Some(eth_conf) = eth_iface.ethernet.as_ref() {
            match eth_conf.auto_neg {
                Some(true) => {
                    flag_need_wired = true;
                    nm_wired_set.auto_negotiate = Some(true);
                    nm_wired_set.speed = None;
                    nm_wired_set.duplex = None;
                }
                Some(false) => {
                    flag_need_wired = true;
                    nm_wired_set.auto_negotiate = Some(false);
                    if let Some(v) = eth_conf.speed {
                        nm_wired_set.speed = Some(v);
                    }
                    if let Some(v) = eth_conf.duplex {
                        nm_wired_set.duplex = Some(format!("{}", v));
                    }
                }
                None => (),
            }
        }
    }

    if let Some(accept_all_mac_addresses) = &base_iface.accept_all_mac_addresses
    {
        if nm_supports_accept_all_mac_addresses_mode().unwrap_or_default() {
            nm_wired_set.accept_all_mac_addresses =
                Some(i32::from(*accept_all_mac_addresses));
            flag_need_wired = true;
        }
    }

    if flag_need_wired {
        nm_conn.wired = Some(nm_wired_set);
    }
}
