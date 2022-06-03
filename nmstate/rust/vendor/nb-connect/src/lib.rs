//! Non-blocking TCP or Unix connect.
//!
//! This crate allows you to create a [`TcpStream`] or a [`UnixStream`] in a non-blocking way,
//! without waiting for the connection to become fully established.
//!
//! [`TcpStream`]: https://doc.rust-lang.org/stable/std/net/struct.TcpStream.html
//! [`UnixStream`]: https://doc.rust-lang.org/stable/std/os/unix/net/struct.UnixStream.html
//!
//! # Examples
//!
//! ```
//! use polling::{Event, Poller};
//! use std::time::Duration;
//!
//! // Create a pending TCP connection.
//! let stream = nb_connect::tcp(([127, 0, 0, 1], 80))?;
//!
//! // Create a poller that waits for the stream to become writable.
//! let poller = Poller::new()?;
//! poller.add(&stream, Event::writable(0))?;
//!
//! // Wait for at most 1 second.
//! if poller.wait(&mut Vec::new(), Some(Duration::from_secs(1)))? == 0 {
//!     println!("timeout");
//! } else if let Some(err) = stream.take_error()? {
//!     println!("error: {}", err);
//! } else {
//!     println!("connected");
//! }
//! # std::io::Result::Ok(())
//! ```

#![warn(missing_docs, missing_debug_implementations, rust_2018_idioms)]
#![deprecated(
    since = "1.2.0",
    note = "This crate is now deprecated in favor of [socket2](https://crates.io/crates/socket2)."
)]

use std::io;
use std::net::{SocketAddr, TcpStream};

use socket2::{Domain, Protocol, SockAddr, Socket, Type};

#[cfg(unix)]
use std::{os::unix::net::UnixStream, path::Path};

fn connect(addr: SockAddr, domain: Domain, protocol: Option<Protocol>) -> io::Result<Socket> {
    let sock_type = Type::STREAM;
    #[cfg(any(
        target_os = "android",
        target_os = "dragonfly",
        target_os = "freebsd",
        target_os = "fuchsia",
        target_os = "illumos",
        target_os = "linux",
        target_os = "netbsd",
        target_os = "openbsd"
    ))]
    // If we can, set nonblocking at socket creation for unix
    let sock_type = sock_type.nonblocking();
    // This automatically handles cloexec on unix, no_inherit on windows and nosigpipe on macos
    let socket = Socket::new(domain, sock_type, protocol)?;
    #[cfg(not(any(
        target_os = "android",
        target_os = "dragonfly",
        target_os = "freebsd",
        target_os = "fuchsia",
        target_os = "illumos",
        target_os = "linux",
        target_os = "netbsd",
        target_os = "openbsd"
    )))]
    // If the current platform doesn't support nonblocking at creation, enable it after creation
    socket.set_nonblocking(true)?;
    match socket.connect(&addr) {
        Ok(_) => {}
        #[cfg(unix)]
        Err(err) if err.raw_os_error() == Some(libc::EINPROGRESS) => {}
        Err(err) if err.kind() == io::ErrorKind::WouldBlock => {}
        Err(err) => return Err(err),
    }
    Ok(socket)
}

/// Creates a pending Unix connection to the specified path.
///
/// The returned Unix stream will be in non-blocking mode and in the process of connecting to the
/// specified path.
///
/// The stream becomes writable when connected.
///
/// # Examples
///
/// ```no_run
/// use polling::{Event, Poller};
/// use std::time::Duration;
///
/// // Create a pending Unix connection.
/// let stream = nb_connect::unix("/tmp/socket")?;
///
/// // Create a poller that waits for the stream to become writable.
/// let poller = Poller::new()?;
/// poller.add(&stream, Event::writable(0))?;
///
/// // Wait for at most 1 second.
/// if poller.wait(&mut Vec::new(), Some(Duration::from_secs(1)))? == 0 {
///     println!("timeout");
/// } else {
///     println!("connected");
/// }
/// # std::io::Result::Ok(())
/// ```
#[cfg(unix)]
pub fn unix<P: AsRef<Path>>(path: P) -> io::Result<UnixStream> {
    let socket = connect(SockAddr::unix(path)?, Domain::UNIX, None)?;
    Ok(socket.into())
}

/// Creates a pending TCP connection to the specified address.
///
/// The returned TCP stream will be in non-blocking mode and in the process of connecting to the
/// specified address.
///
/// The stream becomes writable when connected.
///
/// # Examples
///
/// ```
/// use polling::{Event, Poller};
/// use std::time::Duration;
///
/// // Create a pending TCP connection.
/// let stream = nb_connect::tcp(([127, 0, 0, 1], 80))?;
///
/// // Create a poller that waits for the stream to become writable.
/// let poller = Poller::new()?;
/// poller.add(&stream, Event::writable(0))?;
///
/// // Wait for at most 1 second.
/// if poller.wait(&mut Vec::new(), Some(Duration::from_secs(1)))? == 0 {
///     println!("timeout");
/// } else if let Some(err) = stream.take_error()? {
///     println!("error: {}", err);
/// } else {
///     println!("connected");
/// }
/// # std::io::Result::Ok(())
/// ```
pub fn tcp<A: Into<SocketAddr>>(addr: A) -> io::Result<TcpStream> {
    let addr = addr.into();
    let domain = Domain::for_address(addr);
    let socket = connect(addr.into(), domain, Some(Protocol::TCP))?;
    Ok(socket.into())
}
