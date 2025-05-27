# Make: Makefile prompt for user input #make #shell

[View original Gist on GitHub](https://gist.github.com/Integralist/fdf8374a593eeef8db8fed0fb868d5eb)

## Makefile

```makefile
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=fastly
FULL_PKG_NAME=github.com/fastly/terraform-provider-fastly
VERSION_PLACEHOLDER=version.ProviderVersion
VERSION=$(shell git describe --tags --always)
VERSION_SHORT=$(shell git describe --tags --always --abbrev=0)
DOCS_PROVIDER_VERSION=$(subst v,,$(VERSION_SHORT))

GOHOSTOS ?= $(shell go env GOHOSTOS || echo unknown)
GOHOSTARCH ?= $(shell go env GOHOSTARCH || echo unknown)

TEST_PARALLELISM?=4

default: build

build: clean
	go build -o bin/terraform-provider-$(PKG_NAME)_$(VERSION) -ldflags="-X $(FULL_PKG_NAME)/$(VERSION_PLACEHOLDER)=$(VERSION)"
	@sh -c "'$(CURDIR)/scripts/generate-dev-overrides.sh'"

test:
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=$(TEST_PARALLELISM)

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -parallel=$(TEST_PARALLELISM) -timeout 360m -ldflags="-X=$(FULL_PKG_NAME)/$(VERSION_PLACEHOLDER)=acc"

# WARNING: This target will delete infrastructure.
clean_test:
	@printf 'WARNING: This will delete infrastructure. Continue? (y/n) '; \
	read answer; \
	if echo "$$answer" | grep -iq '^y'; then \
	  SILENCE=true make sweep || true; \
		fastly service list --token $$FASTLY_API_KEY | grep -E '^tf\-' | awk '{print $$2}' | xargs -I % fastly service delete --token $$FASTLY_API_KEY -f -s %; \
		TEST_PARALLELISM=8 make testacc; \
	fi

sweep:
	@if [ "$(SILENCE)" != "true" ]; then \
		echo "WARNING: This will destroy infrastructure. Use only in development accounts."; \
	fi
	go test ./fastly -v -sweep=ALL $(SWEEPARGS) -timeout 30m || true

clean:
	rm -rf ./bin

.PHONY: build test testacc sweep clean
```

