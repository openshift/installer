.PHONY: all
all: build

.PHONY: build
build: node_modules ../assets/frontend/scripts/app-bundle.js

../assets/frontend/scripts/app-bundle.js: node_modules *.jsx *.js $(shell find components -type f -name '*.jsx')
	yarn run build

node_modules: package.json yarn.lock
	yarn install --quiet > /dev/null

.PHONY: test
test:
	yarn test

.PHONY: lint
lint:
	yarn run lint

.PHONY: dev
dev:
	yarn run dev

.PHONY: clean
clean:
	rm -f ../assets/frontend/scripts/app-bundle.js
	rm -fr node_modules
