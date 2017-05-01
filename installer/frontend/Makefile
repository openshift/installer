.PHONY: all
all: build

.PHONY: build
build:
	@yarn install --quiet > /dev/null && yarn run build

.PHONY: test
test:
	@yarn test

.PHONY: lint
lint:
	@yarn run lint

.PHONY: dev
dev:
	@yarn run dev
