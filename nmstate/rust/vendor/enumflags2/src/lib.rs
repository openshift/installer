//! # Enum Flags
//! `enumflags2` defines a `BitFlags<T>` type, which is a `Set<T>`
//! for enums without associated data.
//!
//! ## Example
//! ```
//! use enumflags2::BitFlags;
//!
//! #[derive(BitFlags, Copy, Clone, Debug, PartialEq)]
//! #[repr(u8)]
//! enum Test {
//!     A = 0b0001,
//!     B = 0b0010,
//!     C = 0b0100,
//!     D = 0b1000,
//! }
//!
//! let a_b = Test::A | Test::B; // BitFlags<Test>
//! let a_c = Test::A | Test::C;
//! let b_c_d = Test::C | Test::B | Test::D;
//!
//! // BitFlags<Test>(0b11, [A, B])
//! println!("{:?}", a_b);
//!
//! // BitFlags<Test>(0b1, [A])
//! println!("{:?}", a_b & a_c);
//!
//! // Iterate over the flags like a normal set!
//! assert_eq!(a_b.iter().collect::<Vec<_>>(), &[Test::A, Test::B]);
//!
//! assert!(a_b.contains(Test::A));
//! assert!(b_c_d.contains(Test::B | Test::C));
//! assert!(!(b_c_d.contains(a_b)));
//!
//! assert!(a_b.intersects(a_c));
//! assert!(!(a_b.intersects(Test::C | Test::D)));
//! ```
//!
//! ## Note
//!
//! By default, the `BitFlags` are `usize`-sized. If you want them to be smaller,
//! specify a `repr` on your enum as in the example above.
//!
//! ## Optional Feature Flags
//!
//! - [`serde`](https://serde.rs/) implements `Serialize` and `Deserialize`
//!   for `BitFlags<T>`.
//! - `std` implements `std::error::Error` for `FromBitsError`.
//! - `not_literal` enables a workaround that allows using discriminant
//!   expressions that can't be evaluated at macro expansion time. Notably,
//!   this includes using pre-existing constants.
//!
//!   This is disabled by default because of the high potential for confusing
//!   error messages - if a flag doesn't have exactly one bit set, the error
//!   message will be "attempt to subtract with overflow", pointing at the
//!   relevant flag.
//!
//! ### Migrating from 0.5
//!
//! The minimum rustc version has been bumped to 1.34.0, because of `syn 1.0`. The
//! version policy from now on will be "what's available on Debian stable", [because
//! Debian is famously slow with new software versions][debian-snailpace].
//!
//! You should no longer depend on `enumflags2_derive` directly.
//! Use the reexport from the `enumflags2` crate.
//! semver guarantees will be violated if you depend on the derive crate directly.
//!
//! The derive macro has been renamed to `BitFlags`, to make it clearer what the
//! derive does.
//!
//! The `nostd` feature flag has been removed. The crate now only depends on `libcore`
//! by default. Enable the `std` flag to get an implementation of `std::error::Error`
//! on error types.
//!
//! Flags more than one bit set have been found to have inconsistent semantics.
//! They are now rejected at compile-time. The same applies to flags without any
//! bit set. If you were relying on this in your code, please [open an issue][issue]
//! and explain your usecase.
//!
//! `BitFlags::from_bits` returns a `Result` instead of an `Option`. This might
//! necessitate some minor changes in your code.
//!
//! `BitFlags::not` has been removed. Use the `!` operator instead.
//!
//! [debian-snailpace]: https://www.jwz.org/blog/2016/04/i-would-like-debian-to-stop-shipping-xscreensaver/
//! [issue]: https://github.com/NieDzejkob/enumflags2/issues/new
#![warn(missing_docs)]
#![cfg_attr(all(not(test), not(feature = "std")), no_std)]

use core::{cmp, ops};
use core::iter::FromIterator;

#[allow(unused_imports)]
#[macro_use]
extern crate enumflags2_derive;

#[doc(hidden)]
pub use enumflags2_derive::BitFlags_internal as BitFlags;

/// A trait automatically implemented by `derive(BitFlags)` to make the enum
/// a valid type parameter for `BitFlags<T>`.
pub trait RawBitFlags: Copy + Clone + 'static + _internal::RawBitFlags {}

/// While the module is public, this is only the case because it needs to be
/// accessed by the derive macro. Do not use this directly. Stability guarantees
/// don't apply.
#[doc(hidden)]
pub mod _internal {
    /// A trait automatically implemented by `derive(BitFlags)` to make the enum
    /// a valid type parameter for `BitFlags<T>`.
    pub trait RawBitFlags: Copy + Clone + 'static {
        /// The underlying integer type.
        type Type: BitFlagNum;

        /// Return a value with all flag bits set.
        fn all() -> Self::Type;

        /// Return the bits as a number type.
        fn bits(self) -> Self::Type;

        /// Return a slice that contains each variant exactly one.
        fn flag_list() -> &'static [Self];

        /// Return the name of the type for debug formatting purposes.
        ///
        /// This is typically `BitFlags<EnumName>`
        fn bitflags_type_name() -> &'static str {
            "BitFlags"
        }
    }

    use ::core::ops::{BitAnd, BitOr, BitXor, Not};
    use ::core::cmp::PartialOrd;
    use ::core::fmt;

    pub trait BitFlagNum
        : Default
        + BitOr<Self, Output = Self>
        + BitAnd<Self, Output = Self>
        + BitXor<Self, Output = Self>
        + Not<Output = Self>
        + PartialOrd<Self>
        + fmt::Debug
        + fmt::Binary
        + Copy
        + Clone {
    }

    impl BitFlagNum for u8 {}
    impl BitFlagNum for u16 {}
    impl BitFlagNum for u32 {}
    impl BitFlagNum for u64 {}
    impl BitFlagNum for usize {}

    // Re-export libcore so the macro doesn't inject "extern crate" downstream.
    pub mod core {
        pub use core::{convert, option, ops};
    }
}

// Internal debug formatting implementations
mod formatting;

// impl TryFrom<T::Type> for BitFlags<T>
mod fallible;
pub use crate::fallible::FromBitsError;

/// Represents a set of flags of some type `T`.
/// The type must have the `#[derive(BitFlags)]` attribute applied.
#[derive(Copy, Clone, Eq, Hash)]
#[repr(transparent)]
pub struct BitFlags<T: RawBitFlags> {
    val: T::Type,
}

/// The default value returned is one with all flags unset, i. e. [`empty`][Self::empty].
impl<T> Default for BitFlags<T>
where
    T: RawBitFlags,
{
    fn default() -> Self {
        Self::empty()
    }
}

impl<T: RawBitFlags> From<T> for BitFlags<T> {
    fn from(t: T) -> BitFlags<T> {
        Self::from_flag(t)
    }
}

impl<T> BitFlags<T>
where
    T: RawBitFlags,
{
    /// Create a new BitFlags unsafely, without checking if the bits form
    /// a valid bit pattern for the type.
    ///
    /// Consider using `from_bits` or `from_bits_truncate` instead.
    ///
    /// # Safety
    ///
    /// The argument must not have set bits at positions not corresponding to
    /// any flag.
    pub unsafe fn new(val: T::Type) -> Self {
        BitFlags { val }
    }

    /// Create an empty BitFlags. Empty means `0`.
    pub fn empty() -> Self {
        unsafe { BitFlags::new(T::Type::default()) }
    }

    /// Create a BitFlags with all flags set.
    pub fn all() -> Self {
        unsafe { BitFlags::new(T::all()) }
    }

    /// Returns true if all flags are set
    pub fn is_all(self) -> bool {
        self.val == T::all()
    }

    /// Returns true if no flag is set
    pub fn is_empty(self) -> bool {
        self.val == Self::empty().bits()
    }

    /// Returns the underlying type value
    pub fn bits(self) -> T::Type {
        self.val
    }

    /// Returns true if at least one flag is shared.
    pub fn intersects<B: Into<BitFlags<T>>>(self, other: B) -> bool {
        (self.bits() & other.into().bits()) > Self::empty().bits()
    }

    /// Returns true if all flags are contained.
    pub fn contains<B: Into<BitFlags<T>>>(self, other: B) -> bool {
        let other = other.into();
        (self.bits() & other.bits()) == other.bits()
    }

    /// Returns a `BitFlags<T>` if the raw value provided does not contain
    /// any illegal flags.
    pub fn from_bits(bits: T::Type) -> Result<Self, FromBitsError<T>> {
        let flags = Self::from_bits_truncate(bits);
        if flags.bits() == bits {
            Ok(flags)
        } else {
            Err(FromBitsError {
                flags,
                invalid: bits & !flags.bits(),
            })
        }
    }

    /// Turn a `T` into a `BitFlags<T>`. Also available as `flag.into()`.
    pub fn from_flag(flag: T) -> Self {
        BitFlags { val: flag.bits() }
    }

    /// Truncates flags that are illegal
    pub fn from_bits_truncate(bits: T::Type) -> Self {
        unsafe { BitFlags::new(bits & T::all()) }
    }

    /// Toggles the matching bits
    pub fn toggle<B: Into<BitFlags<T>>>(&mut self, other: B) {
        *self ^= other.into();
    }

    /// Inserts the flags into the BitFlag
    pub fn insert<B: Into<BitFlags<T>>>(&mut self, other: B) {
        *self |= other.into();
    }

    /// Removes the matching flags
    pub fn remove<B: Into<BitFlags<T>>>(&mut self, other: B) {
        *self &= !other.into();
    }

    /// Returns an iterator that yields each set flag
    pub fn iter(self) -> impl Iterator<Item = T> {
        T::flag_list().iter().cloned().filter(move |&flag| self.contains(flag))
    }
}

impl<T, B> cmp::PartialEq<B> for BitFlags<T>
where
    T: RawBitFlags,
    B: Into<BitFlags<T>> + Copy,
{
    fn eq(&self, other: &B) -> bool {
        self.bits() == Into::<Self>::into(*other).bits()
    }
}

impl<T, B> ops::BitOr<B> for BitFlags<T>
where
    T: RawBitFlags,
    B: Into<BitFlags<T>>,
{
    type Output = BitFlags<T>;
    fn bitor(self, other: B) -> BitFlags<T> {
        unsafe { BitFlags::new(self.bits() | other.into().bits()) }
    }
}

impl<T, B> ops::BitAnd<B> for BitFlags<T>
where
    T: RawBitFlags,
    B: Into<BitFlags<T>>,
{
    type Output = BitFlags<T>;
    fn bitand(self, other: B) -> BitFlags<T> {
        unsafe { BitFlags::new(self.bits() & other.into().bits()) }
    }
}

impl<T, B> ops::BitXor<B> for BitFlags<T>
where
    T: RawBitFlags,
    B: Into<BitFlags<T>>,
{
    type Output = BitFlags<T>;
    fn bitxor(self, other: B) -> BitFlags<T> {
        unsafe { BitFlags::new((self.bits() ^ other.into().bits()) & T::all()) }
    }
}

impl<T, B> ops::BitOrAssign<B> for BitFlags<T>
where
    T: RawBitFlags,
    B: Into<BitFlags<T>>,
{
    fn bitor_assign(&mut self, other: B) {
        *self = *self | other;
    }
}

impl<T, B> ops::BitAndAssign<B> for BitFlags<T>
where
    T: RawBitFlags,
    B: Into<BitFlags<T>>,
{
    fn bitand_assign(&mut self, other: B) {
        *self = *self & other;
    }
}
impl<T, B> ops::BitXorAssign<B> for BitFlags<T>
where
    T: RawBitFlags,
    B: Into<BitFlags<T>>,
{
    fn bitxor_assign(&mut self, other: B) {
        *self = *self ^ other;
    }
}

impl<T> ops::Not for BitFlags<T>
where
    T: RawBitFlags,
{
    type Output = BitFlags<T>;
    fn not(self) -> BitFlags<T> {
        unsafe { BitFlags::new(!self.bits() & T::all()) }
    }
}

impl<T, B> FromIterator<B> for BitFlags<T>
where
    T: RawBitFlags,
    B: Into<BitFlags<T>>
{
    fn from_iter<I>(it: I) -> BitFlags<T>
    where
        I: IntoIterator<Item = B>
    {
        it.into_iter().fold(BitFlags::empty(), |acc, flag| acc | flag)
    }
}

impl<T, B> Extend<B> for BitFlags<T>
where
    T: RawBitFlags,
    B: Into<BitFlags<T>>
{
    fn extend<I>(&mut self, it: I)
    where
        I: IntoIterator<Item = B>
    {
        *self = it.into_iter().fold(*self, |acc, flag| acc | flag)
    }
}

#[cfg(feature = "serde")]
mod impl_serde {
    use serde::{Serialize, Deserialize};
    use serde::de::{Error, Unexpected};
    use super::{BitFlags, RawBitFlags};

    impl<'a, T> Deserialize<'a> for BitFlags<T>
    where
        T: RawBitFlags,
        T::Type: Deserialize<'a> + Into<u64>,
    {
        fn deserialize<D: serde::Deserializer<'a>>(d: D) -> Result<Self, D::Error> {
            let val = T::Type::deserialize(d)?;
            Self::from_bits(val)
                .or_else(|_| Err(D::Error::invalid_value(
                    Unexpected::Unsigned(val.into()),
                    &"valid bit representation"
                )))
        }
    }

    impl<T> Serialize for BitFlags<T>
    where
        T: RawBitFlags,
        T::Type: Serialize,
    {
        fn serialize<S: serde::Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
            T::Type::serialize(&self.val, s)
        }
    }
}
