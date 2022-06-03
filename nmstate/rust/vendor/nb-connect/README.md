# nb-connect (deprecated)

[![Build](https://github.com/smol-rs/nb-connect/workflows/Build%20and%20test/badge.svg)](
https://github.com/smol-rs/nb-connect/actions)
[![License](https://img.shields.io/badge/license-Apache--2.0_OR_MIT-blue.svg)](
https://github.com/smol-rs/nb-connect)
[![Cargo](https://img.shields.io/crates/v/nb-connect.svg)](
https://crates.io/crates/nb-connect)
[![Documentation](https://docs.rs/nb-connect/badge.svg)](
https://docs.rs/nb-connect)

**This crate is now deprecated in favor of [socket2](https://crates.io/crates/socket2).**

Non-blocking TCP or Unix connect.

This crate allows you to create a [`TcpStream`] or a [`UnixStream`] in a non-blocking way,
without waiting for the connection to become fully established.

[`TcpStream`]: https://doc.rust-lang.org/stable/std/net/struct.TcpStream.html
[`UnixStream`]: https://doc.rust-lang.org/stable/std/os/unix/net/struct.UnixStream.html

## Examples

```rust
use polling::{Event, Poller};
use std::time::Duration;

// Create a pending TCP connection.
let stream = nb_connect::tcp(([127, 0, 0, 1], 80))?;

// Create a poller that waits for the stream to become writable.
let poller = Poller::new()?;
poller.add(&stream, Event::writable(0))?;

// Wait for at most 1 second.
if poller.wait(&mut Vec::new(), Some(Duration::from_secs(1)))? == 0 {
    println!("timeout");
} else if let Some(err) = stream.take_error()? {
    println!("error: {}", err);
} else {
    println!("connected");
}
```

## License

Licensed under either of

 * Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
 * MIT license ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)

at your option.

#### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall be
dual licensed as above, without any additional terms or conditions.
