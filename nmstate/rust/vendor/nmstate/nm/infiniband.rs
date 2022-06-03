use crate::nm::nm_dbus::NmConnection;

use crate::InfiniBandInterface;

pub(crate) fn gen_nm_ib_setting(
    iface: &InfiniBandInterface,
    nm_conn: &mut NmConnection,
) {
    let mut nm_ib_set =
        nm_conn.infiniband.as_ref().cloned().unwrap_or_default();
    if let Some(ib_conf) = iface.ib.as_ref() {
        nm_ib_set.parent = ib_conf.base_iface.as_ref().and_then(|p| {
            if p.is_empty() {
                None
            } else {
                Some(p.clone())
            }
        });
        nm_ib_set.pkey = ib_conf.pkey.and_then(|p| {
            if p == u16::MAX {
                None
            } else {
                Some(i32::from(p))
            }
        });
        nm_ib_set.mode = Some(ib_conf.mode.to_string());
    }
    if let Some(mtu) = iface.base.mtu {
        nm_ib_set.mtu = Some(mtu as u32);
    }
    nm_conn.infiniband = Some(nm_ib_set)
}
