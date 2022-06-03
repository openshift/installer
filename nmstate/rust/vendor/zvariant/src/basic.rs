// FIXME: Drop this when the deprecated `Basic::ALIGNMENT` is dropped in the next API break.
#![allow(deprecated)]

use crate::{EncodingFormat, Signature, Type};

/// Trait for basic types.
///
/// All basic types are also [`Type`] implementers.
///
/// [`Type`]: trait.Type.html
/// [`Value`]: enum.Value.html
pub trait Basic: Type {
    /// The type signature, as a character.
    const SIGNATURE_CHAR: char;
    /// The type signature, as a string.
    const SIGNATURE_STR: &'static str;
    #[deprecated(since = "2.0.2", note = "Please use the `alignment` function instead")]
    /// The required padding alignment.
    const ALIGNMENT: usize;

    /// The required padding alignment for the given format.
    fn alignment(format: EncodingFormat) -> usize;
}

impl<B: ?Sized> Basic for &B
where
    B: Basic,
{
    const SIGNATURE_CHAR: char = B::SIGNATURE_CHAR;
    const SIGNATURE_STR: &'static str = B::SIGNATURE_STR;
    const ALIGNMENT: usize = B::ALIGNMENT;

    fn alignment(format: EncodingFormat) -> usize {
        B::alignment(format)
    }
}

macro_rules! impl_type {
    ($for:ty) => {
        impl Type for $for {
            fn signature() -> Signature<'static> {
                Signature::from_static_str_unchecked(<$for>::SIGNATURE_STR)
            }
        }
    };
}

macro_rules! alignment_method {
    ($alignment:expr) => {
        alignment_method!($alignment, $alignment);
    };
    ($dbus_alignment:expr, $gvariant_alignment:expr) => {
        fn alignment(format: EncodingFormat) -> usize {
            match format {
                EncodingFormat::DBus => $dbus_alignment,
                #[cfg(feature = "gvariant")]
                EncodingFormat::GVariant => $gvariant_alignment,
            }
        }
    };
}

impl Basic for u8 {
    const SIGNATURE_CHAR: char = 'y';
    const SIGNATURE_STR: &'static str = "y";
    const ALIGNMENT: usize = 1;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(u8);

// No i8 type in D-Bus/GVariant, let's pretend it's i16
impl Basic for i8 {
    const SIGNATURE_CHAR: char = i16::SIGNATURE_CHAR;
    const SIGNATURE_STR: &'static str = i16::SIGNATURE_STR;
    const ALIGNMENT: usize = i16::ALIGNMENT;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(i8);

impl Basic for bool {
    const SIGNATURE_CHAR: char = 'b';
    const SIGNATURE_STR: &'static str = "b";
    const ALIGNMENT: usize = 4;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(bool);

impl Basic for i16 {
    const SIGNATURE_CHAR: char = 'n';
    const SIGNATURE_STR: &'static str = "n";
    const ALIGNMENT: usize = 2;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(i16);

impl Basic for u16 {
    const SIGNATURE_CHAR: char = 'q';
    const SIGNATURE_STR: &'static str = "q";
    const ALIGNMENT: usize = 2;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(u16);

impl Basic for i32 {
    const SIGNATURE_CHAR: char = 'i';
    const SIGNATURE_STR: &'static str = "i";
    const ALIGNMENT: usize = 4;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(i32);

impl Basic for u32 {
    const SIGNATURE_CHAR: char = 'u';
    const SIGNATURE_STR: &'static str = "u";
    const ALIGNMENT: usize = 4;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(u32);

impl Basic for i64 {
    const SIGNATURE_CHAR: char = 'x';
    const SIGNATURE_STR: &'static str = "x";
    const ALIGNMENT: usize = 8;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(i64);

impl Basic for u64 {
    const SIGNATURE_CHAR: char = 't';
    const SIGNATURE_STR: &'static str = "t";
    const ALIGNMENT: usize = 8;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(u64);

// No f32 type in D-Bus/GVariant, let's pretend it's f64
impl Basic for f32 {
    const SIGNATURE_CHAR: char = f64::SIGNATURE_CHAR;
    const SIGNATURE_STR: &'static str = f64::SIGNATURE_STR;
    const ALIGNMENT: usize = f64::ALIGNMENT;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(f32);

impl Basic for f64 {
    const SIGNATURE_CHAR: char = 'd';
    const SIGNATURE_STR: &'static str = "d";
    const ALIGNMENT: usize = 8;

    alignment_method!(Self::ALIGNMENT);
}
impl_type!(f64);

impl Basic for str {
    const SIGNATURE_CHAR: char = 's';
    const SIGNATURE_STR: &'static str = "s";
    const ALIGNMENT: usize = 4;

    alignment_method!(Self::ALIGNMENT, 1);
}
impl_type!(str);

impl Basic for String {
    const SIGNATURE_CHAR: char = 's';
    const SIGNATURE_STR: &'static str = "s";
    const ALIGNMENT: usize = 4;

    alignment_method!(Self::ALIGNMENT, 1);
}
impl_type!(String);

impl Basic for char {
    const SIGNATURE_CHAR: char = <&str>::SIGNATURE_CHAR;
    const SIGNATURE_STR: &'static str = <&str>::SIGNATURE_STR;
    const ALIGNMENT: usize = <&str>::ALIGNMENT;

    alignment_method!(Self::ALIGNMENT, 1);
}
impl_type!(char);
