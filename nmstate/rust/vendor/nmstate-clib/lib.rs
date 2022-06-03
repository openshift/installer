// Copyright 2021 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

mod logger;

use std::ffi::{CStr, CString};
use std::time::SystemTime;

use libc::{c_char, c_int};
use nmstate::NmstateError;
use once_cell::sync::OnceCell;

use crate::logger::MemoryLogger;

const NMSTATE_FLAG_KERNEL_ONLY: u32 = 1 << 1;
const NMSTATE_FLAG_NO_VERIFY: u32 = 1 << 2;
const NMSTATE_FLAG_INCLUDE_STATUS_DATA: u32 = 1 << 3;
const NMSTATE_FLAG_INCLUDE_SECRETS: u32 = 1 << 4;
const NMSTATE_FLAG_NO_COMMIT: u32 = 1 << 5;
const NMSTATE_FLAG_MEMORY_ONLY: u32 = 1 << 6;
const NMSTATE_FLAG_RUNNING_CONFIG_ONLY: u32 = 1 << 7;

const NMSTATE_PASS: c_int = 0;
const NMSTATE_FAIL: c_int = 1;

static INSTANCE: OnceCell<MemoryLogger> = OnceCell::new();

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn nmstate_net_state_retrieve(
    flags: u32,
    state: *mut *mut c_char,
    log: *mut *mut c_char,
    err_kind: *mut *mut c_char,
    err_msg: *mut *mut c_char,
) -> c_int {
    assert!(!state.is_null());
    assert!(!log.is_null());
    assert!(!err_kind.is_null());
    assert!(!err_msg.is_null());

    unsafe {
        *log = std::ptr::null_mut();
        *state = std::ptr::null_mut();
        *err_kind = std::ptr::null_mut();
        *err_msg = std::ptr::null_mut();
    }

    let logger = match init_logger() {
        Ok(l) => l,
        Err(e) => {
            unsafe {
                *err_msg =
                    CString::new(format!("Failed to setup logger: {}", e))
                        .unwrap()
                        .into_raw();
            }
            return NMSTATE_FAIL;
        }
    };
    let now = SystemTime::now();

    let mut net_state = nmstate::NetworkState::new();
    if (flags & NMSTATE_FLAG_KERNEL_ONLY) > 0 {
        net_state.set_kernel_only(true);
    }

    if (flags & NMSTATE_FLAG_INCLUDE_STATUS_DATA) > 0 {
        net_state.set_include_status_data(true);
    }

    if (flags & NMSTATE_FLAG_INCLUDE_SECRETS) > 0 {
        net_state.set_include_secrets(true);
    }

    if (flags & NMSTATE_FLAG_RUNNING_CONFIG_ONLY) > 0 {
        net_state.set_running_config_only(true);
    }

    let result = net_state.retrieve();
    unsafe {
        *log = CString::new(logger.drain(now)).unwrap().into_raw();
    }

    match result {
        Ok(s) => match serde_json::to_string(&s) {
            Ok(state_str) => unsafe {
                *state = CString::new(state_str).unwrap().into_raw();
                NMSTATE_PASS
            },
            Err(e) => unsafe {
                *err_msg = CString::new(format!(
                    "serde_json::to_string failure: {}",
                    e
                ))
                .unwrap()
                .into_raw();
                *err_kind =
                    CString::new(format!("{}", nmstate::ErrorKind::Bug))
                        .unwrap()
                        .into_raw();
                NMSTATE_FAIL
            },
        },
        Err(e) => {
            unsafe {
                *err_msg = CString::new(e.msg()).unwrap().into_raw();
                *err_kind =
                    CString::new(format!("{}", &e.kind())).unwrap().into_raw();
            }
            NMSTATE_FAIL
        }
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn nmstate_net_state_apply(
    flags: u32,
    state: *const c_char,
    rollback_timeout: u32,
    log: *mut *mut c_char,
    err_kind: *mut *mut c_char,
    err_msg: *mut *mut c_char,
) -> c_int {
    assert!(!log.is_null());
    assert!(!err_kind.is_null());
    assert!(!err_msg.is_null());

    unsafe {
        *log = std::ptr::null_mut();
        *err_kind = std::ptr::null_mut();
        *err_msg = std::ptr::null_mut();
    }

    if state.is_null() {
        return NMSTATE_PASS;
    }

    let logger = match init_logger() {
        Ok(l) => l,
        Err(e) => {
            unsafe {
                *err_msg =
                    CString::new(format!("Failed to setup logger: {}", e))
                        .unwrap()
                        .into_raw();
            }
            return NMSTATE_FAIL;
        }
    };
    let now = SystemTime::now();

    let net_state_cstr = unsafe { CStr::from_ptr(state) };

    let net_state_str = match net_state_cstr.to_str() {
        Ok(s) => s,
        Err(e) => {
            unsafe {
                *err_msg = CString::new(format!(
                    "Error on converting C char to rust str: {}",
                    e
                ))
                .unwrap()
                .into_raw();
                *err_kind = CString::new(format!(
                    "{}",
                    nmstate::ErrorKind::InvalidArgument
                ))
                .unwrap()
                .into_raw();
            }
            return NMSTATE_FAIL;
        }
    };

    let mut net_state =
        match nmstate::NetworkState::new_from_json(net_state_str) {
            Ok(n) => n,
            Err(e) => {
                unsafe {
                    *err_msg = CString::new(e.msg()).unwrap().into_raw();
                    *err_kind = CString::new(format!("{}", &e.kind()))
                        .unwrap()
                        .into_raw();
                }
                return NMSTATE_FAIL;
            }
        };
    if (flags & NMSTATE_FLAG_KERNEL_ONLY) > 0 {
        net_state.set_kernel_only(true);
    }

    if (flags & NMSTATE_FLAG_NO_VERIFY) > 0 {
        net_state.set_verify_change(false);
    }

    if (flags & NMSTATE_FLAG_NO_COMMIT) > 0 {
        net_state.set_commit(false);
    }

    if (flags & NMSTATE_FLAG_MEMORY_ONLY) > 0 {
        net_state.set_memory_only(true);
    }

    net_state.set_timeout(rollback_timeout);

    let result = net_state.apply();
    unsafe {
        *log = CString::new(logger.drain(now)).unwrap().into_raw();
    }

    if let Err(e) = result {
        unsafe {
            *err_msg = CString::new(e.msg()).unwrap().into_raw();
            *err_kind =
                CString::new(format!("{}", &e.kind())).unwrap().into_raw();
        }
        NMSTATE_FAIL
    } else {
        NMSTATE_PASS
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn nmstate_checkpoint_commit(
    checkpoint: *const c_char,
    log: *mut *mut c_char,
    err_kind: *mut *mut c_char,
    err_msg: *mut *mut c_char,
) -> c_int {
    assert!(!log.is_null());
    assert!(!err_kind.is_null());
    assert!(!err_msg.is_null());

    unsafe {
        *log = std::ptr::null_mut();
        *err_kind = std::ptr::null_mut();
        *err_msg = std::ptr::null_mut();
    }

    let logger = match init_logger() {
        Ok(l) => l,
        Err(e) => {
            unsafe {
                *err_msg =
                    CString::new(format!("Failed to setup logger: {}", e))
                        .unwrap()
                        .into_raw();
            }
            return NMSTATE_FAIL;
        }
    };
    let now = SystemTime::now();

    let mut checkpoint_str = "";
    if !checkpoint.is_null() {
        let checkpoint_cstr = unsafe { CStr::from_ptr(checkpoint) };
        checkpoint_str = match checkpoint_cstr.to_str() {
            Ok(s) => s,
            Err(e) => {
                unsafe {
                    *err_msg = CString::new(format!(
                        "Error on converting C char to rust str: {}",
                        e
                    ))
                    .unwrap()
                    .into_raw();
                    *err_kind = CString::new(format!(
                        "{}",
                        nmstate::ErrorKind::InvalidArgument
                    ))
                    .unwrap()
                    .into_raw();
                }
                return NMSTATE_FAIL;
            }
        }
    }

    let result = nmstate::NetworkState::checkpoint_commit(checkpoint_str);
    unsafe {
        *log = CString::new(logger.drain(now)).unwrap().into_raw();
    }

    if let Err(e) = result {
        unsafe {
            *err_msg = CString::new(e.msg()).unwrap().into_raw();
            *err_kind =
                CString::new(format!("{}", &e.kind())).unwrap().into_raw();
        }
        NMSTATE_FAIL
    } else {
        NMSTATE_PASS
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn nmstate_checkpoint_rollback(
    checkpoint: *const c_char,
    log: *mut *mut c_char,
    err_kind: *mut *mut c_char,
    err_msg: *mut *mut c_char,
) -> c_int {
    assert!(!log.is_null());
    assert!(!err_kind.is_null());
    assert!(!err_msg.is_null());

    unsafe {
        *log = std::ptr::null_mut();
        *err_kind = std::ptr::null_mut();
        *err_msg = std::ptr::null_mut();
    }

    let logger = match init_logger() {
        Ok(l) => l,
        Err(e) => {
            unsafe {
                *err_msg =
                    CString::new(format!("Failed to setup logger: {}", e))
                        .unwrap()
                        .into_raw();
            }
            return NMSTATE_FAIL;
        }
    };
    let now = SystemTime::now();

    let mut checkpoint_str = "";
    if !checkpoint.is_null() {
        let checkpoint_cstr = unsafe { CStr::from_ptr(checkpoint) };
        checkpoint_str = match checkpoint_cstr.to_str() {
            Ok(s) => s,
            Err(e) => {
                unsafe {
                    *err_msg = CString::new(format!(
                        "Error on converting C char to rust str: {}",
                        e
                    ))
                    .unwrap()
                    .into_raw();
                    *err_kind = CString::new(format!(
                        "{}",
                        nmstate::ErrorKind::InvalidArgument
                    ))
                    .unwrap()
                    .into_raw();
                }
                return NMSTATE_FAIL;
            }
        }
    }

    // TODO: save log to the output pointer
    let result = nmstate::NetworkState::checkpoint_rollback(checkpoint_str);
    unsafe {
        *log = CString::new(logger.drain(now)).unwrap().into_raw();
    }

    if let Err(e) = result {
        unsafe {
            *err_msg = CString::new(e.msg()).unwrap().into_raw();
            *err_kind =
                CString::new(format!("{}", &e.kind())).unwrap().into_raw();
        }
        NMSTATE_FAIL
    } else {
        NMSTATE_PASS
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn nmstate_cstring_free(cstring: *mut c_char) {
    unsafe {
        if !cstring.is_null() {
            drop(CString::from_raw(cstring));
        }
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn nmstate_generate_configurations(
    state: *const c_char,
    configs: *mut *mut c_char,
    log: *mut *mut c_char,
    err_kind: *mut *mut c_char,
    err_msg: *mut *mut c_char,
) -> c_int {
    assert!(!state.is_null());
    assert!(!configs.is_null());
    assert!(!log.is_null());
    assert!(!err_kind.is_null());
    assert!(!err_msg.is_null());

    unsafe {
        *log = std::ptr::null_mut();
        *configs = std::ptr::null_mut();
        *err_kind = std::ptr::null_mut();
        *err_msg = std::ptr::null_mut();
    }

    if state.is_null() {
        return NMSTATE_PASS;
    }

    let logger = match init_logger() {
        Ok(l) => l,
        Err(e) => {
            unsafe {
                *err_msg =
                    CString::new(format!("Failed to setup logger: {}", e))
                        .unwrap()
                        .into_raw();
            }
            return NMSTATE_FAIL;
        }
    };
    let now = SystemTime::now();

    let net_state_cstr = unsafe { CStr::from_ptr(state) };

    let net_state_str = match net_state_cstr.to_str() {
        Ok(s) => s,
        Err(e) => {
            unsafe {
                *err_msg = CString::new(format!(
                    "Error on converting C char to rust str: {}",
                    e
                ))
                .unwrap()
                .into_raw();
                *err_kind = CString::new(format!(
                    "{}",
                    nmstate::ErrorKind::InvalidArgument
                ))
                .unwrap()
                .into_raw();
            }
            return NMSTATE_FAIL;
        }
    };

    let net_state = match nmstate::NetworkState::new_from_json(net_state_str) {
        Ok(n) => n,
        Err(e) => {
            unsafe {
                *err_msg = CString::new(e.msg()).unwrap().into_raw();
                *err_kind =
                    CString::new(format!("{}", &e.kind())).unwrap().into_raw();
            }
            return NMSTATE_FAIL;
        }
    };

    let result = net_state.gen_conf();
    unsafe {
        *log = CString::new(logger.drain(now)).unwrap().into_raw();
    }
    match result {
        Ok(s) => match serde_json::to_string(&s) {
            Ok(cfgs) => unsafe {
                *configs = CString::new(cfgs).unwrap().into_raw();
                NMSTATE_PASS
            },
            Err(e) => unsafe {
                *err_msg = CString::new(format!(
                    "serde_json::to_string failure: {}",
                    e
                ))
                .unwrap()
                .into_raw();
                *err_kind =
                    CString::new(format!("{}", nmstate::ErrorKind::Bug))
                        .unwrap()
                        .into_raw();
                NMSTATE_FAIL
            },
        },
        Err(e) => {
            unsafe {
                *err_msg = CString::new(e.msg()).unwrap().into_raw();
                *err_kind =
                    CString::new(format!("{}", &e.kind())).unwrap().into_raw();
            }
            NMSTATE_FAIL
        }
    }
}

fn init_logger() -> Result<&'static MemoryLogger, NmstateError> {
    match INSTANCE.get() {
        Some(l) => {
            l.add_consumer();
            Ok(l)
        }
        None => {
            if INSTANCE.set(MemoryLogger::new()).is_err() {
                return Err(NmstateError::new(
                    nmstate::ErrorKind::Bug,
                    "Failed to set once_sync for logger".to_string(),
                ));
            }
            if let Some(l) = INSTANCE.get() {
                if let Err(e) = log::set_logger(l) {
                    Err(NmstateError::new(
                        nmstate::ErrorKind::Bug,
                        format!("Failed to log::set_logger: {}", e),
                    ))
                } else {
                    l.add_consumer();
                    log::set_max_level(log::LevelFilter::Debug);
                    Ok(l)
                }
            } else {
                Err(NmstateError::new(
                    nmstate::ErrorKind::Bug,
                    "Failed to get logger from once_sync".to_string(),
                ))
            }
        }
    }
}
