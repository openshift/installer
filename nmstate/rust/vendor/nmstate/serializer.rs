pub(crate) fn is_option_string_empty(data: &Option<String>) -> bool {
    if let Some(s) = data {
        s.is_empty()
    } else {
        true
    }
}
