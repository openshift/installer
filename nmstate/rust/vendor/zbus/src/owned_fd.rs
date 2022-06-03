use std::{
    mem::forget,
    os::unix::io::{AsRawFd, FromRawFd, IntoRawFd, RawFd},
};

/// An owned representation of a file descriptor
///
/// When it is dropped, the underlying frile descriptor will be dropped.
/// You can take ownership of the file descriptor (and avoid it being closed)
/// by using the
/// [`IntoRawFd`](https://doc.rust-lang.org/stable/std/os/unix/io/trait.IntoRawFd.html)
/// implementation.
#[derive(Debug, PartialEq, Eq)]
pub struct OwnedFd {
    inner: RawFd,
}

impl FromRawFd for OwnedFd {
    unsafe fn from_raw_fd(fd: RawFd) -> OwnedFd {
        OwnedFd { inner: fd }
    }
}

impl AsRawFd for OwnedFd {
    fn as_raw_fd(&self) -> RawFd {
        self.inner
    }
}

impl IntoRawFd for OwnedFd {
    fn into_raw_fd(self) -> RawFd {
        let v = self.inner;
        forget(self);
        v
    }
}

impl Drop for OwnedFd {
    fn drop(&mut self) {
        let _ = nix::unistd::close(self.inner);
    }
}
