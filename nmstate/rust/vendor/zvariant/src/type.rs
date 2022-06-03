use crate::{utils::*, Signature};
use serde::de::{Deserialize, DeserializeSeed};
use std::{convert::TryInto, marker::PhantomData};

/// Trait implemented by all serializable types.
///
/// This very simple trait provides the signature for the implementing type. Since the [D-Bus type
/// system] relies on these signatures, our [serialization and deserialization] API requires this
/// trait in addition to [`Serialize`] and [`Deserialize`], respectively.
///
/// Implementation is provided for all the [basic types] and blanket implementations for common
/// container types, such as, arrays, slices, tuples, [`Vec`] and [`HashMap`]. For easy
/// implementation for custom types, use `Type` derive macro from [zvariant_derive] crate.
///
/// If your type's signature cannot be determined statically, you should implement the
/// [DynamicType] trait instead, which is otherwise automatically implemented if you implement this
/// trait.
///
/// [D-Bus type system]: https://dbus.freedesktop.org/doc/dbus-specification.html#type-system
/// [serialization and deserialization]: index.html#functions
/// [`Serialize`]: https://docs.serde.rs/serde/trait.Serialize.html
/// [`Deserialize`]: https://docs.serde.rs/serde/de/trait.Deserialize.html
/// [basic types]: trait.Basic.html
/// [`Vec`]: https://doc.rust-lang.org/std/vec/struct.Vec.html
/// [`HashMap`]: https://doc.rust-lang.org/std/collections/struct.HashMap.html
/// [zvariant_derive]: https://docs.rs/zvariant_derive/2.0.0/zvariant_derive/
pub trait Type {
    /// Get the signature for the implementing type.
    ///
    /// # Example
    ///
    /// ```
    /// use std::collections::HashMap;
    /// use zvariant::Type;
    ///
    /// assert_eq!(u32::signature(), "u");
    /// assert_eq!(String::signature(), "s");
    /// assert_eq!(<(u32, &str, u64)>::signature(), "(ust)");
    /// assert_eq!(<(u32, &str, &[u64])>::signature(), "(usat)");
    /// assert_eq!(<HashMap<u8, &str>>::signature(), "a{ys}");
    /// ```
    fn signature() -> Signature<'static>;
}

/// Types with dynamic signatures.
///
/// Prefer implementing [Type] if possible, but if the actual signature of your type cannot be
/// determined until runtime, you can implement this type to support serialization.  You should
/// also implement [DynamicDeserialize] for deserialization.
pub trait DynamicType {
    /// Get the signature for the implementing type.
    ///
    /// See [Type::signature] for details.
    fn dynamic_signature(&self) -> Signature<'_>;
}

/// Types that deserialize based on dynamic signatures.
///
/// Prefer implementing [Type] and [Deserialize] if possible, but if the actual signature of your
/// type cannot be determined until runtime, you should implement this type to support
/// deserialization given a signature.
pub trait DynamicDeserialize<'de>: DynamicType {
    /// A [DeserializeSeed] implementation for this type.
    type Deserializer: DeserializeSeed<'de, Value = Self> + DynamicType;

    /// Get a deserializer compatible with this signature.
    fn deserializer_for_signature<S>(signature: S) -> zvariant::Result<Self::Deserializer>
    where
        S: TryInto<Signature<'de>>,
        S::Error: Into<zvariant::Error>;
}

impl<T> DynamicType for T
where
    T: Type + ?Sized,
{
    fn dynamic_signature(&self) -> Signature<'_> {
        <T as Type>::signature()
    }
}

impl<T> Type for PhantomData<T>
where
    T: Type + ?Sized,
{
    fn signature() -> Signature<'static> {
        T::signature()
    }
}

impl<'de, T> DynamicDeserialize<'de> for T
where
    T: Type + ?Sized + Deserialize<'de>,
{
    type Deserializer = PhantomData<T>;

    fn deserializer_for_signature<S>(signature: S) -> zvariant::Result<Self::Deserializer>
    where
        S: TryInto<Signature<'de>>,
        S::Error: Into<zvariant::Error>,
    {
        let mut expected = <T as Type>::signature();
        let original = signature.try_into().map_err(Into::into)?;

        if original == expected {
            return Ok(PhantomData);
        }

        let mut signature = original.clone();
        while expected.len() < signature.len()
            && signature.starts_with(STRUCT_SIG_START_CHAR)
            && signature.ends_with(STRUCT_SIG_END_CHAR)
        {
            signature = signature.slice(1..signature.len() - 1);
        }

        while signature.len() < expected.len()
            && expected.starts_with(STRUCT_SIG_START_CHAR)
            && expected.ends_with(STRUCT_SIG_END_CHAR)
        {
            expected = expected.slice(1..expected.len() - 1);
        }

        if signature == expected {
            Ok(PhantomData)
        } else {
            let expected = <T as Type>::signature();
            Err(zvariant::Error::SignatureMismatch(
                original.to_owned(),
                format!("`{}`", expected),
            ))
        }
    }
}

macro_rules! array_type {
    ($arr:ty) => {
        impl<T> Type for $arr
        where
            T: Type,
        {
            #[inline]
            fn signature() -> Signature<'static> {
                Signature::from_string_unchecked(format!("a{}", T::signature()))
            }
        }
    };
}

array_type!([T]);
array_type!(Vec<T>);

#[cfg(feature = "arrayvec")]
impl<A, T> Type for arrayvec::ArrayVec<A>
where
    A: arrayvec::Array<Item = T>,
    T: Type,
{
    #[inline]
    fn signature() -> Signature<'static> {
        <[T]>::signature()
    }
}

#[cfg(feature = "arrayvec")]
impl<A> Type for arrayvec::ArrayString<A>
where
    A: arrayvec::Array<Item = u8> + Copy,
{
    #[inline]
    fn signature() -> Signature<'static> {
        <&str>::signature()
    }
}

// Empty type deserves empty signature
impl Type for () {
    #[inline]
    fn signature() -> Signature<'static> {
        Signature::from_static_str_unchecked("")
    }
}

impl<T> Type for &T
where
    T: ?Sized + Type,
{
    #[inline]
    fn signature() -> Signature<'static> {
        T::signature()
    }
}

#[cfg(feature = "gvariant")]
impl<T> Type for Option<T>
where
    T: Type,
{
    #[inline]
    fn signature() -> Signature<'static> {
        Signature::from_string_unchecked(format!("m{}", T::signature()))
    }
}

////////////////////////////////////////////////////////////////////////////////

macro_rules! tuple_impls {
    ($($len:expr => ($($n:tt $name:ident)+))+) => {
        $(
            impl<$($name),+> Type for ($($name,)+)
            where
                $($name: Type,)+
            {
                #[inline]
                fn signature() -> Signature<'static> {
                    let mut sig = String::with_capacity(255);
                    sig.push(STRUCT_SIG_START_CHAR);
                    $(
                        sig.push_str($name::signature().as_str());
                    )+
                    sig.push(STRUCT_SIG_END_CHAR);

                    Signature::from_string_unchecked(sig)
                }
            }
        )+
    }
}

tuple_impls! {
    1 => (0 T0)
    2 => (0 T0 1 T1)
    3 => (0 T0 1 T1 2 T2)
    4 => (0 T0 1 T1 2 T2 3 T3)
    5 => (0 T0 1 T1 2 T2 3 T3 4 T4)
    6 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5)
    7 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6)
    8 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6 7 T7)
    9 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6 7 T7 8 T8)
    10 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6 7 T7 8 T8 9 T9)
    11 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6 7 T7 8 T8 9 T9 10 T10)
    12 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6 7 T7 8 T8 9 T9 10 T10 11 T11)
    13 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6 7 T7 8 T8 9 T9 10 T10 11 T11 12 T12)
    14 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6 7 T7 8 T8 9 T9 10 T10 11 T11 12 T12 13 T13)
    15 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6 7 T7 8 T8 9 T9 10 T10 11 T11 12 T12 13 T13 14 T14)
    16 => (0 T0 1 T1 2 T2 3 T3 4 T4 5 T5 6 T6 7 T7 8 T8 9 T9 10 T10 11 T11 12 T12 13 T13 14 T14 15 T15)
}

////////////////////////////////////////////////////////////////////////////////

// Arrays are serialized as tuples/structs by Serde so we treat them as such too even though
// it's very strange. Slices and arrayvec::ArrayVec can be used anyway so I guess it's no big
// deal.
// TODO: Mention this fact in the module docs.

macro_rules! array_impls {
    ($($len:tt)+) => {
        $(
            impl<T> Type for [T; $len]
            where
                T: Type,
            {
                #[inline]
                #[allow(clippy::reversed_empty_ranges)]
                fn signature() -> Signature<'static> {
                    let mut sig = String::with_capacity(255);
                    sig.push(STRUCT_SIG_START_CHAR);
                    if $len > 0 {
                        for _ in 0..$len {
                            sig.push_str(T::signature().as_str());
                        }
                    }
                    sig.push(STRUCT_SIG_END_CHAR);

                    Signature::from_string_unchecked(sig)
                }
            }
        )+
    }
}

array_impls! {
    0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32
}

////////////////////////////////////////////////////////////////////////////////

use std::{
    borrow::Cow,
    collections::{BTreeMap, HashMap},
    hash::{BuildHasher, Hash},
};

macro_rules! map_impl {
    ($ty:ident < K $(: $kbound1:ident $(+ $kbound2:ident)*)*, V $(, $typaram:ident : $bound:ident)* >) => {
        impl<K, V $(, $typaram)*> Type for $ty<K, V $(, $typaram)*>
        where
            K: Type $(+ $kbound1 $(+ $kbound2)*)*,
            V: Type,
            $($typaram: $bound,)*
        {
            #[inline]
            fn signature() -> Signature<'static> {
                Signature::from_string_unchecked(format!("a{{{}{}}}", K::signature(), V::signature()))
            }
        }
    }
}

map_impl!(BTreeMap<K: Ord, V>);
map_impl!(HashMap<K: Eq + Hash, V, H: BuildHasher>);

impl<T> Type for Cow<'_, T>
where
    T: ?Sized + Type + ToOwned,
{
    #[inline]
    fn signature() -> Signature<'static> {
        T::signature()
    }
}

// BitFlags
#[cfg(feature = "enumflags2")]
impl<F> Type for enumflags2::BitFlags<F>
where
    F: Type + enumflags2::RawBitFlags,
{
    #[inline]
    fn signature() -> Signature<'static> {
        F::signature()
    }
}

#[cfg(feature = "serde_bytes")]
impl Type for serde_bytes::Bytes {
    fn signature() -> Signature<'static> {
        Signature::from_static_str_unchecked("ay")
    }
}

#[cfg(feature = "serde_bytes")]
impl Type for serde_bytes::ByteBuf {
    fn signature() -> Signature<'static> {
        Signature::from_static_str_unchecked("ay")
    }
}

// TODO: Blanket implementation for more types: https://github.com/serde-rs/serde/blob/master/serde/src/ser/impls.rs
