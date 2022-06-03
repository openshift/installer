use std::collections::HashMap;
use std::convert::TryFrom;

use serde::Deserialize;

use super::super::{connection::DbusDictionary, ErrorKind, NmError};

const VALID_FEATURES: [&str; 58] = [
    "feature-esp-hw-offload",
    "feature-esp-tx-csum-hw-offload",
    "feature-fcoe-mtu",
    "feature-gro",
    "feature-gso",
    "feature-highdma",
    "feature-hw-tc-offload",
    "feature-l2-fwd-offload",
    "feature-loopback",
    "feature-lro",
    "feature-macsec-hw-offload",
    "feature-ntuple",
    "feature-rx",
    "feature-rxhash",
    "feature-rxvlan",
    "feature-rx-all",
    "feature-rx-fcs",
    "feature-rx-gro-hw",
    "feature-rx-gro-list",
    "feature-rx-udp-gro-forwarding",
    "feature-rx-udp_tunnel-port-offload",
    "feature-rx-vlan-filter",
    "feature-rx-vlan-stag-filter",
    "feature-rx-vlan-stag-hw-parse",
    "feature-sg",
    "feature-tls-hw-record",
    "feature-tls-hw-rx-offload",
    "feature-tls-hw-tx-offload",
    "feature-tso",
    "feature-tx",
    "feature-txvlan",
    "feature-tx-checksum-fcoe-crc",
    "feature-tx-checksum-ipv4",
    "feature-tx-checksum-ipv6",
    "feature-tx-checksum-ip-generic",
    "feature-tx-checksum-sctp",
    "feature-tx-esp-segmentation",
    "feature-tx-fcoe-segmentation",
    "feature-tx-gre-csum-segmentation",
    "feature-tx-gre-segmentation",
    "feature-tx-gso-list",
    "feature-tx-gso-partial",
    "feature-tx-gso-robust",
    "feature-tx-ipxip4-segmentation",
    "feature-tx-ipxip6-segmentation",
    "feature-tx-nocache-copy",
    "feature-tx-scatter-gather",
    "feature-tx-scatter-gather-fraglist",
    "feature-tx-sctp-segmentation",
    "feature-tx-tcp6-segmentation",
    "feature-tx-tcp-ecn-segmentation",
    "feature-tx-tcp-mangleid-segmentation",
    "feature-tx-tcp-segmentation",
    "feature-tx-tunnel-remcsum-segmentation",
    "feature-tx-udp-segmentation",
    "feature-tx-udp_tnl-csum-segmentation",
    "feature-tx-udp_tnl-segmentation",
    "feature-tx-vlan-stag-hw-insert",
];

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
#[non_exhaustive]
pub struct NmSettingEthtool {
    pub pause_rx: Option<bool>,
    pub pause_tx: Option<bool>,
    pub pause_autoneg: Option<bool>,
    pub coalesce_adaptive_rx: Option<bool>,
    pub coalesce_adaptive_tx: Option<bool>,
    pub coalesce_pkt_rate_high: Option<u32>,
    pub coalesce_pkt_rate_low: Option<u32>,
    pub coalesce_rx_frames: Option<u32>,
    pub coalesce_rx_frames_high: Option<u32>,
    pub coalesce_rx_frames_low: Option<u32>,
    pub coalesce_rx_frames_irq: Option<u32>,
    pub coalesce_rx_usecs: Option<u32>,
    pub coalesce_rx_usecs_high: Option<u32>,
    pub coalesce_rx_usecs_low: Option<u32>,
    pub coalesce_rx_usecs_irq: Option<u32>,
    pub coalesce_sample_interval: Option<u32>,
    pub coalesce_stats_block_usecs: Option<u32>,
    pub coalesce_tx_frames: Option<u32>,
    pub coalesce_tx_frames_high: Option<u32>,
    pub coalesce_tx_frames_low: Option<u32>,
    pub coalesce_tx_frames_irq: Option<u32>,
    pub coalesce_tx_usecs: Option<u32>,
    pub coalesce_tx_usecs_high: Option<u32>,
    pub coalesce_tx_usecs_low: Option<u32>,
    pub coalesce_tx_usecs_irq: Option<u32>,
    pub features: Option<HashMap<String, bool>>,
    pub ring_rx: Option<u32>,
    pub ring_rx_jumbo: Option<u32>,
    pub ring_rx_mini: Option<u32>,
    pub ring_tx: Option<u32>,
    _other: HashMap<String, zvariant::OwnedValue>,
}

impl TryFrom<DbusDictionary> for NmSettingEthtool {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        let mut features: HashMap<String, bool> = HashMap::new();
        let feature_keys: Vec<String> = v
            .keys()
            .filter(|k| k.starts_with("feature-"))
            .cloned()
            .collect();
        for k in feature_keys {
            if let Some(feature_value) = v.remove(&k) {
                let value = bool::try_from(feature_value)?;
                features.insert(k, value);
            }
        }

        Ok(Self {
            pause_rx: _from_map!(v, "pause-rx", bool::try_from)?,
            pause_tx: _from_map!(v, "pause-tx", bool::try_from)?,
            pause_autoneg: _from_map!(v, "pause-autoneg", bool::try_from)?,
            coalesce_adaptive_rx: _from_map!(
                v,
                "coalesce-adaptive-rx",
                bool::try_from
            )?,
            coalesce_adaptive_tx: _from_map!(
                v,
                "coalesce-adaptive-tx",
                bool::try_from
            )?,
            coalesce_pkt_rate_high: _from_map!(
                v,
                "coalesce-pkt-rate-high",
                u32::try_from
            )?,
            coalesce_pkt_rate_low: _from_map!(
                v,
                "coalesce-pkt-rate-low",
                u32::try_from
            )?,
            coalesce_rx_frames: _from_map!(
                v,
                "coalesce-rx-frames",
                u32::try_from
            )?,
            coalesce_rx_frames_high: _from_map!(
                v,
                "coalesce-rx-frames-high",
                u32::try_from
            )?,
            coalesce_rx_frames_low: _from_map!(
                v,
                "coalesce-rx-frames-low",
                u32::try_from
            )?,
            coalesce_rx_frames_irq: _from_map!(
                v,
                "coalesce-rx-frames-irq",
                u32::try_from
            )?,
            coalesce_rx_usecs: _from_map!(
                v,
                "coalesce-rx-usecs",
                u32::try_from
            )?,
            coalesce_rx_usecs_high: _from_map!(
                v,
                "coalesce-rx-usecs-high",
                u32::try_from
            )?,
            coalesce_rx_usecs_low: _from_map!(
                v,
                "coalesce-rx-usecs-low",
                u32::try_from
            )?,
            coalesce_rx_usecs_irq: _from_map!(
                v,
                "coalesce-rx-usecs-irq",
                u32::try_from
            )?,
            coalesce_tx_frames: _from_map!(
                v,
                "coalesce-tx-frames",
                u32::try_from
            )?,
            coalesce_tx_frames_high: _from_map!(
                v,
                "coalesce-tx-frames-high",
                u32::try_from
            )?,
            coalesce_tx_frames_low: _from_map!(
                v,
                "coalesce-tx-frames-low",
                u32::try_from
            )?,
            coalesce_tx_frames_irq: _from_map!(
                v,
                "coalesce-tx-frames-irq",
                u32::try_from
            )?,
            coalesce_tx_usecs: _from_map!(
                v,
                "coalesce-tx-usecs",
                u32::try_from
            )?,
            coalesce_tx_usecs_high: _from_map!(
                v,
                "coalesce-tx-usecs-high",
                u32::try_from
            )?,
            coalesce_tx_usecs_low: _from_map!(
                v,
                "coalesce-tx-usecs-low",
                u32::try_from
            )?,
            coalesce_tx_usecs_irq: _from_map!(
                v,
                "coalesce-tx-usecs-irq",
                u32::try_from
            )?,
            coalesce_sample_interval: _from_map!(
                v,
                "coalesce-sample-interval",
                u32::try_from
            )?,
            coalesce_stats_block_usecs: _from_map!(
                v,
                "coalesce-stats-block-usecs",
                u32::try_from
            )?,
            features: Some(features),
            ring_rx: _from_map!(v, "ring-rx", u32::try_from)?,
            ring_rx_jumbo: _from_map!(v, "ring-rx-jumbo", u32::try_from)?,
            ring_rx_mini: _from_map!(v, "ring-rx-mini", u32::try_from)?,
            ring_tx: _from_map!(v, "ring-tx", u32::try_from)?,
            _other: v,
        })
    }
}

impl NmSettingEthtool {
    pub fn validate(&self) -> Result<(), NmError> {
        if let Some(features) = self.features.as_ref() {
            for k in features.keys() {
                if !VALID_FEATURES.contains(&k.as_str()) {
                    return Err(NmError::new(
                        ErrorKind::InvalidArgument,
                        format!("Unsupported ethtool feature {}", k),
                    ));
                }
            }
        }
        Ok(())
    }

    pub(crate) fn to_keyfile(
        &self,
    ) -> Result<HashMap<String, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();

        for (k, v) in self.to_value()?.drain() {
            ret.insert(k.to_string(), v);
        }

        Ok(ret)
    }

    pub(crate) fn to_value(
        &self,
    ) -> Result<HashMap<&str, zvariant::Value>, NmError> {
        let mut ret = HashMap::new();
        if let Some(v) = &self.pause_rx {
            ret.insert("pause-rx", zvariant::Value::new(v));
        }
        if let Some(v) = &self.pause_tx {
            ret.insert("pause-tx", zvariant::Value::new(v));
        }
        if let Some(v) = &self.pause_autoneg {
            ret.insert("pause-autoneg", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_adaptive_rx {
            ret.insert("coalesce-adaptive-rx", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_adaptive_tx {
            ret.insert("coalesce-adaptive-tx", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_pkt_rate_high {
            ret.insert("coalesce-pkt-rate-high", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_pkt_rate_low {
            ret.insert("coalesce-pkt-rate-low", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_rx_frames {
            ret.insert("coalesce-rx-frames", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_rx_frames_low {
            ret.insert("coalesce-rx-frames-low", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_rx_frames_high {
            ret.insert("coalesce-rx-frames-high", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_rx_frames_irq {
            ret.insert("coalesce-rx-frames-irq", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_tx_frames {
            ret.insert("coalesce-tx-frames", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_tx_frames_low {
            ret.insert("coalesce-tx-frames-low", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_tx_frames_high {
            ret.insert("coalesce-tx-frames-high", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_tx_frames_irq {
            ret.insert("coalesce-tx-frames-irq", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_rx_usecs {
            ret.insert("coalesce-rx-usecs", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_rx_usecs_low {
            ret.insert("coalesce-rx-usecs-low", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_rx_usecs_high {
            ret.insert("coalesce-rx-usecs-high", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_rx_usecs_irq {
            ret.insert("coalesce-rx-usecs-irq", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_tx_usecs {
            ret.insert("coalesce-tx-usecs", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_tx_usecs_low {
            ret.insert("coalesce-tx-usecs-low", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_tx_usecs_high {
            ret.insert("coalesce-tx-usecs-high", zvariant::Value::new(v));
        }
        if let Some(v) = &self.coalesce_tx_usecs_irq {
            ret.insert("coalesce-tx-usecs-irq", zvariant::Value::new(v));
        }
        if let Some(features) = &self.features {
            for (k, v) in features {
                ret.insert(k, zvariant::Value::new(v));
            }
        }
        if let Some(v) = &self.ring_rx {
            ret.insert("ring-rx", zvariant::Value::new(v));
        }
        if let Some(v) = &self.ring_rx_jumbo {
            ret.insert("ring-rx-jumbo", zvariant::Value::new(v));
        }
        if let Some(v) = &self.ring_rx_mini {
            ret.insert("ring-rx-mini", zvariant::Value::new(v));
        }
        if let Some(v) = &self.ring_tx {
            ret.insert("ring-tx", zvariant::Value::new(v));
        }
        Ok(ret)
    }
}
