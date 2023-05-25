# This Makefile is handling the go build, install, test and clean. It is needed
# to avoid isues with the GOROOT and GOPATH.
TMPDIR := $(shell mktemp -d)
TOP_SRC_DIR:=$(shell dirname $(realpath $(firstword \
	$(MAKEFILE_LIST))))/../../../

include $(TOP_SRC_DIR)/Makefile.inc

RUST_DEBUG_BIN_DIR=$(TOP_SRC_DIR)/target/debug
CLIB_SO_DEV_DEBUG=$(RUST_DEBUG_BIN_DIR)/$(CLIB_SO_DEV)
CLIB_HEADER=nmstate.h
CLIB_HEADER_IN=$(TOP_SRC_DIR)src/clib/$(CLIB_HEADER).in

all: build check

build:
	go build;

$(TMPDIR)/$(CLIB_HEADER): $(CLIB_HEADER_IN)
	install -p -v -D -m644 $(CLIB_HEADER_IN) \
		$(TMPDIR)/$(CLIB_HEADER)
	sed -i -e 's/@_VERSION_MAJOR@/$(CLIB_VERSION_MAJOR)/' \
		$(TMPDIR)/$(CLIB_HEADER)
	sed -i -e 's/@_VERSION_MINOR@/$(CLIB_VERSION_MINOR)/' \
		$(TMPDIR)/$(CLIB_HEADER)
	sed -i -e 's/@_VERSION_MICRO@/$(CLIB_VERSION_MICRO)/' \
		$(TMPDIR)/$(CLIB_HEADER)

$(CLIB_SO_DEV_DEBUG):
	cd $(TOP_SRC_DIR) && cargo build -p nmstate-clib

$(TMPDIR)/$(CLIB_SO_FULL): $(CLIB_SO_DEV_DEBUG)
	install -p -v -D -m755 $(CLIB_SO_DEV_DEBUG) \
		$(TMPDIR)/$(CLIB_SO_FULL)
	ln -sfv $(CLIB_SO_FULL) $(TMPDIR)/$(CLIB_SO_DEV)

check: $(TMPDIR)/$(CLIB_SO_FULL) $(TMPDIR)/$(CLIB_HEADER)
	LD_LIBRARY_PATH=$(TMPDIR) \
	CGO_CFLAGS="-I$(TMPDIR)" \
	CGO_LDFLAGS="-L$(TMPDIR)" \
	go test $(WHAT)
	rm -rf $(TMPDIR)
clean:
	go clean;
