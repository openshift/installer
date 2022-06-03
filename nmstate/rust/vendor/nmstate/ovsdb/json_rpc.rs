use std::io::{Read, Write};
use std::os::unix::net::UnixStream;
use std::time::SystemTime;

use serde::{Deserialize, Serialize};
use serde_json::Value;

use crate::{ErrorKind, NmstateError};

const BUFFER_SIZE: usize = 4096;

#[derive(Debug)]
pub(crate) struct OvsDbJsonRpc {
    socket: UnixStream,
    transaction_id: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone, Default, PartialEq, Eq)]
struct OvsDbRpcRequest {
    method: String,
    params: Value,
    id: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone, Default, PartialEq, Eq)]
struct OvsDbRpcError {
    error: String,
    details: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone, Default, PartialEq, Eq)]
pub(crate) struct OvsDbRpcReply {
    // The result might also contain a error.
    result: Value,
    error: Option<OvsDbRpcError>,
    id: u64,
}

impl OvsDbJsonRpc {
    pub(crate) fn connect(socket_path: &str) -> Result<Self, NmstateError> {
        Ok(Self {
            socket: UnixStream::connect(socket_path).map_err(|e| {
                NmstateError::new(ErrorKind::Bug, format!("socket error {}", e))
            })?,
            transaction_id: get_sec_since_epoch(),
        })
    }

    pub(crate) fn exec(
        &mut self,
        method: &str,
        params: &Value,
    ) -> Result<Value, NmstateError> {
        self.transaction_id += 1;
        let req = OvsDbRpcRequest {
            method: method.to_string(),
            params: params.clone(),
            id: self.transaction_id,
        };
        let buffer = serde_json::to_string(&req)?;
        log::debug!("OVSDB: sending command {}", buffer);
        self.socket
            .write_all(buffer.as_bytes())
            .map_err(parse_socket_io_error)?;
        let reply = self.recv()?;
        if method == "transact" {
            check_transact_error(reply)
        } else {
            Ok(reply)
        }
    }

    fn recv(&mut self) -> Result<Value, NmstateError> {
        let mut response: Vec<u8> = Vec::new();
        loop {
            let mut buffer = [0u8; BUFFER_SIZE];
            let read = self
                .socket
                .read(&mut buffer)
                .map_err(parse_socket_io_error)?;
            log::debug!("OVSDB: recv data {:?}", &buffer[..read]);
            response.extend_from_slice(&buffer[..read]);
            if read < BUFFER_SIZE {
                break;
            }
        }
        let reply_string =
            String::from_utf8(response).map_err(parse_str_parse_error)?;
        log::debug!("OVSDB: recv string {:?}", &reply_string);
        let reply: OvsDbRpcReply = serde_json::from_str(&reply_string)?;
        if reply.id != self.transaction_id {
            let e = NmstateError::new(
                ErrorKind::PluginFailure,
                format!(
                    "Transaction ID mismatch for OVS DB JSON RPC: {:?}",
                    reply
                ),
            );
            log::error!("{}", e);
            Err(e)
        } else if let Some(rpc_error) = reply.error {
            let e = NmstateError::new(
                ErrorKind::PluginFailure,
                format!("OVS DB JSON RPC error: {:?}", rpc_error),
            );
            log::error!("{}", e);
            Err(e)
        } else {
            Ok(reply.result)
        }
    }
}

fn get_sec_since_epoch() -> u64 {
    match SystemTime::now().duration_since(SystemTime::UNIX_EPOCH) {
        Ok(d) => d.as_secs(),
        Err(_) => 0,
    }
}

fn parse_str_parse_error(e: std::string::FromUtf8Error) -> NmstateError {
    NmstateError::new(
        ErrorKind::PluginFailure,
        format!("Reply from OVSDB is not valid UTF-8 string: {}", e),
    )
}

fn parse_socket_io_error(e: std::io::Error) -> NmstateError {
    NmstateError::new(
        ErrorKind::PluginFailure,
        format!("OVSDB Socket error: {}", e),
    )
}

fn check_transact_error(reply: Value) -> Result<Value, NmstateError> {
    if let Some(trans_replies) = reply.as_array() {
        for trans_reply in trans_replies {
            if let Some(error_type) = trans_reply
                .as_object()
                .and_then(|r| r.get("error"))
                .and_then(|e| e.as_str())
            {
                let error_detail = trans_reply
                    .as_object()
                    .and_then(|r| r.get("details"))
                    .and_then(|d| d.as_str())
                    .unwrap_or("");
                let e = NmstateError::new(
                    ErrorKind::PluginFailure,
                    format!(
                        "OVS DB JSON RPC error {}: {}",
                        error_type, error_detail
                    ),
                );
                log::error!("{}", e);
                return Err(e);
            }
        }
    }
    Ok(reply)
}
