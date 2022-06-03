use async_io::Async;
use std::{
    io,
    os::unix::{
        io::{AsRawFd, FromRawFd, RawFd},
        net::UnixStream,
    },
};

use nix::{
    cmsg_space,
    sys::{
        socket::{recvmsg, sendmsg, ControlMessage, ControlMessageOwned, MsgFlags},
        uio::IoVec,
    },
};

use crate::{utils::FDS_MAX, OwnedFd};

/// Trait representing some transport layer over which the DBus protocol can be used
///
/// The crate provides an implementation of it for std's `UnixStream` on unix platforms.
/// You will want to implement this trait to integrate zbus with a async-runtime-aware
/// implementation of the socket, for example.
pub trait Socket {
    /// Whether this transport supports file descriptor passing
    const SUPPORTS_FD_PASSING: bool;

    /// Attempt to receive a message from the socket
    ///
    /// On success, returns the number of bytes read as well as a `Vec` containing
    /// any associated file descriptors.
    ///
    /// This method may return an error of kind `WouldBlock` instead if blocking for
    /// non-blocking sockets.
    fn recvmsg(&mut self, buffer: &mut [u8]) -> io::Result<(usize, Vec<OwnedFd>)>;

    /// Attempt to send a message on the socket
    ///
    /// On success, return the number of bytes written. There may be a partial write, in
    /// which case the caller is responsible of sending the remaining data by calling this
    /// method again until everything is written or it returns an error of kind `WouldBlock`.
    ///
    /// If at least one byte has been written, then all the provided file descriptors will
    /// have been sent as well, and should not be provided again in subsequent calls.
    /// If `Err(Errorkind::Wouldblock)`, none of the provided file descriptors were sent.
    ///
    /// If the underlying transport does not support transmitting file descriptors, this
    /// will return `Err(ErrorKind::InvalidInput)`.
    fn sendmsg(&mut self, buffer: &[u8], fds: &[RawFd]) -> io::Result<usize>;

    /// Close the socket.
    ///
    /// After this call, all reading and writing operations will fail.
    ///
    /// NB: All currently implementations don't block so this method will never return
    /// `Err(Errorkind::Wouldblock)`.
    fn close(&self) -> io::Result<()>;
}

impl Socket for UnixStream {
    const SUPPORTS_FD_PASSING: bool = true;

    fn recvmsg(&mut self, buffer: &mut [u8]) -> io::Result<(usize, Vec<OwnedFd>)> {
        let iov = [IoVec::from_mut_slice(buffer)];
        let mut cmsgspace = cmsg_space!([RawFd; FDS_MAX]);

        match recvmsg(
            self.as_raw_fd(),
            &iov,
            Some(&mut cmsgspace),
            MsgFlags::empty(),
        ) {
            Ok(msg) => {
                let mut fds = vec![];
                for cmsg in msg.cmsgs() {
                    if let ControlMessageOwned::ScmRights(fd) = cmsg {
                        fds.extend(fd.iter().map(|&f| unsafe { OwnedFd::from_raw_fd(f) }));
                    } else {
                        return Err(io::Error::new(
                            io::ErrorKind::InvalidData,
                            "unexpected CMSG kind",
                        ));
                    }
                }
                Ok((msg.bytes, fds))
            }
            Err(nix::Error::Sys(e)) => Err(e.into()),
            _ => Err(io::Error::new(io::ErrorKind::Other, "unhandled nix error")),
        }
    }

    fn sendmsg(&mut self, buffer: &[u8], fds: &[RawFd]) -> io::Result<usize> {
        let cmsg = if !fds.is_empty() {
            vec![ControlMessage::ScmRights(fds)]
        } else {
            vec![]
        };
        let iov = [IoVec::from_slice(buffer)];
        match sendmsg(self.as_raw_fd(), &iov, &cmsg, MsgFlags::empty(), None) {
            // can it really happen?
            Ok(0) => Err(io::Error::new(
                io::ErrorKind::WriteZero,
                "failed to write to buffer",
            )),
            Ok(n) => Ok(n),
            Err(nix::Error::Sys(e)) => Err(e.into()),
            _ => Err(io::Error::new(io::ErrorKind::Other, "unhandled nix error")),
        }
    }

    fn close(&self) -> io::Result<()> {
        self.shutdown(std::net::Shutdown::Both)
    }
}

impl<S> Socket for Async<S>
where
    S: Socket,
{
    const SUPPORTS_FD_PASSING: bool = true;

    fn recvmsg(&mut self, buffer: &mut [u8]) -> io::Result<(usize, Vec<OwnedFd>)> {
        self.get_mut().recvmsg(buffer)
    }

    fn sendmsg(&mut self, buffer: &[u8], fds: &[RawFd]) -> io::Result<usize> {
        self.get_mut().sendmsg(buffer, fds)
    }

    fn close(&self) -> io::Result<()> {
        self.get_ref().close()
    }
}
