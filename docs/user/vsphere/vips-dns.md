# IP Adresses

An installer-provisioned vSphere installation requires two static IP addresses:

* **API** - used to access the cluster API.
* **Ingress** - used for cluster ingress traffic.

A virtual IP address for each of these should be specified in the [install configuration](install.md#create-configuration).

# DNS Records

DNS records must be created for the two IP addresses in whichever DNS server is appropriate for the environment.
The records should have the following values:

| Name                                  | Value       |
| -                                     |  -          |
| `api.<cluster-name>.<base-domain>`    | API VIP     |
| `*.apps.<cluster-name>.<base-domain>` | Ingress VIP |

Note that `cluster-name` and `base-domain` are variables custom to an installation and
must correspond to the values specified in the [install configuration](install.md#create-configuration).