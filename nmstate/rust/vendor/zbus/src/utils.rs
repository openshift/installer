use nix::{
    errno::Errno,
    poll::{poll, PollFd, PollFlags},
};
use std::os::unix::io::RawFd;

pub(crate) const FDS_MAX: usize = 1024; // this is hardcoded in sdbus - nothing in the spec

pub(crate) fn padding_for_8_bytes(value: usize) -> usize {
    padding_for_n_bytes(value, 8)
}

pub(crate) fn padding_for_n_bytes(value: usize, align: usize) -> usize {
    let len_rounded_up = value.wrapping_add(align).wrapping_sub(1) & !align.wrapping_sub(1);

    len_rounded_up.wrapping_sub(value)
}

pub(crate) fn wait_on(fd: RawFd, flags: PollFlags) -> std::io::Result<()> {
    let pollfd = PollFd::new(fd, flags);
    loop {
        match poll(&mut [pollfd], -1) {
            Ok(_) => break,
            Err(nix::Error::Sys(e)) => {
                if e == Errno::EAGAIN || e == Errno::EINTR {
                    // we got interupted, try polling again
                    continue;
                } else {
                    return Err(std::io::Error::from(e));
                }
            }
            _ => {
                return Err(std::io::Error::new(
                    std::io::ErrorKind::Other,
                    "unhandled nix error",
                ))
            }
        }
    }
    Ok(())
}
