//! Futures, streams, and async I/O combinators.
//!
//! This crate is a subset of [futures] that compiles an order of magnitude faster, fixes minor
//! warts in its API, fills in some obvious gaps, and removes almost all unsafe code from it.
//!
//! In short, this crate aims to be more enjoyable than [futures] but still fully compatible with
//! it.
//!
//! [futures]: https://docs.rs/futures
//!
//! # Examples
//!
#![cfg_attr(feature = "std", doc = "```no_run")]
#![cfg_attr(not(feature = "std"), doc = "```ignore")]
//! use futures_lite::future;
//!
//! fn main() {
//!     future::block_on(async {
//!         println!("Hello world!");
//!     })
//! }
//! ```

#![warn(missing_docs, missing_debug_implementations, rust_2018_idioms)]
#![cfg_attr(not(feature = "std"), no_std)]

// TODO: These hidden re-exports are deprecated and should eventually be removed.
#[cfg(feature = "std")]
#[doc(hidden)]
pub use crate::io::{
    AsyncBufRead, AsyncBufReadExt, AsyncRead, AsyncReadExt, AsyncSeek, AsyncSeekExt, AsyncWrite,
    AsyncWriteExt,
};
#[doc(hidden)]
pub use crate::{
    future::{Future, FutureExt},
    stream::{Stream, StreamExt},
};

pub mod future;
pub mod prelude;
pub mod stream;

#[cfg(feature = "std")]
pub mod io;

/// Unwraps `Poll<T>` or returns [`Pending`][`core::task::Poll::Pending`].
///
/// # Examples
///
/// ```
/// use futures_lite::{future, prelude::*, ready};
/// use std::pin::Pin;
/// use std::task::{Context, Poll};
///
/// fn do_poll(cx: &mut Context<'_>) -> Poll<()> {
///     let mut fut = future::ready(42);
///     let fut = Pin::new(&mut fut);
///
///     let num = ready!(fut.poll(cx));
///     # drop(num);
///     // ... use num
///
///     Poll::Ready(())
/// }
/// ```
///
/// The `ready!` call expands to:
///
/// ```
/// # use futures_lite::{future, prelude::*, ready};
/// # use std::pin::Pin;
/// # use std::task::{Context, Poll};
/// #
/// # fn do_poll(cx: &mut Context<'_>) -> Poll<()> {
///     # let mut fut = future::ready(42);
///     # let fut = Pin::new(&mut fut);
///     #
/// let num = match fut.poll(cx) {
///     Poll::Ready(t) => t,
///     Poll::Pending => return Poll::Pending,
/// };
///     # drop(num);
///     # // ... use num
///     #
///     # Poll::Ready(())
/// # }
/// ```
#[macro_export]
macro_rules! ready {
    ($e:expr $(,)?) => {
        match $e {
            core::task::Poll::Ready(t) => t,
            core::task::Poll::Pending => return core::task::Poll::Pending,
        }
    };
}

/// Pins a variable of type `T` on the stack and rebinds it as `Pin<&mut T>`.
///
/// ```
/// use futures_lite::{future, pin};
/// use std::fmt::Debug;
/// use std::future::Future;
/// use std::pin::Pin;
/// use std::time::Instant;
///
/// // Inspects each invocation of `Future::poll()`.
/// async fn inspect<T: Debug>(f: impl Future<Output = T>) -> T {
///     pin!(f);
///     future::poll_fn(|cx| dbg!(f.as_mut().poll(cx))).await
/// }
///
/// # spin_on::spin_on(async {
/// let f = async { 1 + 2 };
/// inspect(f).await;
/// # })
/// ```
#[macro_export]
macro_rules! pin {
    ($($x:ident),* $(,)?) => {
        $(
            let mut $x = $x;
            #[allow(unused_mut)]
            let mut $x = unsafe {
                core::pin::Pin::new_unchecked(&mut $x)
            };
        )*
    }
}
