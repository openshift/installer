use serde_json::Value;

use crate::NetworkState;

fn _get_json_value_difference<'a, 'b>(
    reference: String,
    desire: &'a Value,
    current: &'b Value,
) -> Option<(String, &'a Value, &'b Value)> {
    match (desire, current) {
        (Value::Bool(des), Value::Bool(cur)) => {
            if des != cur {
                Some((reference, desire, current))
            } else {
                None
            }
        }
        (Value::Number(des), Value::Number(cur)) => {
            if des != cur {
                Some((reference, desire, current))
            } else {
                None
            }
        }
        (Value::String(des), Value::String(cur)) => {
            if des != cur {
                if des == NetworkState::PASSWORD_HID_BY_NMSTATE {
                    None
                } else {
                    Some((reference, desire, current))
                }
            } else {
                None
            }
        }
        (Value::Array(des), Value::Array(cur)) => {
            if des.len() != cur.len() {
                Some((reference, desire, current))
            } else {
                for (index, des_element) in des.iter().enumerate() {
                    // The [] is safe as we already checked the length
                    let cur_element = &cur[index];
                    if let Some(difference) = get_json_value_difference(
                        format!("{}[{}]", &reference, index),
                        des_element,
                        cur_element,
                    ) {
                        return Some(difference);
                    }
                }
                None
            }
        }
        (Value::Object(des), Value::Object(cur)) => {
            for (key, des_value) in des.iter() {
                let reference = format!("{}.{}", reference, key);
                if let Some(cur_value) = cur.get(key) {
                    if let Some(difference) = get_json_value_difference(
                        reference.clone(),
                        des_value,
                        cur_value,
                    ) {
                        return Some(difference);
                    }
                } else if des_value != &Value::Null {
                    return Some((reference, des_value, &Value::Null));
                }
            }
            None
        }
        (Value::Null, _) => None,
        (_, _) => Some((reference, desire, current)),
    }
}

pub(crate) fn get_json_value_difference<'a, 'b>(
    reference: String,
    desire: &'a Value,
    current: &'b Value,
) -> Option<(String, &'a Value, &'b Value)> {
    if let Some((reference, desire, current)) =
        _get_json_value_difference(reference, desire, current)
    {
        if should_ignore(reference.as_str(), desire, current) {
            None
        } else {
            Some((reference, desire, current))
        }
    } else {
        None
    }
}

fn should_ignore(reference: &str, desire: &Value, current: &Value) -> bool {
    if reference.contains("interface.link-aggregation.options") {
        // Per oVirt request, bond option difference should not
        // fail verification.
        log::warn!(
            "Bond option miss-match: {} desire '{}', current '{}'",
            reference,
            desire,
            current
        );
        true
    } else {
        false
    }
}
