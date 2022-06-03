#Vendoring of NMState

This directory vendors:
* The [NMState rust C library](https://github.com/nmstate/nmstate/tree/base/rust/src/clib).
* The [Golang bindings of the NMState rust C library](https://github.com/nmstate/nmstate/tree/base/rust/go/nmstate)

## Building

There is a Makefile with `all` and `clean` targets

## Updating to a new version

To update to a new version one needs to:
* Update the VERSION file to match the VERSION in NMState repo in the commit you are syncing to.
* Update the `./Cargo.toml` to update the NMState clib dependenc to the commit/release you choose.
* Run `make rust-vendor` so that Cargo updates the `rust/vendor` directory
* Run `make update-build-cargo-from-vendor` so that the build directory gets an updated Cargo.toml from the vendored nmstate-clib.
