use crate::{Connection, Error, MessageHeader, Proxy, Result};
use std::{
    borrow::Cow,
    collections::HashMap,
    convert::{AsRef, TryFrom},
};
use zvariant::ObjectPath;

#[derive(Hash, Eq, PartialEq)]
struct ProxyKey<'key> {
    interface: Cow<'key, str>,
    path: ObjectPath<'key>,
}

impl<'p, P> From<&P> for ProxyKey<'_>
where
    P: AsRef<Proxy<'p>>,
{
    fn from(proxy: &P) -> Self {
        let proxy = proxy.as_ref();
        ProxyKey {
            interface: Cow::from(proxy.interface().to_owned()),
            path: ObjectPath::try_from(proxy.path())
                .expect("invalid object path")
                .to_owned(),
        }
    }
}

impl<'key> TryFrom<&'key MessageHeader<'_>> for ProxyKey<'key> {
    type Error = Error;

    fn try_from(hdr: &'key MessageHeader<'_>) -> Result<Self> {
        match (hdr.interface()?, hdr.path()?.cloned()) {
            (Some(interface), Some(path)) => Ok(ProxyKey {
                interface: Cow::from(interface),
                path,
            }),
            (_, _) => Err(Error::Message(crate::MessageError::MissingField)),
        }
    }
}

/// Receives signals for [`Proxy`] instances.
///
/// Use this to receive signals on a given connection for a bunch of proxies at the same time.
pub struct SignalReceiver<'r, 'p> {
    conn: Connection,
    proxies: HashMap<ProxyKey<'static>, &'r Proxy<'p>>,
}

impl<'r, 'p> SignalReceiver<'r, 'p> {
    /// Create a new `SignalReceiver` instance.
    pub fn new(conn: Connection) -> Self {
        Self {
            conn,
            proxies: HashMap::new(),
        }
    }

    /// Get a reference to the associated connection.
    pub fn connection(&self) -> &Connection {
        &self.conn
    }

    /// Get a iterator for all the proxies in this receiver.
    pub fn proxies(&self) -> impl Iterator<Item = &&Proxy<'_>> {
        self.proxies.values()
    }

    /// Watch for signals relevant to the `proxy`.
    ///
    /// # Panics
    ///
    /// This method will panic if you try to add a proxy with a different associated connection than
    /// the one associated with this receiver.
    pub fn receive_for<'a: 'p, 'b: 'r, P>(&mut self, proxy: &'b P)
    where
        P: AsRef<Proxy<'a>>,
    {
        let proxy = proxy.as_ref();
        assert_eq!(proxy.connection().unique_name(), self.conn.unique_name());

        let key = ProxyKey::from(proxy);
        self.proxies.insert(key, proxy);
    }

    /// Received and handle the next incoming signal on the associated connection.
    ///
    /// This method will wait for signal messages on the associated connection and call any
    /// handler registered (through [`Proxy::connect_signal`]) with a proxy in this receiver.
    ///
    /// If the signal message was handled by a handler, `Ok(None)` is returned. Otherwise, the
    /// received message is returned.
    pub fn next_signal(&self) -> Result<Option<crate::Message>> {
        let msg = self.conn.receive_specific(|msg| {
            let hdr = msg.header()?;

            if hdr.message_type()? != crate::MessageType::Signal {
                return Ok(false);
            }
            let member = match hdr.member()? {
                Some(m) => m,
                None => return Ok(false),
            };

            let key = ProxyKey::try_from(&hdr)?;

            if let Some(proxy) = self.proxies.get(&key) {
                if proxy.has_signal_handler(member) {
                    return Ok(true);
                }
            }

            Ok(false)
        })?;

        if self.handle_signal(&msg)? {
            Ok(None)
        } else {
            Ok(Some(msg))
        }
    }

    /// Handle the provided signal message.
    ///
    /// Call any handler registered (through [`Proxy::connect_signal`]) with a proxy in this receiver.
    ///
    /// If no errors are encountered, `Ok(true)` is returned if a handler was found and called for,
    /// the signal; `Ok(false)` otherwise.
    pub fn handle_signal(&self, msg: &crate::Message) -> Result<bool> {
        let hdr = msg.header()?;
        let key = ProxyKey::try_from(&hdr)?;

        match self.proxies.get(&key) {
            Some(proxy) => proxy.handle_signal(&msg),
            None => Ok(false),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::{dbus_interface, dbus_proxy, fdo};
    use std::{
        cell::RefCell,
        convert::TryInto,
        rc::Rc,
        sync::{Arc, Mutex},
    };

    fn multiple_signal_iface_test() -> std::result::Result<u32, Box<dyn std::error::Error>> {
        #[dbus_proxy(interface = "org.freedesktop.zbus.MultiSignal")]
        trait MultiSignal {
            #[dbus_proxy(signal)]
            fn some_signal(&self, sig_arg: &str) -> Result<()>;

            fn emit_it(&self, arg: &str) -> Result<()>;
        }
        let conn = Connection::new_session()?;
        let mut receiver = SignalReceiver::new(conn.clone());

        let proxy1 = MultiSignalProxy::new_for(
            &conn,
            "org.freedesktop.zbus.MultiSignal",
            "/org/freedesktop/zbus/MultiSignal/1",
        )?;
        let proxy1_str = Arc::new(Mutex::new(None));
        let clone = proxy1_str.clone();
        proxy1.connect_some_signal(move |s| {
            *clone.lock().unwrap() = Some(s.to_string());

            Ok(())
        })?;
        receiver.receive_for(&proxy1);

        let proxy2 = MultiSignalProxy::new_for(
            &conn,
            "org.freedesktop.zbus.MultiSignal",
            "/org/freedesktop/zbus/MultiSignal/2",
        )?;
        let proxy2_str = Arc::new(Mutex::new(None));
        let clone = proxy2_str.clone();
        proxy2.connect_some_signal(move |s| {
            *clone.lock().unwrap() = Some(s.to_string());

            Ok(())
        })?;
        receiver.receive_for(&proxy2);

        proxy1.emit_it("hi")?;
        proxy2.emit_it("bye")?;

        loop {
            receiver.next_signal()?;
            if proxy1_str.lock().unwrap().is_some() && proxy2_str.lock().unwrap().is_some() {
                break;
            }
        }

        Ok(99)
    }

    #[test]
    fn multiple_proxy_signals() {
        struct MultiSignal {
            times_called: Rc<RefCell<u8>>,
        }

        #[dbus_interface(interface = "org.freedesktop.zbus.MultiSignal")]
        impl MultiSignal {
            #[dbus_interface(signal)]
            fn some_signal(&self, sig_arg: &str) -> Result<()>;

            fn emit_it(&mut self, arg: &str) -> Result<()> {
                *self.times_called.borrow_mut() += 1;
                self.some_signal(arg)
            }
        }

        let conn = Connection::new_session().unwrap();
        fdo::DBusProxy::new(&conn)
            .unwrap()
            .request_name(
                "org.freedesktop.zbus.MultiSignal",
                fdo::RequestNameFlags::ReplaceExisting.into(),
            )
            .unwrap();
        let mut object_server = crate::ObjectServer::new(&conn);
        let times_called = Rc::new(RefCell::new(0));
        let iface = MultiSignal {
            times_called: times_called.clone(),
        };
        object_server
            .at(
                &"/org/freedesktop/zbus/MultiSignal/1".try_into().unwrap(),
                iface,
            )
            .unwrap();
        let iface = MultiSignal {
            times_called: times_called.clone(),
        };
        object_server
            .at(
                &"/org/freedesktop/zbus/MultiSignal/2".try_into().unwrap(),
                iface,
            )
            .unwrap();

        let child = std::thread::spawn(|| multiple_signal_iface_test().unwrap());

        while *times_called.borrow() < 2 {
            object_server.try_handle_next().unwrap();
        }

        let val = child.join().expect("failed to join");
        assert_eq!(val, 99);
    }
}
