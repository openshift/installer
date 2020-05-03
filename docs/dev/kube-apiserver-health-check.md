## Graceful Termination
`kube-apiserver` in OpenShift is fronted by an external and an internal Load Balancer. This document serves as a
guideline on how to properly configure the health check probe of the load balancers so that when a `kube-apiserver`
instance restarts we can ensure:
- The load balancers detect it and takes it out of service in time. No new request should be forwarded to the 
  `kube-apiserver` instance when it has stopped listening.
- Existing connections are not cut off hard, they are allowed to complete gracefully.

## Load Balancer Health Check Probe
`kube-apiserver` provides graceful termination support via the `/readyz` health check endpoint. When `/readyz` reports
`HTTP Status 200 OK` it indicates that the apiserver is ready to serve request(s).

Now let's walk through the events (in chronological order) that unfold when a `kube-apiserver` instance restarts:
* E1: `T+0s`: `kube-apiserver` receives a TERM signal.
* E2: `T+0s`: `/readyz` starts reporting `failure` to signal to the load balancers that a shut down is in progress.
  * The apiserver will continue to accept new request(s).
  * The apiserver waits for certain amount of time (configurable by `shutdown-delay-duration`) before it stops accepting new request(s).
* E3: `T+30s`: `kube-apiserver` (the http server) stops listening:
  * `/healthz` turns red.
  * Default TCP health check probe on port `6443` will fail.
  * Any new request forwarded to it will fail, most likely with a `connection refused` error or `GOAWAY` for http/2.
  * Existing request(s) in-flight are not cut off but are given up to `60s` to complete gracefully.
* E4: `T+30s+60s`: Any existing request(s) that are still in flight are terminated with an error `reason: Timeout message: request did not complete within 60s`.
* E5: `T+30s+60s`: The apiserver process exits.

Please note that after `E3` takes place, there is a scenario where all existing requests in-flight can gracefully complete
before the `60s` timeout. In such a case no request is forcefully terminated (`E4` does not transpire) and `E5` 
can come about well before `T+30s+60s`. 

An important note to consider is that today in OpenShift the time difference between `E3` and `E2` is `70s`. This is known as
`shutdown-delay-duration` and is configurable by the devs only. This is not a knob we allow the end user to tweak. 
```
$ kubectl -n openshift-kube-apiserver get cm config -o json | jq -r '.data."config.yaml"' |
  jq '.apiServerArguments."shutdown-delay-duration"'
[
  "70s"
]
```
In future we will reduce `shutdown-delay-duration` to `30s`. So in this document we will continue with `E3 - E2` is `30s`.

Given the above, we can infer the following:
* The load balancers should use `/readyz` endpoint for `kube-apiserver` health check probe. It must NOT use `/healthz` or
default TCP port probe.
* The time taken by a load balancer (let's say `t` seconds) to deem a `kube-apiserver` instance unhealthy and take it
out of service should not bleed into `E3`. So `E2 + t < E3` must be true so that no new request is forwarded to the 
instance at `E3` or later. 
* In the worst case, a load balancer should take at most `30s` (since `E2` triggers) to take the `kube-apiserver` 
instance out of service.

Below is the health check configuration used by `aws` currently.

```
protocol: HTTPS
path: /readyz
port: 6443
unhealthy threshold: 2
timeout: 10s
interval: 10s
```

Based on aws documentation, the following is true of the ec2 load balancer health check probes:
* Each health check request is independent and lasts the entire interval.
* The time it takes for the instance to respond does not affect the interval for the next health check.

Now let's verify that with the above configuration in effect, a load balancer takes at most `30s` (in the worst case) to
deem a particular `kube-apiserver` instance unhealthy and take it out of service. With that in mind we will plot the 
timeline of the health check probes accordingly. There are three probes `P1`, `P2` and `P3` involved in this worst 
case scenario:
* E1: T+0s:  `P1` kicks off and it immediately gets a `200` response from `/readyz`.
* E2: T+0s:  `/readyz` starts reporting red, immediately after `E1`.
* E3: T+10s: `P2` kicks off.
* E4: T+20s: `P2` times out (we assume the worst case here).
* E5: T+20s: `P3` kicks off (each health check is independent and will be kicked off at every interval).
* E6: T+30s: `P3` times out (we assume the worst case here)
* E7: T+30s: `unhealthy threshold` is satisfied and the load balancer takes the unhealthy `kube-apiserver` instance out 
  of service.

Based on the worst case scenario above we have verified that with the above configuration aws load balancer will take at
most `30s` to detect an unhealthy `kube-apiserver` instance and take it out of service.

If you are working with a different platform please take into consideration relevant health check probe specifics if any
and ensure that the worst case time to detect an unhealthy `kube-apiserver` instance is at most `30s` as explained in
this document.
