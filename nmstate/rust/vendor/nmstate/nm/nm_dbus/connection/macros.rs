macro_rules! _connection_inner_string_member {
    ($self_struct: ident, $member: ident) => {
        $self_struct
            .connection
            .as_ref()
            .map(|conn| conn.$member.as_deref())
            .flatten()
    };
}

macro_rules! _from_map {
    ($map: ident, $remove: expr, $convert: expr) => {
        $map.remove($remove).map($convert).transpose().map_err(|e| {
            let e = NmError::from(e);
            NmError::new(
                e.kind,
                format!("key {} fail to convert: {}", $remove, e.msg),
            )
        })
    };
}

pub(crate) use _from_map;
