use std::collections::HashMap;
use std::fmt::Write;

use log::error;

use super::{ErrorKind, NmError};

pub(crate) const DEFAULT_SEPARATOR: &str = ";";

pub(crate) fn keyfile_sections_to_string(
    sections: &[(&str, HashMap<String, zvariant::Value>)],
) -> Result<String, NmError> {
    let mut ret = String::new();
    for (section_name, data) in sections {
        let _ = writeln!(ret, "[{}]", section_name);
        // Sort the keys
        let mut keys: Vec<&String> = data.keys().collect();
        keys.sort_unstable();
        for key in keys {
            if let Some(v) = data.get(key) {
                let v = zvariant_value_to_string(v)?;
                if !v.is_empty() {
                    let _ = writeln!(ret, "{}={}", key, v);
                }
            }
        }
        ret += "\n";
    }
    if ret.ends_with("\n\n") {
        ret.pop();
    }
    Ok(ret)
}

fn zvariant_value_to_string(
    value: &zvariant::Value,
) -> Result<String, NmError> {
    match value {
        zvariant::Value::Bool(b) => Ok(if *b {
            "true".to_string()
        } else {
            "false".to_string()
        }),
        zvariant::Value::I32(d) => Ok(format!("{}", d)),
        zvariant::Value::U32(d) => Ok(format!("{}", d)),
        zvariant::Value::U8(d) => Ok(format!("{}", d)),
        zvariant::Value::U16(d) => Ok(format!("{}", d)),
        zvariant::Value::I16(d) => Ok(format!("{}", d)),
        zvariant::Value::U64(d) => Ok(format!("{}", d)),
        zvariant::Value::I64(d) => Ok(format!("{}", d)),
        zvariant::Value::Dict(d) => {
            let e = NmError::new(
                ErrorKind::Bug,
                format!("Cannot convert Dict {:?} to key file format", d),
            );
            log::error!("{}", e);
            Err(e)
        }
        zvariant::Value::Array(a) => {
            let mut ret = String::new();
            for item in a.get() {
                ret += &zvariant_value_to_string(item)?;
                ret += DEFAULT_SEPARATOR;
            }
            ret.pop();
            Ok(ret)
        }
        zvariant::Value::Str(s) => Ok(s.as_str().to_string()),
        _ => {
            let e = NmError::new(
                ErrorKind::Bug,
                format!("BUG: Unknown value type in section {:?}", value),
            );
            error!("{}", e);
            Err(e)
        }
    }
}
