// SPDX-License-Identifier: MIT

use futures::stream::TryStreamExt;

#[test]
// CI container normally have a veth for external communication which support link modes of ethtool.
fn test_dump_link_modes() {
    let rt = tokio::runtime::Builder::new_current_thread()
        .enable_io()
        .build()
        .unwrap();
    rt.block_on(dump_link_modes());
}

async fn dump_link_modes() {
    let (connection, mut handle, _) = ethtool::new_connection().unwrap();
    tokio::spawn(connection);

    let mut link_modes_handle = handle.link_mode().get(None).execute().await;

    let mut msgs = Vec::new();
    while let Some(msg) = link_modes_handle.try_next().await.unwrap() {
        msgs.push(msg);
    }
    assert!(!msgs.is_empty());
    let ethtool_msg = &msgs[0].payload;
    println!("ethtool_msg {:?}", &ethtool_msg);

    assert!(ethtool_msg.cmd == ethtool::EthtoolCmd::LinkModeGetReply);
    assert!(ethtool_msg.nlas.len() > 1);
}
