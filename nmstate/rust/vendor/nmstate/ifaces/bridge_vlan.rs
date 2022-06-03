use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "kebab-case", deny_unknown_fields)]
#[non_exhaustive]
pub struct BridgePortVlanConfig {
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_bool_or_string"
    )]
    pub enable_native: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub mode: Option<BridgePortVlanMode>,
    #[serde(
        skip_serializing_if = "Option::is_none",
        default,
        deserialize_with = "crate::deserializer::option_u16_or_string"
    )]
    pub tag: Option<u16>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub trunk_tags: Option<Vec<BridgePortTunkTag>>,
}

impl BridgePortVlanConfig {
    pub fn new() -> Self {
        Self::default()
    }

    pub(crate) fn flatten_vlan_ranges(&mut self) {
        if let Some(trunk_tags) = &self.trunk_tags {
            let mut new_trunk_tags = Vec::new();
            for trunk_tag in trunk_tags {
                match trunk_tag {
                    BridgePortTunkTag::Id(_) => {
                        new_trunk_tags.push(trunk_tag.clone())
                    }
                    BridgePortTunkTag::IdRange(range) => {
                        for i in range.min..range.max + 1 {
                            new_trunk_tags.push(BridgePortTunkTag::Id(i));
                        }
                    }
                };
            }
            self.trunk_tags = Some(new_trunk_tags);
        }
    }

    pub(crate) fn sort_trunk_tags(&mut self) {
        if let Some(trunk_tags) = self.trunk_tags.as_mut() {
            trunk_tags.sort_unstable_by(|tag_a, tag_b| match (tag_a, tag_b) {
                (BridgePortTunkTag::Id(a), BridgePortTunkTag::Id(b)) => {
                    a.cmp(b)
                }
                _ => {
                    log::warn!(
                        "Please call flatten_vlan_ranges() \
                        before sort_port_vlans()"
                    );
                    std::cmp::Ordering::Equal
                }
            })
        }
    }

    pub(crate) fn is_changed(&self, current: &Self) -> bool {
        (self.enable_native.is_some()
            && self.enable_native != current.enable_native)
            || (self.mode.is_some() && self.mode != current.mode)
            || (self.tag.is_some() && self.tag != current.tag)
            || (self.trunk_tags.is_some()
                && self.trunk_tags != current.trunk_tags)
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum BridgePortVlanMode {
    Trunk,
    Access,
}

impl Default for BridgePortVlanMode {
    fn default() -> Self {
        Self::Access
    }
}

impl std::fmt::Display for BridgePortVlanMode {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Self::Trunk => "trunk",
                Self::Access => "access",
            }
        )
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "kebab-case")]
#[non_exhaustive]
pub enum BridgePortTunkTag {
    #[serde(deserialize_with = "crate::deserializer::u16_or_string")]
    Id(u16),
    IdRange(BridgePortVlanRange),
}

impl BridgePortTunkTag {
    pub fn get_vlan_tag_range(&self) -> (u16, u16) {
        match self {
            Self::Id(min) => (*min, *min),
            Self::IdRange(range) => (range.min, range.max),
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct BridgePortVlanRange {
    #[serde(deserialize_with = "crate::deserializer::u16_or_string")]
    pub max: u16,
    #[serde(deserialize_with = "crate::deserializer::u16_or_string")]
    pub min: u16,
}
