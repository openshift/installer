.PHONY: all deps runtime

VERSION = latest
RUBY_VERSION = 2.4.1
RBENV_VERSION = 1.1.1
RUBY_BUILD_VERSION = 20170726

RUNTIME_IMAGE = smokes-runtime:$(VERSION)
BUILDER_IMAGE = smokes-builder:$(VERSION)
RBENV_URL = https://github.com/rbenv/rbenv/archive/v$(RBENV_VERSION).zip
RUBY_BUILD_URL = https://github.com/rbenv/ruby-build/archive/v$(RUBY_BUILD_VERSION).zip

TMPDIR = /tmp
BUILD_ROOT = $(PWD)

all: runtime

clean:
	rm -rf $(TMPDIR)/rbenv.zip $(TMPDIR)/ruby-build.zip rbenv
	docker rmi $(RUNTIME_IMAGE) $(BUILDER_IMAGE)

deps: $(BUILD_ROOT)/rbenv $(BUILD_ROOT)/rbenv/plugins/ruby-build

$(BUILD_ROOT)/rbenv:
	curl -#Lo "$(TMPDIR)/rbenv.zip" $(RBENV_URL)
	unzip "$(TMPDIR)/rbenv.zip"
	mv rbenv-$(RBENV_VERSION) rbenv

$(BUILD_ROOT)/rbenv/plugins/ruby-build: rbenv
	curl -#Lo "$(TMPDIR)/ruby-build.zip" $(RUBY_BUILD_URL)
	mkdir -p ./rbenv/plugins
	unzip -d ./rbenv/plugins/ "$(TMPDIR)/ruby-build.zip"
	mv ./rbenv/plugins/ruby-build-$(RUBY_BUILD_VERSION) ./rbenv/plugins/ruby-build

builder: deps
	docker build -f Dockerfile.builder -t $(BUILDER_IMAGE) .

build-ruby: builder
	docker run -t --rm -v $(PWD)/rbenv:/home/rspec/.rbenv $(BUILDER_IMAGE) ". ~/.bashrc && rbenv install -s $(RUBY_VERSION)"

runtime: build-ruby
	docker build -f Dockerfile.runtime -t $(RUNTIME_IMAGE) .
