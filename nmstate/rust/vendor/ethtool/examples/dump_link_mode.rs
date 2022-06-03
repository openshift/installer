// SPDX-License-Identifier: MIT

use futures::stream::TryStreamExt;

// Once we find a way to load netsimdev kernel module in CI, we can convert this
// to a test
fn main() {
    let rt = tokio::runtime::Builder::new_current_thread()
        .enable_io()
        .build()
        .unwrap();
    rt.block_on(get_link_mode(None));
}

async fn get_link_mode(iface_name: Option<&str>) {
    let (connection, mut handle, _) = ethtool::new_connection().unwrap();
    tokio::spawn(connection);

    let mut link_mode_handle = handle.link_mode().get(iface_name).execute().await;

    let mut msgs = Vec::new();
    while let Some(msg) = link_mode_handle.try_next().await.unwrap() {
        msgs.push(msg);
    }
    assert!(!msgs.is_empty());
    for msg in msgs {
        println!("{:?}", msg);
    }
}
