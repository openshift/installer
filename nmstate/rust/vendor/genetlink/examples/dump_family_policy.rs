// SPDX-License-Identifier: MIT

use std::env::args;

use anyhow::{bail, Error};
use futures::StreamExt;
use genetlink::new_connection;
use netlink_packet_core::{
    NetlinkHeader,
    NetlinkMessage,
    NetlinkPayload,
    NLM_F_DUMP,
    NLM_F_REQUEST,
};
use netlink_packet_generic::{
    ctrl::{nlas::GenlCtrlAttrs, GenlCtrl, GenlCtrlCmd},
    GenlMessage,
};

#[tokio::main]
async fn main() -> Result<(), Error> {
    let argv: Vec<_> = args().collect();

    if argv.len() < 2 {
        eprintln!("Usage: dump_family_policy <family name>");
        bail!("Required arguments not given");
    }

    let nlmsg = NetlinkMessage {
        header: NetlinkHeader {
            flags: NLM_F_REQUEST | NLM_F_DUMP,
            ..Default::default()
        },
        payload: GenlMessage::from_payload(GenlCtrl {
            cmd: GenlCtrlCmd::GetPolicy,
            nlas: vec![GenlCtrlAttrs::FamilyName(argv[1].to_owned())],
        })
        .into(),
    };
    let (conn, mut handle, _) = new_connection()?;
    tokio::spawn(conn);

    let mut responses = handle.request(nlmsg).await?;

    while let Some(result) = responses.next().await {
        let resp = result?;
        match resp.payload {
            NetlinkPayload::InnerMessage(genlmsg) => {
                if genlmsg.payload.cmd == GenlCtrlCmd::GetPolicy {
                    println!("<<< {:?}", genlmsg);
                }
            }
            NetlinkPayload::Error(err) => {
                eprintln!("Received a netlink error message: {:?}", err);
                bail!(err);
            }
            _ => {}
        }
    }

    Ok(())
}
