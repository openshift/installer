use std::{collections::VecDeque, io};

use crate::{message::Message, message_header::MIN_MESSAGE_SIZE, raw::Socket, OwnedFd};

/// A low-level representation of a D-Bus connection
///
/// This wrapper is agnostic on the actual transport, using the `Socket` trait
/// to abstract it. It is compatible with sockets both in blocking or non-blocking
/// mode.
///
/// This wrapper abstracts away the serialization & buffering considerations of the
/// protocol, and allows interaction based on messages, rather than bytes.
#[derive(derivative::Derivative)]
#[derivative(Debug)]
pub struct Connection<S> {
    #[derivative(Debug = "ignore")]
    socket: S,
    raw_in_buffer: Vec<u8>,
    raw_in_fds: Vec<OwnedFd>,
    msg_in_buffer: Option<Message>,
    raw_out_buffer: VecDeque<u8>,
    msg_out_buffer: VecDeque<Message>,
}

impl<S: Socket> Connection<S> {
    pub(crate) fn wrap(socket: S) -> Connection<S> {
        Connection {
            socket,
            raw_in_buffer: vec![],
            raw_in_fds: vec![],
            msg_in_buffer: None,
            raw_out_buffer: VecDeque::new(),
            msg_out_buffer: VecDeque::new(),
        }
    }

    /// Attempt to flush the outgoing buffer
    ///
    /// This will try to write as many messages as possible from the
    /// outgoing buffer into the socket, until an error is encountered.
    ///
    /// This method will thus only block if the socket is in blocking mode.
    pub fn try_flush(&mut self) -> io::Result<()> {
        // first, empty the raw_out_buffer of any partially-sent message
        while !self.raw_out_buffer.is_empty() {
            let (front, _) = self.raw_out_buffer.as_slices();
            // VecDeque should never return an empty front buffer if the VecDeque
            // itself is not empty
            debug_assert!(!front.is_empty());
            let written = self.socket.sendmsg(front, &[])?;
            self.raw_out_buffer.drain(..written);
        }

        // now, try to drain the msg_out_buffer
        while let Some(msg) = self.msg_out_buffer.front() {
            let mut data = msg.as_bytes();
            let fds = msg.fds();
            let written = self.socket.sendmsg(data, &fds)?;
            // at least some part of the message has been sent, see if we can/need to send more
            // now the message must be removed from msg_out_buffer and any leftover bytes
            // must be stored into raw_out_buffer
            let msg = self.msg_out_buffer.pop_front().unwrap();
            data = &msg.as_bytes()[written..];
            while !data.is_empty() {
                match self.socket.sendmsg(data, &[]) {
                    Ok(n) => data = &data[n..],
                    Err(e) => {
                        // an error occured, we cannot send more, store the remaining into
                        // raw_out_buffer and forward the error
                        self.raw_out_buffer.extend(data);
                        return Err(e);
                    }
                }
            }
        }
        Ok(())
    }

    /// Enqueue a message to be sent out to the socket
    ///
    /// This method will *not* write anything to the socket, you need to call
    /// `try_flush()` afterwards so that your message is actually sent out.
    pub fn enqueue_message(&mut self, msg: Message) {
        self.msg_out_buffer.push_back(msg);
    }

    /// Attempt to read a message from the socket
    ///
    /// This methods will read from the socket until either a full D-Bus message is
    /// read or an error is encountered.
    ///
    /// If the socket is in non-blocking mode, it may read a partial message. In such case it
    /// will buffer it internally and try to complete it the next time you call `try_receive_message`.
    pub fn try_receive_message(&mut self) -> crate::Result<Message> {
        if self.msg_in_buffer.is_none() {
            // We don't have enough data to make a proper message header yet.
            // Some partial read may be in raw_in_buffer, so we try to complete it
            // until we have MIN_MESSAGE_SIZE bytes
            //
            // Given that MIN_MESSAGE_SIZE is 16, this codepath is actually extremely unlikely
            // to be taken more than once
            while self.raw_in_buffer.len() < MIN_MESSAGE_SIZE {
                let current_bytes = self.raw_in_buffer.len();
                let mut buf = vec![0; MIN_MESSAGE_SIZE - current_bytes];
                let (read, fds) = self.socket.recvmsg(&mut buf)?;
                self.raw_in_buffer.extend(&buf[..read]);
                self.raw_in_fds.extend(fds);
            }

            // We now have a full message header, so let us construct the Message
            self.msg_in_buffer = Some(Message::from_bytes(&self.raw_in_buffer)?);
            self.raw_in_buffer.clear();
        }

        // At this point, we must have a partial message in self.msg_in_buffer, and we
        // need to complete it
        {
            let msg = self.msg_in_buffer.as_mut().unwrap();
            loop {
                match msg.bytes_to_completion() {
                    Ok(0) => {
                        // the message is now complete, we can return
                        break;
                    }
                    Ok(needed) => {
                        // we need to read more data
                        let mut buf = vec![0; needed];
                        let (read, fds) = self.socket.recvmsg(&mut buf)?;
                        msg.add_bytes(&buf[..read])?;
                        self.raw_in_fds.extend(fds);
                    }
                    Err(e) => {
                        // the message is invalid, return the error
                        return Err(e.into());
                    }
                }
            }
        }

        // If we reach here, the message is complete, return it
        let mut msg = self.msg_in_buffer.take().unwrap();
        msg.set_owned_fds(std::mem::replace(&mut self.raw_in_fds, vec![]));
        Ok(msg)
    }

    /// Close the connection.
    ///
    /// After this call, all reading and writing operations will fail.
    pub fn close(&self) -> crate::Result<()> {
        self.socket().close().map_err(|e| e.into())
    }

    /// Access the underlying socket
    ///
    /// This method is intended to provide access to the socket in order to register it
    /// to you event loop, for async integration.
    ///
    /// You should not try to read or write from it directly, as it may
    /// corrupt the internal state of this wrapper.
    pub fn socket(&self) -> &S {
        &self.socket
    }
}

#[cfg(test)]
mod tests {
    use super::Connection;
    use crate::message::Message;
    use std::os::unix::net::UnixStream;

    #[test]
    fn raw_send_receive() {
        let (p0, p1) = UnixStream::pair().unwrap();

        let mut conn0 = Connection::wrap(p0);
        let mut conn1 = Connection::wrap(p1);

        let msg = Message::method(None, None, "/", Some("org.zbus.p2p"), "Test", &()).unwrap();

        conn0.enqueue_message(msg);
        conn0.try_flush().unwrap();

        let ret = conn1.try_receive_message().unwrap();

        assert_eq!(ret.to_string(), "Method call Test");
    }
}
