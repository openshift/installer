use std::{
    any::{Any, TypeId},
    cell::RefCell,
    collections::{hash_map::Entry, HashMap},
    fmt::Write,
    marker::PhantomData,
    rc::Rc,
};

use scoped_tls::scoped_thread_local;
use zvariant::{ObjectPath, OwnedValue, Value};

use crate::{dbus_interface, fdo, Connection, Error, Message, MessageHeader, MessageType, Result};

scoped_thread_local!(static LOCAL_NODE: Node);
scoped_thread_local!(static LOCAL_CONNECTION: Connection);

/// The trait used to dispatch messages to an interface instance.
///
/// Note: It is not recommended to manually implement this trait. The [`dbus_interface`] macro
/// implements it for you.
///
/// [`dbus_interface`]: attr.dbus_interface.html
pub trait Interface: Any {
    /// Return the name of the interface. Ex: "org.foo.MyInterface"
    fn name() -> &'static str
    where
        Self: Sized;

    /// Get a property value. Returns `None` if the property doesn't exist.
    fn get(&self, property_name: &str) -> Option<fdo::Result<OwnedValue>>;

    /// Return all the properties.
    fn get_all(&self) -> HashMap<String, OwnedValue>;

    /// Set a property value. Returns `None` if the property doesn't exist.
    fn set(&mut self, property_name: &str, value: &Value<'_>) -> Option<fdo::Result<()>>;

    /// Call a `&self` method. Returns `None` if the method doesn't exist.
    fn call(&self, connection: &Connection, msg: &Message, name: &str) -> Option<Result<u32>>;

    /// Call a `&mut self` method. Returns `None` if the method doesn't exist.
    fn call_mut(
        &mut self,
        connection: &Connection,
        msg: &Message,
        name: &str,
    ) -> Option<Result<u32>>;

    /// Write introspection XML to the writer, with the given indentation level.
    fn introspect_to_writer(&self, writer: &mut dyn Write, level: usize);
}

impl dyn Interface {
    /// Return Any of self
    fn downcast_ref<T: Any>(&self) -> Option<&T> {
        if <dyn Interface as Any>::type_id(self) == TypeId::of::<T>() {
            // SAFETY: If type ID matches, it means object is of type T
            Some(unsafe { &*(self as *const dyn Interface as *const T) })
        } else {
            None
        }
    }
}

struct Introspectable;

#[dbus_interface(name = "org.freedesktop.DBus.Introspectable")]
impl Introspectable {
    fn introspect(&self) -> String {
        LOCAL_NODE.with(|node| node.introspect())
    }
}

struct Peer;

#[dbus_interface(name = "org.freedesktop.DBus.Peer")]
impl Peer {
    fn ping(&self) {}

    fn get_machine_id(&self) -> fdo::Result<String> {
        let mut id = match std::fs::read_to_string("/var/lib/dbus/machine-id") {
            Ok(id) => id,
            Err(e) => {
                if let Ok(id) = std::fs::read_to_string("/etc/machine-id") {
                    id
                } else {
                    return Err(fdo::Error::IOError(format!(
                        "Failed to read from /var/lib/dbus/machine-id or /etc/machine-id: {}",
                        e
                    )));
                }
            }
        };

        let len = id.trim_end().len();
        id.truncate(len);
        Ok(id)
    }
}

struct Properties;

#[dbus_interface(name = "org.freedesktop.DBus.Properties")]
impl Properties {
    fn get(&self, interface_name: &str, property_name: &str) -> fdo::Result<OwnedValue> {
        LOCAL_NODE.with(|node| {
            let iface = node.get_interface(interface_name).ok_or_else(|| {
                fdo::Error::UnknownInterface(format!("Unknown interface '{}'", interface_name))
            })?;

            let res = iface.borrow().get(property_name);
            res.ok_or_else(|| {
                fdo::Error::UnknownProperty(format!("Unknown property '{}'", property_name))
            })?
        })
    }

    // TODO: should be able to take a &Value instead (but obscure deserialize error for now..)
    fn set(
        &mut self,
        interface_name: &str,
        property_name: &str,
        value: OwnedValue,
    ) -> fdo::Result<()> {
        LOCAL_NODE.with(|node| {
            let iface = node.get_interface(interface_name).ok_or_else(|| {
                fdo::Error::UnknownInterface(format!("Unknown interface '{}'", interface_name))
            })?;

            let res = iface.borrow_mut().set(property_name, &value);
            res.ok_or_else(|| {
                fdo::Error::UnknownProperty(format!("Unknown property '{}'", property_name))
            })?
        })
    }

    fn get_all(&self, interface_name: &str) -> fdo::Result<HashMap<String, OwnedValue>> {
        LOCAL_NODE.with(|node| {
            let iface = node.get_interface(interface_name).ok_or_else(|| {
                fdo::Error::UnknownInterface(format!("Unknown interface '{}'", interface_name))
            })?;

            let res = iface.borrow().get_all();
            Ok(res)
        })
    }

    #[dbus_interface(signal)]
    fn properties_changed(
        &self,
        interface_name: &str,
        changed_properties: &HashMap<&str, &Value<'_>>,
        invalidated_properties: &[&str],
    ) -> Result<()>;
}

#[derive(Default, derivative::Derivative)]
#[derivative(Debug)]
struct Node {
    path: String,
    children: HashMap<String, Node>,
    #[derivative(Debug = "ignore")]
    interfaces: HashMap<&'static str, Rc<RefCell<dyn Interface>>>,
}

impl Node {
    fn new(path: &str) -> Self {
        let mut node = Self {
            path: path.to_string(),
            ..Default::default()
        };
        node.at(Peer::name(), Peer);
        node.at(Introspectable::name(), Introspectable);
        node.at(Properties::name(), Properties);

        node
    }

    fn get_interface(&self, iface: &str) -> Option<Rc<RefCell<dyn Interface>>> {
        self.interfaces.get(iface).cloned()
    }

    fn remove_interface(&mut self, iface: &str) -> bool {
        self.interfaces.remove(iface).is_some()
    }

    fn is_empty(&self) -> bool {
        self.interfaces
            .keys()
            .find(|k| {
                *k != &Peer::name() && *k != &Introspectable::name() && *k != &Properties::name()
            })
            .is_none()
    }

    fn remove_node(&mut self, node: &str) -> bool {
        self.children.remove(node).is_some()
    }

    fn at<I>(&mut self, name: &'static str, iface: I) -> bool
    where
        I: Interface,
    {
        match self.interfaces.entry(name) {
            Entry::Vacant(e) => e.insert(Rc::new(RefCell::new(iface))),
            Entry::Occupied(_) => return false,
        };

        true
    }

    fn with_iface_func<F, I>(&self, func: F) -> Result<()>
    where
        F: Fn(&I) -> Result<()>,
        I: Interface,
    {
        let iface = self
            .interfaces
            .get(I::name())
            .ok_or(Error::InterfaceNotFound)?
            .borrow();
        let iface = iface.downcast_ref::<I>().ok_or(Error::InterfaceNotFound)?;
        func(iface)
    }

    fn introspect_to_writer<W: Write>(&self, writer: &mut W, level: usize) {
        if level == 0 {
            writeln!(
                writer,
                r#"
<!DOCTYPE node PUBLIC "-//freedesktop//DTD D-BUS Object Introspection 1.0//EN"
 "http://www.freedesktop.org/standards/dbus/1.0/introspect.dtd">
<node>"#
            )
            .unwrap();
        }

        for iface in self.interfaces.values() {
            iface.borrow().introspect_to_writer(writer, level + 2);
        }

        for (path, node) in &self.children {
            let level = level + 2;
            writeln!(
                writer,
                "{:indent$}<node name=\"{}\">",
                "",
                path,
                indent = level
            )
            .unwrap();
            node.introspect_to_writer(writer, level);
            writeln!(writer, "{:indent$}</node>", "", indent = level).unwrap();
        }

        if level == 0 {
            writeln!(writer, "</node>").unwrap();
        }
    }

    fn introspect(&self) -> String {
        let mut xml = String::with_capacity(1024);

        self.introspect_to_writer(&mut xml, 0);

        xml
    }

    fn emit_signal<B>(
        &self,
        dest: Option<&str>,
        iface: &str,
        signal_name: &str,
        body: &B,
    ) -> Result<()>
    where
        B: serde::ser::Serialize + zvariant::Type,
    {
        if !LOCAL_CONNECTION.is_set() {
            panic!("emit_signal: Connection TLS not set");
        }

        LOCAL_CONNECTION.with(|conn| conn.emit_signal(dest, &self.path, iface, signal_name, body))
    }
}

/// An object server, holding server-side D-Bus objects & interfaces.
///
/// Object servers hold interfaces on various object paths, and expose them over D-Bus.
///
/// All object paths will have the standard interfaces implemented on your behalf, such as
/// `org.freedesktop.DBus.Introspectable` or `org.freedesktop.DBus.Properties`.
///
/// **NOTE:** The lifetime `'a` on the `ObjectServer` struct is bogus and only exists for
/// backwards-compatibility and will be dropped in the next major release (i-e 2.0). This means
/// that `'a` can be considered `'static` for all intents and purposes.
///
/// # Example
///
/// This example exposes the `org.myiface.Example.Quit` method on the `/org/zbus/path`
/// path.
///
/// ```no_run
///# use std::error::Error;
///# use std::convert::TryInto;
/// use zbus::{Connection, ObjectServer, dbus_interface};
/// use std::rc::Rc;
/// use std::cell::RefCell;
///
/// struct Example {
///     // Interfaces are owned by the ObjectServer. They can have
///     // `&mut self` methods.
///     //
///     // If you need a shared state, you can use a RefCell for ex:
///     quit: Rc<RefCell<bool>>,
/// }
///
/// impl Example {
///     fn new(quit: Rc<RefCell<bool>>) -> Self {
///         Self { quit }
///     }
/// }
///
/// #[dbus_interface(name = "org.myiface.Example")]
/// impl Example {
///     // This will be the "Quit" D-Bus method.
///     fn quit(&self) {
///         *self.quit.borrow_mut() = true;
///     }
///
///     // See `dbus_interface` documentation to learn
///     // how to expose properties & signals as well.
/// }
///
/// let connection = Connection::new_session()?;
/// let mut object_server = ObjectServer::new(&connection);
/// let quit = Rc::new(RefCell::new(false));
///
/// let interface = Example::new(quit.clone());
/// object_server.at(&"/org/zbus/path".try_into()?, interface)?;
///
/// loop {
///     if let Err(err) = object_server.try_handle_next() {
///         eprintln!("{}", err);
///     }
///
///     if *quit.borrow() {
///         break;
///     }
/// }
///# Ok::<_, Box<dyn Error + Send + Sync>>(())
/// ```
#[derive(Debug)]
pub struct ObjectServer<'a> {
    conn: Connection,
    root: Node,
    phantom: PhantomData<&'a ()>,
}

impl<'a> ObjectServer<'a> {
    /// Creates a new D-Bus `ObjectServer` for a given connection.
    pub fn new(connection: &Connection) -> Self {
        Self {
            conn: connection.clone(),
            root: Node::new("/"),
            phantom: PhantomData,
        }
    }

    // Get the Node at path.
    fn get_node(&self, path: &ObjectPath<'_>) -> Option<&Node> {
        let mut node = &self.root;
        let mut node_path = String::new();

        for i in path.split('/').skip(1) {
            if i.is_empty() {
                continue;
            }
            write!(&mut node_path, "/{}", i).unwrap();
            match node.children.get(i) {
                Some(n) => node = n,
                None => return None,
            }
        }

        Some(node)
    }

    // Get the Node at path. Optionally create one if it doesn't exist.
    fn get_node_mut(&mut self, path: &ObjectPath<'_>, create: bool) -> Option<&mut Node> {
        let mut node = &mut self.root;
        let mut node_path = String::new();

        for i in path.split('/').skip(1) {
            if i.is_empty() {
                continue;
            }
            write!(&mut node_path, "/{}", i).unwrap();
            match node.children.entry(i.into()) {
                Entry::Vacant(e) => {
                    if create {
                        node = e.insert(Node::new(&node_path));
                    } else {
                        return None;
                    }
                }
                Entry::Occupied(e) => node = e.into_mut(),
            }
        }

        Some(node)
    }

    /// Register a D-Bus [`Interface`] at a given path. (see the example above)
    ///
    /// If the interface already exists at this path, returns false.
    ///
    /// [`Interface`]: trait.Interface.html
    pub fn at<I>(&mut self, path: &ObjectPath<'_>, iface: I) -> Result<bool>
    where
        I: Interface,
    {
        Ok(self.get_node_mut(path, true).unwrap().at(I::name(), iface))
    }

    /// Unregister a D-Bus [`Interface`] at a given path.
    ///
    /// If there are no more interfaces left at that path, destroys the object as well.
    /// Returns whether the object was destroyed.
    ///
    /// [`Interface`]: trait.Interface.html
    pub fn remove<I>(&mut self, path: &ObjectPath<'_>) -> Result<bool>
    where
        I: Interface,
    {
        let node = self
            .get_node_mut(path, false)
            .ok_or(Error::InterfaceNotFound)?;
        if !node.remove_interface(I::name()) {
            return Err(Error::InterfaceNotFound);
        }
        if node.is_empty() {
            let mut path_parts = path.rsplit('/').filter(|i| !i.is_empty());
            let last_part = path_parts.next().unwrap();
            let ppath = ObjectPath::from_string_unchecked(
                path_parts.fold(String::new(), |a, p| format!("/{}{}", p, a)),
            );
            self.get_node_mut(&ppath, false)
                .unwrap()
                .remove_node(last_part);
            return Ok(true);
        }
        Ok(false)
    }

    /// Run `func` with the given path & interface.
    ///
    /// Run the function `func` with the interface at path. If the interface was not found, return
    /// `Error::InterfaceNotFound`.
    ///
    /// This function is useful to emit signals outside of a dispatched handler:
    /// ```no_run
    ///# use std::error::Error;
    ///# use std::convert::TryInto;
    ///# use zbus::{Connection, ObjectServer, dbus_interface};
    ///
    ///# struct MyIface;
    ///# #[dbus_interface(name = "org.myiface.MyIface")]
    ///# impl MyIface {
    ///#     #[dbus_interface(signal)]
    ///#     fn emit_signal(&self) -> zbus::Result<()>;
    ///# }
    ///#
    ///# let connection = Connection::new_session()?;
    ///# let mut object_server = ObjectServer::new(&connection);
    ///#
    ///# let path = &"/org/zbus/path".try_into()?;
    ///# object_server.at(path, MyIface)?;
    /// object_server.with(path, |iface: &MyIface| {
    ///   iface.emit_signal()
    /// })?;
    ///#
    ///#
    ///# Ok::<_, Box<dyn Error + Send + Sync>>(())
    /// ```
    pub fn with<F, I>(&self, path: &ObjectPath<'_>, func: F) -> Result<()>
    where
        F: Fn(&I) -> Result<()>,
        I: Interface,
    {
        let node = self.get_node(path).ok_or(Error::InterfaceNotFound)?;
        LOCAL_CONNECTION.set(&self.conn, || {
            LOCAL_NODE.set(node, || node.with_iface_func(func))
        })
    }

    /// Emit a signal on the currently dispatched node.
    ///
    /// This is an internal helper function to emit a signal on on the current node. You shouldn't
    /// call this method directly, rather with the derived signal implementation from
    /// [`dbus_interface`].
    ///
    /// # Panics
    ///
    /// This method will panic if called from outside of a node context. Use [`ObjectServer::with`]
    /// to bring a node into the current context.
    ///
    /// [`dbus_interface`]: attr.dbus_interface.html
    pub fn local_node_emit_signal<B>(
        destination: Option<&str>,
        iface: &str,
        signal_name: &str,
        body: &B,
    ) -> Result<()>
    where
        B: serde::ser::Serialize + zvariant::Type,
    {
        if !LOCAL_NODE.is_set() {
            panic!("emit_signal: Node TLS not set");
        }

        LOCAL_NODE.with(|n| n.emit_signal(destination, iface, signal_name, body))
    }

    fn dispatch_method_call_try(
        &mut self,
        msg_header: &MessageHeader<'_>,
        msg: &Message,
    ) -> fdo::Result<Result<u32>> {
        let conn = self.conn.clone();
        let path = msg_header
            .path()
            .ok()
            .flatten()
            .ok_or_else(|| fdo::Error::Failed("Missing object path".into()))?;
        let iface = msg_header
            .interface()
            .ok()
            .flatten()
            // TODO: In the absence of an INTERFACE field, if two or more interfaces on the same object
            // have a method with the same name, it is undefined which of those methods will be
            // invoked. Implementations may choose to either return an error, or deliver the message
            // as though it had an arbitrary one of those interfaces.
            .ok_or_else(|| fdo::Error::Failed("Missing interface".into()))?;
        let member = msg_header
            .member()
            .ok()
            .flatten()
            .ok_or_else(|| fdo::Error::Failed("Missing member".into()))?;

        let node = self
            .get_node_mut(&path, false)
            .ok_or_else(|| fdo::Error::UnknownObject(format!("Unknown object '{}'", path)))?;
        let iface = node.get_interface(iface).ok_or_else(|| {
            fdo::Error::UnknownInterface(format!("Unknown interface '{}'", iface))
        })?;

        LOCAL_CONNECTION.set(&conn, || {
            LOCAL_NODE.set(node, || {
                let res = iface.borrow().call(&conn, &msg, member);
                res.or_else(|| iface.borrow_mut().call_mut(&conn, &msg, member))
                    .ok_or_else(|| {
                        fdo::Error::UnknownMethod(format!("Unknown method '{}'", member))
                    })
            })
        })
    }

    fn dispatch_method_call(
        &mut self,
        msg_header: &MessageHeader<'_>,
        msg: &Message,
    ) -> Result<u32> {
        match self.dispatch_method_call_try(msg_header, msg) {
            Err(e) => e.reply(&self.conn, msg),
            Ok(r) => r,
        }
    }

    /// Dispatch an incoming message to a registered interface.
    ///
    /// The object server will handle the message by:
    ///
    /// - looking up the called object path & interface,
    ///
    /// - calling the associated method if one exists,
    ///
    /// - returning a message (responding to the caller with either a return or error message) to
    ///   the caller through the associated server connection.
    ///
    /// Returns an error if the message is malformed, true if it's handled, false otherwise.
    ///
    /// # Note
    ///
    /// This API is subject to change, or becoming internal-only once zbus provides a general
    /// mechanism to dispatch messages.
    pub fn dispatch_message(&mut self, msg: &Message) -> Result<bool> {
        let msg_header = msg.header()?;

        match msg_header.message_type()? {
            MessageType::MethodCall => {
                self.dispatch_method_call(&msg_header, &msg)?;
                Ok(true)
            }
            _ => Ok(false),
        }
    }

    /// Receive and handle the next message from the associated connection.
    ///
    /// This function will read the incoming message from
    /// [`receive_message()`](Connection::receive_message) of the associated connection and pass it
    /// to [`dispatch_message()`](Self::dispatch_message). If the message was handled by an an
    /// interface, it returns `Ok(None)`. If not, it returns the received message.
    ///
    /// Returns an error if the message is malformed or an error occured.
    ///
    /// # Note
    ///
    /// This API is subject to change, or becoming internal-only once zbus provides a general
    /// mechanism to dispatch messages.
    pub fn try_handle_next(&mut self) -> Result<Option<Message>> {
        let msg = self.conn.receive_message()?;

        if !self.dispatch_message(&msg)? {
            return Ok(Some(msg));
        }

        Ok(None)
    }
}

#[cfg(test)]
mod tests {
    use std::{cell::Cell, collections::HashMap, convert::TryInto, error::Error, rc::Rc, thread};

    use ntest::timeout;
    use serde::{Deserialize, Serialize};
    use zvariant::derive::Type;

    use crate::{
        dbus_interface, dbus_proxy, fdo, Connection, MessageHeader, MessageType, ObjectServer,
    };

    #[derive(Deserialize, Serialize, Type)]
    pub struct ArgStructTest {
        foo: i32,
        bar: String,
    }

    #[dbus_proxy]
    trait MyIface {
        fn ping(&self) -> zbus::Result<u32>;

        fn quit(&self) -> zbus::Result<()>;

        fn test_header(&self) -> zbus::Result<()>;

        fn test_error(&self) -> zbus::Result<()>;

        fn test_single_struct_arg(&self, arg: ArgStructTest) -> zbus::Result<()>;

        fn test_hashmap_return(&self) -> zbus::Result<HashMap<String, String>>;

        fn create_obj(&self, key: &str) -> zbus::Result<()>;

        fn destroy_obj(&self, key: &str) -> zbus::Result<()>;

        #[dbus_proxy(property)]
        fn count(&self) -> zbus::Result<u32>;

        #[dbus_proxy(property)]
        fn set_count(&self, count: u32) -> zbus::Result<()>;

        #[dbus_proxy(property)]
        fn hash_map(&self) -> zbus::Result<HashMap<String, String>>;
    }

    #[derive(Debug, Clone)]
    enum NextAction {
        Nothing,
        Quit,
        CreateObj(String),
        DestroyObj(String),
    }

    struct MyIfaceImpl {
        action: Rc<Cell<NextAction>>,
        count: u32,
    }

    impl MyIfaceImpl {
        fn new(action: Rc<Cell<NextAction>>) -> Self {
            Self { action, count: 0 }
        }
    }

    #[dbus_interface(interface = "org.freedesktop.MyIface")]
    impl MyIfaceImpl {
        fn ping(&mut self) -> u32 {
            self.count += 1;
            if self.count % 3 == 0 {
                self.alert_count(self.count).expect("Failed to emit signal");
            }
            self.count
        }

        fn quit(&mut self) {
            self.action.set(NextAction::Quit);
        }

        fn test_header(&self, #[zbus(header)] header: MessageHeader<'_>) {
            assert_eq!(header.message_type().unwrap(), MessageType::MethodCall);
            assert_eq!(header.member().unwrap(), Some("TestHeader"));
        }

        fn test_error(&self) -> zbus::fdo::Result<()> {
            Err(zbus::fdo::Error::Failed("error raised".to_string()))
        }

        fn test_single_struct_arg(&self, arg: ArgStructTest) {
            assert_eq!(arg.foo, 1);
            assert_eq!(arg.bar, "TestString");
        }

        fn test_hashmap_return(&self) -> zbus::Result<HashMap<String, String>> {
            let mut map = HashMap::new();
            map.insert("hi".into(), "hello".into());
            map.insert("bye".into(), "now".into());

            Ok(map)
        }

        fn create_obj(&self, key: String) {
            self.action.set(NextAction::CreateObj(key));
        }

        fn destroy_obj(&self, key: String) {
            self.action.set(NextAction::DestroyObj(key));
        }

        #[dbus_interface(property)]
        fn set_count(&mut self, val: u32) -> zbus::fdo::Result<()> {
            if val == 42 {
                return Err(zbus::fdo::Error::InvalidArgs("Tsss tsss!".to_string()));
            }
            self.count = val;
            Ok(())
        }

        #[dbus_interface(property)]
        fn count(&self) -> u32 {
            self.count
        }

        #[dbus_interface(property)]
        fn hash_map(&self) -> HashMap<String, String> {
            self.test_hashmap_return().unwrap()
        }

        #[dbus_interface(signal)]
        fn alert_count(&self, val: u32) -> zbus::Result<()>;
    }

    fn check_hash_map(map: HashMap<String, String>) {
        assert_eq!(map["hi"], "hello");
        assert_eq!(map["bye"], "now");
    }

    fn my_iface_test() -> std::result::Result<u32, Box<dyn Error>> {
        let conn = Connection::new_session()?;
        let proxy = MyIfaceProxy::new_for(
            &conn,
            "org.freedesktop.MyService",
            "/org/freedesktop/MyService",
        )?;

        proxy.ping()?;
        assert_eq!(proxy.count()?, 1);
        proxy.test_header()?;
        proxy.test_single_struct_arg(ArgStructTest {
            foo: 1,
            bar: "TestString".into(),
        })?;
        check_hash_map(proxy.test_hashmap_return()?);
        check_hash_map(proxy.hash_map()?);
        proxy.introspect()?;
        let val = proxy.ping()?;

        proxy.create_obj("MyObj")?;
        let my_obj_proxy =
            MyIfaceProxy::new_for(&conn, "org.freedesktop.MyService", "/zbus/test/MyObj")?;
        my_obj_proxy.ping()?;
        proxy.destroy_obj("MyObj")?;
        assert!(my_obj_proxy.introspect().is_err());
        assert!(my_obj_proxy.ping().is_err());

        proxy.quit()?;
        Ok(val)
    }

    #[test]
    #[timeout(2000)]
    fn basic_iface() {
        let conn = Connection::new_session().unwrap();
        let mut object_server = ObjectServer::new(&conn);
        let action = Rc::new(Cell::new(NextAction::Nothing));

        fdo::DBusProxy::new(&conn)
            .unwrap()
            .request_name(
                "org.freedesktop.MyService",
                fdo::RequestNameFlags::ReplaceExisting.into(),
            )
            .unwrap();

        let iface = MyIfaceImpl::new(action.clone());
        object_server
            .at(&"/org/freedesktop/MyService".try_into().unwrap(), iface)
            .unwrap();

        let child = thread::spawn(|| my_iface_test().expect("child failed"));

        loop {
            let m = conn.receive_message().unwrap();
            if let Err(e) = object_server.dispatch_message(&m) {
                eprintln!("{}", e);
            }

            object_server
                .with(
                    &"/org/freedesktop/MyService".try_into().unwrap(),
                    |iface: &MyIfaceImpl| iface.alert_count(51),
                )
                .unwrap();

            match action.replace(NextAction::Nothing) {
                NextAction::Nothing => (),
                NextAction::Quit => break,
                NextAction::CreateObj(key) => {
                    let path = format!("/zbus/test/{}", key);
                    object_server
                        .at(&path.try_into().unwrap(), MyIfaceImpl::new(action.clone()))
                        .unwrap();
                }
                NextAction::DestroyObj(key) => {
                    let path = format!("/zbus/test/{}", key);
                    object_server
                        .remove::<MyIfaceImpl>(&path.try_into().unwrap())
                        .unwrap();
                }
            }
        }

        let val = child.join().expect("failed to join");
        assert_eq!(val, 2);
    }
}
