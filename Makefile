BIN = bin/ndql
CMD = ./cmd/ndql
THIRD_PARTY_LICENSES = NOTICE
TOOL = ./tools/run.sh

#
# Build
#

.PHONY: $(BIN)
$(BIN):
	./bin/build.sh -o $@ $(CMD)

#
# Test
#

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: test-tree
test-tree:
	go test -race -cover ./pkg/tree/...

#
# Lint
#

.PHONY: lint
lint: check-licenses vet golangci-lint
# lint: check-licenses vet vuln golangci-lint

.PHONY: vuln
vuln:
	"$(TOOL)" govulncheck ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: golangci-lint
golangci-lint:
	"$(TOOL)" golangci-lint config verify -v
	"$(TOOL)" golangci-lint run

.PHONY: check-licenses-diff
check-licenses-diff: $(THIRD_PARTY_LICENSES)
	git diff --exit-code $(THIRD_PARTY_LICENSES)

.PHONY: check-licenses
check-licenses: check-licenses-diff
	./hack/license.sh check

#
# Code generation
#

.PHONY: $(THIRD_PARTY_LICENSES)
$(THIRD_PARTY_LICENSES):
	./hack/license.sh report > $@

.PHONY: generate
generate:
	go generate ./...

.PHONY: clean-generated
clean-generated:
	find . -name "*_generated.go" -type f -delete

#
# etc
#
.PHONY: clean-tools
clean-tools:
	rm -f bin/kubectl bin/kind
