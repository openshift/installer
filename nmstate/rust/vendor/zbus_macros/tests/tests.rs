use zbus::{self, fdo};
use zbus_macros::{dbus_interface, dbus_proxy, DBusError};

#[test]
fn test_proxy() {
    #[dbus_proxy(
        interface = "org.freedesktop.zbus.Test",
        default_service = "org.freedesktop.zbus",
        default_path = "/org/freedesktop/zbus/test"
    )]
    trait Test {
        /// comment for a_test()
        fn a_test(&self, val: &str) -> zbus::Result<u32>;

        #[dbus_proxy(name = "CheckRENAMING")]
        fn check_renaming(&self) -> zbus::Result<Vec<u8>>;

        #[dbus_proxy(property)]
        fn property(&self) -> fdo::Result<Vec<String>>;

        #[dbus_proxy(property)]
        fn set_property(&self, val: u16) -> fdo::Result<()>;
    }
}

#[test]
fn test_derive_error() {
    #[derive(Debug, DBusError)]
    #[dbus_error(prefix = "org.freedesktop.zbus")]
    enum Test {
        ZBus(zbus::Error),
        SomeExcuse,
        #[dbus_error(name = "I.Am.Sorry.Dave")]
        IAmSorryDave(String),
        LetItBe {
            desc: String,
        },
    }
}

#[test]
fn test_interface() {
    use zbus::Interface;

    struct Test<'a, T> {
        something: &'a str,
        generic: T,
    }

    #[dbus_interface(name = "org.freedesktop.zbus.Test")]
    impl<T: 'static> Test<'static, T>
    where
        T: serde::ser::Serialize + zvariant::Type,
    {
        /// Testing `no_arg` documentation is reflected in XML.
        fn no_arg(&self) {
            unimplemented!()
        }

        fn str_u32(&self, val: &str) -> zbus::fdo::Result<u32> {
            val.parse()
                .map_err(|e| zbus::fdo::Error::Failed(format!("Invalid val: {}", e)))
        }

        // TODO: naming output arguments after "RFC: Structural Records #2584"
        fn many_output(&self) -> zbus::fdo::Result<(&T, String)> {
            Ok((&self.generic, self.something.to_string()))
        }

        fn pair_output(&self) -> zbus::fdo::Result<((u32, String),)> {
            unimplemented!()
        }

        #[dbus_interface(name = "CheckVEC")]
        fn check_vec(&self) -> Vec<u8> {
            unimplemented!()
        }

        /// Testing my_prop documentation is reflected in XML.
        ///
        /// And that too.
        #[dbus_interface(property)]
        fn my_prop(&self) -> u16 {
            unimplemented!()
        }

        #[dbus_interface(property)]
        fn set_my_prop(&self, _val: u16) {
            unimplemented!()
        }

        /// Emit a signal.
        #[dbus_interface(signal)]
        fn signal(&self, arg: u8, other: &str) -> zbus::Result<()>;
    }

    const EXPECTED_XML: &str = r#"<interface name="org.freedesktop.zbus.Test">
  <!--
   Testing `no_arg` documentation is reflected in XML.
   -->
  <method name="NoArg">
  </method>
  <method name="StrU32">
    <arg name="val" type="s" direction="in"/>
    <arg type="u" direction="out"/>
  </method>
  <method name="ManyOutput">
    <arg type="u" direction="out"/>
    <arg type="s" direction="out"/>
  </method>
  <method name="PairOutput">
    <arg type="(us)" direction="out"/>
  </method>
  <method name="CheckVEC">
    <arg type="ay" direction="out"/>
  </method>
  <!--
   Emit a signal.
   -->
  <signal name="Signal">
    <arg name="arg" type="y"/>
    <arg name="other" type="s"/>
  </signal>
  <!--
   Testing my_prop documentation is reflected in XML.

   And that too.
   -->
  <property name="MyProp" type="q" access="readwrite"/>
</interface>
"#;
    let t = Test {
        something: &"somewhere",
        generic: 42u32,
    };
    let mut xml = String::new();
    t.introspect_to_writer(&mut xml, 0);
    assert_eq!(xml, EXPECTED_XML);

    assert_eq!(Test::<u32>::name(), "org.freedesktop.zbus.Test");

    if false {
        // check compilation
        let c = zbus::Connection::new_session().unwrap();
        let m = zbus::Message::method(None, None, "/", None, "StrU32", &(42,)).unwrap();
        let _ = t.call(&c, &m, "StrU32").unwrap();
        t.signal(23, "ergo sum").unwrap();
    }
}
