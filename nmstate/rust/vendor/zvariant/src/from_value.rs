#[cfg(feature = "gvariant")]
use crate::Maybe;
use crate::{
    Array, Dict, Error, Fd, ObjectPath, OwnedObjectPath, OwnedSignature, Signature, Str, Structure,
    Value,
};
use std::{collections::HashMap, convert::TryFrom, hash::BuildHasher};

macro_rules! value_try_from {
    ($kind:ident, $to:ty) => {
        impl<'a> TryFrom<Value<'a>> for $to {
            type Error = Error;

            fn try_from(value: Value<'a>) -> Result<Self, Self::Error> {
                if let Value::$kind(value) = value {
                    Ok(value.into())
                } else {
                    Err(Error::IncorrectType)
                }
            }
        }
    };
}

macro_rules! value_try_from_ref {
    ($kind:ident, $to:ty) => {
        impl<'a> TryFrom<&'a Value<'a>> for &'a $to {
            type Error = Error;

            fn try_from(value: &'a Value<'a>) -> Result<Self, Self::Error> {
                if let Value::$kind(value) = value {
                    Ok(value)
                } else {
                    Err(Error::IncorrectType)
                }
            }
        }
    };
}

macro_rules! value_try_from_ref_clone {
    ($kind:ident, $to:ty) => {
        impl<'a> TryFrom<&'a Value<'a>> for $to {
            type Error = Error;

            fn try_from(value: &'a Value<'_>) -> Result<Self, Self::Error> {
                if let Value::$kind(value) = value {
                    Ok(value.clone().into())
                } else {
                    Err(Error::IncorrectType)
                }
            }
        }
    };
}

macro_rules! value_try_from_all {
    ($from:ident, $to:ty) => {
        value_try_from!($from, $to);
        value_try_from_ref!($from, $to);
        value_try_from_ref_clone!($from, $to);
    };
}

value_try_from_all!(U8, u8);
value_try_from_all!(Bool, bool);
value_try_from_all!(I16, i16);
value_try_from_all!(U16, u16);
value_try_from_all!(I32, i32);
value_try_from_all!(U32, u32);
value_try_from_all!(I64, i64);
value_try_from_all!(U64, u64);
value_try_from_all!(F64, f64);
value_try_from_all!(Fd, Fd);

value_try_from_all!(Str, Str<'a>);
value_try_from_all!(Signature, Signature<'a>);
value_try_from_all!(ObjectPath, ObjectPath<'a>);
value_try_from_all!(Structure, Structure<'a>);
value_try_from_all!(Dict, Dict<'a, 'a>);
value_try_from_all!(Array, Array<'a>);
#[cfg(feature = "gvariant")]
value_try_from_all!(Maybe, Maybe<'a>);

value_try_from!(Str, String);
value_try_from_ref!(Str, str);
value_try_from_ref_clone!(Str, String);

impl<'a, T> TryFrom<Value<'a>> for Vec<T>
where
    T: TryFrom<Value<'a>>,
    T::Error: Into<crate::Error>,
{
    type Error = Error;

    fn try_from(value: Value<'a>) -> Result<Self, Self::Error> {
        if let Value::Array(v) = value {
            Self::try_from(v)
        } else {
            Err(Error::IncorrectType)
        }
    }
}

impl TryFrom<Value<'_>> for OwnedObjectPath {
    type Error = Error;

    fn try_from(value: Value<'_>) -> Result<Self, Self::Error> {
        ObjectPath::try_from(value).map(OwnedObjectPath::from)
    }
}

impl TryFrom<Value<'_>> for OwnedSignature {
    type Error = Error;

    fn try_from(value: Value<'_>) -> Result<Self, Self::Error> {
        Signature::try_from(value).map(OwnedSignature::from)
    }
}

// tuple conversions in `structure` module for avoiding code-duplication.

#[cfg(feature = "enumflags2")]
impl<'a, F> TryFrom<Value<'a>> for enumflags2::BitFlags<F>
where
    F: enumflags2::RawBitFlags,
    F::Type: TryFrom<Value<'a>, Error = Error>,
{
    type Error = Error;

    fn try_from(value: Value<'a>) -> Result<Self, Self::Error> {
        Self::from_bits(F::Type::try_from(value)?)
            .map_err(|_| Error::Message("Failed to convert to bitflags".into()))
    }
}

impl<'a, K, V, H> TryFrom<Value<'a>> for HashMap<K, V, H>
where
    K: crate::Basic + TryFrom<Value<'a>> + std::hash::Hash + std::cmp::Eq,
    V: TryFrom<Value<'a>>,
    H: BuildHasher + Default,
    K::Error: Into<crate::Error>,
    V::Error: Into<crate::Error>,
{
    type Error = crate::Error;

    fn try_from(value: Value<'a>) -> Result<Self, Self::Error> {
        if let Value::Dict(v) = value {
            Self::try_from(v)
        } else {
            Err(crate::Error::IncorrectType)
        }
    }
}

// This would be great but somehow it conflicts with some blanket generic implementations from
// core:
//
// impl<'a, T> TryFrom<Value<'a>> for Option<T>
//
// TODO: this could be useful
// impl<'a, 'b, T> TryFrom<&'a Value<'b>> for Vec<T>
// impl<'a, 'b, K, V, H> TryFrom<&'a Value<'v>> for HashMap<K, V, H>
// and more..
