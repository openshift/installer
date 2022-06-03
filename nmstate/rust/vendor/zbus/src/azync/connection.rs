use async_io::Async;
use once_cell::sync::OnceCell;
use std::{
    io::{self, ErrorKind},
    os::unix::{io::AsRawFd, net::UnixStream},
    pin::Pin,
    task::{Context, Poll},
};

use futures::{
    sink::{Sink, SinkExt},
    stream::{Stream, TryStreamExt},
};

use crate::{
    azync::{Authenticated, AuthenticatedType},
    raw::{Connection as RawConnection, Socket},
    Error, Guid, Message, MessageType, Result, DEFAULT_MAX_QUEUED,
};

/// The asynchronous sibling of [`zbus::Connection`].
///
/// Most of the API is very similar to [`zbus::Connection`], except it's asynchronous. However,
/// there are a few differences:
///
/// ### Generic over Socket
///
/// This type is generic over [`zbus::raw::Socket`] so that support for new socket types can be
/// added with the same type easily later on.
///
/// ### Cloning and Mutability
///
/// Unlike [`zbus::Connection`], this type does not implement [`std::clone::Clone`]. The reason is
/// that implementation will be very difficult (and still prone to deadlocks) if connection is
/// owned by multiple tasks/threads. Create separate connection instances or use
/// [`futures::stream::StreamExt::split`] to split reading and writing between two separate async
/// tasks.
///
/// Also notice that unlike [`zbus::Connection`], most methods take a `&mut self`, rather than a
/// `&self`. If they'd take `&self`, `Connection` will need to manage mutability internally, which
/// is not a very good match with the general async/await machinery and runtimes in Rust and could
/// easily lead into some hard-to-debug deadlocks. You can use [`std::cell::Cell`],
/// [`std::sync::Mutex`] or other related API combined with [`std::rc::Rc`] or [`std::sync::Arc`]
/// for sharing a mutable `Connection` instance between different parts of your code (or threads).
///
/// ### Sending Messages
///
/// For sending messages you can either use [`Connection::send_message`] method or make use of the
/// [`Sink`] implementation. For latter, you might find [`SinkExt`] API very useful. Keep in mind
/// that [`Connection`] will not manage the serial numbers (cookies) on the messages for you when
/// they are sent through the [`Sink`] implementation. You can manually assign unique serial numbers
/// to them using the [`Connection::assign_serial_num`] method before sending them off, if needed.
/// Having said that, [`Sink`] is mainly useful for sending out signals, as they do not expect a
/// reply, and serial numbers are not very useful for signals either for the same reason.
///
/// ### Receiving Messages
///
/// Unlike [`zbus::Connection`], there is no direct async equivalent of
/// [`zbus::Connection::receive_message`] method provided. This is because the `futures` crate
/// already provides a nice rich API that makes use of the  [`Stream`] implementation.
///
/// ### Examples
///
/// #### Get the session bus ID
///
/// ```
///# use zvariant::Type;
///#
///# futures::executor::block_on(async {
/// use zbus::azync::Connection;
///
/// let mut connection = Connection::new_session().await?;
///
/// let reply = connection
///     .call_method(
///         Some("org.freedesktop.DBus"),
///         "/org/freedesktop/DBus",
///         Some("org.freedesktop.DBus"),
///         "GetId",
///         &(),
///     )
///     .await?;
///
/// let id: &str = reply.body()?;
/// println!("Unique ID of the bus: {}", id);
///# Ok::<(), zbus::Error>(())
///# });
/// ```
///
/// #### Monitoring all messages
///
/// Let's eavesdrop on the session bus 😈 using the [Monitor] interface:
///
/// ```rust,no_run
///# futures::executor::block_on(async {
/// use futures::TryStreamExt;
/// use zbus::azync::Connection;
///
/// let mut connection = Connection::new_session().await?;
///
/// connection
///     .call_method(
///         Some("org.freedesktop.DBus"),
///         "/org/freedesktop/DBus",
///         Some("org.freedesktop.DBus.Monitoring"),
///         "BecomeMonitor",
///         &(&[] as &[&str], 0u32),
///     )
///     .await?;
///
/// while let Some(msg) = connection.try_next().await? {
///     println!("Got message: {}", msg);
/// }
///
///# Ok::<(), zbus::Error>(())
///# });
/// ```
///
/// This should print something like:
///
/// ```console
/// Got message: Signal NameAcquired from org.freedesktop.DBus
/// Got message: Signal NameLost from org.freedesktop.DBus
/// Got message: Method call GetConnectionUnixProcessID from :1.1324
/// Got message: Error org.freedesktop.DBus.Error.NameHasNoOwner:
///              Could not get PID of name ':1.1332': no such name from org.freedesktop.DBus
/// Got message: Method call AddMatch from :1.918
/// Got message: Method return from org.freedesktop.DBus
/// ```
///
/// [Monitor]: https://dbus.freedesktop.org/doc/dbus-specification.html#bus-messages-become-monitor
#[derive(Debug)]
pub struct Connection<S> {
    server_guid: Guid,
    cap_unix_fd: bool,
    unique_name: OnceCell<String>,

    raw_conn: RawConnection<Async<S>>,
    // Serial number for next outgoing message
    serial: u32,

    // Queue of incoming messages
    incoming_queue: Vec<Message>,

    // Max number of messages to queue
    max_queued: usize,
}

impl<S> Connection<S>
where
    S: AsRawFd + std::fmt::Debug + Unpin + Socket,
    Async<S>: Socket,
{
    /// Create and open a D-Bus connection from the given `stream`.
    ///
    /// The connection may either be set up for a *bus* connection, or not (for peer-to-peer
    /// communications).
    ///
    /// Upon successful return, the connection is fully established and negotiated: D-Bus messages
    /// can be sent and received.
    pub async fn new_client(stream: S, bus_connection: bool) -> Result<Self> {
        // SASL Handshake
        let auth = Authenticated::client(Async::new(stream)?).await?;

        if bus_connection {
            Connection::new_authenticated_bus(auth).await
        } else {
            Ok(Connection::new_authenticated(auth))
        }
    }

    /// Create a server `Connection` for the given `stream` and the server `guid`.
    ///
    /// The connection will wait for incoming client authentication handshake & negotiation messages,
    /// for peer-to-peer communications.
    ///
    /// Upon successful return, the connection is fully established and negotiated: D-Bus messages
    /// can be sent and received.
    pub async fn new_server(stream: S, guid: &Guid) -> Result<Self> {
        use nix::sys::socket::{getsockopt, sockopt::PeerCredentials};

        // FIXME: Could and should this be async?
        let creds = getsockopt(stream.as_raw_fd(), PeerCredentials)
            .map_err(|e| Error::Handshake(format!("Failed to get peer credentials: {}", e)))?;

        let auth = Authenticated::server(Async::new(stream)?, guid.clone(), creds.uid()).await?;

        Ok(Self::new_authenticated(auth))
    }

    /// Send `msg` to the peer.
    ///
    /// Unlike [`Sink`] implementation, this method sets a unique (to this connection) serial
    /// number on the message before sending it off, for you.
    ///
    /// On successfully sending off `msg`, the assigned serial number is returned.
    pub async fn send_message(&mut self, mut msg: Message) -> Result<u32> {
        let serial = self.assign_serial_num(&mut msg)?;

        self.send(msg).await?;

        Ok(serial)
    }

    /// Send a method call.
    ///
    /// Create a method-call message, send it over the connection, then wait for the reply.
    ///
    /// On succesful reply, an `Ok(Message)` is returned. On error, an `Err` is returned. D-Bus
    /// error replies are returned as [`Error::MethodError`].
    pub async fn call_method<B>(
        &mut self,
        destination: Option<&str>,
        path: &str,
        iface: Option<&str>,
        method_name: &str,
        body: &B,
    ) -> Result<Message>
    where
        B: serde::ser::Serialize + zvariant::Type,
    {
        let m = Message::method(
            self.unique_name(),
            destination,
            path,
            iface,
            method_name,
            body,
        )?;
        let serial = self.send_message(m).await?;

        let mut tmp_queue = vec![];

        while let Some(m) = self.try_next().await? {
            let h = m.header()?;

            if h.reply_serial()? != Some(serial) {
                if self.incoming_queue.len() + tmp_queue.len() < self.max_queued() {
                    // We first push to a temporary queue as otherwise it'll create an infinite loop
                    // since subsequent `receive_message` call will pick up the message from the main
                    // queue.
                    tmp_queue.push(m);
                }

                continue;
            } else {
                self.incoming_queue.append(&mut tmp_queue);
            }

            match h.message_type()? {
                MessageType::Error => return Err(m.into()),
                MessageType::MethodReturn => return Ok(m),
                _ => (),
            }
        }

        // If Stream gives us None, that means the socket was closed
        Err(Error::Io(io::Error::new(
            ErrorKind::BrokenPipe,
            "socket closed",
        )))
    }

    /// Emit a signal.
    ///
    /// Create a signal message, and send it over the connection.
    pub async fn emit_signal<B>(
        &mut self,
        destination: Option<&str>,
        path: &str,
        iface: &str,
        signal_name: &str,
        body: &B,
    ) -> Result<()>
    where
        B: serde::ser::Serialize + zvariant::Type,
    {
        let m = Message::signal(
            self.unique_name(),
            destination,
            path,
            iface,
            signal_name,
            body,
        )?;

        self.send_message(m).await.map(|_| ())
    }

    /// Reply to a message.
    ///
    /// Given an existing message (likely a method call), send a reply back to the caller with the
    /// given `body`.
    ///
    /// Returns the message serial number.
    pub async fn reply<B>(&mut self, call: &Message, body: &B) -> Result<u32>
    where
        B: serde::ser::Serialize + zvariant::Type,
    {
        let m = Message::method_reply(self.unique_name(), call, body)?;
        self.send_message(m).await
    }

    /// Reply an error to a message.
    ///
    /// Given an existing message (likely a method call), send an error reply back to the caller
    /// with the given `error_name` and `body`.
    ///
    /// Returns the message serial number.
    pub async fn reply_error<B>(
        &mut self,
        call: &Message,
        error_name: &str,
        body: &B,
    ) -> Result<u32>
    where
        B: serde::ser::Serialize + zvariant::Type,
    {
        let m = Message::method_error(self.unique_name(), call, error_name, body)?;
        self.send_message(m).await
    }

    /// Sets the unique name for this connection.
    ///
    /// This method should only be used when initializing a client *bus* connection with
    /// [`Connection::new_authenticated`]. Setting the unique name to anything other than the return
    /// value of the bus hello is a protocol violation.
    ///
    /// Returns and error if the name has already been set.
    pub fn set_unique_name(self, name: String) -> std::result::Result<Self, String> {
        self.unique_name.set(name).map(|_| self)
    }

    /// Assigns a serial number to `msg` that is unique to this connection.
    ///
    /// This method can fail if `msg` is corrupt.
    pub fn assign_serial_num(&mut self, msg: &mut Message) -> Result<u32> {
        let serial = self.next_serial();
        msg.modify_primary_header(|primary| {
            primary.set_serial_num(serial);

            Ok(())
        })?;

        Ok(serial)
    }

    /// The unique name as assigned by the message bus or `None` if not a message bus connection.
    pub fn unique_name(&self) -> Option<&str> {
        self.unique_name.get().map(|s| s.as_str())
    }

    /// Max number of messages to queue.
    pub fn max_queued(&self) -> usize {
        self.max_queued
    }

    /// Set the max number of messages to queue.
    ///
    /// Since typically you'd want to set this at instantiation time, this method takes ownership
    /// of `self` and returns an owned `Connection` instance so you can use the builder pattern to
    /// set the value.
    ///
    /// # Example
    ///
    /// ```
    ///# use std::error::Error;
    ///# use zbus::azync::Connection;
    /// use futures::executor::block_on;
    ///
    /// let conn = block_on(Connection::new_session())?.set_max_queued(30);
    /// assert_eq!(conn.max_queued(), 30);
    ///
    /// // Do something usefull with `conn`..
    ///# Ok::<_, Box<dyn Error + Send + Sync>>(())
    /// ```
    pub fn set_max_queued(mut self, max: usize) -> Self {
        self.max_queued = max;

        self
    }

    /// The server's GUID.
    pub fn server_guid(&self) -> &str {
        self.server_guid.as_str()
    }

    /// Create a `Connection` from an already authenticated unix socket.
    ///
    /// This method can be used in conjunction with [`crate::azync::Authenticated`] to handle
    /// the initial handshake of the D-Bus connection asynchronously.
    ///
    /// If the aim is to initialize a client *bus* connection, you need to send the client hello and assign
    /// the resulting unique name using [`set_unique_name`] before doing anything else.
    ///
    /// [`set_unique_name`]: struct.Connection.html#method.set_unique_name
    fn new_authenticated(auth: Authenticated<Async<S>>) -> Self {
        let auth = auth.into_inner();
        Self {
            raw_conn: auth.conn,
            server_guid: auth.server_guid,
            cap_unix_fd: auth.cap_unix_fd,
            serial: 1,
            unique_name: OnceCell::new(),
            incoming_queue: vec![],
            max_queued: DEFAULT_MAX_QUEUED,
        }
    }

    async fn new_authenticated_bus(auth: Authenticated<Async<S>>) -> Result<Self> {
        let mut connection = Connection::new_authenticated(auth);

        // Now that the server has approved us, we must send the bus Hello, as per specs
        // TODO: Use fdo module once it's async.
        let name: String = connection
            .call_method(
                Some("org.freedesktop.DBus"),
                "/org/freedesktop/DBus",
                Some("org.freedesktop.DBus"),
                "Hello",
                &(),
            )
            .await?
            .body()?;

        Ok(connection
            .set_unique_name(name)
            // programmer (probably our) error if this fails.
            .expect("Attempted to set unique_name twice"))
    }

    fn next_serial(&mut self) -> u32 {
        let serial = self.serial;
        self.serial = serial + 1;

        serial
    }

    // Used by Sink impl.
    fn flush(&mut self, cx: &mut Context<'_>) -> Poll<Result<()>> {
        loop {
            match self.raw_conn.try_flush() {
                Ok(()) => return Poll::Ready(Ok(())),
                Err(e) => {
                    if e.kind() == ErrorKind::WouldBlock {
                        let poll = self.raw_conn.socket().poll_writable(cx);

                        match poll {
                            Poll::Pending => return Poll::Pending,
                            // Guess socket became ready already so let's try it again.
                            Poll::Ready(Ok(_)) => continue,
                            Poll::Ready(Err(e)) => return Poll::Ready(Err(e.into())),
                        }
                    } else {
                        return Poll::Ready(Err(Error::Io(e)));
                    }
                }
            }
        }
    }
}

impl Connection<UnixStream> {
    /// Create a `Connection` to the session/user message bus.
    ///
    /// Although, session bus hardly ever runs on anything other than UNIX domain sockets, if you
    /// want your code to be able to handle those rare cases, use [`ConnectionType::new_session`]
    /// instead.
    pub async fn new_session() -> Result<Self> {
        Self::new_authenticated_bus(Authenticated::session().await?).await
    }

    /// Create a `Connection` to the system-wide message bus.
    ///
    /// Although, system bus hardly ever runs on anything other than UNIX domain sockets, if you
    /// want your code to be able to handle those rare cases, use [`ConnectionType::new_system`]
    /// instead.
    pub async fn new_system() -> Result<Self> {
        Self::new_authenticated_bus(Authenticated::system().await?).await
    }
}

impl<S> Sink<Message> for Connection<S>
where
    S: AsRawFd + std::fmt::Debug + Unpin + Socket,
    Async<S>: Socket,
{
    type Error = Error;

    fn poll_ready(self: Pin<&mut Self>, _cx: &mut Context<'_>) -> Poll<Result<()>> {
        // TODO: We should have a max queue length in raw::Socket for outgoing messages.
        Poll::Ready(Ok(()))
    }

    fn start_send(self: Pin<&mut Self>, msg: Message) -> Result<()> {
        let conn = self.get_mut();
        if !msg.fds().is_empty() && !conn.cap_unix_fd {
            return Err(Error::Unsupported);
        }

        conn.raw_conn.enqueue_message(msg);

        Ok(())
    }

    fn poll_flush(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Result<()>> {
        self.get_mut().flush(cx)
    }

    fn poll_close(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Result<()>> {
        let conn = self.get_mut();
        match conn.flush(cx) {
            Poll::Ready(Ok(_)) => (),
            Poll::Ready(Err(e)) => return Poll::Ready(Err(e)),
            Poll::Pending => return Poll::Pending,
        }

        Poll::Ready((conn.raw_conn).close())
    }
}

impl<S> Stream for Connection<S>
where
    S: Socket,
    Async<S>: Socket,
{
    type Item = Result<Message>;

    fn poll_next(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Option<Self::Item>> {
        let conn = self.get_mut();

        if let Some(msg) = conn.incoming_queue.pop() {
            return Poll::Ready(Some(Ok(msg)));
        }

        loop {
            match conn.raw_conn.try_receive_message() {
                Ok(m) => return Poll::Ready(Some(Ok(m))),
                Err(Error::Io(e)) if e.kind() == ErrorKind::WouldBlock => {
                    let poll = conn.raw_conn.socket().poll_readable(cx);

                    match poll {
                        Poll::Pending => return Poll::Pending,
                        // Guess socket became ready already so let's try it again.
                        Poll::Ready(Ok(_)) => continue,
                        Poll::Ready(Err(e)) => return Poll::Ready(Some(Err(e.into()))),
                    }
                }
                Err(Error::Io(e)) if e.kind() == ErrorKind::BrokenPipe => return Poll::Ready(None),
                Err(e) => return Poll::Ready(Some(Err(e))),
            }
        }
    }
}

/// Type representing all concrete [`Connection`] types, provided by zbus.
///
/// For maximum portability, use constructor method provided by this type instead of ones provided
/// by [`Connection`].
pub enum ConnectionType {
    Unix(Connection<UnixStream>),
}

impl ConnectionType {
    /// Create a `ConnectionType` for the given [D-Bus address].
    ///
    /// [D-Bus address]: https://dbus.freedesktop.org/doc/dbus-specification.html#addresses
    pub async fn new_for_address(address: &str, bus_connection: bool) -> Result<Self> {
        match AuthenticatedType::for_address(address).await? {
            AuthenticatedType::Unix(auth) => {
                let conn = if bus_connection {
                    Connection::new_authenticated_bus(auth).await?
                } else {
                    Connection::new_authenticated(auth)
                };

                Ok(ConnectionType::Unix(conn))
            }
        }
    }

    /// Create a `ConnectionType` to the session/user message bus.
    pub async fn new_session() -> Result<Self> {
        match AuthenticatedType::session().await? {
            AuthenticatedType::Unix(auth) => {
                let conn = Connection::new_authenticated_bus(auth).await?;

                Ok(ConnectionType::Unix(conn))
            }
        }
    }

    /// Create a `ConnectionType` to the system-wide message bus.
    pub async fn new_system() -> Result<Self> {
        match AuthenticatedType::system().await? {
            AuthenticatedType::Unix(auth) => {
                let conn = Connection::new_authenticated_bus(auth).await?;

                Ok(ConnectionType::Unix(conn))
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use std::os::unix::net::UnixStream;

    use super::*;

    #[test]
    fn unix_p2p() {
        futures::executor::block_on(test_unix_p2p()).unwrap();
    }

    async fn test_unix_p2p() -> Result<()> {
        let guid = Guid::generate();

        let (p0, p1) = UnixStream::pair().unwrap();

        let server = Connection::new_server(p0, &guid);
        let client = Connection::new_client(p1, false);

        let (mut client_conn, mut server_conn) = futures::try_join!(client, server)?;

        let server_future = async {
            let mut method: Option<Message> = None;
            while let Some(m) = server_conn.try_next().await? {
                if m.to_string() == "Method call Test" {
                    method.replace(m);

                    break;
                }
            }
            let method = method.unwrap();

            // Send another message first to check the queueing function on client side.
            server_conn
                .emit_signal(None, "/", "org.zbus.p2p", "ASignalForYou", &())
                .await?;
            server_conn.reply(&method, &("yay")).await
        };

        let client_future = async {
            let reply = client_conn
                .call_method(None, "/", Some("org.zbus.p2p"), "Test", &())
                .await?;
            assert_eq!(reply.to_string(), "Method return");
            // Check we didn't miss the signal that was sent during the call.
            let m = client_conn.try_next().await?.unwrap();
            assert_eq!(m.to_string(), "Signal ASignalForYou");
            reply.body::<String>().map_err(|e| e.into())
        };

        let (val, _) = futures::try_join!(client_future, server_future)?;
        assert_eq!(val, "yay");

        Ok(())
    }
}
