use std::convert::TryFrom;
use std::marker::PhantomData;
use std::str::FromStr;

use serde::{de, de::Visitor, Deserializer};

pub(crate) fn u8_or_string<'de, D>(deserializer: D) -> Result<u8, D::Error>
where
    D: Deserializer<'de>,
{
    option_u8_or_string(deserializer).and_then(|i| {
        if let Some(i) = i {
            Ok(i)
        } else {
            Err(de::Error::custom("Required filed undefined"))
        }
    })
}

pub(crate) fn u16_or_string<'de, D>(deserializer: D) -> Result<u16, D::Error>
where
    D: Deserializer<'de>,
{
    option_u16_or_string(deserializer).and_then(|i| {
        if let Some(i) = i {
            Ok(i)
        } else {
            Err(de::Error::custom("Required filed undefined"))
        }
    })
}

pub(crate) fn u32_or_string<'de, D>(deserializer: D) -> Result<u32, D::Error>
where
    D: Deserializer<'de>,
{
    option_u32_or_string(deserializer).and_then(|i| {
        if let Some(i) = i {
            Ok(i)
        } else {
            Err(de::Error::custom("Required filed undefined"))
        }
    })
}

pub(crate) fn bool_or_string<'de, D>(deserializer: D) -> Result<bool, D::Error>
where
    D: Deserializer<'de>,
{
    option_bool_or_string(deserializer).and_then(|i| {
        if let Some(i) = i {
            Ok(i)
        } else {
            Err(de::Error::custom("Required filed undefined"))
        }
    })
}

pub(crate) fn option_bool_or_string<'de, D>(
    deserializer: D,
) -> Result<Option<bool>, D::Error>
where
    D: Deserializer<'de>,
{
    struct IntegerOrString(PhantomData<fn() -> Option<bool>>);

    impl<'de> Visitor<'de> for IntegerOrString {
        type Value = Option<bool>;

        fn expecting(
            &self,
            formatter: &mut std::fmt::Formatter,
        ) -> std::fmt::Result {
            formatter.write_str(
                "Need to be boolean: 1|0|true|false|yes|no|on|off|y|n",
            )
        }

        fn visit_bool<E>(self, value: bool) -> Result<Option<bool>, E>
        where
            E: de::Error,
        {
            Ok(Some(value))
        }

        fn visit_str<E>(self, value: &str) -> Result<Option<bool>, E>
        where
            E: de::Error,
        {
            match value.to_lowercase().as_str() {
                "1" | "true" | "yes" | "on" | "y" => Ok(Some(true)),
                "0" | "false" | "no" | "off" | "n" => Ok(Some(false)),
                _ => Err(de::Error::custom(
                    "Need to be boolean: 1|0|true|false|yes|no|on|off|y|n",
                )),
            }
        }

        fn visit_u64<E>(self, value: u64) -> Result<Option<bool>, E>
        where
            E: de::Error,
        {
            match value {
                1 => Ok(Some(true)),
                0 => Ok(Some(false)),
                _ => Err(de::Error::custom(
                    "Need to be boolean: 1|0|true|false|yes|no|on|off|y|n",
                )),
            }
        }
    }

    deserializer.deserialize_any(IntegerOrString(PhantomData))
}

pub(crate) fn option_u8_or_string<'de, D>(
    deserializer: D,
) -> Result<Option<u8>, D::Error>
where
    D: Deserializer<'de>,
{
    option_u64_or_string(deserializer).and_then(|i| {
        if let Some(i) = i {
            match u8::try_from(i) {
                Ok(i) => Ok(Some(i)),
                Err(e) => Err(de::Error::custom(e)),
            }
        } else {
            Ok(None)
        }
    })
}

pub(crate) fn option_u16_or_string<'de, D>(
    deserializer: D,
) -> Result<Option<u16>, D::Error>
where
    D: Deserializer<'de>,
{
    option_u64_or_string(deserializer).and_then(|i| {
        if let Some(i) = i {
            match u16::try_from(i) {
                Ok(i) => Ok(Some(i)),
                Err(e) => Err(de::Error::custom(e)),
            }
        } else {
            Ok(None)
        }
    })
}

pub(crate) fn option_u32_or_string<'de, D>(
    deserializer: D,
) -> Result<Option<u32>, D::Error>
where
    D: Deserializer<'de>,
{
    option_u64_or_string(deserializer).and_then(|i| {
        if let Some(i) = i {
            match u32::try_from(i) {
                Ok(i) => Ok(Some(i)),
                Err(e) => Err(de::Error::custom(e)),
            }
        } else {
            Ok(None)
        }
    })
}

// This function is inspired by https://serde.rs/string-or-struct.html
pub(crate) fn option_u64_or_string<'de, D>(
    deserializer: D,
) -> Result<Option<u64>, D::Error>
where
    D: Deserializer<'de>,
{
    struct IntegerOrString(PhantomData<fn() -> Option<u64>>);

    impl<'de> Visitor<'de> for IntegerOrString {
        type Value = Option<u64>;

        fn expecting(
            &self,
            formatter: &mut std::fmt::Formatter,
        ) -> std::fmt::Result {
            formatter.write_str("unsigned integer or string")
        }

        fn visit_str<E>(self, value: &str) -> Result<Option<u64>, E>
        where
            E: de::Error,
        {
            if let Some(prefix_len) = value.strip_prefix("0x") {
                u64::from_str_radix(prefix_len, 16)
                    .map_err(de::Error::custom)
                    .map(Some)
            } else {
                FromStr::from_str(value)
                    .map_err(de::Error::custom)
                    .map(Some)
            }
        }

        fn visit_u64<E>(self, value: u64) -> Result<Option<u64>, E>
        where
            E: de::Error,
        {
            Ok(Some(value))
        }
    }

    deserializer.deserialize_any(IntegerOrString(PhantomData))
}

pub(crate) fn option_i64_or_string<'de, D>(
    deserializer: D,
) -> Result<Option<i64>, D::Error>
where
    D: Deserializer<'de>,
{
    struct IntegerOrString(PhantomData<fn() -> Option<i64>>);

    impl<'de> Visitor<'de> for IntegerOrString {
        type Value = Option<i64>;

        fn expecting(
            &self,
            formatter: &mut std::fmt::Formatter,
        ) -> std::fmt::Result {
            formatter.write_str("signed integer or string")
        }

        fn visit_str<E>(self, value: &str) -> Result<Option<i64>, E>
        where
            E: de::Error,
        {
            FromStr::from_str(value)
                .map_err(de::Error::custom)
                .map(Some)
        }

        fn visit_u64<E>(self, value: u64) -> Result<Option<i64>, E>
        where
            E: de::Error,
        {
            i64::try_from(value).map_err(de::Error::custom).map(Some)
        }

        fn visit_i64<E>(self, value: i64) -> Result<Option<i64>, E>
        where
            E: de::Error,
        {
            Ok(Some(value))
        }
    }

    deserializer.deserialize_any(IntegerOrString(PhantomData))
}
