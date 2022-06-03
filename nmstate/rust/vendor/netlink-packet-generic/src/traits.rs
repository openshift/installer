// SPDX-License-Identifier: MIT

//! Traits for implementing generic netlink family

/// Provide the definition for generic netlink family
///
/// Family payload type should implement this trait to provide necessary
/// informations in order to build the packet headers (`nlmsghdr` and `genlmsghdr`).
///
/// If you are looking for an example implementation, you can refer to the
/// [`crate::ctrl`] module.
pub trait GenlFamily {
    /// Return the unique family name registered in the kernel
    ///
    /// Let the resolver lookup the dynamically assigned ID
    fn family_name() -> &'static str;

    /// Return the assigned family ID
    ///
    /// # Note
    /// The implementation of generic family should assign the ID to `GENL_ID_GENERATE` (0x0).
    /// So the controller can dynamically assign the family ID.
    ///
    /// Regarding to the reason above, you should not have to implement the function
    /// unless the family uses static ID.
    fn family_id(&self) -> u16 {
        0
    }

    /// Return the command type of the current message
    fn command(&self) -> u8;

    /// Indicate the protocol version
    fn version(&self) -> u8;
}
