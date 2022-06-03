// SPDX-License-Identifier: MIT

use crate::{
    error::GenetlinkError,
    message::{map_from_rawgenlmsg, map_to_rawgenlmsg, RawGenlMessage},
    resolver::Resolver,
};
use futures::{lock::Mutex, Stream, StreamExt};
use netlink_packet_core::{DecodeError, NetlinkMessage, NetlinkPayload};
use netlink_packet_generic::{GenlFamily, GenlHeader, GenlMessage};
use netlink_packet_utils::{Emitable, ParseableParametrized};
use netlink_proto::{sys::SocketAddr, ConnectionHandle};
use std::{fmt::Debug, sync::Arc};

/// The generic netlink connection handle
///
/// The handle is used to send messages to the connection. It also resolves
/// the family id automatically before sending messages.
///
/// # Family id resolving
/// There is a resolver with cache inside each connection. When you send generic
/// netlink message, the handle resolves and fills the family id into the message.
///
/// Since the resolver is created in [`new_connection()`](crate::new_connection),
/// the cache state wouldn't share between different connections.
///
/// P.s. The cloned handles use the same connection with the original handle. So,
/// they share the same cache state.
///
/// # Detailed process of sending generic messages
/// 1. Check if the message's family id is resolved. If yes, jump to step 6.
/// 2. Query the family id using the builtin resolver.
/// 3. If the id is in the cache, returning the id in the cache and skip step 4.
/// 4. The resolver sends `CTRL_CMD_GETFAMILY` request to get the id and records it in the cache.
/// 5. fill the family id using [`GenlMessage::set_resolved_family_id()`].
/// 6. Serialize the payload to [`RawGenlMessage`].
/// 7. Send it through the connection.
///     - The family id filled into `message_type` field in [`NetlinkMessage::finalize()`].
/// 8. In the response stream, deserialize the payload back to [`GenlMessage<F>`].
#[derive(Clone, Debug)]
pub struct GenetlinkHandle {
    handle: ConnectionHandle<RawGenlMessage>,
    resolver: Arc<Mutex<Resolver>>,
}

impl GenetlinkHandle {
    pub(crate) fn new(handle: ConnectionHandle<RawGenlMessage>) -> Self {
        Self {
            handle,
            resolver: Arc::new(Mutex::new(Resolver::new())),
        }
    }

    /// Resolve the family id of the given [`GenlFamily`].
    pub async fn resolve_family_id<F>(&self) -> Result<u16, GenetlinkError>
    where
        F: GenlFamily,
    {
        self.resolver
            .lock()
            .await
            .query_family_id(self, F::family_name())
            .await
    }

    /// Clear the resolver's fanily id cache
    pub async fn clear_family_id_cache(&self) {
        self.resolver.lock().await.clear_cache();
    }

    /// Send the generic netlink message and get the response stream
    ///
    /// The function resolves the family id before sending the request. If the
    /// resolving process is failed, the function would return an error.
    pub async fn request<F>(
        &mut self,
        mut message: NetlinkMessage<GenlMessage<F>>,
    ) -> Result<
        impl Stream<Item = Result<NetlinkMessage<GenlMessage<F>>, DecodeError>>,
        GenetlinkError,
    >
    where
        F: GenlFamily + Emitable + ParseableParametrized<[u8], GenlHeader> + Debug,
    {
        self.resolve_message_family_id(&mut message).await?;
        self.send_request(message)
    }

    /// Send the request without resolving family id
    ///
    /// This function is identical to [`request()`](Self::request) but it doesn't
    /// resolve the family id for you.
    pub fn send_request<F>(
        &mut self,
        message: NetlinkMessage<GenlMessage<F>>,
    ) -> Result<
        impl Stream<Item = Result<NetlinkMessage<GenlMessage<F>>, DecodeError>>,
        GenetlinkError,
    >
    where
        F: GenlFamily + Emitable + ParseableParametrized<[u8], GenlHeader> + Debug,
    {
        let raw_msg = map_to_rawgenlmsg(message);

        let stream = self.handle.request(raw_msg, SocketAddr::new(0, 0))?;
        Ok(stream.map(map_from_rawgenlmsg))
    }

    /// Send the generic netlink message without returning the response stream
    pub async fn notify<F>(
        &mut self,
        mut message: NetlinkMessage<GenlMessage<F>>,
    ) -> Result<(), GenetlinkError>
    where
        F: GenlFamily + Emitable + ParseableParametrized<[u8], GenlHeader> + Debug,
    {
        self.resolve_message_family_id(&mut message).await?;
        self.send_notify(message)
    }

    /// Send the notify without resolving family id
    pub fn send_notify<F>(
        &mut self,
        message: NetlinkMessage<GenlMessage<F>>,
    ) -> Result<(), GenetlinkError>
    where
        F: GenlFamily + Emitable + ParseableParametrized<[u8], GenlHeader> + Debug,
    {
        let raw_msg = map_to_rawgenlmsg(message);

        self.handle.notify(raw_msg, SocketAddr::new(0, 0))?;
        Ok(())
    }

    async fn resolve_message_family_id<F>(
        &mut self,
        message: &mut NetlinkMessage<GenlMessage<F>>,
    ) -> Result<(), GenetlinkError>
    where
        F: GenlFamily + Debug,
    {
        if let NetlinkPayload::InnerMessage(genlmsg) = &mut message.payload {
            if genlmsg.family_id() == 0 {
                // The family id is not resolved
                // Resolve it before send it
                let id = self.resolve_family_id::<F>().await?;
                genlmsg.set_resolved_family_id(id);
            }
        }

        Ok(())
    }
}
