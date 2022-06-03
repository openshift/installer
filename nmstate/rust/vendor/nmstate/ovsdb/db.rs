use std::collections::HashMap;
use std::convert::TryFrom;
use std::convert::TryInto;

use serde_json::{Map, Value};

use crate::{
    ovsdb::json_rpc::OvsDbJsonRpc, ErrorKind, NmstateError, OvsDbGlobalConfig,
};

const OVS_DB_NAME: &str = "Open_vSwitch";
pub(crate) const GLOBAL_CONFIG_TABLE: &str = "Open_vSwitch";
const NM_RESERVED_EXTERNAL_ID: &str = "NM.connection.uuid";

const DEFAULT_OVS_DB_SOCKET_PATH: &str = "/run/openvswitch/db.sock";

#[derive(Debug)]
pub(crate) struct OvsDbConnection {
    rpc: OvsDbJsonRpc,
}

#[derive(Debug, Clone, Default, PartialEq, Eq)]
pub(crate) struct OvsDbSelect {
    table: String,
    conditions: Vec<OvsDbCondition>,
    columns: Option<Vec<&'static str>>,
}

#[derive(Debug, Clone, Default, PartialEq, Eq)]
pub(crate) struct OvsDbCondition {
    column: String,
    function: String,
    value: Value,
}

impl OvsDbCondition {
    fn to_value(&self) -> Value {
        Value::Array(vec![
            Value::String(self.column.to_string()),
            Value::String(self.function.to_string()),
            self.value.clone(),
        ])
    }
}

impl OvsDbSelect {
    fn to_value(&self) -> Value {
        let mut ret = Map::new();
        ret.insert("op".to_string(), Value::String("select".to_string()));
        ret.insert("table".to_string(), Value::String(self.table.clone()));
        let condition_values: Vec<Value> =
            self.conditions.iter().map(|c| c.to_value()).collect();
        ret.insert("where".to_string(), Value::Array(condition_values));
        if let Some(columns) = self.columns.as_ref() {
            ret.insert(
                "columns".to_string(),
                Value::Array(
                    columns
                        .as_slice()
                        .iter()
                        .map(|c| Value::String(c.to_string()))
                        .collect(),
                ),
            );
        }
        Value::Object(ret)
    }
}

impl OvsDbConnection {
    // TODO: support environment variable OVS_DB_UNIX_SOCKET_PATH
    pub(crate) fn new() -> Result<Self, NmstateError> {
        Ok(Self {
            rpc: OvsDbJsonRpc::connect(DEFAULT_OVS_DB_SOCKET_PATH)?,
        })
    }

    pub(crate) fn check_connection(&mut self) -> bool {
        if let Ok(reply) = self.rpc.exec("list_dbs", &Value::Array(vec![])) {
            if let Some(dbs) = reply.as_array() {
                dbs.iter().any(|db| db.as_str() == Some(OVS_DB_NAME))
            } else {
                false
            }
        } else {
            false
        }
    }

    fn _get_ovs_ifaec(
        &mut self,
        table_name: &str,
    ) -> Result<Vec<OvsDbIface>, NmstateError> {
        let select = OvsDbSelect {
            table: table_name.to_string(),
            conditions: vec![],
            columns: Some(vec!["external_ids", "name"]),
        };
        let mut ret: Vec<OvsDbIface> = Vec::new();
        match self.rpc.exec(
            "transact",
            &Value::Array(vec![
                Value::String(OVS_DB_NAME.to_string()),
                select.to_value(),
            ]),
        )? {
            Value::Array(reply) => {
                if let Some(ovsdb_ifaces) = reply
                    .get(0)
                    .and_then(|v| v.as_object())
                    .and_then(|v| v.get("rows"))
                    .and_then(|v| v.as_array())
                {
                    for ovsdb_iface in ovsdb_ifaces {
                        ret.push(ovsdb_iface.try_into()?);
                    }
                    Ok(ret)
                } else {
                    let e = NmstateError::new(
                        ErrorKind::PluginFailure,
                        format!(
                            "Invalid reply from OVSDB for querying {} \
                            table: {:?}",
                            table_name, reply
                        ),
                    );
                    log::error!("{}", e);
                    Err(e)
                }
            }
            reply => {
                let e = NmstateError::new(
                    ErrorKind::PluginFailure,
                    format!(
                        "Invalid reply from OVSDB for querying {} table: {:?}",
                        table_name, reply
                    ),
                );
                log::error!("{}", e);
                Err(e)
            }
        }
    }

    pub(crate) fn get_ovs_ifaces(
        &mut self,
    ) -> Result<Vec<OvsDbIface>, NmstateError> {
        self._get_ovs_ifaec("Interface")
    }

    pub(crate) fn get_ovs_bridges(
        &mut self,
    ) -> Result<Vec<OvsDbIface>, NmstateError> {
        self._get_ovs_ifaec("Bridge")
    }

    pub(crate) fn get_ovsdb_global_conf(
        &mut self,
    ) -> Result<OvsDbGlobalConfig, NmstateError> {
        let select = OvsDbSelect {
            table: GLOBAL_CONFIG_TABLE.to_string(),
            conditions: vec![],
            columns: Some(vec!["external_ids", "other_config"]),
        };
        match self.rpc.exec(
            "transact",
            &Value::Array(vec![
                Value::String(OVS_DB_NAME.to_string()),
                select.to_value(),
            ]),
        )? {
            Value::Array(reply) => {
                if let Some(global_conf) = reply
                    .get(0)
                    .and_then(|v| v.as_object())
                    .and_then(|v| v.get("rows"))
                    .and_then(|v| v.as_array())
                    .and_then(|v| v.get(0))
                    .and_then(|v| v.as_object())
                {
                    Ok(global_conf.into())
                } else {
                    let e = NmstateError::new(
                        ErrorKind::PluginFailure,
                        format!(
                        "Invalid reply from OVSDB for querying {} table: {:?}",
                        GLOBAL_CONFIG_TABLE, reply
                    ),
                    );
                    log::error!("{}", e);
                    Err(e)
                }
            }
            reply => {
                let e = NmstateError::new(
                    ErrorKind::PluginFailure,
                    format!(
                        "Invalid reply from OVSDB for querying {} table: {:?}",
                        GLOBAL_CONFIG_TABLE, reply
                    ),
                );
                log::error!("{}", e);
                Err(e)
            }
        }
    }
    pub(crate) fn apply_global_conf(
        &mut self,
        ovs_conf: &OvsDbGlobalConfig,
    ) -> Result<(), NmstateError> {
        let update: OvsDbUpdate = ovs_conf.into();
        self.rpc.exec(
            "transact",
            &Value::Array(vec![
                Value::String(OVS_DB_NAME.to_string()),
                update.to_value(),
            ]),
        )?;
        Ok(())
    }
}

#[derive(Debug, Default)]
pub(crate) struct OvsDbIface {
    pub(crate) name: String,
    pub(crate) external_ids: HashMap<String, String>,
}

impl TryFrom<&Value> for OvsDbIface {
    type Error = NmstateError;
    fn try_from(v: &Value) -> Result<OvsDbIface, Self::Error> {
        let e = NmstateError::new(
            ErrorKind::PluginFailure,
            format!("Failed to parse OVS Interface info from : {:?}", v),
        );
        let mut ret = OvsDbIface::default();
        if let Value::Object(v) = v {
            if let (Some(Value::String(n)), Some(Value::Array(ids))) =
                (v.get("name"), v.get("external_ids"))
            {
                ret.name = n.to_string();
                ret.external_ids = parse_str_map(ids);
                return Ok(ret);
            }
        }
        log::error!("{}", e);
        Err(e)
    }
}

pub(crate) fn parse_str_map(v: &[Value]) -> HashMap<String, String> {
    let mut ret = HashMap::new();
    if let Some(ids) = v.get(1).and_then(|i| i.as_array()) {
        for kv in ids {
            if let Some(kv) = kv.as_array() {
                if let (Some(Value::String(k)), Some(Value::String(v))) =
                    (kv.get(0), kv.get(1))
                {
                    if k == NM_RESERVED_EXTERNAL_ID {
                        continue;
                    }
                    ret.insert(k.to_string(), v.to_string());
                }
            }
        }
    }
    ret
}

#[derive(Debug, Clone, Default, PartialEq, Eq)]
pub(crate) struct OvsDbUpdate {
    pub(crate) table: String,
    pub(crate) conditions: Vec<OvsDbCondition>,
    pub(crate) row: HashMap<String, Value>,
}

impl OvsDbUpdate {
    fn to_value(&self) -> Value {
        let mut ret = Map::new();
        ret.insert("op".to_string(), Value::String("update".to_string()));
        ret.insert("table".to_string(), Value::String(self.table.clone()));
        let condition_values: Vec<Value> =
            self.conditions.iter().map(|c| c.to_value()).collect();
        ret.insert("where".to_string(), Value::Array(condition_values));
        let mut row_map = Map::new();
        for (k, v) in self.row.iter() {
            row_map.insert(k.to_string(), v.clone());
        }
        ret.insert("row".to_string(), Value::Object(row_map));
        Value::Object(ret)
    }
}
