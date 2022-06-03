use anyhow::Context;
use std::net::IpAddr;

use crate::{
    constants::{self, RTA_GATEWAY},
    emit_ip,
    ip_len,
    nlas::{Nla, NlaBuffer},
    parsers::parse_ip,
    traits::{Emitable, Parseable},
    DecodeError,
};

bitflags! {
    pub struct NextHopFlags: u8 {
        const RTNH_F_EMPTY = 0;
        const RTNH_F_DEAD = constants::RTNH_F_DEAD as u8;
        const RTNH_F_PERVASIVE = constants::RTNH_F_PERVASIVE as u8;
        const RTNH_F_ONLINK = constants::RTNH_F_ONLINK as u8;
        const RTNH_F_OFFLOAD = constants::RTNH_F_OFFLOAD as u8;
        const RTNH_F_LINKDOWN = constants::RTNH_F_LINKDOWN as u8;
        const RTNH_F_UNRESOLVED = constants::RTNH_F_UNRESOLVED as u8;
    }
}

buffer!(NextHopBuffer {
    length: (u16, 0..2),
    flags: (u8, 2),
    hops: (u8, 3),
    interface_id: (u32, 4..8),
    gateway_nla: (slice, GATEWAY_OFFSET..),
});

impl<T: AsRef<[u8]>> NextHopBuffer<T> {
    pub fn new_checked(buffer: T) -> Result<Self, DecodeError> {
        let packet = Self::new(buffer);
        packet.check_buffer_length()?;
        Ok(packet)
    }

    fn check_buffer_length(&self) -> Result<(), DecodeError> {
        let len = self.buffer.as_ref().len();
        if len < 8 {
            return Err(format!("invalid NextHopBuffer: length {} < {}", len, 8).into());
        }
        if len < self.length() as usize {
            return Err(format!(
                "invalid NextHopBuffer: length {} < {}",
                len,
                8 + self.length()
            )
            .into());
        }
        Ok(())
    }
}

const GATEWAY_OFFSET: usize = 8;

#[derive(Debug, Clone, Copy, Eq, PartialEq)]
pub struct NextHop {
    /// Next-hop flags (see [`NextHopFlags`])
    pub flags: NextHopFlags,
    /// Next-hop priority
    pub hops: u8,
    /// Interface index for the next-hop
    pub interface_id: u32,
    /// Gateway address (it is actually encoded as an `RTA_GATEWAY` nla)
    pub gateway: Option<IpAddr>,
}

impl<'a, T: AsRef<[u8]>> Parseable<NextHopBuffer<&'a T>> for NextHop {
    fn parse(buf: &NextHopBuffer<&T>) -> Result<NextHop, DecodeError> {
        let gateway = if buf.length() as usize > GATEWAY_OFFSET {
            let gateway_nla_buf = NlaBuffer::new_checked(buf.gateway_nla())
                .context("cannot parse RTA_GATEWAY attribute in next-hop")?;
            if gateway_nla_buf.kind() != RTA_GATEWAY {
                return Err(format!("invalid RTA_GATEWAY attribute in next-hop: expected NLA type to be RTA_GATEWAY ({}), but got {} instead", RTA_GATEWAY, gateway_nla_buf.kind()).into());
            }
            let gateway = parse_ip(gateway_nla_buf.value()).context(
                "invalid RTA_GATEWAY attribute in next-hop: failed to parse NLA value as an IP address",
            )?;
            Some(gateway)
        } else {
            None
        };
        Ok(NextHop {
            flags: NextHopFlags::from_bits_truncate(buf.flags()),
            hops: buf.hops(),
            interface_id: buf.interface_id(),
            gateway,
        })
    }
}

struct GatewayNla<'a>(&'a IpAddr);

impl<'a> Nla for GatewayNla<'a> {
    fn value_len(&self) -> usize {
        ip_len(self.0)
    }
    fn kind(&self) -> u16 {
        RTA_GATEWAY
    }
    fn emit_value(&self, buffer: &mut [u8]) {
        emit_ip(buffer, self.0)
    }
}

impl Emitable for NextHop {
    fn buffer_len(&self) -> usize {
        // len, flags, hops and interface id fields
        8 + self
            .gateway
            .as_ref()
            .map(|ip| {
                // RTA_GATEWAY attribute header (length and type) + value length
                4 + ip_len(ip)
            })
            .unwrap_or(0)
    }

    fn emit(&self, buffer: &mut [u8]) {
        let mut nh_buffer = NextHopBuffer::new(buffer);
        nh_buffer.set_length(self.buffer_len() as u16);
        nh_buffer.set_flags(self.flags.bits());
        nh_buffer.set_hops(self.hops);
        nh_buffer.set_interface_id(self.interface_id);
        if let Some(ref gateway) = self.gateway {
            let gateway_nla = GatewayNla(gateway);
            gateway_nla.emit(nh_buffer.gateway_nla_mut());
        }
    }
}
