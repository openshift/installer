use crate::{
    EthtoolCoalesceConfig, EthtoolConfig, EthtoolPauseConfig, EthtoolRingConfig,
};

pub(crate) fn np_ethtool_to_nmstate(
    np_iface: &nispor::Iface,
) -> Option<EthtoolConfig> {
    np_iface.ethtool.as_ref().map(gen_ethtool_config)
}

fn gen_ethtool_config(ethtool_info: &nispor::EthtoolInfo) -> EthtoolConfig {
    let mut ret = EthtoolConfig::new();
    if let Some(pause) = &ethtool_info.pause {
        let mut pause_config = EthtoolPauseConfig::new();
        pause_config.rx = Some(pause.rx);
        pause_config.tx = Some(pause.tx);
        pause_config.autoneg = Some(pause.auto_negotiate);
        ret.pause = Some(pause_config);
    }
    if let Some(feature) = &ethtool_info.features {
        ret.feature = Some(feature.changeable.clone());
    }
    if let Some(coalesce) = &ethtool_info.coalesce {
        let mut coalesce_config = EthtoolCoalesceConfig::new();
        coalesce_config.pkt_rate_high = coalesce.pkt_rate_high;
        coalesce_config.pkt_rate_low = coalesce.pkt_rate_low;
        coalesce_config.sample_interval = coalesce.rate_sample_interval;
        coalesce_config.rx_frames = coalesce.rx_max_frames;
        coalesce_config.rx_frames_high = coalesce.rx_max_frames_high;
        coalesce_config.rx_frames_low = coalesce.rx_max_frames_low;
        coalesce_config.rx_usecs = coalesce.rx_usecs;
        coalesce_config.rx_usecs_high = coalesce.rx_usecs_high;
        coalesce_config.rx_usecs_irq = coalesce.rx_usecs_irq;
        coalesce_config.rx_usecs_low = coalesce.rx_usecs_low;
        coalesce_config.stats_block_usecs = coalesce.stats_block_usecs;
        coalesce_config.tx_frames = coalesce.tx_max_frames;
        coalesce_config.tx_frames_high = coalesce.tx_max_frames_high;
        coalesce_config.tx_frames_low = coalesce.tx_max_frames_low;
        coalesce_config.tx_frames_irq = coalesce.tx_max_frames_irq;
        coalesce_config.tx_usecs = coalesce.tx_usecs;
        coalesce_config.tx_usecs_high = coalesce.tx_usecs_high;
        coalesce_config.tx_usecs_low = coalesce.tx_usecs_low;
        coalesce_config.tx_usecs_irq = coalesce.tx_usecs_irq;
        coalesce_config.adaptive_rx = coalesce.use_adaptive_rx;
        coalesce_config.adaptive_tx = coalesce.use_adaptive_tx;

        ret.coalesce = Some(coalesce_config);
    }
    if let Some(ring) = &ethtool_info.ring {
        let mut ring_config = EthtoolRingConfig::new();
        ring_config.rx = ring.rx;
        ring_config.rx_max = ring.rx_max;
        ring_config.rx_jumbo = ring.rx_jumbo;
        ring_config.rx_jumbo_max = ring.rx_jumbo_max;
        ring_config.rx_mini = ring.rx_mini;
        ring_config.rx_mini_max = ring.rx_mini_max;
        ring_config.tx = ring.tx;
        ring_config.tx_max = ring.tx_max;

        ret.ring = Some(ring_config);
    }
    ret
}
