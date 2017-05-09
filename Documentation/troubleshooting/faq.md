# Tectonic FAQ

## Tectonic release versioning

The version number of a Tectonic release is a string in the format `W.X.Y-tectonic.z`, where `W.X.Y` stands for the Kubernetes version included with the release, and `z` is an incrementing number, starting at 1, for successive Tectonic releases including that Kubernetes version. When the version number of the included Kubernetes release changes, the incrementing value `z` resets to 1.

For example, if the Kubernetes version is 1.6.2, the first Tectonic production release including it is labeled `1.6.2-tectonic.1`. A second release including the same 1.6.2 version of Kubernetes would be `1.6.2-tectonic.2`. When the version number of the included Kubernetes advances to 1.6.2, the associated Tectonic version number would be `1.6.2-tectonic.1`.

## License and pull-secret formats

When copying your license and pull-secret from account.coreos.com, be sure to choose the correct format. The license format should be "Raw Format" and the pull-secret should be in the "dockercfg" format.

## Can I use my license on multiple clusters?

Yes, a single license can be used on multiple clusters, as long as the total limits on that license are not exceeded.

## What happens when I exceed my license limits?

Once your cluster is larger than your license allows, the Console will prompt you to input an updated license. Enter the new license in the dialog box that is shown. The new license limits will take a few minutes to take effect, but applications on the cluster will continue to run.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12 co-m-screenshot">
    <img src="../img/license-exceeded.png">
    <div class="co-m-screenshot-caption">Console showing a cluster that has exceeded capacity</div>
  </div>
</div>

Visit [coreos.com/contact][contact] for an updated license.

## Domain name can't be changed

The domain names selected for Tectonic and Controller DNS names during install cannot be easily changed later. If a cluster's domain name must change, set up a new cluster with the new domain name and migrate cluster work to it.

## Safari browser on OSX

Tectonic Installer may fail to download `assets.zip` file during the install due to a known issue in Safari/Webkit. OSX users are advised to use Firefox or Chrome to use Tectonic Installer on OSX machines.

## How can I recover from a Tectonic Identity misconfiguration?

If you have misconfigured Tectonic Identity, it is possible to be temporarily locked out of your cluster. To correct the misconfiguration, use the credentials in the [assets bundle][assets] generated during installation to revert or correct your [user and RBAC configuration][user-management].

## What Services are Installed with Tectonic?

Tectonic Services are the applications that are installed into your cluster. These include:

| Name | Description |
|------|-------------|
| Tectonic Console   | Web management console for Kubernetes and the services themselves |
| Tectonic Identity  | Centralized user management for services on your cluster |
| Prometheus         | Complete cluster monitoring: instrumentation, collection, and querying. |

---

## Community Support Forum

Make sure to check out the [community support forum](https://github.com/coreos/tectonic-forum/issues) to work through issues, report bugs, identify documentation requirements, or put in feature requests.

*If you have any unanswered questions about Tectonic Installer visit [coreos.com/contact][contact]*.

*If you have any unanswered questions about Tectonic visit [coreos.com/contact][contact]*.


[assets]: ../admin/assets-zip.md
[user-management]: ../admin/user-management.md
[contact]: https://coreos.com/contact/
[sign-up]: https://account.coreos.com/signup/summary/tectonic-2016-12
