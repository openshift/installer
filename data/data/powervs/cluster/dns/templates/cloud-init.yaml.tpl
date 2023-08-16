
#cloud-config
packages:
  - bind
  - bind-utils
%{ if is_proxy ~}
  - httpd
  - mod_ssl
  - squid
%{ endif ~}
write_files:
%{ if is_proxy ~}
- path: /etc/httpd/conf.d/proxy.conf
  content: |
    SSLProxyEngine on
    ProxyPass / https://s3.direct.${vpc_region}.cloud-object-storage.appdomain.cloud/
%{ endif ~}
- path: /tmp/named-conf-edit.sed
  permissions: '0640'
  content: |
    /^\s*listen-on port 53 /s/127\.0\.0\.1/127\.0\.0\.1; MYIP/
    /^\s*allow-query /s/localhost/any/
    /^\s*dnssec-validation /s/ yes/ no/
    /^\s*type hint;/s/ hint/ forward/
    /^\s*file\s"named.ca";/d
    /^\s*type forward/a \\tforward only;\n\tforwarders { 161.26.0.7; 161.26.0.8; };
runcmd:
  - export MYIP=`hostname -I`; sed -i.bak "s/MYIP/$MYIP/" /tmp/named-conf-edit.sed
  - sed -i.orig -f /tmp/named-conf-edit.sed /etc/named.conf
  - systemctl enable named.service
  - systemctl start named.service
%{ if is_proxy ~}
  - echo 'acl ibm_endpoints dstdomain .cloud.ibm.com s3.${vpc_region}.cloud-object-storage.appdomain.cloud' > /tmp/squid.conf
  - echo 'http_access deny !ibm_endpoints' >> /tmp/squid.conf
  - cat /etc/squid/squid.conf >> /tmp/squid.conf
  - mv -f /tmp/squid.conf /etc/squid/squid.conf
  - service httpd start
  - service squid start
%{ endif ~}
