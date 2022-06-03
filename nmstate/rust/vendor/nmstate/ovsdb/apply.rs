use crate::{ovsdb::db::OvsDbConnection, NetworkState, NmstateError};

pub(crate) fn ovsdb_apply(
    desired: &NetworkState,
    current: &NetworkState,
) -> Result<(), NmstateError> {
    if desired.ovsdb.external_ids.is_some()
        || desired.ovsdb.other_config.is_some()
    {
        let mut cli = OvsDbConnection::new()?;
        let mut desired = desired.ovsdb.clone();
        desired.merge(&current.ovsdb);
        cli.apply_global_conf(&desired)
    } else {
        Ok(())
    }
}
