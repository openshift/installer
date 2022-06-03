// SPDX-License-Identifier: MIT

use crate::{error::GenetlinkError, GenetlinkHandle};
use futures::{future::Either, StreamExt};
use netlink_packet_core::{NetlinkMessage, NetlinkPayload, NLM_F_REQUEST};
use netlink_packet_generic::{
    ctrl::{nlas::GenlCtrlAttrs, GenlCtrl, GenlCtrlCmd},
    GenlMessage,
};
use std::{collections::HashMap, future::Future};

#[derive(Clone, Debug, Default)]
pub struct Resolver {
    cache: HashMap<&'static str, u16>,
}

impl Resolver {
    pub fn new() -> Self {
        Self {
            cache: HashMap::new(),
        }
    }

    pub fn get_cache_by_name(&self, family_name: &str) -> Option<u16> {
        self.cache.get(family_name).copied()
    }

    pub fn query_family_id(
        &mut self,
        handle: &GenetlinkHandle,
        family_name: &'static str,
    ) -> impl Future<Output = Result<u16, GenetlinkError>> + '_ {
        if let Some(id) = self.get_cache_by_name(family_name) {
            Either::Left(futures::future::ready(Ok(id)))
        } else {
            let mut handle = handle.clone();
            Either::Right(async move {
                let mut genlmsg: GenlMessage<GenlCtrl> = GenlMessage::from_payload(GenlCtrl {
                    cmd: GenlCtrlCmd::GetFamily,
                    nlas: vec![GenlCtrlAttrs::FamilyName(family_name.to_owned())],
                });
                genlmsg.finalize();
                // We don't have to set family id here, since nlctrl has static family id (0x10)
                let mut nlmsg = NetlinkMessage::from(genlmsg);
                nlmsg.header.flags = NLM_F_REQUEST;
                nlmsg.finalize();

                let mut res = handle.send_request(nlmsg)?;

                while let Some(result) = res.next().await {
                    let rx_packet = result?;
                    match rx_packet.payload {
                        NetlinkPayload::InnerMessage(genlmsg) => {
                            let family_id = genlmsg
                                .payload
                                .nlas
                                .iter()
                                .find_map(|nla| {
                                    if let GenlCtrlAttrs::FamilyId(id) = nla {
                                        Some(*id)
                                    } else {
                                        None
                                    }
                                })
                                .ok_or_else(|| {
                                    GenetlinkError::AttributeNotFound(
                                        "CTRL_ATTR_FAMILY_ID".to_owned(),
                                    )
                                })?;

                            self.cache.insert(family_name, family_id);
                            return Ok(family_id);
                        }
                        NetlinkPayload::Error(e) => return Err(e.into()),
                        _ => (),
                    }
                }

                Err(GenetlinkError::NoMessageReceived)
            })
        }
    }

    pub fn clear_cache(&mut self) {
        self.cache.clear();
    }
}

#[cfg(test)]
mod test {
    use super::*;
    use crate::new_connection;
    use std::io::ErrorKind;

    #[tokio::test]
    async fn test_resolver_nlctrl() {
        let (conn, handle, _) = new_connection().unwrap();
        tokio::spawn(conn);

        let mut resolver = Resolver::new();
        // nlctrl should always be 0x10
        let nlctrl_fid = resolver.query_family_id(&handle, "nlctrl").await.unwrap();
        assert_eq!(nlctrl_fid, 0x10);
    }

    const TEST_FAMILIES: &[&str] = &[
        "devlink",
        "ethtool",
        "acpi_event",
        "tcp_metrics",
        "TASKSTATS",
        "nl80211",
    ];

    #[tokio::test]
    async fn test_resolver_cache() {
        let (conn, handle, _) = new_connection().unwrap();
        tokio::spawn(conn);

        let mut resolver = Resolver::new();

        // Test if family id cached
        for name in TEST_FAMILIES.iter().copied() {
            let id = resolver
                .query_family_id(&handle, name)
                .await
                .or_else(|e| {
                    if let GenetlinkError::NetlinkError(io_err) = &e {
                        if io_err.kind() == ErrorKind::NotFound {
                            // Ignore non exist entries
                            Ok(0)
                        } else {
                            Err(e)
                        }
                    } else {
                        Err(e)
                    }
                })
                .unwrap();
            if id == 0 {
                eprintln!(
                    "Generic family \"{}\" not exist or not loaded in this environment. Ignored.",
                    name
                );
                continue;
            }

            let cache = resolver.get_cache_by_name(name).unwrap();
            assert_eq!(id, cache);
            eprintln!("{:?}", (name, cache));
        }
    }
}
