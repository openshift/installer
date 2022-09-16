#Introduction

The installer gathers all user configuration and consolidates it in the install-config.yaml file. It then uses this configuration to generate a slew of manifests. The Installer is responsible for bootstrapping the cluster. The bootstrap process starts several OpenShift operators while providing configuration to them via these manifests. More information about the Bootstrap process can be found [here](https://github.com/openshift/installer/blob/master/docs/dev/bootstrap_services.md). The role of the Installer comes to an end after this point. These operators in turn use these manifests to generate the resulting cluster. It then becomes the responsibilty of these Operators to manage changes within the cluster on day-2.

While designing new features for the Installer, it is important to adhere to some design principles so that we achieve the following outcomes in no specific order of their importance:
1. the design of different Installer features are consistent
2. the Installer is able to provide the best user experince
3. communicating the Installer's design goals and principles while colloborating with other teams

# Design Documentation

This is an attempt to document the OpenShift Installer's existing design principles and articulate design goals to guide new feature development and in some cases inform bug fixes. This is expected to be a living document which is continually updated when we gain clarity around some of these principles.

## Design Goals

### Install methods

The OpenShift Installer currently supports two installation types, UPI (User Provisioned Infrastructure) and IPI (Installer Provisionined Infrastructure). IPI is an opinionated method of installing OpenShift. The IPI installation method tries to automate a large portion of the OpenShift Install with the least amount of intervention from the user. While it tries to take care of supporting the most common install scenarios, it cannot be expected to support every install variation that a user wants to use for their cloud.

This is where UPI offers greater flexibility required for highly specialized installs. In this install method, the user is expected to provide several Infrastructure components needed for the install on their own. For this reason, documentation surrounding the UPI flow are especially important.

The OpenShift Installer has the responsibility to validate configuration provided to it. There is a validation step that is common to both the UPI and IPI install methods and additional validation is performed specially for the IPI install method.

### Cluster Configuration

A design goal for the Installer is to provide a user-friendly interface for collecting configuration details. One of the ways that the Installer accomplishes this goal is by minimizing the configuration details required to create an OpenShift cluster on the platform and with capabilities requested by the user. The Installer and OpenShift as a whole, have to strive to define sane defaults for configuration whenever possible and allow the user to override these defaults.

Any information that can be learned from the install environment should not be requested as an input from the user. If another OpenShift entity is capable of learning this information, the Installer would defer this responsibility to this entity. This prevents the Installer from becoming a monolith and allows this component to keep its treatment of this environmental config consistent during the time of cluster bringup and day-2.

Another important goal while adding new features to OpenShift is for a cluster upgraded to a version supporting this feature to behave identically to a freshly Installed cluster of that version. When adding configuration for enabling a new functionality within the OpenShift cluster, in addition to designing for a new Installation, consider how a Cluster upgraded to the version supporting this feature would behave. The Operator providing this new functionality should allow updating its API to enable this feature on day-2. Also, it may be helpful to design this feature to be disabled by default and enabled explicitly by changing this API. 
