# Using a Custom External Load Balancer - post deployment

You can shift api/ingress traffic from the default self-hosted load balancer to a load balancer that you provide. To do so, the instance that it runs from must be able to access every machine in your cluster. You might ensure this access by creating the instance on a subnet that is within your cluster's network.

## External Facing OpenShift Services

Add the following external facing services to your new load balancer:

- The master nodes serve the OpenShift API on port 6443 using TCP.
- The apps hosted on the master and worker nodes are served on ports 80, and 443. They are both served using TCP.

## HAProxy Example Load Balancer Config

The following `HAProxy` config file demonstrates a basic configuration for an external load balancer:

```haproxy

defaults
    mode                    tcp
    log                     global
    timeout connect         30s
    timeout client          1m
    timeout server          1m
frontend <cluster-name>-api-6443
    bind :::6443 v4v6
    default_backend api
frontend <cluster-name>-apps-80
    bind :::80  v4v6
    default_backend ingress
frontend <cluster-name>-apps-443
    bind :::443  v4v6
    default_backend ingress-sec
backend api
    option  httpchk GET /readyz HTTP/1.0
    option  log-health-checks
    balance roundrobin
    server master-0 <master0-IP>:6443 check check-ssl inter 1s fall 2 rise 3 verify none
    server master-1 <master1-IP>:6443 check check-ssl inter 1s fall 2 rise 3 verify none
    server master-2 <master2-IP>:6443 check check-ssl inter 1s fall 2 rise 3 verify none
backend ingress
    option  httpchk GET /healthz/ready  HTTP/1.0
    option  log-health-checks
    balance roundrobin
    server master-0 <master0-IP>:80 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
    server master-1 <master1-IP>:80 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
    server master-2 <master2-IP>:80 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
    server worker-0 <worker0-IP>:80 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
    server worker-1 <worker1-IP>:80 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
backend ingress-sec
    option  httpchk GET /healthz/ready  HTTP/1.0
    option  log-health-checks
    balance roundrobin
    server master-0 <master0-IP>:443 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
    server master-1 <master1-IP>:443 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
    server master-2 <master2-IP>:443 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
    server worker-0 <worker0-IP>:443 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
    server worker-1 <worker1-IP>:443 check check-ssl port 1936 inter 1s fall 2 rise 3 verify none
```

## DNS Lookups

To ensure that your API and apps are accessible through your load balancer, create or update your DNS entries for those endpoints. To use your new load balancing service for external traffic, make sure the IP address for these DNS entries is the IP address your load balancer is reachable at.

```dns
<load balancer ip> api.<cluster-name>.<base domain>
<load balancer ip> apps.<cluster-name>.base domain>
```

## Verifying that the API is Reachable

One good way to test whether or not you can reach the API is to run the `oc` command. If you can't do that easily, you can use this curl command:

```sh
curl https://api.<cluster-name>.<base domain>:6443/version --insecure
```

Result:

```json
{
  "major": "1",
  "minor": "20",
  "gitVersion": "v1.20.0+ba45583",
  "gitCommit": "ba455830ecb91ff61bb61ca4f70b6f3f4a5e3796",
  "gitTreeState": "clean",
  "buildDate": "2021-02-05T22:18:43Z",
  "goVersion": "go1.15.5",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```

Note: The versions in the sample output may differ from your own. As long as you get a JSON payload response, the API is accessible.

## Verifying that Apps Reachable

The simplest way to verify that apps are reachable is to open the OpenShift console in a web browser. If you don't have access to a web browser, query the console with the following curl command:

```sh
curl http://console-openshift-console.apps.<cluster-name>.<base domain> -I -L --insecure
```


Result:

```http
HTTP/1.1 302 Found
Cache-Control: no-cache
Content-length: 0
Location: https://console-openshift-console.apps.<cluster-name>.<base domain>/

HTTP/1.1 200 OK
Referrer-Policy: strict-origin-when-cross-origin
Set-Cookie: csrf-token=ZPwr8qTwPBh/NQjoENlDWxmACNEsLl1PYrQyyX87wnIm5AnBrwv3dEqpZwClwpN4nWlGp2ufBh7KbM0ycwLQpQ==; Path=/; Secure
X-Content-Type-Options: nosniff
X-Dns-Prefetch-Control: off
X-Frame-Options: DENY
X-Xss-Protection: 1; mode=block
Date: Wed, 17 Mar 2021 09:18:06 GMT
Content-Type: text/html; charset=utf-8
Set-Cookie: 1e2670d92730b515ce3a1bb65da45062=1115fecfa3e981219adba594404c9b69; path=/; HttpOnly; Secure; SameSite=None
Cache-control: private
```
