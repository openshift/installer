use crate::EthernetInterface;

#[test]
fn test_ethtool_stringlized_attributes() {
    let iface: EthernetInterface = serde_yaml::from_str(
        r#"---
name: eth1
type: ethernet
state: up
ethtool:
  pause:
    rx: "off"
    tx: "off"
    autoneg: "off"
  feature:
    rx-checksum: "true"
    rx-gro: "true"
    rx-lro: "true"
    rx-vlan-hw-parse: "true"
    tx-vlan-hw-insert: "true"
    rx-ntuple-filter: "true"
    rx-hashing: "true"
    tx-scatter-gather: "true"
    tx-tcp-segmentation: "true"
    tx-generic-segmentation: "true"
  coalesce:
    adaptive-rx: "0"
    adaptive-tx: "0"
    pkt-rate-high: "101"
    pkt-rate-low: "102"
    rx-frames: "103"
    rx-frames-high: "104"
    rx-frames-irq: "105"
    rx-frames-low: "106"
    rx-usecs: "107"
    rx-usecs-high: "108"
    rx-usecs-irq: "109"
    rx-usecs-low: "110"
    sample-interval: "111"
    stats-block-usecs: "112"
    tx-frames: "113"
    tx-frames-high: "114"
    tx-frames-irq: "115"
    tx-frames-low: "116"
    tx-usecs: "117"
    tx-usecs-high: "118"
    tx-usecs-irq: "119"
    tx-usecs-low: "120"
  ring:
    rx: "200"
    rx-max: "201"
    rx-jumbo: "202"
    rx-jumbo-max: "203"
    rx-mini: "204"
    rx-mini-max: "205"
    tx: "206"
    tx-max: "207"

"#,
    )
    .unwrap();

    let ethtool_conf = iface.base.ethtool.unwrap();
    let features = ethtool_conf.feature.as_ref().unwrap();
    let pause = ethtool_conf.pause.as_ref().unwrap();
    let coalesce = ethtool_conf.coalesce.as_ref().unwrap();
    let ring = ethtool_conf.ring.as_ref().unwrap();

    assert_eq!(features.get("rx-checksum"), Some(&true));
    assert_eq!(features.get("rx-gro"), Some(&true));
    assert_eq!(features.get("rx-lro"), Some(&true));
    assert_eq!(features.get("rx-vlan-hw-parse"), Some(&true));
    assert_eq!(features.get("tx-vlan-hw-insert"), Some(&true));
    assert_eq!(features.get("rx-ntuple-filter"), Some(&true));
    assert_eq!(features.get("rx-hashing"), Some(&true));
    assert_eq!(features.get("tx-scatter-gather"), Some(&true));
    assert_eq!(features.get("tx-tcp-segmentation"), Some(&true));
    assert_eq!(features.get("tx-generic-segmentation"), Some(&true));
    assert_eq!(pause.tx, Some(false));
    assert_eq!(pause.rx, Some(false));
    assert_eq!(pause.autoneg, Some(false));

    assert_eq!(coalesce.adaptive_rx, Some(false));
    assert_eq!(coalesce.adaptive_tx, Some(false));
    assert_eq!(coalesce.pkt_rate_high, Some(101));
    assert_eq!(coalesce.pkt_rate_low, Some(102));
    assert_eq!(coalesce.rx_frames, Some(103));
    assert_eq!(coalesce.rx_frames_high, Some(104));
    assert_eq!(coalesce.rx_frames_irq, Some(105));
    assert_eq!(coalesce.rx_frames_low, Some(106));
    assert_eq!(coalesce.rx_usecs, Some(107));
    assert_eq!(coalesce.rx_usecs_high, Some(108));
    assert_eq!(coalesce.rx_usecs_irq, Some(109));
    assert_eq!(coalesce.rx_usecs_low, Some(110));
    assert_eq!(coalesce.sample_interval, Some(111));
    assert_eq!(coalesce.stats_block_usecs, Some(112));
    assert_eq!(coalesce.tx_frames, Some(113));
    assert_eq!(coalesce.tx_frames_high, Some(114));
    assert_eq!(coalesce.tx_frames_irq, Some(115));
    assert_eq!(coalesce.tx_frames_low, Some(116));
    assert_eq!(coalesce.tx_usecs, Some(117));
    assert_eq!(coalesce.tx_usecs_high, Some(118));
    assert_eq!(coalesce.tx_usecs_irq, Some(119));
    assert_eq!(coalesce.tx_usecs_low, Some(120));
    assert_eq!(ring.rx, Some(200));
    assert_eq!(ring.rx_max, Some(201));
    assert_eq!(ring.rx_jumbo, Some(202));
    assert_eq!(ring.rx_jumbo_max, Some(203));
    assert_eq!(ring.rx_mini, Some(204));
    assert_eq!(ring.rx_mini_max, Some(205));
    assert_eq!(ring.tx, Some(206));
    assert_eq!(ring.tx_max, Some(207));
}
