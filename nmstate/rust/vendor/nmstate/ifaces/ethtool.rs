use std::collections::HashMap;
use std::marker::PhantomData;

use serde::{
    de, de::MapAccess, de::Visitor, Deserialize, Deserializer, Serialize,
};

const ETHTOOL_FEATURE_CLI_ALIAS: [(&str, &str); 17] = [
    ("rx", "rx-checksum"),
    ("rx-checksumming", "rx-checksum"),
    ("ufo", "tx-udp-fragmentation"),
    ("gso", "tx-generic-segmentation"),
    ("generic-segmentation-offload", "tx-generic-segmentation"),
    ("gro", "rx-gro"),
    ("generic-receive-offload", "rx-gro"),
    ("lro", "rx-lro"),
    ("large-receive-offload", "rx-lro"),
    ("rxvlan", "rx-vlan-hw-parse"),
    ("rx-vlan-offload", "rx-vlan-hw-parse"),
    ("txvlan", "tx-vlan-hw-insert"),
    ("tx-vlan-offload", "tx-vlan-hw-insert"),
    ("ntuple", "rx-ntuple-filter"),
    ("ntuple-filters", "rx-ntuple-filter"),
    ("rxhash", "rx-hashing"),
    ("receive-hashing", "rx-hashing"),
];

pub type EthtoolFeatureConfig = HashMap<String, bool>;

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Clone, Default)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct EthtoolConfig {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub pause: Option<EthtoolPauseConfig>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "parse_ethtool_feature"
    )]
    pub feature: Option<EthtoolFeatureConfig>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub coalesce: Option<EthtoolCoalesceConfig>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ring: Option<EthtoolRingConfig>,
}

impl EthtoolConfig {
    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn pre_edit_cleanup(&mut self) {
        self.pre_verify_cleanup();
    }

    // There are some alias on ethtool features.
    pub(crate) fn pre_verify_cleanup(&mut self) {
        if let Some(features) = self.feature.as_mut() {
            for (cli_alias, kernel_name) in ETHTOOL_FEATURE_CLI_ALIAS {
                if let Some(v) = features.remove(cli_alias) {
                    features.insert(kernel_name.to_string(), v);
                }
            }
        }
    }
}

#[derive(
    Serialize, Deserialize, Debug, Eq, PartialEq, Clone, Default, Copy,
)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct EthtoolPauseConfig {
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub rx: Option<bool>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub tx: Option<bool>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub autoneg: Option<bool>,
}

impl EthtoolPauseConfig {
    pub fn new() -> Self {
        Self::default()
    }
}

#[derive(
    Serialize, Deserialize, Debug, Eq, PartialEq, Clone, Default, Copy,
)]
#[serde(rename_all = "kebab-case", deny_unknown_fields)]
#[non_exhaustive]
pub struct EthtoolCoalesceConfig {
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub adaptive_rx: Option<bool>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub adaptive_tx: Option<bool>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub pkt_rate_high: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub pkt_rate_low: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_frames: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_frames_high: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_frames_irq: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_frames_low: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_usecs: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_usecs_high: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_usecs_irq: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_usecs_low: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub sample_interval: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub stats_block_usecs: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx_frames: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx_frames_high: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx_frames_irq: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx_frames_low: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx_usecs: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx_usecs_high: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx_usecs_irq: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx_usecs_low: Option<u32>,
}

impl EthtoolCoalesceConfig {
    pub fn new() -> Self {
        Self::default()
    }
}

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Clone, Default)]
#[serde(rename_all = "kebab-case", deny_unknown_fields)]
#[non_exhaustive]
pub struct EthtoolRingConfig {
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_max: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_jumbo: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_jumbo_max: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_mini: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub rx_mini_max: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx: Option<u32>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u32_or_string"
    )]
    pub tx_max: Option<u32>,
}

impl EthtoolRingConfig {
    pub fn new() -> Self {
        Self::default()
    }
}

fn parse_ethtool_feature<'de, D>(
    deserializer: D,
) -> Result<Option<EthtoolFeatureConfig>, D::Error>
where
    D: Deserializer<'de>,
{
    struct FeatureVisitor(PhantomData<fn() -> Option<EthtoolFeatureConfig>>);

    impl<'de> Visitor<'de> for FeatureVisitor {
        type Value = Option<EthtoolFeatureConfig>;

        fn expecting(
            &self,
            formatter: &mut std::fmt::Formatter,
        ) -> std::fmt::Result {
            formatter.write_str("Need to hash map of String:bool")
        }

        fn visit_map<M>(
            self,
            mut access: M,
        ) -> Result<Option<EthtoolFeatureConfig>, M::Error>
        where
            M: MapAccess<'de>,
        {
            let mut ret =
                HashMap::with_capacity(access.size_hint().unwrap_or(0));

            while let Some((key, value)) =
                access.next_entry::<String, serde_json::Value>()?
            {
                match value {
                    serde_json::Value::Bool(b) => {
                        ret.insert(key, b);
                    }
                    serde_json::Value::String(b)
                        if b.to_lowercase().as_str() == "false"
                            || b.as_str() == "0" =>
                    {
                        ret.insert(key, false);
                    }
                    serde_json::Value::String(b)
                        if b.to_lowercase().as_str() == "true"
                            || b.as_str() == "1" =>
                    {
                        ret.insert(key, true);
                    }
                    _ => {
                        return Err(de::Error::custom(
                            "Invalid feature value, should be boolean",
                        ));
                    }
                }
            }

            Ok(Some(ret))
        }
    }

    deserializer.deserialize_any(FeatureVisitor(PhantomData))
}
