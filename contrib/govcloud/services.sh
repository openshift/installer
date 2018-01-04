#!/bin/bash
# TODO: Convert to Ignition + Systemd

# OpenVPN https://github.com/kylemanna/docker-openvpn
OVPN_DATA=basic-data
CLIENT=CI
IMG=kylemanna/openvpn:2.4
SERV_IP=${ip}
CLIENT_DIR=/home/core/vpn-config
mkdir -p $CLIENT_DIR
docker run -v $OVPN_DATA:/etc/openvpn --rm $IMG ovpn_genconfig -u udp://$SERV_IP
docker run -v $OVPN_DATA:/etc/openvpn --rm -e "EASYRSA_BATCH=1" -e "EASYRSA_REQ_CN=CI Test CA" $IMG ovpn_initpki nopass
docker run -v $OVPN_DATA:/etc/openvpn --rm $IMG easyrsa build-client-full $CLIENT nopass
docker run -v $OVPN_DATA:/etc/openvpn --rm $IMG ovpn_getclient $CLIENT | tee $CLIENT_DIR/config.ovpn
docker run --restart always --name "ovpn-test" -v $OVPN_DATA:/etc/openvpn -d -p 1194:1194/udp --cap-add=NET_ADMIN $IMG

# Web server for OpenVPN client Config
NGINX_DIR=/home/core/nginx-config
mkdir -p $NGINX_DIR
echo "
server {
    listen       80;
    server_name  localhost;
    auth_basic           'Administrators Area';
    auth_basic_user_file /etc/nginx/conf.d/htpasswd;
    location / {
        root   /usr/share/nginx/html;
        index  config.ovpn;
    }
}
" > $NGINX_DIR/default.conf

printf "${username}:$(openssl passwd -crypt ${password})\n" >> $NGINX_DIR/htpasswd
docker run -d --restart always \
--name nginx-app \
--net host \
-v $CLIENT_DIR:/usr/share/nginx/html/ \
-v $NGINX_DIR:/etc/nginx/conf.d/ \
nginx:1.13.7-alpine

# PowerDNS
docker run -d --name=mysql --net=host -e MYSQL_ROOT_PASSWORD=powerdns mysql
sleep 5
docker run --net=host \
--name pdns-master -d \
--restart always \
-e PDNS_RECURSOR=10.0.0.2 \
-e PDNS_SOA=10.0.0.2 \
-e PDNS_ALLOW_AXFR_IPS=127.0.0.1 \
-e PDNS_DISTRIBUTOR_THREADS=3 \
-e PDNS_CACHE_TTL=20 \
-e PDNS_RECURSIVE_CACHE_TTL=10 \
-e DB_ENV_MYSQL_ROOT_PASSWORD=powerdns \
-e MYSQL_HOST=127.0.0.1 \
-e MYSQL_PORT="3306" \
-e PDNS_ZONE=${dns_zone} \
quay.io/nicholas_lane/pdns:4.0-1
sleep 20
docker exec pdns-master pdnsutil create-zone ${dns_zone}

# Nat Gateway routing
echo 1 > /proc/sys/net/ipv4/ip_forward
iptables -t nat -A POSTROUTING -s ${private_cidr} -j MASQUERADE

# Disable automatic updates
systemctl stop update-engine
