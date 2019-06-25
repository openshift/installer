### Proxy Testing

This will create an extremely basic configuration of squid to support
the testing of authenticated proxy with `openshift-install`.

NOTE: Make sure TCP/3128 is open


- Create directories and configuration files
```
mkdir -p /srv/squid/{etc,cache}
htpasswd -c /srv/squid/etc/passwords <username>

cat << EOF > /srv/squid/etc/squid.conf
auth_param basic program /usr/lib/squid3/basic_ncsa_auth /etc/squid/passwords
auth_param basic realm proxy
acl authenticated proxy_auth REQUIRED
http_access allow authenticated
http_port 3128
cache_dir ufs /var/spool/squid 100 16 256
coredump_dir /var/spool/squid
EOF

chcon -Rt svirt_sandbox_file_t /srv/squid/
```

- Start container
```
URL=docker.io/datadog/squid:latest
SQUID_CACHE_PATH=/srv/squid/cache
SQUID_ETC_PATH=/srv/squid/etc

podman pull ${URL}
podman rm -f squid

podman run --name squid -d -p 3128:3128 \
        --volume ${SQUID_CACHE_PATH}:/var/spool/squid:Z \
        --volume ${SQUID_ETC_PATH}:/etc/squid:Z \
        ${URL}
```

- install-config.yaml snipit

```yaml
---
apiVersion: v1
baseDomain: devcluster.openshift.com
proxy:
  httpsProxy:  "http://username:password@proxy:port"
  httpProxy: "http://username:password@proxy:port"
```
