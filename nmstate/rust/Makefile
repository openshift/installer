CLIB_HEADER_TEMPLATE=vendor/nmstate-clib/nmstate.h.in
CLIB_HEADER=include/nmstate.h
NMSTATE_VERSION=$(shell cat VERSION)
VERSION_MAJOR=$(shell echo $(NMSTATE_VERSION) | cut -f1 -d.)
VERSION_MINOR=$(shell echo $(NMSTATE_VERSION) | cut -f2 -d.)
VERSION_MICRO=$(shell echo $(NMSTATE_VERSION) | cut -f3 -d.)

.PHONY: all
all: lib/libnmstate.a $(CLIB_HEADER)

lib/libnmstate.a:
	cd build && cargo build
	mkdir -p lib
	cp build/target/debug/libnmstate.a lib/libnmstate.a

$(CLIB_HEADER):
	mkdir include
	sed -e 's/@_VERSION_MAJOR@/$(VERSION_MAJOR)/' \
	-e 's/@_VERSION_MINOR@/$(VERSION_MINOR)/' \
	-e 's/@_VERSION_MICRO@/$(VERSION_MICRO)/' \
		$(CLIB_HEADER_TEMPLATE) > $(CLIB_HEADER)

.PHONY: vendor-sync
vendor-sync:
	cargo vendor ./vendor

.PHONY: update-build-cargo-from-vendor
update-build-cargo-from-vendor:
	sed -e 's/cdylib/staticlib/' \
		-e 's#nmstate = .*#nmstate = { path = "../vendor/nmstate" }#' \
		-e 's#path = "lib.rs"#path = "../vendor/nmstate-clib/lib.rs"#' \
		vendor/nmstate-clib/Cargo.toml > build/Cargo.toml

.PHONY: clean
clean:
	rm -fr lib include build/target
