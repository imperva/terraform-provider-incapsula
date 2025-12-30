TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
HOSTNAME=registry.terraform.io
NAMESPACE=terraform-providers
PKG_NAME=incapsula
BINARY=terraform-provider-${PKG_NAME}
# Whenever bumping provider version, please update the version in incapsula/client.go (line 27) as well.
VERSION=3.35.1

# Mac Intel Chip
OS_ARCH=darwin_amd64
# For Mac M1 Chip
# OS_ARCH=darwin_arm64
# OS_ARCH=linux_amd64

default: install

build: fmtcheck
	export GO111MODULE="on"
	go mod vendor
	go build -o ${BINARY}

build-github: fmtcheck
	go build -o ${BINARY}
	mv ${BINARY} ${GOPATH}/bin

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PKG_NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PKG_NAME}/${VERSION}/${OS_ARCH}

MOCK_SERVER_PORT?=19443
MOCK_SERVER_URL=http://localhost:$(MOCK_SERVER_PORT)

test: fmtcheck check-mock-server
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

check-mock-server:
	@if ! curl -s -o /dev/null -w '' $(MOCK_SERVER_URL)/account 2>/dev/null; then \
		echo ""; \
		echo "ERROR: Mock Imperva API server is not running on $(MOCK_SERVER_URL)"; \
		echo ""; \
		echo "To start the mock server, run in a separate terminal:"; \
		echo ""; \
		echo "  make server"; \
		echo ""; \
		echo "Then set the following environment variables:"; \
		echo ""; \
		echo "  export INCAPSULA_API_ID=mock-api-id"; \
		echo "  export INCAPSULA_API_KEY=mock-api-key"; \
		echo "  export INCAPSULA_BASE_URL=$(MOCK_SERVER_URL)"; \
		echo "  export INCAPSULA_BASE_URL_REV_2=$(MOCK_SERVER_URL)"; \
		echo "  export INCAPSULA_BASE_URL_REV_3=$(MOCK_SERVER_URL)"; \
		echo "  export INCAPSULA_BASE_URL_API=$(MOCK_SERVER_URL)"; \
		echo "  export INCAPSULA_CUSTOM_TEST_DOMAIN=.mock.incaptest.com"; \
		echo ""; \
		exit 1; \
	fi
	@echo "Mock server is running on $(MOCK_SERVER_URL)"

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

clean:
	go clean -cache -modcache -i -r

server:
	@echo "Starting mock Imperva API server..."
	go run ./cmd/mock-server

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test server check-mock-server

