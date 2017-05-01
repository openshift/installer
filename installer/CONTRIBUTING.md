
# Contributing

## Adding Operators

* Add static operator manifests to `installer/assets/candidate/APP`
* Demonstrate that components start quickly and reliably, without an operator
* Provide documentation
    * What are the success criteria? How do we know its installed and working?
    * How to support it? What are the failure modes?
    * Is host state used or modified?
    * How to fully uninstall?
    * Cleanup and migration strategies?
* Move static operator manifests to `installer/assets/APP`
    * Operator pods will be paused by default and can be enabled through the console
    * Show that the app runs by itself
    * Show that the app updates reliably by manually unpausing the operator
    * Show that the app upgrades from old versions and migrates correctly
