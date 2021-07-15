# Load Balancer Setup

The libvirt deployment does not deploy a load balancer for development purposes. 

This doc goes over an example configuration of HAProxy for doing local development.

### Installing HAProxy
These instructions are for Fedora 34+.

Install the RPM for `HAProxy`.
```sh 
sudo dnf install haproxy
```

Configure `haproxy.cfg`. A default configuration follows, replace with the appropriate IP addresses for your environment:

```sh 
#---------------------------------------------------------------------
# Global settings
#---------------------------------------------------------------------
global
    log         127.0.0.1 local2

    chroot      /var/lib/haproxy
    pidfile     /var/run/haproxy.pid
    maxconn     4000
    user        haproxy
    group       haproxy
    daemon

    # turn on stats unix socket
    stats socket /var/lib/haproxy/stats

    # utilize system-wide crypto-policies
    # ssl-default-bind-ciphers PROFILE=SYSTEM
    # ssl-default-server-ciphers PROFILE=SYSTEM

#---------------------------------------------------------------------
# common defaults that all the 'listen' and 'backend' sections will
# use if not designated in their block
#---------------------------------------------------------------------
defaults
    mode                    tcp
    log                     global
    option                  httplog
    option                  dontlognull
    option http-server-close
    option forwardfor       except 127.0.0.0/8
    option                  redispatch
    retries                 3
    timeout http-request    10s
    timeout queue           1m
    timeout connect         10s
    timeout client          1m
    timeout server          1m
    timeout http-keep-alive 10s
    timeout check           10s
    maxconn                 3000

#---------------------------------------------------------------------
# main frontend which proxys to the backends
#---------------------------------------------------------------------

frontend api
    bind <HAProxy Host IP>:6443
    default_backend controlplaneapi

frontend internalapi
    bind <HAProxy Host IP>:22623
    default_backend controlplaneapiinternal

frontend secure
    bind <HAProxy Host IP>:443
    default_backend secure

frontend insecure
    bind <HAProxy Host IP>:80
    default_backend insecure

#---------------------------------------------------------------------
# static backend
#---------------------------------------------------------------------

backend controlplaneapi
    balance source
    server bootstrap <BOOTSTRAP IP>:6443 check     
    server master0 <MASTER 0 IP>:6443 check
    server master1 <MASTER 1 IP>:6443 check
    server master2 <MASTER 2 IP>:6443 check

backend controlplaneapiinternal
    balance source
    server bootstrap <BOOTSTRAP IP>:22623 check     
    server master0 <MASTER 0 IP>:22623 check
    server master1 <MASTER 1 IP>:22623 check
    server master2 <MASTER 2 IP>:22623 check

backend secure
    balance source
    server compute0 <WORKER 0 IP>:443 check
    server compute1 <WORKER 1 IP>:443 check
    server compute2 <WORKER 2 IP>:443 check

backend insecure
    balance source
    server worker0 <WORKER 0 IP>:80 check
    server worker1 <WORKER 1 IP>:80 check
    server worker2 <WORKER 2 IP>:80 check
```

Start and (optionally, enable) the systemd daemon.

```sh 
# If you want it enabled
sudo systemctl enable --now haproxy.service
# If you want to start it manually every time
sudo systemctl start haproxy.service
```

Ensure it's running by checking the systemd journal:

```sh 
# Hit Ctrl+C when done following the logs.
sudo journalctl -f -u haproxy.service
```