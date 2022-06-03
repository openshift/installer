// This is based on the work of https://github.com/gahag/memory_logger
// which is MIT licensed.

use std::sync::{
    atomic::{AtomicU16, Ordering},
    Mutex,
};
use std::time::SystemTime;

use serde::ser::{Serialize, SerializeMap, Serializer};

const INITIAL_VEC_CAPACITY: usize = 256;

#[derive(Debug, PartialEq, Eq, Clone)]
pub(crate) struct LogEntry {
    time: SystemTime,
    level: log::Level,
    file: String,
    msg: String,
}

impl Default for LogEntry {
    fn default() -> Self {
        Self {
            time: SystemTime::UNIX_EPOCH,
            level: log::Level::Debug,
            file: String::new(),
            msg: String::new(),
        }
    }
}

impl Serialize for LogEntry {
    // Serialize is also used for verification.
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut map = serializer.serialize_map(Some(4))?;

        map.serialize_entry(
            "time",
            &format!(
                "{}",
                self.time
                    .duration_since(SystemTime::UNIX_EPOCH)
                    .unwrap_or_default()
                    .as_secs()
            ),
        )?;
        map.serialize_entry("level", self.level.as_str())?;
        map.serialize_entry("file", self.file.as_str())?;
        map.serialize_entry("msg", self.msg.as_str())?;
        map.end()
    }
}

impl From<&log::Record<'_>> for LogEntry {
    fn from(r: &log::Record<'_>) -> Self {
        Self {
            time: SystemTime::now(),
            level: r.level(),
            file: format!(
                "{}:{}:{}",
                r.module_path().unwrap_or_default(),
                r.file().unwrap_or_default(),
                r.line().unwrap_or_default()
            ),
            msg: r.args().to_string(),
        }
    }
}

/// A blocking memory logger. Logging and read operations may block.
///
/// You should have only a single instance of this in your program.
#[derive(Default, Debug)]
pub(crate) struct MemoryLogger {
    consumer_count: AtomicU16,
    logs: Mutex<Vec<LogEntry>>,
}

impl log::Log for MemoryLogger {
    fn enabled(&self, metadata: &log::Metadata) -> bool {
        metadata.target().starts_with("nmstate::")
            || metadata.target().starts_with("nispor::")
    }

    fn log(&self, record: &log::Record) {
        if self.enabled(record.metadata()) {
            let mut logs = self.logs.lock().expect("inner lock poisoned");
            logs.push(LogEntry::from(record))
        }
    }

    fn flush(&self) {}
}

impl MemoryLogger {
    pub(crate) fn new() -> Self {
        Self {
            consumer_count: AtomicU16::new(0),
            logs: Mutex::new(Vec::with_capacity(INITIAL_VEC_CAPACITY)),
        }
    }

    pub(crate) fn add_consumer(&self) {
        self.consumer_count.fetch_add(1, Ordering::SeqCst);
    }

    /// Drain the whole buffered log, only return logs since specified time.
    /// Note that this locks the logger, causing logging to block.
    pub(crate) fn drain(&self, since: SystemTime) -> String {
        let mut logs = self.logs.lock().expect("inner lock poisoned");
        let ret = serde_json::to_string(
            &logs
                .as_slice()
                .iter()
                .filter(|l| l.time >= since)
                .collect::<Vec<&LogEntry>>(),
        )
        .unwrap_or_default();
        if self.consumer_count.fetch_sub(1, Ordering::SeqCst) == 1 {
            logs.clear();
        }
        ret
    }
}
