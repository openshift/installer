// SPDX-License-Identifier: MIT

//! Define constants related to generic netlink
pub const GENL_ID_CTRL: u16 = libc::GENL_ID_CTRL as u16;
pub const GENL_HDRLEN: usize = 4;

pub const CTRL_CMD_UNSPEC: u8 = libc::CTRL_CMD_UNSPEC as u8;
pub const CTRL_CMD_NEWFAMILY: u8 = libc::CTRL_CMD_NEWFAMILY as u8;
pub const CTRL_CMD_DELFAMILY: u8 = libc::CTRL_CMD_DELFAMILY as u8;
pub const CTRL_CMD_GETFAMILY: u8 = libc::CTRL_CMD_GETFAMILY as u8;
pub const CTRL_CMD_NEWOPS: u8 = libc::CTRL_CMD_NEWOPS as u8;
pub const CTRL_CMD_DELOPS: u8 = libc::CTRL_CMD_DELOPS as u8;
pub const CTRL_CMD_GETOPS: u8 = libc::CTRL_CMD_GETOPS as u8;
pub const CTRL_CMD_NEWMCAST_GRP: u8 = libc::CTRL_CMD_NEWMCAST_GRP as u8;
pub const CTRL_CMD_DELMCAST_GRP: u8 = libc::CTRL_CMD_DELMCAST_GRP as u8;
pub const CTRL_CMD_GETMCAST_GRP: u8 = libc::CTRL_CMD_GETMCAST_GRP as u8;
pub const CTRL_CMD_GETPOLICY: u8 = 10;

pub const CTRL_ATTR_UNSPEC: u16 = libc::CTRL_ATTR_UNSPEC as u16;
pub const CTRL_ATTR_FAMILY_ID: u16 = libc::CTRL_ATTR_FAMILY_ID as u16;
pub const CTRL_ATTR_FAMILY_NAME: u16 = libc::CTRL_ATTR_FAMILY_NAME as u16;
pub const CTRL_ATTR_VERSION: u16 = libc::CTRL_ATTR_VERSION as u16;
pub const CTRL_ATTR_HDRSIZE: u16 = libc::CTRL_ATTR_HDRSIZE as u16;
pub const CTRL_ATTR_MAXATTR: u16 = libc::CTRL_ATTR_MAXATTR as u16;
pub const CTRL_ATTR_OPS: u16 = libc::CTRL_ATTR_OPS as u16;
pub const CTRL_ATTR_MCAST_GROUPS: u16 = libc::CTRL_ATTR_MCAST_GROUPS as u16;
pub const CTRL_ATTR_POLICY: u16 = 8;
pub const CTRL_ATTR_OP_POLICY: u16 = 9;
pub const CTRL_ATTR_OP: u16 = 10;

pub const CTRL_ATTR_OP_UNSPEC: u16 = libc::CTRL_ATTR_OP_UNSPEC as u16;
pub const CTRL_ATTR_OP_ID: u16 = libc::CTRL_ATTR_OP_ID as u16;
pub const CTRL_ATTR_OP_FLAGS: u16 = libc::CTRL_ATTR_OP_FLAGS as u16;

pub const CTRL_ATTR_MCAST_GRP_UNSPEC: u16 = libc::CTRL_ATTR_MCAST_GRP_UNSPEC as u16;
pub const CTRL_ATTR_MCAST_GRP_NAME: u16 = libc::CTRL_ATTR_MCAST_GRP_NAME as u16;
pub const CTRL_ATTR_MCAST_GRP_ID: u16 = libc::CTRL_ATTR_MCAST_GRP_ID as u16;

pub const CTRL_ATTR_POLICY_UNSPEC: u16 = 0;
pub const CTRL_ATTR_POLICY_DO: u16 = 1;
pub const CTRL_ATTR_POLICY_DUMP: u16 = 2;

pub const NL_ATTR_TYPE_INVALID: u32 = 0;
pub const NL_ATTR_TYPE_FLAG: u32 = 1;
pub const NL_ATTR_TYPE_U8: u32 = 2;
pub const NL_ATTR_TYPE_U16: u32 = 3;
pub const NL_ATTR_TYPE_U32: u32 = 4;
pub const NL_ATTR_TYPE_U64: u32 = 5;
pub const NL_ATTR_TYPE_S8: u32 = 6;
pub const NL_ATTR_TYPE_S16: u32 = 7;
pub const NL_ATTR_TYPE_S32: u32 = 8;
pub const NL_ATTR_TYPE_S64: u32 = 9;
pub const NL_ATTR_TYPE_BINARY: u32 = 10;
pub const NL_ATTR_TYPE_STRING: u32 = 11;
pub const NL_ATTR_TYPE_NUL_STRING: u32 = 12;
pub const NL_ATTR_TYPE_NESTED: u32 = 13;
pub const NL_ATTR_TYPE_NESTED_ARRAY: u32 = 14;
pub const NL_ATTR_TYPE_BITFIELD32: u32 = 15;

pub const NL_POLICY_TYPE_ATTR_UNSPEC: u16 = 0;
pub const NL_POLICY_TYPE_ATTR_TYPE: u16 = 1;
pub const NL_POLICY_TYPE_ATTR_MIN_VALUE_S: u16 = 2;
pub const NL_POLICY_TYPE_ATTR_MAX_VALUE_S: u16 = 3;
pub const NL_POLICY_TYPE_ATTR_MIN_VALUE_U: u16 = 4;
pub const NL_POLICY_TYPE_ATTR_MAX_VALUE_U: u16 = 5;
pub const NL_POLICY_TYPE_ATTR_MIN_LENGTH: u16 = 6;
pub const NL_POLICY_TYPE_ATTR_MAX_LENGTH: u16 = 7;
pub const NL_POLICY_TYPE_ATTR_POLICY_IDX: u16 = 8;
pub const NL_POLICY_TYPE_ATTR_POLICY_MAXTYPE: u16 = 9;
pub const NL_POLICY_TYPE_ATTR_BITFIELD32_MASK: u16 = 10;
pub const NL_POLICY_TYPE_ATTR_PAD: u16 = 11;
pub const NL_POLICY_TYPE_ATTR_MASK: u16 = 12;
