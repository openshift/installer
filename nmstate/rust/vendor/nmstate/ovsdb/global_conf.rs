use std::collections::HashMap;

use serde_json::{Map, Value};

use crate::{
    ovsdb::db::{parse_str_map, OvsDbUpdate, GLOBAL_CONFIG_TABLE},
    OvsDbGlobalConfig,
};

impl From<&Map<std::string::String, Value>> for OvsDbGlobalConfig {
    fn from(m: &Map<std::string::String, Value>) -> Self {
        let mut ret = Self::default();
        if let (Some(Value::Array(ids)), Some(Value::Array(other_cfg))) =
            (m.get("external_ids"), m.get("other_config"))
        {
            ret.external_ids = Some(convert_map(parse_str_map(ids)));
            ret.other_config = Some(convert_map(parse_str_map(other_cfg)));
        }
        ret
    }
}

// Convert HashMap<String, String> to HashMap<String, Option<String>>
fn convert_map(
    mut m: HashMap<String, String>,
) -> HashMap<String, Option<String>> {
    let mut ret = HashMap::new();
    for (k, v) in m.drain() {
        ret.insert(k, Some(v));
    }
    ret
}

impl From<&OvsDbGlobalConfig> for OvsDbUpdate {
    fn from(ovs_conf: &OvsDbGlobalConfig) -> Self {
        let mut row = HashMap::new();
        let mut value_array = Vec::new();
        for (k, v) in ovs_conf.get_external_ids().iter() {
            value_array.push(Value::Array(vec![
                Value::String(k.to_string()),
                Value::String(v.to_string()),
            ]));
        }
        row.insert(
            "external_ids".to_string(),
            Value::Array(vec![
                Value::String("map".to_string()),
                Value::Array(value_array),
            ]),
        );
        let mut value_array = Vec::new();
        for (k, v) in ovs_conf.get_other_config().iter() {
            value_array.push(Value::Array(vec![
                Value::String(k.to_string()),
                Value::String(v.to_string()),
            ]));
        }
        row.insert(
            "other_config".to_string(),
            Value::Array(vec![
                Value::String("map".to_string()),
                Value::Array(value_array),
            ]),
        );

        OvsDbUpdate {
            table: GLOBAL_CONFIG_TABLE.to_string(),
            conditions: vec![],
            row,
        }
    }
}
