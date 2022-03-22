## Bootstrap Gather via API Endpoint ##

Many user have security restrictions that prevent SSH access to VMs. These restrictions impede the ability to diagnose
bootstrap failures, as the installer relies on SSH to the bootstrap VM in order to gather logs. For a typical
installation, the installer does have ingress to the bootstrap via the API endpoint. The installer can leverage that
endpoint to gather the bootstrap logs instead of using SSH.

### High-Level Design ###

The installer will make a GET request to https://api.<cluster-domain>:6443/bootstrap in order to fetch the bootstrap
logs. This assumes that the installer already had a sticky connection to the bootstrap backend of the API load balancer.

The bootstrap will intercept incoming connections to port 6443 and forward them internally to port 6444. The bootstrap
will run a server listening on port 6444. If the path of the incoming request is /bootstrap, then the server will gather
bootstrap logs and respond to the request with the gathered logs. Otherwise, the server will forward the incoming
request to the kubernetes api-server running locally on port 6443.

### Details ###

#### TCP Termination ####

Since the intercepting server is evaluating the path of the URL, the server needs to terminate the TCP connection. The
intercepting server needs to serve the same certs that the local kubernetes api-server serves (or at least certs 
matching the certs that the installer uses to connect to the API endpoint). When forwarding the request to the local
kubernetes api-server, the intercepting server needs to the certs that would be used by the bootstrap for local queries
to the kubernetes api-server. All of these certs should be present on the bootstrap VM already in /opt/openshift/tls/.

#### Gathering Logs from Container ####

If we use a container for the intercepting server, then we need a way to communicate back to the bootstrap host to
actually gather the logs. A couple options here are (1) using SSH locally or (2) running commands through a named pipe
back to the host.

#### SSH to the Masters ####

Part of the bootstrap log gather process gathering logs from the masters. When using SSH to connect to the bootstrap,
the bootstrap is able to connect via SSH to the masters using the private keys on the install machine. If the
installer is instead connecting via the API endpoint to the bootstrap VM, the private key is not available. Here are two
possible solutions for this.
1. Have the installer send the private key in the body of the GET request. This is undesireable since it exposes the
private key outside of the install machine.
2. Use SSH tunneling over HTTPS (using the API endpoint) to connect to the masters via the bootstrap.
3. Use something like https://www.rutschle.net/tech/sslh/README.html to demux HTTPS and SSH over the API endpoint.

If we use (2), then that could replace the part where the container does a local SSH to the host to gather logs.

If we use (3), then that could replace the entire intercepting server.

#### Health Probes ####

The load balancer for the API endpoint makes this design a bit flaky. There are competing demands where we want the
load balancer to stop forwarding kubernetes api-server requests to the bootstrap VM once the api-servers are up on the
masters but we want the load balancer to forward bootstrap-gather requests to the bootstrap VM so long as the bootstrap
VM is running. One solution could be to have the intercepting server forward to the permanent control plane after the
cluster-bootstrap step (which runs the temporary control plane) completes successfully.

#### Intercepting Server Service ####

The current implementation is running the intercepting server as its own service. However, the service is not starting
automatically successfully. If the service is started manually after the bootstrap VM boots, the service does run
successfully. I suspect there may be issues with permissions from the user used to run the service automatically that
are causing problems. In the end, it may be better to start the intercepting server from the bootkube.sh script. This
would make it easier to transition from forwarding to the local kubernetes api-server to forwarding to the permanent
kubernetes api-server after cluster-bootstrap completes.