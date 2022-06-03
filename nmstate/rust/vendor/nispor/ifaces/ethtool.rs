// Copyright 2021 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

use std::collections::{hash_map::Entry, BTreeMap, HashMap};

use ethtool::{
    EthtoolAttr, EthtoolCoalesceAttr, EthtoolFeatureAttr, EthtoolFeatureBit,
    EthtoolHandle, EthtoolHeader, EthtoolLinkModeAttr, EthtoolPauseAttr,
    EthtoolRingAttr,
};
use futures::stream::TryStreamExt;
use serde::{Deserialize, Serialize, Serializer};

use crate::NisporError;

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct EthtoolPauseInfo {
    pub rx: bool,
    pub tx: bool,
    pub auto_negotiate: bool,
}

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Clone)]
#[serde(rename_all = "snake_case")]
#[non_exhaustive]
pub struct EthtoolFeatureInfo {
    #[serde(serialize_with = "ordered_map")]
    pub fixed: HashMap<String, bool>,
    #[serde(serialize_with = "ordered_map")]
    pub changeable: HashMap<String, bool>,
}

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct EthtoolCoalesceInfo {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub pkt_rate_high: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub pkt_rate_low: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rate_sample_interval: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_max_frames: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_max_frames_high: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_max_frames_irq: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_max_frames_low: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_usecs: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_usecs_high: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_usecs_irq: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_usecs_low: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub stats_block_usecs: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx_max_frames: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx_max_frames_high: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx_max_frames_irq: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx_max_frames_low: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx_usecs: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx_usecs_high: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx_usecs_irq: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx_usecs_low: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub use_adaptive_rx: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub use_adaptive_tx: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct EthtoolRingInfo {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_max: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_jumbo: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_jumbo_max: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_mini: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rx_mini_max: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tx_max: Option<u32>,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Clone)]
#[serde(rename_all = "snake_case")]
#[non_exhaustive]
pub enum EthtoolLinkModeDuplex {
    Half,
    Full,
    Unknown,
    Other(u8),
}

impl From<&ethtool::EthtoolLinkModeDuplex> for EthtoolLinkModeDuplex {
    fn from(v: &ethtool::EthtoolLinkModeDuplex) -> Self {
        match v {
            ethtool::EthtoolLinkModeDuplex::Half => Self::Half,
            ethtool::EthtoolLinkModeDuplex::Full => Self::Full,
            ethtool::EthtoolLinkModeDuplex::Unknown => Self::Unknown,
            ethtool::EthtoolLinkModeDuplex::Other(d) => Self::Other(*d),
        }
    }
}

impl Default for EthtoolLinkModeDuplex {
    fn default() -> Self {
        EthtoolLinkModeDuplex::Unknown
    }
}

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct EthtoolLinkModeInfo {
    pub auto_negotiate: bool,
    pub ours: Vec<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub peer: Option<Vec<String>>,
    pub speed: u32,
    pub duplex: EthtoolLinkModeDuplex,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub controller_subordinate_cfg: Option<u8>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub controller_subordinate_state: Option<u8>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub lanes: Option<u32>,
}

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct EthtoolInfo {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub pause: Option<EthtoolPauseInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub features: Option<EthtoolFeatureInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub coalesce: Option<EthtoolCoalesceInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ring: Option<EthtoolRingInfo>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub link_mode: Option<EthtoolLinkModeInfo>,
}

fn ordered_map<S>(
    value: &HashMap<String, bool>,
    serializer: S,
) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    let ordered: BTreeMap<_, _> = value.iter().collect();
    ordered.serialize(serializer)
}

pub(crate) async fn get_ethtool_infos(
) -> Result<HashMap<String, EthtoolInfo>, NisporError> {
    let mut infos: HashMap<String, EthtoolInfo> = HashMap::new();

    let (connection, mut handle, _) = ethtool::new_connection()?;

    tokio::spawn(connection);

    let mut pause_infos = dump_pause_infos(&mut handle).await?;
    let mut feature_infos = dump_feature_infos(&mut handle).await?;
    let mut coalesce_infos = dump_coalesce_infos(&mut handle).await?;
    let mut ring_infos = dump_ring_infos(&mut handle).await?;
    let mut link_mode_infos = dump_link_mode_infos(&mut handle).await?;

    for (iface_name, pause_info) in pause_infos.drain() {
        infos.insert(
            iface_name,
            EthtoolInfo {
                pause: Some(pause_info),
                ..Default::default()
            },
        );
    }

    for (iface_name, feature_info) in feature_infos.drain() {
        match infos.get_mut(&iface_name) {
            Some(ref mut info) => {
                info.features = Some(feature_info);
            }
            None => {
                infos.insert(
                    iface_name,
                    EthtoolInfo {
                        features: Some(feature_info),
                        ..Default::default()
                    },
                );
            }
        };
    }

    for (iface_name, coalesce_info) in coalesce_infos.drain() {
        match infos.get_mut(&iface_name) {
            Some(ref mut info) => {
                info.coalesce = Some(coalesce_info);
            }
            None => {
                infos.insert(
                    iface_name,
                    EthtoolInfo {
                        coalesce: Some(coalesce_info),
                        ..Default::default()
                    },
                );
            }
        };
    }

    for (iface_name, ring_info) in ring_infos.drain() {
        match infos.get_mut(&iface_name) {
            Some(ref mut info) => {
                info.ring = Some(ring_info);
            }
            None => {
                infos.insert(
                    iface_name,
                    EthtoolInfo {
                        ring: Some(ring_info),
                        ..Default::default()
                    },
                );
            }
        };
    }

    for (iface_name, link_mode_info) in link_mode_infos.drain() {
        match infos.get_mut(&iface_name) {
            Some(ref mut info) => {
                info.link_mode = Some(link_mode_info);
            }
            None => {
                infos.insert(
                    iface_name,
                    EthtoolInfo {
                        link_mode: Some(link_mode_info),
                        ..Default::default()
                    },
                );
            }
        };
    }

    Ok(infos)
}

async fn dump_pause_infos(
    handle: &mut EthtoolHandle,
) -> Result<HashMap<String, EthtoolPauseInfo>, NisporError> {
    let mut infos = HashMap::new();
    let mut pause_handle = handle.pause().get(None).execute().await;
    while let Some(genl_msg) = pause_handle.try_next().await? {
        let ethtool_msg = genl_msg.payload;
        let mut iface_name = None;
        let mut pause_info = EthtoolPauseInfo::default();

        for nla in &ethtool_msg.nlas {
            if let EthtoolAttr::Pause(nla) = nla {
                if let EthtoolPauseAttr::Header(hdrs) = nla {
                    iface_name = get_iface_name_from_header(hdrs);
                } else if let EthtoolPauseAttr::AutoNeg(v) = nla {
                    pause_info.auto_negotiate = *v
                } else if let EthtoolPauseAttr::Rx(v) = nla {
                    pause_info.rx = *v
                } else if let EthtoolPauseAttr::Tx(v) = nla {
                    pause_info.tx = *v
                }
            }
        }
        if let Some(i) = iface_name {
            infos.insert(i, pause_info);
        }
    }
    Ok(infos)
}

async fn dump_feature_infos(
    handle: &mut EthtoolHandle,
) -> Result<HashMap<String, EthtoolFeatureInfo>, NisporError> {
    let mut infos = HashMap::new();
    let mut feature_handle = handle.feature().get(None).execute().await;
    while let Some(genl_msg) = feature_handle.try_next().await? {
        let ethtool_msg = genl_msg.payload;
        let mut iface_name = None;
        let mut fixed_features: HashMap<String, bool> = HashMap::new();
        let mut changeable_features: HashMap<String, bool> = HashMap::new();

        for nla in &ethtool_msg.nlas {
            if let EthtoolAttr::Feature(EthtoolFeatureAttr::NoChange(
                feature_bits,
            )) = nla
            {
                for EthtoolFeatureBit { name, .. } in feature_bits {
                    fixed_features.insert(name.to_string(), false);
                }
            }
        }

        for nla in &ethtool_msg.nlas {
            if let EthtoolAttr::Feature(EthtoolFeatureAttr::Header(hdrs)) = nla
            {
                iface_name = get_iface_name_from_header(hdrs);
                break;
            }
        }

        for nla in &ethtool_msg.nlas {
            if let EthtoolAttr::Feature(EthtoolFeatureAttr::Hw(feature_bits)) =
                nla
            {
                for feature_bit in feature_bits {
                    match feature_bit {
                        EthtoolFeatureBit {
                            index: _,
                            name,
                            value: true,
                        } => {
                            // Dummy interface show `tx-lockless` is
                            // changeable, but EthtoolFeatureAttr::NoChange() says
                            // otherwise. The kernel code
                            // `NETIF_F_NEVER_CHANGE` shows `tx-lockless`
                            // should never been changeable.
                            if let Entry::Occupied(mut e) =
                                fixed_features.entry(name.clone())
                            {
                                e.insert(false);
                            } else {
                                changeable_features
                                    .insert(name.to_string(), false);
                            }
                        }
                        EthtoolFeatureBit {
                            index: _,
                            name,
                            value: false,
                        } => {
                            fixed_features.insert(name.to_string(), false);
                        }
                    }
                }
            }
        }

        for nla in &ethtool_msg.nlas {
            if let EthtoolAttr::Feature(EthtoolFeatureAttr::Active(
                feature_bits,
            )) = nla
            {
                for feature_bit in feature_bits {
                    if let Entry::Occupied(mut e) =
                        fixed_features.entry(feature_bit.name.clone())
                    {
                        e.insert(true);
                    } else if changeable_features
                        .contains_key(&feature_bit.name)
                    {
                        changeable_features
                            .insert(feature_bit.name.clone(), true);
                    }
                }
            }
        }

        if let Some(i) = iface_name {
            infos.insert(
                i.to_string(),
                EthtoolFeatureInfo {
                    fixed: fixed_features,
                    changeable: changeable_features,
                },
            );
        }
    }

    Ok(infos)
}

async fn dump_coalesce_infos(
    handle: &mut EthtoolHandle,
) -> Result<HashMap<String, EthtoolCoalesceInfo>, NisporError> {
    let mut infos = HashMap::new();
    let mut coalesce_handle = handle.coalesce().get(None).execute().await;
    while let Some(genl_msg) = coalesce_handle.try_next().await? {
        let ethtool_msg = genl_msg.payload;
        let mut iface_name = None;
        let mut coalesce_info = EthtoolCoalesceInfo::default();
        for nla in &ethtool_msg.nlas {
            if let EthtoolAttr::Coalesce(nla) = nla {
                match nla {
                    EthtoolCoalesceAttr::Header(hdrs) => {
                        iface_name = get_iface_name_from_header(hdrs)
                    }
                    EthtoolCoalesceAttr::RxUsecs(d) => {
                        coalesce_info.rx_usecs = Some(*d)
                    }
                    EthtoolCoalesceAttr::RxMaxFrames(d) => {
                        coalesce_info.rx_max_frames = Some(*d)
                    }
                    EthtoolCoalesceAttr::RxUsecsIrq(d) => {
                        coalesce_info.rx_usecs_irq = Some(*d)
                    }
                    EthtoolCoalesceAttr::RxMaxFramesIrq(d) => {
                        coalesce_info.rx_max_frames_irq = Some(*d)
                    }
                    EthtoolCoalesceAttr::TxUsecs(d) => {
                        coalesce_info.tx_usecs = Some(*d)
                    }
                    EthtoolCoalesceAttr::TxMaxFrames(d) => {
                        coalesce_info.tx_max_frames = Some(*d)
                    }
                    EthtoolCoalesceAttr::TxUsecsIrq(d) => {
                        coalesce_info.tx_usecs_irq = Some(*d)
                    }
                    EthtoolCoalesceAttr::TxMaxFramesIrq(d) => {
                        coalesce_info.tx_max_frames_irq = Some(*d)
                    }
                    EthtoolCoalesceAttr::StatsBlockUsecs(d) => {
                        coalesce_info.stats_block_usecs = Some(*d)
                    }
                    EthtoolCoalesceAttr::UseAdaptiveRx(d) => {
                        coalesce_info.use_adaptive_rx = Some(*d)
                    }
                    EthtoolCoalesceAttr::UseAdaptiveTx(d) => {
                        coalesce_info.use_adaptive_tx = Some(*d)
                    }
                    EthtoolCoalesceAttr::PktRateLow(d) => {
                        coalesce_info.pkt_rate_low = Some(*d)
                    }
                    EthtoolCoalesceAttr::RxUsecsLow(d) => {
                        coalesce_info.rx_usecs_low = Some(*d)
                    }
                    EthtoolCoalesceAttr::RxMaxFramesLow(d) => {
                        coalesce_info.rx_max_frames_low = Some(*d)
                    }
                    EthtoolCoalesceAttr::TxUsecsLow(d) => {
                        coalesce_info.tx_usecs_low = Some(*d)
                    }
                    EthtoolCoalesceAttr::TxMaxFramesLow(d) => {
                        coalesce_info.tx_max_frames_low = Some(*d)
                    }
                    EthtoolCoalesceAttr::PktRateHigh(d) => {
                        coalesce_info.pkt_rate_high = Some(*d)
                    }
                    EthtoolCoalesceAttr::RxUsecsHigh(d) => {
                        coalesce_info.rx_usecs_high = Some(*d)
                    }
                    EthtoolCoalesceAttr::RxMaxFramesHigh(d) => {
                        coalesce_info.rx_max_frames_high = Some(*d)
                    }
                    EthtoolCoalesceAttr::TxUsecsHigh(d) => {
                        coalesce_info.tx_usecs_high = Some(*d)
                    }
                    EthtoolCoalesceAttr::TxMaxFramesHigh(d) => {
                        coalesce_info.tx_max_frames_high = Some(*d)
                    }
                    EthtoolCoalesceAttr::RateSampleInterval(d) => {
                        coalesce_info.rate_sample_interval = Some(*d)
                    }
                    _ => log::warn!(
                        "WARN: Unsupported EthtoolCoalesceAttr {:?}",
                        nla
                    ),
                }
            }
        }
        if let Some(i) = iface_name {
            infos.insert(i, coalesce_info);
        }
    }
    Ok(infos)
}

async fn dump_ring_infos(
    handle: &mut EthtoolHandle,
) -> Result<HashMap<String, EthtoolRingInfo>, NisporError> {
    let mut infos = HashMap::new();
    let mut ring_handle = handle.ring().get(None).execute().await;
    while let Some(genl_msg) = ring_handle.try_next().await? {
        let ethtool_msg = genl_msg.payload;
        let mut iface_name = None;
        let mut ring_info = EthtoolRingInfo::default();
        for nla in &ethtool_msg.nlas {
            if let EthtoolAttr::Ring(nla) = nla {
                match nla {
                    EthtoolRingAttr::Header(hdrs) => {
                        iface_name = get_iface_name_from_header(hdrs)
                    }
                    EthtoolRingAttr::RxMax(d) => ring_info.rx_max = Some(*d),
                    EthtoolRingAttr::RxMiniMax(d) => {
                        ring_info.rx_mini_max = Some(*d)
                    }
                    EthtoolRingAttr::RxJumboMax(d) => {
                        ring_info.rx_jumbo_max = Some(*d)
                    }
                    EthtoolRingAttr::TxMax(d) => ring_info.tx_max = Some(*d),
                    EthtoolRingAttr::Rx(d) => ring_info.rx = Some(*d),
                    EthtoolRingAttr::RxMini(d) => ring_info.rx_mini = Some(*d),
                    EthtoolRingAttr::RxJumbo(d) => {
                        ring_info.rx_jumbo = Some(*d)
                    }
                    EthtoolRingAttr::Tx(d) => ring_info.tx = Some(*d),
                    _ => log::warn!(
                        "WARN: Unsupported EthtoolRingAttr {:?}",
                        nla
                    ),
                }
            }
        }
        if let Some(i) = iface_name {
            infos.insert(i, ring_info);
        }
    }
    Ok(infos)
}

async fn dump_link_mode_infos(
    handle: &mut EthtoolHandle,
) -> Result<HashMap<String, EthtoolLinkModeInfo>, NisporError> {
    let mut infos = HashMap::new();
    let mut link_mode_handle = handle.link_mode().get(None).execute().await;
    while let Some(genl_msg) = link_mode_handle.try_next().await? {
        let ethtool_msg = genl_msg.payload;
        let mut iface_name = None;
        let mut link_mode_info = EthtoolLinkModeInfo::default();
        for nla in &ethtool_msg.nlas {
            if let EthtoolAttr::LinkMode(nla) = nla {
                match nla {
                    EthtoolLinkModeAttr::Header(hdrs) => {
                        iface_name = get_iface_name_from_header(hdrs)
                    }
                    EthtoolLinkModeAttr::Autoneg(d) => {
                        link_mode_info.auto_negotiate = *d
                    }
                    EthtoolLinkModeAttr::Ours(d) => {
                        link_mode_info.ours = d.clone()
                    }
                    EthtoolLinkModeAttr::Peer(d) => {
                        link_mode_info.peer = Some(d.clone())
                    }
                    EthtoolLinkModeAttr::Speed(d) => {
                        link_mode_info.speed =
                            if *d == std::u32::MAX { 0 } else { *d }
                    }
                    EthtoolLinkModeAttr::Duplex(d) => {
                        link_mode_info.duplex = d.into()
                    }
                    EthtoolLinkModeAttr::ControllerSubordinateCfg(d) => {
                        link_mode_info.controller_subordinate_cfg = Some(*d)
                    }
                    EthtoolLinkModeAttr::ControllerSubordinateState(d) => {
                        link_mode_info.controller_subordinate_state = Some(*d)
                    }
                    EthtoolLinkModeAttr::Lanes(d) => {
                        link_mode_info.lanes = Some(*d)
                    }

                    _ => log::warn!(
                        "WARN: Unsupported EthtoolLinkModeAttr {:?}",
                        nla
                    ),
                }
            }
        }
        if let Some(i) = iface_name {
            infos.insert(i, link_mode_info);
        }
    }
    Ok(infos)
}

fn get_iface_name_from_header(hdrs: &[EthtoolHeader]) -> Option<String> {
    for hdr in hdrs {
        if let EthtoolHeader::DevName(iface_name) = hdr {
            return Some(iface_name.to_string());
        }
    }
    None
}
