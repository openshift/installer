// SPDX-License-Identifier: MIT

#[macro_use]
extern crate thiserror;

mod connection;
mod error;
mod handle;
pub mod message;
mod resolver;

#[cfg(feature = "tokio_socket")]
pub use connection::new_connection;
pub use connection::new_connection_with_socket;
pub use error::GenetlinkError;
pub use handle::GenetlinkHandle;
