use serde::{Deserialize, Deserializer, Serialize};
use static_assertions::assert_impl_all;
use std::{collections::HashMap, convert::TryFrom, hash::BuildHasher};

use crate::{
    derive::Type, Array, Dict, Fd, ObjectPath, OwnedObjectPath, OwnedSignature, Signature, Str,
    Structure, Value,
};

#[cfg(feature = "gvariant")]
use crate::Maybe;

// FIXME: Replace with a generic impl<T: TryFrom<Value>> TryFrom<OwnedValue> for T?
// https://gitlab.freedesktop.org/dbus/zbus/-/issues/138

/// Owned [`Value`](enum.Value.html)
#[derive(Debug, Clone, PartialEq, Serialize, Type)]
pub struct OwnedValue(Value<'static>);

assert_impl_all!(OwnedValue: Send, Sync, Unpin);

impl OwnedValue {
    pub(crate) fn into_inner(self) -> Value<'static> {
        self.0
    }
}

macro_rules! ov_try_from {
    ($to:ty) => {
        impl<'a> TryFrom<OwnedValue> for $to {
            type Error = crate::Error;

            fn try_from(v: OwnedValue) -> Result<Self, Self::Error> {
                <$to>::try_from(v.0)
            }
        }
    };
}

macro_rules! ov_try_from_ref {
    ($to:ty) => {
        impl<'a> TryFrom<&'a OwnedValue> for $to {
            type Error = crate::Error;

            fn try_from(v: &'a OwnedValue) -> Result<Self, Self::Error> {
                <$to>::try_from(&v.0)
            }
        }
    };
}

ov_try_from!(u8);
ov_try_from!(bool);
ov_try_from!(i16);
ov_try_from!(u16);
ov_try_from!(i32);
ov_try_from!(u32);
ov_try_from!(i64);
ov_try_from!(u64);
ov_try_from!(f64);
ov_try_from!(String);
ov_try_from!(Signature<'a>);
ov_try_from!(OwnedSignature);
ov_try_from!(ObjectPath<'a>);
ov_try_from!(OwnedObjectPath);
ov_try_from!(Array<'a>);
ov_try_from!(Dict<'a, 'a>);
#[cfg(feature = "gvariant")]
ov_try_from!(Maybe<'a>);
ov_try_from!(Str<'a>);
ov_try_from!(Structure<'a>);
ov_try_from!(Fd);

ov_try_from_ref!(u8);
ov_try_from_ref!(bool);
ov_try_from_ref!(i16);
ov_try_from_ref!(u16);
ov_try_from_ref!(i32);
ov_try_from_ref!(u32);
ov_try_from_ref!(i64);
ov_try_from_ref!(u64);
ov_try_from_ref!(f64);
ov_try_from_ref!(&'a str);
ov_try_from_ref!(&'a Signature<'a>);
ov_try_from_ref!(&'a ObjectPath<'a>);
ov_try_from_ref!(&'a Array<'a>);
ov_try_from_ref!(&'a Dict<'a, 'a>);
ov_try_from_ref!(&'a Str<'a>);
ov_try_from_ref!(&'a Structure<'a>);
#[cfg(feature = "gvariant")]
ov_try_from_ref!(&'a Maybe<'a>);
ov_try_from_ref!(Fd);

impl<'a, T> TryFrom<OwnedValue> for Vec<T>
where
    T: TryFrom<Value<'a>>,
    T::Error: Into<crate::Error>,
{
    type Error = crate::Error;

    fn try_from(value: OwnedValue) -> Result<Self, Self::Error> {
        if let Value::Array(v) = value.0 {
            Self::try_from(v)
        } else {
            Err(crate::Error::IncorrectType)
        }
    }
}

#[cfg(feature = "enumflags2")]
impl<'a, F> TryFrom<OwnedValue> for enumflags2::BitFlags<F>
where
    F: enumflags2::RawBitFlags,
    F::Type: TryFrom<Value<'a>, Error = crate::Error>,
{
    type Error = crate::Error;

    fn try_from(value: OwnedValue) -> Result<Self, Self::Error> {
        Self::try_from(value.0)
    }
}

impl<'k, 'v, K, V, H> TryFrom<OwnedValue> for HashMap<K, V, H>
where
    K: crate::Basic + TryFrom<Value<'k>> + std::hash::Hash + std::cmp::Eq,
    V: TryFrom<Value<'v>>,
    H: BuildHasher + Default,
    K::Error: Into<crate::Error>,
    V::Error: Into<crate::Error>,
{
    type Error = crate::Error;

    fn try_from(value: OwnedValue) -> Result<Self, Self::Error> {
        if let Value::Dict(v) = value.0 {
            Self::try_from(v)
        } else {
            Err(crate::Error::IncorrectType)
        }
    }
}

// tuple conversions in `structure` module for avoiding code-duplication.

impl<'a> From<Value<'a>> for OwnedValue {
    fn from(v: Value<'a>) -> Self {
        // TODO: add into_owned, avoiding copy if already owned..
        OwnedValue(v.to_owned())
    }
}

impl<'a> From<&Value<'a>> for OwnedValue {
    fn from(v: &Value<'a>) -> Self {
        OwnedValue(v.to_owned())
    }
}

macro_rules! to_value {
    ($from:ty) => {
        impl<'a> From<$from> for OwnedValue {
            fn from(v: $from) -> Self {
                OwnedValue::from(<Value<'a>>::from(v))
            }
        }
    };
}

to_value!(u8);
to_value!(bool);
to_value!(i16);
to_value!(u16);
to_value!(i32);
to_value!(u32);
to_value!(i64);
to_value!(u64);
to_value!(f64);
to_value!(Array<'a>);
to_value!(Dict<'a, 'a>);
#[cfg(feature = "gvariant")]
to_value!(Maybe<'a>);
to_value!(Str<'a>);
to_value!(Signature<'a>);
to_value!(Structure<'a>);
to_value!(ObjectPath<'a>);
to_value!(Fd);

impl From<OwnedValue> for Value<'static> {
    fn from(v: OwnedValue) -> Value<'static> {
        v.into_inner()
    }
}

impl std::ops::Deref for OwnedValue {
    type Target = Value<'static>;

    fn deref(&self) -> &Self::Target {
        &self.0
    }
}

impl<'de> Deserialize<'de> for OwnedValue {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        Ok(Value::deserialize(deserializer)?.into())
    }
}

#[cfg(test)]
mod tests {
    use byteorder::LE;
    use std::{convert::TryFrom, error::Error, result::Result};

    use crate::{from_slice, to_bytes, EncodingContext, OwnedValue, Value};

    #[cfg(feature = "enumflags2")]
    #[test]
    fn bitflags() -> Result<(), Box<dyn Error>> {
        #[repr(u32)]
        #[derive(enumflags2::BitFlags, Copy, Clone, Debug)]
        pub enum Flaggy {
            One = 0x1,
            Two = 0x2,
        }

        let v = Value::from(0x2u32);
        let ov: OwnedValue = v.into();
        assert_eq!(<enumflags2::BitFlags<Flaggy>>::try_from(ov)?, Flaggy::Two);
        Ok(())
    }

    #[test]
    fn from_value() -> Result<(), Box<dyn Error>> {
        let v = Value::from("hi!");
        let ov: OwnedValue = v.into();
        assert_eq!(<&str>::try_from(&ov)?, "hi!");
        Ok(())
    }

    #[test]
    fn serde() -> Result<(), Box<dyn Error>> {
        let ec = EncodingContext::<LE>::new_dbus(0);
        let ov: OwnedValue = Value::from("hi!").into();
        let ser = to_bytes(ec, &ov)?;
        let de: Value<'_> = from_slice(&ser, ec)?;
        assert_eq!(<&str>::try_from(&de)?, "hi!");
        Ok(())
    }
}
