use libc;
#[cfg(not(target_os = "redox"))]
use nix::Error;
use nix::sys::signal::*;
use nix::unistd::*;
use std::convert::TryFrom;
use std::sync::atomic::{AtomicBool, Ordering};

#[test]
fn test_kill_none() {
    kill(getpid(), None).expect("Should be able to send signal to myself.");
}

#[test]
#[cfg(not(target_os = "fuchsia"))]
fn test_killpg_none() {
    killpg(getpgrp(), None)
        .expect("Should be able to send signal to my process group.");
}

#[test]
fn test_old_sigaction_flags() {
    let _m = crate::SIGNAL_MTX.lock().expect("Mutex got poisoned by another test");

    extern "C" fn handler(_: ::libc::c_int) {}
    let act = SigAction::new(
        SigHandler::Handler(handler),
        SaFlags::empty(),
        SigSet::empty(),
    );
    let oact = unsafe { sigaction(SIGINT, &act) }.unwrap();
    let _flags = oact.flags();
    let oact = unsafe { sigaction(SIGINT, &act) }.unwrap();
    let _flags = oact.flags();
}

#[test]
fn test_sigprocmask_noop() {
    sigprocmask(SigmaskHow::SIG_BLOCK, None, None)
        .expect("this should be an effective noop");
}

#[test]
fn test_sigprocmask() {
    let _m = crate::SIGNAL_MTX.lock().expect("Mutex got poisoned by another test");

    // This needs to be a signal that rust doesn't use in the test harness.
    const SIGNAL: Signal = Signal::SIGCHLD;

    let mut old_signal_set = SigSet::empty();
    sigprocmask(SigmaskHow::SIG_BLOCK, None, Some(&mut old_signal_set))
        .expect("expect to be able to retrieve old signals");

    // Make sure the old set doesn't contain the signal, otherwise the following
    // test don't make sense.
    assert_eq!(old_signal_set.contains(SIGNAL), false,
               "the {:?} signal is already blocked, please change to a \
                different one", SIGNAL);

    // Now block the signal.
    let mut signal_set = SigSet::empty();
    signal_set.add(SIGNAL);
    sigprocmask(SigmaskHow::SIG_BLOCK, Some(&signal_set), None)
        .expect("expect to be able to block signals");

    // And test it again, to make sure the change was effective.
    old_signal_set.clear();
    sigprocmask(SigmaskHow::SIG_BLOCK, None, Some(&mut old_signal_set))
        .expect("expect to be able to retrieve old signals");
    assert_eq!(old_signal_set.contains(SIGNAL), true,
               "expected the {:?} to be blocked", SIGNAL);

    // Reset the signal.
    sigprocmask(SigmaskHow::SIG_UNBLOCK, Some(&signal_set), None)
        .expect("expect to be able to block signals");
}

lazy_static! {
    static ref SIGNALED: AtomicBool = AtomicBool::new(false);
}

extern fn test_sigaction_handler(signal: libc::c_int) {
    let signal = Signal::try_from(signal).unwrap();
    SIGNALED.store(signal == Signal::SIGINT, Ordering::Relaxed);
}

#[cfg(not(target_os = "redox"))]
extern fn test_sigaction_action(_: libc::c_int, _: *mut libc::siginfo_t, _: *mut libc::c_void) {}

#[test]
#[cfg(not(target_os = "redox"))]
fn test_signal_sigaction() {
    let _m = crate::SIGNAL_MTX.lock().expect("Mutex got poisoned by another test");

    let action_handler = SigHandler::SigAction(test_sigaction_action);
    assert_eq!(unsafe { signal(Signal::SIGINT, action_handler) }.unwrap_err(), Error::UnsupportedOperation);
}

#[test]
fn test_signal() {
    let _m = crate::SIGNAL_MTX.lock().expect("Mutex got poisoned by another test");

    unsafe { signal(Signal::SIGINT, SigHandler::SigIgn) }.unwrap();
    raise(Signal::SIGINT).unwrap();
    assert_eq!(unsafe { signal(Signal::SIGINT, SigHandler::SigDfl) }.unwrap(), SigHandler::SigIgn);

    let handler = SigHandler::Handler(test_sigaction_handler);
    assert_eq!(unsafe { signal(Signal::SIGINT, handler) }.unwrap(), SigHandler::SigDfl);
    raise(Signal::SIGINT).unwrap();
    assert!(SIGNALED.load(Ordering::Relaxed));
    assert_eq!(unsafe { signal(Signal::SIGINT, SigHandler::SigDfl) }.unwrap(), handler);

    // Restore default signal handler
    unsafe { signal(Signal::SIGINT, SigHandler::SigDfl) }.unwrap();
}
