// SPDX-License-Identifier: MIT

use futures::stream::TryStreamExt;

#[test]
fn test_get_features_of_loopback() {
    let rt = tokio::runtime::Builder::new_current_thread()
        .enable_io()
        .build()
        .unwrap();
    rt.block_on(get_feature(Some("lo")));
}

async fn get_feature(iface_name: Option<&str>) {
    let (connection, mut handle, _) = ethtool::new_connection().unwrap();
    tokio::spawn(connection);

    let mut feature_handle = handle.feature().get(iface_name).execute().await;

    let mut msgs = Vec::new();
    while let Some(msg) = feature_handle.try_next().await.unwrap() {
        msgs.push(msg);
    }
    assert!(msgs.len() == 1);
    let ethtool_msg = &msgs[0].payload;

    assert!(ethtool_msg.cmd == ethtool::EthtoolCmd::FeatureGetReply);
    assert!(ethtool_msg.nlas.len() > 1);
    assert!(
        ethtool_msg.nlas[0]
            == ethtool::EthtoolAttr::Feature(ethtool::EthtoolFeatureAttr::Header(vec![
                ethtool::EthtoolHeader::DevIndex(1),
                ethtool::EthtoolHeader::DevName("lo".into())
            ]))
    );
}
