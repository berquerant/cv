GOMOD = go mod
GOBUILD = go build -trimpath -v
GOTEST = go test -v -cover -race

ROOT = $(shell git rev-parse --show-toplevel)
BIN = dist/cv
CMD = ./cmd/cv

.PHONY: $(BIN)
$(BIN):
	$(GOBUILD) -o $@ $(CMD)

.PHONY: test
test:
	$(GOTEST) ./...

.PHONY: init
init:
	$(GOMOD) tidy

.PHONY: generate
generate: clean-generated
	go generate ./...

.PHONY: clean-generated
clean-generated:
	find . -name "*_generated.go" -type f -delete

.PHONY: vuln
vuln:
	go run golang.org/x/vuln/cmd/govulncheck ./...
