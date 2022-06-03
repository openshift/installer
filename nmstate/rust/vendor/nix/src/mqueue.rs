//! Posix Message Queue functions
//!
//! # Example
//!
// no_run because a kernel module may be required.
//! ```no_run
//! # use std::ffi::CString;
//! # use nix::mqueue::*;
//! use nix::sys::stat::Mode;
//!
//! const MSG_SIZE: mq_attr_member_t = 32;
//! let mq_name= CString::new("/a_nix_test_queue").unwrap();
//!
//! let oflag0 = MQ_OFlag::O_CREAT | MQ_OFlag::O_WRONLY;
//! let mode = Mode::S_IWUSR | Mode::S_IRUSR | Mode::S_IRGRP | Mode::S_IROTH;
//! let mqd0 = mq_open(&mq_name, oflag0, mode, None).unwrap();
//! let msg_to_send = b"msg_1";
//! mq_send(&mqd0, msg_to_send, 1).unwrap();
//!
//! let oflag1 = MQ_OFlag::O_CREAT | MQ_OFlag::O_RDONLY;
//! let mqd1 = mq_open(&mq_name, oflag1, mode, None).unwrap();
//! let mut buf = [0u8; 32];
//! let mut prio = 0u32;
//! let len = mq_receive(&mqd1, &mut buf, &mut prio).unwrap();
//! assert_eq!(prio, 1);
//! assert_eq!(msg_to_send, &buf[0..len]);
//!
//! mq_close(mqd1).unwrap();
//! mq_close(mqd0).unwrap();
//! ```
//! [Further reading and details on the C API](https://man7.org/linux/man-pages/man7/mq_overview.7.html)

use crate::Result;
use crate::errno::Errno;

use libc::{self, c_char, mqd_t, size_t};
use std::ffi::CStr;
use crate::sys::stat::Mode;
use std::mem;

libc_bitflags!{
    pub struct MQ_OFlag: libc::c_int {
        O_RDONLY;
        O_WRONLY;
        O_RDWR;
        O_CREAT;
        O_EXCL;
        O_NONBLOCK;
        O_CLOEXEC;
    }
}

libc_bitflags!{
    pub struct FdFlag: libc::c_int {
        FD_CLOEXEC;
    }
}

#[repr(C)]
#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct MqAttr {
    mq_attr: libc::mq_attr,
}

/// Identifies an open POSIX Message Queue
// A safer wrapper around libc::mqd_t, which is a pointer on some platforms
// Deliberately is not Clone to prevent use-after-close scenarios
#[repr(transparent)]
#[derive(Debug)]
#[allow(missing_copy_implementations)]
pub struct MqdT(mqd_t);

// x32 compatibility
// See https://sourceware.org/bugzilla/show_bug.cgi?id=21279
#[cfg(all(target_arch = "x86_64", target_pointer_width = "32"))]
#[cfg_attr(docsrs, doc(cfg(all())))]
pub type mq_attr_member_t = i64;
#[cfg(not(all(target_arch = "x86_64", target_pointer_width = "32")))]
#[cfg_attr(docsrs, doc(cfg(all())))]
pub type mq_attr_member_t = libc::c_long;

impl MqAttr {
    pub fn new(mq_flags: mq_attr_member_t,
               mq_maxmsg: mq_attr_member_t,
               mq_msgsize: mq_attr_member_t,
               mq_curmsgs: mq_attr_member_t)
               -> MqAttr
    {
        let mut attr = mem::MaybeUninit::<libc::mq_attr>::uninit();
        unsafe {
            let p = attr.as_mut_ptr();
            (*p).mq_flags = mq_flags;
            (*p).mq_maxmsg = mq_maxmsg;
            (*p).mq_msgsize = mq_msgsize;
            (*p).mq_curmsgs = mq_curmsgs;
            MqAttr { mq_attr: attr.assume_init() }
        }
    }

    pub const fn flags(&self) -> mq_attr_member_t {
        self.mq_attr.mq_flags
    }

    /// The max number of messages that can be held by the queue
    pub const fn maxmsg(&self) -> mq_attr_member_t {
        self.mq_attr.mq_maxmsg
    }

    /// The maximum size of each message (in bytes)
    pub const fn msgsize(&self) -> mq_attr_member_t {
        self.mq_attr.mq_msgsize
    }

    /// The number of messages currently held in the queue
    pub const fn curmsgs(&self) -> mq_attr_member_t {
        self.mq_attr.mq_curmsgs
    }
}


/// Open a message queue
///
/// See also [`mq_open(2)`](https://pubs.opengroup.org/onlinepubs/9699919799/functions/mq_open.html)
// The mode.bits cast is only lossless on some OSes
#[allow(clippy::cast_lossless)]
pub fn mq_open(name: &CStr,
               oflag: MQ_OFlag,
               mode: Mode,
               attr: Option<&MqAttr>)
               -> Result<MqdT> {
    let res = match attr {
        Some(mq_attr) => unsafe {
            libc::mq_open(name.as_ptr(),
                          oflag.bits(),
                          mode.bits() as libc::c_int,
                          &mq_attr.mq_attr as *const libc::mq_attr)
        },
        None => unsafe { libc::mq_open(name.as_ptr(), oflag.bits()) },
    };
    Errno::result(res).map(MqdT)
}

/// Remove a message queue
///
/// See also [`mq_unlink(2)`](https://pubs.opengroup.org/onlinepubs/9699919799/functions/mq_unlink.html)
pub fn mq_unlink(name: &CStr) -> Result<()> {
    let res = unsafe { libc::mq_unlink(name.as_ptr()) };
    Errno::result(res).map(drop)
}

/// Close a message queue
///
/// See also [`mq_close(2)`](https://pubs.opengroup.org/onlinepubs/9699919799/functions/mq_close.html)
pub fn mq_close(mqdes: MqdT) -> Result<()> {
    let res = unsafe { libc::mq_close(mqdes.0) };
    Errno::result(res).map(drop)
}

/// Receive a message from a message queue
///
/// See also [`mq_receive(2)`](https://pubs.opengroup.org/onlinepubs/9699919799/functions/mq_receive.html)
pub fn mq_receive(mqdes: &MqdT, message: &mut [u8], msg_prio: &mut u32) -> Result<usize> {
    let len = message.len() as size_t;
    let res = unsafe {
        libc::mq_receive(mqdes.0,
                         message.as_mut_ptr() as *mut c_char,
                         len,
                         msg_prio as *mut u32)
    };
    Errno::result(res).map(|r| r as usize)
}

/// Send a message to a message queue
///
/// See also [`mq_send(2)`](https://pubs.opengroup.org/onlinepubs/9699919799/functions/mq_send.html)
pub fn mq_send(mqdes: &MqdT, message: &[u8], msq_prio: u32) -> Result<()> {
    let res = unsafe {
        libc::mq_send(mqdes.0,
                      message.as_ptr() as *const c_char,
                      message.len(),
                      msq_prio)
    };
    Errno::result(res).map(drop)
}

/// Get message queue attributes
///
/// See also [`mq_getattr(2)`](https://pubs.opengroup.org/onlinepubs/9699919799/functions/mq_getattr.html)
pub fn mq_getattr(mqd: &MqdT) -> Result<MqAttr> {
    let mut attr = mem::MaybeUninit::<libc::mq_attr>::uninit();
    let res = unsafe { libc::mq_getattr(mqd.0, attr.as_mut_ptr()) };
    Errno::result(res).map(|_| unsafe{MqAttr { mq_attr: attr.assume_init() }})
}

/// Set the attributes of the message queue. Only `O_NONBLOCK` can be set, everything else will be ignored
/// Returns the old attributes
/// It is recommend to use the `mq_set_nonblock()` and `mq_remove_nonblock()` convenience functions as they are easier to use
///
/// [Further reading](https://pubs.opengroup.org/onlinepubs/9699919799/functions/mq_setattr.html)
pub fn mq_setattr(mqd: &MqdT, newattr: &MqAttr) -> Result<MqAttr> {
    let mut attr = mem::MaybeUninit::<libc::mq_attr>::uninit();
    let res = unsafe {
        libc::mq_setattr(mqd.0, &newattr.mq_attr as *const libc::mq_attr, attr.as_mut_ptr())
    };
    Errno::result(res).map(|_| unsafe{ MqAttr { mq_attr: attr.assume_init() }})
}

/// Convenience function.
/// Sets the `O_NONBLOCK` attribute for a given message queue descriptor
/// Returns the old attributes
#[allow(clippy::useless_conversion)]    // Not useless on all OSes
pub fn mq_set_nonblock(mqd: &MqdT) -> Result<MqAttr> {
    let oldattr = mq_getattr(mqd)?;
    let newattr = MqAttr::new(mq_attr_member_t::from(MQ_OFlag::O_NONBLOCK.bits()),
                              oldattr.mq_attr.mq_maxmsg,
                              oldattr.mq_attr.mq_msgsize,
                              oldattr.mq_attr.mq_curmsgs);
    mq_setattr(mqd, &newattr)
}

/// Convenience function.
/// Removes `O_NONBLOCK` attribute for a given message queue descriptor
/// Returns the old attributes
pub fn mq_remove_nonblock(mqd: &MqdT) -> Result<MqAttr> {
    let oldattr = mq_getattr(mqd)?;
    let newattr = MqAttr::new(0,
                              oldattr.mq_attr.mq_maxmsg,
                              oldattr.mq_attr.mq_msgsize,
                              oldattr.mq_attr.mq_curmsgs);
    mq_setattr(mqd, &newattr)
}
