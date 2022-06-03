use std::{
    borrow::Cow,
    collections::HashMap,
    convert::{TryFrom, TryInto},
    sync::Mutex,
};
use zvariant::{ObjectPath, OwnedValue, Value};

use crate::{Connection, Error, Message, Result};

use crate::fdo::{self, IntrospectableProxy, PropertiesProxy};

const LOCK_FAIL_MSG: &str = "Failed to lock a mutex or read-write lock";

type SignalHandler = Box<dyn FnMut(&Message) -> Result<()> + Send>;

/// A client-side interface proxy.
///
/// A `Proxy` is a helper to interact with an interface on a remote object.
///
/// # Example
///
/// ```
/// use std::result::Result;
/// use std::error::Error;
/// use zbus::{Connection, Proxy};
///
/// fn main() -> Result<(), Box<dyn Error>> {
///     let connection = Connection::new_session()?;
///     let p = Proxy::new(
///         &connection,
///         "org.freedesktop.DBus",
///         "/org/freedesktop/DBus",
///         "org.freedesktop.DBus",
///     )?;
///     // owned return value
///     let _id: String = p.call("GetId", &())?;
///     // borrowed return value
///     let _id: &str = p.call_method("GetId", &())?.body()?;
///     Ok(())
/// }
/// ```
///
/// # Note
///
/// It is recommended to use the [`dbus_proxy`] macro, which provides a more convenient and
/// type-safe *façade* `Proxy` derived from a Rust trait.
///
/// ## Current limitations:
///
/// At the moment, `Proxy` doesn't:
///
/// * cache properties
/// * track the current name owner
/// * prevent auto-launching
///
/// [`dbus_proxy`]: attr.dbus_proxy.html
pub struct Proxy<'a> {
    conn: Connection,
    destination: Cow<'a, str>,
    path: Cow<'a, str>,
    interface: Cow<'a, str>,
    sig_handlers: Mutex<HashMap<&'static str, SignalHandler>>,
}

impl<'a> Proxy<'a> {
    /// Create a new `Proxy` for the given destination/path/interface.
    pub fn new(
        conn: &Connection,
        destination: &'a str,
        path: &'a str,
        interface: &'a str,
    ) -> Result<Self> {
        Ok(Self {
            conn: conn.clone(),
            destination: Cow::from(destination),
            path: Cow::from(path),
            interface: Cow::from(interface),
            sig_handlers: Mutex::new(HashMap::new()),
        })
    }

    /// Create a new `Proxy` for the given destination/path/interface, taking ownership of all
    /// passed arguments.
    pub fn new_owned(
        conn: Connection,
        destination: String,
        path: String,
        interface: String,
    ) -> Result<Self> {
        Ok(Self {
            conn,
            destination: Cow::from(destination),
            path: Cow::from(path),
            interface: Cow::from(interface),
            sig_handlers: Mutex::new(HashMap::new()),
        })
    }

    /// Get a reference to the associated connection.
    pub fn connection(&self) -> &Connection {
        &self.conn
    }

    /// Get a reference to the destination service name.
    pub fn destination(&self) -> &str {
        &self.destination
    }

    /// Get a reference to the object path.
    pub fn path(&self) -> &str {
        &self.path
    }

    /// Get a reference to the interface.
    pub fn interface(&self) -> &str {
        &self.interface
    }

    /// Introspect the associated object, and return the XML description.
    ///
    /// See the [xml](xml/index.html) module for parsing the result.
    pub fn introspect(&self) -> fdo::Result<String> {
        IntrospectableProxy::new_for(&self.conn, &self.destination, &self.path)?.introspect()
    }

    /// Get the property `property_name`.
    ///
    /// Effectively, call the `Get` method of the `org.freedesktop.DBus.Properties` interface.
    pub fn get_property<T>(&self, property_name: &str) -> fdo::Result<T>
    where
        T: TryFrom<OwnedValue>,
    {
        PropertiesProxy::new_for(&self.conn, &self.destination, &self.path)?
            .get(&self.interface, property_name)?
            .try_into()
            .map_err(|_| Error::InvalidReply.into())
    }

    /// Set the property `property_name`.
    ///
    /// Effectively, call the `Set` method of the `org.freedesktop.DBus.Properties` interface.
    pub fn set_property<'t, T: 't>(&self, property_name: &str, value: T) -> fdo::Result<()>
    where
        T: Into<Value<'t>>,
    {
        PropertiesProxy::new_for(&self.conn, &self.destination, &self.path)?.set(
            &self.interface,
            property_name,
            &value.into(),
        )
    }

    /// Call a method and return the reply.
    ///
    /// Typically, you would want to use [`call`] method instead. Use this method if you need to
    /// deserialize the reply message manually (this way, you can avoid the memory
    /// allocation/copying, by deserializing the reply to an unowned type).
    ///
    /// [`call`]: struct.Proxy.html#method.call
    pub fn call_method<B>(&self, method_name: &str, body: &B) -> Result<Message>
    where
        B: serde::ser::Serialize + zvariant::Type,
    {
        let reply = self.conn.call_method(
            Some(&self.destination),
            &self.path,
            Some(&self.interface),
            method_name,
            body,
        );
        match reply {
            Ok(mut reply) => {
                reply.disown_fds();

                Ok(reply)
            }
            Err(e) => Err(e),
        }
    }

    /// Call a method and return the reply body.
    ///
    /// Use [`call_method`] instead if you need to deserialize the reply manually/separately.
    ///
    /// [`call_method`]: struct.Proxy.html#method.call_method
    pub fn call<B, R>(&self, method_name: &str, body: &B) -> Result<R>
    where
        B: serde::ser::Serialize + zvariant::Type,
        R: serde::de::DeserializeOwned + zvariant::Type,
    {
        Ok(self.call_method(method_name, body)?.body()?)
    }

    /// Register a handler for signal named `signal_name`.
    ///
    /// Once a handler is successfully registered, call [`Self::next_signal`] to wait for the next
    /// signal to arrive and be handled by its registered handler.
    ///
    /// If the associated connnection is to a bus, a match rule is added for the signal on the bus
    /// so that the bus sends us the signals.
    ///
    /// ### Errors
    ///
    /// This method can fail if addition of the relevant match rule on the bus fails. You can
    /// safely `unwrap` the `Result` if you're certain that associated connnection is not a bus
    /// connection.
    pub fn connect_signal<H>(&self, signal_name: &'static str, handler: H) -> fdo::Result<()>
    where
        H: FnMut(&Message) -> Result<()> + Send + 'static,
    {
        if self
            .sig_handlers
            .lock()
            .expect(LOCK_FAIL_MSG)
            .insert(signal_name, Box::new(handler))
            .is_none()
            && self.conn.is_bus()
        {
            let rule = self.match_rule_for_signal(signal_name);
            fdo::DBusProxy::new(&self.conn)?.add_match(&rule)?;
        }

        Ok(())
    }

    /// Deregister the handler for the signal named `signal_name`.
    ///
    /// If the associated connnection is to a bus, the match rule is removed for the signal on the
    /// bus so that the bus stops sending us the signal. This method returns `Ok(true)` if a
    /// handler was registered for `signal_name` and was removed by this call; `Ok(false)`
    /// otherwise.
    ///
    /// ### Errors
    ///
    /// This method can fail if removal of the relevant match rule on the bus fails. You can
    /// safely `unwrap` the `Result` if you're certain that associated connnection is not a bus
    /// connection.
    pub fn disconnect_signal(&self, signal_name: &'static str) -> fdo::Result<bool> {
        if self
            .sig_handlers
            .lock()
            .expect(LOCK_FAIL_MSG)
            .remove(signal_name)
            .is_some()
            && self.conn.is_bus()
        {
            let rule = self.match_rule_for_signal(signal_name);
            fdo::DBusProxy::new(&self.conn)?.remove_match(&rule)?;

            Ok(true)
        } else {
            Ok(false)
        }
    }

    /// Receive and handle the next incoming signal on the associated connection.
    ///
    /// This method will wait for signal messages on the associated connection and call any
    /// handlers registered through the [`Self::connect_signal`] method. Signal handlers can be
    /// registered and deregistered from another threads during the call to this method.
    ///
    /// If the signal message was handled by a handler, `Ok(None)` is returned. Otherwise, the
    /// received message is returned.
    pub fn next_signal(&self) -> Result<Option<Message>> {
        let msg = self.conn.receive_specific(|msg| {
            let handlers = self.sig_handlers.lock().expect(LOCK_FAIL_MSG);
            if handlers.is_empty() {
                // No signal handers associated anymore so no need to continue.
                return Ok(true);
            }

            let hdr = msg.header()?;

            let member = match hdr.member()? {
                Some(m) => m,
                None => return Ok(false),
            };

            Ok(hdr.interface()? == Some(&self.interface)
                && hdr.path()? == Some(&ObjectPath::try_from(self.path.as_ref())?)
                && hdr.message_type()? == crate::MessageType::Signal
                && handlers.contains_key(member))
        })?;

        if self.handle_signal(&msg)? {
            Ok(None)
        } else {
            Ok(Some(msg))
        }
    }

    /// Handle the provided signal message.
    ///
    /// Call any handlers registered through the [`Self::connect_signal`] method for the provided
    /// signal message.
    ///
    /// If no errors are encountered, `Ok(true)` is returned if a handler was found and called for,
    /// the signal; `Ok(false)` otherwise.
    pub fn handle_signal(&self, msg: &Message) -> Result<bool> {
        let mut handlers = self.sig_handlers.lock().expect(LOCK_FAIL_MSG);
        if handlers.is_empty() {
            return Ok(false);
        }

        let hdr = msg.header()?;
        if let Some(name) = hdr.member()? {
            if let Some(handler) = handlers.get_mut(name) {
                handler(&msg)?;

                return Ok(true);
            }
        }

        Ok(false)
    }

    pub(crate) fn has_signal_handler(&self, signal_name: &str) -> bool {
        self.sig_handlers
            .lock()
            .expect(LOCK_FAIL_MSG)
            .contains_key(signal_name)
    }

    fn match_rule_for_signal(&self, signal_name: &'static str) -> String {
        // FIXME: Use the API to create this once we've it (issue#69).
        format!(
            "type='signal',path_namespace='{}',interface='{}',member='{}'",
            self.path, self.interface, signal_name,
        )
    }
}

impl<'asref, 'p: 'asref> std::convert::AsRef<Proxy<'asref>> for Proxy<'p> {
    fn as_ref(&self) -> &Proxy<'asref> {
        &self
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::Arc;

    #[test]
    fn signal() {
        // Register a well-known name with the session bus and ensure we get the appropriate
        // signals called for that.
        let conn = Connection::new_session().unwrap();
        let owner_change_signaled = Arc::new(Mutex::new(false));
        let name_acquired_signaled = Arc::new(Mutex::new(false));

        let proxy = Proxy::new(
            &conn,
            "org.freedesktop.DBus",
            "/org/freedesktop/DBus",
            "org.freedesktop.DBus",
        )
        .unwrap();

        let well_known = "org.freedesktop.zbus.ProxySignalTest";
        let unique_name = conn.unique_name().unwrap().to_string();
        {
            let well_known = well_known.clone();
            let signaled = owner_change_signaled.clone();
            proxy
                .connect_signal("NameOwnerChanged", move |m| {
                    let (name, _, new_owner) = m.body::<(&str, &str, &str)>()?;
                    if name != well_known {
                        // Meant for the other testcase then
                        return Ok(());
                    }
                    assert_eq!(new_owner, unique_name);
                    *signaled.lock().unwrap() = true;

                    Ok(())
                })
                .unwrap();
        }
        {
            let signaled = name_acquired_signaled.clone();
            // `NameAcquired` is emitted twice, first when the unique name is assigned on
            // connection and secondly after we ask for a specific name.
            proxy
                .connect_signal("NameAcquired", move |m| {
                    if m.body::<&str>()? == well_known {
                        *signaled.lock().unwrap() = true;
                    }

                    Ok(())
                })
                .unwrap();
        }

        fdo::DBusProxy::new(&conn)
            .unwrap()
            .request_name(&well_known, fdo::RequestNameFlags::ReplaceExisting.into())
            .unwrap();

        loop {
            proxy.next_signal().unwrap();

            if *owner_change_signaled.lock().unwrap() && *name_acquired_signaled.lock().unwrap() {
                break;
            }
        }
    }
}
