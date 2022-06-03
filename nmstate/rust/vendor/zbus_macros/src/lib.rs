#![deny(rust_2018_idioms)]
#![doc(
    html_logo_url = "https://storage.googleapis.com/fdo-gitlab-uploads/project/avatar/3213/zbus-logomark.png"
)]

//! This crate provides derive macros helpers for zbus.
use proc_macro::TokenStream;
use syn::{parse_macro_input, AttributeArgs, DeriveInput, ItemImpl, ItemTrait};

mod error;
mod iface;
mod proxy;
mod utils;

/// Attribute macro for defining a D-Bus proxy (using zbus [`Proxy`]).
///
/// The macro must be applied on a `trait T`. A matching `impl T` will provide the proxy. The proxy
/// instance can be created with the associated `new()` or `new_for()` methods. The former doesn't take
/// any argument and uses the default service name and path. The later allows you to specify both.
///
/// Each trait method will be expanded to call to the associated D-Bus remote interface.
///
/// Trait methods accept `dbus_proxy` attributes:
///
/// * `name` - override the D-Bus name (pascal case form by default)
///
/// * `property` - expose the method as a property. If the method takes an argument, it must be a
///   setter, with a `set_` prefix. Otherwise, it's a getter.
///
/// * `signal` - declare a signal just like a D-Bus method. The macro will provide a method to
///   register and deregister a handler for the signal, whose signature must match that of the
///   signature declaration.
///
///   NB: Any doc comments provided shall be appended to the ones added by the macro.
///
/// (the expanded `impl` also provides an `introspect()` method, for convenience)
///
/// # Example
///
/// ```
///# use std::error::Error;
/// use zbus_macros::dbus_proxy;
/// use zbus::{Connection, Result, fdo};
/// use zvariant::Value;
///
/// #[dbus_proxy(
///     interface = "org.test.SomeIface",
///     default_service = "org.test.SomeService",
///     default_path = "/org/test/SomeObject"
/// )]
/// trait SomeIface {
///     fn do_this(&self, with: &str, some: u32, arg: &Value) -> Result<bool>;
///
///     #[dbus_proxy(property)]
///     fn a_property(&self) -> fdo::Result<String>;
///
///     #[dbus_proxy(property)]
///     fn set_a_property(&self, a_property: &str) -> fdo::Result<()>;
///
///     #[dbus_proxy(signal)]
///     fn some_signal(&self, arg1: &str, arg2: u32) -> fdo::Result<()>;
/// };
///
/// let connection = Connection::new_session()?;
/// let proxy = SomeIfaceProxy::new(&connection)?;
/// let _ = proxy.do_this("foo", 32, &Value::new(true));
/// let _ = proxy.set_a_property("val");
///
/// proxy.connect_some_signal(|s, u| {
///     println!("arg1: {}, arg2: {}", s, u);
///
///     Ok(())
/// })?;
///
/// // You'll want to make at least a call to `handle_next_signal` before disconnecting the signal.
/// assert!(proxy.disconnect_some_signal()?);
/// assert!(!proxy.disconnect_some_signal()?);
///
///# Ok::<_, Box<dyn Error + Send + Sync>>(())
/// ```
///
/// [`zbus_polkit`] is a good example of how to bind a real D-Bus API.
///
/// [`zbus_polkit`]: https://docs.rs/zbus_polkit/1.0.0/zbus_polkit/policykit1/index.html
/// [`Proxy`]: https://docs.rs/zbus/1.0.0/zbus/struct.Proxy.html
/// [`zbus::SignalReceiver::receive_for`]:
/// https://docs.rs/zbus/1.5.0/zbus/struct.SignalReceiver.html#method.receive_for
#[proc_macro_attribute]
pub fn dbus_proxy(attr: TokenStream, item: TokenStream) -> TokenStream {
    let args = parse_macro_input!(attr as AttributeArgs);
    let input = parse_macro_input!(item as ItemTrait);
    proxy::expand(args, input).into()
}

/// Attribute macro for implementing a D-Bus interface.
///
/// The macro must be applied on an `impl T`. All methods will be exported, either as methods,
/// properties or signal depending on the item attributes. It will implement the [`Interface`] trait
/// `for T` on your behalf, to handle the message dispatching and introspection support.
///
/// The methods accepts the `dbus_interface` attributes:
///
/// * `name` - override the D-Bus name (pascal case form of the method by default)
///
/// * `property` - expose the method as a property. If the method takes an argument, it must be a
///   setter, with a `set_` prefix. Otherwise, it's a getter.
///
/// * `signal` - the method is a "signal". It must be a method declaration (without body). Its code
///   block will be expanded to emit the signal from the object path associated with the interface
///   instance.
///
///   You can call a signal method from a an interface method, or from an [`ObjectServer::with`]
///   function.
///
/// * `struct_return` - the method returns a structure. Although it is very rare for a D-Bus method
///   to return a single structure, it does happen. Since it is not possible for zbus to
///   differentiate this case from multiple out arguments put in a structure, you must either
///   explicitly mark the case of single structure return with this attribute or wrap the structure
///   in another structure or tuple. The latter can be achieve by using `(,)` syntax, for example
///   instead of `MyStruct`, write `(MyStruct,)`.
///
/// The method arguments offers some the following `zbus` attributes:
///
/// * `header` - This marks the method argument to receive the message header associated with the
/// D-Bus method call being handled.
///
/// # Example
///
/// ```
///# use std::error::Error;
/// use zbus_macros::dbus_interface;
/// use zbus::MessageHeader;
///
/// struct Example {
///     some_data: String,
/// }
///
/// #[dbus_interface(name = "org.myservice.Example")]
/// impl Example {
///     // "Quit" method. A method may throw errors.
///     fn quit(&self, #[zbus(header)] hdr: MessageHeader<'_>) -> zbus::fdo::Result<()> {
///         let path = hdr.path()?.unwrap();
///         let msg = format!("You are leaving me on the {} path?", path);
///
///         Err(zbus::fdo::Error::Failed(msg))
///     }
///
///     // "TheAnswer" property (note: the "name" attribute), with its associated getter.
///     #[dbus_interface(property, name = "TheAnswer")]
///     fn answer(&self) -> u32 {
///         2 * 3 * 7
///     }
///
///     // "Notify" signal (note: no implementation body).
///     #[dbus_interface(signal)]
///     fn notify(&self, message: &str) -> zbus::Result<()>;
/// }
///
///# Ok::<_, Box<dyn Error + Send + Sync>>(())
/// ```
///
/// See also [`ObjectServer`] documentation to learn how to export an interface over a `Connection`.
///
/// [`ObjectServer`]: https://docs.rs/zbus/1.0.0/zbus/struct.ObjectServer.html
/// [`ObjectServer::with`]: https://docs.rs/zbus/1.2.0/zbus/struct.ObjectServer.html#method.with
/// [`Connection::emit_signal()`]: https://docs.rs/zbus/1.0.0/zbus/struct.Connection.html#method.emit_signal
/// [`Interface`]: https://docs.rs/zbus/1.0.0/zbus/trait.Interface.html
#[proc_macro_attribute]
pub fn dbus_interface(attr: TokenStream, item: TokenStream) -> TokenStream {
    let args = parse_macro_input!(attr as AttributeArgs);
    let input = syn::parse_macro_input!(item as ItemImpl);
    iface::expand(args, input)
        .unwrap_or_else(|err| err.to_compile_error())
        .into()
}

/// Derive macro for defining a D-Bus error.
///
/// This macro helps to implement an [`Error`] suitable for D-Bus handling with zbus. It will expand
/// an `enum E` with [`Error`] traits implementation, and `From<zbus::Error>`. The latter makes it
/// possible for you to declare proxy methods to directly return this type, rather than
/// [`zbus::Error`]. However, for this to work, we require a variant by the name `ZBus` that
/// contains an unnamed field of type [`zbus::Error`].
///
/// Additionnally, the derived `impl E` will provide the following convenience methods:
///
/// * `name(&self)` - get the associated D-Bus error name.
///
/// * `description(&self)` - get the associated error description (the first argument of an error
///   message)
///
/// * `reply(&self, &zbus::Connection, &zbus::Message)` - send this error as reply to the message.
///
/// Note: it is recommended that errors take a single argument `String` which describes it in
/// a human-friendly fashion (support for other arguments is limited or TODO currently).
///
/// # Example
///
/// ```
/// use zbus_macros::DBusError;
///
/// #[derive(DBusError, Debug)]
/// #[dbus_error(prefix = "org.myservice.App")]
/// enum Error {
///     ZBus(zbus::Error),
///     FileNotFound(String),
///     OutOfMemory,
/// }
/// ```
///
/// [`Error`]: http://doc.rust-lang.org/std/error/trait.Error.html
/// [`zbus::Error`]: https://docs.rs/zbus/1.0.0/zbus/enum.Error.html
#[proc_macro_derive(DBusError, attributes(dbus_error))]
pub fn derive_dbus_error(input: TokenStream) -> TokenStream {
    let input = parse_macro_input!(input as DeriveInput);
    error::expand_derive(input).into()
}
