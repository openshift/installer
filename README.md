[![Build Status](https://travis-ci.org/openshift/installer.svg?branch=master)](https://travis-ci.org/openshift/installer)

# Tectonic Installer

The CoreOS and OpenShift teams are now working together to integrate Tectonic and Open Shift into a converged platform which will be developed in https://github.com/openshift/installer. We'll consider all feature requests for the new converged platform, but will not be adding new features to _this_ repository.

In the meantime, current Tectonic customers will continue to receive support and updates. Any such bugfixes will take place on the [track-1](https://github.com/coreos/tectonic-installer/tree/track-1) branch.

See the CoreOS blog for any additional details:
https://coreos.com/blog/coreos-tech-to-combine-with-red-hat-openshift

*Note*: The master branch of the project reflects a work-in-progress design approach works only on AWS and Libvirt. In order to deploy Tectonic to other platforms, e.g. Azure, bare metal, OpenStack, etc, please checkout the [track-1](https://github.com/coreos/tectonic-installer/tree/track-1) branch of this project, which maintains support for the previous architecture and more platforms.
