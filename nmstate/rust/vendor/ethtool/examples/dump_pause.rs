// SPDX-License-Identifier: MIT

use futures::stream::TryStreamExt;

// Once we find a way to load netsimdev kernel module in CI, we can convert this
// to a test
fn main() {
    env_logger::init();
    let rt = tokio::runtime::Builder::new_current_thread()
        .enable_io()
        .build()
        .unwrap();
    rt.block_on(get_pause(None));
}

async fn get_pause(iface_name: Option<&str>) {
    let (connection, mut handle, _) = ethtool::new_connection().unwrap();
    tokio::spawn(connection);

    let mut pause_handle = handle.pause().get(iface_name).execute().await;

    let mut msgs = Vec::new();
    while let Some(msg) = pause_handle.try_next().await.unwrap() {
        msgs.push(msg);
    }
    assert!(!msgs.is_empty());
    println!("{:?}", msgs);
}
