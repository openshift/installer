// SPDX-License-Identifier: MIT

//! This crate provides the packet of generic netlink family and its controller.
//!
//! The `[GenlMessage]` provides a generic netlink family message which is
//! sub-protocol independant.
//! You can wrap your message into the type, then it can be used in `netlink-proto` crate.
//!
//! # Implementing a generic netlink family
//! A generic netlink family contains several commands, and a version number in
//! the header.
//!
//! The payload usually consists of netlink attributes, carrying the messages to
//! the peer. In order to help you to make your payload into a valid netlink
//! packet, this crate requires the informations about the family id,
//! and the informations in the generic header. So, you need to implement some
//! traits on your types.
//!
//! All the things in the payload including all netlink attributes used
//! and the optional header should be handled by your implementation.
//!
//! ## Serializaion / Deserialization
//! To implement your generic netlink family, you should handle the payload
//! serialization process including its specific header (if any) and the netlink
//! attributes.
//!
//! To achieve this, you should implement [`netlink_packet_utils::Emitable`]
//! trait for the payload type.
//!
//! For deserialization, [`netlink_packet_utils::ParseableParametrized<[u8], GenlHeader>`](netlink_packet_utils::ParseableParametrized)
//! trait should be implemented. As mention above, to provide more scalability,
//! we use the simplest buffer type: `[u8]` here. You can turn it into other
//! buffer type easily during deserializing.
//!
//! ## `GenlFamily` trait
//! The trait is aim to provide some necessary informations in order to build
//! the packet headers of netlink (nlmsghdr) and generic netlink (genlmsghdr).
//!
//! ### `family_name()`
//! The method let the resolver to obtain the name registered in the kernel.
//!
//! ### `family_id()`
//! Few netlink family has static family ID (e.g. controller). The method is
//! mainly used to let those family to return their familt ID.
//!
//! If you don't know what is this, please **DO NOT** implement this method.
//! Since the default implementation return `GENL_ID_GENERATE`, which means
//! the family ID is allocated by the kernel dynamically.
//!
//! ### `command()`
//! This method tells the generic netlink command id of the packet
//! The return value is used to fill the `cmd` field in the generic netlink header.
//!
//! ### `version()`
//! This method return the family version of the payload.
//! The return value is used to fill the `version` field in the generic netlink header.
//!
//! ## Family Header
//! Few family would use a family specific message header. For simplification
//! and scalability, this crate treats it as a part of the payload, and make
//! implementations to handle the header by themselves.
//!
//! If you are implementing such a generic family, note that you should define
//! the header data structure in your payload type and handle the serialization.

#[macro_use]
extern crate netlink_packet_utils;

pub mod buffer;
pub use self::buffer::GenlBuffer;

pub mod constants;

pub mod ctrl;

pub mod header;
pub use self::header::GenlHeader;

pub mod message;
pub use self::message::GenlMessage;

pub mod traits;
pub use self::traits::GenlFamily;
