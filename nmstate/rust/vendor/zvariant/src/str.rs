// FIXME: Drop this when the deprecated `Basic::ALIGNMENT` is dropped in the next API break.
#![allow(deprecated)]

use serde::{Deserialize, Deserializer, Serialize, Serializer};
use static_assertions::assert_impl_all;
use std::hash::{Hash, Hasher};

use crate::{Basic, EncodingFormat, Signature, Type};

/// A string wrapper.
///
/// This is used for keeping strings in a [`Value`]. API is provided to convert from, and to a
/// [`&str`] and [`String`].
///
/// [`Value`]: enum.Value.html#variant.Str
/// [`&str`]: https://doc.rust-lang.org/std/str/index.html
/// [`String`]: https://doc.rust-lang.org/std/string/struct.String.html
#[derive(Debug, Default, PartialEq, Eq, Hash, Clone, Serialize, Deserialize)]
#[serde(rename(serialize = "zvariant::Str", deserialize = "zvariant::Str"))]
pub struct Str<'a>(#[serde(borrow)] Inner<'a>);

#[derive(Debug, Eq, Clone)]
enum Inner<'a> {
    Static(&'static str),
    Borrowed(&'a str),
    Owned(Box<str>),
}

impl<'a> Default for Inner<'a> {
    fn default() -> Self {
        Self::Static("")
    }
}

impl<'a> PartialEq for Inner<'a> {
    fn eq(&self, other: &Inner<'a>) -> bool {
        self.as_str() == other.as_str()
    }
}

impl<'a> Hash for Inner<'a> {
    fn hash<H: Hasher>(&self, h: &mut H) {
        self.as_str().hash(h)
    }
}

impl<'a> Inner<'a> {
    /// The underlying string.
    pub fn as_str(&self) -> &str {
        match self {
            Inner::Static(s) => s,
            Inner::Borrowed(s) => s,
            Inner::Owned(s) => s,
        }
    }
}

impl<'a> Serialize for Inner<'a> {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        s.serialize_str(self.as_str())
    }
}

impl<'de: 'a, 'a> Deserialize<'de> for Inner<'a> {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        <&'a str>::deserialize(deserializer).map(Inner::Borrowed)
    }
}

assert_impl_all!(Str<'_>: Send, Sync, Unpin);

impl<'a> Str<'a> {
    /// An owned string without allocations
    pub fn from_static(s: &'static str) -> Self {
        Str(Inner::Static(s))
    }

    /// A borrowed clone (this never allocates, unlike clone).
    pub fn as_ref(&self) -> Str<'_> {
        match &self.0 {
            Inner::Static(s) => Str(Inner::Static(s)),
            Inner::Borrowed(s) => Str(Inner::Borrowed(s)),
            Inner::Owned(s) => Str(Inner::Borrowed(s)),
        }
    }

    /// The underlying string.
    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }

    /// Creates an owned clone of `self`.
    pub fn to_owned(&self) -> Str<'static> {
        self.clone().into_owned()
    }

    /// Creates an owned clone of `self`.
    pub fn into_owned(self) -> Str<'static> {
        match self.0 {
            Inner::Static(s) => Str(Inner::Static(s)),
            Inner::Borrowed(s) => Str(Inner::Owned(s.into())),
            Inner::Owned(s) => Str(Inner::Owned(s)),
        }
    }
}

impl<'a> Basic for Str<'a> {
    const SIGNATURE_CHAR: char = <&str>::SIGNATURE_CHAR;
    const SIGNATURE_STR: &'static str = <&str>::SIGNATURE_STR;
    const ALIGNMENT: usize = <&str>::ALIGNMENT;

    fn alignment(format: EncodingFormat) -> usize {
        <&str>::alignment(format)
    }
}

impl<'a> Type for Str<'a> {
    fn signature() -> Signature<'static> {
        Signature::from_static_str_unchecked(Self::SIGNATURE_STR)
    }
}

impl<'a> From<&'a str> for Str<'a> {
    fn from(value: &'a str) -> Self {
        Self(Inner::Borrowed(value))
    }
}

impl<'a> From<&'a String> for Str<'a> {
    fn from(value: &'a String) -> Self {
        Self(Inner::Borrowed(value))
    }
}

impl<'a> From<String> for Str<'a> {
    fn from(value: String) -> Self {
        Self(Inner::Owned(value.into()))
    }
}

impl<'a> From<Str<'a>> for String {
    fn from(value: Str<'a>) -> String {
        match value.0 {
            Inner::Static(s) => s.into(),
            Inner::Borrowed(s) => s.into(),
            Inner::Owned(s) => s.into(),
        }
    }
}

impl<'a> From<&'a Str<'a>> for &'a str {
    fn from(value: &'a Str<'a>) -> &'a str {
        value.as_str()
    }
}

impl<'a> std::ops::Deref for Str<'a> {
    type Target = str;

    fn deref(&self) -> &Self::Target {
        self.as_str()
    }
}

impl<'a> PartialEq<str> for Str<'a> {
    fn eq(&self, other: &str) -> bool {
        self.as_str() == other
    }
}

impl<'a> PartialEq<&str> for Str<'a> {
    fn eq(&self, other: &&str) -> bool {
        self.as_str() == *other
    }
}

impl<'a> std::fmt::Display for Str<'a> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        self.as_str().fmt(f)
    }
}

#[cfg(test)]
mod tests {
    use super::Str;

    #[test]
    fn from_string() {
        let string = String::from("value");
        let v = Str::from(&string);
        assert_eq!(v.as_str(), "value");
    }
}
